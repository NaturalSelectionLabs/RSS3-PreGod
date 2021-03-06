package xscan_test

import (
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/xscan"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/stretchr/testify/assert"
)

func TestGetLatestBlockHeight(t *testing.T) {
	// eth
	blockHeight, err := xscan.GetLatestBlockHeight(constants.NetworkIDEthereum)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, blockHeight)

	// polygon
	blockHeight, err = xscan.GetLatestBlockHeight(constants.NetworkIDPolygon)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, blockHeight)
}
