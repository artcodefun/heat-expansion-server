package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/admin/domain"
	"github.com/artcodefun/heat-expansion-server/internal/admin/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/admin/infrastructure/db/mappers"
)

type SessionRepository struct {
	q *gen.Queries
}

func NewSessionRepository(q *gen.Queries) *SessionRepository {
	return &SessionRepository{q: q}
}

func (r *SessionRepository) Tx(tx ports.Transaction) ports.SessionRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &SessionRepository{q: r.q.WithTx(sqlTx)}
	}
	return r
}

func (r *SessionRepository) Create(ctx context.Context, session *domain.Session) error {
	return r.q.CreateSession(ctx, mappers.CreateSessionParamsFromDomain(session))
}

func (r *SessionRepository) FindByToken(ctx context.Context, token string) (*domain.Session, error) {
	row, err := r.q.GetSessionByToken(ctx, gen.GetSessionByTokenParams{
		Token: token,
		Now:   domain.NowUnix(),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.SessionFromDB(row), nil
}

func (r *SessionRepository) Delete(ctx context.Context, token string) error {
	return r.q.DeleteSession(ctx, token)
}
