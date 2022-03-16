package protocol

type NoteList struct {
	SignedBase

	Total int            `json:"total"`
	List  []NoteListItem `json:"list"`
}

type NoteListItem struct {
	Identifier  string                   `json:"identifier"`
	DateCreated string                   `json:"date_created"`
	DateUpdated string                   `json:"date_updated"`
	Links       NoteListItemLinks        `json:"links"`
	Tags        []string                 `json:"tags"`
	Authors     []string                 `json:"authors"`
	Summary     []string                 `json:"summary"`
	Attachments []NoteListItemAttachment `json:"attachments"`
}

type NoteListItemLinks struct {
	IdentifierBack string `json:"identifier_back"`
}

type NoteListItemTag []string

type NoteListItemAttachment struct {
	Address     string `json:"address"`
	MimeType    string `json:"mime_type"`
	SizeInBytes int    `json:"size_in_bytes"`
}
