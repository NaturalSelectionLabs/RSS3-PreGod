package twitter

import (
	"fmt"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/datatype"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/httpx"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	lop "github.com/samber/lo/parallel"
	"github.com/valyala/fastjson"
	"regexp"
	"strings"
)

const endpoint = "https://api.twitter.com/1.1"

func GetUserShow(name string) (*UserShow, error) {
	key := util.GotKey("round-robin", "Twitter", config.Config.Indexer.Twitter.Tokens)
	authorization := fmt.Sprintf("Bearer %s", key)

	var headers = map[string]string{
		"Authorization": authorization,
	}

	url := fmt.Sprintf("%s/users/show.json?screen_name=%s", endpoint, name)

	response, err := httpx.Get(url, headers)
	if err != nil {
		return nil, err
	}

	var parser fastjson.Parser
	parsedJson, err := parser.Parse(string(response))

	if err != nil {
		return nil, err
	}

	userShow := new(UserShow)

	userShow.Name = string(parsedJson.GetStringBytes("name"))
	userShow.ScreenName = string(parsedJson.GetStringBytes("screen_name"))
	userShow.Description = string(parsedJson.GetStringBytes("description"))
	userShow.Entities = string(parsedJson.GetObject("entities").MarshalTo(nil))

	return userShow, nil
}

// TODO: offset?
// See https://developer.twitter.com/en/docs/twitter-api/v1/tweets/timelines/api-reference/get-statuses-user_timeline
func GetTimeline(name string, count uint32) ([]*ContentInfo, error) {
	key := util.GotKey("round-robin", "Twitter", config.Config.Indexer.Twitter.Tokens)
	authorization := fmt.Sprintf("Bearer %s", key)

	var headers = map[string]string{
		"Authorization": authorization,
	}

	url := fmt.Sprintf("%s/statuses/user_timeline.json?screen_name=%s&count=%d&exclude_replies=true", endpoint, name, count)

	response, err := httpx.Get(url, headers)

	if err != nil {
		return nil, err
	}

	contentInfos := make([]*ContentInfo, 0, 100)

	var parser fastjson.Parser

	parsedJson, err := parser.Parse(string(response))
	if err != nil {
		return nil, err
	}

	contentArray, err := parsedJson.Array()
	if err != nil {
		return contentInfos, err
	}

	cs := lop.Map(contentArray, func(contentValue *fastjson.Value, _ int) *ContentInfo {
		contentInfo := new(ContentInfo)

		contentInfo.PreContent, err = formatTweetText(contentValue)
		if err != nil {
			logger.Errorf("format tweet text error: %s", err)
			return nil
		}

		contentInfo.Timestamp = string(contentValue.GetStringBytes("created_at"))
		contentInfo.Hash = string(contentValue.GetStringBytes("id_str"))
		contentInfo.Link = fmt.Sprintf("https://twitter.com/%s/status/%s", name, contentInfo.Hash)
		contentInfo.Attachments = getTweetAttachments(contentValue)
		contentInfo.ScreenName = string(contentValue.GetStringBytes("user", "screen_name"))

		return contentInfo
	})

	contentInfos = append(contentInfos, cs...)

	return contentInfos, nil
}

func getTweetAttachments(contentInfo *fastjson.Value) datatype.Attachments {
	attachments := datatype.Attachments{}

	// media
	extendedEntitiesValue := contentInfo.Get("extended_entities")
	if extendedEntitiesValue != nil {
		medias := extendedEntitiesValue.GetArray("media")

		as := lop.Map(medias, func(media *fastjson.Value, _ int) datatype.Attachment {
			// TODO: video
			mediaUrl := string(media.GetStringBytes("media_url_https"))

			contentHeader, _ := httpx.GetContentHeader(mediaUrl)

			a := datatype.Attachment{
				Type:        "media",
				Address:     mediaUrl,
				MimeType:    contentHeader.MIMEType,
				SizeInBytes: contentHeader.SizeInByte,
			}

			return a
		})

		attachments = append(attachments, as...)
	}

	// quote address
	quotedStatusValue := contentInfo.Get("quoted_status")
	if quotedStatusValue != nil {
		quotedStatusId := string(quotedStatusValue.GetStringBytes("id_str"))
		quotedStatusLink := fmt.Sprintf("https://twitter.com/%s/status/%s", string(quotedStatusValue.GetStringBytes("user", "screen_name")), quotedStatusId)

		qa := datatype.Attachment{
			Type:     "quote_address",
			Address:  quotedStatusLink,
			MimeType: "text/uri-list",
		}

		attachments = append(attachments, qa)

		text := string(quotedStatusValue.GetStringBytes("text"))
		qc := datatype.Attachment{
			Type:     "quote_content",
			Address:  text,
			MimeType: "text/plain",
		}

		attachments = append(attachments, qc)
	}

	return attachments
}

// TODO: is this logic correct?
func formatTweetText(contentValue *fastjson.Value) (string, error) {
	text := contentValue.GetStringBytes("text")

	matched, err := regexp.Match("(https://t.co/[a-zA-Z0-9]+)$", text)
	if err != nil {
		return "", err
	}

	if matched {
		index := strings.Index(string(text), "https://t.co")
		text = text[:index]
	}

	extendedEntitiesValue := contentValue.Get("extended_entities")
	if extendedEntitiesValue != nil {
		media := extendedEntitiesValue.GetArray("media")
		if len(media) > 0 {
			for _, mediaItem := range media {
				mediaUrl := mediaItem.GetStringBytes("media_url_https")
				imageStr := fmt.Sprintf("<img class=\"media\" src=\"%s\">", mediaUrl)
				text = append(text, imageStr...)
			}
		}
	}

	quotedStatusValue := contentValue.Get("quoted_status")
	if quotedStatusValue != nil {
		userValue := quotedStatusValue.Get("user")
		if userValue != nil {
			screenName := userValue.GetStringBytes("screen_name")
			formatTweetStr, err := formatTweetText(quotedStatusValue)

			if err != nil {
				return "", err
			}

			quotedStatusStr := fmt.Sprintf("\nRT @%s:%s ", screenName, formatTweetStr)
			text = append(text, quotedStatusStr...)
		}
	}

	return string(text), nil
}
