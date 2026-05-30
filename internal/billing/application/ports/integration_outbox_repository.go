package ports

import (
	"context"

	billingevents "github.com/artcodefun/heat-expansion-server/contracts/billing/events"
	"github.com/google/uuid"
)

type IntegrationOutboxRepository interface {
	Save(ctx context.Context, event billingevents.IntegrationEvent) error
	Exists(ctx context.Context, originID uuid.UUID, eventType string) (bool, error)
	ClaimUnpublished(ctx context.Context, limit int) ([]billingevents.IntegrationEvent, error)
	MarkPublished(ctx context.Context, id uuid.UUID, publishedAt int64) error
	Tx(tx Transaction) IntegrationOutboxRepository
}
