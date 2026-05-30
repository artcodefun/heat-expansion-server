package services

import (
	"context"

	billingevents "github.com/artcodefun/heat-expansion-server/contracts/billing/events"
	v1 "github.com/artcodefun/heat-expansion-server/contracts/billing/events/v1"
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/billing/domain"
)

type IntegrationProducerService struct {
	outbox ports.IntegrationOutboxRepository
	txMgr  ports.TransactionManager
}

func NewIntegrationProducerService(outbox ports.IntegrationOutboxRepository, txMgr ports.TransactionManager) *IntegrationProducerService {
	return &IntegrationProducerService{outbox: outbox, txMgr: txMgr}
}

func (s *IntegrationProducerService) HandleOrderPaid(ctx context.Context, ev domain.OrderPaidEvent) error {
	return s.txMgr.WithTx(ctx, func(tx ports.Transaction) error {
		outbox := s.outbox.Tx(tx)

		originID := ev.ID()
		eventType := v1.EventCrystalsPurchasedV1

		exists, err := outbox.Exists(ctx, originID, eventType)
		if err != nil {
			return err
		}
		if exists {
			return nil
		}

		integrationEvent := billingevents.NewIntegrationEvent(
			originID,
			ev.OccurredAt(),
			v1.CrystalsPurchasedV1{
				UserID:    ev.UserID,
				OrderID:   ev.OrderID,
				PackageID: ev.PackageID,
				Crystals:  ev.Crystals,
			},
		)
		return outbox.Save(ctx, integrationEvent)
	})
}
