package ens_test

import (
	"log"
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/ens"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
)

func TestGetENSList(t *testing.T) {
	t.Parallel()

	if err := config.Setup(); err != nil {
		log.Fatalf("config.Setup err: %v", err)
	}

	result := ens.GetENSList("0x827431510a5D249cE4fdB7F00C83a3353F471848")

}
