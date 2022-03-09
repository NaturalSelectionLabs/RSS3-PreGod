package model_datatype

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Attachment struct {
	Type        string `json:"type,omitempty"`
	Content     string `json:"content,omitempty"`
	Address     string `json:"address,omitempty"`
	MimeType    string `json:"mime_type,omitempty"`
	SizeInBytes int    `json:"size_in_bytes,omitempty"`
}

type Attachments []Attachment

// Returns Attachments value, implements driver.Valuer interface
func (as Attachments) Value() (driver.Value, error) {
	if as == nil {
		return nil, nil
	}

	bytes, err := json.Marshal(as)

	return string(bytes), err
}

// Scans value into Attachments, implements sql.Scanner interface
func (as *Attachments) Scan(value interface{}) error {
	if value == nil {
		*as = nil

		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal Attachments value:", value))
	}

	result := Attachments{}
	err := json.Unmarshal(bytes, &result)
	*as = result

	return err
}
