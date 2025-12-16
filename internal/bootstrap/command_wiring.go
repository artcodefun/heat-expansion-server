package bootstrap

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	infraevents "github.com/artcodefun/heat-expansion-api/internal/infrastructure/events"
	infrajobs "github.com/artcodefun/heat-expansion-api/internal/infrastructure/jobs"
)

// WireCommandEvents subscribes command-side handlers to domain events on the in-memory publisher.
// It no-ops if the provided publisher does not support subscriptions.
func WireCommandEvents(c *Commands, pub ports.EventPublisher) {
	p, ok := pub.(*infraevents.InMemoryPublisher)
	if !ok {
		return
	}
	p.Subscribe(func(e domain.DomainEvent) {
		switch ev := e.(type) {
		case domain.UserAccountCreatedEvent:
			_ = c.Base.HandleUserAccountCreatedEvent(ev)
		case domain.BuildingProductionStartedEvent:
			_ = c.Building.HandleProductionStartedEvent(&ev)
		case domain.BuildingProductionFinishedEvent:
			_ = c.Building.HandleProductionFinishedEvent(&ev)
		case domain.ArmyProductionStartedEvent:
			_ = c.Army.HandleProductionStartedEvent(&ev)
		case domain.TechResearchStartedEvent:
			_ = c.Tech.HandleTechResearchStartedEvent(&ev)
		case domain.BuffActivatedEvent:
			_ = c.Storage.HandleBuffActivatedEvent(&ev)
		case domain.MilitaryOperationStartedEvent:
			_ = c.Operation.HandleMilitaryOperationStartedEvent(ev)
			_ = c.Activity.HandleMilitaryOperationStartedEvent(ev)
		case domain.MilitaryOperationArrivedEvent:
			_ = c.Operation.HandleMilitaryOperationArrivedEvent(ev)
		case domain.MilitaryOperationResolvedEvent:
			_ = c.Activity.HandleMilitaryOperationResolvedEvent(ev)
		case domain.MilitaryOperationReturnStartedEvent:
			_ = c.Operation.HandleMilitaryOperationReturnStartedEvent(ev)
		case domain.MilitaryOperationReturnArrivedEvent:
			_ = c.Operation.HandleMilitaryOperationReturnArrivedEvent(ev)
		case domain.ScanReportCreatedEvent:
			_ = c.Activity.HandleScanReportCreatedEvent(ev)
		}
	})
}

// WireCommandSchedulerHandler binds a job dispatcher to the in-memory scheduler so payloads route to command handlers.
// It no-ops if the scheduler is not the in-memory implementation.
func WireCommandSchedulerHandler(c *Commands, sch ports.Scheduler) {
	s, ok := sch.(*infrajobs.DBScheduler)
	if !ok {
		return
	}
	s.Subscribe(func(j ports.SchadulableJob) {
		switch job := j.(type) {
		case ports.MoveBuildQueueJob:
			_ = c.Building.HandleMoveBuildQueueJob(job)
		case ports.MoveArmyQueueJob:
			_ = c.Army.HandleMoveArmyQueueJob(job)
		case ports.MoveTechQueueJob:
			_ = c.Tech.HandleMoveTechQueueJob(job)
		case ports.UpdateMilitaryOperationJob:
			_ = c.Operation.HandleUpdateMilitaryOperationJob(job)
		case ports.RadarScanJob:
			_ = c.Radar.HandleRadarScanJob(job)
		case ports.DeleteExpiredBuffJob:
			_, _ = c.Storage.HandleDeleteExpiredBuffJob(job.BaseID)
		case ports.SpawnNearbyLocationsJob:
			_ = c.World.HandleSpawnNearbyLocationsJob(job)
		}
	})
}
