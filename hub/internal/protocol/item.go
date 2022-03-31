package protocol

import "time"

type Item struct {
	Identifier  string           `json:"identifier"`
	DateCreated time.Time        `json:"date_created"`
	DateUpdated time.Time        `json:"date_updated"`
	RelatedURLs []string         `json:"related_urls"`
	Links       string           `json:"links"`
	BackLinks   string           `json:"backlinks"`
	Tags        []string         `json:"tags"`
	Authors     []string         `json:"authors"`
	Title       string           `json:"title"`
	Summary     string           `json:"summary"`
	Attachments []ItemAttachment `json:"attachments"`
}

type ItemAttachment struct {
	Type        string `json:"type"`
	Address     string `json:"address"`
	MimeType    string `json:"mime_type"`
	SizeInBytes int64  `json:"size_in_bytes"`
}
