package poap_test

import (
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/poap"
	"github.com/stretchr/testify/assert"
)

func TestGetActions(t *testing.T) {
	t.Parallel()

	result, err := poap.GetActions("0xBf6f8E4ae37680a60B13C2f02b6437e6737d5203")

	assert.Nil(t, err)

	assert.True(t, len(result) > 0)
}
