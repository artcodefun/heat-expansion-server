package mappers

import (
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
)

func BlackMarketBuildingOfferFromRow(r gen.ListActiveBlackMarketBuildingOffersRow) readmodels.BlackMarketOffer {
	return readmodels.BlackMarketOffer{
		ID:              r.ID,
		Kind:            readmodels.BlackMarketOfferKind(r.Kind),
		PrototypeID:     int(r.ProtoID),
		PriceInCrystals: int(r.PriceInCrystals),
		EndsAt:          nullInt64ToInt64Ptr(r.EndsAt),
		IsLimited:       r.IsLimited,
		Priority:        int(r.Priority),
		Building:        blackMarketBuildPrototypeFromRow(r),
	}
}

func BlackMarketArmyOfferFromRow(r gen.ListActiveBlackMarketArmyOffersRow) readmodels.BlackMarketOffer {
	return readmodels.BlackMarketOffer{
		ID:              r.ID,
		Kind:            readmodels.BlackMarketOfferKind(r.Kind),
		PrototypeID:     int(r.ProtoID),
		PriceInCrystals: int(r.PriceInCrystals),
		EndsAt:          nullInt64ToInt64Ptr(r.EndsAt),
		IsLimited:       r.IsLimited,
		Priority:        int(r.Priority),
		Army:            blackMarketArmyPrototypeFromRow(r),
	}
}

func BlackMarketStorageOfferFromRow(r gen.ListActiveBlackMarketStorageOffersRow) readmodels.BlackMarketOffer {
	return readmodels.BlackMarketOffer{
		ID:              r.ID,
		Kind:            readmodels.BlackMarketOfferKind(r.Kind),
		PrototypeID:     int(r.ProtoID),
		PriceInCrystals: int(r.PriceInCrystals),
		EndsAt:          nullInt64ToInt64Ptr(r.EndsAt),
		IsLimited:       r.IsLimited,
		Priority:        int(r.Priority),
		Storage:         blackMarketStoragePrototypeFromRow(r),
	}
}

func blackMarketBuildPrototypeFromRow(r gen.ListActiveBlackMarketBuildingOffersRow) *readmodels.BuildItemPrototype {
	proto := buildPrototypeFromParts(r.ProtoID, r.Name, r.Category, r.Faction, r.UnlockTechnologyID, r.ShortDescription, r.FullDescription, r.Price, r.ProductionTime, r.Space, r.ImageUrl, r.ControlData, r.ResourcesData, r.DefenseData, r.MilitaryData, r.IntelligenceData)
	return &proto
}

func blackMarketArmyPrototypeFromRow(r gen.ListActiveBlackMarketArmyOffersRow) *readmodels.ArmyItemPrototype {
	proto := armyPrototypeFromParts(r.ProtoID, r.Name, r.Category, r.Faction, r.UnlockTechnologyID, r.ShortDescription, r.FullDescription, r.Price, r.ProductionTime, r.Space, r.ImageUrl, r.Attack, r.Defence, r.Capacity, r.Stealth, r.Speed)
	return &proto
}

func blackMarketStoragePrototypeFromRow(r gen.ListActiveBlackMarketStorageOffersRow) *readmodels.StorageItemPrototype {
	proto := readmodels.StorageItemPrototype{
		ID:               int(r.ProtoID),
		Name:             r.Name,
		Category:         readmodels.StorageCategory(r.Category),
		EstimatedWorth:   int(r.EstimatedWorth),
		ShortDescription: nullString(r.ShortDescription),
		FullDescription:  nullString(r.FullDescription),
		ImageURL:         nullString(r.ImageUrl),
		BuffData:         buffStorageDataFromJSON(r.BuffData),
		IntelData:        intelStorageDataFromJSON(r.IntelData),
		DamagedData:      damagedStorageDataFromJSON(r.DamagedData),
		ArtifactData:     artifactStorageDataFromJSON(r.ArtifactData),
		ConsumableData:   consumableStorageDataFromJSON(r.ConsumableData),
	}
	return &proto
}
