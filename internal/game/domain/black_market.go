package domain

import "fmt"

const (
	BlackMarketCreditsAmountPerCrystal    = 100
	BlackMarketIronAmountPerCrystal       = 25
	BlackMarketTitaniumAmountPerCrystal   = 5
	BlackMarketAntimatterAmountPerCrystal = 1
)

// BlackMarketResourceRate defines the Black Market exchange rate for a resource.
type BlackMarketResourceRate struct {
	ResourceType     ResourceType
	AmountPerCrystal int
}

func ListBlackMarketResourceRates() []BlackMarketResourceRate {
	return []BlackMarketResourceRate{
		{ResourceType: ResourceTypeCredits, AmountPerCrystal: BlackMarketCreditsAmountPerCrystal},
		{ResourceType: ResourceTypeIron, AmountPerCrystal: BlackMarketIronAmountPerCrystal},
		{ResourceType: ResourceTypeTitanium, AmountPerCrystal: BlackMarketTitaniumAmountPerCrystal},
		{ResourceType: ResourceTypeAntimatter, AmountPerCrystal: BlackMarketAntimatterAmountPerCrystal},
	}
}

func BlackMarketResourceRateFor(resourceType ResourceType) (BlackMarketResourceRate, error) {
	for _, rate := range ListBlackMarketResourceRates() {
		if rate.ResourceType == resourceType {
			return rate, nil
		}
	}
	return BlackMarketResourceRate{}, NewError("error.domain.black_market.rate_not_available", H{"resource": fmt.Sprint(resourceType)})
}

func (r BlackMarketResourceRate) ResourcesForCrystals(crystals int) PriceModel {
	amount := r.AmountPerCrystal * crystals
	switch r.ResourceType {
	case ResourceTypeCredits:
		return PriceModel{Credits: amount}
	case ResourceTypeIron:
		return PriceModel{Iron: amount}
	case ResourceTypeTitanium:
		return PriceModel{Titanium: amount}
	case ResourceTypeAntimatter:
		return PriceModel{Antimatter: amount}
	default:
		return PriceModel{}
	}
}
