package dtos

import (
	"strings"

	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
)

// blackMarketResourcesPurchaseBody represents the JSON payload required to buy Black Market resources.
type blackMarketResourcesPurchaseBody struct {
	ResourceType string `json:"resource_type" binding:"required,resource_type"`
	Crystals     int    `json:"crystals" binding:"required,min=1"`
}

// blackMarketOffersQuery contains query params for BlackMarketOffersListRequest.
type blackMarketOffersQuery struct {
	Kind    string `form:"kind" binding:"omitempty,black_market_offer_kind"`
	Limited *bool  `form:"limited"`
}

type BlackMarketOfferURI struct {
	BaseID  int   `uri:"baseId" binding:"required,min=1"`
	OfferID int64 `uri:"offerId" binding:"required,min=1"`
}

// blackMarketOfferPurchaseBody represents the JSON payload for purchasing a Black Market offer.
type blackMarketOfferPurchaseBody struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

type BlackMarketResourcesPurchaseRequest = Request[BaseURI, None, blackMarketResourcesPurchaseBody]
type BlackMarketResourceRatesRequest = Request[BaseURI, None, None]
type BlackMarketOffersListRequest = Request[BaseURI, blackMarketOffersQuery, None]
type BlackMarketOfferPurchaseRequest = Request[BlackMarketOfferURI, None, blackMarketOfferPurchaseBody]

func IsValidResourceType(value string) bool {
	upper := strings.ToUpper(value)
	switch ResourceType(upper) {
	case ResourceTypeCredits, ResourceTypeIron, ResourceTypeTitanium, ResourceTypeAntimatter:
		return true
	default:
		return false
	}
}

func IsValidBlackMarketOfferKind(value string) bool {
	upper := strings.ToUpper(strings.TrimSpace(value))
	switch BlackMarketOfferKind(upper) {
	case BlackMarketOfferKindBuilding, BlackMarketOfferKindArmy, BlackMarketOfferKindStorage:
		return true
	default:
		return false
	}
}

func ResourceTypeFromDTO(value string) domain.ResourceType {
	return domain.ResourceType(ResourceType(strings.ToUpper(strings.TrimSpace(string(value)))))
}

func BlackMarketOfferKindPtrFromDTO(value string) *domain.BlackMarketOfferKind {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	kind := domain.BlackMarketOfferKind(strings.ToUpper(trimmed))
	return &kind
}
