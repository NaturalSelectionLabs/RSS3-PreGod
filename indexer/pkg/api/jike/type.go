package jike

//nolint:tagliatelle // format is required by Jike API
type UserProfileStruct struct {
	ScreenName string `json:"screenName"`
	Bio        string `json:"bio"`
}

type TimelineStruct struct {
	Hash       string `json:"hash"`
	Timestamp  string `json:"timestamp"`
	PreContent string `json:"pre_content"`
	Link       string `json:"link"`
}

//nolint:tagliatelle // format is required by Jike API
type RefreshTokenStruct struct {
	AccessToken  string `json:"x-jike-access-token"`
	RefreshToken string `json:"x-jike-refresh-token"`
	Success      bool   `json:"success"`
}

//nolint:tagliatelle // format is required by Jike API
type TimelineRequestStruct struct {
	OperationName string `json:"operationName"`
	Variables     struct {
		Username string `json:"username"`
	} `json:"variables"`
	Query string `json:"query"`
}
