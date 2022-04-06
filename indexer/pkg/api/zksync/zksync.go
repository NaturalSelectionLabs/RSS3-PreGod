package zksync

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/httpx"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/valyala/fastjson"
)

const endpoint = "https://api.zksync.io"

func GetLatestBlockHeight() (int64, error) {
	url := endpoint + "/api/v0.1/status"
	response, err := httpx.Get(url, nil)

	if err != nil {
		return 0, err
	}

	var parser fastjson.Parser
	parsedJson, parseErr := parser.Parse(string(response))

	if parseErr != nil {
		return 0, nil
	}

	blockHeight := parsedJson.GetInt64("last_verified")

	return blockHeight, nil
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

	var parser fastjson.Parser
	parsedJson, _ := parser.Parse(string(response))

	array := parsedJson.GetArray()
	tokens := make([]Token, len(array))

	for i, arr := range array {
		tokens[i].Id = arr.GetInt64("id")
		tokens[i].Address = string(arr.GetStringBytes("address"))
		tokens[i].Symbol = string(arr.GetStringBytes("symbol"))
		tokens[i].Decimals = arr.GetInt64("decimals")
		tokens[i].Kind = string(arr.GetStringBytes("kind"))
		tokens[i].IsNFT = arr.GetBool("is_nft")
	}

	return tokens, nil
}

func GetTxsByBlock(blockHeight int64) ([]ZKTransaction, error) {
	url := fmt.Sprintf("%s/api/v0.1/blocks/%d/transactions", endpoint, blockHeight)
	logger.Infof("GetTxsByBlock, url: [%s]", url)
	response, err := httpx.Get(url, nil)

	if err != nil {
		return nil, err
	}

	var parser fastjson.Parser
	parsedJson, _ := parser.Parse(string(response))

	array := parsedJson.GetArray()
	trxs := make([]ZKTransaction, len(array))

	for i, arr := range array {
		trxs[i].TxHash = string(arr.GetStringBytes("tx_hash"))
		trxs[i].BlockNumber = arr.GetInt64("block_number")
		trxs[i].Op.From = string(arr.GetStringBytes("op", "from"))
		trxs[i].Op.To = string(arr.GetStringBytes("op", "to"))
		trxs[i].Op.Type = string(arr.GetStringBytes("op", "type"))
		trxs[i].Op.Nonce = arr.GetInt64("op", "nonce")
		trxs[i].Op.TokenId = arr.GetInt64("op", "token")
		trxs[i].Op.Amount = string(arr.GetStringBytes("op", "amount"))
		trxs[i].Op.AccountId = arr.GetInt64("op", "accountId")
		trxs[i].Success = arr.GetBool("success")
		trxs[i].FailReason = string(arr.GetStringBytes("fail_reason"))
		trxs[i].CreatedAt = string(arr.GetStringBytes("created_at"))
		trxs[i].BatchId = string(arr.GetStringBytes("batch_id"))
	}

	return trxs, nil
}
