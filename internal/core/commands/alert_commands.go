package commands

import (
	"log"
	"time"

	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/core/services"
)

type AlertCommands struct {
	AlertRepo ports.AlertRepository
	Access    *services.AccessControlService
}

func NewAlertCommands(repo ports.AlertRepository, access *services.AccessControlService) *AlertCommands {
	return &AlertCommands{
		AlertRepo: repo,
		Access:    access,
	}
}

func (c *AlertCommands) MarkAllAsRead(baseID int, userID int) error {
	if err := c.Access.EnsureBaseOwnership(userID, baseID); err != nil {
		return err
	}
	return c.AlertRepo.MarkAllAsRead(baseID)
}

func (c *AlertCommands) HandleActivityCreatedEvent(e domain.ActivityCreatedEvent) error {
	var kind domain.AlertKind
	var title, content string
	var ttl time.Duration = 24 * 3 * time.Hour // 3 days default

	switch e.Kind {
	case domain.ActivityKindDefense:
		kind = domain.AlertKindCombat
		title = "Base Under Attack"
		if e.Subtype == string(domain.DefenseActivitySubtypeSpy) {
			content = "Spies have been noticed inside the base!"
		} else {
			content = "Foreign army has attacked the base!"
		}
	case domain.ActivityKindScan:
		if e.Subtype != string(domain.ScanActivitySubtypeExternalScanDetected) {
			return nil // No alert for reports produced by the player
		}
		kind = domain.AlertKindIntel
		title = "External Scan Detected"
		content = "Your sensors detected an external scan targeting your base!"
	case domain.ActivityKindRadar:
		kind = domain.AlertKindIntel
		title = "Incoming Threat Detected"
		content = "Radars have detected an incoming threat!"
	default:
		return nil // No alert for other kinds
	}

	alert := domain.NewAlert(e.BaseID, &e.ActivityID, kind, title, content, ttl)
	err := c.AlertRepo.Create(alert)
	if err != nil {
		log.Printf("Failed to create alert for activity %s: %v", e.ActivityID, err)
	}
	return err
}
