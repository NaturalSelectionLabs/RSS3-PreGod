package zksync

import (
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
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

var (
	DeafultGetNextBlockDuration = 500 * time.Millisecond
	GetLatestNextBlockDuration  = 600 * time.Second
)

type crawlerConfig struct {
	FromHeight    int64
	Step          int64
	MinStep       int64
	Confirmations int64
	SleepInterval time.Duration
	NextRoundTime time.Time
	Interrupt     chan os.Signal
}

var DefaultZksyncConfig = &crawlerConfig{
	FromHeight:    0,
	Step:          50,
	MinStep:       10,
	Confirmations: 15,
	SleepInterval: DeafultGetNextBlockDuration,
	NextRoundTime: time.Now(),
	Interrupt:     make(chan os.Signal, 1),
}

type ZkSyncPlatform string

const (
	ZkSync ZkSyncPlatform = "zksync"
)

type crawler struct {
	config           *crawlerConfig
	platform         ZkSyncPlatform
	networkID        constants.NetworkID
	platformID       constants.PlatformID
	metadataIdentity string
}

var ZksTokensCache = map[int]Token{}

type ZkSyncInfo struct {
	From           string
	To             string
	TokenAddress   string
	Amount         string
	Symbol         string
	FormatedAmount *big.Int
	Decimals       int
	Timestamp      string
	TxHash         string
	Type           string
}

var (
	zkCP = crawler{
		config:           DefaultZksyncConfig,
		platform:         ZkSync,
		networkID:        constants.NetworkIDZkSync,
		platformID:       constants.PlatformID(1010),
		metadataIdentity: "zksync",
	}
)

type noteInstanceBuilder struct {
	countMap map[string]int
}
