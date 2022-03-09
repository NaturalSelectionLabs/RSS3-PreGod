package zksync

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util/httpx"
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

func GetTokens() ([]Token, error) {
	url := endpoint + "/api/v0.1/tokens"
	response, err := httpx.Get(url, nil)

	if err != nil {
		return nil, err
	}

	var parser fastjson.Parser
	parsedJson, _ := parser.Parse(string(response))

	arrs := parsedJson.GetArray()
	tokens := make([]Token, len(arrs))

	for i, arr := range arrs {
		tokens[i].Id = arr.GetInt64("id")
		tokens[i].Address = string(arr.GetStringBytes("address"))
		tokens[i].Symbol = string(arr.GetStringBytes("symbol"))
		tokens[i].Decimals = arr.GetInt64("decimals")
		tokens[i].Kind = string(arr.GetStringBytes("kind"))
		tokens[i].IsNFT = arr.GetBool("is_nft")
	}

	return tokens, nil
}

// GetZksToken returns Token by tokenId
func GetZksToken(id int64) Token {
	// TODO: get Token from db
	token := Token{Id: 0,
		Address:  "0x0000000000000000000000000000000000000000",
		Symbol:   "ETH",
		Decimals: 18,
		Kind:     "ERC20",
		IsNFT:    false,
	}

	return token
}

// InactiveAdminAddress checks an admin address is active or not
func InactiveAdminAddress(adminAddress string) bool {
	// TODO: check from db
	return false
}

func GetTxsByBlock(blockHeight int64) ([]ZKTransaction, error) {
	url := fmt.Sprintf("%s/api/v0.1/blocks/%d/transactions", endpoint, blockHeight)
	response, err := httpx.Get(url, nil)

	if err != nil {
		return nil, err
	}

	var parser fastjson.Parser
	parsedJson, _ := parser.Parse(string(response))

	arrs := parsedJson.GetArray()
	trxs := make([]ZKTransaction, len(arrs))

	for i, arr := range arrs {
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
