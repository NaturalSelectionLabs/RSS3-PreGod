package xscan

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/httpx"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fastjson"
)

var jsoni = jsoniter.ConfigCompatibleWithStandardLibrary

func GetApiKey(networkId constants.NetworkID) string {
	var err error
	if err = config.Setup(); err != nil {
		return ""
	}

	var apiKey string
	if networkId == constants.NetworkIDEthereumMainnet {
		apiKey, err = jsoni.MarshalToString(config.Config.Indexer.EtherScan.ApiKey)
	} else if networkId == constants.NetworkIDPolygon {
		apiKey, err = jsoni.MarshalToString(config.Config.Indexer.PolygonScan.ApiKey)
	}

	if err != nil {
		return ""
	}

	return strings.Trim(apiKey, "\"")
}

func GetLatestBlockHeight(networkId constants.NetworkID) (int64, error) {
	apiKey := GetApiKey(networkId)
	if apiKey == "" {
		return 0, fmt.Errorf("no api key")
	}

	var url string
	if networkId == constants.NetworkIDEthereumMainnet {
		url = "https://api.etherscan.io/api/?module=proxy&action=eth_blockNumber&apikey=" + apiKey
	} else if networkId == constants.NetworkIDPolygon {
		url = "https://api.polygonscan.com/api/?module=proxy&action=eth_blockNumber&apikey=" + apiKey
	}

	response, err := httpx.Get(url, nil)
	if err != nil {
		return 0, err
	}

	var parser fastjson.Parser
	parsedJson, parseErr := parser.Parse(string(response))

	if parseErr != nil {
		return 0, parseErr
	}

	msg := string(parsedJson.GetStringBytes("message"))
	result := string(parsedJson.GetStringBytes("result"))

	if msg == "NOTOK" {
		return 0, fmt.Errorf("api error, %s", result)
	}

	blockHeight, err := strconv.ParseUint(result[2:], 16, 64)
	if err != nil {
		return 0, err
	}

	return int64(blockHeight), nil
}
