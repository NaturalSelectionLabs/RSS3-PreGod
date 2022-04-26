package arbitrum

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/httpx"
	"github.com/valyala/fastjson"
)

const endpoint = "https://api.arbiscan.io"

func GetNFTTxs(owner string) ([]byte, error) {
	url := fmt.Sprintf(
		"%s/api?module=account&action=tokennfttx&address=%s&startblock=0&endblock=999999999&sort=asc&apikey=%s",
		endpoint, owner, config.Config.Indexer.Aribtrum.ApiKey)

	response, err := httpx.Get(url, nil)
	if err != nil {
		return nil, err
	}

	return response.Body, nil
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

	array := parsedJson.GetArray("result")

	result := make([]NFTTransferItem, 0, len(array))

	for _, v := range array {
		var item NFTTransferItem
		item.TokenAddress = string(v.GetStringBytes("contractAddress"))
		item.TokenId = string(v.GetStringBytes("tokenID"))
		item.Name = string(v.GetStringBytes("tokenName"))
		item.Symbol = string(v.GetStringBytes("tokenSymbol"))
		item.From = string(v.GetStringBytes("from"))
		item.To = string(v.GetStringBytes("to"))
		item.Timestamp = string(v.GetStringBytes("timeStamp"))
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

	array := parsedJson.GetArray("result")

	nfts := make(map[string]NFTItem)

	for _, v := range array {
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

	result := make([]NFTItem, 0, len(nfts))

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
