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

type ArmyPrototypeRepo struct {
	q *gen.Queries
}

func NewArmyPrototypeRepo(q *gen.Queries) *ArmyPrototypeRepo {
	return &ArmyPrototypeRepo{q: q}
}

func (r *ArmyPrototypeRepo) Tx(tx ports.Transaction) ports.ArmyPrototypeRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &ArmyPrototypeRepo{q: r.q.WithTx(sqlTx)}
	}
	// Fallback to original if type mismatch; better than panicking in production
	return r
}

func (r *ArmyPrototypeRepo) CreatePrototype(proto *domain.ArmyItemPrototype) error {
	// Not implemented yet: write path is rarely used and requires dedicated sqlc queries
	return errors.New("CreatePrototype not implemented for army prototypes (read-only in this service)")
}

func (r *ArmyPrototypeRepo) FindPrototypeByID(id int) (*domain.ArmyItemPrototype, error) {
	p, err := r.q.GetArmyPrototypeByID(context.Background(), int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.ArmyPrototypeFromDB(p), nil
}

func (r *ArmyPrototypeRepo) FindAllPrototypes() ([]*domain.ArmyItemPrototype, error) {
	rows, err := r.q.ListArmyPrototypes(context.Background())
	if err != nil {
		return nil, err
	}
	return mappers.ArmyPrototypesFromDB(rows), nil
}

func (r *ArmyPrototypeRepo) UpdatePrototype(proto *domain.ArmyItemPrototype) error {
	return errors.New("UpdatePrototype not implemented for army prototypes (read-only in this service)")
}

func (r *ArmyPrototypeRepo) DeletePrototype(id int) error {
	return errors.New("DeletePrototype not implemented for army prototypes (read-only in this service)")
}
