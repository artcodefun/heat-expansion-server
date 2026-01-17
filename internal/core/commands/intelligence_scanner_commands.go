package commands

import (
	"math/rand"
	"time"

	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/core/services"
)

type IntelligenceScannerCommands struct {
	BaseRepo          ports.UserBaseRepository
	SectorRepo        ports.SectorRepository
	ResourceRepo      ports.ResourceLocationRepository
	DangerousRepo     ports.DangerousLocationRepository
	ScanReportRepo    ports.ScanReportRepository
	SectorProvisioner *services.SectorProvisioningService
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
		Scheduler:         scheduler,
		Outbox:            outbox,
		TxMgr:             txMgr,
	}
}

func (c *IntelligenceScannerCommands) HandleBuildingProductionFinishedEvent(event *domain.BuildingProductionFinishedEvent) error {
	base, err := c.BaseRepo.FindByID(event.BaseID)
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
				_ = c.Scheduler.Schedule(ports.IntelligenceScanJob{BaseID: event.BaseID, BuildingID: b.ID}, firstAt)
			}
			break
		}
	}
	return nil
}

func (c *IntelligenceScannerCommands) HandleIntelligenceScanJob(job ports.IntelligenceScanJob) error {
	// 1. Idempotency: skip if already scanned recently
	now := domain.NowUnix()
	exists, err := c.ScanReportRepo.RecentReportExistsByScanner(job.BuildingID, now-60) // small buffer for overlapping runs
	if err != nil || exists {
		return err
	}

	base, err := c.BaseRepo.FindByID(job.BaseID)
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

	target := randomSectorInRange(base.Coordinates, rangeTiles)
	sector, err := c.SectorRepo.FindByCoordinates(target.X, target.Y)
	if err == ports.ErrNotFound || sector == nil {
		_ = c.TxMgr.WithTx(func(tx ports.Transaction) error {
			var inErr error
			sector, inErr = c.SectorProvisioner.EnsureSectorExists(c.SectorRepo.Tx(tx), target.X, target.Y)
			return inErr
		})
	} else if err != nil {
		c.reschedule(job, periodSec)
		return nil
	}

	occType, _ := c.SectorRepo.GetLocationTypeByCoordinates(sector.Coordinates.X, sector.Coordinates.Y)
	var report *domain.SectorScanReport
	switch occType {
	case domain.LocationTypeUserBase:
		defenderBase, _ := c.BaseRepo.FindByCoordinates(sector.Coordinates.X, sector.Coordinates.Y)
		report = domain.NewSectorScanReportFromUserBase(base.ID, sector, defenderBase)
	case domain.LocationTypeResourceful:
		res, _ := c.ResourceRepo.FindByCoordinates(sector.Coordinates.X, sector.Coordinates.Y)
		report = domain.NewSectorScanReportFromResourceLocation(base.ID, sector, res)
	case domain.LocationTypeDangerous:
		dl, _ := c.DangerousRepo.FindByCoordinates(sector.Coordinates.X, sector.Coordinates.Y)
		report = domain.NewSectorScanReportFromDangerousLocation(base.ID, sector, dl)
	case domain.LocationTypeEmpty:
		report = domain.NewSectorScanReportFromEmptySector(base.ID, sector)
	default:
		c.reschedule(job, periodSec)
		return nil
	}

	report.SourceScannerID = &job.BuildingID
	err = c.TxMgr.WithTx(func(tx ports.Transaction) error {
		srRepo := c.ScanReportRepo.Tx(tx)
		if err := srRepo.Create(report); err != nil {
			return err
		}
		report.EmitCreated()
		if err := c.Outbox.Tx(tx).Save(report.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})

	c.reschedule(job, periodSec)
	return err
}

func (c *IntelligenceScannerCommands) reschedule(job ports.IntelligenceScanJob, periodSec int64) {
	jitter := int64(rand.Intn(60) - 30)
	_ = c.Scheduler.Schedule(job, time.Now().Unix()+periodSec+jitter)
}

func randomSectorInRange(origin domain.Vector2i, r int) domain.Vector2i {
	if r <= 0 {
		return origin
	}
	dx := rand.Intn(r*2+1) - r
	dy := rand.Intn(r*2+1) - r
	return domain.Vector2i{X: origin.X + dx, Y: origin.Y + dy}
}
