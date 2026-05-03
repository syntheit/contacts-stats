package config

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	VCFPath         string
	CardDAVURL      string
	CardDAVUser     string
	CardDAVPassword string
	Serve           bool
	Port            string
	Command         string // "update" or "serve"
}

func Load() *Config {
	// Try loading .env file, ignore error if not present
	_ = godotenv.Load()

	cfg := &Config{}

	// Flags
	vcfFlag := flag.String("vcf", "", "Path to the VCF file")
	serveFlag := flag.Bool("serve", false, "Start a web server (deprecated, use 'serve' command)")
	portFlag := flag.String("port", "8080", "Port for the web server")
	
	// Check for subcommand
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "serve":
			cfg.Command = "serve"
			// Parse flags after subcommand if needed, but for now we keep it simple
			// We might want to re-parse flags ignoring the first arg
		case "update":
			cfg.Command = "update"
		default:
			// Fallback to old behavior if flags are present
		}
	}
	
	// Parse global flags for backward compatibility or mixed usage
	flag.Parse()

	// Env vars override or fill in
	if *vcfFlag != "" {
		cfg.VCFPath = *vcfFlag
	} else {
		cfg.VCFPath = os.Getenv("VCF_PATH")
	}

	cfg.CardDAVURL = os.Getenv("CARDDAV_URL")
	cfg.CardDAVUser = os.Getenv("CARDDAV_USER")
	cfg.CardDAVPassword = os.Getenv("CARDDAV_PASSWORD")

	cfg.Port = *portFlag
	if os.Getenv("PORT") != "" {
		cfg.Port = os.Getenv("PORT")
	}

	if *serveFlag {
		cfg.Serve = true
	}

	return cfg
}
