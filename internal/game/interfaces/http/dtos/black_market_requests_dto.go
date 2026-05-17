package dtos

import (
	"strings"

	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
)

type blackMarketResourcesPurchaseBody struct {
	ResourceType string `json:"resource_type" binding:"required,resource_type"`
	Crystals     int    `json:"crystals" binding:"required,min=1"`
}

type BlackMarketResourcesPurchaseRequest = Request[BaseURI, None, blackMarketResourcesPurchaseBody]
type BlackMarketResourceRatesRequest = Request[BaseURI, None, None]

func IsValidResourceType(value string) bool {
	upper := strings.ToUpper(value)
	switch ResourceType(upper) {
	case ResourceTypeCredits, ResourceTypeIron, ResourceTypeTitanium, ResourceTypeAntimatter:
		return true
	default:
		return false
	}
}

func ResourceTypeFromDTO(value string) domain.ResourceType {
	return domain.ResourceType(ResourceType(strings.ToUpper(strings.TrimSpace(string(value)))))
}
