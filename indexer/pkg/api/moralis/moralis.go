package moralis

import (
	"context"
	"fmt"
	"net/url"
	"sync"
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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	parser      fastjson.Parser
	jsoni       = jsoniter.ConfigCompatibleWithStandardLibrary
	client      *ethclient.Client
	ensContract = "0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85"
	endpoint    = "https://deep-index.moralis.io"
)

func requestMoralisApi(url string, apiKey string, isCache bool) (httpx.Response, error) {
	var headers = map[string]string{
		"accept":    "application/json",
		"X-API-Key": apiKey,
	}

	var err error

	var response httpx.Response

	if isCache {
		response, err = httpx.Get(url, headers)
	} else {
		response, err = httpx.NoCacheGet(url, headers)
	}

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

func GetNFTs(ctx context.Context, userAddress string, chainType ChainType, fromDate string, apiKey string) (NFTResult, error) {
	_, trace := otel.Tracer("crawler_moralis").Start(ctx, "GetNFTs")
	trace.SetAttributes(
		attribute.String("userAddress", userAddress),
		attribute.String("fromDate", fromDate),
	)

	defer trace.End()

	// Gets all NFT items of user
	requestURL := fmt.Sprintf("%s/api/v2/%s/nft?chain=%s&format=decimal&from_date=%s",
		endpoint, userAddress, chainType, url.QueryEscape(fromDate),
	)

	response, err := requestMoralisApi(requestURL, apiKey, true)

	if err != nil {
		return NFTResult{}, err
	}

	res := new(NFTResult)
	SetMoralisAttributes(&res.MoralisAttributes, response)

	err = jsoni.Unmarshal(response.Body, &res)
	if err != nil {
		return NFTResult{}, err
	}

	trace.SetAttributes(
		attribute.Int("responseSize", len(response.Body)),
	)

	lop.ForEach(res.Result, func(item NFTItem, i int) {
		if item.MetaData == "" && item.TokenURI != "" {
			url := nft_utils.FormatUrl(item.TokenURI)
			if metadataRes, err := httpx.Get(url, nil); err != nil {
				logger.Warnf("http get nft metadata error with url '%s': [%v], moralis token uri: %v", url, err, item.TokenURI)
			} else {
				res.Result[i].MetaData = string(metadataRes.Body)
			}
		}
	})

	return *res, nil
}

func GetNFTTransfers(
	ctx context.Context,
	userAddress string,
	chainType ChainType,
	blockHeight int64,
	fromDate string,
	apiKey string,
) (NFTTransferResult, error) {
	_, trace := otel.Tracer("crawler_moralis").Start(ctx, "GetNFTTransfers")
	trace.SetAttributes(
		attribute.String("userAddress", userAddress),
		attribute.String("fromDate", fromDate),
	)

	defer trace.End()

	// Gets all NFT transfers of user
	requestURL := fmt.Sprintf("%s/api/v2/%s/nft/transfers?chain=%s&from_block=%d&format=decimal&direction=both&from_date=%s",
		endpoint, userAddress, chainType, blockHeight, url.QueryEscape(fromDate),
	)
	response, err := requestMoralisApi(requestURL, apiKey, true)

	if err != nil {
		return NFTTransferResult{}, err
	}

	res := new(NFTTransferResult)
	SetMoralisAttributes(&res.MoralisAttributes, response)

	trace.SetAttributes(
		attribute.Int("responseSize", len(response.Body)),
	)

	err = jsoni.Unmarshal(response.Body, &res)
	if err != nil {
		return NFTTransferResult{}, err
	}

	return *res, nil
}

func GetLogs(fromBlock int64, toBlock int64, address string, topic string, chainType ChainType, apiKey string) (GetLogsResult, error) {
	url := fmt.Sprintf("%s/api/v2/%s/logs?chain=%s&from_block=%d&to_block=%d&topic0=%s",
		endpoint, address, string(chainType), fromBlock, toBlock, topic)
	response, err := requestMoralisApi(url, apiKey, false)

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

	response, err := requestMoralisApi(url, apiKey, true)

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
func GetTxByToken(tokenAddress string, tokenId string, chainType ChainType, apiKey string) (TransferItem, error) {
	url := fmt.Sprintf("%s/api/v2/nft/%s/%s/transfers?chain=%s&format=decimal&limit=1",
		endpoint, tokenAddress, tokenId, chainType)
	response, err := requestMoralisApi(url, apiKey, true)

	if err != nil {
		return TransferItem{}, err
	}

	res := new(TransferItem)
	SetMoralisAttributes(&res.MoralisAttributes, response)

	parsedJson, err := parser.Parse(string(response.Body))
	if err != nil {
		logger.Errorf("GetTxByToken: %v", err)

		return TransferItem{}, err
	}

	if err := jsoni.UnmarshalFromString(string(parsedJson.GetObject("result", "0").MarshalTo(nil)), &res); err != nil {
		return TransferItem{}, err
	}

	return *res, nil
}

func GetMetadataByToken(tokenAddress string, tokenId string, chainType ChainType, apiKey string) (NFTItem, error) {
	url := fmt.Sprintf("%s/api/v2/nft/%s/%s?chain=%s&format=decimal&limit=1",
		endpoint, tokenAddress, tokenId, chainType)
	response, err := requestMoralisApi(url, apiKey, true)

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

func GetErc20Transfers(ctx context.Context, userAddress string, chainType ChainType, fromDate string, apiKey string) ([]ERC20TransferItem, error) {
	_, trace := otel.Tracer("crawler_moralis").Start(ctx, "setERC20")
	trace.SetAttributes(
		attribute.String("fromDate", fromDate),
	)

	defer trace.End()

	offset := 0
	transferItems := make([]ERC20TransferItem, 0)

	for {
		transfer, err := getErc20Once(ctx, userAddress, chainType, fromDate, apiKey, offset)
		if err != nil {
			logger.Errorf("get erc20 once error: %v", err)

			return nil, err
		}

		transferItems = append(transferItems, transfer.Result...)

		if len(transferItems) >= 1000 {
			break
		}

		if len(transfer.Result) < transfer.PageSize {
			break
		}

		offset += transfer.PageSize
	}

	return transferItems, nil
}

func getErc20Once(ctx context.Context, userAddress string, chainType ChainType, fromDate string, apiKey string, offest int) (*ERC20Transfer, error) {
	_, trace := otel.Tracer("crawler_moralis").Start(ctx, "setERC20")
	trace.SetAttributes(
		attribute.String("fromDate", fromDate),
	)

	defer trace.End()

	requestURL := fmt.Sprintf("%s/api/v2/%s/erc20/transfers?chain=%s&from_block=%d&offset=%d&from_date=%s",
		endpoint, userAddress, chainType, 0, offest, url.QueryEscape(fromDate),
	)

	response, err := requestMoralisApi(requestURL, apiKey, true)

	if err != nil {
		return nil, err
	}

	trace.SetAttributes(
		attribute.Int("responseSize", len(response.Body)),
	)

	res := new(ERC20Transfer)
	SetMoralisAttributes(&res.MoralisAttributes, response)

	if err = jsoni.Unmarshal(response.Body, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func GetErc20TokenMetaData(chainType ChainType, addresses []string, apiKey string) (Erc20TokensMap, error) {
	if len(addresses) <= 0 {
		return Erc20TokensMap{}, fmt.Errorf("addresss is empty")
	}

	res := Erc20TokensMap{}

	getErc20TokenMetaDataFromCache(addresses, res)

	addrLen := len(addresses)
	logger.Infof("addresss length: %d", addrLen)

	if len(res) == addrLen {
		return res, nil
	}

	limit := 200

	addressBatch := make([][]string, addrLen/limit+1)
	logger.Debugf("sss:%d", addrLen/limit+1)

	for i := 0; i < addrLen; i += limit {
		addressBatch = append(addressBatch, addresses[i:Min(i+limit, addrLen)])
		logger.Infof("len(addressBatch): %v", len(addressBatch))
	}

	lop.ForEach(addressBatch, func(batch []string, i int) {
		batchRes := Erc20TokensMap{}
		getErc20TokenMetaDataFromUrl(chainType, batch, apiKey, batchRes)

		for _, item := range batchRes {
			res[item.Address] = item
		}
	})

	setErc20TokenMetaDataInCache(res)

	return res, nil
}

func Min(x, y int) int {
	if x < y {
		return x
	}

	return y
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

	response, err := requestMoralisApi(url, apiKey, true)

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

/*
 * About eth handler native assets
 */

func GetEthTransfers(ctx context.Context, userAddress string, chainType ChainType, fromDate string, apiKey string) ([]ETHTransferItem, error) {
	_, trace := otel.Tracer("crawler_moralis").Start(ctx, "GetEthTransfers")
	trace.SetAttributes(
		attribute.String("userAddress", userAddress),
		attribute.String("fromDate", fromDate),
	)

	defer trace.End()

	var (
		transferItems = make([]ETHTransferItem, 0)
		wg            sync.WaitGroup
		errorCh       = make(chan error, 1)
		doneCh        = make(chan bool)
		open          = true
	)

	wg.Add(2)

	for offset := 0; offset < 1000; offset += 500 {
		go func(offset int) {
			defer func() {
				wg.Done()
				recover()
			}()

			transfer, err := getETHOnce(ctx, userAddress, chainType, fromDate, apiKey, offset)
			if err != nil {
				logger.Errorf("get eth once error: %v", err)

				if open && offset > 0 {
					errorCh <- err
				}

				return
			}

			transferItems = append(transferItems, transfer.Result...)
		}(offset)
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	select {
	case <-doneCh:
		break
	case err := <-errorCh:
		open = false

		close(errorCh)

		return []ETHTransferItem{}, err
	}

	return transferItems, nil
}

func getETHOnce(ctx context.Context, userAddress string, chainType ChainType, fromDate string, apiKey string, offest int) (*ETHTransfer, error) {
	_, trace := otel.Tracer("crawler_moralis").Start(ctx, "getETHOnce")
	trace.SetAttributes(
		attribute.String("userAddress", userAddress),
		attribute.String("fromDate", fromDate),
	)

	defer trace.End()

	requestURL := fmt.Sprintf("%s/api/v2/%s?chain=%s&from_block=%d&offset=%d&from_date=%s",
		endpoint, userAddress, chainType, 0, offest, url.QueryEscape(fromDate),
	)

	response, err := requestMoralisApi(requestURL, apiKey, true)

	if err != nil {
		return nil, err
	}

	trace.SetAttributes(
		attribute.Int("responseSize", len(response.Body)),
	)

	res := new(ETHTransfer)
	SetMoralisAttributes(&res.MoralisAttributes, response)

	if err = jsoni.Unmarshal(response.Body, &res); err != nil {
		return nil, err
	}

	return res, nil
}
