package misskey

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/httpx"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fastjson"
)

var (
	jsoni  = jsoniter.ConfigCompatibleWithStandardLibrary
	parser fastjson.Parser
)

func GetUserShow(accountInfo []string) (*UserShow, error) {
	url := "https://" + accountInfo[1] + "/api/users/show"
	logger.Infof("url: %s", url)

	username := fmt.Sprintf(`{"username":"%s"}`, accountInfo[0])

	response, requestErr := httpx.Post(url, nil, username)

	if requestErr != nil {
		return nil, requestErr
	}

	parsedJson, parseErr := parser.Parse(string(response))

	if parseErr != nil || parsedJson == nil {
		return nil, requestErr
	}

	errorObj := parsedJson.Get("error")

	if errorObj != nil {
		errorMsg := string(errorObj.GetStringBytes("message"))

		return nil, fmt.Errorf("Get misskey userinfo error: %s", errorMsg)
	}

	userShow := new(UserShow)
	userShow.Id = string(parsedJson.GetStringBytes("id"))
	userShow.Bios = append(userShow.Bios, string(parsedJson.GetStringBytes("description")))
	fields := parsedJson.GetArray("fields")

	for _, field := range fields {
		userShow.Bios = append(userShow.Bios, string(field.GetStringBytes("value")))
	}

	return userShow, nil
}

func GetUserNoteList(address string, count int, since time.Time) ([]Note, error) {
	accountInfo, err := formatUserAccount(address)

	if err != nil {
		return nil, err
	}

	userShow, getUserIdErr := GetUserShow(accountInfo)

	if getUserIdErr != nil {
		return nil, getUserIdErr
	}

	url := "https://" + accountInfo[1] + "/api/users/notes"

	request := new(TimelineRequest)

	request.UserId = userShow.Id
	request.Limit = count
	request.SinceDate = since.Unix() * 1000
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

	var noteList = make([]Note, 0, 10)

	for _, note := range parsedObject {
		ns := new(Note)

		ns.Summary = string(note.GetStringBytes("text"))
		formatContent(note, ns, accountInfo[1])

		ns.Id = string(note.GetStringBytes("id"))
		ns.Author = string(note.GetStringBytes("userId"))
		ns.Link = fmt.Sprintf("https://%s/notes/%s", accountInfo[1], ns.Id)

		t, timeErr := time.Parse(time.RFC3339, string(note.GetStringBytes("createdAt")))

		if timeErr != nil {
			return nil, timeErr
		}

		ns.CreatedAt = t

		noteList = append(noteList, *ns)
	}

	return noteList, nil
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

	renoteId := string(note.GetStringBytes("renoteId"))

	// format renote if any
	if renoteId != "null" {
		renoteUser := string(note.GetStringBytes("renote", "user", "username"))

		renoteText := string(note.GetStringBytes("renote", "text"))

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
		name := string(emoji.GetStringBytes("name"))
		url := string(emoji.GetStringBytes("url"))

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
		_type := string(image.GetStringBytes("type"))

		if strings.HasPrefix(_type, "image/") {
			url := string(image.GetStringBytes("url"))

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
