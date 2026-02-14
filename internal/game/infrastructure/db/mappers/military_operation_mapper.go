package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
)

func MilitaryOperationFromDB(r gen.MilitaryOperation) *domain.MilitaryOperation {
	var unitDTOs []dtos.MilitaryUnitDTO
	_ = json.Unmarshal(r.Units, &unitDTOs)
	units := make([]domain.MilitaryUnitSnap, 0, len(unitDTOs))
	for _, d := range unitDTOs {
		units = append(units, dtos.MilitaryUnitFromDTO(d))
	}

	var storageDTOs []dtos.StorageItemSnapDTO
	_ = json.Unmarshal(r.StorageSnaps, &storageDTOs)
	storageSnaps := make([]domain.StorageItemSnap, 0, len(storageDTOs))
	for _, d := range storageDTOs {
		storageSnaps = append(storageSnaps, dtos.StorageItemSnapFromDTO(d))
	}

	var milModsDTO dtos.MilitaryModifiersDTO
	_ = json.Unmarshal(r.TotalModifiers, &milModsDTO)

	var spyDTO *dtos.SpyResultDTO
	unmarshalIfValid(r.SpyResult, &spyDTO)
	var attackDTO *dtos.AttackResultDTO
	unmarshalIfValid(r.AttackResult, &attackDTO)

	op := &domain.MilitaryOperation{
		ID:                int(r.ID),
		Type:              domain.MilitaryOperationType(r.Type),
		OwnerUserID:       r.OwnerUserID,
		SourceBaseID:      int(r.SourceBaseID),
		SourceCoordinates: domain.Vector2i{X: int(r.SourceX), Y: int(r.SourceY)},
		TargetCoordinates: domain.Vector2i{X: int(r.TargetX), Y: int(r.TargetY)},
		OutboundDepartAt:  r.OutboundDepartAt,
		OutboundArriveAt:  r.OutboundArriveAt,
		ReturnDepartAt:    r.ReturnDepartAt,
		ReturnArriveAt:    r.ReturnArriveAt,
		CompletedAt:       r.CompletedAt,
		CrystalsSkipPrice: int(r.CrystalsSkipPrice),
		Phase:             domain.MilitaryOperationPhase(r.Phase),
		Result:            domain.MilitaryOperationResult(r.Result),
		Units:             units,
		StorageSnaps:      storageSnaps,
		TotalModifiers:    dtos.MilitaryModifiersFromDTO(milModsDTO),
		SpyResult:         dtos.SpyResultFromDTO(spyDTO),
		AttackResult:      dtos.AttackResultFromDTO(attackDTO),
	}
	return op
}

func InsertMilitaryOperationParamsFromDomain(op *domain.MilitaryOperation) gen.InsertMilitaryOperationParams {
	unitDTOs := make([]dtos.MilitaryUnitDTO, 0, len(op.Units))
	for _, u := range op.Units {
		unitDTOs = append(unitDTOs, dtos.MilitaryUnitDTOFromDomain(u))
	}
	unitsJSON, _ := json.Marshal(unitDTOs)

	storageDTOs := make([]dtos.StorageItemSnapDTO, 0, len(op.StorageSnaps))
	for _, s := range op.StorageSnaps {
		storageDTOs = append(storageDTOs, dtos.StorageItemSnapDTOFromDomain(s))
	}
	storageJSON, _ := json.Marshal(storageDTOs)

	milModsDTO := dtos.MilitaryModifiersDTOFromDomain(op.TotalModifiers)
	milModsJSON, _ := json.Marshal(milModsDTO)

	spyDTO := dtos.SpyResultDTOFromDomain(op.SpyResult)
	attackDTO := dtos.AttackResultDTOFromDomain(op.AttackResult)

	return gen.InsertMilitaryOperationParams{
		Type:              string(op.Type),
		OwnerUserID:       op.OwnerUserID,
		SourceBaseID:      int64(op.SourceBaseID),
		SourceX:           int32(op.SourceCoordinates.X),
		SourceY:           int32(op.SourceCoordinates.Y),
		TargetX:           int32(op.TargetCoordinates.X),
		TargetY:           int32(op.TargetCoordinates.Y),
		OutboundDepartAt:  op.OutboundDepartAt,
		OutboundArriveAt:  op.OutboundArriveAt,
		ReturnDepartAt:    op.ReturnDepartAt,
		ReturnArriveAt:    op.ReturnArriveAt,
		CompletedAt:       op.CompletedAt,
		CrystalsSkipPrice: int32(op.CrystalsSkipPrice),
		Phase:             string(op.Phase),
		Result:            string(op.Result),
		Units:             unitsJSON,
		StorageSnaps:      storageJSON,
		TotalModifiers:    milModsJSON,
		SpyResult:         toNullRawMessage(spyDTO),
		AttackResult:      toNullRawMessage(attackDTO),
	}
}

func UpdateMilitaryOperationParamsFromDomain(op *domain.MilitaryOperation) gen.UpdateMilitaryOperationParams {
	unitDTOs := make([]dtos.MilitaryUnitDTO, 0, len(op.Units))
	for _, u := range op.Units {
		unitDTOs = append(unitDTOs, dtos.MilitaryUnitDTOFromDomain(u))
	}
	unitsJSON, _ := json.Marshal(unitDTOs)

	storageDTOs := make([]dtos.StorageItemSnapDTO, 0, len(op.StorageSnaps))
	for _, s := range op.StorageSnaps {
		storageDTOs = append(storageDTOs, dtos.StorageItemSnapDTOFromDomain(s))
	}
	storageJSON, _ := json.Marshal(storageDTOs)

	milModsDTO := dtos.MilitaryModifiersDTOFromDomain(op.TotalModifiers)
	milModsJSON, _ := json.Marshal(milModsDTO)

	spyDTO := dtos.SpyResultDTOFromDomain(op.SpyResult)
	attackDTO := dtos.AttackResultDTOFromDomain(op.AttackResult)

	return gen.UpdateMilitaryOperationParams{
		ID:                int64(op.ID),
		Type:              string(op.Type),
		OwnerUserID:       op.OwnerUserID,
		SourceBaseID:      int64(op.SourceBaseID),
		SourceX:           int32(op.SourceCoordinates.X),
		SourceY:           int32(op.SourceCoordinates.Y),
		TargetX:           int32(op.TargetCoordinates.X),
		TargetY:           int32(op.TargetCoordinates.Y),
		OutboundDepartAt:  op.OutboundDepartAt,
		OutboundArriveAt:  op.OutboundArriveAt,
		ReturnDepartAt:    op.ReturnDepartAt,
		ReturnArriveAt:    op.ReturnArriveAt,
		CompletedAt:       op.CompletedAt,
		CrystalsSkipPrice: int32(op.CrystalsSkipPrice),
		Phase:             string(op.Phase),
		Result:            string(op.Result),
		Units:             unitsJSON,
		StorageSnaps:      storageJSON,
		TotalModifiers:    milModsJSON,
		SpyResult:         toNullRawMessage(spyDTO),
		AttackResult:      toNullRawMessage(attackDTO),
	}
}
