package fetcher

import "github.com/emersion/go-vcard"

type Fetcher interface {
	Fetch() ([]vcard.Card, error)
}
