package arweave_test

import (
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/arweave"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fastjson"
)

func TestGetLatestBlockHeight(t *testing.T) {
	blockHeight, err := arweave.GetLatestBlockHeight()
	assert.Nil(t, err)
	assert.NotEqual(t, 0, blockHeight)
}

func TestGetContentByTxHash(t *testing.T) {
	hash := "BhM-D1bsQkaqi72EEG1aRVs4Nv5bZZIW-mH8yEdIDWA"
	response, err := arweave.GetContentByTxHash(hash)
	// assert for nil
	assert.Nil(t, err)
	assert.NotEmpty(t, response)

	var parser fastjson.Parser
	parsedJson, parseErr := parser.Parse(string(response))
	assert.Nil(t, parseErr)

	// check title
	title := parsedJson.GetStringBytes("content", "title")
	assert.NotEmpty(t, title)

	// check body
	body := parsedJson.GetStringBytes("content", "body")
	assert.NotEmpty(t, body)

	// check contributor
	contributor := parsedJson.GetStringBytes("authorship", "contributor")
	assert.NotEmpty(t, contributor)
	// check originalDigest
	originalDigest := parsedJson.GetStringBytes("originalDigest")
	assert.NotEmpty(t, originalDigest)
}

func TestGetTransacitons(t *testing.T) {
	response, err := arweave.GetTransactions(877250, 877250, arweave.MirrorUploader, "")
	// assert for nil
	assert.Nil(t, err)
	assert.NotEmpty(t, response)

	var parser fastjson.Parser
	parsedJson, parseErr := parser.Parse(string(response))
	assert.Nil(t, parseErr)

	edges := parsedJson.GetArray("data", "transactions", "edges")
	assert.NotEmpty(t, edges)
}

func TestGetArticles(t *testing.T) {
	articles, err := arweave.GetMirrorContents(877250, 877250, arweave.MirrorUploader)
	// assert for nil
	assert.Nil(t, err)
	assert.NotEmpty(t, articles)

	for _, article := range articles {
		assert.NotEmpty(t, article.Title)
		assert.NotEmpty(t, article.Timestamp)
		assert.NotEmpty(t, article.Content)
		assert.NotEmpty(t, article.Author)
		assert.NotEmpty(t, article.Link)
		assert.NotEmpty(t, article.Digest)
		assert.NotEmpty(t, article.OriginalDigest)
	}
}
