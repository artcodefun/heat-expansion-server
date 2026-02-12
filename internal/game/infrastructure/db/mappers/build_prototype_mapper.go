package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-api/internal/game/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/db/gen"
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
