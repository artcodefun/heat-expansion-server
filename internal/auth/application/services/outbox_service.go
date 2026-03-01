package services

import (
	"context"
	"time"

	"github.com/artcodefun/heat-expansion-server/internal/auth/application/ports"
)

// OutboxService encapsulates the application-level logic for dispatching
// domain events from the transactional outbox to the EventPublisher.
type OutboxService struct {
	Outbox    ports.OutboxEventRepository
	Publisher ports.EventPublisher
	TxMgr     ports.TransactionManager
}

func NewOutboxService(outbox ports.OutboxEventRepository, publisher ports.EventPublisher, txMgr ports.TransactionManager) *OutboxService {
	return &OutboxService{Outbox: outbox, Publisher: publisher, TxMgr: txMgr}
}

// ProcessBatch claims up to the given limit of unpublished events,
// publishes them via the EventPublisher, and marks successfully published
// events as published.
func (s *OutboxService) ProcessBatch(ctx context.Context, limit int) error {
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
				// If publishing fails, skip marking this event as published
				continue
			}

			if err := repo.MarkPublished(ctx, evt.ID(), now); err != nil {
				return err
			}
		}

		return nil
	})
}
