package mappers

import (
	"database/sql"
	"encoding/json"

	"github.com/sqlc-dev/pqtype"
)

func nullStringToString(ns *string, valid bool) string {
	if valid && ns != nil {
		return *ns
	}
	return ""
}

func nullableIntPtr(n int64, valid bool) *int {
	if !valid {
		return nil
	}
	v := int(n)
	return &v
}

func unmarshalIfValid[T any](raw pqtype.NullRawMessage, out *T) {
	if !raw.Valid {
		return
	}
	_ = json.Unmarshal(raw.RawMessage, out)
}

func toNullRawMessage(v any) pqtype.NullRawMessage {
	if v == nil {
		return pqtype.NullRawMessage{Valid: false}
	}
	b, _ := json.Marshal(v)
	return pqtype.NullRawMessage{RawMessage: b, Valid: true}
}

// toNullInt64ZeroAsNull converts an int where 0 represents NULL to sql.NullInt64
func toNullInt64ZeroAsNull(v int) sql.NullInt64 {
	if v == 0 {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: int64(v), Valid: true}
}
