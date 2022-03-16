package protocol

type AssetList struct {
	SignedBase

	Total int             `json:"total"`
	List  []AssetListItem `json:"list"`
}

type AssetListItem struct {
	Identifier  string                    `json:"identifier"`
	DateCreated string                    `json:"date_created"`
	DateUpdated string                    `json:"date_updated"`
	Links       AssetListItemLinks        `json:"links"`
	Tags        []string                  `json:"tags"`
	Authors     []string                  `json:"authors"`
	Summary     []string                  `json:"summary"`
	Attachments []AssetListItemAttachment `json:"attachments"`
}

type AssetListItemLinks struct {
	IdentifierBack string `json:"identifier_back"`
}

type AssetListItemTag []string

type AssetListItemAttachment struct {
	Address     string `json:"address"`
	MimeType    string `json:"mime_type"`
	SizeInBytes int    `json:"size_in_bytes"`
}
