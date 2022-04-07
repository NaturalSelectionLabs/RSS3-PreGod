package database

import (
	"database/sql"
	"encoding/json"
	"log"

	"gorm.io/datatypes"
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

func WrapJSON(value any) (datatypes.JSON, error) {
	bytes, err := json.Marshal(value)

	return bytes, err
}

func UnwrapJSON[T any](value datatypes.JSON) (T, error) {
	var a T

	err := json.Unmarshal(value, &a)

	return a, err
}

func MustWrapJSON(value any) datatypes.JSON {
	bytes, err := WrapJSON(value)
	if err != nil {
		log.Panicln(err)
	}

	return bytes
}
