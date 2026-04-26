package repo

import (
	"context"
	"database/sql"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/mappers"
)

type ActivityRepo struct {
	q *gen.Queries
}

func NewActivityRepo(q *gen.Queries) *ActivityRepo { return &ActivityRepo{q: q} }

func (r *ActivityRepo) Tx(tx ports.Transaction) ports.ActivityRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &ActivityRepo{q: r.q.WithTx(sqlTx)}
	}
	return r
}

func (r *ActivityRepo) Create(ctx context.Context, item *domain.ActivityItem) error {
	params := mappers.InsertActivityParamsFromDomain(item)
	_, err := r.q.InsertActivity(ctx, params)
	return err
}

func (r *ActivityRepo) ExistsForOperation(ctx context.Context, baseID int, kind domain.ActivityKind, opID int) (bool, error) {
	return r.q.ExistsForOperation(ctx, gen.ExistsForOperationParams{
		BaseID: int64(baseID),
		Kind:   string(kind),
		OpID:   int64(opID),
	})
}

func (r *ActivityRepo) ExistsForScanReport(ctx context.Context, reportID int) (bool, error) {
	return r.q.ExistsForScanReport(ctx, int64(reportID))
}
