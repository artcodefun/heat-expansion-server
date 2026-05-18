package dtos

import readmodels "github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"

type ResourceType string

// ResourceType enum values
const (
	ResourceTypeCredits    ResourceType = "CREDITS"
	ResourceTypeIron       ResourceType = "IRON"
	ResourceTypeTitanium   ResourceType = "TITANIUM"
	ResourceTypeAntimatter ResourceType = "ANTIMATTER"
)

type BlackMarketResourceRateDTO struct {
	ResourceType   ResourceType `json:"resource_type"`
	CrystalsCost   int          `json:"crystals_cost"`
	ResourceAmount int          `json:"resource_amount"`
}

func BlackMarketResourceRateDTOFromReadModel(rate *readmodels.BlackMarketResourceRate) BlackMarketResourceRateDTO {
	return BlackMarketResourceRateDTO{
		ResourceType:   ResourceType(rate.ResourceType),
		CrystalsCost:   rate.CrystalsCost,
		ResourceAmount: rate.ResourceAmount,
	}
}

func BlackMarketResourceRateDTOListFromReadModels(rates []*readmodels.BlackMarketResourceRate) []BlackMarketResourceRateDTO {
	items := make([]BlackMarketResourceRateDTO, 0, len(rates))
	for _, rate := range rates {
		items = append(items, BlackMarketResourceRateDTOFromReadModel(rate))
	}
	return items
}
