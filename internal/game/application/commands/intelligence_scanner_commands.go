package commands

import (
	"context"
	"math/rand"
	"time"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
)

type IntelligenceScannerCommands struct {
	BaseRepo          ports.UserBaseRepository
	SectorRepo        ports.SectorRepository
	ResourceRepo      ports.ResourceLocationRepository
	DangerousRepo     ports.DangerousLocationRepository
	ScanReportRepo    ports.ScanReportRepository
	SectorProvisioner *services.SectorProvisioningService
	intelService      *domain.IntelligenceService
	Scheduler         ports.Scheduler
	Outbox            ports.OutboxEventRepository
	TxMgr             ports.TransactionManager
}

func NewIntelligenceScannerCommands(baseRepo ports.UserBaseRepository, sectorRepo ports.SectorRepository, resRepo ports.ResourceLocationRepository, dangerRepo ports.DangerousLocationRepository, scanRepo ports.ScanReportRepository, provisioner *services.SectorProvisioningService, scheduler ports.Scheduler, outbox ports.OutboxEventRepository, txMgr ports.TransactionManager) *IntelligenceScannerCommands {
	return &IntelligenceScannerCommands{
		BaseRepo:          baseRepo,
		SectorRepo:        sectorRepo,
		ResourceRepo:      resRepo,
		DangerousRepo:     dangerRepo,
		ScanReportRepo:    scanRepo,
		SectorProvisioner: provisioner,
		intelService:      domain.NewIntelligenceService(),
		Scheduler:         scheduler,
		Outbox:            outbox,
		TxMgr:             txMgr,
	}
}

func (c *IntelligenceScannerCommands) HandleBuildingProductionFinishedEvent(ctx context.Context, event domain.BuildingProductionFinishedEvent) error {
	base, err := c.BaseRepo.FindByID(ctx, event.BaseID)
	if err != nil {
		return nil
	}
	for _, b := range base.BuildingsPresent {
		if b.ID == event.PresentItemID {
			if b.Prototype.IntelligenceData != nil && b.Prototype.IntelligenceData.Subtype == domain.IntelligenceSubtypeScanner {
				cooldown := b.Prototype.IntelligenceData.ScanCooldown
				if cooldown <= 0 {
					cooldown = 3600
				}
				firstAt := time.Now().Unix() + cooldown
				_ = c.Scheduler.Schedule(ctx, ports.IntelligenceScanJob{BaseID: event.BaseID, BuildingID: b.ID}, firstAt)
			}
			break
		}
	}
	return nil
}

func (c *IntelligenceScannerCommands) HandleIntelligenceScanJob(ctx context.Context, job ports.IntelligenceScanJob) error {
	// 1. Idempotency: skip if already scanned recently
	now := domain.NowUnix()
	exists, err := c.ScanReportRepo.RecentReportExistsByScanner(ctx, job.BuildingID, now-60) // small buffer for overlapping runs
	if err != nil || exists {
		return err
	}

	base, err := c.BaseRepo.FindByID(ctx, job.BaseID)
	if err != nil {
		return err
	}
	var scannerProto *domain.BuildItemPrototype
	for _, b := range base.BuildingsPresent {
		if b.ID == job.BuildingID {
			scannerProto = &b.Prototype
			break
		}
	}
	if scannerProto == nil || scannerProto.IntelligenceData == nil || scannerProto.IntelligenceData.Subtype != domain.IntelligenceSubtypeScanner {
		return nil
	}

	rangeTiles := scannerProto.IntelligenceData.ScanRange
	if rangeTiles <= 0 {
		rangeTiles = 1
	}
	periodSec := scannerProto.IntelligenceData.ScanCooldown
	if periodSec <= 0 {
		periodSec = 3600
	}
	scanStrength := scannerProto.IntelligenceData.StealthStrength
	if scanStrength <= 0 {
		scanStrength = 100
	}

	target := randomSectorInRange(base.Coordinates, rangeTiles)
	sector, err := c.SectorRepo.FindByCoordinates(ctx, target.X, target.Y)
	if err == ports.ErrNotFound || sector == nil {
		_ = c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
			var inErr error
			sector, inErr = c.SectorProvisioner.EnsureSectorExists(ctx, c.SectorRepo.Tx(tx), target.X, target.Y)
			return inErr
		})
	} else if err != nil {
		c.reschedule(ctx, job, periodSec)
		return nil
	}

	occType, _ := c.SectorRepo.GetLocationTypeByCoordinates(ctx, sector.Coordinates.X, sector.Coordinates.Y)
	var report *domain.SectorScanReport

	switch occType {
	case domain.LocationTypeUserBase:
		defenderBase, _ := c.BaseRepo.FindByCoordinates(ctx, sector.Coordinates.X, sector.Coordinates.Y)
		if c.intelService.ResolveScanVisibility(scanStrength, defenderBase) {
			report = domain.NewSectorScanReportFromUserBase(base.ID, sector, defenderBase)
		} else {
			report = domain.NewSectorScanReportFromUserBase(base.ID, sector, nil)
		}
	case domain.LocationTypeResourceful:
		res, _ := c.ResourceRepo.FindByCoordinates(ctx, sector.Coordinates.X, sector.Coordinates.Y)
		report = domain.NewSectorScanReportFromResourceLocation(base.ID, sector, res)
	case domain.LocationTypeDangerous:
		dl, _ := c.DangerousRepo.FindByCoordinates(ctx, sector.Coordinates.X, sector.Coordinates.Y)
		report = domain.NewSectorScanReportFromDangerousLocation(base.ID, sector, dl)
	case domain.LocationTypeEmpty:
		report = domain.NewSectorScanReportFromEmptySector(base.ID, sector)
	default:
		c.reschedule(ctx, job, periodSec)
		return nil
	}

	report.SourceType = domain.ScanReportSourceScanner
	report.SourceID = &job.BuildingID
	err = c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		srRepo := c.ScanReportRepo.Tx(tx)
		if err := srRepo.Create(ctx, report); err != nil {
			return err
		}
		report.EmitCreated()
		if err := c.Outbox.Tx(tx).Save(ctx, report.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})

	c.reschedule(ctx, job, periodSec)
	return err
}

func (c *IntelligenceScannerCommands) reschedule(ctx context.Context, job ports.IntelligenceScanJob, periodSec int64) {
	jitter := int64(rand.Intn(60) - 30)
	_ = c.Scheduler.Schedule(ctx, job, time.Now().Unix()+periodSec+jitter)
}

func randomSectorInRange(origin domain.Vector2i, r int) domain.Vector2i {
	if r <= 0 {
		return origin
	}

	for {
		dx := rand.Intn(r*2+1) - r
		dy := rand.Intn(r*2+1) - r
		if dx == 0 && dy == 0 {
			continue
		}
		target := domain.Vector2i{X: origin.X + dx, Y: origin.Y + dy}
		if origin.DistanceTo(target) <= float64(r) {
			return target
		}
	}
}
