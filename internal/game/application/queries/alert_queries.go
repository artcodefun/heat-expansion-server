package queries

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
)

type AlertQueries struct {
	AlertReadRepo ports.AlertReadRepository
	Access        *services.AccessControlService
}

func NewAlertQueries(readRepo ports.AlertReadRepository, access *services.AccessControlService) *AlertQueries {
	return &AlertQueries{
		AlertReadRepo: readRepo,
		Access:        access,
	}
}

func (q *AlertQueries) ListActiveAlerts(ctx context.Context, actor cqrs.Actor, baseID int) ([]*readmodels.AlertItem, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	return q.AlertReadRepo.ListActiveAlerts(ctx, baseID)
}

func (q *AlertQueries) GetUnreadAlertsCount(ctx context.Context, actor cqrs.Actor, baseID int) (int, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return 0, err
	}
	return q.AlertReadRepo.GetUnreadAlertsCount(ctx, baseID)
}
