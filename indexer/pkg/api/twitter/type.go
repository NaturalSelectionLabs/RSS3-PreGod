package twitter

import "time"

type UserShow struct {
	Name        string
	Description string
	ScreenName  string
}

type ContentInfo struct {
	PreContent string
	Timestamp  string
	Hash       string
	Link       string
}

func (i ContentInfo) GetTsp() (time.Time, error) {
	return time.Parse(time.RubyDate, i.Timestamp)
}
