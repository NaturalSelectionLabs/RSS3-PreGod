package middleware

import (
	"fmt"
	"strings"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/api"
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
			_ = c.Error(api.ErrorInvalidParams)
			c.Abort()

			return
		}

		request.Instance = strings.ToLower(request.Instance)

		instance, err := rss3uri.ParseInstance(request.Instance)
		if err != nil {
			api.SetError(c, api.ErrorInvalidParams, err)
			c.Abort()

			return
		}

		c.Set(KeyInstance, instance)
	}
}

func GetInstance(c *gin.Context) (rss3uri.Instance, error) {
	value, exists := c.Get(KeyInstance)
	if !exists {
		return nil, fmt.Errorf("instance not found")
	}

	instance, ok := value.(rss3uri.Instance)
	if !ok {
		return nil, fmt.Errorf("instance not found")
	}

	return instance, nil
}

func GetPlatformInstance(c *gin.Context) (*rss3uri.PlatformInstance, error) {
	instance, err := GetInstance(c)
	if err != nil {
		return nil, err
	}

	platformInstance, ok := instance.(*rss3uri.PlatformInstance)
	if !ok {
		return nil, fmt.Errorf("instance not found")
	}

	return platformInstance, nil
}

func GetNetworkInstance(c *gin.Context) (*rss3uri.NetworkInstance, error) {
	instance, err := GetInstance(c)
	if err != nil {
		return nil, err
	}

	networkInstance, ok := instance.(*rss3uri.NetworkInstance)
	if !ok {
		return nil, fmt.Errorf("instance not found")
	}

	return networkInstance, nil
}
