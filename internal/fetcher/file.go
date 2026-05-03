package fetcher

import (
	"os"

	"github.com/emersion/go-vcard"
)

type FileFetcher struct {
	Path string
}

func (f *FileFetcher) Fetch() ([]vcard.Card, error) {
	file, err := os.Open(f.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dec := vcard.NewDecoder(file)
	var cards []vcard.Card
	for {
		card, err := dec.Decode()
		if err != nil {
			break
		}
		cards = append(cards, card)
	}
	return cards, nil
}
