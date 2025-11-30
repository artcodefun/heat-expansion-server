package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/mappers"
)

type OperationReadRepo struct{ q *gen.Queries }

func NewOperationReadRepo(q *gen.Queries) *OperationReadRepo { return &OperationReadRepo{q: q} }

func (r *OperationReadRepo) GetOperation(opID int) (*readmodels.MilitaryOperation, error) {
	row, err := r.q.GetOperation(context.Background(), int64(opID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	m := mappers.OperationFromModel(row)
	return &m, nil
}

func (r *OperationReadRepo) ListOperationsByBase(baseID int) ([]*readmodels.MilitaryOperation, error) {
	rows, err := r.q.ListOperationsByBase(context.Background(), int64(baseID))
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.MilitaryOperation, 0, len(rows))
	for _, row := range rows {
		v := mappers.OperationFromModel(row)
		out = append(out, &v)
	}
	return out, nil
}

func (r *OperationReadRepo) ListActiveOperations(baseID int) ([]*readmodels.MilitaryOperation, error) {
	rows, err := r.q.ListActiveOperations(context.Background(), int64(baseID))
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.MilitaryOperation, 0, len(rows))
	for _, row := range rows {
		v := mappers.OperationFromModel(row)
		out = append(out, &v)
	}
	return out, nil
}

func (r *OperationReadRepo) ListRadarDetectedOperations(baseID int) ([]readmodels.RadarActivity, error) {
	// Derived from activities; fetch and filter radar activities
	acts, err := r.q.ListMilitaryActivities(context.Background(), gen.ListMilitaryActivitiesParams{BaseID: int64(baseID), Limit: 100})
	if err != nil {
		return nil, err
	}
	out := []readmodels.RadarActivity{}
	for _, a := range acts {
		item := mappers.ActivityItemFromModel(a)
		if item.Radar != nil {
			out = append(out, *item.Radar)
		}
	}
	return out, nil
}
