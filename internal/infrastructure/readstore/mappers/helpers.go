package mappers

import (
	"database/sql"
	"encoding/json"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/dtos"
	"github.com/sqlc-dev/pqtype"
)

func priceFromJSON(b []byte) readmodels.PriceModel {
	var pd dtos.PriceDTO
	_ = json.Unmarshal(b, &pd)
	return readmodels.PriceModel{Credits: pd.Credits, Iron: pd.Iron, Titanium: pd.Titanium, Antimatter: pd.Antimatter}
}

func technologyEffectsFromJSON(b []byte) []readmodels.TechnologyEffect {
	if len(b) == 0 {
		return nil
	}
	var arrDTO []dtos.TechnologyEffectDTO
	if err := json.Unmarshal(b, &arrDTO); err != nil {
		return nil
	}
	if len(arrDTO) == 0 {
		return nil
	}
	out := make([]readmodels.TechnologyEffect, 0, len(arrDTO))
	for _, d := range arrDTO {
		out = append(out, readmodels.TechnologyEffect{EffectType: readmodels.EffectType(d.Type), Value: d.Value})
	}
	return out
}

func nullInt64(n sql.NullInt64) int {
	if n.Valid {
		return int(n.Int64)
	}
	return 0
}

func nullString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

func jsonToNullRaw[T any](nm pqtype.NullRawMessage) *T {
	if !nm.Valid {
		return nil
	}
	var v T
	if err := json.Unmarshal(nm.RawMessage, &v); err == nil {
		return &v
	}
	return nil
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
