package commands

import (
	"errors"
	"fmt"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/core/services"
)

type OperationCommands struct {
	UserBaseRepo   ports.UserBaseRepository
	SectorRepo     ports.SectorRepository
	OperationRepo  ports.MilitaryOperationRepository
	ResourceRepo   ports.ResourceLocationRepository
	DangerousRepo  ports.DangerousLocationRepository
	ScanReportRepo ports.ScanReportRepository
	Provisioner    *services.SectorProvisioningService
	Scheduler      ports.Scheduler
	EventPublisher ports.EventPublisher
	TxMgr          ports.TransactionManager
	Access         *services.AccessControlService
}

func NewOperationCommands(userBaseRepo ports.UserBaseRepository, sectorRepo ports.SectorRepository, opRepo ports.MilitaryOperationRepository, resRepo ports.ResourceLocationRepository, dangerRepo ports.DangerousLocationRepository, scanRepo ports.ScanReportRepository, provisioner *services.SectorProvisioningService, scheduler ports.Scheduler, publisher ports.EventPublisher, txMgr ports.TransactionManager, access *services.AccessControlService) *OperationCommands {
	return &OperationCommands{UserBaseRepo: userBaseRepo, SectorRepo: sectorRepo, OperationRepo: opRepo, ResourceRepo: resRepo, DangerousRepo: dangerRepo, ScanReportRepo: scanRepo, Provisioner: provisioner, Scheduler: scheduler, EventPublisher: publisher, TxMgr: txMgr, Access: access}
}

func (c *OperationCommands) CreateMilitaryOperation(ctx cqrs.CommandContext, opType domain.MilitaryOperationType, sourceBaseID int, targetX int, targetY int, deployments []domain.ArmyDeploymentRequest) (*domain.MilitaryOperation, error) {
	if err := c.Access.EnsureBaseOwnership(ctx.UserID, sourceBaseID); err != nil {
		return nil, err
	}
	if sourceBaseID <= 0 {
		return nil, errors.New("invalid source or target")
	}
	var createdOp *domain.MilitaryOperation
	var events []domain.DomainEvent
	err := c.TxMgr.WithTx(func(tx ports.Transaction) error {
		bRepo := c.UserBaseRepo.Tx(tx)
		sRepo := c.SectorRepo.Tx(tx)
		oRepo := c.OperationRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(sourceBaseID)
		if err != nil {
			return err
		}
		if len(deployments) == 0 {
			return errors.New("no units provided for operation")
		}
		readyToDeploy, err := base.GetReadyToDeployArmy(deployments)
		if err != nil {
			return err
		}
		units := domain.OperationUnitsFromDeployed(readyToDeploy)
		if len(units) == 0 {
			return errors.New("no units provided for operation")
		}
		if opType == domain.MilitaryOperationTypeSpy {
			for _, u := range units {
				if u.Category != domain.ArmyCategorySpy {
					return errors.New("spy operations require only spy units")
				}
			}
		}
		sourceSector, err := c.Provisioner.EnsureSectorExists(sRepo, base.Coordinates.X, base.Coordinates.Y)
		if err != nil {
			return err
		}
		targetSector, err := c.Provisioner.EnsureSectorExists(sRepo, targetX, targetY)
		if err != nil {
			return err
		}
		switch opType {
		case domain.MilitaryOperationTypeAttack:
			createdOp = domain.NewAttackOperation(base.UserID, sourceBaseID, sourceSector.Coordinates, targetSector.Coordinates, units)
		case domain.MilitaryOperationTypeSpy:
			createdOp = domain.NewSpyOperation(base.UserID, sourceBaseID, sourceSector.Coordinates, targetSector.Coordinates, units)
		default:
			return errors.New("unsupported operation type")
		}
		if err := oRepo.Create(createdOp); err != nil {
			return err
		}
		for _, ready := range readyToDeploy {
			if _, err := base.AllocateArmyToOperation(domain.ArmyDeploymentRequest{PresentItemID: ready.PresentItemID, Count: ready.Count}, createdOp.ID); err != nil {
				return err
			}
		}
		createdOp.Start()
		if err := oRepo.Update(createdOp); err != nil {
			return err
		}
		events = append(events, createdOp.EventProducer.PullEvents()...)
		return nil
	})
	if err != nil {
		return nil, err
	}
	publishEvents(events, c.EventPublisher)
	return createdOp, nil
}

func (c *OperationCommands) HandleUpdateMilitaryOperationJob(cmd ports.UpdateMilitaryOperationJob) error {
	var events []domain.DomainEvent
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
		events = append(events, op.EventProducer.PullEvents()...)
		return nil
	})
	if err != nil {
		return err
	}
	publishEvents(events, c.EventPublisher)
	return nil
}

func (c *OperationCommands) HandleMilitaryOperationStartedEvent(event domain.MilitaryOperationStartedEvent) error {
	return c.Scheduler.Schedule(ports.UpdateMilitaryOperationJob{OperationID: event.OperationID}, event.OutboundArriveAt)
}

func (c *OperationCommands) HandleMilitaryOperationArrivedEvent(event domain.MilitaryOperationArrivedEvent) error {
	var events []domain.DomainEvent
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
		sector, err := c.Provisioner.EnsureSectorExists(sRepo, op.TargetCoordinates.X, op.TargetCoordinates.Y)
		if err != nil {
			return err
		}
		attackerBase, err := bRepo.FindByIDForUpdate(op.SourceBaseID)
		if err != nil {
			return err
		}
		svc := domain.NewMilitaryOperationService(op, attackerBase)
		var report *domain.SectorScanReport
		occType, _ := sRepo.GetLocationTypeByCoordinates(sector.Coordinates.X, sector.Coordinates.Y)
		switch occType {
		case domain.LocationTypeUserBase:
			base, err := bRepo.FindByCoordinatesForUpdate(op.TargetCoordinates.X, op.TargetCoordinates.Y)
			if err != nil {
				return err
			}
			svc.ResolveAgainstUserBase(base)
			if err := bRepo.Update(base); err != nil {
				return err
			}
			if op.Result == domain.OperationResultSuccess {
				if op.Type == domain.MilitaryOperationTypeSpy && op.SpyResult != nil && op.SpyResult.Outcome == domain.SpyOutcomeBlockedByCloaking {
					report = domain.NewSectorScanReportFromUserBase(op.SourceBaseID, op.TargetCoordinates, nil, domain.LocationDetails{})
				} else {
					report = domain.NewSectorScanReportFromUserBase(op.SourceBaseID, op.TargetCoordinates, base, domain.LocationDetails{})
				}
			}
		case domain.LocationTypeResourceful:
			res, err := rRepo.FindByCoordinatesForUpdate(op.TargetCoordinates.X, op.TargetCoordinates.Y)
			if err != nil {
				return err
			}
			svc.ResolveAgainstResourceLocation(res)
			if err := rRepo.Update(res); err != nil {
				return err
			}
			if op.Result == domain.OperationResultSuccess {
				if op.Type == domain.MilitaryOperationTypeSpy && op.SpyResult != nil && op.SpyResult.Outcome == domain.SpyOutcomeBlockedByCloaking {
					report = domain.NewSectorScanReportFromResourceLocation(op.SourceBaseID, op.TargetCoordinates, nil, domain.LocationDetails{})
				} else {
					report = domain.NewSectorScanReportFromResourceLocation(op.SourceBaseID, op.TargetCoordinates, res, domain.LocationDetails{})
				}
			}
		case domain.LocationTypeDangerous:
			dl, err := dRepo.FindByCoordinatesForUpdate(op.TargetCoordinates.X, op.TargetCoordinates.Y)
			if err != nil {
				return err
			}
			svc.ResolveAgainstDangerousLocation(dl)
			if err := dRepo.Update(dl); err != nil {
				return err
			}
			if op.Result == domain.OperationResultSuccess {
				if op.Type == domain.MilitaryOperationTypeSpy && op.SpyResult != nil && op.SpyResult.Outcome == domain.SpyOutcomeBlockedByCloaking {
					report = domain.NewSectorScanReportFromDangerousLocation(op.SourceBaseID, op.TargetCoordinates, nil, domain.LocationDetails{})
				} else {
					report = domain.NewSectorScanReportFromDangerousLocation(op.SourceBaseID, op.TargetCoordinates, dl, domain.LocationDetails{})
				}
			}
		case domain.LocationTypeEmpty:
			svc.ResolveAgainstEmptySector(sector)
			report = domain.NewSectorScanReportFromEmptySector(op.SourceBaseID, op.TargetCoordinates, sector)
		default:
			return fmt.Errorf("unsupported sector classification")
		}
		if report != nil {
			report.SourceOperationID = op.ID
			if err := srRepo.Create(report); err == nil {
				report.EmitCreated()
				events = append(events, report.EventProducer.PullEvents()...)
			}
		}
		if err := bRepo.Update(attackerBase); err != nil {
			return err
		}
		if err := oRepo.Update(op); err != nil {
			return err
		}
		events = append(events, op.EventProducer.PullEvents()...)
		return nil
	})
	if err != nil {
		return err
	}
	publishEvents(events, c.EventPublisher)
	return nil
}

func (c *OperationCommands) HandleMilitaryOperationReturnStartedEvent(event domain.MilitaryOperationReturnStartedEvent) error {
	return c.Scheduler.Schedule(ports.UpdateMilitaryOperationJob{OperationID: event.OperationID}, event.ReturnArriveAt)
}

func (c *OperationCommands) HandleMilitaryOperationReturnArrivedEvent(event domain.MilitaryOperationReturnArrivedEvent) error {
	var events []domain.DomainEvent
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
		}
		if err := bRepo.Update(base); err != nil {
			return err
		}
		events = append(events, base.EventProducer.PullEvents()...)
		return nil
	})
	if err != nil {
		return err
	}
	publishEvents(events, c.EventPublisher)
	return nil
}
