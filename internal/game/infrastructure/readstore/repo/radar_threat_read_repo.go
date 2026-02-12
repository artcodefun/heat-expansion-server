package repo

import (
	"context"
	"database/sql"

	"github.com/artcodefun/heat-expansion-api/internal/game/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/game/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/readstore/mappers"
	"github.com/google/uuid"
)

type RadarThreatReadRepo struct {
	q *gen.Queries
}

func NewRadarThreatReadRepo(q *gen.Queries) *RadarThreatReadRepo {
	return &RadarThreatReadRepo{q: q}
}

func (r *RadarThreatReadRepo) GetRadarThreat(id uuid.UUID) (*readmodels.RadarThreat, error) {
	m, err := r.q.GetRadarThreat(context.Background(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.RadarThreatFromModel(m), nil
}

func (r *RadarThreatReadRepo) ListIncomingThreats(baseID int) ([]*readmodels.RadarThreat, error) {
	rows, err := r.q.ListIncomingThreats(context.Background(), int64(baseID))
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.RadarThreat, 0, len(rows))
	for _, m := range rows {
		out = append(out, mappers.RadarThreatFromModel(m))
	}
	return out, nil
}
