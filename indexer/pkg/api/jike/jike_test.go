package jike_test

import (
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/jike"
	// "github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	// "github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	err := jike.Login()

	assert.Nil(t, err)
	// TODO fix empty
	// assert.NotEmpty(t, jike.AccessToken)

	// TODO fix 401
	// previousRefreshToken := jike.RefreshToken

	// assert.Nil(t, config.Setup())
	// assert.Nil(t, logger.Setup())
	// err = jike.RefreshJikeToken()
	// assert.Nil(t, err)
	// assert.True(t, previousRefreshToken != jike.RefreshToken)
	assert.Nil(t, err)
}

// func TestGetUserProfile(t *testing.T) {
// 	err := jike.Login()

// 	assert.Nil(t, err)
// 	assert.NotEmpty(t, jike.AccessToken)

// 	userId := "C05E4867-4251-4F11-9096-C1D720B41710"

// 	jike.GetUserProfile(userId)
// }
