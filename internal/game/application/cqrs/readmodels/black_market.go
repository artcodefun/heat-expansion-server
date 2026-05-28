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
	ResourceType   ResourceType
	CrystalsCost   int
	ResourceAmount int
}

type BlackMarketOfferKind string

const (
	BlackMarketOfferKindBuilding BlackMarketOfferKind = "BUILDING"
	BlackMarketOfferKindArmy     BlackMarketOfferKind = "ARMY"
	BlackMarketOfferKindStorage  BlackMarketOfferKind = "STORAGE"
)

type BlackMarketOffer struct {
	ID              int64
	Kind            BlackMarketOfferKind
	PrototypeID     int
	PriceInCrystals int
	EndsAt          *int64
	IsLimited       bool
	Priority        int
	Building        *BuildItemPrototype
	Army            *ArmyItemPrototype
	Storage         *StorageItemPrototype
}

func BlackMarketResourceRateFromDomain(rate domain.BlackMarketResourceRate) *BlackMarketResourceRate {
	return &BlackMarketResourceRate{
		ResourceType:   ResourceType(rate.ResourceType),
		CrystalsCost:   rate.CrystalsCost,
		ResourceAmount: rate.ResourceAmount,
	}
}

func BlackMarketResourceRateListFromDomain(rates []domain.BlackMarketResourceRate) []*BlackMarketResourceRate {
	items := make([]*BlackMarketResourceRate, 0, len(rates))
	for _, rate := range rates {
		items = append(items, BlackMarketResourceRateFromDomain(rate))
	}
	return items
}
