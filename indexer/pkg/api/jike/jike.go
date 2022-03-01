package jike

import (
	"fmt"
	"strconv"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/types"
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

	token := new(types.RefreshTokenStruct)

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

func GetUserProfile(name string) (*types.UserProfileStruct, error) {
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

	profile := new(types.UserProfileStruct)

	parsedObject := parsedJson.GetObject("user")

	profile.ScreenName = parsedObject.Get("screenName").String()
	profile.Bio = parsedObject.Get("bio").String()

	return profile, err
}

func GetUserTimeline(name string) ([]types.TimelineStruct, error) {
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

	data := new(types.TimelineRequestStruct)

	data.OperationName = "UserFeeds"
	data.Variables.Username = name

	// nolint:lll // format is required by Jike API
	data.Query = "query UserFeeds($username: String!, $loadMoreKey: JSON) {  userProfile(username: $username) {    username    feeds(loadMoreKey: $loadMoreKey) {      ...BasicFeedItem      __typename    }    __typename  }}fragment BasicFeedItem on FeedsConnection {  pageInfo {    loadMoreKey    hasNextPage    __typename  }  nodes {    ... on ReadSplitBar {      id      type      text      __typename    }    ... on MessageEssential {      ...FeedMessageFragment      __typename    }    ... on UserAction {      id      type      action      actionTime      ... on UserFollowAction {        users {          ...TinyUserFragment          ...TinyUserFragment          ...TinyUserFragment          ...TinyUserFragment          ...TinyUserFragment          ...TinyUserFragment          ...TinyUserFragment          ...TinyUserFragment          __typename        }        allTargetUsers {          ...TinyUserFragment          following          statsCount {            followedCount            __typename          }          ...TinyUserFragment          ...TinyUserFragment          ...TinyUserFragment          ...TinyUserFragment          ...TinyUserFragment          ...TinyUserFragment          ...TinyUserFragment          __typename        }        __typename      }      ... on UserRespectAction {        users {          ...TinyUserFragment          ...TinyUserFragment          ...TinyUserFragment          ...TinyUserFragment          ...TinyUserFragment          ...TinyUserFragment          ...TinyUserFragment          ...TinyUserFragment          __typename        }        targetUsers {          ...TinyUserFragment          ...TinyUserFragment          ...TinyUserFragment          ...TinyUserFragment          ...TinyUserFragment          ...TinyUserFragment          ...TinyUserFragment          ...TinyUserFragment          __typename        }        content        __typename      }      __typename    }    __typename  }  __typename}fragment FeedMessageFragment on MessageEssential {  ...EssentialFragment  ... on OriginalPost {    ...LikeableFragment    ...CommentableFragment    ...RootMessageFragment    ...UserPostFragment    ...MessageInfoFragment    pinned {      personalUpdate      __typename    }    __typename  }  ... on Repost {    ...LikeableFragment    ...CommentableFragment    ...UserPostFragment    ...RepostFragment    pinned {      personalUpdate      __typename    }    __typename  }  ... on Question {    ...UserPostFragment    __typename  }  ... on OfficialMessage {    ...LikeableFragment    ...CommentableFragment    ...MessageInfoFragment    ...RootMessageFragment    __typename  }  __typename}fragment EssentialFragment on MessageEssential {  id  type  content  shareCount  repostCount  createdAt  collected  pictures {    format    watermarkPicUrl    picUrl    thumbnailUrl    smallPicUrl    width    height    __typename  }  urlsInText {    url    originalUrl    title    __typename  }  __typename}fragment LikeableFragment on LikeableMessage {  liked  likeCount  __typename}fragment CommentableFragment on CommentableMessage {  commentCount  __typename}fragment RootMessageFragment on RootMessage {  topic {    id    content    __typename  }  __typename}fragment UserPostFragment on MessageUserPost {  readTrackInfo  user {    ...TinyUserFragment    __typename  }  __typename}fragment TinyUserFragment on UserInfo {  avatarImage {    thumbnailUrl    smallPicUrl    picUrl    __typename  }  isSponsor  username  screenName  briefIntro  __typename}fragment MessageInfoFragment on MessageInfo {  video {    title    type    image {      picUrl      __typename    }    __typename  }  linkInfo {    originalLinkUrl    linkUrl    title    pictureUrl    linkIcon    audio {      title      type      image {        thumbnailUrl        picUrl        __typename      }      author      __typename    }    video {      title      type      image {        picUrl        __typename      }      __typename    }    __typename  }  __typename}fragment RepostFragment on Repost {  target {    ...RepostTargetFragment    __typename  }  targetType  __typename}fragment RepostTargetFragment on RepostTarget {  ... on OriginalPost {    id    type    content    pictures {      thumbnailUrl      __typename    }    topic {      id      content      __typename    }    user {      ...TinyUserFragment      __typename    }    __typename  }  ... on Repost {    id    type    content    pictures {      thumbnailUrl      __typename    }    user {      ...TinyUserFragment      __typename    }    __typename  }  ... on Question {    id    type    content    pictures {      thumbnailUrl      __typename    }    user {      ...TinyUserFragment      __typename    }    __typename  }  ... on Answer {    id    type    content    pictures {      thumbnailUrl      __typename    }    user {      ...TinyUserFragment      __typename    }    __typename  }  ... on OfficialMessage {    id    type    content    pictures {      thumbnailUrl      __typename    }    __typename  }  ... on DeletedRepostTarget {    status    __typename  }  __typename}"

	url := "https://web-api.okjike.com/api/graphql"

	json, _ := jsoni.MarshalToString(data)

	response, err := util.PostRaw(url, headers, json)

	if err != nil {
		logger.Fatalf("Jike GetUserTimeline err: %v", err)

		return nil, err
	}

	parsedJson, err := parser.Parse(string(response.Body()))

	parsedObject := parsedJson.GetArray("data", "userProfile", "feeds", "nodes")

	result := make([]types.TimelineStruct, len(parsedObject))

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
