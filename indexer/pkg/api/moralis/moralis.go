package moralis

import (
	"fmt"
	"strings"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/httpx"
	jsoniter "github.com/json-iterator/go"
)

const endpoint = "https://deep-index.moralis.io"

var jsoni = jsoniter.ConfigCompatibleWithStandardLibrary

func GetApiKey() string {
	if err := config.Setup(); err != nil {
		return ""
	}

	apiKey, err := jsoni.MarshalToString(config.Config.Indexer.Moralis.ApiKey)
	if err != nil {
		return ""
	}

	return strings.Trim(apiKey, "\"")
}

func GetNFTs(userAddress string, chainType ChainType, apiKey string) (NFTResult, error) {
	var headers = map[string]string{
		"accept":    "application/json",
		"X-API-Key": apiKey,
	}

	// Gets all NFT items of user
	url := fmt.Sprintf("%s/api/v2/%s/nft?chain=%s&format=decimal",
		endpoint, userAddress, chainType)

	response, err := httpx.Get(url, headers)
	if err != nil {
		return NFTResult{}, err
	}

	res := new(NFTResult)

	err = jsoni.Unmarshal(response, &res)
	if err != nil {
		return NFTResult{}, err
	}

	return *res, nil
}

func GetNFTTransfers(userAddress string, chainType ChainType, apiKey string) (NFTTransferResult, error) {
	var headers = map[string]string{
		"accept":    "application/json",
		"X-API-Key": apiKey,
	}

	// Gets all NFT transfers of user
	url := fmt.Sprintf("%s/api/v2/%s/nft/transfers?chain=%s&format=decimal&direction=both",
		endpoint, userAddress, chainType)

	response, err := httpx.Get(url, headers)
	if err != nil {
		return NFTTransferResult{}, err
	}

	res := new(NFTTransferResult)

	err = jsoni.Unmarshal(response, &res)
	if err != nil {
		return NFTTransferResult{}, err
	}

	return *res, nil
}

func GetLogs(fromBlock int64, toBlock int64, address string, topic string, chainType string, apiKey string) (*GetLogsResult, error) {
	var headers = map[string]string{
		"accept":    "application/json",
		"X-API-Key": apiKey,
	}

	url := fmt.Sprintf("%s/api/v2/%s/logs?chain=%s&from_block=%d&to_block=%d&topic0=%s",
		endpoint, address, chainType, fromBlock, toBlock, topic)

	response, err := httpx.Get(url, headers)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(response))

	res := new(GetLogsResult)

	err = jsoni.Unmarshal(response, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
