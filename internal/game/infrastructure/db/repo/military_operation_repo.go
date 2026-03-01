package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/mappers"
)

type MilitaryOperationRepo struct {
	q *gen.Queries
}

func NewMilitaryOperationRepo(q *gen.Queries) *MilitaryOperationRepo {
	return &MilitaryOperationRepo{q: q}
}

func (r *MilitaryOperationRepo) Tx(tx ports.Transaction) ports.MilitaryOperationRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &MilitaryOperationRepo{q: r.q.WithTx(sqlTx)}
	}
	return r
}

func (r *MilitaryOperationRepo) Create(ctx context.Context, op *domain.MilitaryOperation) error {
	id, err := r.q.InsertMilitaryOperation(ctx, mappers.InsertMilitaryOperationParamsFromDomain(op))
	if err != nil {
		return err
	}
	op.ID = int(id)
	return nil
}

func (r *MilitaryOperationRepo) FindByID(ctx context.Context, id int) (*domain.MilitaryOperation, error) {
	row, err := r.q.GetMilitaryOperationByID(ctx, int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.MilitaryOperationFromDB(row), nil
}

func (r *MilitaryOperationRepo) FindByIDForUpdate(ctx context.Context, id int) (*domain.MilitaryOperation, error) {
	row, err := r.q.GetMilitaryOperationByIDForUpdate(ctx, int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.MilitaryOperationFromDB(row), nil
}

func (r *MilitaryOperationRepo) Update(ctx context.Context, op *domain.MilitaryOperation) error {
	return r.q.UpdateMilitaryOperation(ctx, mappers.UpdateMilitaryOperationParamsFromDomain(op))
}

func (r *MilitaryOperationRepo) Delete(ctx context.Context, id int) error {
	return r.q.DeleteMilitaryOperation(ctx, int64(id))
}
