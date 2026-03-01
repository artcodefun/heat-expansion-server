package commands

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
)

type RadarThreatCommands struct {
	RadarThreatRepo ports.RadarThreatRepository
	Outbox          ports.OutboxEventRepository
	TxMgr           ports.TransactionManager
}

func NewRadarThreatCommands(radarThreatRepo ports.RadarThreatRepository, outbox ports.OutboxEventRepository, txMgr ports.TransactionManager) *RadarThreatCommands {
	return &RadarThreatCommands{
		RadarThreatRepo: radarThreatRepo,
		Outbox:          outbox,
		TxMgr:           txMgr,
	}
}

// HandleMilitaryOperationArrivedEvent updates the radar threat status when the operation arrives at target.
func (c *RadarThreatCommands) HandleMilitaryOperationArrivedEvent(ctx context.Context, event domain.MilitaryOperationArrivedEvent) error {
	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		threat, err := c.RadarThreatRepo.Tx(tx).FindByOperationID(ctx, event.OperationID)
		if err != nil {
			if err == ports.ErrNotFound {
				return nil
			}
			return err
		}

		threat.MarkArrived()
		return c.RadarThreatRepo.Tx(tx).Update(ctx, threat)
	})
}

// HandleMilitaryOperationCancelledEvent updates the radar threat status when the operation is cancelled.
func (c *RadarThreatCommands) HandleMilitaryOperationCancelledEvent(ctx context.Context, event domain.MilitaryOperationCancelledEvent) error {
	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		threat, err := c.RadarThreatRepo.Tx(tx).FindByOperationID(ctx, event.OperationID)
		if err != nil {
			if err == ports.ErrNotFound {
				return nil
			}
			return err
		}

		threat.MarkLost()
		return c.RadarThreatRepo.Tx(tx).Update(ctx, threat)
	})
}
