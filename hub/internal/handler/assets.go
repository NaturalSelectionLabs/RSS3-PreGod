package handler

import (
	"fmt"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetAssetListHandlerFunc(c *gin.Context) {
	instance, err := middleware.GetPlatformInstance(c)
	if err != nil {
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
