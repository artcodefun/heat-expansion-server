package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/gen"
)

func ResourceLocationFromDB(r gen.ResourceLocation, armyProtos map[int]*domain.ArmyItemPrototype, buildProtos map[int]*domain.BuildItemPrototype) *domain.ResourceLocationModel {
	var resDTO dtos.LocationResourceStatsDTO
	_ = json.Unmarshal(r.Resources, &resDTO)
	res := dtos.LocationResourceStatsFromDTO(resDTO, r.ResourcesCalcTimestamp)

	var armyDTOs []dtos.ArmyStackDTO
	_ = json.Unmarshal(r.Armies, &armyDTOs)
	armies := make([]domain.ArmyStack, 0, len(armyDTOs))
	for _, d := range armyDTOs {
		if p, ok := armyProtos[d.PrototypeID]; ok {
			armies = append(armies, dtos.ArmyStackFromDTO(d, *p))
		}
	}

	var structDTOs []dtos.DefenseStackDTO
	_ = json.Unmarshal(r.Buildings, &structDTOs)
	structs := make([]domain.DefenseStack, 0, len(structDTOs))
	for _, d := range structDTOs {
		if p, ok := buildProtos[d.PrototypeID]; ok {
			structs = append(structs, dtos.DefenseStackFromDTO(d, *p))
		}
	}

	return &domain.ResourceLocationModel{
		ID:          int(r.ID),
		Coordinates: domain.Vector2i{X: int(r.SectorX), Y: int(r.SectorY)},
		LocationDetails: domain.LocationDetails{
			Name:        nullStringToString(&r.Name.String, r.Name.Valid),
			Description: nullStringToString(&r.Description.String, r.Description.Valid),
			ImageURL:    nullStringToString(&r.ImageUrl.String, r.ImageUrl.Valid),
		},
		Type:                r.Type,
		Amount:              int(r.Amount),
		Resources:           res,
		DefendingArmies:     armies,
		DefendingStructures: structs,
	}
}

func InsertResourceLocationParamsFromDomain(loc *domain.ResourceLocationModel) gen.InsertResourceLocationParams {
	resJSON, _ := json.Marshal(dtos.LocationResourceStatsDTOFromDomain(loc.Resources))
	armyDTOs := make([]dtos.ArmyStackDTO, 0, len(loc.DefendingArmies))
	for _, a := range loc.DefendingArmies {
		armyDTOs = append(armyDTOs, dtos.ArmyStackDTOFromDomain(a))
	}
	structsDTO := make([]dtos.DefenseStackDTO, 0, len(loc.DefendingStructures))
	for _, s := range loc.DefendingStructures {
		structsDTO = append(structsDTO, dtos.DefenseStackDTOFromDomain(s))
	}
	armiesJSON, _ := json.Marshal(armyDTOs)
	buildingsJSON, _ := json.Marshal(structsDTO)
	return gen.InsertResourceLocationParams{
		SectorX:                int32(loc.Coordinates.X),
		SectorY:                int32(loc.Coordinates.Y),
		Type:                   loc.Type,
		Amount:                 int32(loc.Amount),
		Name:                   toNullString(loc.Name),
		Description:            toNullString(loc.Description),
		ImageUrl:               toNullString(loc.ImageURL),
		Resources:              resJSON,
		ResourcesCalcTimestamp: loc.Resources.CalculationTimestamp,
		Armies:                 armiesJSON,
		Buildings:              buildingsJSON,
	}
}

func UpdateResourceLocationParamsFromDomain(loc *domain.ResourceLocationModel) gen.UpdateResourceLocationParams {
	resJSON, _ := json.Marshal(dtos.LocationResourceStatsDTOFromDomain(loc.Resources))
	armyDTOs := make([]dtos.ArmyStackDTO, 0, len(loc.DefendingArmies))
	for _, a := range loc.DefendingArmies {
		armyDTOs = append(armyDTOs, dtos.ArmyStackDTOFromDomain(a))
	}
	structsDTO := make([]dtos.DefenseStackDTO, 0, len(loc.DefendingStructures))
	for _, s := range loc.DefendingStructures {
		structsDTO = append(structsDTO, dtos.DefenseStackDTOFromDomain(s))
	}
	armiesJSON, _ := json.Marshal(armyDTOs)
	buildingsJSON, _ := json.Marshal(structsDTO)
	return gen.UpdateResourceLocationParams{
		ID:                     int64(loc.ID),
		Type:                   loc.Type,
		Amount:                 int32(loc.Amount),
		Name:                   toNullString(loc.Name),
		Description:            toNullString(loc.Description),
		ImageUrl:               toNullString(loc.ImageURL),
		Resources:              resJSON,
		ResourcesCalcTimestamp: loc.Resources.CalculationTimestamp,
		Armies:                 armiesJSON,
		Buildings:              buildingsJSON,
	}
}
