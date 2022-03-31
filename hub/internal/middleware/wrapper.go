package middleware

import (
	"net/http"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/api"
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
			code := api.ErrorToCode(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, &WrapperResponse{
				Code:  code,
				Error: api.CodeToError(code).Error(),
			})
		}
	}
}
