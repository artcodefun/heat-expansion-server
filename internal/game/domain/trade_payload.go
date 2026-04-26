package domain

import "github.com/google/uuid"

// TradeArmyItemSnap captures an army prototype and quantity in a trade payload snapshot.
type TradeArmyItemSnap struct {
	PrototypeID int
	Count       int
	Capacity    int
}

// TradeStorageItemSnap captures a concrete storage item and lightweight prototype metadata in a trade payload snapshot.
type TradeStorageItemSnap struct {
	ItemID      uuid.UUID
	PrototypeID int
	Category    StorageCategory
}

// TradePayload represents one side of a trade deal (offered or requested).
// It is immutable after operation creation.
type TradePayload struct {
	Resources PriceModel
	Storage   []TradeStorageItemSnap
	Army      []TradeArmyItemSnap
}

func NewTradePayload(resources PriceModel, storage []TradeStorageItemSnap, army []TradeArmyItemSnap) (TradePayload, error) {
	if resources.Credits < 0 || resources.Iron < 0 || resources.Titanium < 0 || resources.Antimatter < 0 {
		return TradePayload{}, NewError("error.domain.trade.invalid_resources", nil)
	}

	seenStorage := make(map[uuid.UUID]struct{}, len(storage))
	normStorage := make([]TradeStorageItemSnap, 0, len(storage))
	for _, s := range storage {
		if _, ok := seenStorage[s.ItemID]; ok {
			return TradePayload{}, NewError("error.domain.trade.duplicate_storage_item", nil)
		}
		seenStorage[s.ItemID] = struct{}{}
		normStorage = append(normStorage, s)
	}

	normArmy := make([]TradeArmyItemSnap, 0, len(army))
	for _, a := range army {
		if a.Count <= 0 || a.Capacity < 0 {
			return TradePayload{}, NewError("error.domain.trade.invalid_army_item", nil)
		}
		normArmy = append(normArmy, a)
	}

	if resources.Credits == 0 && resources.Iron == 0 && resources.Titanium == 0 && resources.Antimatter == 0 && len(normStorage) == 0 && len(normArmy) == 0 {
		return TradePayload{}, NewError("error.domain.trade.empty_payload", nil)
	}

	return TradePayload{
		Resources: resources,
		Storage:   normStorage,
		Army:      normArmy,
	}, nil
}

// RequiredResourceCapacity returns required carrying capacity points for
// transporting payload resources.
func (p TradePayload) RequiredResourceCapacity() float64 {
	return p.Resources.CreditsWorth() / WorthCapacityMultiplier
}

// ProvidedArmyCapacity returns carrying capacity points provided by army units
// included in this payload (before any active modifiers are applied).
func (p TradePayload) ProvidedArmyCapacity() float64 {
	if len(p.Army) == 0 {
		return 0
	}
	total := 0.0
	for _, a := range p.Army {
		if a.Count <= 0 || a.Capacity <= 0 {
			continue
		}
		total += float64(a.Capacity * a.Count)
	}
	return total
}
