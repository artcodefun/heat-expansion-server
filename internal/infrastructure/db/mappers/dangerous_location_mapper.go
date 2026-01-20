package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/gen"
)

func DangerousLocationFromDB(r gen.DangerousLocation, armyProtos map[int]*domain.ArmyItemPrototype, buildProtos map[int]*domain.BuildItemPrototype) *domain.DangerousLocationModel {
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

	var trophyDTOs []dtos.TrophyDTO
	_ = json.Unmarshal(r.Trophies, &trophyDTOs)
	trophies := make([]domain.TrophyStorageItem, 0, len(trophyDTOs))
	for _, d := range trophyDTOs {
		trophies = append(trophies, dtos.TrophyFromDTO(d))
	}

	return &domain.DangerousLocationModel{
		ID:              int(r.ID),
		Coordinates:     domain.Vector2i{X: int(r.SectorX), Y: int(r.SectorY)},
		DefenderFaction: domain.Faction(r.DefenderFaction),
		TotalWorth:      int(r.TotalWorth),
		LocationDetails: domain.LocationDetails{
			Name:        nullStringToString(&r.Name.String, r.Name.Valid),
			Description: nullStringToString(&r.Description.String, r.Description.Valid),
			ImageURL:    nullStringToString(&r.ImageUrl.String, r.ImageUrl.Valid),
		},
		Resources:           res,
		DefendingArmies:     armies,
		DefendingStructures: structs,
		Trophies:            trophies,
	}
}

func InsertDangerousLocationParamsFromDomain(loc *domain.DangerousLocationModel) gen.InsertDangerousLocationParams {
	resJSON, _ := json.Marshal(dtos.LocationResourceStatsDTOFromDomain(loc.Resources))
	armyDTOs := make([]dtos.ArmyStackDTO, 0, len(loc.DefendingArmies))
	for _, a := range loc.DefendingArmies {
		armyDTOs = append(armyDTOs, dtos.ArmyStackDTOFromDomain(a))
	}
	structsDTO := make([]dtos.DefenseStackDTO, 0, len(loc.DefendingStructures))
	for _, s := range loc.DefendingStructures {
		structsDTO = append(structsDTO, dtos.DefenseStackDTOFromDomain(s))
	}
	trophyDTOs := make([]dtos.TrophyDTO, 0, len(loc.Trophies))
	for _, t := range loc.Trophies {
		trophyDTOs = append(trophyDTOs, dtos.TrophyDTOFromDomain(t))
	}
	armiesJSON, _ := json.Marshal(armyDTOs)
	buildingsJSON, _ := json.Marshal(structsDTO)
	trophiesJSON, _ := json.Marshal(trophyDTOs)
	return gen.InsertDangerousLocationParams{
		SectorX:                int32(loc.Coordinates.X),
		SectorY:                int32(loc.Coordinates.Y),
		DefenderFaction:        string(loc.DefenderFaction),
		TotalWorth:             int32(loc.TotalWorth),
		Name:                   toNullString(loc.Name),
		Description:            toNullString(loc.Description),
		ImageUrl:               toNullString(loc.ImageURL),
		Resources:              resJSON,
		ResourcesCalcTimestamp: loc.Resources.CalculationTimestamp,
		Armies:                 armiesJSON,
		Buildings:              buildingsJSON,
		Trophies:               trophiesJSON,
	}
}

func UpdateDangerousLocationParamsFromDomain(loc *domain.DangerousLocationModel) gen.UpdateDangerousLocationParams {
	resJSON, _ := json.Marshal(dtos.LocationResourceStatsDTOFromDomain(loc.Resources))
	armyDTOs := make([]dtos.ArmyStackDTO, 0, len(loc.DefendingArmies))
	for _, a := range loc.DefendingArmies {
		armyDTOs = append(armyDTOs, dtos.ArmyStackDTOFromDomain(a))
	}
	structsDTO := make([]dtos.DefenseStackDTO, 0, len(loc.DefendingStructures))
	for _, s := range loc.DefendingStructures {
		structsDTO = append(structsDTO, dtos.DefenseStackDTOFromDomain(s))
	}
	trophyDTOs := make([]dtos.TrophyDTO, 0, len(loc.Trophies))
	for _, t := range loc.Trophies {
		trophyDTOs = append(trophyDTOs, dtos.TrophyDTOFromDomain(t))
	}
	armiesJSON, _ := json.Marshal(armyDTOs)
	buildingsJSON, _ := json.Marshal(structsDTO)
	trophiesJSON, _ := json.Marshal(trophyDTOs)
	return gen.UpdateDangerousLocationParams{
		ID:                     int64(loc.ID),
		DefenderFaction:        string(loc.DefenderFaction),
		TotalWorth:             int32(loc.TotalWorth),
		Name:                   toNullString(loc.Name),
		Description:            toNullString(loc.Description),
		ImageUrl:               toNullString(loc.ImageURL),
		Resources:              resJSON,
		ResourcesCalcTimestamp: loc.Resources.CalculationTimestamp,
		Armies:                 armiesJSON,
		Buildings:              buildingsJSON,
		Trophies:               trophiesJSON,
	}
}
