package protocol

import (
	"encoding/json"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/isotime"
)

const (
	Version = "v0.4.0"
)

var _ json.Marshaler = &File{}

type File struct {
	Version        string    `json:"version"`
	DateUpdated    time.Time `json:"date_updated"`
	Identifier     string    `json:"identifier"`
	IdentifierNext string    `json:"identifier_next,omitempty"`
	Total          int       `json:"total"`
	List           any       `json:"list,omitempty"`
}

func (f File) MarshalJSON() ([]byte, error) {
	return json.Marshal(&magicFile{
		Version:        Version,
		DateUpdated:    f.DateUpdated.Format(isotime.ISO8601),
		Identifier:     f.Identifier,
		IdentifierNext: f.IdentifierNext,
		Total:          f.Total,
		List:           f.List,
	})
}

type magicFile struct {
	Version        string `json:"version"`
	DateUpdated    string `json:"date_updated"`
	Identifier     string `json:"identifier"`
	IdentifierNext string `json:"identifier_next,omitempty"`
	Total          int    `json:"total"`
	List           any    `json:"list,omitempty"`
}
