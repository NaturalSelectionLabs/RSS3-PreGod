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
	return cache.ConstructKey("httpx", method, url, data)
}

func getCache(url, method, data string) (string, bool) {
	key := getCacheKey(method, url, data)

	if str, err := cache.GetRaw(context.Background(), key); err != nil {
		if err != cache.CacheMissedError {
			logger.Errorf("Error while getting cache: %v", err)
		}

		return "", false
	} else {
		return str, true
	}
}

func setCache(url, method, data, response string) error {
	key := getCacheKey(method, url, data)

	if err := cache.SetRaw(context.Background(), key, response, expiredTime); err != nil {
		return err
	}

	return nil
}
