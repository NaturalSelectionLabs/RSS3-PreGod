package httpx

import (
	"context"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/cache"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
)

var expiredTime = 5 * time.Minute

func getCacheKey(url string, data string) string {
	return cache.ConstructKey(url, data)
}

func getCache(url string, data string) (string, bool) {
	key := getCacheKey(url, data)

	var response string

	if err := cache.Get(context.Background(), key, response); err != nil {
		if err != cache.CacheMissedError {
			logger.Errorf("Error while getting cache: %v", err)
		}

		return "", false
	}

	return response, true
}

func setCache(url string, data string, response string) error {
	key := getCacheKey(url, data)

	if err := cache.Set(context.Background(), key, response, expiredTime); err != nil {
		return err
	}

	return nil
}
