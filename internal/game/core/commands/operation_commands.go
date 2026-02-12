package commands

import (
	"errors"
	"fmt"

	"github.com/artcodefun/heat-expansion-api/internal/game/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/game/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/game/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/game/core/services"
)

type OperationCommands struct {
	UserBaseRepo   ports.UserBaseRepository
	UserRepo       ports.UserRepository
	SectorRepo     ports.SectorRepository
	OperationRepo  ports.MilitaryOperationRepository
	ResourceRepo   ports.ResourceLocationRepository
	DangerousRepo  ports.DangerousLocationRepository
	ScanReportRepo ports.ScanReportRepository
	StorageProtos  ports.StoragePrototypeRepository
	Provisioner    *services.SectorProvisioningService
	Scheduler      ports.Scheduler
	Outbox         ports.OutboxEventRepository
	TxMgr          ports.TransactionManager
	Access         *services.AccessControlService
	crystalService *domain.CrystalSpendingService
}

func NewOperationCommands(userBaseRepo ports.UserBaseRepository, userRepo ports.UserRepository, sectorRepo ports.SectorRepository, opRepo ports.MilitaryOperationRepository, resRepo ports.ResourceLocationRepository, dangerRepo ports.DangerousLocationRepository, scanRepo ports.ScanReportRepository, storageProtos ports.StoragePrototypeRepository, provisioner *services.SectorProvisioningService, scheduler ports.Scheduler, outbox ports.OutboxEventRepository, txMgr ports.TransactionManager, access *services.AccessControlService) *OperationCommands {
	return &OperationCommands{
		UserBaseRepo:   userBaseRepo,
		UserRepo:       userRepo,
		SectorRepo:     sectorRepo,
		OperationRepo:  opRepo,
		ResourceRepo:   resRepo,
		DangerousRepo:  dangerRepo,
		ScanReportRepo: scanRepo,
		StorageProtos:  storageProtos,
		Provisioner:    provisioner,
		Scheduler:      scheduler,
		Outbox:         outbox,
		TxMgr:          txMgr,
		Access:         access,
		crystalService: domain.NewCrystalSpendingService(),
	}
}

func (c *OperationCommands) CreateMilitaryOperation(ctx cqrs.CommandContext, opType domain.MilitaryOperationType, sourceBaseID int, targetX int, targetY int, deployments []domain.ArmyDeploymentRequest) (*domain.MilitaryOperation, error) {
	if err := c.Access.EnsureBaseOwnership(ctx.UserID, sourceBaseID); err != nil {
		return nil, err
	}
	var createdOp *domain.MilitaryOperation
	err := c.TxMgr.WithTx(func(tx ports.Transaction) error {
		bRepo := c.UserBaseRepo.Tx(tx)
		sRepo := c.SectorRepo.Tx(tx)
		oRepo := c.OperationRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(sourceBaseID)
		if err != nil {
			return repoErr(err)
		}
		readyToDeploy, err := base.GetReadyToDeployArmy(deployments)
		if err != nil {
			return cqrs.NewDomainError(err)
		}

		snaps := base.ActiveStorageSnaps()
		units := domain.MilitaryUnitsFromDeployed(readyToDeploy)
		sourceSector, err := c.Provisioner.EnsureSectorExists(sRepo, base.Coordinates.X, base.Coordinates.Y)
		if err != nil {
			return err
		}
		targetSector, err := c.Provisioner.EnsureSectorExists(sRepo, targetX, targetY)
		if err != nil {
			return err
		}
		var opCreationErr error
		switch opType {
		case domain.MilitaryOperationTypeAttack:
			createdOp, opCreationErr = domain.NewAttackOperation(base.UserID, sourceBaseID, sourceSector.Coordinates, targetSector.Coordinates, units, snaps)
		case domain.MilitaryOperationTypeSpy:
			createdOp, opCreationErr = domain.NewSpyOperation(base.UserID, sourceBaseID, sourceSector.Coordinates, targetSector.Coordinates, units, snaps)
		default:
			return errors.New("unsupported operation type")
		}
		if opCreationErr != nil {
			return cqrs.NewDomainError(opCreationErr)
		}
		if err := oRepo.Create(createdOp); err != nil {
			return err
		}
		for _, ready := range readyToDeploy {
			if _, err := base.AllocateArmyToOperation(domain.ArmyDeploymentRequest{PresentItemID: ready.PresentItemID, Count: ready.Count}, createdOp.ID); err != nil {
				return cqrs.NewDomainError(err)
			}
		}
		if err := bRepo.Update(base); err != nil {
			return err
		}
		createdOp.Start()
		if err := oRepo.Update(createdOp); err != nil {
			return err
		}
		if err := c.Outbox.Tx(tx).Save(createdOp.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return createdOp, nil
}

func (c *OperationCommands) CancelMilitaryOperation(ctx cqrs.CommandContext, operationID int) error {
	err := c.TxMgr.WithTx(func(tx ports.Transaction) error {
		oRepo := c.OperationRepo.Tx(tx)
		op, err := oRepo.FindByIDForUpdate(operationID)
		if err != nil {
			return repoErr(err)
		}

		if err := c.Access.EnsureBaseOwnership(ctx.UserID, op.SourceBaseID); err != nil {
			return err
		}

		if err := op.Cancel(); err != nil {
			return cqrs.NewDomainError(err)
		}

		if err := oRepo.Update(op); err != nil {
			return err
		}
		if err := c.Outbox.Tx(tx).Save(op.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})
	return err
}

// SpeedUpOperationWithCrystals allows a user to spend crystals to fast-forward
// an in-flight military operation (outbound or returning) to its arrival.
func (c *OperationCommands) SpeedUpOperationWithCrystals(ctx cqrs.CommandContext, operationID int) error {
	err := c.TxMgr.WithTx(func(tx ports.Transaction) error {
		oRepo := c.OperationRepo.Tx(tx)
		uRepo := c.UserRepo.Tx(tx)

		op, err := oRepo.FindByIDForUpdate(operationID)
		if err != nil {
			return repoErr(err)
		}

		// Ensure the caller owns the source base for this operation.
		if err := c.Access.EnsureBaseOwnership(ctx.UserID, op.SourceBaseID); err != nil {
			return err
		}

		user, err := uRepo.FindByIDForUpdate(ctx.UserID)
		if err != nil {
			return repoErr(err)
		}

		if err := c.crystalService.SpeedUpOperation(user, op); err != nil {
			return cqrs.NewDomainError(err)
		}

		if err := uRepo.Update(user); err != nil {
			return err
		}
		if err := oRepo.Update(op); err != nil {
			return err
		}
		if err := c.Outbox.Tx(tx).Save(op.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (c *OperationCommands) HandleUpdateMilitaryOperationJob(cmd ports.UpdateMilitaryOperationJob) error {
	err := c.TxMgr.WithTx(func(tx ports.Transaction) error {
		oRepo := c.OperationRepo.Tx(tx)
		op, err := oRepo.FindByIDForUpdate(cmd.OperationID)
		if err != nil {
			return err
		}
		op.UpdatePhaseBasedOnTime()
		if err := oRepo.Update(op); err != nil {
			return err
		}
		if err := c.Outbox.Tx(tx).Save(op.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (c *OperationCommands) HandleMilitaryOperationStartedEvent(event domain.MilitaryOperationStartedEvent) error {
	return c.Scheduler.Schedule(ports.UpdateMilitaryOperationJob{OperationID: event.OperationID}, event.OutboundArriveAt)
}

func (c *OperationCommands) HandleMilitaryOperationArrivedEvent(event domain.MilitaryOperationArrivedEvent) error {
	err := c.TxMgr.WithTx(func(tx ports.Transaction) error {
		oRepo := c.OperationRepo.Tx(tx)
		sRepo := c.SectorRepo.Tx(tx)
		bRepo := c.UserBaseRepo.Tx(tx)
		rRepo := c.ResourceRepo.Tx(tx)
		dRepo := c.DangerousRepo.Tx(tx)
		srRepo := c.ScanReportRepo.Tx(tx)
		op, err := oRepo.FindByIDForUpdate(event.OperationID)
		if err != nil {
			return err
		}
		if op.Phase != domain.OperationPhaseAtTarget {
			return nil // Already handled or inconsistent state
		}
		sector, err := c.Provisioner.EnsureSectorExists(sRepo, op.TargetCoordinates.X, op.TargetCoordinates.Y)
		if err != nil {
			return err
		}
		attackerBase, err := bRepo.FindByIDForUpdate(op.SourceBaseID)
		if err != nil {
			return err
		}
		svc := domain.NewMilitaryOperationService(op, attackerBase)
		events, report, err := c.resolveOperationAtTarget(op, sector, svc, sRepo, bRepo, rRepo, dRepo)
		if err != nil {
			return err
		}
		if report != nil {
			report.SourceOperationID = op.ID
			if err := srRepo.Create(report); err == nil {
				report.EmitCreated()
			}
		}
		if err := bRepo.Update(attackerBase); err != nil {
			return err
		}
		if err := oRepo.Update(op); err != nil {
			return err
		}
		var allEvents []domain.DomainEvent
		allEvents = append(allEvents, op.EventProducer.PullEvents()...)
		allEvents = append(allEvents, events...)
		if report != nil {
			allEvents = append(allEvents, report.EventProducer.PullEvents()...)
		}
		if err := c.Outbox.Tx(tx).Save(allEvents); err != nil {
			return err
		}
		return nil
	})
	return err
}

// resolveOperationAtTarget encapsulates location-type-specific resolution and optional scan report creation.
// It mutates the target location through the provided repositories and returns an in-memory scan report (if any).
func (_ *OperationCommands) resolveOperationAtTarget(
	op *domain.MilitaryOperation,
	sector *domain.SectorModel,
	svc domain.MilitaryOperationService,
	sRepo ports.SectorRepository,
	bRepo ports.UserBaseRepository,
	rRepo ports.ResourceLocationRepository,
	dRepo ports.DangerousLocationRepository,
) ([]domain.DomainEvent, *domain.SectorScanReport, error) {
	occType, err := sRepo.GetLocationTypeByCoordinates(sector.Coordinates.X, sector.Coordinates.Y)
	if err != nil {
		return nil, nil, repoErr(err)
	}
	var report *domain.SectorScanReport
	var events []domain.DomainEvent
	switch occType {
	case domain.LocationTypeUserBase:
		base, err := bRepo.FindByCoordinatesForUpdate(op.TargetCoordinates.X, op.TargetCoordinates.Y)
		if err != nil {
			return nil, nil, err
		}
		svc.ResolveAgainstUserBase(base)
		if err := bRepo.Update(base); err != nil {
			return nil, nil, err
		}
		events = append(events, base.EventProducer.PullEvents()...)
		if op.Result == domain.OperationResultSuccess {
			if op.Type == domain.MilitaryOperationTypeSpy && op.SpyResult != nil && op.SpyResult.Outcome == domain.SpyOutcomeBlockedByCloaking {
				report = domain.NewSectorScanReportFromUserBase(op.SourceBaseID, sector, nil)
			} else {
				report = domain.NewSectorScanReportFromUserBase(op.SourceBaseID, sector, base)
			}
		}
	case domain.LocationTypeResourceful:
		res, err := rRepo.FindByCoordinatesForUpdate(op.TargetCoordinates.X, op.TargetCoordinates.Y)
		if err != nil {
			return nil, nil, err
		}
		svc.ResolveAgainstResourceLocation(res)
		if err := rRepo.Update(res); err != nil {
			return nil, nil, err
		}
		events = append(events, res.EventProducer.PullEvents()...)
		if op.Result == domain.OperationResultSuccess {
			if op.Type == domain.MilitaryOperationTypeSpy && op.SpyResult != nil && op.SpyResult.Outcome == domain.SpyOutcomeBlockedByCloaking {
				report = domain.NewSectorScanReportFromResourceLocation(op.SourceBaseID, sector, nil)
			} else {
				report = domain.NewSectorScanReportFromResourceLocation(op.SourceBaseID, sector, res)
			}
		}
	case domain.LocationTypeDangerous:
		dl, err := dRepo.FindByCoordinatesForUpdate(op.TargetCoordinates.X, op.TargetCoordinates.Y)
		if err != nil {
			return nil, nil, err
		}
		svc.ResolveAgainstDangerousLocation(dl)
		if err := dRepo.Update(dl); err != nil {
			return nil, nil, err
		}
		events = append(events, dl.EventProducer.PullEvents()...)
		if op.Result == domain.OperationResultSuccess {
			if op.Type == domain.MilitaryOperationTypeSpy && op.SpyResult != nil && op.SpyResult.Outcome == domain.SpyOutcomeBlockedByCloaking {
				report = domain.NewSectorScanReportFromDangerousLocation(op.SourceBaseID, sector, nil)
			} else {
				report = domain.NewSectorScanReportFromDangerousLocation(op.SourceBaseID, sector, dl)
			}
		}
	case domain.LocationTypeEmpty:
		svc.ResolveAgainstEmptySector(sector)
		report = domain.NewSectorScanReportFromEmptySector(op.SourceBaseID, sector)
	default:
		return nil, nil, fmt.Errorf("unsupported sector classification")
	}
	return events, report, nil
}

func (c *OperationCommands) HandleMilitaryOperationReturnStartedEvent(event domain.MilitaryOperationReturnStartedEvent) error {
	return c.Scheduler.Schedule(ports.UpdateMilitaryOperationJob{OperationID: event.OperationID}, event.ReturnArriveAt)
}

func (c *OperationCommands) HandleMilitaryOperationReturnArrivedEvent(event domain.MilitaryOperationReturnArrivedEvent) error {
	err := c.TxMgr.WithTx(func(tx ports.Transaction) error {
		oRepo := c.OperationRepo.Tx(tx)
		bRepo := c.UserBaseRepo.Tx(tx)
		op, err := oRepo.FindByID(event.OperationID)
		if err != nil {
			return err
		}
		base, err := bRepo.FindByIDForUpdate(op.SourceBaseID)
		if err != nil {
			return err
		}
		base.ReturnAllDeployedFromOperation(op.ID)
		if op.AttackResult != nil {
			base.CreditLoot(op.AttackResult.Loot)

			if len(op.AttackResult.Trophies) > 0 {
				protos, err := c.StorageProtos.Tx(tx).FindAllPrototypes()
				if err == nil {
					protoMap := make(map[int]domain.StorageItemPrototype, len(protos))
					for _, p := range protos {
						protoMap[p.ID] = *p
					}
					base.AddTrophies(op.AttackResult.Trophies, protoMap)
				}
			}
		}
		if err := bRepo.Update(base); err != nil {
			return err
		}
		if err := c.Outbox.Tx(tx).Save(base.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})
	return err
}
