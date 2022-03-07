package router

import (
	"net/http"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/router/api"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/router/doc"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/router/monitor"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/router/ping"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/status"
	"github.com/gin-gonic/gin"
)

const (
	API_PATH    = "/api"
	API_VERSION = "v0.4.0"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	// Apply middlewares
	r.Use(gin.Recovery())

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, status.Error(status.ErrorNotFound))
	})

	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, status.Error(status.ErrorMethodNotAllowed))
	})

	// === APIs ===

	// instance
	apis := r.Group(API_PATH)
	apis.Use(middleware.Logger())
	{
		// Instance
		apis.GET("/:instance", api.GetInstance)

		// items
		apis.GET("/:instance/:item_type/:item_uuid", api.GetItem)
		apis.GET("/:instance/list/:item_type/:page_index", api.GetItemPagedList)
		//r.GET("/:authority/list/:item_type", api.GetItemList)

		// // links
		// r.GET("/:authority/list/links/:link_type/:page_index", api.GetLinkList)
		// r.GET("/:authority/list/backlinks/:link_type", api.GetBacklinkList)
	}

	// === Monitor ===
	r.GET("/ping", ping.Ping)
	r.GET("/debug/statsviz/*filepath", monitor.Statsviz)

	// === Static ===
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://rss3.io/favicon.ico")
	})

	// === Docs ===
	r.GET("/docs/*any", doc.Doc(API_PATH, API_VERSION))

	return r
}
