package protocol

import (
	"encoding/json"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/isotime"
	"time"
)

var _ json.Marshaler = &LinkItem{}

type LinkItem struct {
	DateCreated time.Time        `json:"date_created"`
	From        string           `json:"from"`
	To          string           `json:"to"`
	Source      string           `json:"source"`
	Metadata    LinkItemMetadata `json:"metadata"`
}

func (l LinkItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(&magicLinkItem{
		DateCreated: l.DateCreated.Format(isotime.ISO8601),
		From:        l.From,
		To:          l.To,
		Source:      l.Source,
		Metadata:    l.Metadata,
	})
}

type magicLinkItem struct {
	DateCreated string           `json:"date_created"`
	From        string           `json:"from"`
	To          string           `json:"to"`
	Source      string           `json:"source"`
	Metadata    LinkItemMetadata `json:"metadata"`
}

type LinkItemMetadata struct {
	Network string `json:"network"`
	Proof   string `json:"proof"`
}
