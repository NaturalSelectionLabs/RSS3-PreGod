package zksync

import (
	"fmt"
)

type Token struct {
	Id       int    `json:"id"`
	Address  string `json:"address"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
	Kind     string `json:"kind"`
	IsNft    bool   `json:"is_nft"`
}

func (t Token) String() string {
	return fmt.Sprintf(`Id: %d, TokenAddress: %s, Symbol: %s, Decimals: %d, Kind: %s, IsNFT: %v`,
		t.Id, t.Address, t.Symbol, t.Decimals, t.Kind, t.IsNft)
}

// nolint:tagliatelle // returned by zksync api
type Op struct {
	To        string `json:"to"`
	Fee       string `json:"fee"`
	From      string `json:"from"`
	Type      string `json:"type"`
	Nonce     int64  `json:"nonce"`
	TokenId   int    `json:"token"`
	Amount    string `json:"amount"`
	AccountId int64  `json:"accountId"`
	Signature struct {
		PubKey    string `json:"pubKey"`
		Signature string `json:"signature"`
	} `json:"signature"`
}

func (o Op) String() string {
	return fmt.Sprintf(`From: %s, To: %s, Type: %s, TokenId: %d, Amount: %s`,
		o.From, o.To, o.Type, o.TokenId, o.Amount)
}

type ZKTransaction struct {
	TxHash      string      `json:"tx_hash"`
	BlockIndex  int64       `json:"block_index"`
	BlockNumber int64       `json:"block_number"`
	Op          Op          `json:"op"`
	Success     bool        `json:"success"`
	FailReason  interface{} `json:"fail_reason"`
	CreatedAt   string      `json:"created_at"`
	BatchId     interface{} `json:"batch_id"`
}

func (t ZKTransaction) String() string {
	return fmt.Sprintf(`TxHash: %s, BlockIndex: %d, BlockNumber: %d, Op: %s, Success: %v,CreatedAt:%v `,
		t.TxHash, t.BlockIndex, t.BlockNumber, t.Op, t.Success, t.CreatedAt)
}

type StatusResult struct {
	NextBlockAtMax    interface{} `json:"next_block_at_max"`
	LastCommitted     int64       `json:"last_committed"`
	LastVerified      int64       `json:"last_verified"`
	TotalTransactions int64       `json:"total_transactions"`
	OutstandingTxs    int64       `json:"outstanding_txs"`
	MempoolSize       int64       `json:"mempool_size"`
	CoreStatus        struct {
		MainDatabaseAvailable    bool `json:"main_database_available"`
		ReplicaDatabaseAvailable bool `json:"replica_database_available"`
		Web3Available            bool `json:"web3_available"`
	} `json:"core_status"`
}
