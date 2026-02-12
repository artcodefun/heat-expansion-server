package bootstrap

import (
	"github.com/artcodefun/heat-expansion-api/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-api/internal/game/domain"
	infraevents "github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/events"
	infrajobs "github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/jobs"
)

// WireCommandEvents subscribes command-side handlers to domain events on the in-memory publisher.
// It no-ops if the provided publisher does not support subscriptions.
func WireCommandEvents(c *Commands, pub ports.EventPublisher) {
	p, ok := pub.(*infraevents.SimplePublisher)
	if !ok {
		return
	}
	p.Listen(func(e domain.DomainEvent) error {
		switch ev := e.(type) {
		case domain.ActivityCreatedEvent:
			return c.Alert.HandleActivityCreatedEvent(ev)
		case domain.UserAccountCreatedEvent:
			return c.Base.HandleUserAccountCreatedEvent(ev)
		case domain.UserBaseCreatedEvent:
			return c.World.HandleUserBaseCreatedEvent(ev)
		case domain.BuildingProductionStartedEvent:
			return c.Building.HandleProductionStartedEvent(&ev)
		case domain.BuildingProductionFinishedEvent:
			return c.Scanner.HandleBuildingProductionFinishedEvent(&ev)
		case domain.ArmyProductionStartedEvent:
			return c.Army.HandleProductionStartedEvent(&ev)
		case domain.TechResearchStartedEvent:
			return c.Tech.HandleTechResearchStartedEvent(&ev)
		case domain.BuffActivatedEvent:
			return c.Storage.HandleBuffActivatedEvent(&ev)
		case domain.IntelDecryptionStartedEvent:
			return c.Storage.HandleIntelDecryptionStartedEvent(&ev)
		case domain.DamagedItemRestorationStartedEvent:
			return c.Storage.HandleDamagedItemRestorationStartedEvent(&ev)
		case domain.MilitaryOperationStartedEvent:
			if err := c.Operation.HandleMilitaryOperationStartedEvent(ev); err != nil {
				return err
			}
			if err := c.Radar.HandleMilitaryOperationStartedEvent(ev); err != nil {
				return err
			}
			return c.Activity.HandleMilitaryOperationStartedEvent(ev)
		case domain.MilitaryOperationArrivedEvent:
			if err := c.Operation.HandleMilitaryOperationArrivedEvent(ev); err != nil {
				return err
			}
			return c.RadarThreat.HandleMilitaryOperationArrivedEvent(ev)
		case domain.MilitaryOperationCancelledEvent:
			return c.RadarThreat.HandleMilitaryOperationCancelledEvent(ev)
		case domain.RadarThreatDetectedEvent:
			return c.Activity.HandleRadarThreatDetectedEvent(ev)
		case domain.MilitaryOperationResolvedEvent:
			return c.Activity.HandleMilitaryOperationResolvedEvent(ev)
		case domain.MilitaryOperationReturnStartedEvent:
			return c.Operation.HandleMilitaryOperationReturnStartedEvent(ev)
		case domain.MilitaryOperationReturnArrivedEvent:
			return c.Operation.HandleMilitaryOperationReturnArrivedEvent(ev)
		case domain.ScanReportCreatedEvent:
			return c.Activity.HandleScanReportCreatedEvent(ev)
		case domain.LocationDrainedEvent:
			return c.World.HandleLocationDrainedEvent(ev)
		}
		return nil
	})
}

// WireCommandSchedulerHandler binds a job dispatcher to the in-memory scheduler so payloads route to command handlers.
// It no-ops if the scheduler is not the in-memory implementation.
func WireCommandSchedulerHandler(c *Commands, sch ports.Scheduler) {
	s, ok := sch.(*infrajobs.DBScheduler)
	if !ok {
		return
	}
	s.Listen(func(j ports.SchadulableJob) error {
		switch job := j.(type) {
		case ports.MoveBuildQueueJob:
			return c.Building.HandleMoveBuildQueueJob(job)
		case ports.MoveArmyQueueJob:
			return c.Army.HandleMoveArmyQueueJob(job)
		case ports.MoveTechQueueJob:
			return c.Tech.HandleMoveTechQueueJob(job)
		case ports.UpdateMilitaryOperationJob:
			return c.Operation.HandleUpdateMilitaryOperationJob(job)
		case ports.IntelligenceScanJob:
			return c.Scanner.HandleIntelligenceScanJob(job)
		case ports.IntelligenceRadarJob:
			return c.Radar.HandleIntelligenceRadarJob(job)
		case ports.DeleteExpiredBuffJob:
			return c.Storage.HandleDeleteExpiredBuffJob(job)
		case ports.DecryptIntelItemJob:
			return c.Storage.HandleDecryptIntelItemJob(job)
		case ports.RestoreDamagedItemJob:
			return c.Storage.HandleRestoreDamagedItemJob(job)
		case ports.SpawnNearbyLocationsJob:
			return c.World.HandleSpawnNearbyLocationsJob(job)
		}
		return nil
	})
}
