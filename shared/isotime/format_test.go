package isotime_test

import (
	"testing"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/isotime"
)

func TestName(t *testing.T) {
	result, err := time.Parse(isotime.ISO8601, "2022-03-22T11:52:22.865Z")
	if err != nil {
		t.Error(err)
	}

	t.Log(result)
}
