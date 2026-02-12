package mappers

import (
	"database/sql"
	"encoding/json"

	"github.com/artcodefun/heat-expansion-api/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/db/dtos"
)

func priceFromJSON(b []byte) readmodels.PriceModel {
	var pd dtos.PriceDTO
	_ = json.Unmarshal(b, &pd)
	return readmodels.PriceModel{Credits: pd.Credits, Iron: pd.Iron, Titanium: pd.Titanium, Antimatter: pd.Antimatter}
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
