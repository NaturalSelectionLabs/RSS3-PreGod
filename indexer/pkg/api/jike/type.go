package jike

import (
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
)

//nolint:tagliatelle // format is required by Jike API
type UserProfile struct {
	ScreenName string `json:"screenName"`
	Bio        string `json:"bio"`
}

type Timeline struct {
	Id          string
	Timestamp   time.Time
	Summary     string
	Author      string
	Attachments []model.Attachment
	Link        string
}

//nolint:tagliatelle // format is required by Jike API
type RefreshTokenStruct struct {
	AccessToken  string `json:"x-jike-access-token"`
	RefreshToken string `json:"x-jike-refresh-token"`
	Success      bool   `json:"success"`
}

//nolint:tagliatelle // format is required by Jike API
type TimelineRequest struct {
	OperationName string `json:"operationName"`
	Variables     struct {
		Username string `json:"username"`
	} `json:"variables"`
	Query string `json:"query"`
}
