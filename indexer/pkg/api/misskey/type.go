package misskey

import "time"

//nolint:tagliatelle // format is required by Misskey API
type TimelineRequestStruct struct {
	UserId         string `json:"userId"`
	IncludeReplies bool   `json:"includeReplies"`
	Renote         bool   `json:"renote"`
	UntilDate      int64  `json:"untilDate"`
	Limit          int    `json:"limit"`
	ExcludeNsfw    bool   `json:"excludeNsfw"`
}

type NoteStruct struct {
	Id        string
	CreatedAt time.Time
	Text      string
	Author    string
}
