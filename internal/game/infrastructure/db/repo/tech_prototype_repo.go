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

func (r *TechPrototypeRepo) CreatePrototype(ctx context.Context, proto *domain.TechItemPrototype) error {
	return r.q.CreateTechPrototype(ctx, mappers.TechPrototypeToCreateParams(proto))
}

func (r *TechPrototypeRepo) FindPrototypeByID(ctx context.Context, id int) (*domain.TechItemPrototype, error) {
	p, err := r.q.GetTechPrototypeByID(ctx, int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.TechPrototypeFromDB(p), nil
}

func (r *TechPrototypeRepo) FindAllPrototypes(ctx context.Context) ([]*domain.TechItemPrototype, error) {
	rows, err := r.q.ListTechPrototypes(ctx)
	if err != nil {
		return nil, err
	}
	return mappers.TechPrototypesFromDB(rows), nil
}

func (r *TechPrototypeRepo) UpdatePrototype(ctx context.Context, proto *domain.TechItemPrototype) error {
	_, err := r.q.UpdateTechPrototype(ctx, mappers.TechPrototypeToUpdateParams(proto))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ports.ErrNotFound
		}
		return err
	}
	return nil
}

func (r *TechPrototypeRepo) DeletePrototype(_ context.Context, id int) error {
	return errors.New("DeletePrototype not implemented for tech prototypes (read-only in this service)")
}
