package gitcoin

import (
	"math/big"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/moralis"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/zksync"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/util"
	"github.com/valyala/fastjson"
)

const grantUrl = "https://gitcoin.co/grants/grants.json"
const grantsApi = "https://gitcoin.co/api/v0.1/grants/"
const donationSentTopic = "0x3bb7428b25f9bdad9bd2faa4c6a7a9e5d5882657e96c1d24cc41c1d6c1910a98"
const bulkCheckoutAddress = "0x7d655c57f71464B6f83811C55D84009Cd9f5221C"

type tokenMeta struct {
	decimal int64
	symbol  string
}

var token = map[string]tokenMeta{
	"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee": {18, "ETH"},
	"0x6b175474e89094c44da98b954eedeac495271d0f": {18, "DAI"},
	"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48": {6, "USDC"},
	"0xdac17f958d2ee523a2206206994597c13d831ec7": {6, "USDT"},
	"0x514910771af9ca656af840dff83e8264ecf986ca": {18, "LINK"},
	"0xdd1ad9a21ce722c151a836373babe42c868ce9a4": {18, "UBI"},
	"0xd56dac73a4d6766464b38ec6d91eb45ce7457c44": {18, "PAN"},
	"0xb64ef51c888972c908cfacf59b47c1afbc0ab8ac": {8, "STORJ"},
	"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2": {18, "WETH"},
	"0x1f9840a85d5af5bf1d1762f925bdaddc4201f984": {18, "UNI"},
	"0x7d1afa7b718fb893db30a3abc0cfc608aacfebb0": {18, "MATIC"},
	"0xe4815ae53b124e7263f08dcdbbb757d41ed658c6": {18, "ZKS"},
	"0x57ab1ec28d129707052df4df418d58a2d46d5f51": {18, "sUSD"},
	"0x58b6a8a3302369daec383334672404ee733ab239": {18, "LPT"},
	"0x3472a5a71965499acd81997a54bba8d852c6e53d": {18, "BADGER"},
	"0x0d8775f648430679a709e98d2b0cb6250d2887ef": {18, "BAT"},
	"0x03ab458634910aad20ef5f1c8ee96f1d6ac54919": {18, "RAI"},
	"0x12b19d3e2ccc14da04fae33e63652ce469b3f2fd": {12, "GRID"},
	"0x8dd5fbce2f6a956c3022ba3663759011dd51e73e": {18, "TUSD"},
	"0xe41d2489571d322189246dafa5ebde1f4699f498": {18, "ZRX"},
	"0x1cf4592ebffd730c7dc92c1bdffdfc3b9efcf29a": {18, "WAVES"},
	"0x408e41876cccdc0f92210600ef50372656052a38": {18, "REN"},
	"0xbbbbca6a901c926f240b89eacb641d8aec7aeafd": {18, "LRC"},
	"0x69af81e73a73b40adf4f3d4223cd9b1ece623074": {18, "MASK"},
	"0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359": {18, "SNX"},
	"0xc011a72400e58ecd99ee497cf89e3775d4bd732f": {18, "SNX"},
	"0xdfe691f37b6264a90ff507eb359c45d55037951c": {4, "KARMA"},
	"0x36f3fd68e7325a35eb768f1aedaae9ea0689d723": {18, "ESD"},
	"0xa4e8c3ec456107ea67d3075bf9e3df3a75823db0": {18, "LOOM"},
	"0x9992ec3cf6a55b00978cddf2b27bc6882d88d1ec": {18, "POLY"},
	"0x6b3595068778dd592e39a122f4f5a5cf09c90fe2": {18, "SUSHI"},
	"0x2260fac5e5542a773aa44fbcfedf7c193bc2c599": {8, "WBTC"},
	"0x85eee30c52b0b379b046fb0f85f4f3dc3009afec": {18, "KEEP"},
	"0x5732046a883704404f284ce41ffadd5b007fd668": {18, "BLZ"},
	"0xc944e90c64b2c07662a292be6244bdf05cda44a7": {18, "GRT"},
	"0x55296f69f40ea6d20e478533c15a6b08b654e758": {18, "XYO"},
	"0xd26114cd6ee289accf82350c8d8487fedb8a0c07": {18, "OMG"},
	"0x0f5d2fb29fb7d3cfee444a200298f468908cc942": {18, "MANA"},
	"0x84ca8bc7997272c7cfb4d0cd3d55cd942b3c9419": {18, "DIA"},
	"0xdd974d5c2e2928dea5f71b9825b8b646686bd200": {18, "KNC"},
	"0x0000000000004946c0e9f43f4dee607b0ef1fa1c": {0, "CHI"},
	"0x0e29e5abbb5fd88e28b2d355774e73bd47de3bcd": {18, "HAKKA"},
	"0x4e352cf164e64adcbad318c3a1e222e9eba4ce42": {18, "MCB"},
	"0xf1f955016ecbcd7321c7266bccfb96c68ea5e49b": {18, "RLY"},
	"0x1f573d6fb3f13d689ff844b4ce37794d79a7ff1c": {18, "BNT"},
	"0x67c5870b4a41d4ebef24d2456547a03f1f3e094b": {2, "G$"},
	"0x491604c0fdf08347dd1fa4ee062a822a5dd06b5d": {18, "CTSI"},
	"0xb97048628db6b661d4c2aa833e95dbe1a905b280": {18, "PAY"},
	"0x744d70fdbe2ba4cf95131626614a1763df805b9e": {18, "SNT"},
	"0xba100000625a3754423978a60c9317c58a424e3d": {18, "BAL"},
	"0x0e2298e3b3390e3b945a5456fbf59ecc3f55da16": {18, "YAM"},
	"0x2bf91c18cd4ae9c2f2858ef9fe518180f7b5096d": {8, "KIWI"},
	"0xb6ed7644c69416d67b522e20bc294a9a9b405b31": {8, "0xBTC"},
	"0xa19a40fbd7375431fab013a4b08f00871b9a2791": {4, "SWAGG"},
	"0x1776e1f26f98b1a5df9cd347953a26dd3cb46671": {18, "NMR"},
	"0xfc1e690f61efd961294b3e1ce3313fbd8aa4f85d": {18, "aDAI"},
	"0x5a98fcbea516cf06857215779fd812ca3bef1b32": {18, "LDO"},
	"0x4e3fbd56cd56c3e72c1403e103b45db9da5b9d2b": {18, "CVX"},
	"0x875773784af8135ea0ef43b5a374aad105c5d39e": {18, "IDLE"},
	"0xbc396689893d065f41bc2c6ecbee5e0085233447": {18, "PERP"},
	"0x6810e776880c02933d47db1b9fc05908e5386b96": {18, "GNO"},
	"0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2": {18, "MKR"},
}

// GetGrantsInfo returns grant info from gitcoin
func GetGrantsInfo() ([]GrantInfo, error) {
	content, err := util.Get(grantUrl, nil)
	if err != nil {
		return nil, err
	}

	var parser fastjson.Parser
	parsedJson, parseErr := parser.Parse(string(content))

	if parseErr != nil {
		return nil, nil
	}

	grantArrs := parsedJson.GetArray()
	grants := make([]GrantInfo, 0)

	for _, grant := range grantArrs {
		projects := grant.GetArray()

		item := GrantInfo{Title: projects[0].String(), AdminAddress: projects[1].String()}
		grants = append(grants, item)
	}

	return grants, nil
}

// GetProjectsInfo returns project infos from gitcoin
func GetProjectsInfo(adminAddress string, title string) (ProjectInfo, error) {
	var project ProjectInfo

	headers := make(map[string]string)
	util.SetCommonHeader(headers)

	url := grantsApi + "?admin_address=" + adminAddress
	content, err := util.Get(url, headers)

	if err != nil {
		return project, nil
	}

	var parser fastjson.Parser
	parsedJson, parseErr := parser.Parse(string(content))

	if parseErr != nil {
		return project, nil
	}

	if "[]" == string(content) {
		// project is inactive
		project.Active = false
		project.AdminAddress = adminAddress
		project.Title = title
	} else {
		// project is active
		project.Active = true
		project.AdminAddress = adminAddress
		project.Title = title
		project.Id = parsedJson.GetInt64("id")
		project.Slug = string(parsedJson.GetStringBytes("slug"))
		project.Description = string(parsedJson.GetStringBytes("description"))
		project.ReferUrl = string(parsedJson.GetStringBytes("reference_url"))
		project.Logo = string(parsedJson.GetStringBytes("logo"))
		project.TokenAddress = string(parsedJson.GetStringBytes("token_address"))
		project.TokenSymbol = string(parsedJson.GetStringBytes("token_symbol"))
		project.ContractAddress = string(parsedJson.GetStringBytes("contract_address"))
	}

	return project, nil
}

// GetZkSyncDonations returns donations from zksync
func GetZkSyncDonations(fromBlock, toBlock int64) ([]DonationInfo, error) {
	donations := make([]DonationInfo, 0)

	for i := fromBlock; i <= toBlock; i++ {
		trxs, err := zksync.GetTxsByBlock(i)
		if err != nil {
			return nil, err
		}

		for _, tx := range trxs {
			if tx.Op.Type != "Transfer" ||
				!tx.Success ||
				tx.Op.To == "" ||
				zksync.InactiveAdminAddress(tx.Op.To) {
				continue
			}

			tokenId := tx.Op.TokenId
			token := zksync.GetZksToken(tokenId)

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
				Approach:       DonationApproachZksync,
			}
			donations = append(donations, d)
		}
	}

	return donations, nil
}

// GetEthDonations returns donations from ethereum and polygon
func GetEthDonations(fromBlock int64, toBlock int64, chainType ChainType) ([]DonationInfo, error) {
	apiKey := moralis.GetApiKey()
	logs, err := moralis.GetLogs(fromBlock, toBlock, bulkCheckoutAddress, donationSentTopic, string(chainType), apiKey)

	if err != nil {
		return nil, err
	}

	donations := make([]DonationInfo, len(logs.Result))

	for _, item := range logs.Result {
		donor := "0x" + item.Topic3[26:]
		tokenAddress := "0x" + item.Topic1[26:]
		adminAddress := "0x" + item.Data[26:]
		amount := item.Topic2
		formatedAmount := big.NewInt(1)
		formatedAmount.SetString(amount[2:], 16)

		symbol := token[tokenAddress].symbol
		decimal := token[tokenAddress].decimal

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
			Approach:       DonationApproachStandard,
		}

		donations = append(donations, donation)
	}

	return donations, nil
}
