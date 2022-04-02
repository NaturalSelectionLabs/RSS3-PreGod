package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/api"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/indexer"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type GetNoteListRequest struct {
	Limit         int       `form:"limit"`
	LastTime      time.Time `json:"last_time"`
	Tags          []string  `json:"tags"`
	MimeTypes     []string  `json:"mime_types"`
	ItemSources   []string  `json:"item_sources"`
	LinkSource    string    `json:"link_source"`
	LinkType      string    `json:"link_type"`
	ProfileSource string    `json:"profile_source"`
}

func GetNoteListHandlerFunc(c *gin.Context) {
	instance, err := middleware.GetPlatformInstance(c)
	if err != nil {
		return
	}

	request := GetNoteListRequest{}
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(api.ErrorNotFound)

		return
	}

	client := resty.New()
	// proof=kallydev&platform_id=6&network_id=12&item_type=note
	// TODO
	//profiles, err := database.QueryProfiles(database.DB, "", 1, []int{1})
	//if err != nil {
	//	_ = c.Error(api.ErrorNotFound)
	//
	//	return
	//}

	var indexerResponse indexer.Response

	_, err = client.NewRequest().
		EnableTrace().
		SetQueryParams(map[string]string{
			// TODO
			"proof":       "kallydev",
			"platform_id": "6",
			"network_id":  "12",
			"item_type":   "note",
		}).
		SetResult(&indexerResponse).
		Get(indexer.EndpointItem)
	if err != nil {
		_ = c.Error(err)

		return
	}

	if indexerResponse.Error.Code != 0 {
		_ = c.Error(err)

		return
	}

	noteList := make([]protocol.Item, 0)

	for _, note := range indexerResponse.Data.Note {
		attachments := make([]protocol.ItemAttachment, 0)
		for _, attachment := range note.Attachments {
			attachments = append(attachments, protocol.ItemAttachment{
				Type:        "object",
				Address:     attachment.Address,
				MimeType:    attachment.MimeType,
				SizeInBytes: attachment.SizeInBytes,
			})
		}

		// TODO
		//dateCrated, err := time.Parse("", note.DateCreated)
		//if err != nil {
		//	_ = c.Error(err)
		//
		//	return
		//}

		uri := rss3uri.New(instance)

		authors := make([]string, 0)
		for _, author := range note.Authors {
			authors = append(authors, author)
		}

		noteInstance, err := rss3uri.NewInstance(
			string(constants.ItemTypeNote),
			note.ItemId.Proof,
			constants.NetworkID(note.ItemId.NetworkId).Symbol().String(),
		)
		if err != nil {
			_ = c.Error(err)

			return
		}

		noteList = append(noteList, protocol.Item{
			Identifier: rss3uri.New(noteInstance).String(),
			// TODO
			DateCreated: time.Now(),
			DateUpdated: time.Time{},
			RelatedURLs: nil,
			Links:       fmt.Sprintf("%s/links", uri),
			BackLinks:   fmt.Sprintf("%s/backlinks", uri),
			Tags:        note.Tags,
			Authors:     authors,
			Title:       note.Title,
			Summary:     note.Summary,
			Attachments: attachments,
		})
	}

	if len(noteList) == 0 {
		_ = c.Error(api.ErrorNotFound)

		return
	}

	c.JSON(http.StatusOK, protocol.File{
		Identifier: fmt.Sprintf("%s/notes", rss3uri.New(instance)),
		Total:      len(noteList),
		List:       noteList,
	})
}
