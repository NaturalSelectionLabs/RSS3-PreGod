package middleware

import (
	"fmt"
	"net/http"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/status"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/web"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
)

const (
	KeyInstance = "instance"
)

type InstanceUri struct {
	Instance string `uri:"instance" binding:"required"`
}

func Instance() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := InstanceUri{}
		if err := c.ShouldBindUri(&request); err != nil {
			w := web.Gin{C: c}
			w.JSONResponse(http.StatusBadRequest, status.CodeInvalidParams, nil)
			c.Abort()

			return
		}

		instance, err := rss3uri.ParseInstance(request.Instance)
		if err != nil {
			w := web.Gin{C: c}
			w.JSONResponse(http.StatusBadRequest, status.CodeInvalidParams, nil)
			c.Abort()

			return
		}

		c.Set(KeyInstance, instance)
	}
}

func GetPlatformInstance(c *gin.Context) (*rss3uri.PlatformInstance, error) {
	value, exists := c.Get(KeyInstance)
	if !exists {
		return nil, fmt.Errorf("instance not found")
	}

	platformInstance, ok := value.(*rss3uri.PlatformInstance)
	if !ok {
		return nil, fmt.Errorf("instance not found")
	}

	return platformInstance, nil
}
