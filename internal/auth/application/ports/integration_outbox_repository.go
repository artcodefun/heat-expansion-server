package ports

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/contracts/auth"
	"github.com/google/uuid"
)

type IntegrationOutboxRepository interface {
	Save(ctx context.Context, event auth.IntegrationEvent) error
	Exists(ctx context.Context, originID uuid.UUID, eventType string) (bool, error)
	ClaimUnpublished(ctx context.Context, limit int) ([]auth.IntegrationEvent, error)
	MarkPublished(ctx context.Context, id uuid.UUID, publishedAt int64) error
	Tx(tx Transaction) IntegrationOutboxRepository
}
