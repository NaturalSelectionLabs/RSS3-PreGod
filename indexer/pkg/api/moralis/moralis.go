package moralis

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/httpx"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fastjson"
)

const endpoint = "https://deep-index.moralis.io"

var (
	jsoni  = jsoniter.ConfigCompatibleWithStandardLibrary
	parser fastjson.Parser
)

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

	res := new(GetLogsResult)

	err = jsoni.Unmarshal(response, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// this function is used by ENS indexer
func GetNFTByContract(userAddress string, contactAddress string, chainType ChainType, apiKey string) (NFTResult, error) {
	var headers = map[string]string{
		"accept":    "application/json",
		"X-API-Key": apiKey,
	}

	// Gets all NFT items of user
	url := fmt.Sprintf("%s/api/v2/%s/nft?chain=%s&format=decimal&token_addresses=%s",
		endpoint, userAddress, chainType, contactAddress)

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

// this function is used by ENS indexer
func GetTxByToken(tokenAddress string, tokenId string, chainType ChainType, apiKey string) (NFTTransferItem, error) {
	var headers = map[string]string{
		"accept":    "application/json",
		"X-API-Key": apiKey,
	}

	url := fmt.Sprintf("%s/api/v2/nft/%s/%s/transfers?chain=%s&format=decimal&limit=1",
		endpoint, tokenAddress, tokenId, chainType)

	res := new(NFTTransferItem)

	response, err := httpx.Get(url, headers)
	if err != nil {
		return *res, err
	}

	parsedJson, err := parser.Parse(string(response))

	parsedObject := parsedJson.GetArray("result")[0]

	res.BlockTimestamp = string(parsedObject.GetStringBytes("block_timestamp"))
	res.TransactionHash = string(parsedObject.GetStringBytes("transaction_hash"))

	if err != nil {
		logger.Errorf("GetTxByToken: %v", err)
	}

	return *res, nil
}
