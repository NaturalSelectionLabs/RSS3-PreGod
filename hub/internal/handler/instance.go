package handler

import (
	"net/http"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
)

func GetInstanceHandlerFunc(c *gin.Context) {
	instance, err := middleware.GetPlatformInstance(c)
	if err != nil {
		return
	}

	instanceList := protocol.NewInstanceList(instance)

	c.JSON(http.StatusOK, protocol.File{
		Version: protocol.Version,
		// TODO
		DateUpdated: time.Now(),
		Identifier:  rss3uri.New(instance).String(),
		Total:       len(instanceList),
		List:        instanceList,
	})
}