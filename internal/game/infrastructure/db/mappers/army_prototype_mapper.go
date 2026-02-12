package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-api/internal/game/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/db/gen"
)

func ArmyPrototypeFromDB(p gen.ArmyItemPrototype) *domain.ArmyItemPrototype {
	var priceDTO dtos.PriceDTO
	_ = json.Unmarshal(p.Price, &priceDTO)
	price := dtos.PriceFromDTO(priceDTO)

	proto := &domain.ArmyItemPrototype{
		ID:                 int(p.ID),
		Name:               p.Name,
		Category:           domain.ArmyCategory(p.Category),
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
