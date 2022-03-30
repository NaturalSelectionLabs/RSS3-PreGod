package router

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/handler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/gin-gonic/gin"
)

func InitializeRouter() (router *gin.Engine) {
	if config.Config.HubServer.RunMode == "debug" {
		router = gin.Default()
	} else {
		router = gin.New()
	}

	// Latest version API
	apiRouter := router.Group(fmt.Sprintf("/%s", protocol.Version))
	apiRouter.Use(middleware.Instance())
	{
		apiRouter.GET("/:instance", handler.GetIndexHandlerFunc)
		apiRouter.GET("/:instance/profiles")
		apiRouter.GET("/:instance/links")
		apiRouter.GET("/:instance/backlinks")
	}

	// Older version API

	return router
}
