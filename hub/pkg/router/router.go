package router

import (
	"net/http"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/router/api"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/router/doc"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/router/ping"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/gin-gonic/gin"
)

const (
	API_PATH    = "/api"
	API_VERSION = "v0.4.0"
)

func InitRouter() (engine *gin.Engine) {
	if config.Config.HubServer.RunMode == "debug" {
		engine = gin.Default()
	} else {
		engine = gin.New()
	}

	// Apply middlewares
	engine.Use(gin.Recovery())
	engine.Use(middleware.Logger())

	// === Error handler ===
	engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "not found",
		})
	})

	engine.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"message": "method not allowed",
		})
	})

	// === APIs ===
	apiRouter := engine.Group(API_PATH)
	apiRouter.Use(middleware.Logger())
	apiRouter.Use(middleware.Instance())
	{
		// rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum
		// Instance File
		apiRouter.GET("/:instance", api.GetIndexHandlerFunc)
		apiRouter.PUT("/:instance", api.PutIndexHandlerFunc)

		// Link List File
		// rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum/list/link/following/1
		apiRouter.GET("/:instance/list/link/:link_type/:page_index", api.GetLinkListHandlerFunc)

		// Back Link List File
		apiRouter.GET("/:instance/list/backlink", api.GetBackLinkListHandlerFunc)

		// Asset List File
		// rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum/list/asset/0
		apiRouter.GET("/:instance/list/asset/:page_index", api.GetAssetListRequestHandlerFunc)

		// Note List File
		// rss3://account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@ethereum/list/note/0
		apiRouter.GET("/:instance/list/note/:page_index", api.GetNoteListRequestHandlerFunc)
	}

	//// === Monitor ===
	engine.GET("/ping", ping.Ping)
	//r.GET("/debug/statsviz/*filepath", monitor.Statsviz)
	//
	//// === Static ===
	engine.GET("/favicon.ico", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://rss3.io/favicon.ico")
	})
	//
	//// === Docs ===
	engine.GET("/docs/*any", doc.Doc(API_PATH, API_VERSION))

	return engine
}
