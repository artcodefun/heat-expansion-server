package services

import (
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

func (s *IntegrationOutboxService) ProcessBatch(limit int) error {
	if limit <= 0 {
		limit = 100
	}

	return s.TxMgr.WithTx(func(tx ports.Transaction) error {
		repo := s.Outbox.Tx(tx)

		events, err := repo.ClaimUnpublished(limit)
		if err != nil {
			return err
		}
		if len(events) == 0 {
			return nil
		}

		now := time.Now().Unix()
		for _, evt := range events {
			if err := s.Publisher.Publish(evt); err != nil {
				continue
			}

			if err := repo.MarkPublished(evt.ID, now); err != nil {
				return err
			}
		}

		return nil
	})
}
