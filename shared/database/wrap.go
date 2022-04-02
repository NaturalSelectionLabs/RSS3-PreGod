package database

import "database/sql"

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
