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

func requestMoralisApi(url string, apiKey string) ([]byte, error) {
	var headers = map[string]string{
		"accept":    "application/json",
		"X-API-Key": apiKey,
	}

	response, err := httpx.Get(url, headers)
	if err != nil {
		logger.Errorf("http get error with url '%s': [%v]. response: %s",
			url, err, string(response))

		return nil, err
	}

	return response, nil
}

func GetNFTs(userAddress string, chainType ChainType, apiKey string) (NFTResult, error) {
	// Gets all NFT items of user
	url := fmt.Sprintf("%s/api/v2/%s/nft?chain=%s&format=decimal",
		endpoint, userAddress, chainType)

	response, err := requestMoralisApi(url, apiKey)

	if err != nil {
		return NFTResult{}, err
	}

	res := new(NFTResult)

	err = jsoni.Unmarshal(response, &res)
	if err != nil {
		return NFTResult{}, err
	}

	lop.ForEach(res.Result, func(item NFTItem, i int) {
		if item.MetaData == "" && item.TokenURI != "" {
			if metadataRes, err := httpx.Get(nft_utils.FormatUrl(item.TokenURI), nil); err != nil {
				logger.Warnf("http get nft metadata error with url '%s': [%v]", item.TokenURI, err)
			} else {
				res.Result[i].MetaData = string(metadataRes)
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

	err = jsoni.Unmarshal(response, &res)
	if err != nil {
		return NFTTransferResult{}, err
	}

	return *res, nil
}

func GetLogs(fromBlock int64, toBlock int64, address string, topic string, chainType ChainType, apiKey string) (*GetLogsResult, error) {
	url := fmt.Sprintf("%s/api/v2/%s/logs?chain=%s&from_block=%d&to_block=%d&topic0=%s",
		endpoint, address, string(chainType), fromBlock, toBlock, topic)
	response, err := requestMoralisApi(url, apiKey)

	if err != nil {
		return nil, err
	}

	res := new(GetLogsResult)

	err = jsoni.Unmarshal(response, &res)
	if err != nil {
		logger.Errorf("unmarshal error: [%v]", err)

		return nil, err
	}

	return res, nil
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

	err = jsoni.Unmarshal(response, &res)
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

	parsedJson, err := parser.Parse(string(response))
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

	err = jsoni.Unmarshal(response, &res)
	if err != nil {
		return NFTItem{}, nil
	}

	return *res, nil
}

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
