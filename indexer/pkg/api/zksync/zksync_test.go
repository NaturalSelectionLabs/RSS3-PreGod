package zksync_test

import (
	"log"
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/zksync"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/stretchr/testify/assert"
)

func init() {
	if err := database.Setup(); err != nil {
		log.Fatalf("database.Setup err: %v", err)
	}

	if err := database.DB.AutoMigrate(
		&model.Cache{},
	); err != nil {
		log.Fatalf("database.Setup err: %v", err)
	}
}

func TestGetLatestBlockHeight(t *testing.T) {
	blockHeight, err := zksync.GetLatestBlockHeight()

	assert.Nil(t, err)
	assert.NotEqual(t, 0, blockHeight)
}

func TestGetTokens(t *testing.T) {
	res, err := zksync.GetTokens()

	assert.Nil(t, err)
	assert.NotEmpty(t, res)
	assert.True(t, len(res) > 0)
}

func TestGetTxsByBlock(t *testing.T) {
	res, err := zksync.GetTxsByBlock(1000, true)

	assert.Nil(t, err)
	assert.True(t, len(res) > 0)
}
