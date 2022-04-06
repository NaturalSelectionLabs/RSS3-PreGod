package handler

import (
	"net/http"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/api"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
)

func GetInstanceHandlerFunc(c *gin.Context) {
	instance, err := middleware.GetPlatformInstance(c)
	if err != nil {
		return
	}

	if err := database.QueryInstance(
		database.DB,
		instance.Identity,
		constants.ProfileSourceIDCrossbell.Int(),
	); err != nil {
		_ = c.Error(api.ErrorDatabase)

		return
	}

	instanceList := protocol.NewInstanceList(instance)

	c.JSON(http.StatusOK, protocol.File{
		Version:    protocol.Version,
		Identifier: rss3uri.New(instance).String(),
		Total:      len(instanceList),
		List:       instanceList,
	})
}
