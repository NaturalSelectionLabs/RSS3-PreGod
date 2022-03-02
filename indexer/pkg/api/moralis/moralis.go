package moralis

import (
	"fmt"
	"os"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/types"
	"github.com/joho/godotenv"
	jsoniter "github.com/json-iterator/go"
)

const endpoint = "https://deep-index.moralis.io"

var jsoni = jsoniter.ConfigCompatibleWithStandardLibrary

func GetMoralisApiKey() string {
	if err := godotenv.Load(".env"); err != nil {
		return ""
	}

	return os.Getenv("MoralisApiKey")
}

func GetNFTs(userAddress string, chainType MoralisChainType, apiKey string) (types.MoralisNFTResult, error) {
	var headers = map[string]string{
		"accept":    "application/json",
		"X-API-Key": apiKey,
	}

	// Gets all NFT items of user
	url := fmt.Sprintf("%s/api/v2/%s/nft?chain=%s&format=decimal",
		endpoint, userAddress, chainType)
	response, _ := util.Get(url, headers)

	res := new(types.MoralisNFTResult)

	err := jsoni.Unmarshal(response, &res)
	if err != nil {
		return types.MoralisNFTResult{}, err
	}

	return *res, nil
}

func GetNFTTransfers(userAddress string, chainType MoralisChainType, apiKey string) (types.MoralisNFTTransferResult, error) {
	var headers = map[string]string{
		"accept":    "application/json",
		"X-API-Key": apiKey,
	}

	// Gets all NFT transfers of user
	url := fmt.Sprintf("%s/api/v2/%s/nft/transfers?chain=%s&format=decimal&direction=both",
		endpoint, userAddress, chainType)
	response, _ := util.Get(url, headers)

	res := new(types.MoralisNFTTransferResult)

	err := jsoni.Unmarshal(response, &res)
	if err != nil {
		return types.MoralisNFTTransferResult{}, err
	}

	return *res, nil
}
