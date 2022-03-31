package moralis_test

import (
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/moralis"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	_ "github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/stretchr/testify/assert"
)

var (
	tokenId      = "69122868356010038918278537874891361194569907163152093427587761621557332847656"
	tokenAddress = "0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85"
	ensContract  = "0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85"
)

func TestGetNFT(t *testing.T) {
	t.Parallel()

	result, err := moralis.GetNFTs("0x3b6d02a24df681ffdf621d35d70aba7adaac07c1", "eth", config.Config.Indexer.Moralis.ApiKey)
	assert.NotEmpty(t, result.Result)
	// assert for nil
	assert.Nil(t, err)
}

func TestGetNFTTransfers(t *testing.T) {
	t.Parallel()

	result, err := moralis.GetNFTTransfers("0x3b6d02a24df681ffdf621d35d70aba7adaac07c1", "eth", config.Config.Indexer.Moralis.ApiKey)

	assert.NotEmpty(t, result.Result)
	// assert for nil
	assert.Nil(t, err)
}

func TestGetLogs(t *testing.T) {
	t.Parallel()

	result, err := moralis.GetLogs(
		12605342,
		12605343,
		"0x7d655c57f71464B6f83811C55D84009Cd9f5221C",
		"0x3bb7428b25f9bdad9bd2faa4c6a7a9e5d5882657e96c1d24cc41c1d6c1910a98",
		"eth",
		config.Config.Indexer.Moralis.ApiKey)
	// assert for nil
	assert.Nil(t, err)
	assert.NotEmpty(t, result.Result)

	for _, item := range result.Result {
		assert.NotEmpty(t, item)
	}
}

func TestGetTxByToken(t *testing.T) {
	t.Parallel()

	result, err := moralis.GetTxByToken(
		tokenAddress, tokenId,
		"eth",
		config.Config.Indexer.Moralis.ApiKey)

	assert.Equal(t, result.TransactionHash, "0x44ea5a47fa51ada626874ac5c243e78ee485e354d5b337ea673d7f117eb8b6c3")
	assert.Equal(t, result.BlockTimestamp, "2022-01-02T11:16:35.000Z")

	assert.Nil(t, err)
}

func TestGetNFTByContract(t *testing.T) {
	t.Parallel()

	result, err := moralis.GetNFTByContract(
		"0x827431510a5D249cE4fdB7F00C83a3353F471848", ensContract,
		"eth",
		config.Config.Indexer.Moralis.ApiKey)

	assert.Equal(t, len(result.Result), 1)

	ens := result.Result[0]

	assert.Equal(t, ens.TokenAddress, tokenAddress)
	assert.Equal(t, ens.TokenId, tokenId)

	assert.Nil(t, err)
}
