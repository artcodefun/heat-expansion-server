package mappers

import (
	"database/sql"
	"encoding/json"

	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/dtos"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

func nullStringToString(ns *string, valid bool) string {
	if valid && ns != nil {
		return *ns
	}
	return ""
}

// stringToNullString is the write-direction counterpart to nullStringToString:
// an empty string is stored as SQL NULL for nullable TEXT columns.
func stringToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
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

func creationSourcesFromJSON(raw []byte) []domain.CreationSource {
	if len(raw) == 0 {
		return nil
	}
	var sources []domain.CreationSource
	_ = json.Unmarshal(raw, &sources)
	return sources
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

func nullableBaseID(v *int) sql.NullInt64 {
	if v == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: int64(*v), Valid: true}
}

func nullInt64ToBaseIDPtr(v sql.NullInt64) *int {
	if !v.Valid {
		return nil
	}
	value := int(v.Int64)
	return &value
}

func nullableUUID(v *uuid.UUID) uuid.NullUUID {
	if v == nil {
		return uuid.NullUUID{Valid: false}
	}
	return uuid.NullUUID{UUID: *v, Valid: true}
}

func nullUUIDToUUIDPtr(v uuid.NullUUID) *uuid.UUID {
	if !v.Valid {
		return nil
	}
	value := v.UUID
	return &value
}

func priceToJSON(p domain.PriceModel) json.RawMessage {
	b, _ := json.Marshal(dtos.PriceDTOFromDomain(p))
	return b
}

// creationSourcesToJSON marshals the sources, emitting an empty JSON array
// rather than null so the NOT NULL jsonb column stays a valid list.
func creationSourcesToJSON(sources []domain.CreationSource) json.RawMessage {
	if len(sources) == 0 {
		return json.RawMessage("[]")
	}
	b, _ := json.Marshal(sources)
	return b
}
