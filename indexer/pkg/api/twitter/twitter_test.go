package twitter_test

import (
	"log"
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/twitter"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func init() {
	if err := config.Setup(); err != nil {
		log.Fatalf("config.Setup err: %v", err)
	}

	if err := logger.Setup(); err != nil {
		log.Fatalf("config.Setup err: %v", err)
	}
}

func TestGetUsersShow(t *testing.T) {
	t.Parallel()

	result, err := twitter.GetUsersShow("@rss3_")

	assert.Nil(t, err)

	assert.NotEmpty(t, result.Name)
	assert.NotEmpty(t, result.ScreenName)
	assert.NotEmpty(t, result.Description)
}

func TestGetTimeline(t *testing.T) {
	t.Parallel()

	result, err := twitter.GetTimeline("@rss3_", 200)

	assert.Nil(t, err)

	assert.True(t, len(result) > 0)
}
