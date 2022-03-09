package constants_test

import (
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/stretchr/testify/assert"
)

func TestIsValidPlatformSymbol(t *testing.T) {
	t.Parallel()

	assert.Equal(t, constants.IsValidPlatformSymbol("ethereum"), true)
	assert.Equal(t, constants.IsValidPlatformSymbol("unknown"), false)
	assert.Equal(t, constants.IsValidPlatformSymbol("foobar"), false)
}

func TestPlatformID_IsSignable(t *testing.T) {
	t.Parallel()

	assert.Equal(t, constants.PlatformIDEthereum.IsSignable(), true)
	assert.Equal(t, constants.PlatformIDUnknown.IsSignable(), false)
	assert.Equal(t, constants.PlatformIDJike.IsSignable(), false)
}

func TestPlatformID_Symbol(t *testing.T) {
	t.Parallel()

	assert.Equal(t, constants.PlatformIDEthereum.Symbol(), constants.PlatformSymbolEthereum)
	assert.Equal(t, constants.PlatformIDUnknown.Symbol(), constants.PlatformSymbolUnknown)
	assert.Equal(t, constants.PlatformIDJike.Symbol(), constants.PlatformSymbolJike)
}

func TestPlatformSymbol_ID(t *testing.T) {
	t.Parallel()

	assert.Equal(t, constants.PlatformSymbolEthereum.ID(), constants.PlatformIDEthereum)
	assert.Equal(t, constants.PlatformSymbolUnknown.ID(), constants.PlatformIDUnknown)
	assert.Equal(t, constants.PlatformSymbolJike.ID(), constants.PlatformIDJike)
}
