package jike

import (
	"fmt"
	"strconv"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/datatype"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/httpx"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/robfig/cron/v3"
	lop "github.com/samber/lo/parallel"
	"github.com/valyala/fastjson"
	"golang.org/x/sync/errgroup"
)

var (
	jsoni        = jsoniter.ConfigCompatibleWithStandardLibrary
	AccessToken  string
	RefreshToken string
	parser       fastjson.Parser
)

func init() {
	initLogin()
}

func initLogin() {
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

	err = jsoni.Unmarshal(response.Body, &token)

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

	parsedJson, err := parser.Parse(string(response.Body))

	if err != nil {
		logger.Errorf("Jike GetUserProfile err: %v", "error parsing response")

		return nil, err
	}

	profile := new(UserProfile)

	parsedObject := parsedJson.Get("user")

	profile.ScreenName = string(parsedObject.GetStringBytes("screenName"))
	profile.Bio = string(parsedObject.GetStringBytes("bio"))

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

	response, err := httpx.Post(url, headers, json)

	if err != nil {
		logger.Errorf("Jike GetUserTimeline err: %v", err)

		return nil, err
	}

	parsedJson, err := parser.Parse(string(response.Body))
	if err != nil {
		logger.Errorf("Jike Parsed Json:%s", err)

		return nil, err
	}

	author := string(parsedJson.GetStringBytes("data", "userProfile", "username"))

	parsedObject := parsedJson.GetArray("data", "userProfile", "feeds", "nodes")

	result := make([]Timeline, len(parsedObject))

	lop.ForEach(parsedObject, func(node *fastjson.Value, i int) {
		id := string(node.GetStringBytes("id"))
		result[i].Id = id

		t, timeErr := time.Parse(time.RFC3339, string(node.GetStringBytes("createdAt")))
		if timeErr != nil {
			logger.Errorf("Jike GetUserTimeline timestamp parsing err: %v", timeErr)

			t = time.Time{} // set to zero value
		}

		result[i].Author = author
		result[i].Timestamp = t
		result[i].Summary = string(node.GetStringBytes("content"))
		result[i].Link = fmt.Sprintf("https://web.okjike.com/originalPost/%s", id)
		result[i].Attachments = getAttachment(node)
	})

	if err != nil {
		logger.Errorf("Jike GetUserTimeline err: %v", "error parsing response")

		return nil, err
	}

	return result, err
}

// func formatFeed(node *fastjson.Value) string {
// 	text := string(node.GetStringBytes("content"))
//
// 	if node.Exists("pictures") {
// 		for _, picture := range node.GetArray("pictures") {
// 			var url string
//
// 			if picture.Exists("picUrl") {
// 				url = string(picture.GetStringBytes("picUrl"))
// 			}
//
// 			if picture.Exists("thumbnailUrl") {
// 				url = string(picture.GetStringBytes("thumbnailUrl"))
// 			}
//
// 			text += fmt.Sprintf("<img class=\"media\" src=\"%s\">", string(url))
// 		}
// 	}
//
// 	if node.Exists("target") && string(node.GetStringBytes("type")) == "REPOST" {
// 		target := node.Get("target")
// 		// a status key means the feed is unavailable, e.g, DELETED
// 		if !target.Exists("status") {
// 			var user string
// 			if target.Exists("user", "screenName") {
// 				user = string(target.GetStringBytes("user", "screenName"))
// 			}
//
// 			text += fmt.Sprintf("\nRT %s: %s", user, formatFeed(target))
// 		}
// 	}
//
// 	return text
// }

func getAttachment(node *fastjson.Value) []datatype.Attachment {
	var content string

	attachments := make([]datatype.Attachment, 0)

	// process the original post attachments
	g := new(errgroup.Group)

	g.Go(func() error {
		attachments = append(attachments, getPicture(node)...)

		return nil
	})

	g.Go(func() error {
		attachments = append(attachments, getVideo(node)...)

		return nil
	})

	// a 'status' field often means the report target is unavailable, e.g, DELETED
	if !node.Exists("target", "status") {
		if node.Exists("target") {
			node = node.Get("target")

			// store quote_address
			qAddress := datatype.Attachment{
				Content:  "https://web.okjike.com/originalPost/" + string(node.GetStringBytes("id")),
				MimeType: "text/uri-list",
				Type:     "quote_address",
			}

			// store quote_text
			content = string(node.GetStringBytes("content"))
			qText := datatype.Attachment{
				Content:  content,
				MimeType: "text/plain",
				Type:     "quote_text",
			}

			attachments = append(attachments, qAddress, qText)

			// store quote_media

			if node.Exists("pictures") {
				attachments = append(attachments, getPicture(node)...)
			}
		}
	}

	_ = g.Wait()

	return attachments
}

func getPicture(node *fastjson.Value) []datatype.Attachment {
	pics := node.GetArray("pictures")

	result := make([]datatype.Attachment, len(pics))

	lop.ForEach(pics, func(pic *fastjson.Value, i int) {
		var url string

		if pic.Exists("picUrl") {
			url = string(pic.GetStringBytes("picUrl"))
		} else if pic.Exists("thumbnailUrl") {
			url = string(pic.GetStringBytes("thumbnailUrl"))
		}

		contentHeader, _ := httpx.GetContentHeader(url)

		// qMedia := model.NewAttachment(url, address, contentHeader.MIMEType, "quote_media", contentHeader.SizeInByte, time.Now())
		qMedia := datatype.Attachment{
			Type:        "quote_media",
			Address:     url,
			MimeType:    contentHeader.MIMEType,
			SizeInBytes: contentHeader.SizeInByte,
		}

		result[i] = qMedia
	})

	return result
}

func getVideo(node *fastjson.Value) []datatype.Attachment {
	videoPic := string(node.GetStringBytes("video", "image", "picUrl"))

	if videoPic != "" {
		contentHeader, _ := httpx.GetContentHeader(videoPic)

		videoAttachment := datatype.Attachment{
			Type:        "quote_media",
			Address:     videoPic,
			MimeType:    contentHeader.MIMEType,
			SizeInBytes: contentHeader.SizeInByte,
		}

		return []datatype.Attachment{videoAttachment}
	}

	return nil
}
