package repo

import (
	"context"
	"database/sql"

	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/mappers"
)

type ActivityRepo struct {
	q *gen.Queries
}

func NewActivityRepo(q *gen.Queries) *ActivityRepo { return &ActivityRepo{q: q} }

func (r *ActivityRepo) Tx(tx ports.Transaction) ports.ActivityRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &ActivityRepo{q: r.q.WithTx(sqlTx)}
	}
	return r
}

func (r *ActivityRepo) Create(item *domain.ActivityItem) error {
	params := mappers.InsertActivityParamsFromDomain(item)
	_, err := r.q.InsertActivity(context.Background(), params)
	return err
}
