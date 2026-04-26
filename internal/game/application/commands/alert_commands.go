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
	TradeRepo ports.TradeOperationRepository
	TxMgr     ports.TransactionManager
}

func NewAlertCommands(repo ports.AlertRepository, tradeRepo ports.TradeOperationRepository, txMgr ports.TransactionManager) *AlertCommands {
	return &AlertCommands{
		AlertRepo: repo,
		TradeRepo: tradeRepo,
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

func (c *AlertCommands) HandleTradeOperationCreatedEvent(ctx context.Context, e domain.TradeOperationCreatedEvent) error {
	return c.createTradeAlerts(ctx, e.OperationID, domain.TradeAlertKindCreated)
}

func (c *AlertCommands) HandleTradeOperationAcceptedEvent(ctx context.Context, e domain.TradeOperationAcceptedEvent) error {
	return c.createTradeAlerts(ctx, e.OperationID, domain.TradeAlertKindAccepted)
}

func (c *AlertCommands) HandleTradeOperationDeclinedEvent(ctx context.Context, e domain.TradeOperationDeclinedEvent) error {
	return c.createTradeAlerts(ctx, e.OperationID, domain.TradeAlertKindDeclined)
}

func (c *AlertCommands) HandleTradeOperationCancelledByInitiatorEvent(ctx context.Context, e domain.TradeOperationCancelledByInitiatorEvent) error {
	return c.createTradeAlerts(ctx, e.OperationID, domain.TradeAlertKindCancelled)
}

func (c *AlertCommands) HandleTradeOperationExpiredEvent(ctx context.Context, e domain.TradeOperationExpiredEvent) error {
	return c.createTradeAlerts(ctx, e.OperationID, domain.TradeAlertKindExpired)
}

func (c *AlertCommands) HandleTradeOperationReturnArrivedEvent(ctx context.Context, e domain.TradeOperationReturnArrivedEvent) error {
	return c.createTradeAlerts(ctx, e.OperationID, domain.TradeAlertKindCompleted)
}

func (c *AlertCommands) createTradeAlerts(ctx context.Context, operationID int, kind domain.TradeAlertKind) error {
	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		op, err := c.TradeRepo.Tx(tx).FindByID(ctx, operationID)
		if err != nil {
			if err == ports.ErrNotFound {
				return nil
			}
			return repoErr(err)
		}
		alerts := tradeAlertsForOperation(op, kind)
		for _, alert := range alerts {
			if err := c.AlertRepo.Tx(tx).Create(ctx, alert); err != nil {
				log.Printf("Failed to create trade alert for operation %d: %v", operationID, err)
				return err
			}
		}
		return nil
	})
}

func tradeAlertsForOperation(op *domain.TradeOperation, kind domain.TradeAlertKind) []*domain.Alert {
	if op == nil {
		return nil
	}

	switch kind {
	case domain.TradeAlertKindCreated:
		return []*domain.Alert{domain.NewTradeAlert(op, false, kind)}
	case domain.TradeAlertKindAccepted, domain.TradeAlertKindDeclined:
		return []*domain.Alert{domain.NewTradeAlert(op, true, kind)}
	case domain.TradeAlertKindCancelled:
		return []*domain.Alert{domain.NewTradeAlert(op, false, kind)}
	case domain.TradeAlertKindExpired, domain.TradeAlertKindCompleted:
		return []*domain.Alert{
			domain.NewTradeAlert(op, true, kind),
			domain.NewTradeAlert(op, false, kind),
		}
	default:
		return nil
	}
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
