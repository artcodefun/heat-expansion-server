package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/billing/domain"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/db/mappers"
	"github.com/google/uuid"
)

type UserRepo struct {
	q *gen.Queries
}

func NewUserRepo(q *gen.Queries) *UserRepo {
	return &UserRepo{q: q}
}

func (r *UserRepo) Tx(tx ports.Transaction) ports.UserRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &UserRepo{q: r.q.WithTx(sqlTx)}
	}
	return r
}

func (r *UserRepo) Upsert(ctx context.Context, user *domain.User) error {
	return r.q.UpsertUser(ctx, gen.UpsertUserParams{
		ID:    user.ID,
		Email: user.Email,
	})
}

func (r *UserRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	row, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.UserFromRow(row), nil
}
