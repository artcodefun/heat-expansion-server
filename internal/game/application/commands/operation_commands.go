package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
)

type OperationCommands struct {
	UserBaseRepo   ports.UserBaseRepository
	UserRepo       ports.UserRepository
	DiplomacyRepo  ports.DiplomaticRelationshipRepository
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

func NewOperationCommands(userBaseRepo ports.UserBaseRepository, userRepo ports.UserRepository, diplomacyRepo ports.DiplomaticRelationshipRepository, sectorRepo ports.SectorRepository, opRepo ports.MilitaryOperationRepository, resRepo ports.ResourceLocationRepository, dangerRepo ports.DangerousLocationRepository, scanRepo ports.ScanReportRepository, storageProtos ports.StoragePrototypeRepository, provisioner *services.SectorProvisioningService, scheduler ports.Scheduler, outbox ports.OutboxEventRepository, txMgr ports.TransactionManager, access *services.AccessControlService) *OperationCommands {
	return &OperationCommands{
		UserBaseRepo:   userBaseRepo,
		UserRepo:       userRepo,
		DiplomacyRepo:  diplomacyRepo,
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

func (c *OperationCommands) CreateMilitaryOperation(ctx context.Context, actor cqrs.Actor, opType domain.MilitaryOperationType, sourceBaseID int, targetX int, targetY int, deployments []domain.ArmyDeploymentRequest) (*domain.MilitaryOperation, error) {
	if err := c.Access.EnsureBaseOwnership(ctx, actor.UserID, sourceBaseID); err != nil {
		return nil, err
	}
	var createdOp *domain.MilitaryOperation
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		bRepo := c.UserBaseRepo.Tx(tx)
		sRepo := c.SectorRepo.Tx(tx)
		oRepo := c.OperationRepo.Tx(tx)
		diplomacyRepo := c.DiplomacyRepo.Tx(tx)
		scanRepo := c.ScanReportRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(ctx, sourceBaseID)
		if err != nil {
			return repoErr(err)
		}
		readyToDeploy, err := base.GetReadyToDeployArmy(deployments)
		if err != nil {
			return err
		}

		snaps := base.ActiveStorageSnaps()
		units := domain.MilitaryUnitsFromDeployed(readyToDeploy)
		sourceSector, err := c.Provisioner.EnsureSectorExists(ctx, sRepo, base.Coordinates.X, base.Coordinates.Y)
		if err != nil {
			return err
		}
		targetSector, err := c.Provisioner.EnsureSectorExists(ctx, sRepo, targetX, targetY)
		if err != nil {
			return err
		}
		if opType == domain.MilitaryOperationTypeAttack {
			knownTargetBase, err := c.isKnownPlayerBaseTarget(ctx, sRepo, scanRepo, sourceBaseID, targetSector.Coordinates.X, targetSector.Coordinates.Y)
			if err != nil {
				return err
			}
			if knownTargetBase {
				defenderBase, err := bRepo.FindByCoordinates(ctx, targetSector.Coordinates.X, targetSector.Coordinates.Y)
				if err != nil {
					return repoErr(err)
				}
				if defenderBase.UserID != base.UserID {
					rel, err := diplomacyRepo.FindBetweenUsers(ctx, base.UserID, defenderBase.UserID)
					if err != nil {
						if !errors.Is(err, ports.ErrNotFound) {
							return repoErr(err)
						}
						rel, err = domain.NewUnknownRelationship(base.UserID, defenderBase.UserID)
						if err != nil {
							return err
						}
					}
					if err := rel.CanPerformAttackOperation(); err != nil {
						return err
					}
				}
			}
		}
		var opCreationErr error
		switch opType {
		case domain.MilitaryOperationTypeAttack:
			createdOp, opCreationErr = domain.NewAttackOperation(base.UserID, sourceBaseID, sourceSector.Coordinates, targetSector.Coordinates, units, snaps)
		case domain.MilitaryOperationTypeSpy:
			createdOp, opCreationErr = domain.NewSpyOperation(base.UserID, sourceBaseID, sourceSector.Coordinates, targetSector.Coordinates, units, snaps)
		default:
			return cqrs.NewAppError(cqrs.KindInvalidInput, "error.application.invalid_operation_type")
		}
		if opCreationErr != nil {
			return opCreationErr
		}
		if err := oRepo.Create(ctx, createdOp); err != nil {
			return err
		}
		for _, ready := range readyToDeploy {
			if _, err := base.AllocateArmyToOperation(domain.ArmyDeploymentRequest{PresentItemID: ready.PresentItemID, Count: ready.Count}, createdOp.ID); err != nil {
				return err
			}
		}
		if err := bRepo.Update(ctx, base); err != nil {
			return err
		}
		createdOp.Start()
		if err := oRepo.Update(ctx, createdOp); err != nil {
			return err
		}
		if err := c.Outbox.Tx(tx).Save(ctx, createdOp.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return createdOp, nil
}

func (c *OperationCommands) CancelMilitaryOperation(ctx context.Context, actor cqrs.Actor, operationID int) error {
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		oRepo := c.OperationRepo.Tx(tx)
		op, err := oRepo.FindByIDForUpdate(ctx, operationID)
		if err != nil {
			return repoErr(err)
		}

		if err := c.Access.EnsureBaseOwnership(ctx, actor.UserID, op.SourceBaseID); err != nil {
			return err
		}

		if err := op.Cancel(); err != nil {
			return err
		}

		if err := oRepo.Update(ctx, op); err != nil {
			return err
		}
		if err := c.Outbox.Tx(tx).Save(ctx, op.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})
	return err
}

// SpeedUpOperationWithCrystals allows a user to spend crystals to fast-forward
// an in-flight military operation (outbound or returning) to its arrival.
func (c *OperationCommands) SpeedUpOperationWithCrystals(ctx context.Context, actor cqrs.Actor, operationID int) error {
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		oRepo := c.OperationRepo.Tx(tx)
		uRepo := c.UserRepo.Tx(tx)

		op, err := oRepo.FindByIDForUpdate(ctx, operationID)
		if err != nil {
			return repoErr(err)
		}

		// Ensure the caller owns the source base for this operation.
		if err := c.Access.EnsureBaseOwnership(ctx, actor.UserID, op.SourceBaseID); err != nil {
			return err
		}

		user, err := uRepo.FindByIDForUpdate(ctx, actor.UserID)
		if err != nil {
			return repoErr(err)
		}

		if err := c.crystalService.SpeedUpOperation(user, op); err != nil {
			return err
		}

		if err := uRepo.Update(ctx, user); err != nil {
			return err
		}
		if err := oRepo.Update(ctx, op); err != nil {
			return err
		}
		if err := c.Outbox.Tx(tx).Save(ctx, op.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (_ *OperationCommands) isKnownPlayerBaseTarget(ctx context.Context, sectorRepo ports.SectorRepository, scanRepo ports.ScanReportRepository, sourceBaseID, targetX, targetY int) (bool, error) {
	occType, err := sectorRepo.GetLocationTypeByCoordinates(ctx, targetX, targetY)
	if err != nil {
		return false, err
	}
	if occType != domain.LocationTypeUserBase {
		return false, nil
	}
	reports, err := scanRepo.FindByBaseAndCoordinates(ctx, sourceBaseID, targetX, targetY)
	if err != nil {
		return false, err
	}
	return len(reports) > 0, nil
}

func (c *OperationCommands) HandleUpdateMilitaryOperationJob(ctx context.Context, cmd ports.UpdateMilitaryOperationJob) error {
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		oRepo := c.OperationRepo.Tx(tx)
		op, err := oRepo.FindByIDForUpdate(ctx, cmd.OperationID)
		if err != nil {
			return err
		}
		op.UpdatePhaseBasedOnTime()
		if err := oRepo.Update(ctx, op); err != nil {
			return err
		}
		if err := c.Outbox.Tx(tx).Save(ctx, op.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (c *OperationCommands) HandleMilitaryOperationStartedEvent(ctx context.Context, event domain.MilitaryOperationStartedEvent) error {
	return c.Scheduler.Schedule(ctx, ports.UpdateMilitaryOperationJob{OperationID: event.OperationID}, event.OutboundArriveAt)
}

func (c *OperationCommands) HandleMilitaryOperationArrivedEvent(ctx context.Context, event domain.MilitaryOperationArrivedEvent) error {
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		oRepo := c.OperationRepo.Tx(tx)
		sRepo := c.SectorRepo.Tx(tx)
		bRepo := c.UserBaseRepo.Tx(tx)
		rRepo := c.ResourceRepo.Tx(tx)
		dRepo := c.DangerousRepo.Tx(tx)
		srRepo := c.ScanReportRepo.Tx(tx)
		op, err := oRepo.FindByIDForUpdate(ctx, event.OperationID)
		if err != nil {
			return err
		}
		if op.Phase != domain.OperationPhaseAtTarget {
			return nil // Already handled or inconsistent state
		}
		sector, err := c.Provisioner.EnsureSectorExists(ctx, sRepo, op.TargetCoordinates.X, op.TargetCoordinates.Y)
		if err != nil {
			return err
		}
		attackerBase, err := bRepo.FindByIDForUpdate(ctx, op.SourceBaseID)
		if err != nil {
			return err
		}
		svc := domain.NewMilitaryOperationService(op, attackerBase)
		events, report, err := c.resolveOperationAtTarget(ctx, op, sector, svc, sRepo, bRepo, rRepo, dRepo)
		if err != nil {
			return err
		}
		if report != nil {
			report.SourceOperationID = op.ID
			if err := srRepo.Create(ctx, report); err == nil {
				report.EmitCreated()
			}
		}
		if err := bRepo.Update(ctx, attackerBase); err != nil {
			return err
		}
		if err := oRepo.Update(ctx, op); err != nil {
			return err
		}
		var allEvents []domain.DomainEvent
		allEvents = append(allEvents, op.EventProducer.PullEvents()...)
		allEvents = append(allEvents, events...)
		if report != nil {
			allEvents = append(allEvents, report.EventProducer.PullEvents()...)
		}
		if err := c.Outbox.Tx(tx).Save(ctx, allEvents); err != nil {
			return err
		}
		return nil
	})
	return err
}

// resolveOperationAtTarget encapsulates location-type-specific resolution and optional scan report creation.
// It mutates the target location through the provided repositories and returns an in-memory scan report (if any).
func (_ *OperationCommands) resolveOperationAtTarget(
	ctx context.Context,
	op *domain.MilitaryOperation,
	sector *domain.SectorModel,
	svc domain.MilitaryOperationService,
	sRepo ports.SectorRepository,
	bRepo ports.UserBaseRepository,
	rRepo ports.ResourceLocationRepository,
	dRepo ports.DangerousLocationRepository,
) ([]domain.DomainEvent, *domain.SectorScanReport, error) {
	occType, err := sRepo.GetLocationTypeByCoordinates(ctx, sector.Coordinates.X, sector.Coordinates.Y)
	if err != nil {
		return nil, nil, repoErr(err)
	}
	var report *domain.SectorScanReport
	var events []domain.DomainEvent
	switch occType {
	case domain.LocationTypeUserBase:
		base, err := bRepo.FindByCoordinatesForUpdate(ctx, op.TargetCoordinates.X, op.TargetCoordinates.Y)
		if err != nil {
			return nil, nil, err
		}
		svc.ResolveAgainstUserBase(base)
		if err := bRepo.Update(ctx, base); err != nil {
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
		res, err := rRepo.FindByCoordinatesForUpdate(ctx, op.TargetCoordinates.X, op.TargetCoordinates.Y)
		if err != nil {
			return nil, nil, err
		}
		svc.ResolveAgainstResourceLocation(res)
		if err := rRepo.Update(ctx, res); err != nil {
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
		dl, err := dRepo.FindByCoordinatesForUpdate(ctx, op.TargetCoordinates.X, op.TargetCoordinates.Y)
		if err != nil {
			return nil, nil, err
		}
		svc.ResolveAgainstDangerousLocation(dl)
		if err := dRepo.Update(ctx, dl); err != nil {
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

func (c *OperationCommands) HandleMilitaryOperationReturnStartedEvent(ctx context.Context, event domain.MilitaryOperationReturnStartedEvent) error {
	return c.Scheduler.Schedule(ctx, ports.UpdateMilitaryOperationJob{OperationID: event.OperationID}, event.ReturnArriveAt)
}

func (c *OperationCommands) HandleMilitaryOperationReturnArrivedEvent(ctx context.Context, event domain.MilitaryOperationReturnArrivedEvent) error {
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		oRepo := c.OperationRepo.Tx(tx)
		bRepo := c.UserBaseRepo.Tx(tx)
		op, err := oRepo.FindByID(ctx, event.OperationID)
		if err != nil {
			return err
		}
		base, err := bRepo.FindByIDForUpdate(ctx, op.SourceBaseID)
		if err != nil {
			return err
		}
		base.ReturnAllDeployedFromOperation(op.ID)
		if op.AttackResult != nil {
			base.CreditLoot(op.AttackResult.Loot)

			if len(op.AttackResult.Trophies) > 0 {
				protos, err := c.StorageProtos.Tx(tx).FindAllPrototypes(ctx)
				if err == nil {
					protoMap := make(map[int]domain.StorageItemPrototype, len(protos))
					for _, p := range protos {
						protoMap[p.ID] = *p
					}
					base.AddTrophies(op.AttackResult.Trophies, protoMap)
				}
			}
		}
		if err := bRepo.Update(ctx, base); err != nil {
			return err
		}
		if err := c.Outbox.Tx(tx).Save(ctx, base.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})
	return err
}
