package domain

import "slices"

// BlackMarketService handles purchasing resources from the Black Market.
type BlackMarketService struct{}

func NewBlackMarketService() *BlackMarketService {
	return &BlackMarketService{}
}

// PurchaseResources spends crystals and grants the corresponding resources to a base.
func (s *BlackMarketService) PurchaseResources(user *User, base *UserBaseModel, resourceType ResourceType, crystals int) error {
	rate, err := BlackMarketResourceRateFor(resourceType)
	if err != nil {
		return err
	}

	resourceAmount, err := rate.ResourceAmountForCrystals(crystals)
	if err != nil {
		return err
	}

	if err := user.SpendCrystals(crystals); err != nil {
		return err
	}
	if err := base.ReceiveResource(resourceType, resourceAmount); err != nil {
		return err
	}

	return nil
}

func (s *BlackMarketService) PurchaseBuildingOffer(user *User, base *UserBaseModel, offer BlackMarketOffer, proto *BuildItemPrototype) error {
	if err := validateBlackMarketOffer(offer, BlackMarketOfferKindBuilding, 1); err != nil {
		return err
	}
	if err := validateBlackMarketBuildPrototype(proto); err != nil {
		return err
	}
	if err := user.SpendCrystals(offer.PriceInCrystals); err != nil {
		return err
	}
	if err := base.ReceiveBuilding(*proto); err != nil {
		return err
	}
	return nil
}

func (s *BlackMarketService) PurchaseArmyOffer(user *User, base *UserBaseModel, offer BlackMarketOffer, proto *ArmyItemPrototype, quantity int) error {
	if err := validateBlackMarketOffer(offer, BlackMarketOfferKindArmy, quantity); err != nil {
		return err
	}
	if err := validateBlackMarketArmyPrototype(proto); err != nil {
		return err
	}
	if err := user.SpendCrystals(offer.PriceInCrystals * quantity); err != nil {
		return err
	}
	if err := base.ReceiveArmyItems(*proto, quantity); err != nil {
		return err
	}
	return nil
}

func (s *BlackMarketService) PurchaseStorageOffer(user *User, base *UserBaseModel, offer BlackMarketOffer, proto *StorageItemPrototype) error {
	if err := validateBlackMarketOffer(offer, BlackMarketOfferKindStorage, 1); err != nil {
		return err
	}
	if err := validateBlackMarketStoragePrototype(proto); err != nil {
		return err
	}
	if err := user.SpendCrystals(offer.PriceInCrystals); err != nil {
		return err
	}
	if err := base.ReceiveStorageItem(*proto); err != nil {
		return err
	}
	return nil
}

func validateBlackMarketOffer(offer BlackMarketOffer, expectedKind BlackMarketOfferKind, quantity int) error {
	if quantity < 1 {
		return NewError("error.domain.black_market.invalid_offer_quantity", H{"quantity": quantity})
	}
	if offer.Kind != expectedKind {
		return NewError("error.domain.black_market.offer_kind_mismatch", nil)
	}
	if !offer.IsActive(NowUnix()) {
		return NewError("error.domain.black_market.offer_not_active", H{"offer_id": offer.ID})
	}
	if expectedKind != BlackMarketOfferKindArmy && quantity != 1 {
		return NewError("error.domain.black_market.invalid_offer_quantity", H{"quantity": quantity})
	}
	return nil
}

func validateBlackMarketBuildPrototype(proto *BuildItemPrototype) error {
	if proto == nil || !slices.Contains(proto.CreationSources, CreationSourceBlackMarket) {
		prototypeID := 0
		if proto != nil {
			prototypeID = proto.ID
		}
		return NewError("error.domain.black_market.prototype_not_available", H{"prototype": prototypeID})
	}
	return nil
}

func validateBlackMarketArmyPrototype(proto *ArmyItemPrototype) error {
	if proto == nil || !slices.Contains(proto.CreationSources, CreationSourceBlackMarket) {
		prototypeID := 0
		if proto != nil {
			prototypeID = proto.ID
		}
		return NewError("error.domain.black_market.prototype_not_available", H{"prototype": prototypeID})
	}
	return nil
}

func validateBlackMarketStoragePrototype(proto *StorageItemPrototype) error {
	if proto == nil || !slices.Contains(proto.CreationSources, CreationSourceBlackMarket) {
		prototypeID := 0
		if proto != nil {
			prototypeID = proto.ID
		}
		return NewError("error.domain.black_market.prototype_not_available", H{"prototype": prototypeID})
	}
	return nil
}
