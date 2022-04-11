package moralis

import (
	"fmt"
	"log"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/httpx"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	jsoniter "github.com/json-iterator/go"
	lop "github.com/samber/lo/parallel"
	"github.com/valyala/fastjson"
)

const endpoint = "https://deep-index.moralis.io"

var (
	jsoni  = jsoniter.ConfigCompatibleWithStandardLibrary
	parser fastjson.Parser
)

func requestMoralisApi(url string, apiKey string) ([]byte, error) {
	var headers = map[string]string{
		"accept":    "application/json",
		"X-API-Key": apiKey,
	}

	response, err := httpx.Get(url, headers)
	if err != nil {
		logger.Errorf("http get error: [%v]", err)

		return nil, err
	}

	return response, nil
}

func GetNFTs(userAddress string, chainType ChainType, apiKey string) (NFTResult, error) {
	// Gets all NFT items of user
	url := fmt.Sprintf("%s/api/v2/%s/nft?chain=%s&format=decimal",
		endpoint, userAddress, chainType)

	response, err := requestMoralisApi(url, apiKey)

	if err != nil {
		return NFTResult{}, err
	}

	res := new(NFTResult)

	err = jsoni.Unmarshal(response, &res)
	if err != nil {
		return NFTResult{}, err
	}

	lop.ForEach(res.Result, func(item NFTItem, i int) {
		if item.MetaData == "" {
			if metadataRes, err := httpx.Get(item.TokenURI, nil); err != nil {
				logger.Warnf("http get nft metadata error: [%v]", err)
			} else {
				res.Result[i].MetaData = string(metadataRes)
			}
		}
	})

	return *res, nil
}

func GetNFTTransfers(userAddress string, chainType ChainType, blockHeight int, apiKey string) (NFTTransferResult, error) {
	// Gets all NFT transfers of user
	url := fmt.Sprintf("%s/api/v2/%s/nft/transfers?chain=%s&from_block=%d&format=decimal&direction=both",
		endpoint, userAddress, chainType, blockHeight)
	response, err := requestMoralisApi(url, apiKey)

	log.Println(string(response))

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

func GetLogs(fromBlock int64, toBlock int64, address string, topic string, chainType ChainType, apiKey string) (*GetLogsResult, error) {
	url := fmt.Sprintf("%s/api/v2/%s/logs?chain=%s&from_block=%d&to_block=%d&topic0=%s",
		endpoint, address, string(chainType), fromBlock, toBlock, topic)
	response, err := requestMoralisApi(url, apiKey)

	if err != nil {
		return nil, err
	}

	res := new(GetLogsResult)

	err = jsoni.Unmarshal(response, &res)
	if err != nil {
		logger.Errorf("unmarshal error: [%v]", err)

		return nil, err
	}

	return res, nil
}

// Gets all NFT items of user
func GetNFTByContract(userAddress string, contactAddress string, chainType ChainType, apiKey string) (NFTResult, error) {
	// this function is used by ENS indexer.
	url := fmt.Sprintf("%s/api/v2/%s/nft?chain=%s&format=decimal&token_addresses=%s",
		endpoint, userAddress, chainType, contactAddress)

	response, err := requestMoralisApi(url, apiKey)

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

// GetTxByToken is used by ENS indexer
func GetTxByToken(tokenAddress string, tokenId string, chainType ChainType, apiKey string) (NFTTransferItem, error) {
	url := fmt.Sprintf("%s/api/v2/nft/%s/%s/transfers?chain=%s&format=decimal&limit=1",
		endpoint, tokenAddress, tokenId, chainType)
	response, err := requestMoralisApi(url, apiKey)

	if err != nil {
		return NFTTransferItem{}, err
	}

	res := new(NFTTransferItem)

	parsedJson, err := parser.Parse(string(response))

	jsoni.UnmarshalFromString(string(parsedJson.GetObject("result", "0").MarshalTo(nil)), &res)

	if err != nil {
		logger.Errorf("GetTxByToken: %v", err)
	}

	return *res, nil
}

func GetMetadataByToken(tokenAddress string, tokenId string, chainType ChainType, apiKey string) (NFTItem, error) {
	url := fmt.Sprintf("%s/api/v2/nft/%s/%s?chain=%s&format=decimal&limit=1",
		endpoint, tokenAddress, tokenId, chainType)
	response, err := requestMoralisApi(url, apiKey)

	if err != nil {
		return NFTItem{}, err
	}

	res := new(NFTItem)

	err = jsoni.Unmarshal(response, &res)
	if err != nil {
		return NFTItem{}, nil
	}

	return *res, nil
}
