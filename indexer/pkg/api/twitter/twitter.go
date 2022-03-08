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

func GetUsersShow(name string) (*UserShow, error) {
	key := util.GotKey("round-robin", "Twitter", config.Config.Indexer.Twitter.Tokens)
	authorization := fmt.Sprintf("Bearer %s", key)
	logger.Infof("authorization: %s", authorization)

	var headers = map[string]string{
		"Authorization": authorization,
	}

	url := fmt.Sprintf("%s/users/show.json?screen_name=%s", endpoint, name)

	response, err := util.Get(url, headers)
	if err != nil {
		return nil, err
	}

	parsedJson, err := parser.Parse(string(response))

	if err != nil {
		return nil, err
	}

	userShow := new(UserShow)

	userShow.Name = string(parsedJson.GetStringBytes("name"))
	userShow.ScreenName = string(parsedJson.GetStringBytes("screen_name"))
	userShow.Description = string(parsedJson.GetStringBytes("description"))

	return userShow, nil
}

func GetTimeline(name string, count uint32) (*ContentInfo, error) {
	key := util.GotKey("round-robin", "Twitter", config.Config.Indexer.Twitter.Tokens)
	authorization := fmt.Sprintf("Bearer %s", key)
	logger.Infof("authorization: %s", authorization)

	var headers = map[string]string{
		"Authorization": authorization,
	}

	url := fmt.Sprintf("%s/statuses/user_timeline.json?screen_name=%scount=%d&exclude_replies=true", endpoint, name, count)
	logger.Infof("url: %s", url)

	response, err := util.Get(url, headers)
	if err != nil {
		return nil, err
	}

	logger.Infof("response: %s", string(response))

	ContentInfo := new(ContentInfo)

	return ContentInfo, nil
}

// func userTimeline(name string, count int, useCache bool) {
// 	url := fmt.Sprintf("%s/statuses/user_timeline.json?screen_name=%s&count=%d&exclude_replies=true", endpoint, name, count)
// }

// func formatTweetText(tweet *Tweet) string {
// 	return fmt.Sprintf("%s: %s", tweet.User.Name, tweet.Text)
// }
