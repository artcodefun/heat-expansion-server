package domain

import (
	"github.com/google/uuid"
)

// EventProducer records domain events for aggregates/entities.
type EventProducer struct {
	events []DomainEvent
}

func (ep *EventProducer) AddEvent(event DomainEvent) {
	ep.events = append(ep.events, event)
}

func (ep *EventProducer) PullEvents() []DomainEvent {
	events := ep.events
	ep.events = nil
	return events
}

// DomainEvent is the interface for all domain events.
// We keep only the occurrence time; consumers can type-switch on concrete event types.
type DomainEvent interface {
	OccurredAt() int64
}

// BasicEvent carries the timestamp for all domain events.
type BasicEvent struct {
	timestamp int64
}

func (e BasicEvent) OccurredAt() int64 { return e.timestamp }

func NewBasicEvent() BasicEvent {
	return BasicEvent{timestamp: NowUnix()}
}

// User account creation event
type UserAccountCreatedEvent struct {
	BasicEvent
	UserID int
}

func NewUserAccountCreatedEvent(userID int) UserAccountCreatedEvent {
	return UserAccountCreatedEvent{
		BasicEvent: NewBasicEvent(),
		UserID:     userID,
	}
}

// Building-related domain events

type BuildingProductionStartedEvent struct {
	BasicEvent
	BaseID         int
	ItemID         uuid.UUID
	CompletionDate int64
}

func NewBuildingProductionStartedEvent(baseID int, itemID uuid.UUID, completionDate int64) BuildingProductionStartedEvent {
	return BuildingProductionStartedEvent{
		BasicEvent:     NewBasicEvent(),
		BaseID:         baseID,
		ItemID:         itemID,
		CompletionDate: completionDate,
	}
}

type BuildingProductionFinishedEvent struct {
	BasicEvent
	BaseID        int
	ItemID        uuid.UUID // ID of the in-production item that just finished
	PresentItemID uuid.UUID // ID of the newly created present building instance
}

func NewBuildingProductionFinishedEvent(baseID int, itemID uuid.UUID, presentItemID uuid.UUID) BuildingProductionFinishedEvent {
	return BuildingProductionFinishedEvent{
		BasicEvent:    NewBasicEvent(),
		BaseID:        baseID,
		ItemID:        itemID,
		PresentItemID: presentItemID,
	}
}

type BuildingProductionCancelledEvent struct {
	BasicEvent
	BaseID int
	ItemID uuid.UUID
}

func NewBuildingProductionCancelledEvent(baseID int, itemID uuid.UUID) BuildingProductionCancelledEvent {
	return BuildingProductionCancelledEvent{
		BasicEvent: NewBasicEvent(),
		BaseID:     baseID,
		ItemID:     itemID,
	}
}

type BuildingProductionSpeedupEvent struct {
	BasicEvent
	BaseID int
	ItemID uuid.UUID
}

func NewBuildingProductionSpeedupEvent(baseID int, itemID uuid.UUID) BuildingProductionSpeedupEvent {
	return BuildingProductionSpeedupEvent{
		BasicEvent: NewBasicEvent(),
		BaseID:     baseID,
		ItemID:     itemID,
	}
}

type BuildingPresentDeletedEvent struct {
	BasicEvent
	BaseID int
	ItemID uuid.UUID
}

func NewBuildingPresentDeletedEvent(baseID int, itemID uuid.UUID) BuildingPresentDeletedEvent {
	return BuildingPresentDeletedEvent{
		BasicEvent: NewBasicEvent(),
		BaseID:     baseID,
		ItemID:     itemID,
	}
}

// Army-related domain events
type ArmyProductionPendingEvent struct {
	BasicEvent
	BaseID int
	ItemID uuid.UUID
	Count  int
}

func NewArmyProductionPendingEvent(baseID int, itemID uuid.UUID, count int) ArmyProductionPendingEvent {
	return ArmyProductionPendingEvent{
		BasicEvent: NewBasicEvent(),
		BaseID:     baseID,
		ItemID:     itemID,
		Count:      count,
	}
}

type ArmyProductionStartedEvent struct {
	BasicEvent
	BaseID         int
	ItemID         uuid.UUID
	CompletionDate int64
}

func NewArmyProductionStartedEvent(baseID int, itemID uuid.UUID, completionDate int64) ArmyProductionStartedEvent {
	return ArmyProductionStartedEvent{
		BasicEvent:     NewBasicEvent(),
		BaseID:         baseID,
		ItemID:         itemID,
		CompletionDate: completionDate,
	}
}

type ArmyProductionFinishedEvent struct {
	BasicEvent
	BaseID int
	ItemID uuid.UUID
}

func NewArmyProductionFinishedEvent(baseID int, itemID uuid.UUID) ArmyProductionFinishedEvent {
	return ArmyProductionFinishedEvent{
		BasicEvent: NewBasicEvent(),
		BaseID:     baseID,
		ItemID:     itemID,
	}
}

type ArmyProductionCancelledEvent struct {
	BasicEvent
	BaseID int
	ItemID uuid.UUID
	Count  int
}

func NewArmyProductionCancelledEvent(baseID int, itemID uuid.UUID, count int) ArmyProductionCancelledEvent {
	return ArmyProductionCancelledEvent{
		BasicEvent: NewBasicEvent(),
		BaseID:     baseID,
		ItemID:     itemID,
		Count:      count,
	}
}

type ArmyProductionSpeedupEvent struct {
	BasicEvent
	BaseID int
	ItemID uuid.UUID
}

func NewArmyProductionSpeedupEvent(baseID int, itemID uuid.UUID) ArmyProductionSpeedupEvent {
	return ArmyProductionSpeedupEvent{
		BasicEvent: NewBasicEvent(),
		BaseID:     baseID,
		ItemID:     itemID,
	}
}

type ArmyPresentDeletedEvent struct {
	BasicEvent
	BaseID int
	ItemID uuid.UUID
	Count  int
}

func NewArmyPresentDeletedEvent(baseID int, itemID uuid.UUID, count int) ArmyPresentDeletedEvent {
	return ArmyPresentDeletedEvent{
		BasicEvent: NewBasicEvent(),
		BaseID:     baseID,
		ItemID:     itemID,
		Count:      count,
	}
}

// Technology-related domain events
type TechResearchStartedEvent struct {
	BasicEvent
	BaseID         int
	ItemID         uuid.UUID
	PrototypeID    int
	CompletionDate int64
}

func NewTechResearchStartedEvent(baseID int, itemID uuid.UUID, prototypeID int, completionDate int64) TechResearchStartedEvent {
	return TechResearchStartedEvent{
		BasicEvent:     NewBasicEvent(),
		BaseID:         baseID,
		ItemID:         itemID,
		PrototypeID:    prototypeID,
		CompletionDate: completionDate,
	}
}

type TechResearchFinishedEvent struct {
	BasicEvent
	BaseID      int
	ItemID      uuid.UUID
	PrototypeID int
}

func NewTechResearchFinishedEvent(baseID int, itemID uuid.UUID, prototypeID int) TechResearchFinishedEvent {
	return TechResearchFinishedEvent{
		BasicEvent:  NewBasicEvent(),
		BaseID:      baseID,
		ItemID:      itemID,
		PrototypeID: prototypeID,
	}
}

type TechResearchSpeedupEvent struct {
	BasicEvent
	BaseID      int
	ItemID      uuid.UUID
	PrototypeID int
}

func NewTechResearchSpeedupEvent(baseID int, itemID uuid.UUID, prototypeID int) TechResearchSpeedupEvent {
	return TechResearchSpeedupEvent{
		BasicEvent:  NewBasicEvent(),
		BaseID:      baseID,
		ItemID:      itemID,
		PrototypeID: prototypeID,
	}
}

// Storage-related domain events
type StorageItemPresentDeletedEvent struct {
	BasicEvent
	BaseID int
	ItemID uuid.UUID
}

func NewStorageItemPresentDeletedEvent(baseID int, itemID uuid.UUID) StorageItemPresentDeletedEvent {
	return StorageItemPresentDeletedEvent{
		BasicEvent: NewBasicEvent(),
		BaseID:     baseID,
		ItemID:     itemID,
	}
}

type BuffActivatedEvent struct {
	BasicEvent
	BaseID int
	ItemID uuid.UUID
}

func NewBuffActivatedEvent(baseID int, itemID uuid.UUID) BuffActivatedEvent {
	return BuffActivatedEvent{
		BasicEvent: NewBasicEvent(),
		BaseID:     baseID,
		ItemID:     itemID,
	}
}

type IntelDecryptionStartedEvent struct {
	BasicEvent
	BaseID int
	ItemID uuid.UUID
}

func NewIntelDecryptionStartedEvent(baseID int, itemID uuid.UUID) IntelDecryptionStartedEvent {
	return IntelDecryptionStartedEvent{
		BasicEvent: NewBasicEvent(),
		BaseID:     baseID,
		ItemID:     itemID,
	}
}

type IntelDecryptionFinishedEvent struct {
	BasicEvent
	BaseID    int
	ItemID    uuid.UUID
	IntelType HiddenLocationType
}

func NewIntelDecryptionFinishedEvent(baseID int, itemID uuid.UUID, intelType HiddenLocationType) IntelDecryptionFinishedEvent {
	return IntelDecryptionFinishedEvent{
		BasicEvent: NewBasicEvent(),
		BaseID:     baseID,
		ItemID:     itemID,
		IntelType:  intelType,
	}
}

type DamagedItemRestorationStartedEvent struct {
	BasicEvent
	BaseID int
	ItemID uuid.UUID
}

func NewDamagedItemRestorationStartedEvent(baseID int, itemID uuid.UUID) DamagedItemRestorationStartedEvent {
	return DamagedItemRestorationStartedEvent{
		BasicEvent: NewBasicEvent(),
		BaseID:     baseID,
		ItemID:     itemID,
	}
}

type DamagedItemRestoredEvent struct {
	BasicEvent
	BaseID int
	ItemID uuid.UUID
}

func NewDamagedItemRestoredEvent(baseID int, itemID uuid.UUID) DamagedItemRestoredEvent {
	return DamagedItemRestoredEvent{
		BasicEvent: NewBasicEvent(),
		BaseID:     baseID,
		ItemID:     itemID,
	}
}

type ArtifactActivatedEvent struct {
	BasicEvent
	BaseID int
	ItemID uuid.UUID
}

func NewArtifactActivatedEvent(baseID int, itemID uuid.UUID) ArtifactActivatedEvent {
	return ArtifactActivatedEvent{
		BasicEvent: NewBasicEvent(),
		BaseID:     baseID,
		ItemID:     itemID,
	}
}

type ArtifactDeactivatedEvent struct {
	BasicEvent
	BaseID int
	ItemID uuid.UUID
}

func NewArtifactDeactivatedEvent(baseID int, itemID uuid.UUID) ArtifactDeactivatedEvent {
	return ArtifactDeactivatedEvent{
		BasicEvent: NewBasicEvent(),
		BaseID:     baseID,
		ItemID:     itemID,
	}
}

// Military operation-related domain events

type MilitaryOperationStartedEvent struct {
	BasicEvent
	OperationID      int
	OutboundArriveAt int64
}

func NewMilitaryOperationStartedEvent(operationID int, outboundArriveAt int64) MilitaryOperationStartedEvent {
	return MilitaryOperationStartedEvent{
		BasicEvent:       NewBasicEvent(),
		OperationID:      operationID,
		OutboundArriveAt: outboundArriveAt,
	}
}

type MilitaryOperationArrivedEvent struct {
	BasicEvent
	OperationID int
}

func NewMilitaryOperationArrivedEvent(operationID int) MilitaryOperationArrivedEvent {
	return MilitaryOperationArrivedEvent{
		BasicEvent:  NewBasicEvent(),
		OperationID: operationID,
	}
}

type MilitaryOperationResolvedEvent struct {
	BasicEvent
	OperationID int
	Result      MilitaryOperationResult
}

func NewMilitaryOperationResolvedEvent(operationID int, result MilitaryOperationResult) MilitaryOperationResolvedEvent {
	return MilitaryOperationResolvedEvent{
		BasicEvent:  NewBasicEvent(),
		OperationID: operationID,
		Result:      result,
	}
}

type MilitaryOperationReturnStartedEvent struct {
	BasicEvent
	OperationID    int
	ReturnArriveAt int64
}

func NewMilitaryOperationReturnStartedEvent(operationID int, returnArriveAt int64) MilitaryOperationReturnStartedEvent {
	return MilitaryOperationReturnStartedEvent{
		BasicEvent:     NewBasicEvent(),
		OperationID:    operationID,
		ReturnArriveAt: returnArriveAt,
	}
}

type MilitaryOperationReturnArrivedEvent struct {
	BasicEvent
	OperationID int
}

func NewMilitaryOperationReturnArrivedEvent(operationID int) MilitaryOperationReturnArrivedEvent {
	return MilitaryOperationReturnArrivedEvent{
		BasicEvent:  NewBasicEvent(),
		OperationID: operationID,
	}
}

type MilitaryOperationCancelledEvent struct {
	BasicEvent
	OperationID int
}

func NewMilitaryOperationCancelledEvent(operationID int) MilitaryOperationCancelledEvent {
	return MilitaryOperationCancelledEvent{
		BasicEvent:  NewBasicEvent(),
		OperationID: operationID,
	}
}

// Scan report-related domain events

// ScanReportCreatedEvent is emitted when a SectorScanReport has been persisted and is visible to the domain.
// This signals read models (e.g., Activities) to project a SCAN entry.
type ScanReportCreatedEvent struct {
	BasicEvent
	ReportID          int
	BaseID            int // source base that initiated the scan-producing operation
	SourceOperationID int // optional link back to the operation
}

func NewScanReportCreatedEvent(reportID int, baseID int, sourceOperationID int) ScanReportCreatedEvent {
	return ScanReportCreatedEvent{
		BasicEvent:        NewBasicEvent(),
		ReportID:          reportID,
		BaseID:            baseID,
		SourceOperationID: sourceOperationID,
	}
}

// Radar-related domain events

type RadarThreatDetectedEvent struct {
	BasicEvent
	RadarThreatID uuid.UUID
	OwnerBaseID   int
	OperationID   int
}

func NewRadarThreatDetectedEvent(threatID uuid.UUID, baseID int, opID int) RadarThreatDetectedEvent {
	return RadarThreatDetectedEvent{
		BasicEvent:    NewBasicEvent(),
		RadarThreatID: threatID,
		OwnerBaseID:   baseID,
		OperationID:   opID,
	}
}
