package ports

import (
"context"
"github.com/google/uuid"
"github.com/artcodefun/heat-expansion-server/internal/auth/domain"
)

type AccountRepository interface {
	Create(ctx context.Context, account *domain.Account) error
	FindByEmail(ctx context.Context, email string) (*domain.Account, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Account, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, newHash string) error
	Tx(tx Transaction) AccountRepository
}

type PasswordResetRepository interface {
	Create(ctx context.Context, token *domain.PasswordResetToken) error
	FindByAccountAndTokenHash(ctx context.Context, accountID uuid.UUID, tokenHash string) (*domain.PasswordResetToken, error)
	MarkUsed(ctx context.Context, id uuid.UUID, usedAt int64) error
	InvalidateByAccount(ctx context.Context, accountID uuid.UUID, usedAt int64) error
	Tx(tx Transaction) PasswordResetRepository
}
