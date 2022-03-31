package constants_test

import (
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/stretchr/testify/assert"
)

func TestIsValidNetworkName(t *testing.T) {
	assert.Equal(t, constants.IsValidNetworkName("ethereum"), true)
	assert.Equal(t, constants.IsValidNetworkName("unknown"), false)
	assert.Equal(t, constants.IsValidNetworkName("foobar"), false)
}
