package protocol

import "time"

type Profile struct {
	DateCreated       time.Time           `json:"date_created"`
	DateUpdated       time.Time           `json:"date_updated"`
	Name              string              `json:"name"`
	Avatars           []string            `json:"avatars"`
	Bio               string              `json:"bio"`
	Attachments       []ProfileAttachment `json:"attachments"`
	ConnectedAccounts []string            `json:"connected_accounts"`
	Source            string              `json:"source"`
	Metadata          ProfileMetadata     `json:"metadata"`
}

type ProfileAttachment struct {
	Type     string `json:"type"`
	Content  string `json:"content"`
	MimeType string `json:"mime_type"`
}

type ProfileMetadata struct {
	Network string `json:"network"`
	Proof   string `json:"proof"`
}
