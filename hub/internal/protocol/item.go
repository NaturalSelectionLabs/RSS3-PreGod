package protocol

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/timex"
)

type Item struct {
	Identifier  string                 `json:"identifier"`
	DateCreated timex.Time             `json:"date_created"`
	DateUpdated timex.Time             `json:"date_updated"`
	RelatedURLs []string               `json:"related_urls,omitempty"`
	Links       string                 `json:"links"`
	BackLinks   string                 `json:"backlinks"`
	Tags        []string               `json:"tags,omitempty"`
	Authors     []string               `json:"authors"`
	Title       string                 `json:"title,omitempty"`
	Summary     string                 `json:"summary,omitempty"`
	Attachments []ItemAttachment       `json:"attachments,omitempty"`
	Source      string                 `json:"source"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type ItemAttachment struct {
	Type        string `json:"type,omitempty"`
	Content     string `json:"content,omitempty"`
	Address     string `json:"address,omitempty"`
	MimeType    string `json:"mime_type"`
	SizeInBytes int64  `json:"size_in_bytes,omitempty"`
}
