package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
)

func TradeOperationFromModel(m gen.TradeOperation) readmodels.TradeOperation {
	return readmodels.TradeOperation{
		ID:                int(m.ID),
		UUID:              m.OperationUuid,
		CreatedAt:         m.CreatedAt,
		SenderUserID:      m.SenderUserID,
		SenderBaseID:      int(m.SenderBaseID),
		ReceiverUserID:    m.ReceiverUserID,
		ReceiverBaseID:    int(m.ReceiverBaseID),
		SourceCoordinates: readmodels.Vector2i{X: int(m.SourceX), Y: int(m.SourceY)},
		TargetCoordinates: readmodels.Vector2i{X: int(m.TargetX), Y: int(m.TargetY)},
		OfferedPayload:    tradePayloadFromJSON(m.OfferedPayload),
		RequestedPayload:  tradePayloadFromJSON(m.RequestedPayload),
		TransportUnits:    tradeUnitsFromJSON(m.TransportUnits),
		StorageSnaps:      tradeOperationStorageSnapsFromJSON(m.StorageSnaps),
		TotalModifiers:    tradeModifiersFromJSON(m.TotalModifiers),
		ExpiresAt:         m.ExpiresAt,
		OutboundDepartAt:  m.OutboundDepartAt,
		OutboundArriveAt:  m.OutboundArriveAt,
		ArrivedAtTargetAt: m.ArrivedAtTargetAt,
		ReturnDepartAt:    m.ReturnDepartAt,
		ReturnArriveAt:    m.ReturnArriveAt,
		CompletedAt:       m.CompletedAt,
		Phase:             readmodels.TradeOperationPhase(m.Phase),
		Result:            readmodels.TradeOperationResult(m.Result),
		CrystalsSkipPrice: int(m.CrystalsSkipPrice),
	}
}

func tradePayloadFromJSON(raw json.RawMessage) readmodels.TradePayload {
	if len(raw) == 0 {
		return readmodels.TradePayload{Resources: readmodels.PriceModel{}, Storage: []readmodels.TradeStorageItemSnap{}, Army: []readmodels.TradeArmyItemSnap{}}
	}
	var dto dtos.TradePayloadDTO
	if err := json.Unmarshal(raw, &dto); err != nil {
		return readmodels.TradePayload{Resources: readmodels.PriceModel{}, Storage: []readmodels.TradeStorageItemSnap{}, Army: []readmodels.TradeArmyItemSnap{}}
	}
	return readmodels.TradePayload{
		Resources: readmodels.PriceModel{Credits: dto.Resources.Credits, Iron: dto.Resources.Iron, Titanium: dto.Resources.Titanium, Antimatter: dto.Resources.Antimatter},
		Storage:   tradePayloadStorageSnapsFromDTOs(dto.Storage),
		Army:      tradeArmySnapsFromDTOs(dto.Army),
	}
}

func tradeUnitsFromJSON(raw json.RawMessage) []readmodels.MilitaryUnitSnap {
	if len(raw) == 0 {
		return []readmodels.MilitaryUnitSnap{}
	}
	var unitDTOs []dtos.MilitaryUnitDTO
	if err := json.Unmarshal(raw, &unitDTOs); err != nil {
		return []readmodels.MilitaryUnitSnap{}
	}
	out := make([]readmodels.MilitaryUnitSnap, 0, len(unitDTOs))
	for _, d := range unitDTOs {
		out = append(out, readmodels.MilitaryUnitSnap{
			PrototypeID: d.PrototypeID,
			Category:    readmodels.ArmyCategory(d.Category),
			Attack:      d.Attack,
			Defence:     d.Defence,
			Capacity:    d.Capacity,
			Stealth:     d.Stealth,
			Speed:       d.Speed,
			Count:       d.Count,
		})
	}
	return out
}

func tradeOperationStorageSnapsFromJSON(raw json.RawMessage) []readmodels.StorageItemSnap {
	if len(raw) == 0 {
		return []readmodels.StorageItemSnap{}
	}
	var snapDTOs []dtos.TradeStorageItemSnapDTO
	if err := json.Unmarshal(raw, &snapDTOs); err != nil {
		return []readmodels.StorageItemSnap{}
	}
	return tradeOperationStorageSnapsFromDTOs(snapDTOs)
}

func tradePayloadStorageSnapsFromDTOs(items []dtos.TradeStorageItemSnapDTO) []readmodels.TradeStorageItemSnap {
	out := make([]readmodels.TradeStorageItemSnap, 0, len(items))
	for _, s := range items {
		out = append(out, readmodels.TradeStorageItemSnap{ItemID: s.ItemID, PrototypeID: s.PrototypeID, Category: readmodels.StorageCategory(s.Category)})
	}
	return out
}

func tradeOperationStorageSnapsFromDTOs(items []dtos.TradeStorageItemSnapDTO) []readmodels.StorageItemSnap {
	out := make([]readmodels.StorageItemSnap, 0, len(items))
	for _, s := range items {
		out = append(out, readmodels.StorageItemSnap{PrototypeID: s.PrototypeID, Category: readmodels.StorageCategory(s.Category)})
	}
	return out
}

func tradeArmySnapsFromDTOs(items []dtos.TradeArmyItemSnapDTO) []readmodels.TradeArmyItemSnap {
	out := make([]readmodels.TradeArmyItemSnap, 0, len(items))
	for _, a := range items {
		out = append(out, readmodels.TradeArmyItemSnap{PrototypeID: a.PrototypeID, Count: a.Count, Capacity: a.Capacity})
	}
	return out
}

// EnrichTradeOperationItems populates CurrentPrototype on transport units, storage snaps,
// and payload army/storage items using the provided prototype maps.
func EnrichTradeOperationItems(op *readmodels.TradeOperation, armyMap map[int]readmodels.ArmyItemPrototype, storageMap map[int]readmodels.StorageItemPrototype) {
	for i := range op.TransportUnits {
		if proto, ok := armyMap[op.TransportUnits[i].PrototypeID]; ok {
			op.TransportUnits[i].CurrentPrototype = proto
			op.TransportUnits[i].Name = proto.Name
			op.TransportUnits[i].ImageURL = proto.ImageURL
			op.TransportUnits[i].Space = proto.Space
		}
	}
	for i := range op.StorageSnaps {
		if proto, ok := storageMap[op.StorageSnaps[i].PrototypeID]; ok {
			op.StorageSnaps[i].CurrentPrototype = proto
			op.StorageSnaps[i].Name = proto.Name
			op.StorageSnaps[i].ShortDescription = proto.ShortDescription
			op.StorageSnaps[i].ImageURL = proto.ImageURL
		}
	}
	enrichTradePayloadSnaps(&op.OfferedPayload, armyMap, storageMap)
	enrichTradePayloadSnaps(&op.RequestedPayload, armyMap, storageMap)
}

func enrichTradePayloadSnaps(payload *readmodels.TradePayload, armyMap map[int]readmodels.ArmyItemPrototype, storageMap map[int]readmodels.StorageItemPrototype) {
	for i := range payload.Army {
		if proto, ok := armyMap[payload.Army[i].PrototypeID]; ok {
			payload.Army[i].CurrentPrototype = proto
		}
	}
	for i := range payload.Storage {
		if proto, ok := storageMap[payload.Storage[i].PrototypeID]; ok {
			payload.Storage[i].CurrentPrototype = proto
		}
	}
}

func tradeModifiersFromJSON(raw json.RawMessage) readmodels.MilitaryModifiers {
	if len(raw) == 0 {
		return readmodels.MilitaryModifiers{}
	}
	var dto dtos.MilitaryModifiersDTO
	if err := json.Unmarshal(raw, &dto); err != nil {
		return readmodels.MilitaryModifiers{}
	}
	return readmodels.MilitaryModifiers{AttackMul: float32(dto.AttackMul), DefenceMul: float32(dto.DefenceMul), StealthMul: float32(dto.StealthMul), CapacityMul: float32(dto.CapacityMul), SpeedMul: float32(dto.SpeedMul)}
}
