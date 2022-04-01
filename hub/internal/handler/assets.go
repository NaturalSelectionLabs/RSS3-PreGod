package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/api"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
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
		_ = c.Error(api.ErrorNotFound)

		return
	}

	noteList := make([]protocol.Item, 0)

	if len(noteList) == 0 {
		_ = c.Error(api.ErrorNotFound)

		return
	}

	c.JSON(http.StatusOK, protocol.File{
		Identifier: fmt.Sprintf("%s/assets", rss3uri.New(instance)),
		Total:      len(noteList),
		List:       noteList,
	})
}
