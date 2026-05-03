package stats

import (
	"fmt"
	"strings"

	"github.com/emersion/go-vcard"
	"github.com/nyaruka/phonenumbers"
)

type Stats struct {
	Countries    map[string]int `json:"countries"`
	Subdivisions map[string]int `json:"subdivisions"`
	Detailed     []DetailedStat `json:"detailed"`
	Total        int            `json:"total"`
}

type DetailedStat struct {
	Area  string  `json:"area"`
	Count int     `json:"count"`
	Pct   float64 `json:"pct"`
}

func FromCards(cards []vcard.Card) *Stats {
	stats := &Stats{
		Countries:    make(map[string]int),
		Subdivisions: make(map[string]int),
		Detailed:     []DetailedStat{},
	}

	areaCounts := make(map[string]int)
	totalNumbers := 0

	for _, card := range cards {
		phones := card.Values(vcard.FieldTelephone)
		for _, p := range phones {
			raw := p
			num, err := phonenumbers.Parse(raw, "")
			if err != nil {
				continue
			}

			totalNumbers++
			region := phonenumbers.GetRegionCodeForNumber(num)
			stats.Countries[region]++

			geocoding, _ := phonenumbers.GetGeocodingForNumber(num, "en")
			if geocoding == "" {
				geocoding = region
			}

			// Try to identify subdivision (State/Province)
			sub := identifySubdivision(region, geocoding)
			if sub != "" {
				stats.Subdivisions[sub]++
			}

			// Detailed area info
			areaKey := fmt.Sprintf("%s (%s)", geocoding, region)
			areaCounts[areaKey]++
		}
	}

	stats.Total = totalNumbers
	for area, count := range areaCounts {
		stats.Detailed = append(stats.Detailed, DetailedStat{
			Area:  area,
			Count: count,
			Pct:   float64(count) * 100 / float64(totalNumbers),
		})
	}

	return stats
}

func identifySubdivision(country, geocoding string) string {
	if country == "US" {
		// Geocoding might be "Chicago, IL" or just "Illinois"
		parts := strings.Split(geocoding, ", ")
		if len(parts) > 1 {
			state := parts[len(parts)-1]
			if len(state) == 2 {
				return "US-" + strings.ToUpper(state)
			}
		}
		// Fallback for state names
		states := map[string]string{
			"Alabama": "AL", "Alaska": "AK", "Arizona": "AZ", "Arkansas": "AR", "California": "CA",
			"Colorado": "CO", "Connecticut": "CT", "Delaware": "DE", "Florida": "FL", "Georgia": "GA",
			"Hawaii": "HI", "Idaho": "ID", "Illinois": "IL", "Indiana": "IN", "Iowa": "IA",
			"Kansas": "KS", "Kentucky": "KY", "Louisiana": "LA", "Maine": "ME", "Maryland": "MD",
			"Massachusetts": "MA", "Michigan": "MI", "Minnesota": "MN", "Mississippi": "MS", "Missouri": "MO",
			"Montana": "MT", "Nebraska": "NE", "Nevada": "NV", "New Hampshire": "NH", "New Jersey": "NJ",
			"New Mexico": "NM", "New York": "NY", "North Carolina": "NC", "North Dakota": "ND", "Ohio": "OH",
			"Oklahoma": "OK", "Oregon": "OR", "Pennsylvania": "PA", "Rhode Island": "RI", "South Carolina": "SC",
			"South Dakota": "SD", "Tennessee": "TN", "Texas": "TX", "Utah": "UT", "Vermont": "VT",
			"Virginia": "VA", "Washington": "WA", "West Virginia": "WV", "Wisconsin": "WI", "Wyoming": "WY",
		}
		if code, ok := states[geocoding]; ok {
			return "US-" + code
		}
	}

	if country == "CA" {
		provinces := map[string]string{
			"Alberta": "AB", "British Columbia": "BC", "Manitoba": "MB", "New Brunswick": "NB",
			"Newfoundland and Labrador": "NL", "Nova Scotia": "NS", "Ontario": "ON", "Prince Edward Island": "PE",
			"Quebec": "QC", "Saskatchewan": "SK", "Northwest Territories": "NT", "Nunavut": "NU", "Yukon": "YT",
		}
		if code, ok := provinces[geocoding]; ok {
			return "CA-" + code
		}
		// Also check if it's "City, Province"
		parts := strings.Split(geocoding, ", ")
		if len(parts) > 1 {
			prov := parts[len(parts)-1]
			for name, code := range provinces {
				if prov == name || prov == code {
					return "CA-" + code
				}
			}
		}
	}

	return ""
}
