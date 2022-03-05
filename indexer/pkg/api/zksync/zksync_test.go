package zksync_test

import (
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/zksync"
	"github.com/stretchr/testify/assert"
)

func TestGetLatestBlockHeight(t *testing.T) {
	t.Parallel()

	blockHeight, err := zksync.GetLatestBlockHeight()

	assert.Nil(t, err)
	assert.NotEqual(t, 0, blockHeight)
}

func TestGetTokens(t *testing.T) {
	t.Parallel()

	res, err := zksync.GetTokens()

	assert.Nil(t, err)
	assert.NotEmpty(t, res)
	assert.True(t, len(res) > 0)
}

func TestGetTxsByBlock(t *testing.T) {
	t.Parallel()

	res, err := zksync.GetTxsByBlock(1000)

	assert.Nil(t, err)
	assert.True(t, len(res) > 0)
}
