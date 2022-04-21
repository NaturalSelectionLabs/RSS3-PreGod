package gitcoin

import (
	"fmt"
	"math/big"
	"os"
	"time"
)

type GitcoinPlatform string

const (
	Unknown GitcoinPlatform = "unknown"

	ETH     GitcoinPlatform = "eth"
	Polygon GitcoinPlatform = "polygon"
	ZkSync  GitcoinPlatform = "zksync"
)

func (p GitcoinPlatform) getContractAddress() string {
	if p == ETH {
		return bulkCheckoutAddressETH
	}

	if p == Polygon {
		return bulkCheckoutAddressPolygon
	}

	return ""
}

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
	FromHeight:    2600,
	Step:          50,
	MinStep:       10,
	Confirmations: 15,
	SleepInterval: 600 * time.Second,
	NextRoundTime: time.Now(),
	Interrupt:     make(chan os.Signal, 1),
}

var DefaultEthConfig = &crawlerConfig{
	FromHeight:    10245999, // gitcoin bulkCheckout contract was created at block #10245999
	Step:          50,
	MinStep:       10,
	Confirmations: 15,
	SleepInterval: 600 * time.Second,
	NextRoundTime: time.Now(),
	Interrupt:     make(chan os.Signal, 1),
}

var DefaultPolygonConfig = &crawlerConfig{
	FromHeight:    18682002, // gitcoin bulkCheckout contract was created at block #10245999
	Step:          50,
	MinStep:       10,
	Confirmations: 120,
	SleepInterval: 600 * time.Second,
	NextRoundTime: time.Now(),
	Interrupt:     make(chan os.Signal, 1),
}

type DonationApproach string

const (
	DonationApproachEthereum = "Standard"
	DonationApproachPolygon  = "Polygon"
	DonationApproachZkSync   = "zkSync"
)

type GrantInfo struct {
	Title        string
	AdminAddress string
}

type ProjectInfo struct {
	Active          bool   `gorm:"column:active"`
	Id              int    `gorm:"column:id"`
	Title           string `gorm:"column:title"`
	Slug            string `gorm:"column:slug"`
	Description     string `gorm:"column:description"`
	ReferUrl        string `gorm:"column:reference_url"`
	Logo            string `gorm:"column:logo"`
	AdminAddress    string `gorm:"column:admin_address"`
	TokenAddress    string `gorm:"column:token_address"`
	TokenSymbol     string `gorm:"column:token_symbol"`
	ContractAddress string `gorm:"column:contract_address"`
}

func (ProjectInfo) TableName() string {
	return "reptile-gitcoin.data"
}

type DonationInfo struct {
	Donor          string
	AdminAddress   string
	TokenAddress   string
	Amount         string
	Symbol         string
	FormatedAmount *big.Int
	Decimals       int
	Timestamp      string
	TxHash         string
	Approach       DonationApproach
}

func (d DonationInfo) String() string {
	return fmt.Sprintf(`Donor: %s, AdminAddress: %s, TokenAddress: %s, Amount: %s, Symbol: %s, TxHash: %s`,
		d.Donor, d.AdminAddress, d.TokenAddress, d.Amount, d.Symbol, d.TxHash)
}

func (d DonationInfo) GetTxTo() string {
	if d.Approach == DonationApproachEthereum {
		return bulkCheckoutAddressETH
	}

	if d.Approach == DonationApproachPolygon {
		return bulkCheckoutAddressPolygon
	}

	if d.Approach == DonationApproachZkSync {
		return d.AdminAddress
	}

	return ""
}

type TokenMeta struct {
	Address string `json:"address"`
	Decimal int    `json:"decimal"`
	Symbol  string `json:"symbol"`
}
