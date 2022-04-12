package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type WrapperResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func Wrapper() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if err := c.Errors.Last(); err != nil {
			c.AbortWithStatusJSON(http.StatusOK, &WrapperResponse{
				Code:  c.GetInt("code"),
				Error: err.Error(),
			})
		}
	}
}
