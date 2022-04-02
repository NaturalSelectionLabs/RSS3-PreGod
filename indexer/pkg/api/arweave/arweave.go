package arweave

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/httpx"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/valyala/fastjson"
)

const arweaveEndpoint string = "https://arweave.net"
const arweaveGraphqlEndpoint string = "https://arweave.net/graphql"
const mirrorHost = "https://mirror.xyz/"

// GetLatestBlockHeight gets the latest block height for arweave
func GetLatestBlockHeight() (int64, error) {
	response, err := httpx.Get(arweaveEndpoint, nil)
	if err != nil {
		return 0, nil
	}

	var parser fastjson.Parser
	parsedJson, parseErr := parser.Parse(string(response))

	if parseErr != nil {
		return 0, nil
	}

	blockHeight := parsedJson.GetInt64("height")

	return blockHeight, nil
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

	return httpx.Get(arweaveEndpoint+"/"+hash, headers)
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

	return httpx.Post(arweaveGraphqlEndpoint, headers, data)
}

// GetMirrorContents gets all articles from arweave using filters.
func GetMirrorContents(from, to int64, owner ArAccount) ([]MirrorContent, error) {
	response, err := GetTransactions(from, to, owner)
	if err != nil {
		logger.Errorf("GetTransactions error: [%v]", err)

		return nil, nil
	}

	//log.Println(string(response))

	var parser fastjson.Parser

	parsedJson, parseErr := parser.Parse(string(response))
	if parseErr != nil {
		logger.Errorf("Parse json  error: [%v]", parseErr)

		return nil, nil
	}

	// edges
	edges := parsedJson.GetArray("data", "transactions", "edges")
	result := make([]MirrorContent, len(edges))

	for i := 0; i < len(edges); i++ {
		result[i], err = parseGraphqlNode(edges[i].String())
		if err != nil {
			return nil, nil
		}
	}

	return result, nil
}

func parseGraphqlNode(node string) (MirrorContent, error) {
	var parser fastjson.Parser

	parsedJson, err := parser.Parse(node)
	if err != nil {
		return MirrorContent{}, err
	}

	article := MirrorContent{}

	tags := parsedJson.GetArray("node", "tags")
	for _, tag := range tags {
		// only parse tags with "MirrorXYZ"
		appName := string(tag.GetStringBytes("App-Name"))
		if appName != "MirrorXYZ" {
			continue
		}

		name := string(tag.GetStringBytes("name"))
		value := string(tag.GetStringBytes("value"))

		switch name {
		case "Contributor":
			article.Author = value
		case "Content-Digest":
			article.Digest = value
		case "Original-Content-Digest":
			article.OriginalDigest = value
		}
	}

	id := parsedJson.GetStringBytes("node", "id")
	content, err := GetContentByTxHash(string(id))

	if err != nil {
		logger.Errorf("GetContentByTxHash error: [%v]", err)

		return article, err
	}

	//log.Println(string(content))

	parsedJson, err = parser.Parse(string(content))
	if err != nil {
		return article, err
	}

	// title
	article.Title = string(parsedJson.GetStringBytes("content", "title"))
	// timestamp
	article.Timestamp = parsedJson.GetInt64("content", "timestamp")
	// content
	article.Content = string(parsedJson.GetStringBytes("content", "body")) // timestamp
	// txHash
	article.TxHash = string(id)

	// parse digest and author, in case of empty fields in tags, eg: arweave block #591647
	// Author
	if article.Author == "" {
		article.Author = string(parsedJson.GetStringBytes("authorship", "contributor"))
	}
	// parse Content-Digest
	if article.Digest == "" {
		article.Digest = string(parsedJson.GetStringBytes("digest"))
	}
	// parse OriginalDigest
	if article.OriginalDigest == "" {
		article.OriginalDigest = string(parsedJson.GetStringBytes("originalDigest"))
	}
	// Link
	if article.Author != "" && article.OriginalDigest != "" {
		article.Link = mirrorHost + "/" + article.Author + "/" + article.OriginalDigest
	}

	return article, nil
}
