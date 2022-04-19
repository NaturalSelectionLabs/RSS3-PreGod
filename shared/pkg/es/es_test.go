package es_test

import (
	"context"
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/es"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	err := es.Ping(context.Background())
	assert.Nil(t, err)
}
