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
	Timestamp int64
}

func (e BasicEvent) OccurredAt() int64 { return e.Timestamp }

func NewBasicEvent() BasicEvent {
	return BasicEvent{Timestamp: NowUnix()}
}

// User account creation event
type UserAccountCreatedEvent struct {
	BasicEvent
	UserID uuid.UUID
}

func NewUserAccountCreatedEvent(userID uuid.UUID) UserAccountCreatedEvent {
	return UserAccountCreatedEvent{
		BasicEvent: NewBasicEvent(),
		UserID:     userID,
	}
}

// User base creation event
type UserBaseCreatedEvent struct {
	BasicEvent
	BaseID  int
	OwnerID uuid.UUID
}

func NewUserBaseCreatedEvent(baseID int, ownerID uuid.UUID) UserBaseCreatedEvent {
	return UserBaseCreatedEvent{
		BasicEvent: NewBasicEvent(),
		BaseID:     baseID,
		OwnerID:    ownerID,
	}
}

// LocationDrainedEvent is emitted when a location has no resources, trophies, or defenders left.
type LocationDrainedEvent struct {
	BasicEvent
	X    int
	Y    int
	Type LocationType // "RESOURCEFUL" or "DANGEROUS"
}

func NewLocationDrainedEvent(x, y int, locType LocationType) LocationDrainedEvent {
	return LocationDrainedEvent{
		BasicEvent: NewBasicEvent(),
		X:          x,
		Y:          y,
		Type:       locType,
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

// Trade operation-related domain events

type TradeOperationCreatedEvent struct {
	BasicEvent
	OperationID     int
	OperationUUID   uuid.UUID
	SenderBaseID    int
	ReceiverBaseID  int
	ExpirationAtSec int64
}

func NewTradeOperationCreatedEvent(operationID int, operationUUID uuid.UUID, senderBaseID int, receiverBaseID int, expirationAtSec int64) TradeOperationCreatedEvent {
	return TradeOperationCreatedEvent{
		BasicEvent:      NewBasicEvent(),
		OperationID:     operationID,
		OperationUUID:   operationUUID,
		SenderBaseID:    senderBaseID,
		ReceiverBaseID:  receiverBaseID,
		ExpirationAtSec: expirationAtSec,
	}
}

type TradeOperationAcceptedEvent struct {
	BasicEvent
	OperationID      int
	OutboundArriveAt int64
}

func NewTradeOperationAcceptedEvent(operationID int, outboundArriveAt int64) TradeOperationAcceptedEvent {
	return TradeOperationAcceptedEvent{
		BasicEvent:       NewBasicEvent(),
		OperationID:      operationID,
		OutboundArriveAt: outboundArriveAt,
	}
}

type TradeOperationDeclinedEvent struct {
	BasicEvent
	OperationID int
}

func NewTradeOperationDeclinedEvent(operationID int) TradeOperationDeclinedEvent {
	return TradeOperationDeclinedEvent{
		BasicEvent:  NewBasicEvent(),
		OperationID: operationID,
	}
}

type TradeOperationExpiredEvent struct {
	BasicEvent
	OperationID int
}

func NewTradeOperationExpiredEvent(operationID int) TradeOperationExpiredEvent {
	return TradeOperationExpiredEvent{
		BasicEvent:  NewBasicEvent(),
		OperationID: operationID,
	}
}

type TradeOperationCancelledByInitiatorEvent struct {
	BasicEvent
	OperationID int
}

func NewTradeOperationCancelledByInitiatorEvent(operationID int) TradeOperationCancelledByInitiatorEvent {
	return TradeOperationCancelledByInitiatorEvent{
		BasicEvent:  NewBasicEvent(),
		OperationID: operationID,
	}
}

type TradeOperationOutboundEvent struct {
	BasicEvent
	OperationID      int
	OutboundArriveAt int64
}

func NewTradeOperationOutboundEvent(operationID int, outboundArriveAt int64) TradeOperationOutboundEvent {
	return TradeOperationOutboundEvent{
		BasicEvent:       NewBasicEvent(),
		OperationID:      operationID,
		OutboundArriveAt: outboundArriveAt,
	}
}

type TradeOperationArrivedEvent struct {
	BasicEvent
	OperationID int
}

func NewTradeOperationArrivedEvent(operationID int) TradeOperationArrivedEvent {
	return TradeOperationArrivedEvent{
		BasicEvent:  NewBasicEvent(),
		OperationID: operationID,
	}
}

type TradeOperationReturningEvent struct {
	BasicEvent
	OperationID    int
	ReturnArriveAt int64
}

func NewTradeOperationReturningEvent(operationID int, returnArriveAt int64) TradeOperationReturningEvent {
	return TradeOperationReturningEvent{
		BasicEvent:     NewBasicEvent(),
		OperationID:    operationID,
		ReturnArriveAt: returnArriveAt,
	}
}

// TradeOperationReturnArrivedEvent is emitted whenever trade operation assets
// can be safely released back to the involved bases, regardless of completion reason.
type TradeOperationReturnArrivedEvent struct {
	BasicEvent
	OperationID int
}

func NewTradeOperationReturnArrivedEvent(operationID int) TradeOperationReturnArrivedEvent {
	return TradeOperationReturnArrivedEvent{
		BasicEvent:  NewBasicEvent(),
		OperationID: operationID,
	}
}

// Scan report-related domain events

// ScanReportCreatedEvent is emitted when a SectorScanReport has been persisted and is visible to the domain.
// This signals read models (e.g., Activities) to project a SCAN entry.
type ScanReportCreatedEvent struct {
	BasicEvent
	ReportID   int
	BaseID     int // source base that initiated the scan-producing operation
	SourceType ScanReportSourceType
	SourceID   *uuid.UUID
}

func NewScanReportCreatedEvent(reportID int, baseID int, sourceType ScanReportSourceType, sourceID *uuid.UUID) ScanReportCreatedEvent {
	return ScanReportCreatedEvent{
		BasicEvent: NewBasicEvent(),
		ReportID:   reportID,
		BaseID:     baseID,
		SourceType: sourceType,
		SourceID:   sourceID,
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

// Activity-related domain events

type ActivityCreatedEvent struct {
	BasicEvent
	ActivityID uuid.UUID
	UserID     uuid.UUID
	BaseID     int
	Kind       ActivityKind
	Subtype    string
}

func NewActivityCreatedEvent(activityID uuid.UUID, userID uuid.UUID, baseID int, kind ActivityKind, subtype string) ActivityCreatedEvent {
	return ActivityCreatedEvent{
		BasicEvent: NewBasicEvent(),
		ActivityID: activityID,
		UserID:     userID,
		BaseID:     baseID,
		Kind:       kind,
		Subtype:    subtype,
	}
}

// Diplomacy-related domain events

type DiplomaticMessageSentEvent struct {
	BasicEvent
	MessageID      uuid.UUID
	SenderUserID   uuid.UUID
	ReceiverUserID uuid.UUID
	ReceiverBaseID *int
	Content        TranslationKey
}

func NewDiplomaticMessageSentEvent(messageID uuid.UUID, senderUserID, receiverUserID uuid.UUID, receiverBaseID *int, content TranslationKey) DiplomaticMessageSentEvent {
	return DiplomaticMessageSentEvent{
		BasicEvent:     NewBasicEvent(),
		MessageID:      messageID,
		SenderUserID:   senderUserID,
		ReceiverUserID: receiverUserID,
		ReceiverBaseID: receiverBaseID,
		Content:        content,
	}
}

type DiplomaticRequestCreatedEvent struct {
	BasicEvent
	RequestID      uuid.UUID
	SenderUserID   uuid.UUID
	ReceiverUserID uuid.UUID
	ReceiverBaseID *int
	Kind           DiplomaticRequestKind
}

func NewDiplomaticRequestCreatedEvent(requestID uuid.UUID, senderUserID, receiverUserID uuid.UUID, receiverBaseID *int, kind DiplomaticRequestKind) DiplomaticRequestCreatedEvent {
	return DiplomaticRequestCreatedEvent{
		BasicEvent:     NewBasicEvent(),
		RequestID:      requestID,
		SenderUserID:   senderUserID,
		ReceiverUserID: receiverUserID,
		ReceiverBaseID: receiverBaseID,
		Kind:           kind,
	}
}

// DiplomaticRelationshipCreatedEvent is emitted when a previously unknown diplomatic relationship becomes known.
// This lets projections create first-contact scan reports for bases that were not yet revealed.
type DiplomaticRelationshipCreatedEvent struct {
	BasicEvent
	RelationshipID  uuid.UUID
	UserAID         uuid.UUID
	UserBID         uuid.UUID
	Status          DiplomaticStatus
	ChangedByUserID uuid.UUID
}

func NewDiplomaticRelationshipCreatedEvent(relationshipID, userAID, userBID uuid.UUID, status DiplomaticStatus, changedByUserID uuid.UUID) DiplomaticRelationshipCreatedEvent {
	return DiplomaticRelationshipCreatedEvent{
		BasicEvent:      NewBasicEvent(),
		RelationshipID:  relationshipID,
		UserAID:         userAID,
		UserBID:         userBID,
		Status:          status,
		ChangedByUserID: changedByUserID,
	}
}
