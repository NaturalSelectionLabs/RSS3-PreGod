package gitcoin

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/moralis"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/zksync"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/httpx"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
)

const gitCoinTokensUrl = "https://gitcoin.co/api/v1/tokens"
const donationSentTopic = "0x3bb7428b25f9bdad9bd2faa4c6a7a9e5d5882657e96c1d24cc41c1d6c1910a98"
const bulkCheckoutAddressETH = "0x7d655c57f71464B6f83811C55D84009Cd9f5221C"
const bulkCheckoutAddressPolygon = "0xb99080b9407436eBb2b8Fe56D45fFA47E9bb8877"

type tokenMeta struct {
	decimal int
	symbol  string
}

var (
	token = map[string]tokenMeta{
		"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee": {18, "ETH"},
	}

	jsoni = jsoniter.ConfigCompatibleWithStandardLibrary
)

func UpdateEthAndPolygonTokens() error {
	url := gitCoinTokensUrl
	response, err := httpx.Get(url, nil)

	if err != nil {
		return fmt.Errorf("get eth and polygon token err:%s", err)
	}

	var result []TokenMeta
	if err = jsoni.UnmarshalFromString(string(response.Body), &result); err != nil {
		return fmt.Errorf("get eth and polygon token err:%s", err)
	}

	if len(result) > 0 {
		for _, v := range result {
			meta := tokenMeta{
				decimal: v.Decimal,
				symbol:  v.Symbol,
			}
			address := strings.ToLower(v.Address)
			token[address] = meta
		}
	}

	return nil
}

type DonationsResult struct {
	Donations      []DonationInfo
	AdminAddresses []string
}

type EthDonationsResult struct {
	DonationsResult
	MinRateLimit     int
	MinRateLimitUsed int
}

func NewEthDonationsResult() *EthDonationsResult {
	return &EthDonationsResult{
		DonationsResult: DonationsResult{
			Donations:      make([]DonationInfo, 0),
			AdminAddresses: make([]string, 0),
		},
		MinRateLimit:     0,
		MinRateLimitUsed: 0,
	}
}

// GetEthDonations returns donations from ethereum and polygon
func GetEthDonations(fromBlock int64, toBlock int64, chainType GitcoinPlatform) (*EthDonationsResult, error) {
	var checkoutAddress string

	var donationApproach DonationApproach

	ethDonationsResult := NewEthDonationsResult()

	if chainType == ETH {
		checkoutAddress = bulkCheckoutAddressETH
		donationApproach = DonationApproachEthereum
	} else if chainType == Polygon {
		checkoutAddress = bulkCheckoutAddressPolygon
		donationApproach = DonationApproachPolygon
	} else {
		return nil, fmt.Errorf("invalid chainType %s", string(chainType))
	}

	// at most 1000 results in one response. But our default step is only 50, safe.
	logs, err := moralis.GetLogs(fromBlock, toBlock, checkoutAddress, donationSentTopic,
		moralis.ChainType(chainType), config.Config.Indexer.Moralis.ApiKey)

	if err != nil {
		// MinRateLimit must be sent here
		return ethDonationsResult, fmt.Errorf("getLogs error: [%v]", err)
	}

	ethDonationsResult.MinRateLimit = logs.MinRateLimit
	ethDonationsResult.MinRateLimitUsed = logs.MinRateLimitUsed

	for _, item := range logs.Result {
		donor := "0x" + item.Topic3[26:]
		tokenAddress := "0x" + item.Topic1[26:]
		adminAddress := "0x" + item.Data[26:]
		amount := item.Topic2
		formatedAmount := big.NewInt(1)
		formatedAmount.SetString(amount[2:], 16)

		t, ok := token[tokenAddress]
		if !ok {
			logger.Warnf("token address doesn't exist: %s", tokenAddress)

			continue
		}

		symbol := t.symbol
		decimal := t.decimal

		donation := DonationInfo{
			Donor:          donor,
			AdminAddress:   adminAddress,
			TokenAddress:   tokenAddress,
			Amount:         amount,
			Symbol:         symbol,
			FormatedAmount: formatedAmount,
			Decimals:       decimal,
			Timestamp:      item.BlockTimestamp,
			TxHash:         item.TransactionHash,
			Approach:       donationApproach,
		}

		ethDonationsResult.Donations = append(ethDonationsResult.Donations, donation)
		ethDonationsResult.AdminAddresses = append(ethDonationsResult.AdminAddresses, adminAddress)
	}

	return ethDonationsResult, nil
}

// Asynchronous Zk query start
var ZksTokensCache = map[int]zksync.Token{}

func UpdateZksToken() error {
	tokens, err := zksync.GetTokens()
	if err != nil {
		logger.Errorf("zksync get tokens error: %v", err)

		return err
	}

	for _, token := range tokens {
		ZksTokensCache[token.Id] = token
	}

	return nil
}

func GetZksToken(id int) zksync.Token {
	return ZksTokensCache[id]
}

type ZkSyncDonationResult struct {
	DonationsResult
}

func NewZkSyncDonationResult() *ZkSyncDonationResult {
	return &ZkSyncDonationResult{
		DonationsResult: DonationsResult{
			Donations:      make([]DonationInfo, 0),
			AdminAddresses: make([]string, 0),
		},
	}
}

// GetZkSyncDonations returns donations from zksync
func GetZkSyncDonations(fromBlock int64, toBlock int64) (*ZkSyncDonationResult, error) {
	ethDonationsResult := NewZkSyncDonationResult()

	for i := fromBlock; i <= toBlock; i++ {
		trxs, err := zksync.GetTxsByBlock(i)
		if err != nil {
			logger.Errorf("get txs by block error: [%v]", err)

			return nil, err
		}

		for _, tx := range trxs {
			if tx.Op.Type != "Transfer" || !tx.Success {
				continue
			}

			// admin address empty
			adminAddress := strings.ToLower(tx.Op.To)
			if adminAddress == "" ||
				adminAddress == "0x0" ||
				adminAddress == "0x0000000000000000000000000000000000000000" {
				continue
			}

			tokenId := tx.Op.TokenId
			token := GetZksToken(tokenId)

			formatedAmount := big.NewInt(1)
			formatedAmount.SetString(tx.Op.Amount, 10)

			d := DonationInfo{
				Donor:          tx.Op.From,
				AdminAddress:   tx.Op.To,
				TokenAddress:   token.Address,
				Amount:         tx.Op.Amount,
				Symbol:         token.Symbol,
				FormatedAmount: formatedAmount,
				Decimals:       token.Decimals,
				Timestamp:      tx.CreatedAt,
				TxHash:         tx.TxHash,
				Approach:       DonationApproachZkSync,
			}
			ethDonationsResult.Donations = append(ethDonationsResult.Donations, d)
			ethDonationsResult.AdminAddresses = append(ethDonationsResult.AdminAddresses, adminAddress)
		}
	}

	return ethDonationsResult, nil
}

// GetProjectsInfo returns project info from gitcoin

func queryProjectsInfo(db *gorm.DB, adminAddresses []string) (map[string]ProjectInfo, error) {
	projects := make([]ProjectInfo, 0)
	projectsMap := make(map[string]ProjectInfo)

	if err := db.Where(
		"admin_address in (?)", adminAddresses).Find(&projects).Error; err != nil {
		return nil, err
	}

	for _, project := range projects {
		projectsMap[project.AdminAddress] = project
	}

	return projectsMap, nil
}

func GetProjectsInfo(adminAddresses []string) (map[string]ProjectInfo, error) {
	projects, err := queryProjectsInfo(database.DB, adminAddresses)
	if err != nil {
		return nil, fmt.Errorf("get project info from db false:[%s]", err)
	}

	return projects, nil
}
