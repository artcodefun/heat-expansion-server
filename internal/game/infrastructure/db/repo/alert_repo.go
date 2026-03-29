package repo

import (
	"context"
	"database/sql"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/mappers"
	"github.com/google/uuid"
)

type AlertRepo struct {
	q *gen.Queries
}

func NewAlertRepo(q *gen.Queries) *AlertRepo {
	return &AlertRepo{
		q: q,
	}
}

func (r *AlertRepo) Tx(tx ports.Transaction) ports.AlertRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &AlertRepo{
			q: r.q.WithTx(sqlTx),
		}
	}
	return r
}

func (r *AlertRepo) Create(ctx context.Context, alert *domain.Alert) error {
	params := mappers.InsertAlertParamsFromDomain(alert)
	return r.q.InsertAlert(ctx, params)
}

func (r *AlertRepo) ExistsForActivity(ctx context.Context, activityID uuid.UUID) (bool, error) {
	return r.q.ExistsForActivity(ctx, uuid.NullUUID{UUID: activityID, Valid: true})
}

func (r *AlertRepo) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return r.q.MarkAllAlertsAsReadByUser(ctx, userID)
}

func (r *AlertRepo) DeleteExpired(ctx context.Context, now int64) error {
	return r.q.DeleteExpiredAlerts(ctx, now)
}
