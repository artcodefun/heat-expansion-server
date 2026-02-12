package queries

import (
	"github.com/artcodefun/heat-expansion-api/internal/game/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/game/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/game/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/game/core/services"
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

func (q *AlertQueries) ListActiveAlerts(ctx cqrs.QueryContext, baseID int) ([]*readmodels.AlertItem, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	return q.AlertReadRepo.ListActiveAlerts(baseID)
}

func (q *AlertQueries) GetUnreadAlertsCount(ctx cqrs.QueryContext, baseID int) (int, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return 0, err
	}
	return q.AlertReadRepo.GetUnreadAlertsCount(baseID)
}
