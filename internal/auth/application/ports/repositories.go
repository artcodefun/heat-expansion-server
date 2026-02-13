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
	Tx(tx Transaction) AccountRepository
}
