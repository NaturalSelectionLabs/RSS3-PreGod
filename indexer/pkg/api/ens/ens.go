package ens

import (
	"fmt"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/httpx"
	"github.com/ethereum/go-ethereum/ethclient"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fastjson"
	goens "github.com/wealdtech/go-ens/v3"
)

var (
	jsoni    = jsoniter.ConfigCompatibleWithStandardLibrary
	endpoint = "https://api.thegraph.com/subgraphs/name/ensdomains/ens"
	parser   fastjson.Parser
	client   *ethclient.Client
)

func getClient() {
	gateway := config.Config.Indexer.Infura.Gateway + config.Config.Indexer.Infura.ApiKey
	c, err := ethclient.Dial(gateway)

	if err != nil {

	}
	client = c
}

// returns a list of ENS domains with non-empty text records
func GetENSList(address string) []ENSTextRecord {
	getClient()

	result := []ENSTextRecord{}

	parsedJson := sendTheGraphRequest(`{
		domains( where : { owner:"` + address + `" } ) {
			name
			createdAt
			resolver {
				texts
			}
			events{
				transactionID
			}
		}
	}`)

	parsedList := parsedJson.GetArray("data", "domains")

	for _, item := range parsedList {
		domain := string(item.GetStringBytes("name"))
		record := ENSTextRecord{
			domain:    domain,
			createdAt: time.Unix(item.GetInt64("createdAt"), 0),
			txHash:    string(item.GetStringBytes("events", "transactionID")),
		}

		texts := item.GetArray("resolver", "texts")
		if len(texts) > 0 {
			record.text = make(map[string]string, len(texts))
			for _, text := range texts {
				key := string(text.GetStringBytes())
				record.text[key] = getENSTextValue(domain, key)
			}
		}

		result = append(result, record)
	}
	return result
}

// returns the text record value for a given ENS and text record key
func getENSTextValue(domain string, text string) string {
	r, err := goens.NewResolver(client, domain)

	if err != nil {
		return ""
	}

	t, err := r.Text(text)

	if err != nil {
		return ""
	}

	return t
}

// sends a request to the TheGraph's graphql endpoint
func sendTheGraphRequest(query string) *fastjson.Value {
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	queryMarshalled, _ := jsoni.MarshalToString(map[string]string{
		"query": query,
	})

	response, err := httpx.PostRaw(endpoint, headers, queryMarshalled)

	if err != nil {
		fmt.Printf("Thegraph API request err: %v", err)
		return nil
	}

	result, err := parser.Parse(string(response.Body()))

	if err != nil {
		fmt.Printf("Thegraph API invalid response: %v", err)
		return nil
	}

	return result
}
