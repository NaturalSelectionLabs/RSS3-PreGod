package misskey

import (
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/datatype"
)

//nolint:tagliatelle // format is required by Misskey API
type TimelineRequest struct {
	UserId         string `json:"userId"`
	IncludeReplies bool   `json:"includeReplies"`
	Renote         bool   `json:"renote"`
	UntilDate      int64  `json:"untilDate"`
	SinceDate      int64  `json:"sinceDate"`
	Limit          int    `json:"limit"`
	ExcludeNsfw    bool   `json:"excludeNsfw"`
}

type Note struct {
	Id          string
	CreatedAt   time.Time
	Summary     string
	Author      string
	Attachments datatype.Attachments
	Link        string
	Host        string // misskey host
}

type UserShow struct {
	Id   string
	Bios []string
}
