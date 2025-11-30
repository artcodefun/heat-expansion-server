package commands

import (
	"math/rand"
	"time"

	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/core/services"
)

type WorldGenerationCommands struct {
	UserBases            ports.UserBaseRepository
	Sectors              ports.SectorRepository
	ResourceLocations    ports.ResourceLocationRepository
	DangerousLocations   ports.DangerousLocationRepository
	Content              ports.ContentGenerator
	Provisioner          *services.SectorProvisioningService
	Scheduler            ports.Scheduler
	TxMgr                ports.TransactionManager
	SpawnRadius          int
	MaxResourcefulNearby int
	MaxDangerousNearby   int
	RespawnPeriodSeconds int64
	SpawnAttemptsPerJob  int
}

func NewWorldGenerationCommands(bases ports.UserBaseRepository, sectors ports.SectorRepository, res ports.ResourceLocationRepository, dang ports.DangerousLocationRepository, content ports.ContentGenerator, provisioner *services.SectorProvisioningService, scheduler ports.Scheduler, txMgr ports.TransactionManager) *WorldGenerationCommands {
	return &WorldGenerationCommands{UserBases: bases, Sectors: sectors, ResourceLocations: res, DangerousLocations: dang, Content: content, Provisioner: provisioner, Scheduler: scheduler, TxMgr: txMgr, SpawnRadius: 10, MaxResourcefulNearby: 12, MaxDangerousNearby: 6, RespawnPeriodSeconds: 3600, SpawnAttemptsPerJob: 20}
}

func (c *WorldGenerationCommands) HandleSpawnNearbyLocationsJob(job ports.SpawnNearbyLocationsJob) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bases, err := c.UserBases.FindAll()
	if err != nil || len(bases) == 0 {
		c.reschedule()
		return nil
	}
	base := bases[r.Intn(len(bases))]
	baseSector, err := c.Sectors.FindByCoordinates(base.Coordinates.X, base.Coordinates.Y)
	if err != nil || baseSector == nil {
		c.reschedule()
		return nil
	}
	sectors, err := c.Sectors.FindAll()
	if err != nil {
		c.reschedule()
		return nil
	}
	r2 := c.SpawnRadius * c.SpawnRadius
	resourceCount := 0
	dangerousCount := 0
	for _, s := range sectors {
		dx := s.Coordinates.X - baseSector.Coordinates.X
		dy := s.Coordinates.Y - baseSector.Coordinates.Y
		if dx*dx+dy*dy <= r2 {
			lt, _ := c.Sectors.GetLocationTypeByCoordinates(s.Coordinates.X, s.Coordinates.Y)
			switch lt {
			case domain.LocationTypeResourceful:
				resourceCount++
			case domain.LocationTypeDangerous:
				dangerousCount++
			}
		}
	}
	if resourceCount >= c.MaxResourcefulNearby && dangerousCount >= c.MaxDangerousNearby {
		c.reschedule()
		return nil
	}
	for i := 0; i < c.SpawnAttemptsPerJob; i++ {
		dx := r.Intn(2*c.SpawnRadius+1) - c.SpawnRadius
		dy := r.Intn(2*c.SpawnRadius+1) - c.SpawnRadius
		if dx*dx+dy*dy > r2 {
			continue
		}
		targetX := baseSector.Coordinates.X + dx
		targetY := baseSector.Coordinates.Y + dy
		tSector, _ := c.Sectors.FindByCoordinates(targetX, targetY)
		if tSector == nil {
			_ = c.TxMgr.WithTx(func(tx ports.Transaction) error {
				var err error
				tSector, err = c.Provisioner.EnsureSectorExists(c.Sectors.Tx(tx), targetX, targetY)
				return err
			})
		}
		lt, _ := c.Sectors.GetLocationTypeByCoordinates(targetX, targetY)
		if lt != domain.LocationTypeEmpty {
			continue
		}
		roll := r.Float64()
		if resourceCount < c.MaxResourcefulNearby && (roll < 0.7 || dangerousCount >= c.MaxDangerousNearby) {
			loc := &domain.ResourceLocationModel{Coordinates: tSector.Coordinates, Type: "GENERIC_NODE", Amount: 0, Resources: domain.LocationResourceStats{Credits: r.Intn(500) + 100, Iron: r.Intn(400) + 80, Titanium: r.Intn(300) + 60, Antimatter: r.Intn(200) + 40, CalculationTimestamp: domain.NowUnix()}}
			_ = c.persistResourceful(loc)
			break
		} else if dangerousCount < c.MaxDangerousNearby {
			loc := &domain.DangerousLocationModel{Coordinates: tSector.Coordinates, DangerLevel: 1, Resources: domain.LocationResourceStats{Credits: r.Intn(600) + 150, Iron: r.Intn(500) + 120, Titanium: r.Intn(400) + 90, Antimatter: r.Intn(300) + 70, CalculationTimestamp: domain.NowUnix()}}
			_ = c.persistDangerous(loc)
			break
		}
	}
	c.reschedule()
	return nil
}

func (c *WorldGenerationCommands) persistResourceful(loc *domain.ResourceLocationModel) error {
	return c.TxMgr.WithTx(func(tx ports.Transaction) error {
		sRepo := c.Sectors.Tx(tx)
		rRepo := c.ResourceLocations.Tx(tx)
		return c.Provisioner.CreateResourceLocationIfEmpty(sRepo, rRepo, loc)
	})
}

func (c *WorldGenerationCommands) persistDangerous(loc *domain.DangerousLocationModel) error {
	return c.TxMgr.WithTx(func(tx ports.Transaction) error {
		sRepo := c.Sectors.Tx(tx)
		dRepo := c.DangerousLocations.Tx(tx)
		return c.Provisioner.CreateDangerousLocationIfEmpty(sRepo, dRepo, loc)
	})
}

func (c *WorldGenerationCommands) reschedule() {
	jitter := int64(rand.Intn(300))
	_ = c.Scheduler.Schedule(ports.SpawnNearbyLocationsJob{}, time.Now().Unix()+c.RespawnPeriodSeconds+jitter)
}
