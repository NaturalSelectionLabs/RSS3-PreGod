package router

import (
	"net/http"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/router/api"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	// Apply middlewares
	r.Use(gin.Recovery())

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

	return r
}
