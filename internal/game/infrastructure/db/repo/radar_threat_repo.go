package repo

import (
	"context"
	"database/sql"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/mappers"
	"github.com/google/uuid"
)

type RadarThreatRepo struct {
	q *gen.Queries
}

func NewRadarThreatRepo(q *gen.Queries) *RadarThreatRepo {
	return &RadarThreatRepo{q: q}
}

func (r *RadarThreatRepo) Tx(tx ports.Transaction) ports.RadarThreatRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &RadarThreatRepo{q: r.q.WithTx(sqlTx)}
	}
	return r
}

func (r *RadarThreatRepo) Create(threat *domain.RadarThreat) error {
	params := mappers.InsertRadarThreatParamsFromDomain(threat)
	_, err := r.q.InsertRadarThreat(context.Background(), params)
	return err
}

func (r *RadarThreatRepo) Update(threat *domain.RadarThreat) error {
	params := mappers.UpdateRadarThreatParamsFromDomain(threat)
	_, err := r.q.UpdateRadarThreat(context.Background(), params)
	return err
}

func (r *RadarThreatRepo) FindByID(id uuid.UUID) (*domain.RadarThreat, error) {
	m, err := r.q.GetRadarThreat(context.Background(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.RadarThreatFromModel(m), nil
}

func (r *RadarThreatRepo) FindByOperationID(opID int) (*domain.RadarThreat, error) {
	m, err := r.q.GetRadarThreatByOperationID(context.Background(), int64(opID))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.RadarThreatFromModel(m), nil
}

func (r *RadarThreatRepo) RadarThreatExists(ownerBaseID int, opID int) (bool, error) {
	return r.q.RadarThreatExists(context.Background(), gen.RadarThreatExistsParams{
		OwnerBaseID: int64(ownerBaseID),
		OperationID: int64(opID),
	})
}
