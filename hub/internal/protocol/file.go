package protocol

import (
	"encoding/json"
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/timex"
)

const (
	Version = "v0.4.0"
)

var _ json.Marshaler = version{}

type version []byte

func (v version) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, Version)), nil
}

type File struct {
	Version        version     `json:"version"`
	DateUpdated    *timex.Time `json:"date_updated,omitempty"`
	Identifier     string      `json:"identifier"`
	IdentifierNext string      `json:"identifier_next,omitempty"`
	Total          int64       `json:"total"`
	List           any         `json:"list"`
}
