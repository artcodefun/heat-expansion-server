package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/artcodefun/heat-expansion-server/internal/auth/application/ports"
)

type IntegrationOutboxService struct {
	Outbox    ports.IntegrationOutboxRepository
	Publisher ports.IntegrationEventPublisher
	TxMgr     ports.TransactionManager
}

func NewIntegrationOutboxService(
	outbox ports.IntegrationOutboxRepository,
	publisher ports.IntegrationEventPublisher,
	txMgr ports.TransactionManager,
) *IntegrationOutboxService {
	return &IntegrationOutboxService{
		Outbox:    outbox,
		Publisher: publisher,
		TxMgr:     txMgr,
	}
}

func (s *IntegrationOutboxService) ProcessBatch(ctx context.Context, limit int) error {
	if limit <= 0 {
		limit = 100
	}

	return s.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		repo := s.Outbox.Tx(tx)

		events, err := repo.ClaimUnpublished(ctx, limit)
		if err != nil {
			return err
		}
		if len(events) == 0 {
			return nil
		}

		now := time.Now().Unix()
		for _, evt := range events {
			if err := s.Publisher.Publish(ctx, evt); err != nil {
				slog.WarnContext(ctx, "auth integration outbox event publish failed; leaving event unpublished for retry",
					"event_id", evt.ID.String(),
					"event_type", evt.Type,
					"origin_id", evt.OriginID.String(),
					"error", err.Error(),
				)
				continue
			}

			if err := repo.MarkPublished(ctx, evt.ID, now); err != nil {
				return err
			}
		}

		return nil
	})
}
