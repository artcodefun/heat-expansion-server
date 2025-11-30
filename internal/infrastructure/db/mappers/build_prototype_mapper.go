package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/gen"
)

func BuildPrototypeFromDB(p gen.BuildItemPrototype) *domain.BuildItemPrototype {
	var price domain.PriceModel
	_ = json.Unmarshal(p.Price, &price)

	var ctrl *domain.ControlBuildingData
	if p.ControlData.Valid {
		var tmp domain.ControlBuildingData
		unmarshalIfValid(p.ControlData, &tmp)
		ctrl = &tmp
	}
	var res *domain.ResourcesBuildingData
	if p.ResourcesData.Valid {
		var tmp domain.ResourcesBuildingData
		unmarshalIfValid(p.ResourcesData, &tmp)
		res = &tmp
	}
	var def *domain.DefenseBuildingData
	if p.DefenseData.Valid {
		var tmp domain.DefenseBuildingData
		unmarshalIfValid(p.DefenseData, &tmp)
		def = &tmp
	}
	var mil *domain.MilitaryBuildingData
	if p.MilitaryData.Valid {
		var tmp domain.MilitaryBuildingData
		unmarshalIfValid(p.MilitaryData, &tmp)
		mil = &tmp
	}
	var intel *domain.IntelligenceBuildingData
	if p.IntelligenceData.Valid {
		var tmp domain.IntelligenceBuildingData
		unmarshalIfValid(p.IntelligenceData, &tmp)
		intel = &tmp
	}

	proto := &domain.BuildItemPrototype{
		ID:                 int(p.ID),
		Name:               p.Name,
		Category:           domain.BuildCategory(p.Category),
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
