package zksync

import (
	"encoding/json"
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
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

func GetTxsByBlock(blockHeight int64, isSaveDB bool) ([]ZKTransaction, error) {
	var zkTxs []ZKTransaction

	zkTxs, err := getTxsFromDb(blockHeight)
	if err == nil && len(zkTxs) > 0 {
		return zkTxs, nil
	} else if err != nil {
		logger.Warnf("zksync get txs by db error: %v", err)
	}

	zkTxs, err = getTxsByBlockByUrl(blockHeight)
	if err != nil {
		logger.Warnf("get txs By block by url error: %v", err)
	}

	if len(zkTxs) == 0 {
		return zkTxs, nil
	}

	if isSaveDB {
		err = saveTxsInDB(zkTxs, blockHeight)
		if err != nil {
			logger.Warnf("zksync save txs in db error: %v", err)
		}
	}

	return zkTxs, nil
}

func getTxsByBlockByUrl(blockHeight int64) ([]ZKTransaction, error) {
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

func getTxsFromDb(blockHeight int64) ([]ZKTransaction, error) {
	caches, err := database.QueryCaches(
		database.DB, constants.NetworkSymbolZkSync.String(), endpoint, blockHeight, blockHeight)

	if err != nil {
		return nil, err
	}

	var zkTxs = []ZKTransaction{}

	for _, cache := range caches {
		zkTx := ZKTransaction{}
		if err = jsoni.UnmarshalFromString(string(cache.Data), &zkTx); err != nil {
			return nil, fmt.Errorf("get tokens unmarshal from string error: %v", err)
		}

		zkTxs = append(zkTxs, zkTx)
	}

	return zkTxs, nil
}

func saveTxsInDB(zkTxs []ZKTransaction, blockHeight int64) error {
	caches := []model.Cache{}

	for _, zkTx := range zkTxs {
		zkTxJson, err := json.Marshal(zkTx)
		if err != nil {
			logger.Warnf("zksync[%s] save txs in db error: %v", zkTx.TxHash, err)

			continue
		}

		cache := model.Cache{
			Key:      zkTx.TxHash,
			Network:  constants.NetworkSymbolZkSync.String(),
			Source:   endpoint,
			BlockNum: blockHeight,
			LogIndex: 0,
			Data:     zkTxJson,
		}

		caches = append(caches, cache)
	}

	if err := database.CreateCaches(database.DB, caches, true); err != nil {
		return fmt.Errorf("zksync block height[%d] save txs in db error: %v", blockHeight, err)
	}

	return nil
}

func UpdateZksToken() error {
	tokens, err := GetTokens()
	if err != nil {
		logger.Errorf("zksync get tokens error: %v", err)

		return err
	}

	for _, token := range tokens {
		ZksTokensCache[token.Id] = token
	}

	return nil
}
