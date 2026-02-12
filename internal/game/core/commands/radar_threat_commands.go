package commands

import (
	"github.com/artcodefun/heat-expansion-api/internal/game/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/game/core/ports"
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
func (c *RadarThreatCommands) HandleMilitaryOperationArrivedEvent(event domain.MilitaryOperationArrivedEvent) error {
	return c.TxMgr.WithTx(func(tx ports.Transaction) error {
		threat, err := c.RadarThreatRepo.Tx(tx).FindByOperationID(event.OperationID)
		if err != nil {
			if err == ports.ErrNotFound {
				return nil
			}
			return err
		}

		threat.MarkArrived()
		return c.RadarThreatRepo.Tx(tx).Update(threat)
	})
}

// HandleMilitaryOperationCancelledEvent updates the radar threat status when the operation is cancelled.
func (c *RadarThreatCommands) HandleMilitaryOperationCancelledEvent(event domain.MilitaryOperationCancelledEvent) error {
	return c.TxMgr.WithTx(func(tx ports.Transaction) error {
		threat, err := c.RadarThreatRepo.Tx(tx).FindByOperationID(event.OperationID)
		if err != nil {
			if err == ports.ErrNotFound {
				return nil
			}
			return err
		}

		threat.MarkLost()
		return c.RadarThreatRepo.Tx(tx).Update(threat)
	})
}
