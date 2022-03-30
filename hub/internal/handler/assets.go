package handler

import (
	"errors"
	"fmt"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type GetAssetListRequest struct {
	Limit         int       `form:"limit"`
	LastTime      time.Time `json:"last_time"`
	Tags          []string  `json:"tags"`
	MimeTypes     []string  `json:"mime_types"`
	ItemSources   []string  `json:"item_sources"`
	LinkSource    string    `json:"link_source"`
	LinkType      string    `json:"link_type"`
	ProfileSource string    `json:"profile_source"`
}

func GetAssetListHandlerFunc(c *gin.Context) {
	instance, err := middleware.GetPlatformInstance(c)
	if err != nil {
		return
	}

	request := GetAssetListRequest{}
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(errors.New("invalid params"))

		return
	}

	noteList := []protocol.Item{
		{
			Identifier:  "",
			DateCreated: time.Time{},
			DateUpdated: time.Time{},
			RelatedURLs: nil,
			Links:       "",
			BackLinks:   "",
			Tags:        nil,
			Authors:     nil,
			Title:       "",
			Summary:     "",
			Attachments: nil,
		},
	}

	c.JSON(http.StatusOK, protocol.File{
		Identifier: fmt.Sprintf("%s/assets", instance.String()),
		Total:      len(noteList),
		List:       noteList,
	})
}
