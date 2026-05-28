package commands

import (
	"context"
	"math/rand"
	"time"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
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
		SpawnRadius:          7,
		MaxResourcefulNearby: 10,
		MaxDangerousNearby:   3,
		RespawnPeriodSeconds: 3600,
		SpawnAttemptsPerJob:  20,
	}
}

func (c *WorldGenerationCommands) HandleSpawnNearbyLocationsJob(ctx context.Context, job ports.SpawnNearbyLocationsJob) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	if job.BaseID == 0 {
		return nil
	}

	base, err := c.UserBases.FindByID(ctx, job.BaseID)
	if err != nil {
		// If base is gone, stop rescheduling
		return nil
	}

	center := base.Coordinates

	// 2. Count locations in range using SQL
	resourceCount, dangerousCount, err := c.Sectors.CountLocationsInRange(ctx, center.X, center.Y, c.SpawnRadius)
	if err != nil {
		c.reschedule(ctx, job.BaseID)
		return nil
	}

	if resourceCount >= c.MaxResourcefulNearby && dangerousCount >= c.MaxDangerousNearby {
		c.reschedule(ctx, job.BaseID)
		return nil
	}

	storageProtos, err := c.StoragePrototypes.FindAllPrototypes(ctx)
	if err != nil {
		c.reschedule(ctx, job.BaseID)
		return nil
	}
	storageProtos = domain.FilterStorageItemPrototypesByCreationSource(storageProtos, domain.CreationSourceNPCLocation)
	armyProtos, err := c.ArmyPrototypes.FindAllPrototypes(ctx)
	if err != nil {
		c.reschedule(ctx, job.BaseID)
		return nil
	}
	armyProtos = domain.FilterArmyItemPrototypesByCreationSource(armyProtos, domain.CreationSourceNPCLocation)
	buildProtos, err := c.BuildPrototypes.FindAllPrototypes(ctx)
	if err != nil {
		c.reschedule(ctx, job.BaseID)
		return nil
	}
	buildProtos = domain.FilterBuildItemPrototypesByCreationSource(buildProtos, domain.CreationSourceNPCLocation)

	r2 := c.SpawnRadius * c.SpawnRadius
	for i := 0; i < c.SpawnAttemptsPerJob; i++ {
		dx := r.Intn(2*c.SpawnRadius+1) - c.SpawnRadius
		dy := r.Intn(2*c.SpawnRadius+1) - c.SpawnRadius
		if dx*dx+dy*dy > r2 {
			continue
		}
		targetX := center.X + dx
		targetY := center.Y + dy
		tSector, _ := c.Sectors.FindByCoordinates(ctx, targetX, targetY)
		if tSector == nil {
			_ = c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
				var err error
				tSector, err = c.Provisioner.EnsureSectorExists(ctx, c.Sectors.Tx(tx), targetX, targetY)
				return err
			})
		}
		lt, err := c.Sectors.GetLocationTypeByCoordinates(ctx, targetX, targetY)
		if err != nil || lt != domain.LocationTypeEmpty {
			continue
		}
		roll := r.Float64()
		if resourceCount < c.MaxResourcefulNearby && (roll < 0.7 || dangerousCount >= c.MaxDangerousNearby) {
			resTypes := []domain.ResourceType{domain.ResourceTypeCredits, domain.ResourceTypeIron, domain.ResourceTypeTitanium, domain.ResourceTypeAntimatter}
			resType := resTypes[r.Intn(len(resTypes))]
			defense := domain.AppropriateLocationDefense(base.Stats, domain.LocationTypeResourceful) * (1.0 + r.Float64()*2.0)
			worth := domain.WorthFromDefense(defense)
			loc := domain.NewResourceLocation(
				tSector.Coordinates,
				resType,
				domain.FactionForResourceType(resType),
				worth,
				defense,
				armyProtos,
				buildProtos,
			)

			_ = c.persistResourceful(ctx, loc)
			break
		} else if dangerousCount < c.MaxDangerousNearby {
			dangFactions := []domain.Faction{
				domain.FactionCustodianProtocol,
				domain.FactionScorchWalkers,
				domain.FactionObsidianSentinels,
				domain.FactionNeuralWormApex,
			}
			faction := dangFactions[r.Intn(len(dangFactions))]
			defense := domain.AppropriateLocationDefense(base.Stats, domain.LocationTypeDangerous) * (1.0 + r.Float64()*1.5)
			worth := domain.WorthFromDefense(defense)
			loc := domain.NewDangerousLocation(
				tSector.Coordinates,
				faction,
				worth,
				defense,
				storageProtos,
				armyProtos,
				buildProtos,
			)

			_ = c.persistDangerous(ctx, loc)
			break
		}
	}
	c.reschedule(ctx, job.BaseID)
	return nil
}

func (c *WorldGenerationCommands) HandleUserBaseCreatedEvent(ctx context.Context, ev domain.UserBaseCreatedEvent) error {
	base, err := c.UserBases.FindByID(ctx, ev.BaseID)
	if err == nil {
		_ = c.spawnInitialLocations(ctx, base) // best-effort; recurring job fills any gaps
	}
	jitter := int64(rand.Intn(60))
	return c.Scheduler.Schedule(ctx, ports.SpawnNearbyLocationsJob{BaseID: ev.BaseID}, time.Now().Unix()+jitter)
}

func (c *WorldGenerationCommands) spawnInitialLocations(ctx context.Context, base *domain.UserBaseModel) error {
	armyProtos, err := c.ArmyPrototypes.FindAllPrototypes(ctx)
	if err != nil {
		return err
	}
	armyProtos = domain.FilterArmyItemPrototypesByCreationSource(armyProtos, domain.CreationSourceNPCLocation)
	buildProtos, err := c.BuildPrototypes.FindAllPrototypes(ctx)
	if err != nil {
		return err
	}
	buildProtos = domain.FilterBuildItemPrototypesByCreationSource(buildProtos, domain.CreationSourceNPCLocation)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	center := base.Coordinates

	tasks := []struct {
		resType domain.ResourceType
		dist    int
	}{
		{domain.ResourceTypeCredits, 1},
		{domain.ResourceTypeIron, 1},
		{domain.ResourceTypeCredits, 2},
		{domain.ResourceTypeIron, 2},
		{domain.ResourceTypeTitanium, 2},
		{domain.ResourceTypeAntimatter, 2},
	}

	for _, task := range tasks {
		candidates := domain.HexRing(center, task.dist)
		r.Shuffle(len(candidates), func(i, j int) {
			candidates[i], candidates[j] = candidates[j], candidates[i]
		})

		found := false
		for _, target := range candidates {
			targetX, targetY := target.X, target.Y

			_ = c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
				sRepo := c.Sectors.Tx(tx)
				rRepo := c.ResourceLocations.Tx(tx)

				tSector, err := c.Provisioner.EnsureSectorExists(ctx, sRepo, targetX, targetY)
				if err != nil {
					return err
				}

				lt, err := sRepo.GetLocationTypeByCoordinates(ctx, targetX, targetY)
				if err != nil || lt != domain.LocationTypeEmpty {
					return nil
				}

				defense := domain.AppropriateLocationDefense(base.Stats, domain.LocationTypeResourceful)
				if task.dist == 1 {
					defense *= 0.5
				}
				worth := domain.WorthFromDefense(defense)
				loc := domain.NewResourceLocation(
					tSector.Coordinates,
					task.resType,
					domain.FactionForResourceType(task.resType),
					worth,
					defense,
					armyProtos,
					buildProtos,
				)

				if err := c.Provisioner.CreateResourceLocationIfEmpty(ctx, sRepo, rRepo, loc); err != nil {
					return err
				}

				found = true
				return nil
			})

			if found {
				break
			}
		}
	}
	return nil
}

func (c *WorldGenerationCommands) persistResourceful(ctx context.Context, loc *domain.ResourceLocationModel) error {
	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		sRepo := c.Sectors.Tx(tx)
		rRepo := c.ResourceLocations.Tx(tx)
		return c.Provisioner.CreateResourceLocationIfEmpty(ctx, sRepo, rRepo, loc)
	})
}

func (c *WorldGenerationCommands) persistDangerous(ctx context.Context, loc *domain.DangerousLocationModel) error {
	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		sRepo := c.Sectors.Tx(tx)
		dRepo := c.DangerousLocations.Tx(tx)
		return c.Provisioner.CreateDangerousLocationIfEmpty(ctx, sRepo, dRepo, loc)
	})
}

func (c *WorldGenerationCommands) HandleLocationDrainedEvent(ctx context.Context, event domain.LocationDrainedEvent) error {
	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		// 1. Delete the location record
		var err error
		if event.Type == domain.LocationTypeResourceful {
			err = c.ResourceLocations.Tx(tx).DeleteByCoordinates(ctx, event.X, event.Y)
		} else {
			err = c.DangerousLocations.Tx(tx).DeleteByCoordinates(ctx, event.X, event.Y)
		}
		if err != nil {
			return err
		}

		// Note: We don't need to update the sector table because location presence
		// is checked via EXISTS on resource/dangerous tables.
		return nil
	})
}

func (c *WorldGenerationCommands) reschedule(ctx context.Context, baseID int) {
	jitter := int64(rand.Intn(300))
	_ = c.Scheduler.Schedule(ctx, ports.SpawnNearbyLocationsJob{BaseID: baseID}, time.Now().Unix()+c.RespawnPeriodSeconds+jitter)
}
