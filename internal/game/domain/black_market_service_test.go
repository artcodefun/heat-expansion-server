package domain

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestBlackMarketService_PurchaseResources_SpendsCrystalsAndCreditsBase(t *testing.T) {
	service := NewBlackMarketService()
	user := &User{ID: uuid.New(), Crystals: 10}
	base := newBaseWithDefaults(1)
	base.Stats.Credits = 0

	if err := service.PurchaseResources(user, base, ResourceTypeCredits, 2); err != nil {
		t.Fatalf("PurchaseResources error: %v", err)
	}

	if user.Crystals != 8 {
		t.Fatalf("expected 8 crystals remaining, got %d", user.Crystals)
	}
	if base.Stats.Credits != 200 {
		t.Fatalf("expected credits to increase by 200, got %v", base.Stats.Credits)
	}
}

func TestBlackMarketService_PurchaseResources_AntimatterCostsTenCrystals(t *testing.T) {
	service := NewBlackMarketService()
	user := &User{ID: uuid.New(), Crystals: 20}
	base := newBaseWithDefaults(3)
	base.Stats.Antimatter = 0

	if err := service.PurchaseResources(user, base, ResourceTypeAntimatter, 10); err != nil {
		t.Fatalf("PurchaseResources error: %v", err)
	}

	if user.Crystals != 10 {
		t.Fatalf("expected 10 crystals remaining, got %d", user.Crystals)
	}
	if base.Stats.Antimatter != 1 {
		t.Fatalf("expected antimatter to increase by 1, got %v", base.Stats.Antimatter)
	}
}

func TestBlackMarketService_PurchaseResources_RejectsWhenBaseCannotReceive(t *testing.T) {
	service := NewBlackMarketService()
	user := &User{ID: uuid.New(), Crystals: 10}
	base := newBaseWithDefaults(2)
	base.Stats.Iron = float64(base.Stats.IronCapacity)
	startingCrystals := user.Crystals
	startingIron := base.Stats.Iron

	err := service.PurchaseResources(user, base, ResourceTypeIron, 1)
	if err == nil {
		t.Fatal("expected capacity error")
	}
	var domainErr Error
	if !errors.As(err, &domainErr) {
		t.Fatalf("expected domain error, got %T", err)
	}
	if domainErr.Key != "error.domain.base.resource_capacity_reached" {
		t.Fatalf("expected base capacity key, got %s", domainErr.Key)
	}
	if user.Crystals != startingCrystals {
		t.Fatalf("expected crystals unchanged, got %d", user.Crystals)
	}
	if base.Stats.Iron != startingIron {
		t.Fatalf("expected iron unchanged, got %v", base.Stats.Iron)
	}
}

func TestBlackMarketService_PurchaseResources_RejectsInvalidCrystalMultiple(t *testing.T) {
	service := NewBlackMarketService()
	user := &User{ID: uuid.New(), Crystals: 10}
	base := newBaseWithDefaults(4)

	err := service.PurchaseResources(user, base, ResourceTypeAntimatter, 5)
	if err == nil {
		t.Fatal("expected invalid crystal amount error")
	}
	var domainErr Error
	if !errors.As(err, &domainErr) {
		t.Fatalf("expected domain error, got %T", err)
	}
	if domainErr.Key != "error.domain.black_market.invalid_crystal_amount" {
		t.Fatalf("expected invalid crystal amount key, got %s", domainErr.Key)
	}
	if user.Crystals != 10 {
		t.Fatalf("expected crystals unchanged, got %d", user.Crystals)
	}
}

func TestBlackMarketService_PurchaseBuildingOffer_AddsBuildingAndSpendsCrystals(t *testing.T) {
	service := NewBlackMarketService()
	user := &User{ID: uuid.New(), Crystals: 12}
	base := newBaseWithDefaults(5)
	offer := BlackMarketOffer{ID: 1, Kind: BlackMarketOfferKindBuilding, PrototypeID: 100, PriceInCrystals: 4}
	proto := &BuildItemPrototype{
		ID:              100,
		Category:        BuildCategoryResources,
		CreationSources: []CreationSource{CreationSourceBlackMarket},
		Faction:         FactionExoCoalition,
		Price:           PriceModel{Credits: 100},
	}

	if err := service.PurchaseBuildingOffer(user, base, offer, proto); err != nil {
		t.Fatalf("PurchaseBuildingOffer error: %v", err)
	}

	if user.Crystals != 8 {
		t.Fatalf("expected 8 crystals remaining, got %d", user.Crystals)
	}
	if len(base.BuildingsPresent) != 1 {
		t.Fatalf("expected 1 present building, got %d", len(base.BuildingsPresent))
	}
}

func TestBlackMarketService_PurchaseBuildingOffer_RejectsWhenUserCannotAfford(t *testing.T) {
	service := NewBlackMarketService()
	user := &User{ID: uuid.New(), Crystals: 3}
	base := newBaseWithDefaults(10)
	offer := BlackMarketOffer{ID: 6, Kind: BlackMarketOfferKindBuilding, PrototypeID: 101, PriceInCrystals: 4}
	proto := &BuildItemPrototype{
		ID:              101,
		Category:        BuildCategoryResources,
		CreationSources: []CreationSource{CreationSourceBlackMarket},
		Faction:         FactionExoCoalition,
		Price:           PriceModel{Credits: 100},
	}

	err := service.PurchaseBuildingOffer(user, base, offer, proto)
	if err == nil {
		t.Fatal("expected not enough crystals error")
	}
	var domainErr Error
	if !errors.As(err, &domainErr) {
		t.Fatalf("expected domain error, got %T", err)
	}
	if domainErr.Key != "error.domain.user.not_enough_crystals" {
		t.Fatalf("expected not enough crystals key, got %s", domainErr.Key)
	}
	if user.Crystals != 3 {
		t.Fatalf("expected crystals unchanged, got %d", user.Crystals)
	}
}

func TestBlackMarketService_PurchaseArmyOffer_ScalesPriceByQuantity(t *testing.T) {
	service := NewBlackMarketService()
	user := &User{ID: uuid.New(), Crystals: 20}
	base := newBaseWithDefaults(6)
	offer := BlackMarketOffer{ID: 2, Kind: BlackMarketOfferKindArmy, PrototypeID: 200, PriceInCrystals: 3}
	proto := &ArmyItemPrototype{
		ID:              200,
		Category:        ArmyCategoryInfantry,
		CreationSources: []CreationSource{CreationSourceBlackMarket},
		Faction:         FactionExoCoalition,
		Price:           PriceModel{Credits: 50},
		Space:           1,
	}

	if err := service.PurchaseArmyOffer(user, base, offer, proto, 2); err != nil {
		t.Fatalf("PurchaseArmyOffer error: %v", err)
	}

	if user.Crystals != 14 {
		t.Fatalf("expected 14 crystals remaining, got %d", user.Crystals)
	}
	if len(base.ArmiesPresent) != 1 || base.ArmiesPresent[0].Count != 2 {
		t.Fatalf("expected 2 present army units, got %+v", base.ArmiesPresent)
	}
}

func TestBlackMarketService_PurchaseStorageOffer_AddsStorageItem(t *testing.T) {
	service := NewBlackMarketService()
	user := &User{ID: uuid.New(), Crystals: 9}
	base := newBaseWithDefaults(7)
	offer := BlackMarketOffer{ID: 3, Kind: BlackMarketOfferKindStorage, PrototypeID: 300, PriceInCrystals: 2}
	proto := &StorageItemPrototype{
		ID:              300,
		Category:        StorageCategoryBuff,
		CreationSources: []CreationSource{CreationSourceBlackMarket},
		BuffData:        &BuffStorageData{Type: BuffTypeAttackIncrease, Value: 1.1, DurationSeconds: 60},
	}

	if err := service.PurchaseStorageOffer(user, base, offer, proto); err != nil {
		t.Fatalf("PurchaseStorageOffer error: %v", err)
	}

	if user.Crystals != 7 {
		t.Fatalf("expected 7 crystals remaining, got %d", user.Crystals)
	}
	if len(base.StorageItemsPresent) != 1 {
		t.Fatalf("expected 1 present storage item, got %d", len(base.StorageItemsPresent))
	}
}

func TestBlackMarketService_PurchaseArmyOffer_RejectsInactiveOffer(t *testing.T) {
	service := NewBlackMarketService()
	user := &User{ID: uuid.New(), Crystals: 20}
	base := newBaseWithDefaults(8)
	expiredAt := NowUnix() - 1
	offer := BlackMarketOffer{ID: 4, Kind: BlackMarketOfferKindArmy, PrototypeID: 201, PriceInCrystals: 3, IsLimited: true, EndsAt: &expiredAt}
	proto := &ArmyItemPrototype{
		ID:              201,
		Category:        ArmyCategoryInfantry,
		CreationSources: []CreationSource{CreationSourceBlackMarket},
		Faction:         FactionExoCoalition,
		Space:           1,
	}

	err := service.PurchaseArmyOffer(user, base, offer, proto, 1)
	if err == nil {
		t.Fatal("expected inactive offer error")
	}
	var domainErr Error
	if !errors.As(err, &domainErr) {
		t.Fatalf("expected domain error, got %T", err)
	}
	if domainErr.Key != "error.domain.black_market.offer_not_active" {
		t.Fatalf("expected inactive offer key, got %s", domainErr.Key)
	}
}

func TestBlackMarketService_PurchaseStorageOffer_RejectsPrototypeWithoutBlackMarketSource(t *testing.T) {
	service := NewBlackMarketService()
	user := &User{ID: uuid.New(), Crystals: 9}
	base := newBaseWithDefaults(9)
	offer := BlackMarketOffer{ID: 5, Kind: BlackMarketOfferKindStorage, PrototypeID: 301, PriceInCrystals: 2}
	proto := &StorageItemPrototype{
		ID:              301,
		Category:        StorageCategoryBuff,
		CreationSources: []CreationSource{CreationSourcePlayerBase},
		BuffData:        &BuffStorageData{Type: BuffTypeAttackIncrease, Value: 1.1, DurationSeconds: 60},
	}

	err := service.PurchaseStorageOffer(user, base, offer, proto)
	if err == nil {
		t.Fatal("expected prototype source error")
	}
	var domainErr Error
	if !errors.As(err, &domainErr) {
		t.Fatalf("expected domain error, got %T", err)
	}
	if domainErr.Key != "error.domain.black_market.prototype_not_available" {
		t.Fatalf("expected prototype not available key, got %s", domainErr.Key)
	}
}
