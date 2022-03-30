package ens_test

import (
	"log"
	"testing"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/ens"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/stretchr/testify/assert"
)

func init() {
	if err := config.Setup(); err != nil {
		log.Fatalf("config.Setup err: %v", err)
	}
}

func TestGetENSList(t *testing.T) {
	t.Parallel()

	result, getErr := ens.GetENSList("0x827431510a5D249cE4fdB7F00C83a3353F471848")

	ens := result[0]

	assert.Nil(t, getErr)
	assert.Equal(t, len(result), 1)

	assert.Equal(t, ens.Domain, "henryqw.eth")
	assert.Equal(t, ens.Description, "henryqw.eth, an ENS name.")
	assert.Equal(t, ens.TxHash, "0x44ea5a47fa51ada626874ac5c243e78ee485e354d5b337ea673d7f117eb8b6c3")

	time, timeErr := time.Parse(time.RFC3339, "2022-01-02T11:16:35.000Z")
	assert.Nil(t, timeErr)

	assert.Equal(t, ens.CreatedAt, time)
}
