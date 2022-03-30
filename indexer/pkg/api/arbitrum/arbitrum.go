package arbitrum

import (
	"fmt"
	"strings"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/httpx"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fastjson"
)

const endpoint = "https://api.arbiscan.io"

var jsoni = jsoniter.ConfigCompatibleWithStandardLibrary

func GetApiKey() string {
	if err := config.Setup(); err != nil {
		return ""
	}

	apiKey, err := jsoni.MarshalToString(config.Config.Indexer.Aribtrum.ApiKey)
	if err != nil {
		return ""
	}

	return strings.Trim(apiKey, "\"")
}

func GetNFTTxs(owner string) ([]byte, error) {
	apiKey := GetApiKey()
	url := fmt.Sprintf(
		"%s/api?module=account&action=tokennfttx&address=%s&startblock=0&endblock=999999999&sort=asc&apikey=%s",
		endpoint, owner, apiKey)

	response, err := httpx.Get(url, nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func GetNFTTransfers(owner string) ([]NFTTransferItem, error) {
	response, err := GetNFTTxs(owner)
	if err != nil {
		return nil, err
	}

	var parser fastjson.Parser

	parsedJson, parseErr := parser.Parse(string(response))
	if parseErr != nil {
		return nil, parseErr
	}

	result := make([]NFTTransferItem, 0)

	arrys := parsedJson.GetArray("result")
	for _, v := range arrys {
		var item NFTTransferItem
		item.TokenAddress = string(v.GetStringBytes("contractAddress"))
		item.TokenId = string(v.GetStringBytes("tokenID"))
		item.Name = string(v.GetStringBytes("tokenName"))
		item.Symbol = string(v.GetStringBytes("tokenSymbol"))
		item.From = string(v.GetStringBytes("from"))
		item.To = string(v.GetStringBytes("to"))
		item.TimeStamp = string(v.GetStringBytes("timeStamp"))
		item.Hash = string(v.GetStringBytes("hash"))

		result = append(result, item)
	}

	return result, nil
}

func GetNFTs(owner string) ([]NFTItem, error) {
	response, err := GetNFTTxs(owner)
	if err != nil {
		return nil, err
	}

	var parser fastjson.Parser

	parsedJson, parseErr := parser.Parse(string(response))
	if parseErr != nil {
		return nil, parseErr
	}

	nfts := make(map[string]NFTItem)

	arrys := parsedJson.GetArray("result")
	for _, v := range arrys {
		var nft NFTItem
		nft.TokenAddress = string(v.GetStringBytes("contractAddress"))
		nft.TokenId = string(v.GetStringBytes("tokenID"))
		nft.Name = string(v.GetStringBytes("tokenName"))
		nft.Symbol = string(v.GetStringBytes("tokenSymbol"))

		from := string(v.GetStringBytes("from"))
		to := string(v.GetStringBytes("to"))

		if to == owner {
			nft.Valid = true
		} else if from == owner {
			nft.Valid = false
		}

		nfts[nft.TokenId] = nft
	}

	result := make([]NFTItem, 0)

	for _, v := range nfts {
		v.TokenURI = GetTokenURI(v.TokenAddress)
		result = append(result, v)
	}

	return result, nil
}

func GetTokenURI(contractAddress string) string {
	// TODO: get tokenURI
	return ""
}
