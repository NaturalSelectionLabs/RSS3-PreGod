package protocol

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/timex"
)

type Link struct {
	DateCreated timex.Time   `json:"date_created"`
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
