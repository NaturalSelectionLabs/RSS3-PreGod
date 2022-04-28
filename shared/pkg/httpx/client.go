package httpx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/go-resty/resty/v2"
	"github.com/vincent-petithory/dataurl"
)

type ContentHeader struct {
	MIMEType   string
	SizeInByte int
}

type Response struct {
	Body   []byte
	Header http.Header
}

func NewResponse() *Response {
	return &Response{
		Body:   []byte{},
		Header: http.Header{},
	}
}

func NoCacheGet(url string, headers map[string]string) (*Response, error) {
	return get(url, headers, false)
}

func Get(url string, headers map[string]string) (*Response, error) {
	return get(url, headers, true)
}

func get(url string, headers map[string]string, useCache bool) (*Response, error) {
	resp := NewResponse()

	if useCache {
		// get from cache fist
		cacheResp, ok := getCache(url, methodGet, "")
		if ok {
			resp.Body = []byte(cacheResp)

			return resp, nil
		}
	}

	// nolint: nestif // should be nested if
	if strings.HasPrefix(url, "data:") {
		if strings.Contains(url, "base64") {
			dataUrl, err := dataurl.DecodeString(url)
			if err != nil {
				return nil, err
			}

			resp.Body = dataUrl.Data
		} else {
			// normal data url
			strArr := strings.Split(url, ",")
			if len(strArr) != 2 {
				return nil, fmt.Errorf("invalid data url: %s", url)
			}

			resp.Body = []byte(strArr[1])
		}
	} else {
		client := getClient()

		if headers != nil {
			SetCommonHeader(headers)
		}

		request := client.R().EnableTrace().SetHeaders(headers)

		urlResp, err := request.Get(url)
		if err != nil {
			return nil, err
		}

		if urlResp.StatusCode() != 200 {
			return nil, fmt.Errorf("StatusCode [%d]", urlResp.StatusCode())
		}

		resp.Body = urlResp.Body()
		resp.Header = urlResp.Header()
	}

	if useCache {
		if cacheErr := setCache(url, methodGet, "", string(resp.Body)); cacheErr != nil {
			logger.Errorf("Failed to set cache for url [%s]. err: %+v", url, cacheErr)
		}
	}

	return resp, nil
}

func NoCachePost(url string, headers map[string]string, data string) (*Response, error) {
	return post(url, headers, data, false)
}

func Post(url string, headers map[string]string, data string) (*Response, error) {
	return post(url, headers, data, true)
}

func post(url string, headers map[string]string, data string, useCache bool) (*Response, error) {
	resp := NewResponse()

	// get from cache fist
	cacheResp, ok := getCache(url, methodPost, "")
	if ok {
		resp.Body = []byte(cacheResp)

		return resp, nil
	}

	client := getClient()

	if headers != nil {
		SetCommonHeader(headers)
	}

	request := client.R().EnableTrace().SetHeaders(headers).SetBody(data)

	// Post url
	urlResp, err := request.Post(url)

	if cacheErr := setCache(url, methodPost, data, string(urlResp.Body())); cacheErr != nil {
		logger.Errorf("Failed to set cache for url [%s]. err: %+v", url, cacheErr)
	}

	resp.Body = urlResp.Body()
	resp.Header = urlResp.Header()

	return resp, err
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
	// get from cache fist
	response, ok := getCache(url, methodPost, "")
	if ok {
		jsonBytes := []byte(response)

		var header http.Header
		if err := json.Unmarshal(jsonBytes, &header); err != nil {
			return nil, err
		}

		return header, nil
	}

	client := getClient()

	headers := make(map[string]string)

	SetCommonHeader(headers)

	request := client.R().EnableTrace()

	resp, err := request.Head(url)

	// set cache
	headerMap := map[string][]string(resp.Header())

	jsonBytes, jsonErr := json.Marshal(headerMap)
	if jsonErr != nil {
		logger.Errorf("Failed to marshal header map. err: %+v", jsonErr)
	}

	if cacheErr := setCache(url, methodHead, "", string(jsonBytes)); cacheErr != nil {
		logger.Errorf("Failed to set cache for url [%s]. err: %+v", url, cacheErr)
	}

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
