package misskey_test

import (
	"testing"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/misskey"
	"github.com/stretchr/testify/assert"
)

// TODO: add more tests maybe
func TestGetNoteList(t *testing.T) {
	t.Parallel()

	noteList, err := misskey.GetUserNoteList("cororonnxx@misskey.dev", 100, time.Now().Add(-time.Hour*24*10))

	assert.Nil(t, err)
	assert.NotEmpty(t, noteList)
}
