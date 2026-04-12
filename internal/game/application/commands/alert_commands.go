package commands

import (
	"context"
	"log"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/google/uuid"
)

type AlertCommands struct {
	AlertRepo ports.AlertRepository
	TxMgr     ports.TransactionManager
}

func NewAlertCommands(repo ports.AlertRepository, txMgr ports.TransactionManager) *AlertCommands {
	return &AlertCommands{
		AlertRepo: repo,
		TxMgr:     txMgr,
	}
}

func (c *AlertCommands) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	if userID == uuid.Nil {
		return cqrs.ErrForbidden
	}
	return c.AlertRepo.MarkAllAsRead(ctx, userID)
}

func (c *AlertCommands) HandleActivityCreatedEvent(ctx context.Context, e domain.ActivityCreatedEvent) error {
	alert, ok := domain.NewActivityAlert(e)
	if !ok {
		return nil
	}

	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		repo := c.AlertRepo.Tx(tx)
		if exists, _ := repo.ExistsForActivity(ctx, e.ActivityID); exists {
			return nil
		}
		err := repo.Create(ctx, alert)
		if err != nil {
			log.Printf("Failed to create alert for activity %s: %v", e.ActivityID, err)
		}
		return err
	})
}

func (c *AlertCommands) HandleDiplomaticMessageSentEvent(ctx context.Context, e domain.DiplomaticMessageSentEvent) error {
	alert, ok := domain.NewDiplomaticMessageAlert(e)
	if !ok {
		return nil
	}

	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		return c.AlertRepo.Tx(tx).Create(ctx, alert)
	})
}

func (c *AlertCommands) HandleDiplomaticRequestCreatedEvent(ctx context.Context, e domain.DiplomaticRequestCreatedEvent) error {
	alert, ok := domain.NewDiplomaticRequestAlert(e)
	if !ok {
		return nil
	}

	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		return c.AlertRepo.Tx(tx).Create(ctx, alert)
	})
}
