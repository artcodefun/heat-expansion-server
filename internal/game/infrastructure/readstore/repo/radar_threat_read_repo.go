package repo

import (
	"context"
	"database/sql"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/mappers"
	"github.com/google/uuid"
)

type RadarThreatReadRepo struct {
	q *gen.Queries
}

func NewRadarThreatReadRepo(q *gen.Queries) *RadarThreatReadRepo {
	return &RadarThreatReadRepo{q: q}
}

func (r *RadarThreatReadRepo) GetRadarThreat(ctx context.Context, id uuid.UUID) (*readmodels.RadarThreat, error) {
	m, err := r.q.GetRadarThreat(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.RadarThreatFromModel(m), nil
}

func (r *RadarThreatReadRepo) ListIncomingThreats(ctx context.Context, baseID int) ([]*readmodels.RadarThreat, error) {
	rows, err := r.q.ListIncomingThreats(ctx, int64(baseID))
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.RadarThreat, 0, len(rows))
	for _, m := range rows {
		out = append(out, mappers.RadarThreatFromModel(m))
	}
	return out, nil
}
