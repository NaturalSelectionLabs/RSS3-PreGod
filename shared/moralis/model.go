package moralis

import (
	"encoding/json"
)

type Response struct {
	Total    int64           `json:"total"`
	Page     int64           `json:"page"`
	PageSize int64           `json:"page_size"`
	Cursor   string          `json:"cursor"`
	Result   json.RawMessage `json:"result"`
}

type Transaction struct {
	Hash                     string `json:"hash"`
	Nonce                    string `json:"nonce"`
	TransactionIndex         string `json:"transaction_index"`
	FromAddress              string `json:"from_address"`
	ToAddress                string `json:"to_address"`
	Value                    string `json:"value"`
	Gas                      string `json:"gas"`
	GasPrice                 string `json:"gas_price"`
	Input                    string `json:"input"`
	ReceiptCumulativeGasUsed string `json:"receipt_cumulative_gas_used"`
	ReceiptGasUsed           string `json:"receipt_gas_used"`
	ReceiptContractAddress   string `json:"receipt_contract_address"`
	ReceiptRoot              string `json:"receipt_root"`
	ReceiptStatus            string `json:"receipt_status"`
	BlockTimestamp           string `json:"block_timestamp"`
	BlockNumber              string `json:"block_number"`
	BlockHash                string `json:"block_hash"`
}

type NativeBalance struct {
	Balance string `json:"balance"`
}

type TokenTransfer struct {
	TransactionHash string `json:"transaction_hash"`
	Address         string `json:"address"`
	BlockTimestamp  string `json:"block_timestamp"`
	BlockNumber     string `json:"block_number"`
	BlockHash       string `json:"block_hash"`
	ToAddress       string `json:"to_address"`
	FromAddress     string `json:"from_address"`
	Value           string `json:"value"`
}

type NFTTransfer struct {
	TokenAddress     string `json:"token_address"`
	TokenId          string `json:"token_id"`
	FromAddress      string `json:"from_address"`
	ToAddress        string `json:"to_address"`
	Amount           string `json:"amount"`
	ContractType     string `json:"contract_type"`
	BlockNumber      string `json:"block_number"`
	BlockTimestamp   string `json:"block_timestamp"`
	BlockHash        string `json:"block_hash"`
	TransactionHash  string `json:"transaction_hash"`
	TransactionType  string `json:"transaction_type"`
	TransactionIndex int    `json:"transaction_index"`
	LogIndex         int    `json:"log_index"`
}
