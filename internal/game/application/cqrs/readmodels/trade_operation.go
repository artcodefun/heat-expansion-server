package readmodels

import "github.com/google/uuid"

type TradeOperationPhase string

const (
	TradePhasePending   TradeOperationPhase = "PENDING"
	TradePhaseOutbound  TradeOperationPhase = "OUTBOUND"
	TradePhaseArrived   TradeOperationPhase = "ARRIVED"
	TradePhaseReturning TradeOperationPhase = "RETURNING"
	TradePhaseCompleted TradeOperationPhase = "COMPLETED"
)

type TradeOperationResult string

const (
	TradeResultUnknown  TradeOperationResult = "UNKNOWN"
	TradeResultSuccess  TradeOperationResult = "SUCCESS"
	TradeResultCanceled TradeOperationResult = "CANCELED"
	TradeResultDeclined TradeOperationResult = "DECLINED"
	TradeResultExpired  TradeOperationResult = "EXPIRED"
	TradeResultFailure  TradeOperationResult = "FAILURE"
)

type TradeArmyItemSnap struct {
	PrototypeID      int
	CurrentPrototype ArmyItemPrototype
	Count            int
	Capacity         int
}

type TradeStorageItemSnap struct {
	ItemID           uuid.UUID
	PrototypeID      int
	CurrentPrototype StorageItemPrototype
	Category         StorageCategory
}

type TradePayload struct {
	Resources PriceModel
	Storage   []TradeStorageItemSnap
	Army      []TradeArmyItemSnap
}

type TradeInfo struct {
	Resources PriceModel
	Army      []*ArmyItemPresent
	Storage   []*StorageItemPresent
}

type TradeOperation struct {
	ID                   int
	UUID                 uuid.UUID
	CreatedAt            int64
	SenderUserID         uuid.UUID
	SenderBaseID         int
	ReceiverUserID       uuid.UUID
	ReceiverBaseID       int
	SourceCoordinates    Vector2i
	TargetCoordinates    Vector2i
	OfferedPayload       TradePayload
	RequestedPayload     TradePayload
	TransportUnits       []MilitaryUnitSnap
	StorageSnaps         []StorageItemSnap
	TotalModifiers       MilitaryModifiers
	ExpiresAt            int64
	OutboundDepartAt     int64
	OutboundArriveAt     int64
	ArrivedAtTargetAt    int64
	ReturnDepartAt       int64
	ReturnArriveAt       int64
	CompletedAt          int64
	CrystalsSkipPrice    int
	Phase                TradeOperationPhase
	Result               TradeOperationResult
	SenderScanOfReceiver *SectorScanReport
	ReceiverScanOfSender *SectorScanReport
}
