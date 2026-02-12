package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-api/internal/game/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/game/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/db/mappers"
)

type TechPrototypeRepo struct {
	q *gen.Queries
}

func NewTechPrototypeRepo(q *gen.Queries) *TechPrototypeRepo {
	return &TechPrototypeRepo{q: q}
}

func (r *TechPrototypeRepo) Tx(tx ports.Transaction) ports.TechPrototypeRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &TechPrototypeRepo{q: r.q.WithTx(sqlTx)}
	}
	return r
}

func (r *TechPrototypeRepo) CreatePrototype(proto *domain.TechItemPrototype) error {
	return errors.New("CreatePrototype not implemented for tech prototypes (read-only in this service)")
}

func (r *TechPrototypeRepo) FindPrototypeByID(id int) (*domain.TechItemPrototype, error) {
	p, err := r.q.GetTechPrototypeByID(context.Background(), int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.TechPrototypeFromDB(p), nil
}

func (r *TechPrototypeRepo) FindAllPrototypes() ([]*domain.TechItemPrototype, error) {
	rows, err := r.q.ListTechPrototypes(context.Background())
	if err != nil {
		return nil, err
	}
	return mappers.TechPrototypesFromDB(rows), nil
}

func (r *TechPrototypeRepo) UpdatePrototype(proto *domain.TechItemPrototype) error {
	return errors.New("UpdatePrototype not implemented for tech prototypes (read-only in this service)")
}

func (r *TechPrototypeRepo) DeletePrototype(id int) error {
	return errors.New("DeletePrototype not implemented for tech prototypes (read-only in this service)")
}
