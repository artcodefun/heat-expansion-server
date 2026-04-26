package mappers

import (
	"encoding/json"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/dtos"
)

// Job kind identifiers used for serialization. These are an implementation
// detail of the DB outbox; core code only sees typed jobs.
const (
	jobKindMoveBuildQueue       = "MOVE_BUILD_QUEUE"
	jobKindMoveArmyQueue        = "MOVE_ARMY_QUEUE"
	jobKindMoveTechQueue        = "MOVE_TECH_QUEUE"
	jobKindDeleteExpiredBuff    = "DELETE_EXPIRED_BUFF"
	jobKindRestoreDamagedItem   = "RESTORE_DAMAGED_ITEM"
	jobKindDecryptIntelItem     = "DECRYPT_INTEL_ITEM"
	jobKindUpdateMilitaryOp     = "UPDATE_MILITARY_OPERATION"
	jobKindUpdateTradeOp        = "UPDATE_TRADE_OPERATION"
	jobKindExpireDiplomaticReq  = "EXPIRE_DIPLOMATIC_REQUEST"
	jobKindExpireTradeOp        = "EXPIRE_TRADE_OPERATION"
	jobKindSpawnNearbyLocations = "SPAWN_NEARBY_LOCATIONS"
	jobKindIntelligenceScan     = "INTELLIGENCE_SCAN"
	jobKindIntelligenceRadar    = "INTELLIGENCE_RADAR"
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
	case ports.RestoreDamagedItemJob:
		payload, err = json.Marshal(dtos.RestoreDamagedItemJobDTOFromDomain(j))
		return jobKindRestoreDamagedItem, payload, err
	case ports.DecryptIntelItemJob:
		payload, err = json.Marshal(dtos.DecryptIntelItemJobDTOFromDomain(j))
		return jobKindDecryptIntelItem, payload, err
	case ports.UpdateMilitaryOperationJob:
		payload, err = json.Marshal(dtos.UpdateMilitaryOperationJobDTOFromDomain(j))
		return jobKindUpdateMilitaryOp, payload, err
	case ports.UpdateTradeOperationJob:
		payload, err = json.Marshal(dtos.UpdateTradeOperationJobDTOFromDomain(j))
		return jobKindUpdateTradeOp, payload, err
	case ports.ExpireDiplomaticRequestJob:
		payload, err = json.Marshal(dtos.ExpireDiplomaticRequestJobDTOFromDomain(j))
		return jobKindExpireDiplomaticReq, payload, err
	case ports.ExpireTradeOperationJob:
		payload, err = json.Marshal(dtos.ExpireTradeOperationJobDTOFromDomain(j))
		return jobKindExpireTradeOp, payload, err
	case ports.SpawnNearbyLocationsJob:
		payload, err = json.Marshal(dtos.SpawnNearbyLocationsJobDTOFromDomain(j))
		return jobKindSpawnNearbyLocations, payload, err
	case ports.IntelligenceScanJob:
		payload, err = json.Marshal(dtos.IntelligenceScanJobDTOFromDomain(j))
		return jobKindIntelligenceScan, payload, err
	case ports.IntelligenceRadarJob:
		payload, err = json.Marshal(dtos.IntelligenceRadarJobDTOFromDomain(j))
		return jobKindIntelligenceRadar, payload, err
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
	case jobKindRestoreDamagedItem:
		var dto dtos.RestoreDamagedItemJobDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.RestoreDamagedItemJobFromDTO(dto), nil
	case jobKindDecryptIntelItem:
		var dto dtos.DecryptIntelItemJobDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.DecryptIntelItemJobFromDTO(dto), nil
	case jobKindUpdateMilitaryOp:
		var dto dtos.UpdateMilitaryOperationJobDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.UpdateMilitaryOperationJobFromDTO(dto), nil
	case jobKindUpdateTradeOp:
		var dto dtos.UpdateTradeOperationJobDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.UpdateTradeOperationJobFromDTO(dto), nil
	case jobKindExpireDiplomaticReq:
		var dto dtos.ExpireDiplomaticRequestJobDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.ExpireDiplomaticRequestJobFromDTO(dto), nil
	case jobKindExpireTradeOp:
		var dto dtos.ExpireTradeOperationJobDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.ExpireTradeOperationJobFromDTO(dto), nil
	case jobKindSpawnNearbyLocations:
		var dto dtos.SpawnNearbyLocationsJobDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.SpawnNearbyLocationsJobFromDTO(dto), nil
	case jobKindIntelligenceScan:
		var dto dtos.IntelligenceScanJobDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.IntelligenceScanJobFromDTO(dto), nil
	case jobKindIntelligenceRadar:
		var dto dtos.IntelligenceRadarJobDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.IntelligenceRadarJobFromDTO(dto), nil
	default:
		return nil, errors.New("unknown scheduler job kind")
	}
}

// Event kind identifiers used for serialization. These are an implementation
// detail of the DB outbox; core code only sees typed events.
const (
	evKindUserAccountCreated             = "USER_ACCOUNT_CREATED"
	evKindUserBaseCreated                = "USER_BASE_CREATED"
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
	evKindIntelDecryptionStarted         = "INTEL_DECRYPTION_STARTED"
	evKindIntelDecryptionFinished        = "INTEL_DECRYPTION_FINISHED"
	evKindDamagedItemRestorationStarted  = "DAMAGED_ITEM_RESTORATION_STARTED"
	evKindDamagedItemRestored            = "DAMAGED_ITEM_RESTORED"
	evKindArtifactActivated              = "ARTIFACT_ACTIVATED"
	evKindArtifactDeactivated            = "ARTIFACT_DEACTIVATED"
	evKindMilitaryOperationStarted       = "MILITARY_OPERATION_STARTED"
	evKindMilitaryOperationArrived       = "MILITARY_OPERATION_ARRIVED"
	evKindMilitaryOperationResolved      = "MILITARY_OPERATION_RESOLVED"
	evKindMilitaryOperationReturnStarted = "MILITARY_OPERATION_RETURN_STARTED"
	evKindMilitaryOperationReturnArrived = "MILITARY_OPERATION_RETURN_ARRIVED"
	evKindMilitaryOperationCancelled     = "MILITARY_OPERATION_CANCELLED"
	evKindTradeOperationCreated          = "TRADE_OPERATION_CREATED"
	evKindTradeOperationAccepted         = "TRADE_OPERATION_ACCEPTED"
	evKindTradeOperationDeclined         = "TRADE_OPERATION_DECLINED"
	evKindTradeOperationExpired          = "TRADE_OPERATION_EXPIRED"
	evKindTradeOperationCancelled        = "TRADE_OPERATION_CANCELLED_BY_INITIATOR"
	evKindTradeOperationOutbound         = "TRADE_OPERATION_OUTBOUND"
	evKindTradeOperationArrived          = "TRADE_OPERATION_ARRIVED"
	evKindTradeOperationReturning        = "TRADE_OPERATION_RETURNING"
	evKindTradeOperationReturnArrived    = "TRADE_OPERATION_RETURN_ARRIVED"
	evKindScanReportCreated              = "SCAN_REPORT_CREATED"
	evKindRadarThreatDetected            = "RADAR_THREAT_DETECTED"
	evKindActivityCreated                = "ACTIVITY_CREATED"
	evKindDiplomaticMessageSent          = "DIPLOMATIC_MESSAGE_SENT"
	evKindDiplomaticRequestCreated       = "DIPLOMATIC_REQUEST_CREATED"
	evKindDiplomaticRelationshipCreated  = "DIPLOMATIC_RELATIONSHIP_CREATED"
	evKindLocationDrained                = "LOCATION_DRAINED"
)

// EncodeEvent converts a typed DomainEvent into a kind label and JSON payload
// suitable for persistence.
func EncodeEvent(ev domain.DomainEvent) (kind string, payload []byte, err error) {
	switch e := ev.(type) {
	case domain.UserAccountCreatedEvent:
		payload, err = json.Marshal(dtos.UserAccountCreatedEventDTOFromDomain(e))
		return evKindUserAccountCreated, payload, err

	case domain.UserBaseCreatedEvent:
		payload, err = json.Marshal(dtos.UserBaseCreatedEventDTOFromDomain(e))
		return evKindUserBaseCreated, payload, err

	case domain.LocationDrainedEvent:
		payload, err = json.Marshal(dtos.LocationDrainedEventDTOFromDomain(e))
		return evKindLocationDrained, payload, err

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
	case domain.IntelDecryptionStartedEvent:
		payload, err = json.Marshal(dtos.IntelDecryptionStartedEventDTOFromDomain(e))
		return evKindIntelDecryptionStarted, payload, err
	case domain.IntelDecryptionFinishedEvent:
		payload, err = json.Marshal(dtos.IntelDecryptionFinishedEventDTOFromDomain(e))
		return evKindIntelDecryptionFinished, payload, err
	case domain.DamagedItemRestorationStartedEvent:
		payload, err = json.Marshal(dtos.DamagedItemRestorationStartedEventDTOFromDomain(e))
		return evKindDamagedItemRestorationStarted, payload, err
	case domain.DamagedItemRestoredEvent:
		payload, err = json.Marshal(dtos.DamagedItemRestoredEventDTOFromDomain(e))
		return evKindDamagedItemRestored, payload, err
	case domain.ArtifactActivatedEvent:
		payload, err = json.Marshal(dtos.ArtifactActivatedEventDTOFromDomain(e))
		return evKindArtifactActivated, payload, err
	case domain.ArtifactDeactivatedEvent:
		payload, err = json.Marshal(dtos.ArtifactDeactivatedEventDTOFromDomain(e))
		return evKindArtifactDeactivated, payload, err

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
	case domain.MilitaryOperationCancelledEvent:
		payload, err = json.Marshal(dtos.MilitaryOperationCancelledEventDTOFromDomain(e))
		return evKindMilitaryOperationCancelled, payload, err

	case domain.TradeOperationCreatedEvent:
		payload, err = json.Marshal(dtos.TradeOperationCreatedEventDTOFromDomain(e))
		return evKindTradeOperationCreated, payload, err
	case domain.TradeOperationAcceptedEvent:
		payload, err = json.Marshal(dtos.TradeOperationAcceptedEventDTOFromDomain(e))
		return evKindTradeOperationAccepted, payload, err
	case domain.TradeOperationDeclinedEvent:
		payload, err = json.Marshal(dtos.TradeOperationDeclinedEventDTOFromDomain(e))
		return evKindTradeOperationDeclined, payload, err
	case domain.TradeOperationExpiredEvent:
		payload, err = json.Marshal(dtos.TradeOperationExpiredEventDTOFromDomain(e))
		return evKindTradeOperationExpired, payload, err
	case domain.TradeOperationCancelledByInitiatorEvent:
		payload, err = json.Marshal(dtos.TradeOperationCancelledByInitiatorEventDTOFromDomain(e))
		return evKindTradeOperationCancelled, payload, err
	case domain.TradeOperationOutboundEvent:
		payload, err = json.Marshal(dtos.TradeOperationOutboundEventDTOFromDomain(e))
		return evKindTradeOperationOutbound, payload, err
	case domain.TradeOperationArrivedEvent:
		payload, err = json.Marshal(dtos.TradeOperationArrivedEventDTOFromDomain(e))
		return evKindTradeOperationArrived, payload, err
	case domain.TradeOperationReturningEvent:
		payload, err = json.Marshal(dtos.TradeOperationReturningEventDTOFromDomain(e))
		return evKindTradeOperationReturning, payload, err
	case domain.TradeOperationReturnArrivedEvent:
		payload, err = json.Marshal(dtos.TradeOperationReturnArrivedEventDTOFromDomain(e))
		return evKindTradeOperationReturnArrived, payload, err

	case domain.ScanReportCreatedEvent:
		payload, err = json.Marshal(dtos.ScanReportCreatedEventDTOFromDomain(e))
		return evKindScanReportCreated, payload, err
	case domain.RadarThreatDetectedEvent:
		payload, err = json.Marshal(dtos.RadarThreatDetectedEventDTOFromDomain(e))
		return evKindRadarThreatDetected, payload, err

	case domain.ActivityCreatedEvent:
		payload, err = json.Marshal(dtos.ActivityCreatedEventDTOFromDomain(e))
		return evKindActivityCreated, payload, err
	case domain.DiplomaticMessageSentEvent:
		payload, err = json.Marshal(dtos.DiplomaticMessageSentEventDTOFromDomain(e))
		return evKindDiplomaticMessageSent, payload, err
	case domain.DiplomaticRequestCreatedEvent:
		payload, err = json.Marshal(dtos.DiplomaticRequestCreatedEventDTOFromDomain(e))
		return evKindDiplomaticRequestCreated, payload, err
	case domain.DiplomaticRelationshipCreatedEvent:
		payload, err = json.Marshal(dtos.DiplomaticRelationshipCreatedEventDTOFromDomain(e))
		return evKindDiplomaticRelationshipCreated, payload, err

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

	case evKindUserBaseCreated:
		var dto dtos.UserBaseCreatedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.UserBaseCreatedEventFromDTO(dto), nil
	case evKindDiplomaticRequestCreated:
		var dto dtos.DiplomaticRequestCreatedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.DiplomaticRequestCreatedEventFromDTO(dto), nil

	case evKindLocationDrained:
		var dto dtos.LocationDrainedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.LocationDrainedEventFromDTO(dto), nil

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

	case evKindIntelDecryptionStarted:
		var dto dtos.IntelDecryptionStartedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.IntelDecryptionStartedEventFromDTO(dto), nil

	case evKindIntelDecryptionFinished:
		var dto dtos.IntelDecryptionFinishedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.IntelDecryptionFinishedEventFromDTO(dto), nil

	case evKindDamagedItemRestorationStarted:
		var dto dtos.DamagedItemRestorationStartedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.DamagedItemRestorationStartedEventFromDTO(dto), nil

	case evKindDamagedItemRestored:
		var dto dtos.DamagedItemRestoredEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.DamagedItemRestoredEventFromDTO(dto), nil

	case evKindArtifactActivated:
		var dto dtos.ArtifactActivatedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.ArtifactActivatedEventFromDTO(dto), nil

	case evKindArtifactDeactivated:
		var dto dtos.ArtifactDeactivatedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.ArtifactDeactivatedEventFromDTO(dto), nil

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
	case evKindMilitaryOperationCancelled:
		var dto dtos.MilitaryOperationCancelledEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.MilitaryOperationCancelledEventFromDTO(dto), nil

	case evKindTradeOperationCreated:
		var dto dtos.TradeOperationCreatedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.TradeOperationCreatedEventFromDTO(dto), nil
	case evKindTradeOperationAccepted:
		var dto dtos.TradeOperationAcceptedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.TradeOperationAcceptedEventFromDTO(dto), nil
	case evKindTradeOperationDeclined:
		var dto dtos.TradeOperationDeclinedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.TradeOperationDeclinedEventFromDTO(dto), nil
	case evKindTradeOperationExpired:
		var dto dtos.TradeOperationExpiredEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.TradeOperationExpiredEventFromDTO(dto), nil
	case evKindTradeOperationCancelled:
		var dto dtos.TradeOperationCancelledByInitiatorEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.TradeOperationCancelledByInitiatorEventFromDTO(dto), nil
	case evKindTradeOperationOutbound:
		var dto dtos.TradeOperationOutboundEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.TradeOperationOutboundEventFromDTO(dto), nil
	case evKindTradeOperationArrived:
		var dto dtos.TradeOperationArrivedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.TradeOperationArrivedEventFromDTO(dto), nil
	case evKindTradeOperationReturning:
		var dto dtos.TradeOperationReturningEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.TradeOperationReturningEventFromDTO(dto), nil
	case evKindTradeOperationReturnArrived:
		var dto dtos.TradeOperationReturnArrivedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.TradeOperationReturnArrivedEventFromDTO(dto), nil

	case evKindScanReportCreated:
		var dto dtos.ScanReportCreatedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.ScanReportCreatedEventFromDTO(dto), nil

	case evKindRadarThreatDetected:
		var dto dtos.RadarThreatDetectedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.RadarThreatDetectedEventFromDTO(dto), nil

	case evKindActivityCreated:
		var dto dtos.ActivityCreatedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.ActivityCreatedEventFromDTO(dto), nil

	case evKindDiplomaticMessageSent:
		var dto dtos.DiplomaticMessageSentEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.DiplomaticMessageSentEventFromDTO(dto), nil

	case evKindDiplomaticRelationshipCreated:
		var dto dtos.DiplomaticRelationshipCreatedEventDTO
		if err := json.Unmarshal(payload, &dto); err != nil {
			return nil, err
		}
		return dtos.DiplomaticRelationshipCreatedEventFromDTO(dto), nil

	default:
		return nil, errors.New("unknown domain event kind")
	}
}
