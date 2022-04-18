package httpx

import (
	"context"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/cache"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
)

var expiredTime = 5 * time.Minute

var methodGet = "get"
var methodPost = "post"
var methodHead = "head"

func getCacheKey(method, url, data string) string {
	return cache.ConstructKey(method, url, data)
}

func getCache(url, method, data string) (string, bool) {
	key := getCacheKey(method, url, data)

	var response string

	if err := cache.Get(context.Background(), key, response); err != nil {
		if err != cache.CacheMissedError {
			logger.Errorf("Error while getting cache: %v", err)
		}

		return "", false
	}

	return response, true
}

func setCache(url, method, data, response string) error {
	key := getCacheKey(method, url, data)

	if err := cache.Set(context.Background(), key, response, expiredTime); err != nil {
		return err
	}

	return nil
}
