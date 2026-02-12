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

type BuildPrototypeRepo struct {
	q *gen.Queries
}

func NewBuildPrototypeRepo(q *gen.Queries) *BuildPrototypeRepo {
	return &BuildPrototypeRepo{q: q}
}

func (r *BuildPrototypeRepo) Tx(tx ports.Transaction) ports.BuildPrototypeRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &BuildPrototypeRepo{q: r.q.WithTx(sqlTx)}
	}
	return r
}

func (r *BuildPrototypeRepo) CreatePrototype(proto *domain.BuildItemPrototype) error {
	return errors.New("CreatePrototype not implemented for build prototypes (read-only in this service)")
}

func (r *BuildPrototypeRepo) FindPrototypeByID(id int) (*domain.BuildItemPrototype, error) {
	p, err := r.q.GetBuildPrototypeByID(context.Background(), int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.BuildPrototypeFromDB(p), nil
}

func (r *BuildPrototypeRepo) FindAllPrototypes() ([]*domain.BuildItemPrototype, error) {
	rows, err := r.q.ListBuildPrototypes(context.Background())
	if err != nil {
		return nil, err
	}
	return mappers.BuildPrototypesFromDB(rows), nil
}

func (r *BuildPrototypeRepo) UpdatePrototype(proto *domain.BuildItemPrototype) error {
	return errors.New("UpdatePrototype not implemented for build prototypes (read-only in this service)")
}

func (r *BuildPrototypeRepo) DeletePrototype(id int) error {
	return errors.New("DeletePrototype not implemented for build prototypes (read-only in this service)")
}
