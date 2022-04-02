package database

import (
	"database/sql"
	"encoding/json"
	"gorm.io/datatypes"
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

func WrapJSON(value any) (datatypes.JSON, error) {
	bytes, err := json.Marshal(value)

	return bytes, err
}

func MustWrapJSON(value any) datatypes.JSON {
	bytes, err := WrapJSON(value)
	if err != nil {
		log.Panicln(err)
	}

	return bytes
}
