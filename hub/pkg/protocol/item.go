package protocol

type Item struct {
	Base

	Links Links `json:"links"`

	Tags    []string `json:"tags,omitempty"`
	Authors []string `json:"authors,omitempty"`
	Title   string   `json:"title,omitempty"`
	Summary string   `json:"summary,omitempty"`

	Attachments []ItemAttachment `json:"attachments,omitempty"`

	Metadata ItemMetadata `json:"metadata,omitempty"`

	Auto bool `json:"auto"`
}

type ItemAttachment struct {
	Content  string `json:"content,omitempty"`       // Actual content, mutually exclusive with `address`.
	Address  string `json:"address,omitempty"`       // URIs of same resource pointing to third parties, mutually exclusive with content.
	MimeType string `json:"mime_type"`               // [MIME type](https://en.wikipedia.org/wiki/Media_type)
	Name     string `json:"name,omitempty"`          // Name of the attachment.
	Size     int    `json:"size_in_bytes,omitempty"` // Size of the attachment in bytes.
}

type ItemMetadata struct {
	Proof string `json:"proof"`
	Type  string `json:"type"`
	From  string `json:"from,omitempty"`
	To    string `json:"to,omitempty"`
	Id    string `json:"id"`
}
