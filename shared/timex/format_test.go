package timex_test

import (
	"testing"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/timex"
)

func TestName(t *testing.T) {
	result, err := time.Parse(timex.ISO8601, "2022-03-22T11:52:22.865Z")
	if err != nil {
		t.Error(err)
	}

	t.Log(result)
}
