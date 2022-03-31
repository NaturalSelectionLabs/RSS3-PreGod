package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetIndexHandlerFunc(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func NoRouterHandlerFunc(c *gin.Context) {
	_ = c.Error(ErrorNoRouter)
}

func NoMethodHandlerFunc(c *gin.Context) {
	_ = c.Error(ErrorNoRouter)
}
