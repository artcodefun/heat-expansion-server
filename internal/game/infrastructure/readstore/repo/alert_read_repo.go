package repo

import (
	"context"
	"time"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/mappers"
	"github.com/google/uuid"
)

type AlertReadRepository struct {
	q *gen.Queries
}

func NewAlertReadRepository(q *gen.Queries) *AlertReadRepository {
	return &AlertReadRepository{
		q: q,
	}
}

func (r *AlertReadRepository) ListActiveAlerts(ctx context.Context, userID uuid.UUID) ([]*readmodels.AlertItem, error) {
	now := time.Now().Unix()
	rows, err := r.q.ListAlertsByUser(ctx, gen.ListAlertsByUserParams{
		UserID:    userID,
		ExpiresAt: now,
	})
	if err != nil {
		return nil, err
	}

	return mappers.AlertItemsFromModels(rows), nil
}

func (r *AlertReadRepository) GetUnreadAlertsCount(ctx context.Context, userID uuid.UUID) (int, error) {
	now := time.Now().Unix()
	count, err := r.q.CountUnreadAlertsByUser(ctx, gen.CountUnreadAlertsByUserParams{
		UserID:    userID,
		ExpiresAt: now,
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
