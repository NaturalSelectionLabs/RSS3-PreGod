package jike_test

import (
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/jike"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	err := jike.Login()

	assert.Nil(t, err)
	assert.NotEmpty(t, jike.AccessToken)

	previousRefreshToken := jike.RefreshToken

	err = jike.RefreshJikeToken()
	assert.Nil(t, err)

	assert.True(t,
		previousRefreshToken != jike.RefreshToken)
}

// func TestGetUserProfile(t *testing.T) {
// 	t.Parallel()

// 	err := jike.Login()

// 	assert.Nil(t, err)
// 	assert.NotEmpty(t, jike.AccessToken)

// 	userId := "C05E4867-4251-4F11-9096-C1D720B41710"

// 	jike.GetUserProfile(userId)
// }
