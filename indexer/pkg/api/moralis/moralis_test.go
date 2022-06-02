package moralis_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/moralis"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	_ "github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/stretchr/testify/assert"
)

var (
	tokenId      = "69122868356010038918278537874891361194569907163152093427587761621557332847656"
	tokenAddress = "0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85"
	ensContract  = "0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85"
)

func init() {
	if err := database.Setup(); err != nil {
		log.Fatalf("database.Setup err: %v", err)
	}
}

func TestGetNFT(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	result, err := moralis.GetNFTs(ctx,
		"0x3b6d02a24df681ffdf621d35d70aba7adaac07c1", "eth", time.Unix(0, 0).String(), config.Config.Indexer.Moralis.ApiKey,
	)
	assert.NotEmpty(t, result.Result)
	// assert for nil
	assert.Nil(t, err)
}

func TestGetNFTTransfers(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	result, err := moralis.GetNFTTransfers(ctx,
		"0x3b6d02a24df681ffdf621d35d70aba7adaac07c1", "eth", 0, time.Unix(0, 0).String(), config.Config.Indexer.Moralis.ApiKey,
	)

	assert.NotEmpty(t, result.Result)
	// assert for nil
	assert.Nil(t, err)
}

func TestGetLogs(t *testing.T) {
	t.Parallel()

	result, err := moralis.GetLogs(
		context.Background(),
		12605342,
		12605343,
		"0x7d655c57f71464B6f83811C55D84009Cd9f5221C",
		"0x3bb7428b25f9bdad9bd2faa4c6a7a9e5d5882657e96c1d24cc41c1d6c1910a98",
		"eth",
		config.Config.Indexer.Moralis.ApiKey,
		"moralis-gitcoin")
	// assert for nil
	assert.Nil(t, err)

	if assert.NotEmpty(t, result) && assert.NotEmpty(t, result.Result) {
		for _, item := range result.Result {
			assert.NotEmpty(t, item)
		}
	}
}

func TestGetTxByToken(t *testing.T) {
	t.Parallel()

	result, err := moralis.GetTxByToken(
		context.Background(),
		tokenAddress, tokenId,
		"eth",
		config.Config.Indexer.Moralis.ApiKey)

	assert.Equal(t, "0x44ea5a47fa51ada626874ac5c243e78ee485e354d5b337ea673d7f117eb8b6c3", result.TransactionHash)
	assert.Equal(t, "2022-01-02T11:16:35.000Z", result.BlockTimestamp)

	assert.Nil(t, err)
}

func TestGetNFTByContract(t *testing.T) {
	t.Parallel()

	result, err := moralis.GetNFTByContract(
		context.Background(),
		"0x827431510a5D249cE4fdB7F00C83a3353F471848", ensContract,
		"eth",
		config.Config.Indexer.Moralis.ApiKey)

	assert.Equal(t, 1, len(result.Result))

	if len(result.Result) > 0 {
		ens := result.Result[0]

		assert.Equal(t, tokenAddress, ens.TokenAddress)
		assert.Equal(t, tokenId, ens.TokenId)

		assert.Nil(t, err)
	}
}

// func TestGetENSList(t *testing.T) {
// 	t.Parallel()

// 	result, getErr := moralis.GetENSList("0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944")

// 	if assert.NotEmpty(t, result) {
// 		ens := result[0]

// 		assert.Nil(t, getErr)
// 		assert.Equal(t, 1, len(result))

// 		assert.Equal(t, "diygod.eth", ens.Domain)
// 		assert.Equal(t, "diygod.eth, an ENS name.", ens.Description)
// 		assert.Equal(t, "0xc600982712df36668321bfc782deacb17a1c32f09165eb1e66d1d76294db6156", ens.TxHash)

// 		time, timeErr := time.Parse(time.RFC3339, "2021-11-16T05:54:43.000Z")
// 		assert.Nil(t, timeErr)

// 		assert.Equal(t, time, ens.CreatedAt)
// 	}
// }
