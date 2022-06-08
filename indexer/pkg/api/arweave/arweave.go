package arweave

import (
	"fmt"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/httpx"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	jsoniter "github.com/json-iterator/go"
)

const arweaveEndpoint string = "https://arweave.net"
const arweaveGraphqlEndpoint string = "https://arweave.net/graphql"
const mirrorHost = "https://mirror.xyz/"

var (
	jsoni = jsoniter.ConfigCompatibleWithStandardLibrary
)

// GetLatestBlockHeight gets the latest block height for arweave
func GetLatestBlockHeight() (int64, error) {
	response, err := httpx.Get(arweaveEndpoint, nil)
	if err != nil {
		return 0, err
	}

	latestBlockResult := new(ArLatestBlockResult)
	if err := jsoni.UnmarshalFromString(string(response.Body), latestBlockResult); err != nil {
		logger.Errorf("arweave GetLatestBlockHeight unmarshalFromString error: %v", err)

		return 0, err
	}

	return latestBlockResult.Height, nil
}

func GetLatestBlockHeightWithConfirmations(confirmations int64) (int64, error) {
	// get latest block height
	latestBlockHeight, err := GetLatestBlockHeight()
	if err != nil {
		return 0, err
	}

	return latestBlockHeight - confirmations, nil
}

// GetContentByTxHash gets transaction content by tx hash.
func GetContentByTxHash(hash string) ([]byte, error) {
	var headers = map[string]string{
		"Origin":  "https://viewblock.io",
		"Referer": "https://viewblock.io",
	}

	resp, err := httpx.Get(arweaveEndpoint+"/"+hash, headers)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

// GetTransactions gets all transactions using filters.
func GetTransactions(from, to int64, owner ArAccount, cursor string) ([]byte, error) {
	var headers = map[string]string{
		"content-type": "application/json",
		"accept":       "*/*",
		"origin":       "https://arweave.net",
		"referer":      "https://arweave.net/graphql",
		"user-agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.3",
	}

	queryString := `query {
		transactions(
			block: { min: %d, max: %d }
			owners: ["%s"]
			sort: HEIGHT_ASC
			first: 100
			after: "%s"
		) {
			pageInfo {
				hasNextPage
			}
			edges {
				cursor
				node {
					id
					tags {
						name
						value
					}
				}
			}
		}
	}`

	data := map[string]string{
		"query": fmt.Sprintf(queryString, from, to, owner, cursor),
	}

	json, _ := jsoni.MarshalToString(data)

	resp, err := httpx.Post(arweaveGraphqlEndpoint, headers, json)

	return resp.Body, err
}

// GetMirrorContents gets all articles from arweave using filters.
func GetMirrorContents(from, to int64, owner ArAccount) ([]MirrorContent, error) {
	lastCursor := ""
	results := make([]MirrorContent, 0)

	for {
		response, err := GetTransactions(from, to, owner, lastCursor)
		if err != nil {
			logger.Errorf("GetTransactions error: [%v]", err)

			return nil, nil
		}

		graphqlResult := new(GraphqlResult)
		if err := jsoni.UnmarshalFromString(string(response), graphqlResult); err != nil {
			logger.Errorf("arweave unmarshalFromString error: %v for response %s", err, string(response))

			return nil, err
		}
		// edges
		edges := graphqlResult.Data.Transactions.Edges
		l := len(edges)

		hasNextPage := graphqlResult.Data.Transactions.PageInfo.HasNextPage

		for i := 0; i < l; i++ {
			res, err := parseGraphqlNode(edges[i])
			if err != nil {
				logger.Errorf("parseGraphqlNode error %v for edge", err, edges[i])

				continue
			}

			if res != nil {
				results = append(results, *res)
			}
		}

		if !hasNextPage || l == 0 {
			break
		}

		lastCursor = edges[l-1].Cursor

		time.Sleep(DefaultCrawlConfig.sleepInterval)
	}

	logger.Infof("Getting transactions from [%d] to [%d], [%d] results in total. ", from, to, len(results))

	return results, nil
}

func parseGraphqlNode(node GraphqlResultEdges) (*MirrorContent, error) {
	article := new(MirrorContent)

	var appName string

	tags := node.Node.Tags
	for _, tag := range tags {
		switch tag.Name {
		case "App-Name":
			appName = tag.Value
		case "Contributor":
			article.Author = tag.Value
		case "Content-Digest":
			article.Digest = tag.Value
		case "Original-Content-Digest":
			article.OriginalDigest = tag.Value
		}
	}

	// only parse tags with "MirrorXYZ"
	if appName != "MirrorXYZ" {
		return nil, nil
	}

	id := node.Node.Id
	if id == "" {
		return nil, nil
	}

	content, err := GetContentByTxHash(id)

	if err != nil {
		logger.Errorf("GetContentByTxHash error: [%v]", err)

		return nil, err
	}

	originalMirrorContent := new(OriginalMirrorContent)
	if err := jsoni.UnmarshalFromString(string(content), &originalMirrorContent); err != nil {
		logger.Errorf("arweave unmarshalFromString error: %v", err)

		return nil, err
	}

	// title
	article.Title = originalMirrorContent.Content.Title
	// timestamp
	article.Timestamp = originalMirrorContent.Content.Timestamp
	// content
	article.Content = originalMirrorContent.Content.Body
	// txHash
	article.TxHash = string(id)
	// Author
	article.Author = originalMirrorContent.Authorship.Contributor
	// parse Content-Digest
	article.Digest = originalMirrorContent.Digest
	// parse OriginalDigest
	article.OriginalDigest = originalMirrorContent.OriginalDigest
	// Link
	if article.Author != "" && article.OriginalDigest != "" {
		article.Link = mirrorHost + "/" + article.Author + "/" + article.OriginalDigest
	}

	return article, nil
}
