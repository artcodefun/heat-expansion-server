package commands

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
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

func (c *ActivityCommands) HandleMilitaryOperationStartedEvent(event domain.MilitaryOperationStartedEvent) error {
	op, err := c.OperationRepo.FindByID(event.OperationID)
	if err != nil {
		return err
	}

	item := domain.NewActivityFromOffenseOperation(op.SourceBaseID, op)
	return c.TxMgr.WithTx(func(tx ports.Transaction) error {
		if err := c.ActivityRepo.Tx(tx).Create(&item); err != nil {
			return err
		}
		return c.OutboxEvents.Tx(tx).Save(item.PullEvents())
	})
}

func (c *ActivityCommands) HandleMilitaryOperationResolvedEvent(event domain.MilitaryOperationResolvedEvent) error {
	op, err := c.OperationRepo.FindByID(event.OperationID)
	if err != nil {
		return err
	}
	occType, _ := c.SectorRepo.GetLocationTypeByCoordinates(op.TargetCoordinates.X, op.TargetCoordinates.Y)
	if occType != domain.LocationTypeUserBase {
		return nil
	}
	base, err := c.UserBaseRepo.FindByCoordinates(op.TargetCoordinates.X, op.TargetCoordinates.Y)
	if err != nil {
		return err
	}
	item := domain.NewActivityFromDefenseOperation(base.ID, op)
	if ts := event.OccurredAt(); ts != 0 {
		item.CreatedAt = ts
	}
	return c.TxMgr.WithTx(func(tx ports.Transaction) error {
		if err := c.ActivityRepo.Tx(tx).Create(&item); err != nil {
			return err
		}
		return c.OutboxEvents.Tx(tx).Save(item.PullEvents())
	})
}

func (c *ActivityCommands) HandleScanReportCreatedEvent(event domain.ScanReportCreatedEvent) error {
	report, err := c.ScanRepo.FindByID(event.ReportID)
	if err != nil {
		return err
	}

	attackerActivity := domain.NewActivityFromScan(event.BaseID, report)

	// Defender side detection
	var defenderActivity *domain.ActivityItem
	occType, _ := c.SectorRepo.GetLocationTypeByCoordinates(report.Coordinates.X, report.Coordinates.Y)
	if occType == domain.LocationTypeUserBase {
		defenderBase, err := c.UserBaseRepo.FindByCoordinates(report.Coordinates.X, report.Coordinates.Y)
		attackerBase, _ := c.UserBaseRepo.FindByID(event.BaseID)

		if err == nil && defenderBase != nil && attackerBase != nil {
			interceptInfo := c.intelService.TriangulateScanSource(attackerBase.Coordinates, defenderBase, !report.IsCloaked)

			da := domain.NewActivityFromScanIntercept(defenderBase.ID, interceptInfo)
			defenderActivity = &da
		}
	}

	return c.TxMgr.WithTx(func(tx ports.Transaction) error {
		repo := c.ActivityRepo.Tx(tx)
		outbox := c.OutboxEvents.Tx(tx)

		if err := repo.Create(&attackerActivity); err != nil {
			return err
		}
		if err := outbox.Save(attackerActivity.PullEvents()); err != nil {
			return err
		}

		if defenderActivity != nil {
			if err := repo.Create(defenderActivity); err != nil {
				return err
			}
			if err := outbox.Save(defenderActivity.PullEvents()); err != nil {
				return err
			}
		}
		return nil
	})
}

func (c *ActivityCommands) HandleRadarThreatDetectedEvent(event domain.RadarThreatDetectedEvent) error {
	threat, err := c.RadarThreatRepo.FindByID(event.RadarThreatID)
	if err != nil {
		return err
	}

	activity := domain.NewActivityFromRadarThreat(threat)
	return c.TxMgr.WithTx(func(tx ports.Transaction) error {
		if err := c.ActivityRepo.Tx(tx).Create(&activity); err != nil {
			return err
		}
		return c.OutboxEvents.Tx(tx).Save(activity.PullEvents())
	})
}
