package moralis

import (
	"fmt"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/nft_utils"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/datatype"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/httpx"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	jsoniter "github.com/json-iterator/go"
	lop "github.com/samber/lo/parallel"
	"github.com/valyala/fastjson"
	goens "github.com/wealdtech/go-ens/v3"
)

var (
	parser      fastjson.Parser
	jsoni       = jsoniter.ConfigCompatibleWithStandardLibrary
	client      *ethclient.Client
	ensContract = "0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85"
	endpoint    = "https://deep-index.moralis.io"
)

func requestMoralisApi(url string, apiKey string) (httpx.Response, error) {
	var headers = map[string]string{
		"accept":    "application/json",
		"X-API-Key": apiKey,
	}

	response, err := httpx.NoCacheGet(url, headers)
	if err != nil {
		logger.Errorf("http get error with url '%s': [%v]. response: %s",
			url, err, string(response.Body))

		return response, err
	}

	return response, nil
}

/*
 * About nft handler
 */

func GetNFTs(userAddress string, chainType ChainType, apiKey string) (NFTResult, error) {
	// Gets all NFT items of user
	url := fmt.Sprintf("%s/api/v2/%s/nft?chain=%s&format=decimal",
		endpoint, userAddress, chainType)

	response, err := requestMoralisApi(url, apiKey)

	if err != nil {
		return NFTResult{}, err
	}

	res := new(NFTResult)
	SetMoralisAttributes(&res.MoralisAttributes, response)

	err = jsoni.Unmarshal(response.Body, &res)
	if err != nil {
		return NFTResult{}, err
	}

	lop.ForEach(res.Result, func(item NFTItem, i int) {
		if item.MetaData == "" && item.TokenURI != "" {
			if metadataRes, err := httpx.Get(nft_utils.FormatUrl(item.TokenURI), nil); err != nil {
				logger.Warnf("http get nft metadata error with url '%s': [%v]", item.TokenURI, err)
			} else {
				res.Result[i].MetaData = string(metadataRes.Body)
			}
		}
	})

	return *res, nil
}

func GetNFTTransfers(userAddress string, chainType ChainType, blockHeight int64, apiKey string) (NFTTransferResult, error) {
	// Gets all NFT transfers of user
	url := fmt.Sprintf("%s/api/v2/%s/nft/transfers?chain=%s&from_block=%d&format=decimal&direction=both",
		endpoint, userAddress, chainType, blockHeight)
	response, err := requestMoralisApi(url, apiKey)

	if err != nil {
		return NFTTransferResult{}, err
	}

	res := new(NFTTransferResult)
	SetMoralisAttributes(&res.MoralisAttributes, response)

	err = jsoni.Unmarshal(response.Body, &res)
	if err != nil {
		return NFTTransferResult{}, err
	}

	return *res, nil
}

func GetLogs(fromBlock int64, toBlock int64, address string, topic string, chainType ChainType, apiKey string) (GetLogsResult, error) {
	url := fmt.Sprintf("%s/api/v2/%s/logs?chain=%s&from_block=%d&to_block=%d&topic0=%s",
		endpoint, address, string(chainType), fromBlock, toBlock, topic)
	response, err := requestMoralisApi(url, apiKey)

	if err != nil {
		return GetLogsResult{}, err
	}

	res := new(GetLogsResult)
	SetMoralisAttributes(&res.MoralisAttributes, response)

	err = jsoni.Unmarshal(response.Body, &res)
	if err != nil {
		logger.Errorf("unmarshal error: [%v]", err)

		return *res, err
	}

	return *res, nil
}

// Gets all NFT items of user
func GetNFTByContract(userAddress string, contactAddress string, chainType ChainType, apiKey string) (NFTResult, error) {
	// this function is used by ENS indexer.
	url := fmt.Sprintf("%s/api/v2/%s/nft?chain=%s&format=decimal&token_addresses=%s",
		endpoint, userAddress, chainType, contactAddress)

	response, err := requestMoralisApi(url, apiKey)

	if err != nil {
		return NFTResult{}, err
	}

	res := new(NFTResult)
	SetMoralisAttributes(&res.MoralisAttributes, response)

	err = jsoni.Unmarshal(response.Body, &res)
	if err != nil {
		return NFTResult{}, err
	}

	return *res, nil
}

// GetTxByToken is used by ENS indexer
func GetTxByToken(tokenAddress string, tokenId string, chainType ChainType, apiKey string) (NFTTransferItem, error) {
	url := fmt.Sprintf("%s/api/v2/nft/%s/%s/transfers?chain=%s&format=decimal&limit=1",
		endpoint, tokenAddress, tokenId, chainType)
	response, err := requestMoralisApi(url, apiKey)

	if err != nil {
		return NFTTransferItem{}, err
	}

	res := new(NFTTransferItem)
	SetMoralisAttributes(&res.MoralisAttributes, response)

	parsedJson, err := parser.Parse(string(response.Body))
	if err != nil {
		logger.Errorf("GetTxByToken: %v", err)

		return NFTTransferItem{}, err
	}

	if err := jsoni.UnmarshalFromString(string(parsedJson.GetObject("result", "0").MarshalTo(nil)), &res); err != nil {
		return NFTTransferItem{}, err
	}

	return *res, nil
}

func GetMetadataByToken(tokenAddress string, tokenId string, chainType ChainType, apiKey string) (NFTItem, error) {
	url := fmt.Sprintf("%s/api/v2/nft/%s/%s?chain=%s&format=decimal&limit=1",
		endpoint, tokenAddress, tokenId, chainType)
	response, err := requestMoralisApi(url, apiKey)

	if err != nil {
		return NFTItem{}, err
	}

	res := new(NFTItem)
	SetMoralisAttributes(&res.MoralisAttributes, response)

	err = jsoni.Unmarshal(response.Body, &res)
	if err != nil {
		return NFTItem{}, nil
	}

	return *res, nil
}

/*
 * About erc20 handler
 */

type Erc20TokensMap map[string]Erc20TokenMetaDataItem

var erc20TokensCache = Erc20TokensMap{}

var erc20TokensPackageSize = 200

func GetErc20Transfers(userAddress string, chainType ChainType, apiKey string) ([]ERC20TransferItem, error) {
	offset := 0
	transferItems := make([]ERC20TransferItem, 0)
	var lastTransfer *ERC20Transfer

	for {
		transfer, err := getErc20Once(userAddress, chainType, apiKey, offset)
		if err != nil {
			logger.Warnf("get erc20 once error: %v", err)
			continue
		}

		// Since there is a problem with the page-turning function of Moralis,
		// it is necessary to check whether the page is turned to the end from the previous block result each time.
		if transferCompare(transfer, lastTransfer) {
			break
		}

		transferItems = append(transferItems, transfer.Result...)

		if len(transfer.Result) < transfer.PageSize {
			break
		}

		// Due to a problem with the interface of Moralis,
		// the situation where there may be a page in the same block is filtered out here.
		lastTransfer = transfer
		offset += transfer.PageSize
	}

	return transferItems, nil
}

func transferCompare(currTrans *ERC20Transfer, lastTrans *ERC20Transfer) bool {
	if currTrans == nil || lastTrans == nil {
		return false
	}

	if currTrans.Total != lastTrans.Total {
		return false
	}

	if currTrans.Page != lastTrans.Page {
		return false
	}

	if currTrans.PageSize != lastTrans.PageSize {
		return false
	}

	if currTrans.Cursor != lastTrans.Cursor {
		return false
	}

	for i, item := range currTrans.Result {
		if lastTrans.Result[i] != item {
			return false
		}
	}

	return true
}

func getErc20Once(userAddress string, chainType ChainType, apiKey string, offest int) (*ERC20Transfer, error) {
	url := fmt.Sprintf("%s/api/v2/%s/erc20/transfers?chain=%s&from_block=%d&offset=%d",
		endpoint, userAddress, chainType, 0, offest)
	logger.Debugf("get erc20 once url: %s", url)

	// if toBlock > 0 {
	// 	url = fmt.Sprintf("%s&to_block=%d", url, toBlock)
	// }
	response, err := requestMoralisApi(url, apiKey)

	if err != nil {
		return nil, err
	}

	res := new(ERC20Transfer)
	SetMoralisAttributes(&res.MoralisAttributes, response)

	if err = jsoni.Unmarshal(response.Body, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func GetErc20TokenMetaData(chainType ChainType, addresses []string, apiKey string) (Erc20TokensMap, error) {
	logger.Debugf("GetErc20TokenMetaData: %v", addresses)
	if len(addresses) <= 0 {
		return Erc20TokensMap{}, fmt.Errorf("addresss is empty")
	}

	res := Erc20TokensMap{}

	getErc20TokenMetaDataFromCache(addresses, res)

	if len(res) == len(addresses) {
		return res, nil
	}

	getErc20TokenMetaDataFromUrl(chainType, addresses, apiKey, res)

	setErc20TokenMetaDataInCache(res)

	return res, nil
}

func getErc20TokenMetaDataFromCache(addresses []string, res Erc20TokensMap) {
	for _, address := range addresses {
		if v, ok := erc20TokensCache[address]; ok {
			res[address] = v
		}
	}
}

func getErc20TokenMetaDataFromUrl(chainType ChainType, addresses []string, apiKey string, res Erc20TokensMap) error {
	url := fmt.Sprintf("%s/api/v2/erc20/metadata?chain=%s",
		endpoint, chainType)

	for _, address := range addresses {
		url += fmt.Sprintf("&addresses=%s", address)
	}
	logger.Debugf("url: %s", url)

	response, err := requestMoralisApi(url, apiKey)

	if err != nil {
		return err
	}

	resp := make([]Erc20TokenMetaDataItem, 0)
	attributes := new(MoralisAttributes)
	SetMoralisAttributes(attributes, response)

	err = jsoni.Unmarshal(response.Body, &resp)
	if err != nil {
		return err
	}

	if len(resp) > 0 {
		for _, item := range resp {
			res[item.Address] = item
		}
	}

	return nil
}

func setErc20TokenMetaDataInCache(res Erc20TokensMap) {
	for address, metaData := range res {
		_, ok := erc20TokensCache[address]
		if !ok {
			res[address] = metaData
		}
	}
}

/*
 * About ens handler
 */

// returns a list of ENS domains with non-empty text records
func GetENSList(address string) ([]ENSTextRecord, error) {
	getGatewayClient()

	result := []ENSTextRecord{}

	domain, err := goens.ReverseResolve(client, common.HexToAddress(address))

	if err != nil {
		logger.Errorf("goens.ReverseResolve: %v", err)

		return nil, err
	}

	record := ENSTextRecord{
		Domain: domain,
	}

	err = getENSDetail(address, &record)

	if err != nil {
		return nil, err
	}

	err = getENSTextValue(domain, &record)
	if err != nil {
		return nil, err
	}

	result = append(result, record)

	return result, err
}

// reads the text record value for a given domain with the list of predefined keys
func getENSTextValue(domain string, record *ENSTextRecord) error {
	r, err := goens.NewResolver(client, domain)

	if err != nil {
		logger.Errorf("getENSTextValue NewResolver: %v", err)

		return err
	}

	record.Text = make(map[string]string)

	for _, key := range getTextRecordKeyList() {
		text, err := r.Text(key)

		if err != nil {
			logger.Errorf("getENSTextValue read text: %v", err)

			return err
		}

		record.Text[key] = text

		// append attachments
		switch key {
		case "url":
			if text != "" {
				a := datatype.Attachment{
					Type:     "websites",
					MimeType: "text/uri-list",
					Content:  text,
				}
				record.Attachments = append(record.Attachments, a)
			}
		}
	}

	return nil
}

// returns ENS details from moralis
func getENSDetail(address string, record *ENSTextRecord) error {
	ensList, err := GetNFTByContract(address, ensContract, ETH, getApiKey())

	if err != nil {
		logger.Errorf("getENSDetail GetNFTByContract: %v", err)

		return err
	}

	for _, ens := range ensList.Result {
		// moralis sometimes returns empty metadata
		if ens.MetaData != "" {
			meta := new(NFTMetadata)

			err = jsoni.UnmarshalFromString(ens.MetaData, &meta)

			if err != nil {
				logger.Errorf("getENSDetail unmarshall metadata: %v", err)

				return err
			}

			// an address might have multiple ENS domains
			// if the one is the current ENS domain
			if meta.Name == record.Domain {
				record.Description = meta.Description

				avatar := "https://metadata.ens.domains/mainnet/" + ensContract + "/" + ens.TokenId + "/image"

				record.Attachments = append(record.Attachments, datatype.Attachment{
					Type:    "banner",
					Address: avatar,
				})

				record.Avatar = avatar

				return getENSTransaction(ens, record)
			}
		}
	}

	return nil
}

func getENSTransaction(ens NFTItem, record *ENSTextRecord) error {
	// get TxHash and Tsp with TokenId from Moralis
	t, err := GetTxByToken(ens.TokenAddress, ens.TokenId, ETH, getApiKey())

	if err != nil {
		logger.Errorf("getENSDetail transaction: %v", err)

		return err
	}

	record.TxHash = t.TransactionHash
	record.CreatedAt, err = time.Parse(time.RFC3339, t.BlockTimestamp)

	if err != nil {
		logger.Errorf("getENSDetail transaction: %v", err)

		return err
	}

	return nil
}
