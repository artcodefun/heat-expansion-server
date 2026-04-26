package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
)

func TradeOperationFromDB(r gen.TradeOperation) *domain.TradeOperation {
	var offeredDTO dtos.TradePayloadDTO
	_ = json.Unmarshal(r.OfferedPayload, &offeredDTO)

	var requestedDTO dtos.TradePayloadDTO
	_ = json.Unmarshal(r.RequestedPayload, &requestedDTO)

	var unitDTOs []dtos.MilitaryUnitDTO
	_ = json.Unmarshal(r.TransportUnits, &unitDTOs)
	transportUnits := make([]domain.MilitaryUnitSnap, 0, len(unitDTOs))
	for _, d := range unitDTOs {
		transportUnits = append(transportUnits, dtos.MilitaryUnitFromDTO(d))
	}

	var storageDTOs []dtos.StorageItemSnapDTO
	_ = json.Unmarshal(r.StorageSnaps, &storageDTOs)
	storageSnaps := make([]domain.StorageItemSnap, 0, len(storageDTOs))
	for _, d := range storageDTOs {
		storageSnaps = append(storageSnaps, dtos.StorageItemSnapFromDTO(d))
	}

	var modifiersDTO dtos.MilitaryModifiersDTO
	_ = json.Unmarshal(r.TotalModifiers, &modifiersDTO)

	return &domain.TradeOperation{
		ID:                int(r.ID),
		UUID:              r.OperationUuid,
		CreatedAt:         r.CreatedAt,
		SenderUserID:      r.SenderUserID,
		SenderBaseID:      int(r.SenderBaseID),
		ReceiverUserID:    r.ReceiverUserID,
		ReceiverBaseID:    int(r.ReceiverBaseID),
		SourceCoordinates: domain.Vector2i{X: int(r.SourceX), Y: int(r.SourceY)},
		TargetCoordinates: domain.Vector2i{X: int(r.TargetX), Y: int(r.TargetY)},
		OfferedPayload:    dtos.TradePayloadFromDTO(offeredDTO),
		RequestedPayload:  dtos.TradePayloadFromDTO(requestedDTO),
		TransportUnits:    transportUnits,
		StorageSnaps:      storageSnaps,
		TotalModifiers:    dtos.MilitaryModifiersFromDTO(modifiersDTO),
		ExpiresAt:         r.ExpiresAt,
		OutboundDepartAt:  r.OutboundDepartAt,
		OutboundArriveAt:  r.OutboundArriveAt,
		ArrivedAtTargetAt: r.ArrivedAtTargetAt,
		ReturnDepartAt:    r.ReturnDepartAt,
		ReturnArriveAt:    r.ReturnArriveAt,
		CompletedAt:       r.CompletedAt,
		Phase:             domain.TradeOperationPhase(r.Phase),
		Result:            domain.TradeOperationResult(r.Result),
		CrystalsSkipPrice: int(r.CrystalsSkipPrice),
	}
}

func InsertTradeOperationParamsFromDomain(op *domain.TradeOperation) gen.InsertTradeOperationParams {
	offeredJSON, _ := json.Marshal(dtos.TradePayloadDTOFromDomain(op.OfferedPayload))
	requestedJSON, _ := json.Marshal(dtos.TradePayloadDTOFromDomain(op.RequestedPayload))

	unitDTOs := make([]dtos.MilitaryUnitDTO, 0, len(op.TransportUnits))
	for _, u := range op.TransportUnits {
		unitDTOs = append(unitDTOs, dtos.MilitaryUnitDTOFromDomain(u))
	}
	unitsJSON, _ := json.Marshal(unitDTOs)

	storageDTOs := make([]dtos.StorageItemSnapDTO, 0, len(op.StorageSnaps))
	for _, s := range op.StorageSnaps {
		storageDTOs = append(storageDTOs, dtos.StorageItemSnapDTOFromDomain(s))
	}
	storageJSON, _ := json.Marshal(storageDTOs)

	modifiersJSON, _ := json.Marshal(dtos.MilitaryModifiersDTOFromDomain(op.TotalModifiers))

	return gen.InsertTradeOperationParams{
		OperationUuid:     op.UUID,
		CreatedAt:         op.CreatedAt,
		SenderUserID:      op.SenderUserID,
		SenderBaseID:      int64(op.SenderBaseID),
		ReceiverUserID:    op.ReceiverUserID,
		ReceiverBaseID:    int64(op.ReceiverBaseID),
		SourceX:           int32(op.SourceCoordinates.X),
		SourceY:           int32(op.SourceCoordinates.Y),
		TargetX:           int32(op.TargetCoordinates.X),
		TargetY:           int32(op.TargetCoordinates.Y),
		OfferedPayload:    offeredJSON,
		RequestedPayload:  requestedJSON,
		TransportUnits:    unitsJSON,
		StorageSnaps:      storageJSON,
		TotalModifiers:    modifiersJSON,
		ExpiresAt:         op.ExpiresAt,
		OutboundDepartAt:  op.OutboundDepartAt,
		OutboundArriveAt:  op.OutboundArriveAt,
		ArrivedAtTargetAt: op.ArrivedAtTargetAt,
		ReturnDepartAt:    op.ReturnDepartAt,
		ReturnArriveAt:    op.ReturnArriveAt,
		CompletedAt:       op.CompletedAt,
		Phase:             string(op.Phase),
		Result:            string(op.Result),
		CrystalsSkipPrice: int64(op.CrystalsSkipPrice),
	}
}

func UpdateTradeOperationParamsFromDomain(op *domain.TradeOperation) gen.UpdateTradeOperationParams {
	offeredJSON, _ := json.Marshal(dtos.TradePayloadDTOFromDomain(op.OfferedPayload))
	requestedJSON, _ := json.Marshal(dtos.TradePayloadDTOFromDomain(op.RequestedPayload))

	unitDTOs := make([]dtos.MilitaryUnitDTO, 0, len(op.TransportUnits))
	for _, u := range op.TransportUnits {
		unitDTOs = append(unitDTOs, dtos.MilitaryUnitDTOFromDomain(u))
	}
	unitsJSON, _ := json.Marshal(unitDTOs)

	storageDTOs := make([]dtos.StorageItemSnapDTO, 0, len(op.StorageSnaps))
	for _, s := range op.StorageSnaps {
		storageDTOs = append(storageDTOs, dtos.StorageItemSnapDTOFromDomain(s))
	}
	storageJSON, _ := json.Marshal(storageDTOs)

	modifiersJSON, _ := json.Marshal(dtos.MilitaryModifiersDTOFromDomain(op.TotalModifiers))

	return gen.UpdateTradeOperationParams{
		ID:                int64(op.ID),
		OperationUuid:     op.UUID,
		CreatedAt:         op.CreatedAt,
		SenderUserID:      op.SenderUserID,
		SenderBaseID:      int64(op.SenderBaseID),
		ReceiverUserID:    op.ReceiverUserID,
		ReceiverBaseID:    int64(op.ReceiverBaseID),
		SourceX:           int32(op.SourceCoordinates.X),
		SourceY:           int32(op.SourceCoordinates.Y),
		TargetX:           int32(op.TargetCoordinates.X),
		TargetY:           int32(op.TargetCoordinates.Y),
		OfferedPayload:    offeredJSON,
		RequestedPayload:  requestedJSON,
		TransportUnits:    unitsJSON,
		StorageSnaps:      storageJSON,
		TotalModifiers:    modifiersJSON,
		ExpiresAt:         op.ExpiresAt,
		OutboundDepartAt:  op.OutboundDepartAt,
		OutboundArriveAt:  op.OutboundArriveAt,
		ArrivedAtTargetAt: op.ArrivedAtTargetAt,
		ReturnDepartAt:    op.ReturnDepartAt,
		ReturnArriveAt:    op.ReturnArriveAt,
		CompletedAt:       op.CompletedAt,
		Phase:             string(op.Phase),
		Result:            string(op.Result),
		CrystalsSkipPrice: int64(op.CrystalsSkipPrice),
	}
}
