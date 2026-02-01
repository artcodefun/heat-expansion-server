package repo

import (
	"context"
	"database/sql"

	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/mappers"
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

func (r *AlertRepo) Create(alert *domain.Alert) error {
	params := mappers.InsertAlertParamsFromDomain(alert)
	return r.q.InsertAlert(context.Background(), params)
}

func (r *AlertRepo) ExistsForActivity(activityID uuid.UUID) (bool, error) {
	return r.q.ExistsForActivity(context.Background(), uuid.NullUUID{UUID: activityID, Valid: true})
}

func (r *AlertRepo) MarkAllAsRead(baseID int) error {
	return r.q.MarkAllAlertsAsRead(context.Background(), int64(baseID))
}

func (r *AlertRepo) DeleteExpired(now int64) error {
	return r.q.DeleteExpiredAlerts(context.Background(), now)
}
