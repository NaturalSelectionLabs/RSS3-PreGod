package util

import "github.com/go-resty/resty/v2"

func Get(url string, headers map[string]string) ([]byte, error) {
	// Create a Resty Client
	client := resty.New()
	request := client.R().EnableTrace().SetHeaders(headers)

	// Get url
	resp, err := request.Get(url)

	return resp.Body(), err
}

func Post(url string, headers map[string]string, data string) ([]byte, error) {
	// Create a Resty Client
	client := resty.New()
	request := client.R().EnableTrace().SetHeaders(headers).SetBody(data)

	// Get url
	resp, err := request.Post(url)

	return resp.Body(), err
}

var keyOffset map[string]int

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
