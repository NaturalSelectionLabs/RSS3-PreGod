package twitter

import (
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/datatype"
)

type UserShow struct {
	Name        string
	Description string
	ScreenName  string
	Entities    string
}

type ContentInfo struct {
	PreContent  string
	Timestamp   string
	Hash        string
	Link        string
	ScreenName  string
	Attachments datatype.Attachments
}

func (i ContentInfo) GetTsp() (time.Time, error) {
	return time.Parse(time.RubyDate, i.Timestamp)
}
