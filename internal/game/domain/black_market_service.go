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

	resourceAmount, err := rate.ResourceAmountForCrystals(crystals)
	if err != nil {
		return err
	}
	if err := base.ReceiveResource(resourceType, resourceAmount); err != nil {
		return err
	}

	if err := user.SpendCrystals(crystals); err != nil {
		return err
	}

	return nil
}
