package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/auth/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/auth/domain"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/db/gen"
	"github.com/google/uuid"
)

type PasswordResetRepository struct {
	q *gen.Queries
}

func NewPasswordResetRepository(q *gen.Queries) *PasswordResetRepository {
	return &PasswordResetRepository{q: q}
}

func (r *PasswordResetRepository) Tx(tx ports.Transaction) ports.PasswordResetRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &PasswordResetRepository{q: r.q.WithTx(sqlTx)}
	}
	return r
}

func (r *PasswordResetRepository) Create(ctx context.Context, token *domain.PasswordResetToken) error {
	return r.q.CreatePasswordResetToken(ctx, gen.CreatePasswordResetTokenParams{
		ID:        token.ID,
		AccountID: token.AccountID,
		TokenHash: token.TokenHash,
		ExpiresAt: token.ExpiresAt,
	})
}

func (r *PasswordResetRepository) FindByAccountAndTokenHash(ctx context.Context, accountID uuid.UUID, tokenHash string) (*domain.PasswordResetToken, error) {
	row, err := r.q.GetActivePasswordResetToken(ctx, gen.GetActivePasswordResetTokenParams{
		AccountID: accountID,
		TokenHash: tokenHash,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}

	t := &domain.PasswordResetToken{
		ID:        row.ID,
		AccountID: row.AccountID,
		TokenHash: row.TokenHash,
		ExpiresAt: row.ExpiresAt,
	}
	if row.UsedAt.Valid {
		t.UsedAt = &row.UsedAt.Int64
	}
	return t, nil
}

func (r *PasswordResetRepository) MarkUsed(ctx context.Context, id uuid.UUID, usedAt int64) error {
	return r.q.MarkPasswordResetTokenUsed(ctx, gen.MarkPasswordResetTokenUsedParams{
		ID:     id,
		UsedAt: sql.NullInt64{Int64: usedAt, Valid: true},
	})
}

func (r *PasswordResetRepository) InvalidateByAccount(ctx context.Context, accountID uuid.UUID, usedAt int64) error {
	return r.q.InvalidateAccountPasswordResetTokens(ctx, gen.InvalidateAccountPasswordResetTokensParams{
		AccountID: accountID,
		UsedAt:    sql.NullInt64{Int64: usedAt, Valid: true},
	})
}
