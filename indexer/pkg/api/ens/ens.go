package ens

import (
	"strings"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/moralis"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/datatype"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/httpx"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	jsoniter "github.com/json-iterator/go"
	goens "github.com/wealdtech/go-ens/v3"
)

var (
	jsoni         = jsoniter.ConfigCompatibleWithStandardLibrary
	client        *ethclient.Client
	ensContract   = "0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85"
	infuraGateway = "https://mainnet.infura.io/v3"
)

func getClient() {
	c, err := ethclient.Dial(infuraGateway + "/" + config.Config.Indexer.Infura.ApiKey)

	if err != nil {
		logger.Errorf("connect to Infura: %v", err)
	}

	client = c
}

// returns a list of ENS domains with non-empty text records
func GetENSList(address string) ([]ENSTextRecord, error) {
	getClient()

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

// reads the text record value for a given domain with the list of predefined keys from infura
func getENSTextValue(domain string, record *ENSTextRecord) error {
	r, err := goens.NewResolver(client, domain)

	if err != nil {
		logger.Errorf("getENSTextValue NewResolver: %v", err)

		return err
	}

	record.Text = make(map[string]string)

	var attachments datatype.Attachments

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
			a := datatype.Attachment{
				Type:     "websites",
				MimeType: "text/uri-list",
				Content:  text,
			}
			attachments = append(attachments, a)
		case "avatar":
			// only get content headers if it's http for now
			if strings.HasPrefix(text, "http") {
				contentHeader, err := httpx.GetContentHeader(text)

				if err != nil {
					logger.Errorf("GetContentHeader err: %v", err)
				}

				a := datatype.Attachment{
					Type:     "banner",
					MimeType: contentHeader.MIMEType,
					Content:  text,
				}
				attachments = append(attachments, a)
			}
		}
	}

	record.Attachments = attachments

	return nil
}

// returns ENS details from moralis
func getENSDetail(address string, record *ENSTextRecord) error {
	ensList, err := moralis.GetNFTByContract(address, ensContract, moralis.ETH, config.Config.Indexer.Moralis.ApiKey)

	if err != nil {
		logger.Errorf("getENSDetail GetNFTByContract: %v", err)

		return err
	}

	for _, ens := range ensList.Result {
		// moralis sometimes returns empty metadata
		if ens.MetaData != "" {
			meta := new(moralis.NFTMetadata)

			err = jsoni.UnmarshalFromString(ens.MetaData, &meta)

			if err != nil {
				logger.Errorf("getENSDetail unmarshall metadata: %v", err)

				return err
			}

			// an address might have multiple ENS domains
			// if the one is the current ENS domain
			if meta.Name == record.Domain {
				record.Description = meta.Description

				return getENSTransaction(ens, record)
			}
		}
	}

	return nil
}

func getENSTransaction(ens moralis.NFTItem, record *ENSTextRecord) error {
	// get TxHash and Tsp with TokenId from Moralis
	t, err := moralis.GetTxByToken(ens.TokenAddress, ens.TokenId, moralis.ETH, config.Config.Indexer.Moralis.ApiKey)

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
