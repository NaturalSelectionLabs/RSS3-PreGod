package router

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/api"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/handler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/gin-gonic/gin"
)

func Initialize() (router *gin.Engine) {
	if config.Config.Hub.Server.RunMode == "debug" {
		router = gin.Default()
	} else {
		router = gin.New()
	}

	// Response wrapper
	router.Use(middleware.Wrapper())

	router.NoRoute(api.NoRouterHandlerFunc)
	router.NoMethod(api.NoMethodHandlerFunc)
	router.GET("/", api.GetIndexHandlerFunc)

	// Latest version API
	apiRouter := router.Group(fmt.Sprintf("/%s", protocol.Version))
	{
		instanceMiddleware := middleware.Instance()

		apiRouter.Use(middleware.ListLimit())

		apiRouter.GET("/:instance", instanceMiddleware, handler.GetInstanceHandlerFunc)
		apiRouter.GET("/:instance/profiles", instanceMiddleware, handler.GetProfileListHandlerFunc)
		apiRouter.GET("/:instance/links", instanceMiddleware, handler.GetLinkListHandlerFunc)
		apiRouter.GET("/:instance/backlinks", instanceMiddleware, handler.GetBackLinkListHandlerFunc)
		apiRouter.GET("/:instance/assets", instanceMiddleware, handler.GetAssetListHandlerFunc)
		apiRouter.GET("/:instance/notes", instanceMiddleware, handler.GetNoteListHandlerFunc)
		apiRouter.POST("/notes", handler.BatchGetNoteListHandlerFunc)
	}

	// Older version API

	return router
}
