package moralis_test

import (
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/moralis"
	"github.com/stretchr/testify/assert"
)

func TestGetNFT(t *testing.T) {
	t.Parallel()

	apiKey := moralis.GetApiKey()
	result, err := moralis.GetNFTs("0x3b6d02a24df681ffdf621d35d70aba7adaac07c1", "eth", apiKey)
	// assert for nil
	assert.Nil(t, err)

	assert.NotEmpty(t, result.Result)
}

func TestGetNFTTransfers(t *testing.T) {
	t.Parallel()

	apiKey := moralis.GetApiKey()
	result, err := moralis.GetNFTTransfers("0x3b6d02a24df681ffdf621d35d70aba7adaac07c1", "eth", apiKey)
	// assert for nil
	assert.Nil(t, err)

	assert.NotEmpty(t, result.Result)
}

func TestGetLogs(t *testing.T) {
	t.Parallel()

	apiKey := moralis.GetApiKey()

	result, err := moralis.GetLogs(
		12605342,
		12605343,
		"0x7d655c57f71464B6f83811C55D84009Cd9f5221C",
		"0x3bb7428b25f9bdad9bd2faa4c6a7a9e5d5882657e96c1d24cc41c1d6c1910a98",
		"eth",
		apiKey)
	// assert for nil
	assert.Nil(t, err)
	assert.NotEmpty(t, result.Result)

	for _, item := range result.Result {
		assert.NotEmpty(t, item)
	}
}
