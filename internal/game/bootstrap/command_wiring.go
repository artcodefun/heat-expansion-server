package bootstrap

import (
	"context"
	"log"

	"github.com/artcodefun/heat-expansion-server/contracts/auth"
	v1 "github.com/artcodefun/heat-expansion-server/contracts/auth/v1"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	infraevents "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/events"
	infrajobs "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/jobs"
	amqp "github.com/rabbitmq/amqp091-go"
)

// WireCommandIntegrationEvents wires external integration events to command handlers.
func WireCommandIntegrationEvents(c *Commands, consumer *infraevents.RabbitMQConsumer, authExchange, authQueue string) {
	consumer.Subscribe(authExchange, authQueue, "auth.#", func(ctx context.Context, d amqp.Delivery) error {
		envelope, err := auth.Unmarshal(d.Body)
		if err != nil {
			return err
		}

		switch ev := envelope.Payload.(type) {
		case *v1.AccountRegisteredV1:
			return c.User.HandleAccountRegisteredV1Event(ctx, *ev)
		default:
			log.Printf("Received unknown identity integration event type: %T", ev)
		}
		return nil
	})
}

// WireCommandEvents subscribes command-side handlers to domain events on the in-memory publisher.
// It no-ops if the provided publisher does not support subscriptions.
func WireCommandEvents(c *Commands, pub ports.EventPublisher) {
	p, ok := pub.(*infraevents.SimplePublisher)
	if !ok {
		return
	}
	p.Listen(func(ctx context.Context, e domain.DomainEvent) error {
		switch ev := e.(type) {
		case domain.ActivityCreatedEvent:
			return c.Alert.HandleActivityCreatedEvent(ctx, ev)
		case domain.DiplomaticMessageSentEvent:
			if err := c.Diplomacy.HandleDiplomaticMessageSentEvent(ctx, ev); err != nil {
				return err
			}
			return c.Alert.HandleDiplomaticMessageSentEvent(ctx, ev)
		case domain.DiplomaticRequestCreatedEvent:
			if err := c.Diplomacy.HandleDiplomaticRequestCreatedEvent(ctx, ev); err != nil {
				return err
			}
			return c.Alert.HandleDiplomaticRequestCreatedEvent(ctx, ev)
		case domain.UserAccountCreatedEvent:
			return c.Base.HandleUserAccountCreatedEvent(ctx, ev)
		case domain.UserBaseCreatedEvent:
			return c.World.HandleUserBaseCreatedEvent(ctx, ev)
		case domain.BuildingProductionStartedEvent:
			return c.Building.HandleProductionStartedEvent(ctx, ev)
		case domain.BuildingProductionFinishedEvent:
			return c.Scanner.HandleBuildingProductionFinishedEvent(ctx, ev)
		case domain.ArmyProductionStartedEvent:
			return c.Army.HandleProductionStartedEvent(ctx, ev)
		case domain.TechResearchStartedEvent:
			return c.Tech.HandleTechResearchStartedEvent(ctx, ev)
		case domain.BuffActivatedEvent:
			return c.Storage.HandleBuffActivatedEvent(ctx, ev)
		case domain.IntelDecryptionStartedEvent:
			return c.Storage.HandleIntelDecryptionStartedEvent(ctx, ev)
		case domain.DamagedItemRestorationStartedEvent:
			return c.Storage.HandleDamagedItemRestorationStartedEvent(ctx, ev)
		case domain.MilitaryOperationStartedEvent:
			if err := c.Operation.HandleMilitaryOperationStartedEvent(ctx, ev); err != nil {
				return err
			}
			if err := c.Radar.HandleMilitaryOperationStartedEvent(ctx, ev); err != nil {
				return err
			}
			return c.Activity.HandleMilitaryOperationStartedEvent(ctx, ev)
		case domain.MilitaryOperationArrivedEvent:
			if err := c.Operation.HandleMilitaryOperationArrivedEvent(ctx, ev); err != nil {
				return err
			}
			return c.RadarThreat.HandleMilitaryOperationArrivedEvent(ctx, ev)
		case domain.MilitaryOperationCancelledEvent:
			return c.RadarThreat.HandleMilitaryOperationCancelledEvent(ctx, ev)
		case domain.RadarThreatDetectedEvent:
			return c.Activity.HandleRadarThreatDetectedEvent(ctx, ev)
		case domain.MilitaryOperationResolvedEvent:
			if err := c.Diplomacy.HandleMilitaryOperationResolvedEvent(ctx, ev); err != nil {
				return err
			}
			return c.Activity.HandleMilitaryOperationResolvedEvent(ctx, ev)
		case domain.MilitaryOperationReturnStartedEvent:
			return c.Operation.HandleMilitaryOperationReturnStartedEvent(ctx, ev)
		case domain.MilitaryOperationReturnArrivedEvent:
			return c.Operation.HandleMilitaryOperationReturnArrivedEvent(ctx, ev)
		case domain.ScanReportCreatedEvent:
			return c.Activity.HandleScanReportCreatedEvent(ctx, ev)
		case domain.LocationDrainedEvent:
			return c.World.HandleLocationDrainedEvent(ctx, ev)
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
	s.Listen(func(ctx context.Context, j ports.SchadulableJob) error {
		switch job := j.(type) {
		case ports.MoveBuildQueueJob:
			return c.Building.HandleMoveBuildQueueJob(ctx, job)
		case ports.MoveArmyQueueJob:
			return c.Army.HandleMoveArmyQueueJob(ctx, job)
		case ports.MoveTechQueueJob:
			return c.Tech.HandleMoveTechQueueJob(ctx, job)
		case ports.UpdateMilitaryOperationJob:
			return c.Operation.HandleUpdateMilitaryOperationJob(ctx, job)
		case ports.ExpireDiplomaticRequestJob:
			return c.Diplomacy.HandleExpireDiplomaticRequestJob(ctx, job)
		case ports.IntelligenceScanJob:
			return c.Scanner.HandleIntelligenceScanJob(ctx, job)
		case ports.IntelligenceRadarJob:
			return c.Radar.HandleIntelligenceRadarJob(ctx, job)
		case ports.DeleteExpiredBuffJob:
			return c.Storage.HandleDeleteExpiredBuffJob(ctx, job)
		case ports.DecryptIntelItemJob:
			return c.Storage.HandleDecryptIntelItemJob(ctx, job)
		case ports.RestoreDamagedItemJob:
			return c.Storage.HandleRestoreDamagedItemJob(ctx, job)
		case ports.SpawnNearbyLocationsJob:
			return c.World.HandleSpawnNearbyLocationsJob(ctx, job)
		}
		return nil
	})
}
