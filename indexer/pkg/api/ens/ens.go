package ens

import (
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/moralis"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	goens "github.com/wealdtech/go-ens/v3"
)

var (
	client      *ethclient.Client
	ensContract = "0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85"
)

func getClient() {
	gateway := config.Config.Indexer.Infura.Gateway + config.Config.Indexer.Infura.ApiKey
	c, err := ethclient.Dial(gateway)

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

	for _, key := range getTextRecordKeyList() {
		t, err := r.Text(key)

		if err != nil {
			logger.Errorf("getENSTextValue read text: %v", err)

			return err
		}

		record.Text[key] = t
	}

	return nil
}

// returns ENS details from moralis
func getENSDetail(address string, record *ENSTextRecord) error {
	ensList, err := moralis.GetNFTByContract(address, ensContract, moralis.ETH, config.Config.Indexer.Moralis.ApiKey)

	if err != nil {
		logger.Errorf("getENSDetail GetNFTByContract: %v", err)

		return err
	}

	// an address might have multiple ENS domains
	for _, ens := range ensList.Result {
		// if the one is the current ENS domain
		if ens.MetaData.Name == record.Domain {
			// get TxHash and Tsp with TokenId from Moralis
			t, err := moralis.GetTxByToken(ens.TokenAddress, ens.TokenId, moralis.ETH, config.Config.Indexer.Moralis.ApiKey)

			if err != nil {
				logger.Errorf("getENSDetail transaction: %v", err)

				return err
			}

			record.TxHash = t.TransactionHash
			record.CreatedAt, err = time.Parse(time.RFC3339, t.BlockTimestamp)
			record.Description = ens.MetaData.Description

			if err != nil {
				logger.Errorf("getENSDetail transaction: %v", err)

				return err
			}

			return nil
		}
	}

	return nil
}
