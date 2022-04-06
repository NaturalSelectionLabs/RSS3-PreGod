package constants_test

import (
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri/constants"
)

func TestIsValidNetworkName(t *testing.T) {
	networkID := constants.NetworkID(-1)

	t.Log(networkID.Int())
	t.Log(networkID.Symbol())
	t.Log(networkID.Name())

	t.Log(networkID.Symbol().String())
	t.Log(networkID.Symbol().ID())
	t.Log(networkID.Symbol().Name())

	t.Log(networkID.Name().String())
	t.Log(networkID.Name().ID())
	t.Log(networkID.Name().Symbol())
}
