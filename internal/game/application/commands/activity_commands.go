package commands

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
)

type ActivityCommands struct {
	ActivityRepo    ports.ActivityRepository
	OperationRepo   ports.MilitaryOperationRepository
	RadarThreatRepo ports.RadarThreatRepository
	SectorRepo      ports.SectorRepository
	UserBaseRepo    ports.UserBaseRepository
	ScanRepo        ports.ScanReportRepository
	intelService    *domain.IntelligenceService
	OutboxEvents    ports.OutboxEventRepository
	TxMgr           ports.TransactionManager
}

func NewActivityCommands(
	activityRepo ports.ActivityRepository,
	opRepo ports.MilitaryOperationRepository,
	radarThreatRepo ports.RadarThreatRepository,
	sectorRepo ports.SectorRepository,
	baseRepo ports.UserBaseRepository,
	scanRepo ports.ScanReportRepository,
	outboxEvents ports.OutboxEventRepository,
	txMgr ports.TransactionManager,
) *ActivityCommands {
	return &ActivityCommands{
		ActivityRepo:    activityRepo,
		OperationRepo:   opRepo,
		RadarThreatRepo: radarThreatRepo,
		SectorRepo:      sectorRepo,
		UserBaseRepo:    baseRepo,
		ScanRepo:        scanRepo,
		intelService:    domain.NewIntelligenceService(),
		OutboxEvents:    outboxEvents,
		TxMgr:           txMgr,
	}
}

func (c *ActivityCommands) HandleMilitaryOperationStartedEvent(ctx context.Context, event domain.MilitaryOperationStartedEvent) error {
	op, err := c.OperationRepo.FindByID(ctx, event.OperationID)
	if err != nil {
		return err
	}

	item := domain.NewActivityFromOffenseOperation(op.SourceBaseID, op)
	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		repo := c.ActivityRepo.Tx(tx)
		if ok, _ := repo.ExistsForOperation(ctx, op.SourceBaseID, string(domain.ActivityKindOffense), event.OperationID); ok {
			return nil
		}
		if err := repo.Create(ctx, &item); err != nil {
			return err
		}
		return c.OutboxEvents.Tx(tx).Save(ctx, item.PullEvents())
	})
}

func (c *ActivityCommands) HandleMilitaryOperationResolvedEvent(ctx context.Context, event domain.MilitaryOperationResolvedEvent) error {
	op, err := c.OperationRepo.FindByID(ctx, event.OperationID)
	if err != nil {
		return err
	}
	occType, _ := c.SectorRepo.GetLocationTypeByCoordinates(ctx, op.TargetCoordinates.X, op.TargetCoordinates.Y)
	if occType != domain.LocationTypeUserBase {
		return nil
	}
	base, err := c.UserBaseRepo.FindByCoordinates(ctx, op.TargetCoordinates.X, op.TargetCoordinates.Y)
	if err != nil {
		return err
	}
	item := domain.NewActivityFromDefenseOperation(base.UserID, base.ID, op)
	if ts := event.OccurredAt(); ts != 0 {
		item.CreatedAt = ts
	}
	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		repo := c.ActivityRepo.Tx(tx)
		if ok, _ := repo.ExistsForOperation(ctx, base.ID, string(domain.ActivityKindDefense), event.OperationID); ok {
			return nil
		}
		if err := repo.Create(ctx, &item); err != nil {
			return err
		}
		return c.OutboxEvents.Tx(tx).Save(ctx, item.PullEvents())
	})
}

func (c *ActivityCommands) HandleScanReportCreatedEvent(ctx context.Context, event domain.ScanReportCreatedEvent) error {
	report, err := c.ScanRepo.FindByID(ctx, event.ReportID)
	if err != nil {
		return err
	}

	attackerBase, err := c.UserBaseRepo.FindByID(ctx, event.BaseID)
	if err != nil {
		return err
	}
	attackerActivity := domain.NewActivityFromScan(attackerBase.UserID, event.BaseID, report)

	// Defender side detection
	var defenderActivity *domain.ActivityItem
	occType, _ := c.SectorRepo.GetLocationTypeByCoordinates(ctx, report.Coordinates.X, report.Coordinates.Y)
	if report.SourceType != domain.ScanReportSourceDiplomaticReveal && occType == domain.LocationTypeUserBase {
		defenderBase, err := c.UserBaseRepo.FindByCoordinates(ctx, report.Coordinates.X, report.Coordinates.Y)

		if err == nil && defenderBase != nil && attackerBase != nil {
			interceptInfo := c.intelService.TriangulateScanSource(attackerBase.Coordinates, defenderBase, !report.IsCloaked)

			da := domain.NewActivityFromScanIntercept(defenderBase.UserID, defenderBase.ID, interceptInfo)
			defenderActivity = &da
		}
	}

	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		repo := c.ActivityRepo.Tx(tx)
		outbox := c.OutboxEvents.Tx(tx)

		if ok, _ := repo.ExistsForScanReport(ctx, event.ReportID); ok {
			return nil
		}

		if err := repo.Create(ctx, &attackerActivity); err != nil {
			return err
		}
		if err := outbox.Save(ctx, attackerActivity.PullEvents()); err != nil {
			return err
		}

		if defenderActivity != nil {
			if err := repo.Create(ctx, defenderActivity); err != nil {
				return err
			}
			if err := outbox.Save(ctx, defenderActivity.PullEvents()); err != nil {
				return err
			}
		}
		return nil
	})
}

func (c *ActivityCommands) HandleRadarThreatDetectedEvent(ctx context.Context, event domain.RadarThreatDetectedEvent) error {
	threat, err := c.RadarThreatRepo.FindByID(ctx, event.RadarThreatID)
	if err != nil {
		return err
	}

	ownerID, err := c.UserBaseRepo.GetOwnerID(ctx, threat.OwnerBaseID)
	if err != nil {
		return err
	}
	activity := domain.NewActivityFromRadarThreat(ownerID, threat)
	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		repo := c.ActivityRepo.Tx(tx)
		if ok, _ := repo.ExistsForOperation(ctx, event.OwnerBaseID, string(domain.ActivityKindRadar), event.OperationID); ok {
			return nil
		}
		if err := repo.Create(ctx, &activity); err != nil {
			return err
		}
		return c.OutboxEvents.Tx(tx).Save(ctx, activity.PullEvents())
	})
}
