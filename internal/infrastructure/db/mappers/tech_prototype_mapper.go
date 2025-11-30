package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/gen"
)

func TechPrototypeFromDB(p gen.TechItemPrototype) *domain.TechItemPrototype {
	var priceDTO dtos.PriceDTO
	_ = json.Unmarshal(p.Price, &priceDTO)
	price := dtos.PriceFromDTO(priceDTO)

	var effDTOs []dtos.TechnologyEffectDTO
	_ = json.Unmarshal(p.Effects, &effDTOs)
	effects := make([]domain.TechnologyEffect, 0, len(effDTOs))
	for _, d := range effDTOs {
		effects = append(effects, dtos.TechnologyEffectFromDTO(d))
	}

	proto := &domain.TechItemPrototype{
		ID:                 int(p.ID),
		Name:               p.Name,
		Category:           domain.TechCategory(p.Category),
		UnlockTechnologyID: nullableIntPtr(p.UnlockTechnologyID.Int64, p.UnlockTechnologyID.Valid),
		ShortDescription:   nullStringToString(&p.ShortDescription.String, p.ShortDescription.Valid),
		FullDescription:    nullStringToString(&p.FullDescription.String, p.FullDescription.Valid),
		Price:              price,
		ResearchTime:       p.ResearchTime,
		ImageURL:           nullStringToString(&p.ImageUrl.String, p.ImageUrl.Valid),
		Effects:            effects,
	}
	return proto
}

func TechPrototypesFromDB(src []gen.TechItemPrototype) []*domain.TechItemPrototype {
	dst := make([]*domain.TechItemPrototype, len(src))
	for i, p := range src {
		dst[i] = TechPrototypeFromDB(p)
	}
	return dst
}
