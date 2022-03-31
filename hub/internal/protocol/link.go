package protocol

import (
	"encoding/json"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/isotime"
)

var _ json.Marshaler = &Link{}

type Link struct {
	DateCreated time.Time    `json:"date_created"`
	From        string       `json:"from"`
	To          string       `json:"to"`
	Source      string       `json:"source"`
	Metadata    LinkMetadata `json:"metadata"`
}

func (l Link) MarshalJSON() ([]byte, error) {
	return json.Marshal(&magicLink{
		DateCreated: l.DateCreated.Format(isotime.ISO8601),
		From:        l.From,
		To:          l.To,
		Source:      l.Source,
		Metadata:    l.Metadata,
	})
}

type magicLink struct {
	DateCreated string       `json:"date_created"`
	From        string       `json:"from"`
	To          string       `json:"to"`
	Source      string       `json:"source"`
	Metadata    LinkMetadata `json:"metadata"`
}

type LinkMetadata struct {
	Network string `json:"network"`
	Proof   string `json:"proof"`
}
