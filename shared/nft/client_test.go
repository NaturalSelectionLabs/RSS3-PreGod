package nft_test

import (
	"math/big"
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/nft"
	"github.com/ethereum/go-ethereum/common"
)

func TestGetMetadata(t *testing.T) {
	if _, err := nft.GetMetadata(
		nft.NetworkEthereum,
		common.HexToAddress("0xacbe98efe2d4d103e221e04c76d7c55db15c8e89"),
		big.NewInt(1),
	); err != nil {
		t.Fatal(err)
	}
}
