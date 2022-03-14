package misskey

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util/httpx"
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

	response, requestErr := httpx.Post(url, nil, username)

	if requestErr != nil {
		return "", requestErr
	}

	parsedJson, parseErr := parser.Parse(string(response))

	if parseErr != nil {
		return "", requestErr
	}

	return util.TrimQuote(parsedJson.Get("id").String()), nil
}

func GetUserNoteList(address string, count int, tsp time.Time) ([]Note, error) {
	accountInfo, err := formatUserAccount(address)

	if err == nil {
		userId, getUserIdErr := GetUserId(accountInfo)

		if getUserIdErr != nil {
			return nil, getUserIdErr
		}

		url := "https://" + accountInfo[1] + "/api/users/notes"

		request := new(TimelineRequest)

		request.UserId = userId
		request.Limit = count
		request.UntilDate = tsp.Unix() * 1000
		request.ExcludeNsfw = true
		request.Renote = true
		request.IncludeReplies = false

		json, _ := jsoni.MarshalToString(request)

		response, requestErr := httpx.Post(url, nil, json)

		if requestErr != nil {
			return nil, requestErr
		}

		parsedJson, parseErr := parser.Parse(string(response))

		if parseErr != nil {
			return nil, parseErr
		}

		parsedObject := parsedJson.GetArray()

		var noteList []Note

		for _, note := range parsedObject {
			ns := new(Note)

			ns.Summary = util.TrimQuote(note.Get("text").String())
			formatContent(note, ns, accountInfo[1])

			ns.Id = util.TrimQuote(note.Get("id").String())
			ns.Author = util.TrimQuote(note.Get("userId").String())
			ns.Link = fmt.Sprintf("https://%s/notes/%s", accountInfo[1], ns.Id)

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

func formatContent(note *fastjson.Value, ns *Note, instance string) {
	// add emojis into text
	if len(note.GetArray("emojis")) > 0 {
		formatEmoji(note.GetArray("emojis"), ns)
	}

	// add images into text
	if len(note.GetArray("files")) > 0 {
		formatImage(note.GetArray("files"), ns)
	}

	renoteId := note.Get("renoteId").String()

	// format renote if any
	if renoteId != "null" {
		renoteUser := util.TrimQuote(note.Get("renote", "user", "username").String())

		renoteText := util.TrimQuote(note.Get("renote", "text").String())

		ns.Summary = fmt.Sprintf("%s Renote @%s: %s", ns.Summary, renoteUser, renoteText)

		formatContent(note.Get("renote"), ns, instance)

		quoteText := *model.NewAttachment(renoteText, nil, "text/plain", "quote_text", 0, time.Now())

		address := fmt.Sprintf("https://%s/@%s/%s", instance, renoteUser, renoteId)

		quoteAddress := *model.NewAttachment(address, nil, "text/uri-list", "quote_address", 0, time.Now())

		ns.Attachments = append(ns.Attachments, quoteText, quoteAddress)
	}
}

func formatEmoji(emojiList []*fastjson.Value, ns *Note) {
	for _, emoji := range emojiList {
		name := util.TrimQuote(emoji.Get("name").String())
		url := util.TrimQuote(emoji.Get("url").String())

		ns.Summary = strings.Replace(ns.Summary, name, fmt.Sprintf("<img class=\"emoji\" src=\"%s\" alt=\":%s:\">", url, name), -1)

		content := fmt.Sprintf("{\"name\":\"%s\",\"url\":\"%s\"}", name, url)

		attachment := *model.NewAttachment(content, nil, "text/json", "emojis", 0, time.Now())

		ns.Attachments = append(ns.Attachments, attachment)
	}
}

func formatImage(imageList []*fastjson.Value, ns *Note) {
	var mime string

	var sizeInBytes = 0

	for _, image := range imageList {
		_type := util.TrimQuote(image.Get("type").String())

		if strings.HasPrefix(_type, "image/") {
			url := util.TrimQuote(image.Get("url").String())

			ns.Summary += fmt.Sprintf("<img class=\"media\" src=\"%s\">", url)

			res, err := httpx.Head(url)

			if err == nil {
				sizeInBytes, _ = strconv.Atoi(res.Get("Content-Length"))
				mime = res.Get("Content-Type")
			}

			attachment := *model.NewAttachment(url, nil, mime, "quote_file", sizeInBytes, time.Now())

			ns.Attachments = append(ns.Attachments, attachment)
		}
	}
}

// returns [username, instance]
func formatUserAccount(address string) ([]string, error) {
	res := strings.Split(address, "@")

	if len(res) < 2 {
		err := fmt.Errorf("invalid misskey address: %s", address)
		logger.Errorf("%v", err)

		return nil, err
	}

	return res, nil
}
