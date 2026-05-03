package fetcher

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/emersion/go-vcard"
)

type CardDAVFetcher struct {
	URL      string
	User     string
	Password string
}

func (f *CardDAVFetcher) Fetch() ([]vcard.Card, error) {
	if f.URL == "" {
		return nil, fmt.Errorf("CardDAV URL is required")
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 1. PROPFIND to list files
	req, err := http.NewRequest("PROPFIND", f.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.SetBasicAuth(f.User, f.Password)
	req.Header.Set("Depth", "1")
	req.Header.Set("Content-Type", "application/xml; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to CardDAV server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMultiStatus && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("CardDAV PROPFIND failed with status: %s", resp.Status)
	}

	var ms multistatus
	if err := xml.NewDecoder(resp.Body).Decode(&ms); err != nil {
		return nil, fmt.Errorf("failed to decode PROPFIND response: %w", err)
	}

	fmt.Printf("Found %d items in CardDAV collection\n", len(ms.Responses))

	var cards []vcard.Card

	// 2. Iterate and GET .vcf files
	for i, r := range ms.Responses {
		// Skip directories or non-vcf files (simplistic check)
		if !strings.HasSuffix(r.Href, ".vcf") {
			continue
		}

		targetURL := f.resolveURL(r.Href)
		fmt.Printf("[%d/%d] Fetching: %s\n", i+1, len(ms.Responses), targetURL)
		
		card, err := f.fetchOne(client, targetURL)
		if err != nil {
			fmt.Printf("Warning: failed to fetch %s: %v\n", targetURL, err)
			continue
		}
		if card != nil {
			cards = append(cards, *card)
		}
	}

	return cards, nil
}

func (f *CardDAVFetcher) resolveURL(href string) string {
	u, err := url.Parse(f.URL)
	if err != nil {
		return href // Should not happen if f.URL is valid
	}
	
	ref, err := url.Parse(href)
	if err != nil {
		return href
	}
	
	return u.ResolveReference(ref).String()
}

func (f *CardDAVFetcher) fetchOne(client *http.Client, urlStr string) (*vcard.Card, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(f.User, f.Password)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %s", resp.Status)
	}

	dec := vcard.NewDecoder(resp.Body)
	card, err := dec.Decode()
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}
	return &card, nil
}

// XML structures
type multistatus struct {
	XMLName   xml.Name   `xml:"multistatus"`
	Responses []response `xml:"response"`
}

type response struct {
	Href string `xml:"href"`
}
