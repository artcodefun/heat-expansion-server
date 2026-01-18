package repo

import (
	"context"
	"database/sql"

	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/mappers"
)

type AlertRepository struct {
	q *gen.Queries
}

func NewAlertRepository(q *gen.Queries) *AlertRepository {
	return &AlertRepository{
		q: q,
	}
}

func (r *AlertRepository) Tx(tx ports.Transaction) ports.AlertRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &AlertRepository{
			q: r.q.WithTx(sqlTx),
		}
	}
	return r
}

func (r *AlertRepository) Create(alert *domain.Alert) error {
	params := mappers.InsertAlertParamsFromDomain(alert)
	return r.q.InsertAlert(context.Background(), params)
}

func (r *AlertRepository) MarkAllAsRead(baseID int) error {
	return r.q.MarkAllAlertsAsRead(context.Background(), int64(baseID))
}

func (r *AlertRepository) DeleteExpired(now int64) error {
	return r.q.DeleteExpiredAlerts(context.Background(), now)
}
