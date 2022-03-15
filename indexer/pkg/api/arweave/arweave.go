package arweave

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/httpx"
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

// GetContentByTxHash gets transaction content by tx hash.
func GetContentByTxHash(hash string) ([]byte, error) {
	var headers = map[string]string{
		"Origin":  "https://viewblock.io",
		"Referer": "https://viewblock.io",
	}

	return httpx.Get(arweaveEndpoint+"/"+hash, headers)
}

// GetTransactions gets all transactions using filters.
func GetTransactions(from, to uint64, owner string) ([]byte, error) {
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

// GetArticles gets all articles from arweave using filters.
func GetArticles(from, to uint64, owner string) ([]MirrorArticle, error) {
	response, err := GetTransactions(from, to, owner)
	if err != nil {
		return nil, nil
	}

	var parser fastjson.Parser

	parsedJson, parseErr := parser.Parse(string(response))
	if parseErr != nil {
		return nil, nil
	}

	// edges
	edges := parsedJson.GetArray("data", "transactions", "edges")
	result := make([]MirrorArticle, len(edges))

	for i := 0; i < len(edges); i++ {
		result[i], err = parseGraphqlNode(edges[i].String())
		if err != nil {
			return nil, nil
		}
	}

	return result, nil
}

func parseGraphqlNode(node string) (MirrorArticle, error) {
	var parser fastjson.Parser

	parsedJson, err := parser.Parse(node)
	if err != nil {
		return MirrorArticle{}, err
	}

	article := MirrorArticle{}

	tags := parsedJson.GetArray("node", "tags")
	for _, tag := range tags {
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

		article.Link = mirrorHost + "/" + article.Author + "/" + article.OriginalDigest
	}

	id := parsedJson.GetStringBytes("node", "id")
	content, err := GetContentByTxHash(string(id))

	if err != nil {
		return article, err
	}

	parsedJson, err = parser.Parse(string(content))
	if err != nil {
		return article, err
	}

	// title
	article.Title = string(parsedJson.GetStringBytes("content", "title"))
	// timestamp
	article.TimeStamp = parsedJson.GetInt64("content", "timestamp")
	// content
	article.Content = string(parsedJson.GetStringBytes("content", "body")) // timestamp

	return article, nil
}
