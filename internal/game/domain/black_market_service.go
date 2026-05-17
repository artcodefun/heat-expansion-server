package domain

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

	resourceAmount := rate.AmountPerCrystal * crystals
	if !base.Stats.CanReceiveResourceAmount(resourceType, resourceAmount) {
		return NewError("error.domain.black_market.resource_capacity_reached", H{"resource": resourceType})
	}

	if err := user.SpendCrystals(crystals); err != nil {
		return err
	}
	base.CreditLoot(rate.ResourcesForCrystals(crystals))

	return nil
}
