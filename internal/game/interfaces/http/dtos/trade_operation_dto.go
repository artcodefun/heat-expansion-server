package dtos

import (
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/google/uuid"
)

type TradeOperationPhase string

const (
	TradeOperationPhasePending   TradeOperationPhase = "PENDING"
	TradeOperationPhaseOutbound  TradeOperationPhase = "OUTBOUND"
	TradeOperationPhaseArrived   TradeOperationPhase = "ARRIVED"
	TradeOperationPhaseReturning TradeOperationPhase = "RETURNING"
	TradeOperationPhaseCompleted TradeOperationPhase = "COMPLETED"
)

type TradeOperationResult string

const (
	TradeOperationResultUnknown  TradeOperationResult = "UNKNOWN"
	TradeOperationResultSuccess  TradeOperationResult = "SUCCESS"
	TradeOperationResultCanceled TradeOperationResult = "CANCELED"
	TradeOperationResultDeclined TradeOperationResult = "DECLINED"
	TradeOperationResultExpired  TradeOperationResult = "EXPIRED"
	TradeOperationResultFailure  TradeOperationResult = "FAILURE"
)

type TradeOperationDTO struct {
	ID                   int                  `json:"id"`
	UUID                 uuid.UUID            `json:"uuid"`
	CreatedAt            int64                `json:"created_at"`
	SenderUserID         uuid.UUID            `json:"sender_user_id"`
	SenderBaseID         int                  `json:"sender_base_id"`
	ReceiverUserID       uuid.UUID            `json:"receiver_user_id"`
	ReceiverBaseID       int                  `json:"receiver_base_id"`
	SourceCoordinates    Vector2iDTO          `json:"source_coordinates"`
	TargetCoordinates    Vector2iDTO          `json:"target_coordinates"`
	OfferedPayload       TradePayloadDTO      `json:"offered_payload"`
	RequestedPayload     TradePayloadDTO      `json:"requested_payload"`
	TransportUnits       []MilitaryUnitDTO    `json:"transport_units"`
	StorageSnaps         []StorageItemSnapDTO `json:"storage_snaps"`
	TotalModifiers       MilitaryModifiersDTO `json:"total_modifiers"`
	ExpiresAt            int64                `json:"expires_at"`
	OutboundDepartAt     int64                `json:"outbound_depart_at"`
	OutboundArriveAt     int64                `json:"outbound_arrive_at"`
	ArrivedAtTargetAt    int64                `json:"arrived_at_target_at"`
	ReturnDepartAt       int64                `json:"return_depart_at"`
	ReturnArriveAt       int64                `json:"return_arrive_at"`
	CompletedAt          int64                `json:"completed_at"`
	CrystalsSkipPrice    int                  `json:"crystals_skip_price"`
	Phase                TradeOperationPhase  `json:"phase"`
	Result               TradeOperationResult `json:"result"`
	SenderScanOfReceiver *SectorDTO           `json:"sender_scan_of_receiver,omitempty"`
	ReceiverScanOfSender *SectorDTO           `json:"receiver_scan_of_sender,omitempty"`
}

type TradeArmyItemSnapDTO struct {
	PrototypeID      int                  `json:"prototype_id"`
	CurrentPrototype ArmyItemPrototypeDTO `json:"current_prototype"`
	Count            int                  `json:"count"`
	Capacity         int                  `json:"capacity"`
}

type TradeStorageItemSnapDTO struct {
	ItemID           uuid.UUID               `json:"item_id"`
	PrototypeID      int                     `json:"prototype_id"`
	CurrentPrototype StorageItemPrototypeDTO `json:"current_prototype"`
	Category         StorageCategory         `json:"category"`
}

type TradePayloadDTO struct {
	Resources PriceModelDTO             `json:"resources"`
	Storage   []TradeStorageItemSnapDTO `json:"storage"`
	Army      []TradeArmyItemSnapDTO    `json:"army"`
}

type TradeInfoDTO struct {
	Resources PriceModelDTO           `json:"resources"`
	Army      []ArmyItemPresentDTO    `json:"army"`
	Storage   []StorageItemPresentDTO `json:"storage"`
}

func TradeInfoFromReadModel(m *readmodels.TradeInfo, tr ports.Translator, locale string) TradeInfoDTO {
	return TradeInfoDTO{
		Resources: PriceModelFromReadModel(m.Resources),
		Army:      ArmyItemsPresentFromReadModels(m.Army, tr, locale),
		Storage:   StorageItemsPresentFromReadModels(m.Storage, tr, locale),
	}
}

func TradeOperationFromReadModel(m *readmodels.TradeOperation, tr ports.Translator, locale string) TradeOperationDTO {
	offeredStorage := tradeStorageSnapsFromReadModel(m.OfferedPayload.Storage, tr, locale)
	offeredArmy := tradeArmySnapsFromReadModel(m.OfferedPayload.Army, tr, locale)
	requestedStorage := tradeStorageSnapsFromReadModel(m.RequestedPayload.Storage, tr, locale)
	requestedArmy := tradeArmySnapsFromReadModel(m.RequestedPayload.Army, tr, locale)
	dto := TradeOperationDTO{
		ID:                m.ID,
		UUID:              m.UUID,
		CreatedAt:         m.CreatedAt,
		SenderUserID:      m.SenderUserID,
		SenderBaseID:      m.SenderBaseID,
		ReceiverUserID:    m.ReceiverUserID,
		ReceiverBaseID:    m.ReceiverBaseID,
		SourceCoordinates: Vector2iFromReadModel(m.SourceCoordinates),
		TargetCoordinates: Vector2iFromReadModel(m.TargetCoordinates),
		OfferedPayload: TradePayloadDTO{
			Resources: PriceModelFromReadModel(m.OfferedPayload.Resources),
			Storage:   offeredStorage,
			Army:      offeredArmy,
		},
		RequestedPayload: TradePayloadDTO{
			Resources: PriceModelFromReadModel(m.RequestedPayload.Resources),
			Storage:   requestedStorage,
			Army:      requestedArmy,
		},
		TransportUnits:    MilitaryUnitsFromReadModel(m.TransportUnits, tr, locale),
		StorageSnaps:      storageItemSnapsFromReadModel(m.StorageSnaps, tr, locale),
		TotalModifiers:    MilitaryModifiersFromReadModel(m.TotalModifiers),
		ExpiresAt:         m.ExpiresAt,
		OutboundDepartAt:  m.OutboundDepartAt,
		OutboundArriveAt:  m.OutboundArriveAt,
		ArrivedAtTargetAt: m.ArrivedAtTargetAt,
		ReturnDepartAt:    m.ReturnDepartAt,
		ReturnArriveAt:    m.ReturnArriveAt,
		CompletedAt:       m.CompletedAt,
		CrystalsSkipPrice: m.CrystalsSkipPrice,
		Phase:             TradeOperationPhase(m.Phase),
		Result:            TradeOperationResult(m.Result),
	}
	if m.SenderScanOfReceiver != nil {
		report := SectorScanReportFromReadModel(m.SenderScanOfReceiver, tr, locale)
		dto.SenderScanOfReceiver = &report
	}
	if m.ReceiverScanOfSender != nil {
		report := SectorScanReportFromReadModel(m.ReceiverScanOfSender, tr, locale)
		dto.ReceiverScanOfSender = &report
	}
	return dto
}

func TradeOperationsFromReadModels(items []*readmodels.TradeOperation, tr ports.Translator, locale string) []TradeOperationDTO {
	out := make([]TradeOperationDTO, 0, len(items))
	for _, item := range items {
		out = append(out, TradeOperationFromReadModel(item, tr, locale))
	}
	return out
}

func tradeArmySnapsFromReadModel(snaps []readmodels.TradeArmyItemSnap, tr ports.Translator, locale string) []TradeArmyItemSnapDTO {
	out := make([]TradeArmyItemSnapDTO, 0, len(snaps))
	for _, s := range snaps {
		out = append(out, TradeArmyItemSnapDTO{
			PrototypeID:      s.PrototypeID,
			CurrentPrototype: mapArmyPrototype(s.CurrentPrototype, tr, locale),
			Count:            s.Count,
			Capacity:         s.Capacity,
		})
	}
	return out
}

func tradeStorageSnapsFromReadModel(snaps []readmodels.TradeStorageItemSnap, tr ports.Translator, locale string) []TradeStorageItemSnapDTO {
	out := make([]TradeStorageItemSnapDTO, 0, len(snaps))
	for _, s := range snaps {
		out = append(out, TradeStorageItemSnapDTO{
			ItemID:           s.ItemID,
			PrototypeID:      s.PrototypeID,
			CurrentPrototype: mapStorageItemPrototype(s.CurrentPrototype, tr, locale),
			Category:         StorageCategory(s.Category),
		})
	}
	return out
}
