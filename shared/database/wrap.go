package database

import (
	"database/sql"
	"encoding/json"
	"log"
)

func WrapNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}

	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func UnwrapNullString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}

	return ""
}

func WrapJSON(value any) (json.RawMessage, error) {
	bytes, err := json.Marshal(value)
	return bytes, err
}

func MustWrapJSON(value any) json.RawMessage {
	bytes, err := WrapJSON(value)
	if err != nil {
		log.Panicln(err)
	}

	return bytes
}
