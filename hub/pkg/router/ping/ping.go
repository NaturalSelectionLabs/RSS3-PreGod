package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string `json:"message"`
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Message: "pong",
	})
}
