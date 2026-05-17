package readmodels

import "github.com/artcodefun/heat-expansion-server/internal/game/domain"

type ResourceType string

const (
	ResourceTypeCredits    ResourceType = "CREDITS"
	ResourceTypeIron       ResourceType = "IRON"
	ResourceTypeTitanium   ResourceType = "TITANIUM"
	ResourceTypeAntimatter ResourceType = "ANTIMATTER"
)

type BlackMarketResourceRate struct {
	ResourceType     ResourceType
	AmountPerCrystal int
}

func BlackMarketResourceRateFromDomain(rate domain.BlackMarketResourceRate) *BlackMarketResourceRate {
	return &BlackMarketResourceRate{
		ResourceType:     ResourceType(rate.ResourceType),
		AmountPerCrystal: rate.AmountPerCrystal,
	}
}

func BlackMarketResourceRateListFromDomain(rates []domain.BlackMarketResourceRate) []*BlackMarketResourceRate {
	items := make([]*BlackMarketResourceRate, 0, len(rates))
	for _, rate := range rates {
		items = append(items, BlackMarketResourceRateFromDomain(rate))
	}
	return items
}
