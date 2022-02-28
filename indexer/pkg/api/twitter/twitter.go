package twitter

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
)

const endpoint = "https://api.twitter.com/1.1"

func GetFields() {

}

func GetTimeline() {

}

func userShow(name string) (string, error) {
	key := util.GotKey("round-robin", "Twitter")
	authorization := fmt.Sprintf("Bearer %s", key)

	var headers = map[string]string{
		"accept":        "application/json",
		"Authorization": authorization,
	}
	url := fmt.Sprintf("%s/users/show.json?screen_name=%s", endpoint, name)

	response, err := util.Get(url, headers)
	if err != nil {
		return nil, err
	}

}

func userTimeline(name string, count int, useCache bool) {
	url := fmt.Sprintf("%s/statuses/user_timeline.json?screen_name=%s&count=%d&exclude_replies=true", endpoint, name, count)
}

func formatTweetText(tweet *Tweet) string {
	return fmt.Sprintf("%s: %s", tweet.User.Name, tweet.Text)
}
