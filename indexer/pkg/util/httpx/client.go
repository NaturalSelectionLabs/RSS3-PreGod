package httpx

import (
	"net/http"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/go-resty/resty/v2"
)

func Get(url string, headers map[string]string) ([]byte, error) {
	// Create a Resty Client
	client := getClient()

	if headers != nil {
		SetCommonHeader(headers)
	}

	request := client.R().EnableTrace().SetHeaders(headers)

	// Get url
	resp, err := request.Get(url)

	return resp.Body(), err
}

func Post(url string, headers map[string]string, data string) ([]byte, error) {
	client := getClient()

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
	client := getClient()

	if headers != nil {
		SetCommonHeader(headers)
	}

	request := client.R().EnableTrace().SetHeaders(headers).SetBody(data)

	// Post url
	resp, err := request.Post(url)

	return resp, err
}

func Head(url string) (http.Header, error) {
	client := getClient()

	headers := make(map[string]string)

	SetCommonHeader(headers)

	request := client.R().EnableTrace()

	resp, err := request.Head(url)

	return resp.Header(), err
}

func getClient() *resty.Client {
	client := resty.New()

	if len(config.Config.Network.Proxy) != 0 {
		client.SetProxy(config.Config.Network.Proxy)
	}

	client.SetTimeout(1 * time.Second * 10)

	return client
}

func SetCommonHeader(headers map[string]string) {
	headers["User-Agent"] = "RSS3-PreGod"
}
