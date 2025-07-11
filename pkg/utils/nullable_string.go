package utils

import (
	"database/sql"
	"encoding/json"
)

type NullableString sql.NullString

func (ns NullableString) MarshalJSON() ([]byte, error) {
	if !sql.NullString(ns).Valid {
		return []byte("null"), nil
	}
	return json.Marshal(sql.NullString(ns).String)
}
