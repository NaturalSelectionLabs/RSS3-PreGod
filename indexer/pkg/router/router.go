package router

import (
	"net/http"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/router/api"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/router/monitor"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	// Apply middlewares
	r.Use(gin.Recovery())

	if config.Config.Sentry.DSN != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:        config.Config.Sentry.DSN,
			ServerName: config.Config.Sentry.ServerName,
		}); err != nil {
			panic(err)
		}

		r.Use(sentrygin.New(sentrygin.Options{
			Repanic: true,
		}))
	}

	// === Error handler ===
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "not found",
		})
	})

	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"message": "method not allowed",
		})
	})

	r.GET("/item", api.GetItemHandlerFunc)
	r.GET("/bio", api.GetBioHandlerFunc)
	r.GET("/debug/statsviz/*filepath", monitor.Statsviz)

	return r
}
