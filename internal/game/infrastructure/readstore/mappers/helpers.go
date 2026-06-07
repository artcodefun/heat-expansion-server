package mappers

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/dtos"
	"github.com/google/uuid"
)

func priceFromJSON(b []byte) readmodels.PriceModel {
	var pd dtos.PriceDTO
	_ = json.Unmarshal(b, &pd)
	return readmodels.PriceModel{Credits: pd.Credits, Iron: pd.Iron, Titanium: pd.Titanium, Antimatter: pd.Antimatter}
}

func nullString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

// IdsToInt64 converts a slice of int to a slice of int64 for sqlc-generated queries.
func IdsToInt64(ids []int) []int64 {
	if len(ids) == 0 {
		return nil
	}
	out := make([]int64, 0, len(ids))
	for _, id := range ids {
		out = append(out, int64(id))
	}
	return out
}

func nullInt64ToInt64Ptr(n sql.NullInt64) *int64 {
	if !n.Valid {
		return nil
	}
	v := n.Int64
	return &v
}

func nullBaseIDPtr(v sql.NullInt64) *int {
	if !v.Valid {
		return nil
	}
	value := int(v.Int64)
	return &value
}

func nullUUIDPtr(v uuid.NullUUID) *uuid.UUID {
	if !v.Valid {
		return nil
	}
	value := v.UUID
	return &value
}

func interfaceUUID(v interface{}) uuid.UUID {
	switch value := v.(type) {
	case uuid.UUID:
		return value
	case []byte:
		parsed, err := uuid.ParseBytes(value)
		if err != nil {
			panic(fmt.Sprintf("invalid uuid bytes from readstore: %v", err))
		}
		return parsed
	case string:
		parsed, err := uuid.Parse(value)
		if err != nil {
			panic(fmt.Sprintf("invalid uuid string from readstore: %v", err))
		}
		return parsed
	default:
		panic(fmt.Sprintf("unsupported uuid source type %T", v))
	}
}

func creationSourcesFromJSON(raw json.RawMessage) []readmodels.CreationSource {
	if len(raw) == 0 {
		return nil
	}
	var ss []string
	if err := json.Unmarshal(raw, &ss); err != nil {
		return nil
	}
	out := make([]readmodels.CreationSource, len(ss))
	for i, s := range ss {
		out[i] = readmodels.CreationSource(s)
	}
	return out
}
