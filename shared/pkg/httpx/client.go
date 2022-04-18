package httpx

import (
	"fmt"
	"net/http"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/go-resty/resty/v2"
)

type ContentHeader struct {
	MIMEType   string
	SizeInByte int
}

func Get(url string, headers map[string]string) ([]byte, error) {
	// get from cache fist
	response, ok := getCache(url, "")
	if ok {
		return []byte(response), nil
	}

	client := getClient()

	if headers != nil {
		SetCommonHeader(headers)
	}

	request := client.R().EnableTrace().SetHeaders(headers)

	// Get url
	resp, err := request.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("StatusCode [%d]", resp.StatusCode())
	}

	if cacheErr := setCache(url, "", string(resp.Body())); cacheErr != nil {
		logger.Errorf("Failed to set cache for url [%s]. err: %+v", url, cacheErr)
	}

	return resp.Body(), err
}

func Post(url string, headers map[string]string, data string) ([]byte, error) {
	// get from cache fist
	response, ok := getCache(url, "")
	if ok {
		return []byte(response), nil
	}

	client := getClient()

	if headers != nil {
		SetCommonHeader(headers)
	}

	request := client.R().EnableTrace().SetHeaders(headers).SetBody(data)

	// Post url
	resp, err := request.Post(url)

	if cacheErr := setCache(url, data, string(resp.Body())); cacheErr != nil {
		logger.Errorf("Failed to set cache for url [%s]. err: %+v", url, cacheErr)
	}

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

// TODO: add cache
func Head(url string) (http.Header, error) {
	client := getClient()

	headers := make(map[string]string)

	SetCommonHeader(headers)

	request := client.R().EnableTrace()

	resp, err := request.Head(url)

	return resp.Header(), err
}

var client *resty.Client

func init() {
	client = resty.New()

	if len(config.Config.Network.Proxy) != 0 {
		client.SetProxy(config.Config.Network.Proxy)
	}

	client.SetTimeout(6 * time.Second * 10)
	client.SetRetryCount(3)
}

func getClient() *resty.Client {
	return client
}

func SetCommonHeader(headers map[string]string) {
	headers["User-Agent"] = "RSS3-PreGod"
}
