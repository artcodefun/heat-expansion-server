package queries

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/readmodels"
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

func (q *AlertQueries) ListActiveAlerts(ctx context.Context, actor application.Actor) ([]*readmodels.AlertItem, error) {
	if actor.UserID == uuid.Nil {
		return nil, application.ErrForbidden
	}
	return q.AlertReadRepo.ListActiveAlerts(ctx, actor.UserID)
}

func (q *AlertQueries) GetUnreadAlertsCount(ctx context.Context, actor application.Actor) (int, error) {
	if actor.UserID == uuid.Nil {
		return 0, application.ErrForbidden
	}
	return q.AlertReadRepo.GetUnreadAlertsCount(ctx, actor.UserID)
}
