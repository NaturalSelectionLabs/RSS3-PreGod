package protocol

type Index struct {
	SignedBase

	Profile IndexProfile `json:"profile"`
	Links   IndexLinks   `json:"links"`
	Items   IndexItems   `json:"items"`
}

type IndexProfile struct {
	Name        string            `json:"name"`
	Avatars     []string          `json:"avatars"`
	Bio         string            `json:"bio"`
	Attachments []IndexAttachment `json:"attachment,omitempty"`
	Accounts    []IndexAccount    `json:"accounts"`
}

type IndexAttachment struct {
	Type     string `json:"type"`
	Content  string `json:"content"`
	MimeType string `json:"mime_type"`
}

type IndexAccount struct {
	Identifier string `json:"identifier"`
	Signature  string `json:"signature,omitempty"`
}

type IndexLinks struct {
	Identifiers    []IndexLinkIdentifier `json:"identifiers"`
	IdentifierBack string                `json:"identifier_back"`
}

type IndexLinkIdentifier struct {
	Type             string `json:"type"`
	IdentifierCustom string `json:"identifier_custom"`
	Identifier       string `json:"identifier"`
}

type IndexItems struct {
	Notes  IndexItemsNotes  `json:"notes"`
	Assets IndexItemsAssets `json:"assets"`
}

type IndexItemsNotes struct {
	IdentifierCustom string                  `json:"identifier_custom"`
	Identifier       string                  `json:"identifier"`
	Filters          *IndexItemsNotesFilters `json:"filters,omitempty"`
}

type IndexItemsNotesFilters struct {
	AllowList []string `json:"allowlist,omitempty"`
	BlockList []string `json:"blocklist,omitempty"`
}

type IndexItemsAssets struct {
	IdentifierCustom string                   `json:"identifier_custom"`
	Identifier       string                   `json:"identifier"`
	Filters          *IndexItemsAssetsFilters `json:"filters,omitempty"`
}

type IndexItemsAssetsFilters struct {
	AllowList []string `json:"allowlist,omitempty"`
	BlockList []string `json:"blocklist,omitempty"`
}
