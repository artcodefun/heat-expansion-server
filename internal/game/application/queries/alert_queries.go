package queries

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/google/uuid"
)

type AlertQueries struct {
	AlertReadRepo ports.AlertReadRepository
}

func NewAlertQueries(readRepo ports.AlertReadRepository) *AlertQueries {
	return &AlertQueries{
		AlertReadRepo: readRepo,
	}
}

func (q *AlertQueries) ListActiveAlerts(ctx context.Context, actor cqrs.Actor) ([]*readmodels.AlertItem, error) {
	if actor.UserID == uuid.Nil {
		return nil, cqrs.ErrForbidden
	}
	return q.AlertReadRepo.ListActiveAlerts(ctx, actor.UserID)
}

func (q *AlertQueries) GetUnreadAlertsCount(ctx context.Context, actor cqrs.Actor) (int, error) {
	if actor.UserID == uuid.Nil {
		return 0, cqrs.ErrForbidden
	}
	return q.AlertReadRepo.GetUnreadAlertsCount(ctx, actor.UserID)
}
