package arweave

import (
	"fmt"

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
	if err := jsoni.UnmarshalFromString(string(response.RespBody), latestBlockResult); err != nil {
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

	return resp.RespBody, nil
}

// GetTransactions gets all transactions using filters.
func GetTransactions(from, to int64, owner ArAccount) ([]byte, error) {
	var headers = map[string]string{
		"Accept-Encoding": "gzip, deflate, br",
		"Content-Type":    "application/json",
		"Accept":          "application/json",
		"Origin":          "https://arweave.net",
	}

	queryVariables :=
		"{\"query\":\"query { transactions( " +
			"block: { min: %d, max: %d } " +
			"owners: [\\\"%s\\\"] " +
			"sort: HEIGHT_ASC ) { edges { node {id tags { name value } } } }" +
			"}\"}"
	data := fmt.Sprintf(queryVariables, from, to, owner)

	resp, err := httpx.Post(arweaveGraphqlEndpoint, headers, data)
	if err != nil {
		return nil, err
	}

	return httpx.Post(arweaveGraphqlEndpoint, headers, data)
}

// GetMirrorContents gets all articles from arweave using filters.
func GetMirrorContents(from, to int64, owner ArAccount) ([]MirrorContent, error) {
	response, err := GetTransactions(from, to, owner)
	if err != nil {
		logger.Errorf("GetTransactions error: [%v]", err)

		return nil, nil
	}

	graphqlResult := new(GraphqlResult)
	if err := jsoni.UnmarshalFromString(string(response), graphqlResult); err != nil {
		logger.Errorf("arweave unmarshalFromString error: %v", err)

		return nil, err
	}
	// edges
	edges := graphqlResult.Data.Transactions.Edges
	results := make([]MirrorContent, 0)

	for i := 0; i < len(edges); i++ {
		res, err := parseGraphqlNode(edges[i])
		if err != nil {
			return nil, nil
		}

		if res != nil {
			results = append(results, *res)
		}
	}

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
