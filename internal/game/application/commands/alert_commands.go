package commands

import (
	"context"
	"log"
	"time"

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
	var kind domain.AlertKind
	var title, content string
	var ttl time.Duration = 24 * 3 * time.Hour // 3 days default

	switch e.Kind {
	case domain.ActivityKindDefense:
		kind = domain.AlertKindCombat
		title = "alert.combat.attack.title"
		if e.Subtype == string(domain.DefenseActivitySubtypeSpy) {
			content = "alert.combat.spy.content"
		} else {
			content = "alert.combat.attack.content"
		}
	case domain.ActivityKindScan:
		if e.Subtype != string(domain.ScanActivitySubtypeExternalScanDetected) {
			return nil // No alert for reports produced by the player
		}
		kind = domain.AlertKindIntel
		title = "alert.intel.scan.title"
		content = "alert.intel.scan.content"
	case domain.ActivityKindRadar:
		kind = domain.AlertKindIntel
		title = "alert.intel.threat.title"
		content = "alert.intel.threat.content"
	default:
		return nil // No alert for other kinds
	}

	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		repo := c.AlertRepo.Tx(tx)
		if exists, _ := repo.ExistsForActivity(ctx, e.ActivityID); exists {
			return nil
		}
		alert := domain.NewAlert(e.UserID, e.BaseID, &e.ActivityID, kind, title, content, ttl)
		err := repo.Create(ctx, alert)
		if err != nil {
			log.Printf("Failed to create alert for activity %s: %v", e.ActivityID, err)
		}
		return err
	})
}
