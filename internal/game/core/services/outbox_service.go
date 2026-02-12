package services

import (
	"time"

	"github.com/artcodefun/heat-expansion-api/internal/game/core/ports"
)

// OutboxService encapsulates the application-level logic for dispatching
// domain events from the transactional outbox to the EventPublisher.
//
// It is intended to be run by a background worker that periodically calls
// ProcessBatch.
type OutboxService struct {
	Outbox    ports.OutboxEventRepository
	Publisher ports.EventPublisher
	TxMgr     ports.TransactionManager
}

func NewOutboxService(outbox ports.OutboxEventRepository, publisher ports.EventPublisher, txMgr ports.TransactionManager) *OutboxService {
	return &OutboxService{Outbox: outbox, Publisher: publisher, TxMgr: txMgr}
}

// ProcessBatch claims up to the given limit of unpublished events using
// database-level locking, publishes them via the EventPublisher, and marks
// successfully published events as published. It executes all operations
// within a single transaction provided by the TransactionManager.
func (s *OutboxService) ProcessBatch(limit int) error {
	if limit <= 0 {
		limit = 100
	}

	return s.TxMgr.WithTx(func(tx ports.Transaction) error {
		repo := s.Outbox.Tx(tx)

		records, err := repo.ClaimUnpublished(limit)
		if err != nil {
			return err
		}
		if len(records) == 0 {
			return nil
		}

		now := time.Now().Unix()
		for _, r := range records {
			if err := s.Publisher.Publish(r.Event); err != nil {
				// If publishing fails, skip marking this event as published so it
				// can be retried in a subsequent batch.
				continue
			}

			if err := repo.MarkPublished(r.ID, now); err != nil {
				return err
			}
		}

		return nil
	})
}
