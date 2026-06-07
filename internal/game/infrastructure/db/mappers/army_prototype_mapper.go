package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
)

func ArmyPrototypeFromDB(p gen.ArmyItemPrototype) *domain.ArmyItemPrototype {
	var priceDTO dtos.PriceDTO
	_ = json.Unmarshal(p.Price, &priceDTO)
	price := dtos.PriceFromDTO(priceDTO)

	proto := &domain.ArmyItemPrototype{
		ID:                 int(p.ID),
		Name:               p.Name,
		Category:           domain.ArmyCategory(p.Category),
		CreationSources:    creationSourcesFromJSON(p.CreationSources),
		Faction:            domain.Faction(p.Faction),
		UnlockTechnologyID: nullableIntPtr(p.UnlockTechnologyID.Int64, p.UnlockTechnologyID.Valid),
		ShortDescription:   nullStringToString(&p.ShortDescription.String, p.ShortDescription.Valid),
		FullDescription:    nullStringToString(&p.FullDescription.String, p.FullDescription.Valid),
		Price:              price,
		ProductionTime:     p.ProductionTime,
		Space:              int(p.Space),
		ImageURL:           nullStringToString(&p.ImageUrl.String, p.ImageUrl.Valid),
		Attack:             int(p.Attack),
		Defence:            int(p.Defence),
		Capacity:           int(p.Capacity),
		Stealth:            int(p.Stealth),
		Speed:              int(p.Speed),
	}
	return proto
}

func ArmyPrototypesFromDB(src []gen.ArmyItemPrototype) []*domain.ArmyItemPrototype {
	dst := make([]*domain.ArmyItemPrototype, len(src))
	for i, p := range src {
		dst[i] = ArmyPrototypeFromDB(p)
	}
	return dst
}

// ArmyPrototypeToCreateParams maps a domain prototype to the sqlc insert params.
func ArmyPrototypeToCreateParams(p *domain.ArmyItemPrototype) gen.CreateArmyPrototypeParams {
	return gen.CreateArmyPrototypeParams{
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
		Attack:             int32(p.Attack),
		Defence:            int32(p.Defence),
		Capacity:           int32(p.Capacity),
		Stealth:            int32(p.Stealth),
		Speed:              int32(p.Speed),
		CreationSources:    creationSourcesToJSON(p.CreationSources),
	}
}

// ArmyPrototypeToUpdateParams maps a domain prototype to the sqlc update params,
// keyed by p.ID.
func ArmyPrototypeToUpdateParams(p *domain.ArmyItemPrototype) gen.UpdateArmyPrototypeParams {
	return gen.UpdateArmyPrototypeParams{
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
		Attack:             int32(p.Attack),
		Defence:            int32(p.Defence),
		Capacity:           int32(p.Capacity),
		Stealth:            int32(p.Stealth),
		Speed:              int32(p.Speed),
		CreationSources:    creationSourcesToJSON(p.CreationSources),
	}
}
