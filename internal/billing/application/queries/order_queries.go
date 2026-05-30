package queries

import (
	"context"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	"github.com/google/uuid"
)

type OrderQueries struct {
	Orders ports.OrderReadRepository
}

func NewOrderQueries(orders ports.OrderReadRepository) *OrderQueries {
	return &OrderQueries{Orders: orders}
}

func (q *OrderQueries) GetOrder(ctx context.Context, actor cqrs.Actor, orderID uuid.UUID) (*readmodels.PurchaseOrder, error) {
	order, err := q.Orders.FindByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			return nil, cqrs.ErrOrderNotFound
		}
		return nil, err
	}
	if order.UserID != actor.UserID {
		return nil, cqrs.ErrForbidden
	}
	return order, nil
}
