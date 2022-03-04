package util

import (
	"time"

	"github.com/go-resty/resty/v2"
)

func Get(url string, headers map[string]string) ([]byte, error) {
	// Create a Resty Client
	client := resty.New()
	client.SetTimeout(1 * time.Second * 10)

	if headers == nil {
		headers = make(map[string]string)
	}

	SetCommonHeader(headers)

	request := client.R().EnableTrace().SetHeaders(headers)

	// Get url
	resp, err := request.Get(url)

	return resp.Body(), err
}

func Post(url string, headers map[string]string, data string) ([]byte, error) {
	// Create a Resty Client
	client := resty.New()
	client.SetTimeout(1 * time.Second * 10)

	if headers == nil {
		headers = make(map[string]string)
	}

	SetCommonHeader(headers)

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

	if headers == nil {
		headers = make(map[string]string)
	}

	SetCommonHeader(headers)

	request := client.R().EnableTrace().SetHeaders(headers).SetBody(data)

	// Post url
	resp, err := request.Post(url)

	return resp, err
}

func SetCommonHeader(headers map[string]string) {
	headers["User-Agent"] = "RSS3-PreGod"
}
