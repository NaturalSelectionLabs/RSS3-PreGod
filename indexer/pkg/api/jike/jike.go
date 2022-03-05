package jike

import (
	"fmt"
	"strconv"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
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

func init() {
	c := cron.New()
	// everyday at 00:00, refresh Jike tokens
	c.AddFunc("0 0 * * *", func() { Login() })
	c.Start()
}

func Login() error {
	json, err := jsoni.MarshalToString(config.Config.Indexer.Jike)

	if err != nil {
		logger.Fatalf("Jike Config read err: %v", err)

		return err
	}

	headers := map[string]string{
		"App-Version":  config.Config.Indexer.Jike.AppVersion,
		"Content-Type": "application/json",
	}

	url := "https://api.ruguoapp.com/1.0/users/loginWithPhoneAndPassword"

	response, err := util.PostRaw(url, headers, json)

	if err != nil {
		logger.Fatalf("Jike Login err: %v", err)

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

	response, err := util.Get(url, headers)

	if err != nil {
		logger.Fatalf("Jike RefreshToken err: %v", err)

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
			logger.Fatalf("Jike RefreshToken err: %v", "Jike refresh token endpoint returned a failed response")

			return err
		}
	} else {
		logger.Fatalf("Jike RefreshToken err: %v", err)

		return err
	}
}

func GetUserProfile(name string) (*UserProfileStruct, error) {
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

	response, err := util.Get(url, headers)

	if err != nil {
		logger.Fatalf("Jike GetUserProfile err: %v", err)

		return nil, err
	}

	parsedJson, err := parser.Parse(string(response))

	if err != nil {
		logger.Fatalf("Jike GetUserProfile err: %v", "error parsing response")

		return nil, err
	}

	profile := new(UserProfileStruct)

	parsedObject := parsedJson.GetObject("user")

	profile.ScreenName = parsedObject.Get("screenName").String()
	profile.Bio = parsedObject.Get("bio").String()

	return profile, err
}

// nolint:funlen // format is required by Jike API
func GetUserTimeline(name string) ([]TimelineStruct, error) {
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

	data := new(TimelineRequestStruct)

	data.OperationName = "UserFeeds"
	data.Variables.Username = name

	data.Query = `query UserFeeds($username: String!) {
					userProfile(username: $username) {
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

	response, err := util.PostRaw(url, headers, json)

	if err != nil {
		logger.Fatalf("Jike GetUserTimeline err: %v", err)

		return nil, err
	}

	parsedJson, err := parser.Parse(string(response.Body()))

	parsedObject := parsedJson.GetArray("data", "userProfile", "feeds", "nodes")

	result := make([]TimelineStruct, len(parsedObject))

	for i, node := range parsedObject {
		t, timeErr := time.Parse(time.RFC3339, trimQuote(node.Get("createdAt").String()))
		if err != nil {
			logger.Fatalf("Jike GetUserTimeline timestamp parsing err: %v", timeErr)

			return nil, timeErr
		}

		id := trimQuote(node.Get("id").String())
		result[i].Hash = id
		result[i].Timestamp = fmt.Sprintf("0x%x", t.Unix())
		result[i].PreContent = formatFeed(node)
		result[i].Link = fmt.Sprintf("https://web.okjike.com/originalPost/%s", id)
	}

	if err != nil {
		logger.Fatalf("Jike GetUserTimeline err: %v", "error parsing response")

		return nil, err
	}

	return result, err
}

func formatFeed(node *fastjson.Value) string {
	text := trimQuote(node.Get("content").String())

	if node.Exists("pictures") {
		for _, picture := range node.GetArray("pictures") {
			var url string

			if picture.Exists("picUrl") {
				url = picture.Get("picUrl").String()
			}

			if picture.Exists("thumbnailUrl") {
				url = picture.Get("thumbnailUrl").String()
			}

			text += fmt.Sprintf("<img class=\"media\" src=\"%s\">", trimQuote(url))
		}
	}

	if node.Exists("target") {
		target := node.Get("target")
		// a status key means the feed is unavailable, e.g, DELETED
		if !target.Exists("status") {
			var user string
			if target.Exists("user", "screenName") {
				user = trimQuote(target.Get("user", "screenName").String())
			}

			text += fmt.Sprintf("\nRT %s: %s", user, formatFeed(target))
		}
	}

	return text
}

func trimQuote(s string) string {
	return s[1 : len(s)-1]
}
