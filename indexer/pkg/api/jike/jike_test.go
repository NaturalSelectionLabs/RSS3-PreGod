package jike_test

import (
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/jike"
	_ "github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	assert.Nil(t, logger.Setup())

	err := jike.Login()
	assert.Nil(t, err)
}

func TestGetUserProfile(t *testing.T) {
	assert.Nil(t, logger.Setup())
	assert.Nil(t, jike.Login())

	userId := "C05E4867-4251-4F11-9096-C1D720B41710"

	profile, err := jike.GetUserProfile(userId)

	assert.Equal(t, profile.ScreenName, "Henry.rss3")
	assert.Equal(t, profile.Bio, "henryqw.eth")
	assert.Nil(t, err)
}
