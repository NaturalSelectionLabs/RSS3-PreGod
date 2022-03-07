package twitter

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
)

const endpoint = "https://api.twitter.com/1.1"

func GetTimeline() {

}

func GetFields(name string) (string, error) {
	key := util.GotKey("round-robin", "Twitter", config.Config.Indexer.Twitter.Tokens)
	authorization := fmt.Sprintf("Bearer %s", key)
	var result string

	var headers = map[string]string{
		"accept":        "application/json",
		"Authorization": authorization,
	}

	url := fmt.Sprintf("%s/users/show.json?screen_name=%s", endpoint, name)

	response, err := util.Get(url, headers)
		return result, err
		return nil, err
	}

	// 这里缺一个转化函数，response转化为struct结构体

	return result, nil

}

// func userTimeline(name string, count int, useCache bool) {
// 	url := fmt.Sprintf("%s/statuses/user_timeline.json?screen_name=%s&count=%d&exclude_replies=true", endpoint, name, count)
// }

// func formatTweetText(tweet *Tweet) string {
// 	return fmt.Sprintf("%s: %s", tweet.User.Name, tweet.Text)
// }
