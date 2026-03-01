package services

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/contracts/auth"
	v1 "github.com/artcodefun/heat-expansion-server/contracts/auth/v1"
	"github.com/artcodefun/heat-expansion-server/internal/auth/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/auth/domain"
)

type IntegrationProducerService struct {
	outbox ports.IntegrationOutboxRepository
	txMgr  ports.TransactionManager
}

func NewIntegrationProducerService(outbox ports.IntegrationOutboxRepository, txMgr ports.TransactionManager) *IntegrationProducerService {
	return &IntegrationProducerService{outbox: outbox, txMgr: txMgr}
}

func (s *IntegrationProducerService) HandleAccountRegistered(ctx context.Context, ev domain.AccountRegisteredEvent) error {
	return s.txMgr.WithTx(ctx, func(tx ports.Transaction) error {
		outbox := s.outbox.Tx(tx)

		originID := ev.ID()
		eventType := v1.EventAccountRegisteredV1

		// Check idempotency
		exists, err := outbox.Exists(ctx, originID, eventType)
		if err != nil {
			return err
		}
		if exists {
			return nil
		}

		integrationEvent := auth.NewIntegrationEvent(
			originID,
			ev.OccurredAt(),
			v1.AccountRegisteredV1{
				UserID: ev.AccountID,
				Name:   ev.Name,
				Email:  ev.Email,
			},
		)
		return outbox.Save(ctx, integrationEvent)
	})
}
