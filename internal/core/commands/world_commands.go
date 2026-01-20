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
	StoragePrototypes    ports.StoragePrototypeRepository
	ArmyPrototypes       ports.ArmyPrototypeRepository
	BuildPrototypes      ports.BuildPrototypeRepository
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

func NewWorldGenerationCommands(
	bases ports.UserBaseRepository,
	sectors ports.SectorRepository,
	res ports.ResourceLocationRepository,
	dang ports.DangerousLocationRepository,
	storage ports.StoragePrototypeRepository,
	army ports.ArmyPrototypeRepository,
	build ports.BuildPrototypeRepository,
	content ports.ContentGenerator,
	provisioner *services.SectorProvisioningService,
	scheduler ports.Scheduler,
	txMgr ports.TransactionManager,
) *WorldGenerationCommands {
	return &WorldGenerationCommands{
		UserBases:            bases,
		Sectors:              sectors,
		ResourceLocations:    res,
		DangerousLocations:   dang,
		StoragePrototypes:    storage,
		ArmyPrototypes:       army,
		BuildPrototypes:      build,
		Content:              content,
		Provisioner:          provisioner,
		Scheduler:            scheduler,
		TxMgr:                txMgr,
		SpawnRadius:          10,
		MaxResourcefulNearby: 12,
		MaxDangerousNearby:   6,
		RespawnPeriodSeconds: 3600,
		SpawnAttemptsPerJob:  20,
	}
}

func (c *WorldGenerationCommands) HandleSpawnNearbyLocationsJob(job ports.SpawnNearbyLocationsJob) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	if job.BaseID == 0 {
		return nil
	}

	base, err := c.UserBases.FindByID(job.BaseID)
	if err != nil {
		// If base is gone, stop rescheduling
		return nil
	}

	center := base.Coordinates

	// 2. Count locations in range using SQL
	resourceCount, dangerousCount, err := c.Sectors.CountLocationsInRange(center.X, center.Y, c.SpawnRadius)
	if err != nil {
		c.reschedule(job.BaseID)
		return nil
	}

	if resourceCount >= c.MaxResourcefulNearby && dangerousCount >= c.MaxDangerousNearby {
		c.reschedule(job.BaseID)
		return nil
	}

	// Development scaling: Use MaxSpace as a proxy for base level (default is 100)
	scaleFactor := float64(base.Stats.MaxSpace) / float64(domain.DefaultMaxSpace)
	if scaleFactor < 1.0 {
		scaleFactor = 1.0
	}

	storageProtos, _ := c.StoragePrototypes.FindAllPrototypes()
	armyProtos, _ := c.ArmyPrototypes.FindAllPrototypes()
	buildProtos, _ := c.BuildPrototypes.FindAllPrototypes()

	r2 := c.SpawnRadius * c.SpawnRadius
	for i := 0; i < c.SpawnAttemptsPerJob; i++ {
		dx := r.Intn(2*c.SpawnRadius+1) - c.SpawnRadius
		dy := r.Intn(2*c.SpawnRadius+1) - c.SpawnRadius
		if dx*dx+dy*dy > r2 {
			continue
		}
		targetX := center.X + dx
		targetY := center.Y + dy
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
			resTypes := []domain.ResourceType{domain.ResourceTypeCredits, domain.ResourceTypeIron, domain.ResourceTypeTitanium, domain.ResourceTypeAntimatter}
			resType := resTypes[r.Intn(len(resTypes))]

			// Assign faction associated with the resource
			faction := domain.FactionMarauders
			switch resType {
			case domain.ResourceTypeIron:
				faction = domain.FactionFerrousSwarm
			case domain.ResourceTypeTitanium:
				faction = domain.FactionTitanArachnids
			case domain.ResourceTypeAntimatter:
				faction = domain.FactionVoidEcho
			}

			worth := int(float64(500+r.Intn(1000)) * scaleFactor)
			loc := domain.NewResourceLocation(
				tSector.Coordinates,
				resType,
				faction,
				worth,
				armyProtos,
				buildProtos,
			)

			_ = c.persistResourceful(loc)
			break
		} else if dangerousCount < c.MaxDangerousNearby {
			// Pick a random dangerous NPC faction
			dangFactions := []domain.Faction{
				domain.FactionCustodianProtocol,
				domain.FactionScorchWalkers,
				domain.FactionObsidianSentinels,
				domain.FactionNeuralWormApex,
			}
			faction := dangFactions[r.Intn(len(dangFactions))]

			worth := int(float64(1000+r.Intn(2000)) * scaleFactor)
			loc := domain.NewDangerousLocation(
				tSector.Coordinates,
				faction,
				worth,
				storageProtos,
				armyProtos,
				buildProtos,
			)

			_ = c.persistDangerous(loc)
			break
		}
	}
	c.reschedule(job.BaseID)
	return nil
}

func (c *WorldGenerationCommands) HandleUserBaseCreatedEvent(ev domain.UserBaseCreatedEvent) error {
	jitter := int64(rand.Intn(60))
	return c.Scheduler.Schedule(ports.SpawnNearbyLocationsJob{BaseID: ev.BaseID}, time.Now().Unix()+jitter)
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

func (c *WorldGenerationCommands) HandleLocationDrainedEvent(event domain.LocationDrainedEvent) error {
	return c.TxMgr.WithTx(func(tx ports.Transaction) error {
		// 1. Delete the location record
		var err error
		if event.Type == domain.LocationTypeResourceful {
			err = c.ResourceLocations.Tx(tx).DeleteByCoordinates(event.X, event.Y)
		} else {
			err = c.DangerousLocations.Tx(tx).DeleteByCoordinates(event.X, event.Y)
		}
		if err != nil {
			return err
		}

		// Note: We don't need to update the sector table because location presence
		// is checked via EXISTS on resource/dangerous tables.
		return nil
	})
}

func (c *WorldGenerationCommands) reschedule(baseID int) {
	jitter := int64(rand.Intn(300))
	_ = c.Scheduler.Schedule(ports.SpawnNearbyLocationsJob{BaseID: baseID}, time.Now().Unix()+c.RespawnPeriodSeconds+jitter)
}
