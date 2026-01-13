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

// toNullRawMessage serializes a pointer to JSONB, treating a nil pointer as SQL NULL.
// Using a generic pointer type avoids the "typed nil in interface" pitfall.
func toNullRawMessage[T any](v *T) pqtype.NullRawMessage {
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

func int64PtrToNullInt64(p *int64) sql.NullInt64 {
	if p == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: *p, Valid: true}
}

func nullInt64ToInt64Ptr(n sql.NullInt64) *int64 {
	if !n.Valid {
		return nil
	}
	v := n.Int64
	return &v
}
