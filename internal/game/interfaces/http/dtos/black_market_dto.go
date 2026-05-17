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
	ResourceType     ResourceType `json:"resource_type"`
	AmountPerCrystal int          `json:"amount_per_crystal"`
}

func BlackMarketResourceRateDTOFromReadModel(rate *readmodels.BlackMarketResourceRate) BlackMarketResourceRateDTO {
	return BlackMarketResourceRateDTO{
		ResourceType:     ResourceType(rate.ResourceType),
		AmountPerCrystal: rate.AmountPerCrystal,
	}
}

func BlackMarketResourceRateDTOListFromReadModels(rates []*readmodels.BlackMarketResourceRate) []BlackMarketResourceRateDTO {
	items := make([]BlackMarketResourceRateDTO, 0, len(rates))
	for _, rate := range rates {
		items = append(items, BlackMarketResourceRateDTOFromReadModel(rate))
	}
	return items
}
