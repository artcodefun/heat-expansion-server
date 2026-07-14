package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
)

func TechPrototypeFromDB(p gen.TechItemPrototype) *domain.TechItemPrototype {
	var priceDTO dtos.PriceDTO
	_ = json.Unmarshal(p.Price, &priceDTO)
	price := dtos.PriceFromDTO(priceDTO)

	var improvement *domain.TechImprovement
	if p.Improvement.Valid {
		var impDTO dtos.TechImprovementDTO
		_ = json.Unmarshal(p.Improvement.RawMessage, &impDTO)
		imp := dtos.TechImprovementFromDTO(&impDTO)
		improvement = imp
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
		Improvement:        improvement,
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

// TechPrototypeToCreateParams maps a domain prototype to the sqlc insert params.
// improvement is serialized via its DTO shape; a nil pointer becomes SQL NULL.
func TechPrototypeToCreateParams(p *domain.TechItemPrototype) gen.CreateTechPrototypeParams {
	return gen.CreateTechPrototypeParams{
		ID:                 int64(p.ID),
		Name:               string(p.Name),
		Category:           string(p.Category),
		UnlockTechnologyID: nullableBaseID(p.UnlockTechnologyID),
		ShortDescription:   stringToNullString(string(p.ShortDescription)),
		FullDescription:    stringToNullString(string(p.FullDescription)),
		Price:              priceToJSON(p.Price),
		ResearchTime:       p.ResearchTime,
		ImageUrl:           stringToNullString(p.ImageURL),
		Improvement:        toNullRawMessage(dtos.TechImprovementDTOFromDomain(p.Improvement)),
	}
}

// TechPrototypeToUpdateParams maps a domain prototype to the sqlc update params,
// keyed by p.ID.
func TechPrototypeToUpdateParams(p *domain.TechItemPrototype) gen.UpdateTechPrototypeParams {
	return gen.UpdateTechPrototypeParams{
		ID:                 int64(p.ID),
		Name:               string(p.Name),
		Category:           string(p.Category),
		UnlockTechnologyID: nullableBaseID(p.UnlockTechnologyID),
		ShortDescription:   stringToNullString(string(p.ShortDescription)),
		FullDescription:    stringToNullString(string(p.FullDescription)),
		Price:              priceToJSON(p.Price),
		ResearchTime:       p.ResearchTime,
		ImageUrl:           stringToNullString(p.ImageURL),
		Improvement:        toNullRawMessage(dtos.TechImprovementDTOFromDomain(p.Improvement)),
	}
}
