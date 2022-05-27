package middleware

import (
	"strconv"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/api"
	"github.com/gin-gonic/gin"
)

type ListLimitRequest struct {
	// TODO Validator
	Limit *int `form:"limit"`
}

const MaxListLimit = 1000
const DefaultListLimit = 100

func ListLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := ListLimitRequest{}
		if err := c.ShouldBindQuery(&request); err != nil {
			api.SetError(c, api.ErrorInvalidParams, err)
			c.Abort()

			return
		}

		if request.Limit == nil || *request.Limit <= 0 {
			query := c.Request.URL.Query()
			query.Set("limit", strconv.Itoa(DefaultListLimit))
			c.Request.URL.RawQuery = query.Encode()

			return
		}

		if *request.Limit > MaxListLimit {
			query := c.Request.URL.Query()
			query.Set("limit", strconv.Itoa(MaxListLimit))
			c.Request.URL.RawQuery = query.Encode()

			return
		}
	}
}
