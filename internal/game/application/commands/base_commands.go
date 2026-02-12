package commands

import (
	"fmt"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
)

// BaseCommands encapsulates state-changing base operations.
type BaseCommands struct {
	UserBaseRepo      ports.UserBaseRepository
	SectorRepo        ports.SectorRepository
	BuildRepo         ports.BuildPrototypeRepository
	ArmyRepo          ports.ArmyPrototypeRepository
	ContentGenerator  ports.ContentGenerator
	SectorProvisioner *services.SectorProvisioningService
	Outbox            ports.OutboxEventRepository
	basePlacement     *domain.BasePlacementService
	TxMgr             ports.TransactionManager
}

func NewBaseCommands(userBaseRepo ports.UserBaseRepository, sectorRepo ports.SectorRepository, buildRepo ports.BuildPrototypeRepository, armyRepo ports.ArmyPrototypeRepository, generator ports.ContentGenerator, provisioner *services.SectorProvisioningService, outbox ports.OutboxEventRepository, txMgr ports.TransactionManager) *BaseCommands {
	return &BaseCommands{UserBaseRepo: userBaseRepo, SectorRepo: sectorRepo, BuildRepo: buildRepo, ArmyRepo: armyRepo, ContentGenerator: generator, SectorProvisioner: provisioner, Outbox: outbox, basePlacement: domain.NewBasePlacementService(), TxMgr: txMgr}
}

// CreateBase creates a new base for a user.
func (c *BaseCommands) CreateBase(ctx cqrs.CommandContext, userID int) error {
	return c.TxMgr.WithTx(func(tx ports.Transaction) error {
		sRepo := c.SectorRepo.Tx(tx)
		bRepo := c.UserBaseRepo.Tx(tx)
		biRepo := c.BuildRepo.Tx(tx)
		aiRepo := c.ArmyRepo.Tx(tx)

		// 1. Fetch starter prototypes
		allBuildProtos, err := biRepo.FindAllPrototypes()
		if err != nil {
			return repoErr(err)
		}
		allArmyProtos, err := aiRepo.FindAllPrototypes()
		if err != nil {
			return repoErr(err)
		}

		const maxAttempts = 10
		for attempt := 0; attempt < maxAttempts; attempt++ {
			occupied, err := sRepo.ListOccupiedCoordinates()
			if err != nil {
				return repoErr(err)
			}
			x, y := c.basePlacement.FindFreeChunkForBase(occupied)
			base := domain.NewUserBaseModel(0, userID, domain.Vector2i{X: x, Y: y})

			// Add starter buildings and units via domain logic
			base.EnsureStartingBuildingsPresent(allBuildProtos)
			base.EnsureStartingArmyPresent(allArmyProtos)
			// Fill up starter resources
			base.FillStarterResources()

			created, err := c.SectorProvisioner.CreateUserBaseIfEmpty(sRepo, bRepo, base)
			if err != nil {
				return err
			}
			if created {
				base.EmitCreated()
				return c.Outbox.Tx(tx).Save(base.PullEvents())
			}
		}
		return fmt.Errorf("no free sector after attempts")
	})
}

// HandleUserAccountCreatedEvent reacts to user creation.
func (c *BaseCommands) HandleUserAccountCreatedEvent(ev domain.UserAccountCreatedEvent) error {
	return c.CreateBase(cqrs.CommandContext{}, ev.UserID)
}
