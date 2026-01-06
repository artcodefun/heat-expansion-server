package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/gen"
	"github.com/sqlc-dev/pqtype"
)

// ArmyPrototypeSnapshot and BuildPrototypeSnapshot are lightweight views used to enrich readmodels
// with display-only prototype data like name and image URL.
type ArmyPrototypeSnapshot struct {
	Name     string
	ImageURL string
}

type BuildPrototypeSnapshot struct {
	Name     string
	ImageURL string
}

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
		SpyResult:          spyResultFromJSON(m.SpyResult),
		AttackResult:       attackResultFromJSON(m.AttackResult),
		ProducedScanReport: nil,
	}
}

// JSON helpers: DTO shape -> readmodels.*

func militaryUnitsFromJSON(raw json.RawMessage) []readmodels.MilitaryUnit {
	if len(raw) == 0 {
		return []readmodels.MilitaryUnit{}
	}
	var unitDTOs []dtos.MilitaryUnitDTO
	if err := json.Unmarshal(raw, &unitDTOs); err != nil {
		return []readmodels.MilitaryUnit{}
	}
	if len(unitDTOs) == 0 {
		return []readmodels.MilitaryUnit{}
	}
	unites := make([]readmodels.MilitaryUnit, 0, len(unitDTOs))
	for _, d := range unitDTOs {
		unites = append(unites, militaryUnitFromDTO(d))
	}
	return unites
}

func militaryUnitFromDTO(d dtos.MilitaryUnitDTO) readmodels.MilitaryUnit {
	return readmodels.MilitaryUnit{
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

func defenseStructureFromDTO(d dtos.DefenseStructureDTO) readmodels.DefenseStructure {
	return readmodels.DefenseStructure{
		PrototypeID: d.PrototypeID,
		Defence:     d.Defence,
		Count:       d.Count,
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
		Outcome: readmodels.SpyOutcome(d.Outcome),
	}
	if len(d.AttackerRemaining) > 0 {
		res.AttackerRemaining = make([]readmodels.MilitaryUnit, 0, len(d.AttackerRemaining))
		for _, u := range d.AttackerRemaining {
			res.AttackerRemaining = append(res.AttackerRemaining, militaryUnitFromDTO(u))
		}
	}
	if len(d.DefenderRemaining) > 0 {
		res.DefenderRemaining = make([]readmodels.MilitaryUnit, 0, len(d.DefenderRemaining))
		for _, u := range d.DefenderRemaining {
			res.DefenderRemaining = append(res.DefenderRemaining, militaryUnitFromDTO(u))
		}
	}
	if len(d.DefendersBefore) > 0 {
		res.DefendersBefore = make([]readmodels.MilitaryUnit, 0, len(d.DefendersBefore))
		for _, u := range d.DefendersBefore {
			res.DefendersBefore = append(res.DefendersBefore, militaryUnitFromDTO(u))
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
	}
	if len(d.AttackerRemaining) > 0 {
		res.AttackerRemaining = make([]readmodels.MilitaryUnit, 0, len(d.AttackerRemaining))
		for _, u := range d.AttackerRemaining {
			res.AttackerRemaining = append(res.AttackerRemaining, militaryUnitFromDTO(u))
		}
	}
	if len(d.DefenderRemaining) > 0 {
		res.DefenderRemaining = make([]readmodels.MilitaryUnit, 0, len(d.DefenderRemaining))
		for _, u := range d.DefenderRemaining {
			res.DefenderRemaining = append(res.DefenderRemaining, militaryUnitFromDTO(u))
		}
	}
	if len(d.RemainingStructures) > 0 {
		res.RemainingStructures = make([]readmodels.DefenseStructure, 0, len(d.RemainingStructures))
		for _, s := range d.RemainingStructures {
			res.RemainingStructures = append(res.RemainingStructures, defenseStructureFromDTO(s))
		}
	}
	if len(d.DefendersBefore) > 0 {
		res.DefendersBefore = make([]readmodels.MilitaryUnit, 0, len(d.DefendersBefore))
		for _, u := range d.DefendersBefore {
			res.DefendersBefore = append(res.DefendersBefore, militaryUnitFromDTO(u))
		}
	}
	if len(d.StructuresBefore) > 0 {
		res.StructuresBefore = make([]readmodels.DefenseStructure, 0, len(d.StructuresBefore))
		for _, s := range d.StructuresBefore {
			res.StructuresBefore = append(res.StructuresBefore, defenseStructureFromDTO(s))
		}
	}
	return res
}

// EnrichOperationUnitsAndStructures builds a MilitaryOperation readmodel and enriches its
// units/structures with prototype-derived name and image data.
func EnrichOperationUnitsAndStructures(op *readmodels.MilitaryOperation, armyMap map[int]ArmyPrototypeSnapshot, buildMap map[int]BuildPrototypeSnapshot) {
	for i := range op.Units {
		if proto, ok := armyMap[op.Units[i].PrototypeID]; ok {
			op.Units[i].Name = proto.Name
			op.Units[i].ImageURL = proto.ImageURL
		}
	}
	if op.SpyResult != nil {
		for i := range op.SpyResult.AttackerRemaining {
			if proto, ok := armyMap[op.SpyResult.AttackerRemaining[i].PrototypeID]; ok {
				op.SpyResult.AttackerRemaining[i].Name = proto.Name
				op.SpyResult.AttackerRemaining[i].ImageURL = proto.ImageURL
			}
		}
		for i := range op.SpyResult.DefenderRemaining {
			if proto, ok := armyMap[op.SpyResult.DefenderRemaining[i].PrototypeID]; ok {
				op.SpyResult.DefenderRemaining[i].Name = proto.Name
				op.SpyResult.DefenderRemaining[i].ImageURL = proto.ImageURL
			}
		}
		for i := range op.SpyResult.DefendersBefore {
			if proto, ok := armyMap[op.SpyResult.DefendersBefore[i].PrototypeID]; ok {
				op.SpyResult.DefendersBefore[i].Name = proto.Name
				op.SpyResult.DefendersBefore[i].ImageURL = proto.ImageURL
			}
		}
	}
	if op.AttackResult != nil {
		for i := range op.AttackResult.AttackerRemaining {
			if proto, ok := armyMap[op.AttackResult.AttackerRemaining[i].PrototypeID]; ok {
				op.AttackResult.AttackerRemaining[i].Name = proto.Name
				op.AttackResult.AttackerRemaining[i].ImageURL = proto.ImageURL
			}
		}
		for i := range op.AttackResult.DefenderRemaining {
			if proto, ok := armyMap[op.AttackResult.DefenderRemaining[i].PrototypeID]; ok {
				op.AttackResult.DefenderRemaining[i].Name = proto.Name
				op.AttackResult.DefenderRemaining[i].ImageURL = proto.ImageURL
			}
		}
		for i := range op.AttackResult.RemainingStructures {
			if proto, ok := buildMap[op.AttackResult.RemainingStructures[i].PrototypeID]; ok {
				op.AttackResult.RemainingStructures[i].Name = proto.Name
				op.AttackResult.RemainingStructures[i].ImageURL = proto.ImageURL
			}
		}
		for i := range op.AttackResult.DefendersBefore {
			if proto, ok := armyMap[op.AttackResult.DefendersBefore[i].PrototypeID]; ok {
				op.AttackResult.DefendersBefore[i].Name = proto.Name
				op.AttackResult.DefendersBefore[i].ImageURL = proto.ImageURL
			}
		}
		for i := range op.AttackResult.StructuresBefore {
			if proto, ok := buildMap[op.AttackResult.StructuresBefore[i].PrototypeID]; ok {
				op.AttackResult.StructuresBefore[i].Name = proto.Name
				op.AttackResult.StructuresBefore[i].ImageURL = proto.ImageURL
			}
		}
	}
}
