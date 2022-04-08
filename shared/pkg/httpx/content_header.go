package httpx

import (
	"context"
	"strconv"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/cache"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
)

// returns required fields for an Attachment
func GetContentHeader(url string) (*ContentHeader, error) {
	ckey := getCacheKey(url)

	var contentHeader = &ContentHeader{
		SizeInByte: 0,
		MIMEType:   "",
	}

	if err := cache.Get(context.Background(), ckey, contentHeader); err != nil {
		if err != cache.Nil {
			logger.Error(err)

			return nil, err
		}
	} else {
		return contentHeader, nil
	}

	// if there is no cache, get the content header from the url and store it in cache
	// as this information is not very important, we do it in a separate goroutine
	// so that we don't block the main thread
	go func() {
		res := new(ContentHeader)

		header, err := Head(url)

		if err != nil {
			logger.Errorf("cannot read content type of url: %s. error is : %v", url, err)
		}

		contentLength := header.Get("Content-Length")

		if contentLength != "" {
			if sizeInBytes, err := strconv.Atoi(contentLength); err != nil {
				logger.Errorf("cannot convert content length of url: %s to int: %s. error is : %v", url, contentLength, err)
			} else {
				res.SizeInByte = sizeInBytes
			}
		} else {
			res.SizeInByte = 0
		}

		res.MIMEType = header.Get("Content-Type")

		// store in cache
		if err := cache.Set(context.Background(), ckey, res, 0); err != nil {
			logger.Error(err)
		}
	}()

	// return a fake "nil" content header
	return contentHeader, nil
}

// TODO: should we manage cache keys in the cache package?
func getCacheKey(url string) string {
	return "httpx:content_header:" + url
}
