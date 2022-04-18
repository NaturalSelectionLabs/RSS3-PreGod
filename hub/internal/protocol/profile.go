package protocol

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/timex"
)

type Profile struct {
	DateCreated       timex.Time          `json:"date_created"`
	DateUpdated       timex.Time          `json:"date_updated"`
	Name              string              `json:"name,omitempty"`
	Avatars           []string            `json:"avatars,omitempty"`
	Bio               string              `json:"bio"`
	Tags              []string            `json:"tags,omitempty"`
	RelatedURLs       []string            `json:"related_urls,omitempty"`
	Attachments       []ProfileAttachment `json:"attachments"`
	ConnectedAccounts []string            `json:"connected_accounts,omitempty"`
	Source            string              `json:"source"`
	Metadata          ProfileMetadata     `json:"metadata"`
}

type ProfileAttachment struct {
	Type     string `json:"type"`
	Content  string `json:"content,omitempty"`
	Address  string `json:"address,omitempty"`
	MimeType string `json:"mime_type"`
}

type ProfileMetadata struct {
	Network string `json:"network"`
	Proof   string `json:"proof"`
}
