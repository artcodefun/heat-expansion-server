package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
	"github.com/sqlc-dev/pqtype"
)

// OperationFromModel converts a database row to a read model.
func OperationFromModel(m gen.MilitaryOperation) readmodels.MilitaryOperation {
	return readmodels.MilitaryOperation{
		ID:                 int(m.ID),
		Type:               readmodels.MilitaryOperationType(m.Type),
		OwnerUserID:        int(m.OwnerUserID),
		SourceBaseID:       int(m.SourceBaseID),
		SourceCoordinates:  readmodels.Vector2i{X: int(m.SourceX), Y: int(m.SourceY)},
		TargetCoordinates:  readmodels.Vector2i{X: int(m.TargetX), Y: int(m.TargetY)},
		OutboundDepartAt:   m.OutboundDepartAt,
		OutboundArriveAt:   m.OutboundArriveAt,
		ReturnDepartAt:     m.ReturnDepartAt,
		ReturnArriveAt:     m.ReturnArriveAt,
		CompletedAt:        m.CompletedAt,
		CrystalsSkipPrice:  int(m.CrystalsSkipPrice),
		Phase:              readmodels.MilitaryOperationPhase(m.Phase),
		Result:             readmodels.MilitaryOperationResult(m.Result),
		Units:              militaryUnitsFromJSON(m.Units),
		StorageSnaps:       storageSnapsFromJSON(m.StorageSnaps),
		TotalModifiers:     militaryModifiersFromJSON(m.TotalModifiers),
		SpyResult:          spyResultFromJSON(m.SpyResult),
		AttackResult:       attackResultFromJSON(m.AttackResult),
		ProducedScanReport: nil,
	}
}

// JSON helpers: DTO shape -> readmodels.*

func storageSnapsFromJSON(raw json.RawMessage) []readmodels.StorageItemSnap {
	if len(raw) == 0 {
		return []readmodels.StorageItemSnap{}
	}
	var snapsDTO []dtos.StorageItemSnapDTO
	if err := json.Unmarshal(raw, &snapsDTO); err != nil {
		return []readmodels.StorageItemSnap{}
	}
	out := make([]readmodels.StorageItemSnap, 0, len(snapsDTO))
	for _, s := range snapsDTO {
		out = append(out, storageItemSnapFromDTO(s))
	}
	return out
}

func storageItemSnapFromDTO(d dtos.StorageItemSnapDTO) readmodels.StorageItemSnap {
	var buffData *readmodels.BuffStorageData
	if d.BuffData != nil {
		buffData = &readmodels.BuffStorageData{
			Type:            readmodels.BuffType(d.BuffData.Type),
			Value:           d.BuffData.Value,
			DurationSeconds: d.BuffData.DurationSeconds,
		}
	}
	var artifactData *readmodels.ArtifactStorageData
	if d.ArtifactData != nil {
		artifactData = &readmodels.ArtifactStorageData{
			Type:  readmodels.ArtifactEffectType(d.ArtifactData.Type),
			Value: d.ArtifactData.Value,
		}
	}
	return readmodels.StorageItemSnap{
		PrototypeID:  d.PrototypeID,
		Category:     readmodels.StorageCategory(d.Category),
		BuffData:     buffData,
		ArtifactData: artifactData,
	}
}

func militaryModifiersFromJSON(raw json.RawMessage) readmodels.MilitaryModifiers {
	if len(raw) == 0 {
		return readmodels.MilitaryModifiers{}
	}
	var d dtos.MilitaryModifiersDTO
	if err := json.Unmarshal(raw, &d); err != nil {
		return readmodels.MilitaryModifiers{}
	}
	return militaryModifiersFromDTO(d)
}

func militaryModifiersFromDTO(d dtos.MilitaryModifiersDTO) readmodels.MilitaryModifiers {
	return readmodels.MilitaryModifiers{
		AttackMul:   float32(d.AttackMul),
		DefenceMul:  float32(d.DefenceMul),
		StealthMul:  float32(d.StealthMul),
		CapacityMul: float32(d.CapacityMul),
		SpeedMul:    float32(d.SpeedMul),
	}
}

func militaryUnitsFromJSON(raw json.RawMessage) []readmodels.MilitaryUnitSnap {
	if len(raw) == 0 {
		return []readmodels.MilitaryUnitSnap{}
	}
	var unitDTOs []dtos.MilitaryUnitDTO
	if err := json.Unmarshal(raw, &unitDTOs); err != nil {
		return []readmodels.MilitaryUnitSnap{}
	}
	if len(unitDTOs) == 0 {
		return []readmodels.MilitaryUnitSnap{}
	}
	unites := make([]readmodels.MilitaryUnitSnap, 0, len(unitDTOs))
	for _, d := range unitDTOs {
		unites = append(unites, militaryUnitFromDTO(d))
	}
	return unites
}

func militaryUnitFromDTO(d dtos.MilitaryUnitDTO) readmodels.MilitaryUnitSnap {
	return readmodels.MilitaryUnitSnap{
		PrototypeID: d.PrototypeID,
		Category:    readmodels.ArmyCategory(d.Category),
		Attack:      d.Attack,
		Defence:     d.Defence,
		Capacity:    d.Capacity,
		Stealth:     d.Stealth,
		Speed:       d.Speed,
		Count:       d.Count,
	}
}

func defenseStructureFromDTO(d dtos.DefenseStructureDTO) readmodels.DefenseStructureSnap {
	return readmodels.DefenseStructureSnap{
		PrototypeID: d.PrototypeID,
		Defence:     d.Defence,
		Count:       d.Count,
	}
}

func trophyFromDTO(d dtos.TrophyDTO) readmodels.TrophyStorageItem {
	return readmodels.TrophyStorageItem{
		Prototype: readmodels.StorageItemPrototype{
			ID: d.PrototypeID,
		},
	}
}

func spyResultFromJSON(nm pqtype.NullRawMessage) *readmodels.SpyResult {
	if !nm.Valid {
		return nil
	}
	var d dtos.SpyResultDTO
	if err := json.Unmarshal(nm.RawMessage, &d); err != nil {
		return nil
	}
	res := &readmodels.SpyResult{
		Outcome:                readmodels.SpyOutcome(d.Outcome),
		TotalDefenderModifiers: militaryModifiersFromDTO(d.TotalDefenderModifiers),
	}
	if len(d.AttackerRemaining) > 0 {
		res.AttackerRemaining = make([]readmodels.MilitaryUnitSnap, 0, len(d.AttackerRemaining))
		for _, u := range d.AttackerRemaining {
			res.AttackerRemaining = append(res.AttackerRemaining, militaryUnitFromDTO(u))
		}
	}
	if len(d.DefenderRemaining) > 0 {
		res.DefenderRemaining = make([]readmodels.MilitaryUnitSnap, 0, len(d.DefenderRemaining))
		for _, u := range d.DefenderRemaining {
			res.DefenderRemaining = append(res.DefenderRemaining, militaryUnitFromDTO(u))
		}
	}
	if len(d.DefendersBefore) > 0 {
		res.DefendersBefore = make([]readmodels.MilitaryUnitSnap, 0, len(d.DefendersBefore))
		for _, u := range d.DefendersBefore {
			res.DefendersBefore = append(res.DefendersBefore, militaryUnitFromDTO(u))
		}
	}
	if len(d.DefenderStorageSnaps) > 0 {
		res.DefenderStorageSnaps = make([]readmodels.StorageItemSnap, 0, len(d.DefenderStorageSnaps))
		for _, s := range d.DefenderStorageSnaps {
			res.DefenderStorageSnaps = append(res.DefenderStorageSnaps, storageItemSnapFromDTO(s))
		}
	}
	return res
}

func attackResultFromJSON(nm pqtype.NullRawMessage) *readmodels.AttackResult {
	if !nm.Valid {
		return nil
	}
	var d dtos.AttackResultDTO
	if err := json.Unmarshal(nm.RawMessage, &d); err != nil {
		return nil
	}
	res := &readmodels.AttackResult{
		Outcome: readmodels.AttackOutcome(d.Outcome),
		Loot: readmodels.PriceModel{
			Credits:    d.Loot.Credits,
			Iron:       d.Loot.Iron,
			Titanium:   d.Loot.Titanium,
			Antimatter: d.Loot.Antimatter,
		},
		TotalDefenderModifiers: militaryModifiersFromDTO(d.TotalDefenderModifiers),
	}
	if len(d.AttackerRemaining) > 0 {
		res.AttackerRemaining = make([]readmodels.MilitaryUnitSnap, 0, len(d.AttackerRemaining))
		for _, u := range d.AttackerRemaining {
			res.AttackerRemaining = append(res.AttackerRemaining, militaryUnitFromDTO(u))
		}
	}
	if len(d.DefenderRemaining) > 0 {
		res.DefenderRemaining = make([]readmodels.MilitaryUnitSnap, 0, len(d.DefenderRemaining))
		for _, u := range d.DefenderRemaining {
			res.DefenderRemaining = append(res.DefenderRemaining, militaryUnitFromDTO(u))
		}
	}
	if len(d.RemainingStructures) > 0 {
		res.RemainingStructures = make([]readmodels.DefenseStructureSnap, 0, len(d.RemainingStructures))
		for _, s := range d.RemainingStructures {
			res.RemainingStructures = append(res.RemainingStructures, defenseStructureFromDTO(s))
		}
	}
	if len(d.Trophies) > 0 {
		res.Trophies = make([]readmodels.TrophyStorageItem, 0, len(d.Trophies))
		for _, t := range d.Trophies {
			res.Trophies = append(res.Trophies, trophyFromDTO(t))
		}
	}
	if len(d.DefendersBefore) > 0 {
		res.DefendersBefore = make([]readmodels.MilitaryUnitSnap, 0, len(d.DefendersBefore))
		for _, u := range d.DefendersBefore {
			res.DefendersBefore = append(res.DefendersBefore, militaryUnitFromDTO(u))
		}
	}
	if len(d.StructuresBefore) > 0 {
		res.StructuresBefore = make([]readmodels.DefenseStructureSnap, 0, len(d.StructuresBefore))
		for _, s := range d.StructuresBefore {
			res.StructuresBefore = append(res.StructuresBefore, defenseStructureFromDTO(s))
		}
	}
	if len(d.DefenderStorageSnaps) > 0 {
		res.DefenderStorageSnaps = make([]readmodels.StorageItemSnap, 0, len(d.DefenderStorageSnaps))
		for _, s := range d.DefenderStorageSnaps {
			res.DefenderStorageSnaps = append(res.DefenderStorageSnaps, storageItemSnapFromDTO(s))
		}
	}
	return res
}

// EnrichOperationUnitsAndStructures builds a MilitaryOperation readmodel and enriches its
// units/structures with full prototype data.
func EnrichOperationUnitsAndStructures(op *readmodels.MilitaryOperation, armyMap map[int]readmodels.ArmyItemPrototype, buildMap map[int]readmodels.BuildItemPrototype, storageMap map[int]readmodels.StorageItemPrototype) {
	for i := range op.Units {
		if proto, ok := armyMap[op.Units[i].PrototypeID]; ok {
			op.Units[i].CurrentPrototype = proto
			op.Units[i].Name = proto.Name
			op.Units[i].ImageURL = proto.ImageURL
			op.Units[i].Space = proto.Space
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
	if op.SpyResult != nil {
		for i := range op.SpyResult.AttackerRemaining {
			if proto, ok := armyMap[op.SpyResult.AttackerRemaining[i].PrototypeID]; ok {
				op.SpyResult.AttackerRemaining[i].CurrentPrototype = proto
				op.SpyResult.AttackerRemaining[i].Name = proto.Name
				op.SpyResult.AttackerRemaining[i].ImageURL = proto.ImageURL
				op.SpyResult.AttackerRemaining[i].Space = proto.Space
			}
		}
		for i := range op.SpyResult.DefenderRemaining {
			if proto, ok := armyMap[op.SpyResult.DefenderRemaining[i].PrototypeID]; ok {
				op.SpyResult.DefenderRemaining[i].CurrentPrototype = proto
				op.SpyResult.DefenderRemaining[i].Name = proto.Name
				op.SpyResult.DefenderRemaining[i].ImageURL = proto.ImageURL
				op.SpyResult.DefenderRemaining[i].Space = proto.Space
			}
		}
		for i := range op.SpyResult.DefendersBefore {
			if proto, ok := armyMap[op.SpyResult.DefendersBefore[i].PrototypeID]; ok {
				op.SpyResult.DefendersBefore[i].CurrentPrototype = proto
				op.SpyResult.DefendersBefore[i].Name = proto.Name
				op.SpyResult.DefendersBefore[i].ImageURL = proto.ImageURL
				op.SpyResult.DefendersBefore[i].Space = proto.Space
			}
		}
		for i := range op.SpyResult.DefenderStorageSnaps {
			if proto, ok := storageMap[op.SpyResult.DefenderStorageSnaps[i].PrototypeID]; ok {
				op.SpyResult.DefenderStorageSnaps[i].CurrentPrototype = proto
				op.SpyResult.DefenderStorageSnaps[i].Name = proto.Name
				op.SpyResult.DefenderStorageSnaps[i].ShortDescription = proto.ShortDescription
				op.SpyResult.DefenderStorageSnaps[i].ImageURL = proto.ImageURL
			}
		}
	}
	if op.AttackResult != nil {
		for i := range op.AttackResult.AttackerRemaining {
			if proto, ok := armyMap[op.AttackResult.AttackerRemaining[i].PrototypeID]; ok {
				op.AttackResult.AttackerRemaining[i].CurrentPrototype = proto
				op.AttackResult.AttackerRemaining[i].Name = proto.Name
				op.AttackResult.AttackerRemaining[i].ImageURL = proto.ImageURL
				op.AttackResult.AttackerRemaining[i].Space = proto.Space
			}
		}
		for i := range op.AttackResult.DefenderRemaining {
			if proto, ok := armyMap[op.AttackResult.DefenderRemaining[i].PrototypeID]; ok {
				op.AttackResult.DefenderRemaining[i].CurrentPrototype = proto
				op.AttackResult.DefenderRemaining[i].Name = proto.Name
				op.AttackResult.DefenderRemaining[i].ImageURL = proto.ImageURL
				op.AttackResult.DefenderRemaining[i].Space = proto.Space
			}
		}
		for i := range op.AttackResult.RemainingStructures {
			if proto, ok := buildMap[op.AttackResult.RemainingStructures[i].PrototypeID]; ok {
				op.AttackResult.RemainingStructures[i].CurrentPrototype = proto
				op.AttackResult.RemainingStructures[i].Name = proto.Name
				op.AttackResult.RemainingStructures[i].ImageURL = proto.ImageURL
				op.AttackResult.RemainingStructures[i].Space = proto.Space
			}
		}
		for i := range op.AttackResult.Trophies {
			if proto, ok := storageMap[op.AttackResult.Trophies[i].Prototype.ID]; ok {
				op.AttackResult.Trophies[i].Prototype = proto
			}
		}
		for i := range op.AttackResult.DefendersBefore {
			if proto, ok := armyMap[op.AttackResult.DefendersBefore[i].PrototypeID]; ok {
				op.AttackResult.DefendersBefore[i].CurrentPrototype = proto
				op.AttackResult.DefendersBefore[i].Name = proto.Name
				op.AttackResult.DefendersBefore[i].ImageURL = proto.ImageURL
				op.AttackResult.DefendersBefore[i].Space = proto.Space
			}
		}
		for i := range op.AttackResult.StructuresBefore {
			if proto, ok := buildMap[op.AttackResult.StructuresBefore[i].PrototypeID]; ok {
				op.AttackResult.StructuresBefore[i].CurrentPrototype = proto
				op.AttackResult.StructuresBefore[i].Name = proto.Name
				op.AttackResult.StructuresBefore[i].ImageURL = proto.ImageURL
				op.AttackResult.StructuresBefore[i].Space = proto.Space
			}
		}
		for i := range op.AttackResult.DefenderStorageSnaps {
			if proto, ok := storageMap[op.AttackResult.DefenderStorageSnaps[i].PrototypeID]; ok {
				op.AttackResult.DefenderStorageSnaps[i].CurrentPrototype = proto
				op.AttackResult.DefenderStorageSnaps[i].Name = proto.Name
				op.AttackResult.DefenderStorageSnaps[i].ShortDescription = proto.ShortDescription
				op.AttackResult.DefenderStorageSnaps[i].ImageURL = proto.ImageURL
			}
		}
	}
}
