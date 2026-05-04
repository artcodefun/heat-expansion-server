package repo

import (
	"context"
	"database/sql"

	"github.com/artcodefun/heat-expansion-server/internal/auth/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/auth/domain"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/db/gen"
	"github.com/google/uuid"
)

type PasswordResetRepository struct {
	db *gen.Queries
}

func NewPasswordResetRepository(db *sql.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: gen.New(db)}
}

func (r *PasswordResetRepository) Tx(tx ports.Transaction) ports.PasswordResetRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &PasswordResetRepository{db: r.db.WithTx(sqlTx)}
	}
	return r
}

func (r *PasswordResetRepository) Create(ctx context.Context, token *domain.PasswordResetToken) error {
	return r.db.CreatePasswordResetToken(ctx, gen.CreatePasswordResetTokenParams{
		ID:        token.ID,
		AccountID: token.AccountID,
		TokenHash: token.TokenHash,
		ExpiresAt: token.ExpiresAt,
	})
}

func (r *PasswordResetRepository) FindByAccountAndTokenHash(ctx context.Context, accountID uuid.UUID, tokenHash string) (*domain.PasswordResetToken, error) {
	row, err := r.db.GetActivePasswordResetToken(ctx, gen.GetActivePasswordResetTokenParams{
		AccountID: accountID,
		TokenHash: tokenHash,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
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
	return r.db.MarkPasswordResetTokenUsed(ctx, gen.MarkPasswordResetTokenUsedParams{
		ID:     id,
		UsedAt: sql.NullInt64{Int64: usedAt, Valid: true},
	})
}
