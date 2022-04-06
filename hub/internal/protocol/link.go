package protocol

import (
	"encoding/json"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/timex"
)

var _ json.Marshaler = &Link{}

type Link struct {
	DateCreated time.Time    `json:"date_created"`
	From        string       `json:"from"`
	To          string       `json:"to"`
	Type        string       `json:"type"`
	Source      string       `json:"source"`
	Metadata    LinkMetadata `json:"metadata"`
}

func (l Link) MarshalJSON() ([]byte, error) {
	return json.Marshal(&magicLink{
		DateCreated: l.DateCreated.Format(timex.ISO8601),
		From:        l.From,
		To:          l.To,
		Type:        l.Type,
		Source:      l.Source,
		Metadata:    l.Metadata,
	})
}

type magicLink struct {
	DateCreated string       `json:"date_created"`
	From        string       `json:"from"`
	To          string       `json:"to"`
	Type        string       `json:"type"`
	Source      string       `json:"source"`
	Metadata    LinkMetadata `json:"metadata"`
}

type LinkMetadata struct {
	Network string `json:"network"`
	Proof   string `json:"proof"`
}
