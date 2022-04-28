package zksync

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/httpx"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	jsoniter "github.com/json-iterator/go"
)

var (
	jsoni = jsoniter.ConfigCompatibleWithStandardLibrary
)

const endpoint = "https://api.zksync.io"

func GetLatestBlockHeight() (int64, error) {
	url := endpoint + "/api/v0.1/status"
	response, err := httpx.Get(url, nil)

	if err != nil {
		return 0, err
	}

	statusResult := new(StatusResult)
	if err := jsoni.UnmarshalFromString(string(response.Body), statusResult); err != nil {
		logger.Errorf("zksync GetLatestBlockHeight unmarshalFromString error: %v", err)

		return 0, err
	}

	return statusResult.LastVerified, nil
}

func GetLatestBlockHeightWithConfirmations(confirmations int64) (int64, error) {
	// get latest block height
	latestBlockHeight, err := GetLatestBlockHeight()
	if err != nil {
		return 0, err
	}

	return latestBlockHeight - confirmations, nil
}

func GetTokens() ([]Token, error) {
	url := endpoint + "/api/v0.1/tokens"
	response, err := httpx.Get(url, nil)

	if err != nil {
		return nil, err
	}

	var tokens []Token
	if err = jsoni.UnmarshalFromString(string(response.Body), &tokens); err != nil {
		return nil, fmt.Errorf("GetTokens UnmarshalFromString error: [%v]", err)
	}

	return tokens, nil
}

func GetTxsByBlock(blockHeight int64) ([]ZKTransaction, error) {
	url := fmt.Sprintf("%s/api/v0.1/blocks/%d/transactions", endpoint, blockHeight)
	response, err := httpx.Get(url, nil)

	if err != nil {
		return nil, err
	}

	var zkTxs []ZKTransaction
	if err = jsoni.UnmarshalFromString(string(response.Body), &zkTxs); err != nil {
		return nil, fmt.Errorf("GetTokens UnmarshalFromString error: [%v]", err)
	}

	return zkTxs, nil
}
