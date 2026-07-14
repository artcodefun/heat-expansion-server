package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
)

func BuildPrototypeFromDB(p gen.BuildItemPrototype) *domain.BuildItemPrototype {
	var priceDTO dtos.PriceDTO
	_ = json.Unmarshal(p.Price, &priceDTO)
	price := dtos.PriceFromDTO(priceDTO)

	var ctrl *domain.ControlBuildingData
	if p.ControlData.Valid {
		var dto dtos.ControlBuildingDataDTO
		unmarshalIfValid(p.ControlData, &dto)
		ctrl = dtos.ControlBuildingDataFromDTO(&dto)
	}
	var res *domain.ResourcesBuildingData
	if p.ResourcesData.Valid {
		var dto dtos.ResourcesBuildingDataDTO
		unmarshalIfValid(p.ResourcesData, &dto)
		res = dtos.ResourcesBuildingDataFromDTO(&dto)
	}
	var def *domain.DefenseBuildingData
	if p.DefenseData.Valid {
		var dto dtos.DefenseBuildingDataDTO
		unmarshalIfValid(p.DefenseData, &dto)
		def = dtos.DefenseBuildingDataFromDTO(&dto)
	}
	var mil *domain.MilitaryBuildingData
	if p.MilitaryData.Valid {
		var dto dtos.MilitaryBuildingDataDTO
		unmarshalIfValid(p.MilitaryData, &dto)
		mil = dtos.MilitaryBuildingDataFromDTO(&dto)
	}
	var intel *domain.IntelligenceBuildingData
	if p.IntelligenceData.Valid {
		var dto dtos.IntelligenceBuildingDataDTO
		unmarshalIfValid(p.IntelligenceData, &dto)
		intel = dtos.IntelligenceBuildingDataFromDTO(&dto)
	}

	proto := &domain.BuildItemPrototype{
		ID:                 int(p.ID),
		Name:               p.Name,
		Category:           domain.BuildCategory(p.Category),
		CreationSources:    creationSourcesFromJSON(p.CreationSources),
		Faction:            domain.Faction(p.Faction),
		UnlockTechnologyID: nullableIntPtr(p.UnlockTechnologyID.Int64, p.UnlockTechnologyID.Valid),
		ShortDescription:   nullStringToString(&p.ShortDescription.String, p.ShortDescription.Valid),
		FullDescription:    nullStringToString(&p.FullDescription.String, p.FullDescription.Valid),
		Price:              price,
		ProductionTime:     p.ProductionTime,
		Space:              int(p.Space),
		ImageURL:           nullStringToString(&p.ImageUrl.String, p.ImageUrl.Valid),
		ControlData:        ctrl,
		ResourcesData:      res,
		DefenseData:        def,
		MilitaryData:       mil,
		IntelligenceData:   intel,
	}
	return proto
}

func BuildPrototypesFromDB(src []gen.BuildItemPrototype) []*domain.BuildItemPrototype {
	dst := make([]*domain.BuildItemPrototype, len(src))
	for i, p := range src {
		dst[i] = BuildPrototypeFromDB(p)
	}
	return dst
}

// BuildPrototypeToCreateParams maps a domain prototype to the sqlc insert params.
// The category-specific data blocks are serialized via their DTO shapes; a nil
// pointer becomes SQL NULL.
func BuildPrototypeToCreateParams(p *domain.BuildItemPrototype) gen.CreateBuildPrototypeParams {
	return gen.CreateBuildPrototypeParams{
		ID:                 int64(p.ID),
		Name:               string(p.Name),
		Category:           string(p.Category),
		Faction:            string(p.Faction),
		UnlockTechnologyID: nullableBaseID(p.UnlockTechnologyID),
		ShortDescription:   stringToNullString(string(p.ShortDescription)),
		FullDescription:    stringToNullString(string(p.FullDescription)),
		Price:              priceToJSON(p.Price),
		ProductionTime:     p.ProductionTime,
		Space:              int32(p.Space),
		ImageUrl:           stringToNullString(p.ImageURL),
		ControlData:        toNullRawMessage(dtos.ControlBuildingDataDTOFromDomain(p.ControlData)),
		ResourcesData:      toNullRawMessage(dtos.ResourcesBuildingDataDTOFromDomain(p.ResourcesData)),
		DefenseData:        toNullRawMessage(dtos.DefenseBuildingDataDTOFromDomain(p.DefenseData)),
		MilitaryData:       toNullRawMessage(dtos.MilitaryBuildingDataDTOFromDomain(p.MilitaryData)),
		IntelligenceData:   toNullRawMessage(dtos.IntelligenceBuildingDataDTOFromDomain(p.IntelligenceData)),
		CreationSources:    creationSourcesToJSON(p.CreationSources),
	}
}

// BuildPrototypeToUpdateParams maps a domain prototype to the sqlc update params,
// keyed by p.ID.
func BuildPrototypeToUpdateParams(p *domain.BuildItemPrototype) gen.UpdateBuildPrototypeParams {
	return gen.UpdateBuildPrototypeParams{
		ID:                 int64(p.ID),
		Name:               string(p.Name),
		Category:           string(p.Category),
		Faction:            string(p.Faction),
		UnlockTechnologyID: nullableBaseID(p.UnlockTechnologyID),
		ShortDescription:   stringToNullString(string(p.ShortDescription)),
		FullDescription:    stringToNullString(string(p.FullDescription)),
		Price:              priceToJSON(p.Price),
		ProductionTime:     p.ProductionTime,
		Space:              int32(p.Space),
		ImageUrl:           stringToNullString(p.ImageURL),
		ControlData:        toNullRawMessage(dtos.ControlBuildingDataDTOFromDomain(p.ControlData)),
		ResourcesData:      toNullRawMessage(dtos.ResourcesBuildingDataDTOFromDomain(p.ResourcesData)),
		DefenseData:        toNullRawMessage(dtos.DefenseBuildingDataDTOFromDomain(p.DefenseData)),
		MilitaryData:       toNullRawMessage(dtos.MilitaryBuildingDataDTOFromDomain(p.MilitaryData)),
		IntelligenceData:   toNullRawMessage(dtos.IntelligenceBuildingDataDTOFromDomain(p.IntelligenceData)),
		CreationSources:    creationSourcesToJSON(p.CreationSources),
	}
}
