package dtos

import (
	readmodels "github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
)

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

type BlackMarketOfferKind string

const (
	BlackMarketOfferKindBuilding BlackMarketOfferKind = "BUILDING"
	BlackMarketOfferKindArmy     BlackMarketOfferKind = "ARMY"
	BlackMarketOfferKindStorage  BlackMarketOfferKind = "STORAGE"
)

type BlackMarketOfferDTO struct {
	ID              int64                    `json:"id"`
	Kind            BlackMarketOfferKind     `json:"kind"`
	PrototypeID     int                      `json:"prototype_id"`
	PriceInCrystals int                      `json:"price_in_crystals"`
	EndsAt          *int64                   `json:"ends_at,omitempty"`
	IsLimited       bool                     `json:"is_limited"`
	Priority        int                      `json:"priority"`
	Building        *BuildItemPrototypeDTO   `json:"building,omitempty"`
	Army            *ArmyItemPrototypeDTO    `json:"army,omitempty"`
	Storage         *StorageItemPrototypeDTO `json:"storage,omitempty"`
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

func BlackMarketOfferDTOFromReadModel(item *readmodels.BlackMarketOffer, tr ports.Translator, locale string) BlackMarketOfferDTO {
	dto := BlackMarketOfferDTO{
		ID:              item.ID,
		Kind:            BlackMarketOfferKind(item.Kind),
		PrototypeID:     item.PrototypeID,
		PriceInCrystals: item.PriceInCrystals,
		EndsAt:          item.EndsAt,
		IsLimited:       item.IsLimited,
		Priority:        item.Priority,
	}
	if item.Building != nil {
		building := mapBuildItemPrototype(*item.Building, tr, locale)
		dto.Building = &building
	}
	if item.Army != nil {
		army := mapArmyPrototype(*item.Army, tr, locale)
		dto.Army = &army
	}
	if item.Storage != nil {
		storage := mapStorageItemPrototype(*item.Storage, tr, locale)
		dto.Storage = &storage
	}
	return dto
}

func BlackMarketOfferDTOListFromReadModels(items []*readmodels.BlackMarketOffer, tr ports.Translator, locale string) []BlackMarketOfferDTO {
	out := make([]BlackMarketOfferDTO, 0, len(items))
	for _, item := range items {
		out = append(out, BlackMarketOfferDTOFromReadModel(item, tr, locale))
	}
	return out
}
