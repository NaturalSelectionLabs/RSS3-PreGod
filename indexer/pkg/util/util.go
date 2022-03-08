package util

import (
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/go-resty/resty/v2"
)

func Get(url string, headers map[string]string) ([]byte, error) {
	// Create a Resty Client
	client := resty.New()

	if len(config.Config.Network.Proxy) != 0 {
		client.SetProxy(config.Config.Network.Proxy)
	}

	client.SetTimeout(1 * time.Second * 10)

	if headers != nil {
		SetCommonHeader(headers)
	}

	request := client.R().EnableTrace().SetHeaders(headers)

	// Get url
	resp, err := request.Get(url)

	return resp.Body(), err
}

func Post(url string, headers map[string]string, data string) ([]byte, error) {
	// Create a Resty Client
	client := resty.New()

	if len(config.Config.Network.Proxy) != 0 {
		client.SetProxy(config.Config.Network.Proxy)
	}

	client.SetTimeout(1 * time.Second * 10)

	if headers != nil {
		SetCommonHeader(headers)
	}

	request := client.R().EnableTrace().SetHeaders(headers).SetBody(data)

	// Post url
	resp, err := request.Post(url)

	return resp.Body(), err
}

// PostRaw returns raw *resty.Response for Jike
func PostRaw(url string, headers map[string]string, data string) (*resty.Response, error) {
	// Create a Resty Client
	client := resty.New()
	client.SetTimeout(1 * time.Second * 10)

	if headers != nil {
		SetCommonHeader(headers)
	}

	request := client.R().EnableTrace().SetHeaders(headers).SetBody(data)

	// Post url
	resp, err := request.Post(url)

	return resp, err
}

func SetCommonHeader(headers map[string]string) {
	headers["User-Agent"] = "RSS3-PreGod"
}

var keyOffset = make(map[string]int)

func GotKey(strategy string, indexer_id string, keys []string) string {
	if len(strategy) == 0 {
		strategy = "round-robin"
	}

	if len(indexer_id) == 0 {
		indexer_id = "."
	}

	var offset int

	var key string

	if strategy == "first-always" {
		key = "Bearer " + indexer_id
	} else {
		count, ok := keyOffset[indexer_id]

		if !ok {
			keyOffset[indexer_id] = 0
		}

		offset = count % len(keys)
		keyOffset[indexer_id] = count + 1
		key = keys[offset]
	}

	return key
}
