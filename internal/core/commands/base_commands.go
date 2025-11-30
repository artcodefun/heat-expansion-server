package commands

import (
	"fmt"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/core/services"
)

// BaseCommands encapsulates state-changing base operations.
type BaseCommands struct {
	UserBaseRepo      ports.UserBaseRepository
	SectorRepo        ports.SectorRepository
	ContentGenerator  ports.ContentGenerator
	SectorProvisioner *services.SectorProvisioningService
	basePlacement     *domain.BasePlacementService
	TxMgr             ports.TransactionManager
}

func NewBaseCommands(userBaseRepo ports.UserBaseRepository, sectorRepo ports.SectorRepository, generator ports.ContentGenerator, provisioner *services.SectorProvisioningService, txMgr ports.TransactionManager) *BaseCommands {
	return &BaseCommands{UserBaseRepo: userBaseRepo, SectorRepo: sectorRepo, ContentGenerator: generator, SectorProvisioner: provisioner, basePlacement: domain.NewBasePlacementService(), TxMgr: txMgr}
}

// CreateBase creates a new base for a user.
func (c *BaseCommands) CreateBase(ctx cqrs.CommandContext, userID int) error {
	return c.TxMgr.WithTx(func(tx ports.Transaction) error {
		sRepo := c.SectorRepo.Tx(tx)
		bRepo := c.UserBaseRepo.Tx(tx)
		const maxAttempts = 10
		for attempt := 0; attempt < maxAttempts; attempt++ {
			occupied, err := sRepo.ListOccupiedCoordinates()
			if err != nil {
				return err
			}
			x, y := c.basePlacement.FindFreeChunkForBase(occupied)
			base := domain.NewUserBaseModel(0, userID, domain.Vector2i{X: x, Y: y})
			created, err := c.SectorProvisioner.CreateUserBaseIfEmpty(sRepo, bRepo, base)
			if err != nil {
				return err
			}
			if created {
				return nil
			}
		}
		return fmt.Errorf("no free sector after attempts")
	})
}

// HandleUserAccountCreatedEvent reacts to user creation.
func (c *BaseCommands) HandleUserAccountCreatedEvent(ev domain.UserAccountCreatedEvent) error {
	return c.CreateBase(cqrs.CommandContext{}, ev.UserID)
}
