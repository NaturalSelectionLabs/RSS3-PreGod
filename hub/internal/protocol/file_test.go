package protocol_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
)

func TestFile_MarshalJSON(t *testing.T) {
	f := &protocol.File{
		Version:        protocol.Version,
		DateUpdated:    time.Now(),
		Identifier:     "",
		IdentifierNext: "",
		Total:          0,
		List:           nil,
	}

	data, err := json.Marshal(f)
	if err != nil {
		t.Error(err)
	}

	t.Log(string(data))
}
