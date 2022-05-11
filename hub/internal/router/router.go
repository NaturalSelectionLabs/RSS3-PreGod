package router

import (
	"fmt"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/api"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/handler"
	middlewarex "github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Initialize() (router *gin.Engine) {
	if config.Config.Hub.Server.RunMode == "debug" {
		router = gin.Default()
	} else {
		router = gin.New()
	}

	if config.Config.Sentry.DSN != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:        config.Config.Sentry.DSN,
			ServerName: config.Config.Sentry.ServerName,
		}); err != nil {
			panic(err)
		}

		router.Use(sentrygin.New(sentrygin.Options{}))
	}

	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	// Response wrapper
	router.Use(middlewarex.Wrapper())

	router.NoRoute(api.NoRouterHandlerFunc)
	router.NoMethod(api.NoMethodHandlerFunc)
	router.GET("/", api.GetIndexHandlerFunc)

	// Latest version API
	apiRouter := router.Group(fmt.Sprintf("/%s", protocol.Version))
	{
		instanceMiddleware := middlewarex.Instance()

		apiRouter.Use(middlewarex.ListLimit())

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
