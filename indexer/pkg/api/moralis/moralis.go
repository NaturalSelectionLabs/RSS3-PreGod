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
	"go.opentelemetry.io/otel/codes"
)

const (
	TracerNameCrawlerMoralis = "crawler_moralis"
)

var (
	parser      fastjson.Parser
	jsoni       = jsoniter.ConfigCompatibleWithStandardLibrary
	client      *ethclient.Client
	ensContract = "0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85"
	endpoint    = "https://deep-index.moralis.io"
)

func requestMoralisApi(ctx context.Context, url string, apiKey string, isCache bool) (httpx.Response, error) {
	tracer := otel.Tracer(TracerNameCrawlerMoralis)

	_, moralisSnap := tracer.Start(ctx, "moralis")

	defer moralisSnap.End()

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

// nolint:funlen // TODO
func GetNFTs(ctx context.Context, userAddress string, chainType ChainType, fromDate string, apiKey string) (NFTResult, error) {
	tracer := otel.Tracer(TracerNameCrawlerMoralis)

	ctx, getNFTListSnap := tracer.Start(ctx, "get_nft_list")
	getNFTListSnap.SetAttributes(
		attribute.String("user_address", userAddress),
		attribute.String("from_date", fromDate),
	)

	defer getNFTListSnap.End()

	// Gets all NFT items of user
	requestURL := fmt.Sprintf("%s/api/v2/%s/nft?chain=%s&format=decimal&from_date=%s",
		endpoint, userAddress, chainType, url.QueryEscape(fromDate),
	)

	response, err := requestMoralisApi(ctx, requestURL, apiKey, true)

	if err != nil {
		getNFTListSnap.RecordError(err)
		getNFTListSnap.SetStatus(codes.Error, err.Error())

		return NFTResult{}, err
	}

	res := new(NFTResult)
	SetMoralisAttributes(&res.MoralisAttributes, response)

	err = jsoni.Unmarshal(response.Body, &res)
	if err != nil {
		getNFTListSnap.RecordError(err)
		getNFTListSnap.SetStatus(codes.Error, err.Error())

		return NFTResult{}, err
	}

	getNFTListSnap.SetAttributes(
		attribute.Int("response_size", len(response.Body)),
	)

	var wg sync.WaitGroup

	for i, item := range res.Result {
		if item.MetaData == "" && item.TokenURI != "" {
			wg.Add(1)

			go func(item NFTItem, i int) {
				_, getNFTItemSnap := tracer.Start(ctx, "get_nft_item")

				defer func() {
					wg.Done()

					getNFTItemSnap.End()
				}()

				done := make(chan bool, 1)

				go func() {
					url := nft_utils.FormatUrl(item.TokenURI)

					if metadataRes, err := httpx.Get(url, nil); err != nil {
						getNFTListSnap.RecordError(err)
						getNFTListSnap.SetStatus(codes.Error, err.Error())

						logger.Warnf("http get nft metadata error with url '%s': [%v], moralis token uri: %v", url, err, item.TokenURI)
					} else {
						res.Result[i].MetaData = string(metadataRes.Body)
					}

					close(done)
				}()

				select {
				case <-done:
					return
				case <-time.After(time.Second * 5):
					return
				}
			}(item, i)
		}
	}

	wg.Wait()

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
	ctx, trace := otel.Tracer(TracerNameCrawlerMoralis).Start(ctx, "get_nft_transfer_list")
	trace.SetAttributes(
		attribute.String("user_address", userAddress),
		attribute.String("from_date", fromDate),
	)

	defer trace.End()

	// Gets all NFT transfers of user
	requestURL := fmt.Sprintf("%s/api/v2/%s/nft/transfers?chain=%s&from_block=%d&format=decimal&direction=both&from_date=%s",
		endpoint, userAddress, chainType, blockHeight, url.QueryEscape(fromDate),
	)
	response, err := requestMoralisApi(ctx, requestURL, apiKey, true)

	if err != nil {
		return NFTTransferResult{}, err
	}

	res := new(NFTTransferResult)
	SetMoralisAttributes(&res.MoralisAttributes, response)

	trace.SetAttributes(
		attribute.Int("response_size", len(response.Body)),
	)

	err = jsoni.Unmarshal(response.Body, &res)
	if err != nil {
		return NFTTransferResult{}, err
	}

	return *res, nil
}

func GetLogs(
	ctx context.Context, fromBlock int64, toBlock int64, address string, topic string, chainType ChainType, apiKey string,
) (GetLogsResult, error) {
	tracer := otel.Tracer(TracerNameCrawlerMoralis)

	ctx, getLogListSnap := tracer.Start(ctx, "get_erc20_token_metadata_data_from_url")

	defer getLogListSnap.End()

	url := fmt.Sprintf(
		"%s/api/v2/%s/logs?chain=%s&from_block=%d&to_block=%d&topic0=%s",
		endpoint, address, string(chainType), fromBlock, toBlock, topic,
	)

	response, err := requestMoralisApi(ctx, url, apiKey, false)

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
func GetNFTByContract(ctx context.Context, userAddress string, contactAddress string, chainType ChainType, apiKey string) (NFTResult, error) {
	tracer := otel.Tracer(TracerNameCrawlerMoralis)

	ctx, getNFTByContractSnap := tracer.Start(ctx, "get_nft_by_contract")

	defer getNFTByContractSnap.End()

	// this function is used by ENS indexer.
	url := fmt.Sprintf("%s/api/v2/%s/nft?chain=%s&format=decimal&token_addresses=%s",
		endpoint, userAddress, chainType, contactAddress)

	response, err := requestMoralisApi(ctx, url, apiKey, true)

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
func GetTxByToken(ctx context.Context, tokenAddress string, tokenId string, chainType ChainType, apiKey string) (TransferItem, error) {
	tracer := otel.Tracer(TracerNameCrawlerMoralis)

	ctx, getTxByTokenSnap := tracer.Start(ctx, "get_tx_by_token")

	defer getTxByTokenSnap.End()

	url := fmt.Sprintf("%s/api/v2/nft/%s/%s/transfers?chain=%s&format=decimal&limit=1",
		endpoint, tokenAddress, tokenId, chainType)
	response, err := requestMoralisApi(ctx, url, apiKey, true)

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

func GetMetadataByToken(ctx context.Context, tokenAddress string, tokenId string, chainType ChainType, apiKey string) (NFTItem, error) {
	tracer := otel.Tracer(TracerNameCrawlerMoralis)

	ctx, getMetadataByTokenSnap := tracer.Start(ctx, "get_metadata_by_token")

	defer getMetadataByTokenSnap.End()

	url := fmt.Sprintf("%s/api/v2/nft/%s/%s?chain=%s&format=decimal&limit=1",
		endpoint, tokenAddress, tokenId, chainType)
	response, err := requestMoralisApi(ctx, url, apiKey, true)

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
	ctx, trace := otel.Tracer(TracerNameCrawlerMoralis).Start(ctx, "get_erc20_transfer_list")
	trace.SetAttributes(
		attribute.String("from_date", fromDate),
	)

	defer trace.End()

	var (
		transferItems = make([]ERC20TransferItem, 0)
		count         = 0
		cursor        = ""
	)

	for {
		// Pull up to 500 at a time
		if count == 5 {
			break
		}

		transfer, err := getErc20Once(ctx, userAddress, chainType, fromDate, apiKey, cursor)
		if err != nil {
			logger.Errorf("get erc20 once error: %v", err)

			return nil, err
		}

		transferItems = append(transferItems, transfer.Result...)

		if transfer.Cursor == "" {
			break
		}

		cursor = transfer.Cursor

		count += 1
	}

	return transferItems, nil
}

func getErc20Once(
	ctx context.Context,
	userAddress string,
	chainType ChainType,
	fromDate string,
	apiKey string,
	cursor string) (*ERC20Transfer, error) {
	_, trace := otel.Tracer(TracerNameCrawlerMoralis).Start(ctx, "get_erc20_once")

	trace.SetAttributes(
		attribute.String("from_date", fromDate),
	)

	defer trace.End()

	// requestURL := fmt.Sprintf("%s/api/v2/%s/erc20/transfers?chain=%s&from_block=%d&offset=%d&from_date=%s",
	// 	endpoint, userAddress, chainType, 0, offest, url.QueryEscape(fromDate),
	// )

	requestURL := fmt.Sprintf("%s/api/v2/%s/erc20/transfers?chain=%s&from_block=%d&from_date=%s&cursor=%s",
		endpoint, userAddress, chainType, 0, url.QueryEscape(fromDate), cursor)

	response, err := requestMoralisApi(ctx, requestURL, apiKey, true)

	if err != nil {
		return nil, err
	}

	trace.SetAttributes(
		attribute.Int("response_size", len(response.Body)),
	)

	res := new(ERC20Transfer)
	SetMoralisAttributes(&res.MoralisAttributes, response)

	if err = jsoni.Unmarshal(response.Body, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func GetErc20TokenMetaData(ctx context.Context, chainType ChainType, addresses []string, apiKey string) (Erc20TokensMap, error) {
	tracer := otel.Tracer(TracerNameCrawlerMoralis)

	ctx, getEER20TokenMetadataSnap := tracer.Start(ctx, "get_erc20_token_metadata")

	defer getEER20TokenMetadataSnap.End()

	if len(addresses) <= 0 {
		return Erc20TokensMap{}, fmt.Errorf("addresss is empty")
	}

	res := Erc20TokensMap{}

	getErc20TokenMetaDataFromCache(addresses, res)

	addrLen := len(addresses)

	if len(res) == addrLen {
		return res, nil
	}

	limit := 150

	addressBatch := make([][]string, 0)

	for i := 0; i < addrLen; i += limit {
		addressBatch = append(addressBatch, addresses[i:Min(i+limit, addrLen)])
	}

	lop.ForEach(addressBatch, func(batch []string, i int) {
		batchRes := Erc20TokensMap{}
		getErc20TokenMetaDataFromUrl(ctx, chainType, batch, apiKey, batchRes)

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

func getErc20TokenMetaDataFromUrl(ctx context.Context, chainType ChainType, addresses []string, apiKey string, res Erc20TokensMap) error {
	tracer := otel.Tracer(TracerNameCrawlerMoralis)

	ctx, getEER20TokenMetadataFromURLSnap := tracer.Start(ctx, "get_erc20_token_metadata_data_from_url")

	defer getEER20TokenMetadataFromURLSnap.End()

	url := fmt.Sprintf("%s/api/v2/erc20/metadata?chain=%s",
		endpoint, chainType)

	for _, address := range addresses {
		url += fmt.Sprintf("&addresses=%s", address)
	}

	response, err := requestMoralisApi(ctx, url, apiKey, true)

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
func GetENSList(ctx context.Context, address string) ([]ENSTextRecord, error) {
	tracer := otel.Tracer(TracerNameCrawlerMoralis)

	ctx, getENSListSnap := tracer.Start(ctx, "get_ens_list")

	defer getENSListSnap.End()

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

	err = getENSDetail(ctx, address, &record)

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
func getENSDetail(ctx context.Context, address string, record *ENSTextRecord) error {
	tracer := otel.Tracer(TracerNameCrawlerMoralis)

	ctx, getENSDetailSnap := tracer.Start(ctx, "get_erc20_token_metadata_data_from_url")

	defer getENSDetailSnap.End()

	ensList, err := GetNFTByContract(ctx, address, ensContract, ETH, getApiKey())

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

				return getENSTransaction(ctx, ens, record)
			}
		}
	}

	return nil
}

func getENSTransaction(ctx context.Context, ens NFTItem, record *ENSTextRecord) error {
	tracer := otel.Tracer(TracerNameCrawlerMoralis)

	ctx, getENSTransactionSnap := tracer.Start(ctx, "get_ens_transaction")

	defer getENSTransactionSnap.End()

	// get TxHash and Tsp with TokenId from Moralis
	t, err := GetTxByToken(ctx, ens.TokenAddress, ens.TokenId, ETH, getApiKey())

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
	ctx, trace := otel.Tracer(TracerNameCrawlerMoralis).Start(ctx, "get_eth_transfer_list")
	trace.SetAttributes(
		attribute.String("user_address", userAddress),
		attribute.String("from_date", fromDate),
	)

	defer trace.End()

	var (
		transferItems = make([]ETHTransferItem, 0)
		count         = 0
		cursor        = ""
	)

	for {
		// Pull up to 500 at a time
		if count == 5 {
			break
		}

		transfer, err := getETHOnce(ctx, userAddress, chainType, fromDate, apiKey, cursor)
		if err != nil {
			logger.Errorf("get eth once error: %v", err)

			break
		}

		transferItems = append(transferItems, transfer.Result...)

		if transfer.Cursor == "" {
			break
		}

		cursor = transfer.Cursor

		count += 1
	}

	return transferItems, nil
}

func getETHOnce(
	ctx context.Context,
	userAddress string,
	chainType ChainType,
	fromDate string,
	apiKey string,
	cursor string) (*ETHTransfer, error) {
	ctx, trace := otel.Tracer(TracerNameCrawlerMoralis).Start(ctx, "get_eth_once")
	trace.SetAttributes(
		attribute.String("user_address", userAddress),
		attribute.String("from_date", fromDate),
	)

	defer trace.End()

	requestURL := fmt.Sprintf("%s/api/v2/%s?chain=%s&from_block=%d&from_date=%s&cursor=%s",
		endpoint, userAddress, chainType, 0, url.QueryEscape(fromDate), cursor)

	response, err := requestMoralisApi(ctx, requestURL, apiKey, true)

	if err != nil {
		return nil, err
	}

	trace.SetAttributes(
		attribute.Int("response_size", len(response.Body)),
	)

	res := new(ETHTransfer)
	SetMoralisAttributes(&res.MoralisAttributes, response)

	if err = jsoni.Unmarshal(response.Body, &res); err != nil {
		return nil, err
	}

	return res, nil
}
