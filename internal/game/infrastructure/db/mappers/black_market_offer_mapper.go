package mappers

import (
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
)

func BlackMarketOfferFromDB(row gen.BlackMarketOffer) *domain.BlackMarketOffer {
	return &domain.BlackMarketOffer{
		ID:              row.ID,
		Kind:            domain.BlackMarketOfferKind(row.Kind),
		PrototypeID:     int(row.PrototypeID),
		PriceInCrystals: int(row.PriceInCrystals),
		EndsAt:          nullInt64ToInt64Ptr(row.EndsAt),
		IsLimited:       row.IsLimited,
		Priority:        int(row.Priority),
	}
}

func BlackMarketOffersFromDB(rows []gen.BlackMarketOffer) []*domain.BlackMarketOffer {
	items := make([]*domain.BlackMarketOffer, 0, len(rows))
	for _, row := range rows {
		items = append(items, BlackMarketOfferFromDB(row))
	}
	return items
}

func InsertBlackMarketOfferParamsFromDomain(offer *domain.BlackMarketOffer) gen.InsertBlackMarketOfferParams {
	return gen.InsertBlackMarketOfferParams{
		Kind:            string(offer.Kind),
		PrototypeID:     int64(offer.PrototypeID),
		PriceInCrystals: int64(offer.PriceInCrystals),
		EndsAt:          int64PtrToNullInt64(offer.EndsAt),
		IsLimited:       offer.IsLimited,
		Priority:        int64(offer.Priority),
	}
}

func UpdateBlackMarketOfferParamsFromDomain(offer *domain.BlackMarketOffer) gen.UpdateBlackMarketOfferParams {
	return gen.UpdateBlackMarketOfferParams{
		ID:              offer.ID,
		Kind:            string(offer.Kind),
		PrototypeID:     int64(offer.PrototypeID),
		PriceInCrystals: int64(offer.PriceInCrystals),
		EndsAt:          int64PtrToNullInt64(offer.EndsAt),
		IsLimited:       offer.IsLimited,
		Priority:        int64(offer.Priority),
	}
}
