package ports

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/auth/domain"
	"github.com/google/uuid"
)

type OutboxEventRepository interface {
	Save(ctx context.Context, events []domain.DomainEvent) error
	ClaimUnpublished(limit int) ([]domain.DomainEvent, error)
	MarkPublished(id uuid.UUID, publishedAt int64) error
	Tx(tx Transaction) OutboxEventRepository
}
