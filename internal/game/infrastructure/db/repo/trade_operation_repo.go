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

type TradeOperationRepo struct {
	q *gen.Queries
}

func NewTradeOperationRepo(q *gen.Queries) *TradeOperationRepo {
	return &TradeOperationRepo{q: q}
}

func (r *TradeOperationRepo) Tx(tx ports.Transaction) ports.TradeOperationRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &TradeOperationRepo{q: r.q.WithTx(sqlTx)}
	}
	return r
}

func (r *TradeOperationRepo) Create(ctx context.Context, op *domain.TradeOperation) error {
	id, err := r.q.InsertTradeOperation(ctx, mappers.InsertTradeOperationParamsFromDomain(op))
	if err != nil {
		return err
	}
	op.ID = int(id)
	return nil
}

func (r *TradeOperationRepo) FindByID(ctx context.Context, id int) (*domain.TradeOperation, error) {
	row, err := r.q.GetTradeOperationByID(ctx, int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.TradeOperationFromDB(row), nil
}

func (r *TradeOperationRepo) FindByIDForUpdate(ctx context.Context, id int) (*domain.TradeOperation, error) {
	row, err := r.q.GetTradeOperationByIDForUpdate(ctx, int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.TradeOperationFromDB(row), nil
}

func (r *TradeOperationRepo) Update(ctx context.Context, op *domain.TradeOperation) error {
	return r.q.UpdateTradeOperation(ctx, mappers.UpdateTradeOperationParamsFromDomain(op))
}

func (r *TradeOperationRepo) Delete(ctx context.Context, id int) error {
	return r.q.DeleteTradeOperation(ctx, int64(id))
}
