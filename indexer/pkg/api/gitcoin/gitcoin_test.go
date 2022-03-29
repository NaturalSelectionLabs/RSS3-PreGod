package gitcoin_test

import (
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/gitcoin"
	"github.com/stretchr/testify/assert"
)

func TestGetGrantsInfo(t *testing.T) {
	grants, err := gitcoin.GetGrantsInfo()
	assert.Nil(t, err)
	// TODO fix empty response
	// assert.NotEmpty(t, grants)

	for _, item := range grants {
		if item.AdminAddress != "\"0x0\"" {
			// check title
			assert.NotEmpty(t, item.Title)
			// check address
			assert.NotEmpty(t, item.AdminAddress)
		}
	}
}

func TestGetProjectsInfo(t *testing.T) {
	res, err := gitcoin.GetProjectsInfo("0x8c23B96f2fb77AaE1ac2832debEE30f09da7af3C", "RSS3")
	assert.Nil(t, err)
	assert.NotEmpty(t, res)
}

func TestGetEthDonations(t *testing.T) {
	// TODO fix http error and empty
	// res, err := gitcoin.GetEthDonations(12605342, 12605343, gitcoin.ETH)
	// assert.Nil(t, err)
	// assert.NotEmpty(t, res)
}

func TestGetZkSyncDonations(t *testing.T) {
	ethParam := gitcoin.Param{FromHeight: 1, Step: 10000, MinStep: 10, Confirmations: 10, SleepInterval: 600}
	gc := gitcoin.NewGitcoinCrawler(ethParam, ethParam, ethParam)
	err := gc.UpdateZksToken()
	assert.Nil(t, err)

	_, err = gc.GetZkSyncDonations(1000, 1001)
	// TODO fix empty
	// assert.NotEmpty(t, res)
	assert.Nil(t, err)
}
