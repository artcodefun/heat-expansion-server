package commands

import (
	"math/rand"
	"time"

	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/core/services"
)

type RadarCommands struct {
	BaseRepo          ports.UserBaseRepository
	SectorRepo        ports.SectorRepository
	ResourceRepo      ports.ResourceLocationRepository
	DangerousRepo     ports.DangerousLocationRepository
	ScanReportRepo    ports.ScanReportRepository
	SectorProvisioner *services.SectorProvisioningService
	Scheduler         ports.Scheduler
	EventPublisher    ports.EventPublisher
	TxMgr             ports.TransactionManager
}

func NewRadarCommands(baseRepo ports.UserBaseRepository, sectorRepo ports.SectorRepository, resRepo ports.ResourceLocationRepository, dangerRepo ports.DangerousLocationRepository, scanRepo ports.ScanReportRepository, provisioner *services.SectorProvisioningService, scheduler ports.Scheduler, publisher ports.EventPublisher, txMgr ports.TransactionManager) *RadarCommands {
	return &RadarCommands{BaseRepo: baseRepo, SectorRepo: sectorRepo, ResourceRepo: resRepo, DangerousRepo: dangerRepo, ScanReportRepo: scanRepo, SectorProvisioner: provisioner, Scheduler: scheduler, EventPublisher: publisher, TxMgr: txMgr}
}

func (c *RadarCommands) HandleRadarScanJob(job ports.RadarScanJob) error {
	base, err := c.BaseRepo.FindByID(job.BaseID)
	if err != nil {
		return err
	}
	var radarProto *domain.BuildItemPrototype
	for _, b := range base.BuildingsPresent {
		if b.ID == job.BuildingID {
			radarProto = &b.Prototype
			break
		}
	}
	if radarProto == nil || radarProto.IntelligenceData == nil || radarProto.IntelligenceData.Subtype != domain.IntelligenceSubtypeRadar {
		return nil
	}
	rangeTiles := radarProto.IntelligenceData.ScanRange
	if rangeTiles <= 0 {
		rangeTiles = 1
	}
	periodSec := radarProto.IntelligenceData.ScanCooldown
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
		report = domain.NewSectorScanReportFromUserBase(base.ID, sector.Coordinates, defenderBase, domain.LocationDetails{})
	case domain.LocationTypeResourceful:
		res, _ := c.ResourceRepo.FindByCoordinates(sector.Coordinates.X, sector.Coordinates.Y)
		report = domain.NewSectorScanReportFromResourceLocation(base.ID, sector.Coordinates, res, domain.LocationDetails{})
	case domain.LocationTypeDangerous:
		dl, _ := c.DangerousRepo.FindByCoordinates(sector.Coordinates.X, sector.Coordinates.Y)
		report = domain.NewSectorScanReportFromDangerousLocation(base.ID, sector.Coordinates, dl, domain.LocationDetails{})
	case domain.LocationTypeEmpty:
		report = domain.NewSectorScanReportFromEmptySector(base.ID, sector.Coordinates, sector)
	default:
		c.reschedule(job, periodSec)
		return nil
	}
	report.SourceOperationID = 0
	var events []domain.DomainEvent
	err = c.TxMgr.WithTx(func(tx ports.Transaction) error {
		srRepo := c.ScanReportRepo.Tx(tx)
		if err := srRepo.Create(report); err != nil {
			return nil
		}
		report.EmitCreated()
		events = append(events, report.EventProducer.PullEvents()...)
		return nil
	})
	if err == nil {
		publishEvents(events, c.EventPublisher)
	}
	c.reschedule(job, periodSec)
	return nil
}

func (c *RadarCommands) reschedule(job ports.RadarScanJob, periodSec int64) {
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
