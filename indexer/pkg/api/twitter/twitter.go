package twitter

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/valyala/fastjson"
)

const endpoint = "https://api.twitter.com/1.1"

var parser fastjson.Parser

func GetTimeline() {

}

func GetUsersShow(name string) (string, error) {
	key := util.GotKey("round-robin", "Twitter", config.Config.Indexer.Twitter.Tokens)
	authorization := fmt.Sprintf("Bearer %s", key)
	logger.Infof("authorization: %s", authorization)
	var result string

	var headers = map[string]string{
		"Authorization": authorization,
	}

	url := fmt.Sprintf("%s/users/show.json?screen_name=%s", endpoint, name)

	response, err := util.Get(url, headers)
	if err != nil {
		return result, err
	}

	userShow = new(UserShow)

	parsedJson, err := parser.Parse(string(response))

	if err != nil {
		return result, err
	}

	userShow.Name = parsedJson.GetObject("name").GetString("name")

	// 这里缺一个转化函数，response转化为struct结构体

	return result, nil

}

// func userTimeline(name string, count int, useCache bool) {
// 	url := fmt.Sprintf("%s/statuses/user_timeline.json?screen_name=%s&count=%d&exclude_replies=true", endpoint, name, count)
// }

// func formatTweetText(tweet *Tweet) string {
// 	return fmt.Sprintf("%s: %s", tweet.User.Name, tweet.Text)
// }
