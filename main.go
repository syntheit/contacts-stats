package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"contacts-stats/internal/config"
	"contacts-stats/internal/fetcher"
	"contacts-stats/internal/stats"
)

func main() {
	cfg := config.Load()

	// Default to legacy behavior if no specific command and flags are present
	// or if Serve flag is set
	mode := cfg.Command
	if mode == "" {
		if cfg.Serve {
			mode = "serve_legacy"
		} else {
			mode = "update"
		}
	}

	switch mode {
	case "serve":
		runServer(cfg.Port)
	case "serve_legacy":
		// Run update then serve (legacy behavior)
		runUpdate(cfg)
		runServer(cfg.Port)
	case "update":
		runUpdate(cfg)
	default:
		// Check if the first arg was "serve" or "update" which is handled in config
		if len(os.Args) > 1 && (os.Args[1] == "serve" || os.Args[1] == "update") {
			// already handled
		} else {
			fmt.Println("Usage: contacts-stats [update|serve]")
			os.Exit(1)
		}
	}
}

func runUpdate(cfg *config.Config) {
	var f fetcher.Fetcher

	if cfg.CardDAVURL != "" {
		fmt.Printf("Using CardDAV source: %s\n", cfg.CardDAVURL)
		if cfg.CardDAVUser == "" || cfg.CardDAVPassword == "" {
			log.Fatal("Error: CARDDAV_USER and CARDDAV_PASSWORD are required when CARDDAV_URL is set")
		}
		f = &fetcher.CardDAVFetcher{
			URL:      cfg.CardDAVURL,
			User:     cfg.CardDAVUser,
			Password: cfg.CardDAVPassword,
		}
	} else if cfg.VCFPath != "" {
		fmt.Printf("Using VCF file: %s\n", cfg.VCFPath)
		if _, err := os.Stat(cfg.VCFPath); os.IsNotExist(err) {
			log.Fatalf("Error: VCF file not found at %s", cfg.VCFPath)
		}
		f = &fetcher.FileFetcher{Path: cfg.VCFPath}
	} else {
		log.Fatal("Error: Either VCF_PATH/--vcf or CARDDAV_URL env var must be provided")
	}

	cards, err := f.Fetch()
	if err != nil {
		log.Fatalf("Error fetching contacts: %v", err)
	}

	fmt.Printf("Processed %d cards\n", len(cards))

	s := stats.FromCards(cards)
	statsJSON, _ := json.MarshalIndent(s, "", "  ")
	err = os.WriteFile("stats.json", statsJSON, 0644)
	if err != nil {
		log.Fatalf("Error writing stats.json: %v", err)
	}
	fmt.Println("Statistics compiled into stats.json")
}

func runServer(port string) {
	http.Handle("/", http.FileServer(http.Dir(".")))
	fmt.Printf("Serving at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
