package domain

import "fmt"

const (
	BlackMarketCreditsCrystalsPerUnit    = 1
	BlackMarketCreditsAmountPerUnit      = 100
	BlackMarketIronCrystalsPerUnit       = 1
	BlackMarketIronAmountPerUnit         = 25
	BlackMarketTitaniumCrystalsPerUnit   = 1
	BlackMarketTitaniumAmountPerUnit     = 5
	BlackMarketAntimatterCrystalsPerUnit = 10
	BlackMarketAntimatterAmountPerUnit   = 1
)

// BlackMarketResourceRate defines the Black Market exchange rate for a resource.
type BlackMarketResourceRate struct {
	ResourceType   ResourceType
	CrystalsCost   int
	ResourceAmount int
}

func ListBlackMarketResourceRates() []BlackMarketResourceRate {
	return []BlackMarketResourceRate{
		{ResourceType: ResourceTypeCredits, CrystalsCost: BlackMarketCreditsCrystalsPerUnit, ResourceAmount: BlackMarketCreditsAmountPerUnit},
		{ResourceType: ResourceTypeIron, CrystalsCost: BlackMarketIronCrystalsPerUnit, ResourceAmount: BlackMarketIronAmountPerUnit},
		{ResourceType: ResourceTypeTitanium, CrystalsCost: BlackMarketTitaniumCrystalsPerUnit, ResourceAmount: BlackMarketTitaniumAmountPerUnit},
		{ResourceType: ResourceTypeAntimatter, CrystalsCost: BlackMarketAntimatterCrystalsPerUnit, ResourceAmount: BlackMarketAntimatterAmountPerUnit},
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

func (r BlackMarketResourceRate) ResourceAmountForCrystals(crystals int) (int, error) {
	if crystals <= 0 || crystals%r.CrystalsCost != 0 {
		return 0, NewError("error.domain.black_market.invalid_crystal_amount", H{"resource": fmt.Sprint(r.ResourceType), "crystals_cost": r.CrystalsCost})
	}
	return (crystals / r.CrystalsCost) * r.ResourceAmount, nil
}
