package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	dbgen "github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/readstore/mappers"
	"github.com/google/uuid"
)

type OrderReadRepo struct {
	q *dbgen.Queries
}

func NewOrderReadRepo(q *dbgen.Queries) *OrderReadRepo {
	return &OrderReadRepo{q: q}
}

func (r *OrderReadRepo) FindByID(ctx context.Context, id uuid.UUID) (*readmodels.PurchaseOrder, error) {
	row, err := r.q.GetOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	rm := mappers.OrderReadModelFromRow(row)
	return &rm, nil
}
