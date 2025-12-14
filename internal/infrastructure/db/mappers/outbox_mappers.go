package mappers

import (
	"encoding/json"
	"errors"

	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/dtos"
)

// Job kind identifiers used for serialization. These are an implementation
// detail of the DB outbox; core code only sees typed jobs.
const (
	jobKindMoveBuildQueue       = "MOVE_BUILD_QUEUE"
	jobKindMoveArmyQueue        = "MOVE_ARMY_QUEUE"
	jobKindMoveTechQueue        = "MOVE_TECH_QUEUE"
	jobKindDeleteExpiredBuff    = "DELETE_EXPIRED_BUFF"
	jobKindUpdateMilitaryOp     = "UPDATE_MILITARY_OPERATION"
	jobKindSpawnNearbyLocations = "SPAWN_NEARBY_LOCATIONS"
	jobKindRadarScan            = "RADAR_SCAN"
)

// EncodeJob converts a typed Scheduler job into a kind label and JSON payload
// suitable for persistence.
func EncodeJob(job ports.SchadulableJob) (kind string, payload []byte, err error) {
	switch j := job.(type) {
	case ports.MoveBuildQueueJob:
		payload, err = json.Marshal(dtos.MoveBuildQueueJobDTOFromDomain(j))
		return jobKindMoveBuildQueue, payload, err
	case ports.MoveArmyQueueJob:
		payload, err = json.Marshal(dtos.MoveArmyQueueJobDTOFromDomain(j))
		return jobKindMoveArmyQueue, payload, err
	case ports.MoveTechQueueJob:
		payload, err = json.Marshal(dtos.MoveTechQueueJobDTOFromDomain(j))
		return jobKindMoveTechQueue, payload, err
	case ports.DeleteExpiredBuffJob:
		payload, err = json.Marshal(dtos.DeleteExpiredBuffJobDTOFromDomain(j))
		return jobKindDeleteExpiredBuff, payload, err
	case ports.UpdateMilitaryOperationJob:
		payload, err = json.Marshal(dtos.UpdateMilitaryOperationJobDTOFromDomain(j))
		return jobKindUpdateMilitaryOp, payload, err
	case ports.SpawnNearbyLocationsJob:
		payload, err = json.Marshal(dtos.SpawnNearbyLocationsJobDTOFromDomain(j))
		return jobKindSpawnNearbyLocations, payload, err
	case ports.RadarScanJob:
		payload, err = json.Marshal(dtos.RadarScanJobDTOFromDomain(j))
		return jobKindRadarScan, payload, err
	default:
		return "", nil, errors.New("unsupported scheduler job type")
	}
}

// DecodeJob reconstructs a typed Scheduler job from its kind label and JSON
// payload.
func DecodeJob(kind string, payload []byte) (ports.SchadulableJob, error) {
	switch kind {
	case jobKindMoveBuildQueue:
		var dto dtos.MoveBuildQueueJobDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.MoveBuildQueueJobFromDTO(dto), nil
	case jobKindMoveArmyQueue:
		var dto dtos.MoveArmyQueueJobDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.MoveArmyQueueJobFromDTO(dto), nil
	case jobKindMoveTechQueue:
		var dto dtos.MoveTechQueueJobDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.MoveTechQueueJobFromDTO(dto), nil
	case jobKindDeleteExpiredBuff:
		var dto dtos.DeleteExpiredBuffJobDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.DeleteExpiredBuffJobFromDTO(dto), nil
	case jobKindUpdateMilitaryOp:
		var dto dtos.UpdateMilitaryOperationJobDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.UpdateMilitaryOperationJobFromDTO(dto), nil
	case jobKindSpawnNearbyLocations:
		var dto dtos.SpawnNearbyLocationsJobDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.SpawnNearbyLocationsJobFromDTO(dto), nil
	case jobKindRadarScan:
		var dto dtos.RadarScanJobDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.RadarScanJobFromDTO(dto), nil
	default:
		return nil, errors.New("unknown scheduler job kind")
	}
}

// Event kind identifiers used for serialization. These are an implementation
// detail of the DB outbox; core code only sees typed events.
const (
	evKindUserAccountCreated             = "USER_ACCOUNT_CREATED"
	evKindBuildingProductionStarted      = "BUILDING_PRODUCTION_STARTED"
	evKindBuildingProductionFinished     = "BUILDING_PRODUCTION_FINISHED"
	evKindBuildingProductionCancelled    = "BUILDING_PRODUCTION_CANCELLED"
	evKindBuildingProductionSpeedup      = "BUILDING_PRODUCTION_SPEEDUP"
	evKindBuildingPresentDeleted         = "BUILDING_PRESENT_DELETED"
	evKindArmyProductionPending          = "ARMY_PRODUCTION_PENDING"
	evKindArmyProductionStarted          = "ARMY_PRODUCTION_STARTED"
	evKindArmyProductionFinished         = "ARMY_PRODUCTION_FINISHED"
	evKindArmyProductionCancelled        = "ARMY_PRODUCTION_CANCELLED"
	evKindArmyProductionSpeedup          = "ARMY_PRODUCTION_SPEEDUP"
	evKindArmyPresentDeleted             = "ARMY_PRESENT_DELETED"
	evKindTechResearchStarted            = "TECH_RESEARCH_STARTED"
	evKindTechResearchFinished           = "TECH_RESEARCH_FINISHED"
	evKindTechResearchSpeedup            = "TECH_RESEARCH_SPEEDUP"
	evKindStorageItemPresentDeleted      = "STORAGE_ITEM_PRESENT_DELETED"
	evKindBuffActivated                  = "BUFF_ACTIVATED"
	evKindMilitaryOperationStarted       = "MILITARY_OPERATION_STARTED"
	evKindMilitaryOperationArrived       = "MILITARY_OPERATION_ARRIVED"
	evKindMilitaryOperationResolved      = "MILITARY_OPERATION_RESOLVED"
	evKindMilitaryOperationReturnStarted = "MILITARY_OPERATION_RETURN_STARTED"
	evKindMilitaryOperationReturnArrived = "MILITARY_OPERATION_RETURN_ARRIVED"
	evKindScanReportCreated              = "SCAN_REPORT_CREATED"
)

// EncodeEvent converts a typed DomainEvent into a kind label and JSON payload
// suitable for persistence.
func EncodeEvent(ev domain.DomainEvent) (kind string, payload []byte, err error) {
	switch e := ev.(type) {
	case domain.UserAccountCreatedEvent:
		payload, err = json.Marshal(dtos.UserAccountCreatedEventDTOFromDomain(e))
		return evKindUserAccountCreated, payload, err

	case domain.BuildingProductionStartedEvent:
		payload, err = json.Marshal(dtos.BuildingProductionStartedEventDTOFromDomain(e))
		return evKindBuildingProductionStarted, payload, err
	case domain.BuildingProductionFinishedEvent:
		payload, err = json.Marshal(dtos.BuildingProductionFinishedEventDTOFromDomain(e))
		return evKindBuildingProductionFinished, payload, err
	case domain.BuildingProductionCancelledEvent:
		payload, err = json.Marshal(dtos.BuildingProductionCancelledEventDTOFromDomain(e))
		return evKindBuildingProductionCancelled, payload, err
	case domain.BuildingProductionSpeedupEvent:
		payload, err = json.Marshal(dtos.BuildingProductionSpeedupEventDTOFromDomain(e))
		return evKindBuildingProductionSpeedup, payload, err
	case domain.BuildingPresentDeletedEvent:
		payload, err = json.Marshal(dtos.BuildingPresentDeletedEventDTOFromDomain(e))
		return evKindBuildingPresentDeleted, payload, err

	case domain.ArmyProductionPendingEvent:
		payload, err = json.Marshal(dtos.ArmyProductionPendingEventDTOFromDomain(e))
		return evKindArmyProductionPending, payload, err
	case domain.ArmyProductionStartedEvent:
		payload, err = json.Marshal(dtos.ArmyProductionStartedEventDTOFromDomain(e))
		return evKindArmyProductionStarted, payload, err
	case domain.ArmyProductionFinishedEvent:
		payload, err = json.Marshal(dtos.ArmyProductionFinishedEventDTOFromDomain(e))
		return evKindArmyProductionFinished, payload, err
	case domain.ArmyProductionCancelledEvent:
		payload, err = json.Marshal(dtos.ArmyProductionCancelledEventDTOFromDomain(e))
		return evKindArmyProductionCancelled, payload, err
	case domain.ArmyProductionSpeedupEvent:
		payload, err = json.Marshal(dtos.ArmyProductionSpeedupEventDTOFromDomain(e))
		return evKindArmyProductionSpeedup, payload, err
	case domain.ArmyPresentDeletedEvent:
		payload, err = json.Marshal(dtos.ArmyPresentDeletedEventDTOFromDomain(e))
		return evKindArmyPresentDeleted, payload, err

	case domain.TechResearchStartedEvent:
		payload, err = json.Marshal(dtos.TechResearchStartedEventDTOFromDomain(e))
		return evKindTechResearchStarted, payload, err
	case domain.TechResearchFinishedEvent:
		payload, err = json.Marshal(dtos.TechResearchFinishedEventDTOFromDomain(e))
		return evKindTechResearchFinished, payload, err
	case domain.TechResearchSpeedupEvent:
		payload, err = json.Marshal(dtos.TechResearchSpeedupEventDTOFromDomain(e))
		return evKindTechResearchSpeedup, payload, err

	case domain.StorageItemPresentDeletedEvent:
		payload, err = json.Marshal(dtos.StorageItemPresentDeletedEventDTOFromDomain(e))
		return evKindStorageItemPresentDeleted, payload, err
	case domain.BuffActivatedEvent:
		payload, err = json.Marshal(dtos.BuffActivatedEventDTOFromDomain(e))
		return evKindBuffActivated, payload, err

	case domain.MilitaryOperationStartedEvent:
		payload, err = json.Marshal(dtos.MilitaryOperationStartedEventDTOFromDomain(e))
		return evKindMilitaryOperationStarted, payload, err
	case domain.MilitaryOperationArrivedEvent:
		payload, err = json.Marshal(dtos.MilitaryOperationArrivedEventDTOFromDomain(e))
		return evKindMilitaryOperationArrived, payload, err
	case domain.MilitaryOperationResolvedEvent:
		payload, err = json.Marshal(dtos.MilitaryOperationResolvedEventDTOFromDomain(e))
		return evKindMilitaryOperationResolved, payload, err
	case domain.MilitaryOperationReturnStartedEvent:
		payload, err = json.Marshal(dtos.MilitaryOperationReturnStartedEventDTOFromDomain(e))
		return evKindMilitaryOperationReturnStarted, payload, err
	case domain.MilitaryOperationReturnArrivedEvent:
		payload, err = json.Marshal(dtos.MilitaryOperationReturnArrivedEventDTOFromDomain(e))
		return evKindMilitaryOperationReturnArrived, payload, err

	case domain.ScanReportCreatedEvent:
		payload, err = json.Marshal(dtos.ScanReportCreatedEventDTOFromDomain(e))
		return evKindScanReportCreated, payload, err

	default:
		return "", nil, errors.New("unsupported domain event type")
	}
}

// DecodeEvent reconstructs a DomainEvent from its kind label and JSON payload.
func DecodeEvent(kind string, payload []byte) (domain.DomainEvent, error) {
	switch kind {
	case evKindUserAccountCreated:
		var dto dtos.UserAccountCreatedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.UserAccountCreatedEventFromDTO(dto), nil

	case evKindBuildingProductionStarted:
		var dto dtos.BuildingProductionStartedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.BuildingProductionStartedEventFromDTO(dto), nil
	case evKindBuildingProductionFinished:
		var dto dtos.BuildingProductionFinishedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.BuildingProductionFinishedEventFromDTO(dto), nil
	case evKindBuildingProductionCancelled:
		var dto dtos.BuildingProductionCancelledEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.BuildingProductionCancelledEventFromDTO(dto), nil
	case evKindBuildingProductionSpeedup:
		var dto dtos.BuildingProductionSpeedupEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.BuildingProductionSpeedupEventFromDTO(dto), nil
	case evKindBuildingPresentDeleted:
		var dto dtos.BuildingPresentDeletedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.BuildingPresentDeletedEventFromDTO(dto), nil

	case evKindArmyProductionPending:
		var dto dtos.ArmyProductionPendingEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.ArmyProductionPendingEventFromDTO(dto), nil
	case evKindArmyProductionStarted:
		var dto dtos.ArmyProductionStartedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.ArmyProductionStartedEventFromDTO(dto), nil
	case evKindArmyProductionFinished:
		var dto dtos.ArmyProductionFinishedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.ArmyProductionFinishedEventFromDTO(dto), nil
	case evKindArmyProductionCancelled:
		var dto dtos.ArmyProductionCancelledEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.ArmyProductionCancelledEventFromDTO(dto), nil
	case evKindArmyProductionSpeedup:
		var dto dtos.ArmyProductionSpeedupEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.ArmyProductionSpeedupEventFromDTO(dto), nil
	case evKindArmyPresentDeleted:
		var dto dtos.ArmyPresentDeletedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.ArmyPresentDeletedEventFromDTO(dto), nil

	case evKindTechResearchStarted:
		var dto dtos.TechResearchStartedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.TechResearchStartedEventFromDTO(dto), nil
	case evKindTechResearchFinished:
		var dto dtos.TechResearchFinishedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.TechResearchFinishedEventFromDTO(dto), nil
	case evKindTechResearchSpeedup:
		var dto dtos.TechResearchSpeedupEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.TechResearchSpeedupEventFromDTO(dto), nil

	case evKindStorageItemPresentDeleted:
		var dto dtos.StorageItemPresentDeletedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.StorageItemPresentDeletedEventFromDTO(dto), nil
	case evKindBuffActivated:
		var dto dtos.BuffActivatedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.BuffActivatedEventFromDTO(dto), nil

	case evKindMilitaryOperationStarted:
		var dto dtos.MilitaryOperationStartedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.MilitaryOperationStartedEventFromDTO(dto), nil
	case evKindMilitaryOperationArrived:
		var dto dtos.MilitaryOperationArrivedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.MilitaryOperationArrivedEventFromDTO(dto), nil
	case evKindMilitaryOperationResolved:
		var dto dtos.MilitaryOperationResolvedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.MilitaryOperationResolvedEventFromDTO(dto), nil
	case evKindMilitaryOperationReturnStarted:
		var dto dtos.MilitaryOperationReturnStartedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.MilitaryOperationReturnStartedEventFromDTO(dto), nil
	case evKindMilitaryOperationReturnArrived:
		var dto dtos.MilitaryOperationReturnArrivedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.MilitaryOperationReturnArrivedEventFromDTO(dto), nil

	case evKindScanReportCreated:
		var dto dtos.ScanReportCreatedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.ScanReportCreatedEventFromDTO(dto), nil

	default:
		return nil, errors.New("unknown domain event kind")
	}
}
