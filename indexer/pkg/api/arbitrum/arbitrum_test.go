package arbitrum_test

import (
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/arbitrum"
	"github.com/stretchr/testify/assert"
)

func TestGetNFTTransfers(t *testing.T) {
	t.Parallel()

	res, err := arbitrum.GetNFTTransfers("0xc661572db4d55e5cd96c9813f19f92f694f79814")
	assert.Nil(t, err)
	assert.NotEqual(t, 0, res)
}

func TestGetNFTs(t *testing.T) {
	t.Parallel()

	res, err := arbitrum.GetNFTs("0xc661572db4d55e5cd96c9813f19f92f694f79814")
	assert.Nil(t, err)
	assert.NotEqual(t, 0, res)
	assert.True(t, len(res) > 0)
}
