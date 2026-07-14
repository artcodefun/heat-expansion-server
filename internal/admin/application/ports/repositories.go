package ports

import (
	"context"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/admin/domain"
	"github.com/google/uuid"
)

var ErrNotFound = errors.New("not found")

type AdminRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Admin, error)
	FindByUsername(ctx context.Context, username string) (*domain.Admin, error)
	Save(ctx context.Context, admin *domain.Admin) error
	Tx(tx Transaction) AdminRepository
}

type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	FindByToken(ctx context.Context, token string) (*domain.Session, error)
	Delete(ctx context.Context, token string) error
	Tx(tx Transaction) SessionRepository
}
