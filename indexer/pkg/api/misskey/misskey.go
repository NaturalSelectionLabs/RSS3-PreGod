package misskey

import (
	"fmt"
	"strings"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fastjson"
)

var (
	jsoni  = jsoniter.ConfigCompatibleWithStandardLibrary
	parser fastjson.Parser
)

func GetUserId(accountInfo []string) (string, error) {
	url := "https://" + accountInfo[1] + "/api/users/show"

	username := fmt.Sprintf(`{"username":"%s"}`, accountInfo[0])

	response, requestErr := util.Post(url, nil, username)

	if requestErr != nil {
		return "", requestErr
	}

	parsedJson, parseErr := parser.Parse(string(response))

	if parseErr != nil {
		return "", requestErr
	}

	return util.TrimQuote(parsedJson.Get("id").String()), nil
}

func GetUserNoteList(address string, count int, tsp time.Time) ([]NoteStruct, error) {
	accountInfo, err := formatUserAccount(address)

	if err == nil {
		userId, getUserIdErr := GetUserId(accountInfo)

		if getUserIdErr != nil {
			return nil, getUserIdErr
		}

		url := "https://" + accountInfo[1] + "/api/users/notes"

		request := new(TimelineRequestStruct)

		request.UserId = userId
		request.Limit = count
		request.UntilDate = tsp.Unix() * 1000
		request.ExcludeNsfw = true
		request.Renote = true
		request.IncludeReplies = false

		json, _ := jsoni.MarshalToString(request)

		response, requestErr := util.Post(url, nil, json)

		if requestErr != nil {
			return nil, requestErr
		}

		parsedJson, parseErr := parser.Parse(string(response))

		if parseErr != nil {
			return nil, parseErr
		}

		parsedObject := parsedJson.GetArray()

		var noteList []NoteStruct

		for _, note := range parsedObject {
			ns := new(NoteStruct)

			ns.Text = util.TrimQuote(note.Get("text").String())
			formatContent(note, ns)

			ns.Id = util.TrimQuote(note.Get("id").String())
			ns.Author = util.TrimQuote(note.Get("userId").String())

			t, timeErr := time.Parse(time.RFC3339, util.TrimQuote(note.Get("createdAt").String()))

			if timeErr != nil {
				return nil, timeErr
			}

			ns.CreatedAt = t

			noteList = append(noteList, *ns)
		}

		return noteList, nil
	}

	return nil, err
}

func formatContent(note *fastjson.Value, ns *NoteStruct) {
	// add emojis into text
	if len(note.GetArray("emojis")) > 0 {
		formatEmoji(note.GetArray("emojis"), ns)
	}

	// add images into text
	if len(note.GetArray("files")) > 0 {
		formatImage(note.GetArray("files"), ns)
	}

	// format renote if any
	if note.Get("renoteId").String() != "null" {
		renoteUser := util.TrimQuote(note.Get("renote", "user", "username").String())

		renoteText := util.TrimQuote(note.Get("renote", "text").String())

		ns.Text = fmt.Sprintf("%s Renote @%s: %s", ns.Text, renoteUser, renoteText)

		formatContent(note.Get("renote"), ns)
	}
}

func formatEmoji(emojiList []*fastjson.Value, ns *NoteStruct) {
	for _, emoji := range emojiList {
		name := util.TrimQuote(emoji.Get("name").String())
		url := util.TrimQuote(emoji.Get("url").String())

		ns.Text = strings.Replace(ns.Text, name, fmt.Sprintf("<img class=\"emoji\" src=\"%s\" alt=\":%s:\">", url, name), -1)
	}
}

func formatImage(imageList []*fastjson.Value, ns *NoteStruct) {
	for _, image := range imageList {
		_type := util.TrimQuote(image.Get("type").String())

		if strings.HasPrefix(_type, "image/") {
			url := util.TrimQuote(image.Get("url").String())

			ns.Text += fmt.Sprintf("<img class=\"media\" src=\"%s\">", url)
		}
	}
}

func formatUserAccount(address string) ([]string, error) {
	res := strings.Split(address, "@")

	if len(res) < 2 {
		err := fmt.Errorf("invalid address: %s", address)
		logger.Errorf("%v", err)

		return nil, err
	}

	return res, nil
}
