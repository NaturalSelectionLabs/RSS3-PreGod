package jike

import (
	"fmt"
	"strconv"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util/httpx"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/robfig/cron/v3"
	"github.com/valyala/fastjson"
)

var (
	jsoni        = jsoniter.ConfigCompatibleWithStandardLibrary
	AccessToken  string
	RefreshToken string
	parser       fastjson.Parser
)

func New() {
	Login()

	// everyday at 00:00, refresh Jike tokens
	c := cron.New()
	c.AddFunc("0 0 * * *", func() { Login() })
	c.Start()
}

func Login() error {
	json, err := jsoni.MarshalToString(config.Config.Indexer.Jike)

	if err != nil {
		logger.Errorf("Jike Config read err: %v", err)

		return err
	}

	headers := map[string]string{
		"App-Version":  config.Config.Indexer.Jike.AppVersion,
		"Content-Type": "application/json",
	}

	url := "https://api.ruguoapp.com/1.0/users/loginWithPhoneAndPassword"

	response, err := httpx.PostRaw(url, headers, json)

	if err != nil {
		logger.Errorf("Jike Login err: %v", err)

		return err
	}

	AccessToken = string(response.Header().Get("x-jike-access-token"))
	RefreshToken = string(response.Header().Get("x-jike-refresh-token"))

	return nil
}

func RefreshJikeToken() error {
	headers := map[string]string{
		"App-Version":          config.Config.Indexer.Jike.AppVersion,
		"Content-Type":         "application/json",
		"x-jike-refresh-token": RefreshToken,
	}

	url := "https://api.ruguoapp.com/app_auth_tokens.refresh"

	response, err := httpx.Get(url, headers)

	if err != nil {
		logger.Errorf("Jike RefreshToken err: %v", err)

		return err
	}

	token := new(RefreshTokenStruct)

	err = jsoni.Unmarshal(response, &token)

	if err == nil {
		if token.Success {
			AccessToken = token.AccessToken
			RefreshToken = token.RefreshToken

			return nil
		} else {
			logger.Errorf("Jike RefreshToken err: %v", "Jike refresh token endpoint returned a failed response")

			return err
		}
	} else {
		logger.Errorf("Jike RefreshToken err: %v", err)

		return err
	}
}

func GetUserProfile(name string) (*UserProfile, error) {
	refreshErr := RefreshJikeToken()

	if refreshErr != nil {
		return nil, refreshErr
	}

	headers := map[string]string{
		"App-Version":         config.Config.Indexer.Jike.AppVersion,
		"Content-Type":        "application/json",
		"x-jike-access-token": AccessToken,
	}

	url := "https://api.ruguoapp.com/1.0/users/profile?username=" + name

	response, err := httpx.Get(url, headers)

	if err != nil {
		logger.Errorf("Jike GetUserProfile err: %v", err)

		return nil, err
	}

	parsedJson, err := parser.Parse(string(response))

	if err != nil {
		logger.Errorf("Jike GetUserProfile err: %v", "error parsing response")

		return nil, err
	}

	profile := new(UserProfile)

	parsedObject := parsedJson.GetObject("user")

	profile.ScreenName = parsedObject.Get("screenName").String()
	profile.Bio = parsedObject.Get("bio").String()

	return profile, err
}

// nolint:funlen // format is required by Jike API
func GetUserTimeline(name string) ([]Timeline, error) {
	refreshErr := RefreshJikeToken()

	if refreshErr != nil {
		return nil, refreshErr
	}

	headers := map[string]string{
		"App-Version":  config.Config.Indexer.Jike.AppVersion,
		"Content-Type": "application/json",
		// nolint:lll // format is required by Jike API
		"cookie": "fetchRankedUpdate=" + strconv.FormatInt(time.Now().UnixNano(), 10) + "; x-jike-access-token=" + AccessToken + "; x-jike-refresh-token=" + RefreshToken,
	}

	data := new(TimelineRequest)

	data.OperationName = "UserFeeds"
	data.Variables.Username = name

	data.Query = `query UserFeeds($username: String!) {
					userProfile(username: $username) {
						username
						screenName
						briefIntro
						feeds {
						...BasicFeedItem
						}
					}
				}

				fragment BasicFeedItem on FeedsConnection {
					nodes {
						... on ReadSplitBar {
							id
							type
							text
						}
						... on MessageEssential {
							...FeedMessageFragment
						}
					}
				}

				fragment FeedMessageFragment on MessageEssential {
					...EssentialFragment
					... on OriginalPost {
						...MessageInfoFragment
					}
					... on Repost {
						...RepostFragment
					}
				}

				fragment EssentialFragment on MessageEssential {
					id
					type
					content
					createdAt
					pictures {
						format
						picUrl
						thumbnailUrl
					}
				}

				fragment TinyUserFragment on UserInfo {
					screenName
				}

				fragment MessageInfoFragment on MessageInfo {
					video {
						title
						type
						image {
							picUrl
						}
					}
				}

				fragment RepostFragment on Repost {
					target {
						...RepostTargetFragment
					}
				}

				fragment RepostTargetFragment on RepostTarget {
					... on OriginalPost {
						id
						type
						content
						pictures {
							thumbnailUrl
						}
						user {
							...TinyUserFragment
						}
					}
					... on Repost {
						id
						type
						content
						pictures {
							thumbnailUrl
						}
					}
					... on DeletedRepostTarget {
						status
					}
				}
`

	url := "https://web-api.okjike.com/api/graphql"

	json, _ := jsoni.MarshalToString(data)

	response, err := httpx.PostRaw(url, headers, json)

	if err != nil {
		logger.Errorf("Jike GetUserTimeline err: %v", err)

		return nil, err
	}

	parsedJson, err := parser.Parse(string(response.Body()))

	parsedObject := parsedJson.GetArray("data", "userProfile", "feeds", "nodes")

	result := make([]Timeline, len(parsedObject))

	for i, node := range parsedObject {
		id := util.TrimQuote(node.Get("id").String())
		result[i].Id = id

		t, timeErr := time.Parse(time.RFC3339, util.TrimQuote(node.Get("createdAt").String()))
		if err != nil {
			logger.Errorf("Jike GetUserTimeline timestamp parsing err: %v", timeErr)

			return nil, timeErr
		}

		result[i].Author = util.TrimQuote(parsedJson.Get("username").String())
		result[i].Timestamp = t
		result[i].Summary = util.TrimQuote(node.Get("content").String())
		result[i].Link = fmt.Sprintf("https://web.okjike.com/originalPost/%s", id)
		result[i].Attachments = *getAttachment(node)
	}

	if err != nil {
		logger.Errorf("Jike GetUserTimeline err: %v", "error parsing response")

		return nil, err
	}

	return result, err
}

//nolint:unused // might need it in the future
func formatFeed(node *fastjson.Value) string {
	text := util.TrimQuote(node.Get("content").String())

	if node.Exists("pictures") {
		for _, picture := range node.GetArray("pictures") {
			var url string

			if picture.Exists("picUrl") {
				url = picture.Get("picUrl").String()
			}

			if picture.Exists("thumbnailUrl") {
				url = picture.Get("thumbnailUrl").String()
			}

			text += fmt.Sprintf("<img class=\"media\" src=\"%s\">", util.TrimQuote(url))
		}
	}

	if node.Exists("target") && node.Get("type").String() == "REPOST" {
		target := node.Get("target")
		// a status key means the feed is unavailable, e.g, DELETED
		if !target.Exists("status") {
			var user string
			if target.Exists("user", "screenName") {
				user = util.TrimQuote(target.Get("user", "screenName").String())
			}

			text += fmt.Sprintf("\nRT %s: %s", user, formatFeed(target))
		}
	}

	return text
}

// TODO: handle video attachments
func getAttachment(node *fastjson.Value) *[]model.Attachment {
	var content string

	attachments := make([]model.Attachment, 0)

	// process the original post attachments
	attachments = append(attachments, *getPicture(node.Get("pictures"))...)

	// a 'status' field often means the report target is unavailable, e.g, DELETED
	if !node.Exists("target", "status") {
		if node.Exists("target") {
			node = node.Get("target")
			// store quote_address

			content = "https://web.okjike.com/originalPost/" + util.TrimQuote(node.Get("id").String())

			syncAt := time.Now()

			qAddress := *model.NewAttachment(content, nil, "text/uri-list", "quote_address", 0, syncAt)

			// store quote_text

			content = util.TrimQuote(node.Get("content").String())
			qText := *model.NewAttachment(content, nil, "text/plain", "quote_text", 0, syncAt)

			attachments = append(attachments, qAddress, qText)

			// store quote_media

			if node.Exists("pictures") {
				attachments = append(attachments, *getPicture(node)...)
			}
		}
	}

	return &attachments
}

func getPicture(node *fastjson.Value) *[]model.Attachment {
	address := make([]string, 1)

	var mime, content string

	var sizeInBytes = 0

	pictues := node.GetArray("pictures")

	result := make([]model.Attachment, len(pictues))

	for i, picture := range pictues {
		var url string

		if picture.Exists("thumbnailUrl") {
			url = util.TrimQuote(picture.Get("thumbnailUrl").String())
		}

		if picture.Exists("picUrl") {
			url = util.TrimQuote(picture.Get("picUrl").String())
		}

		header, err := httpx.Head(url)

		if err == nil {
			sizeInBytes, _ = strconv.Atoi(header.Get("Content-Length"))
			mime = header.Get("Content-Type")
		}

		address = append(address, url)

		qMedia := *model.NewAttachment(content, address, mime, "quote_media", sizeInBytes, time.Now())
		result[i] = qMedia
	}

	return &result
}
