package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-api/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-api/internal/game/domain"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/db/mappers"
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

func (r *MilitaryOperationRepo) Create(op *domain.MilitaryOperation) error {
	id, err := r.q.InsertMilitaryOperation(context.Background(), mappers.InsertMilitaryOperationParamsFromDomain(op))
	if err != nil {
		return err
	}
	op.ID = int(id)
	return nil
}

func (r *MilitaryOperationRepo) FindByID(id int) (*domain.MilitaryOperation, error) {
	row, err := r.q.GetMilitaryOperationByID(context.Background(), int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.MilitaryOperationFromDB(row), nil
}

func (r *MilitaryOperationRepo) FindByIDForUpdate(id int) (*domain.MilitaryOperation, error) {
	row, err := r.q.GetMilitaryOperationByIDForUpdate(context.Background(), int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.MilitaryOperationFromDB(row), nil
}

func (r *MilitaryOperationRepo) Update(op *domain.MilitaryOperation) error {
	return r.q.UpdateMilitaryOperation(context.Background(), mappers.UpdateMilitaryOperationParamsFromDomain(op))
}

func (r *MilitaryOperationRepo) Delete(id int) error {
	return r.q.DeleteMilitaryOperation(context.Background(), int64(id))
}
