package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/billing/domain"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/db/mappers"
	"github.com/google/uuid"
)

type OrderRepo struct {
	q *gen.Queries
}

func NewOrderRepo(q *gen.Queries) *OrderRepo {
	return &OrderRepo{q: q}
}

func (r *OrderRepo) Tx(tx ports.Transaction) ports.PurchaseOrderRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &OrderRepo{q: r.q.WithTx(sqlTx)}
	}
	return r
}

func (r *OrderRepo) Save(ctx context.Context, order *domain.PurchaseOrder) error {
	return r.q.InsertOrder(ctx, mappers.InsertOrderParamsFromDomain(order))
}

func (r *OrderRepo) Update(ctx context.Context, order *domain.PurchaseOrder) error {
	return r.q.UpdateOrder(ctx, mappers.UpdateOrderParamsFromDomain(order))
}

func (r *OrderRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.PurchaseOrder, error) {
	row, err := r.q.GetOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.OrderFromRow(row), nil
}

func (r *OrderRepo) FindByProviderOrderID(ctx context.Context, providerOrderID string) (*domain.PurchaseOrder, error) {
	row, err := r.q.GetOrderByProviderOrderID(ctx, providerOrderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.OrderFromRow(row), nil
}

// FindByProviderOrderIDForUpdate locks the order row FOR UPDATE; it must be
// called on a transaction-bound repository via Tx(tx).
func (r *OrderRepo) FindByProviderOrderIDForUpdate(ctx context.Context, providerOrderID string) (*domain.PurchaseOrder, error) {
	row, err := r.q.GetOrderByProviderOrderIDForUpdate(ctx, providerOrderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.OrderFromRow(row), nil
}
