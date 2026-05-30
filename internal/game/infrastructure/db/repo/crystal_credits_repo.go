package repo

import (
	"context"
	"database/sql"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
	"github.com/google/uuid"
)

type CrystalCreditsRepo struct {
	q *gen.Queries
}

func NewCrystalCreditsRepo(q *gen.Queries) *CrystalCreditsRepo {
	return &CrystalCreditsRepo{q: q}
}

func (r *CrystalCreditsRepo) Tx(tx ports.Transaction) ports.CrystalCreditsRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &CrystalCreditsRepo{q: r.q.WithTx(sqlTx)}
	}
	return r
}

func (r *CrystalCreditsRepo) Insert(ctx context.Context, orderID uuid.UUID, userID uuid.UUID, crystals int, creditedAt int64) error {
	return r.q.InsertCrystalCredit(ctx, gen.InsertCrystalCreditParams{
		OrderID:    orderID,
		UserID:     userID,
		Crystals:   int32(crystals),
		CreditedAt: creditedAt,
	})
}

func (r *CrystalCreditsRepo) Exists(ctx context.Context, orderID uuid.UUID) (bool, error) {
	return r.q.CrystalCreditExists(ctx, orderID)
}
