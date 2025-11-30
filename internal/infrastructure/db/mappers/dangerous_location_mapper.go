package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/gen"
)

func DangerousLocationFromDB(r gen.DangerousLocation) *domain.DangerousLocationModel {
	var resDTO dtos.LocationResourceStatsDTO
	_ = json.Unmarshal(r.Resources, &resDTO)
	res := dtos.LocationResourceStatsFromDTO(resDTO, r.ResourcesCalcTimestamp)

	var unitDTOs []dtos.MilitaryUnitDTO
	_ = json.Unmarshal(r.Units, &unitDTOs)
	units := make([]domain.MilitaryUnit, 0, len(unitDTOs))
	for _, d := range unitDTOs {
		units = append(units, dtos.MilitaryUnitFromDTO(d))
	}

	var structDTOs []dtos.DefenseStructureDTO
	_ = json.Unmarshal(r.Structures, &structDTOs)
	structs := make([]domain.DefenseStructure, 0, len(structDTOs))
	for _, d := range structDTOs {
		structs = append(structs, dtos.DefenseStructureFromDTO(d))
	}

	return &domain.DangerousLocationModel{
		ID:          int(r.ID),
		Coordinates: domain.Vector2i{X: int(r.SectorX), Y: int(r.SectorY)},
		DangerLevel: int(r.DangerLevel),
		LocationDetails: domain.LocationDetails{
			Name:        nullStringToString(&r.Name.String, r.Name.Valid),
			Description: nullStringToString(&r.Description.String, r.Description.Valid),
			ImageURL:    nullStringToString(&r.ImageUrl.String, r.ImageUrl.Valid),
		},
		Resources:      res,
		DefendingUnits: units,
		Structures:     structs,
	}
}

func InsertDangerousLocationParamsFromDomain(loc *domain.DangerousLocationModel) gen.InsertDangerousLocationParams {
	resJSON, _ := json.Marshal(dtos.LocationResourceStatsDTOFromDomain(loc.Resources))
	unitDTOs := make([]dtos.MilitaryUnitDTO, 0, len(loc.DefendingUnits))
	for _, u := range loc.DefendingUnits {
		unitDTOs = append(unitDTOs, dtos.MilitaryUnitDTOFromDomain(u))
	}
	structsDTO := make([]dtos.DefenseStructureDTO, 0, len(loc.Structures))
	for _, s := range loc.Structures {
		structsDTO = append(structsDTO, dtos.DefenseStructureDTOFromDomain(s))
	}
	unitsJSON, _ := json.Marshal(unitDTOs)
	structsJSON, _ := json.Marshal(structsDTO)
	return gen.InsertDangerousLocationParams{
		SectorX:                int32(loc.Coordinates.X),
		SectorY:                int32(loc.Coordinates.Y),
		DangerLevel:            int32(loc.DangerLevel),
		Name:                   toNullString(loc.Name),
		Description:            toNullString(loc.Description),
		ImageUrl:               toNullString(loc.ImageURL),
		Resources:              resJSON,
		ResourcesCalcTimestamp: loc.Resources.CalculationTimestamp,
		Units:                  unitsJSON,
		Structures:             structsJSON,
	}
}

func UpdateDangerousLocationParamsFromDomain(loc *domain.DangerousLocationModel) gen.UpdateDangerousLocationParams {
	resJSON, _ := json.Marshal(dtos.LocationResourceStatsDTOFromDomain(loc.Resources))
	unitDTOs := make([]dtos.MilitaryUnitDTO, 0, len(loc.DefendingUnits))
	for _, u := range loc.DefendingUnits {
		unitDTOs = append(unitDTOs, dtos.MilitaryUnitDTOFromDomain(u))
	}
	structsDTO := make([]dtos.DefenseStructureDTO, 0, len(loc.Structures))
	for _, s := range loc.Structures {
		structsDTO = append(structsDTO, dtos.DefenseStructureDTOFromDomain(s))
	}
	unitsJSON, _ := json.Marshal(unitDTOs)
	structsJSON, _ := json.Marshal(structsDTO)
	return gen.UpdateDangerousLocationParams{
		ID:                     int64(loc.ID),
		DangerLevel:            int32(loc.DangerLevel),
		Name:                   toNullString(loc.Name),
		Description:            toNullString(loc.Description),
		ImageUrl:               toNullString(loc.ImageURL),
		Resources:              resJSON,
		ResourcesCalcTimestamp: loc.Resources.CalculationTimestamp,
		Units:                  unitsJSON,
		Structures:             structsJSON,
	}
}
