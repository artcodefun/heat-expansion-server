package domain

import (
	"math/rand"

	"github.com/google/uuid"
)

// BoxRewardOutcome describes a single reward rolled from a consumable box.
type BoxRewardOutcome struct {
	Type  ConsumableBoxContents
	Value interface{} // int for resources/crystals, uuid.UUID for new items
}

// ConsumableRewardService handles the logic of rolling for and applying rewards from consumable boxes.
type ConsumableRewardService struct {
}

func NewConsumableRewardService() *ConsumableRewardService {
	return &ConsumableRewardService{}
}

// OpenBox consumes a box from the base's storage and applies its contents.
// It returns the list of outcomes for reporting or event emission.
func (s *ConsumableRewardService) OpenBox(
	ub *UserBaseModel,
	user *User,
	boxID uuid.UUID,
	buffProtos []StorageItemPrototype,
	intelProtos []StorageItemPrototype,
	damagedProtos []StorageItemPrototype,
	artifactProtos []StorageItemPrototype,
) ([]BoxRewardOutcome, error) {
	// 1. Find the box in storage
	var box *StorageItemPresent
	for _, it := range ub.StorageItemsPresent {
		if it.ID == boxID {
			box = &it
			break
		}
	}

	if box == nil {
		return nil, NewError("error.domain.storage.box_not_found", nil)
	}

	if box.Prototype.Category != StorageCategoryConsumable || box.Prototype.ConsumableData == nil {
		return nil, NewError("error.domain.storage.not_a_consumable_box", nil)
	}

	data := box.Prototype.ConsumableData
	outcomes := make([]BoxRewardOutcome, 0, data.BoxSize)

	// 2. Roll for rewards
	for i := 0; i < data.BoxSize; i++ {
		outcome := s.RollSingleReward(ub, user, data.BoxContents, buffProtos, intelProtos, damagedProtos, artifactProtos)
		if outcome != (BoxRewardOutcome{}) {
			outcomes = append(outcomes, outcome)
		}
	}

	// 3. Remove the box from storage
	if err := ub.DeletePresentStorageItemByID(boxID); err != nil {
		return nil, err
	}

	return outcomes, nil
}

// RollSingleReward picks a random content type from the pool and applies it.
func (s *ConsumableRewardService) RollSingleReward(
	ub *UserBaseModel,
	user *User,
	pool []ConsumableBoxContents,
	buffProtos []StorageItemPrototype,
	intelProtos []StorageItemPrototype,
	damagedProtos []StorageItemPrototype,
	artifactProtos []StorageItemPrototype,
) BoxRewardOutcome {
	if len(pool) == 0 {
		return BoxRewardOutcome{}
	}

	contentType := pool[rand.Intn(len(pool))]

	switch contentType {
	case ConsumableContentsCredits:
		val := int(500 + rand.Float64()*1500) // ~500-2000 credits worth
		ub.CreditLoot(PriceModel{Credits: val})
		return BoxRewardOutcome{Type: contentType, Value: val}

	case ConsumableContentsIron:
		// 1 Iron = 4 Credits. Reward: ~125-500 Iron
		val := int((500 + rand.Float64()*1500) / WorthIron)
		ub.CreditLoot(PriceModel{Iron: val})
		return BoxRewardOutcome{Type: contentType, Value: val}

	case ConsumableContentsTitanium:
		// 1 Titanium = 20 Credits. Reward: ~25-100 Titanium
		val := int((500 + rand.Float64()*1500) / WorthTitanium)
		ub.CreditLoot(PriceModel{Titanium: val})
		return BoxRewardOutcome{Type: contentType, Value: val}

	case ConsumableContentsAntimatter:
		// 1 Antimatter = 333 Credits. Reward: ~1.5-6 Antimatter
		val := int((500 + rand.Float64()*1500) / WorthAntimatter)
		if val < 1 {
			val = 1
		}
		ub.CreditLoot(PriceModel{Antimatter: val})
		return BoxRewardOutcome{Type: contentType, Value: val}

	case ConsumableContentsCrystals:
		val := 1 + rand.Intn(4)
		user.Crystals += val
		return BoxRewardOutcome{Type: contentType, Value: val}

	case ConsumableContentsBuff:
		if len(buffProtos) == 0 {
			return BoxRewardOutcome{}
		}
		proto := buffProtos[rand.Intn(len(buffProtos))]
		itemID := ub.AddStorageItem(proto, nil)
		return BoxRewardOutcome{Type: contentType, Value: itemID}

	case ConsumableContentsIntel:
		if len(intelProtos) == 0 {
			return BoxRewardOutcome{}
		}
		proto := intelProtos[rand.Intn(len(intelProtos))]
		itemID := ub.AddStorageItem(proto, nil)
		return BoxRewardOutcome{Type: contentType, Value: itemID}

	case ConsumableContentsDamaged:
		if len(damagedProtos) == 0 {
			return BoxRewardOutcome{}
		}
		proto := damagedProtos[rand.Intn(len(damagedProtos))]
		itemID := ub.AddStorageItem(proto, nil)
		return BoxRewardOutcome{Type: contentType, Value: itemID}

	case ConsumableContentsArtifact:
		if len(artifactProtos) == 0 {
			return BoxRewardOutcome{}
		}
		proto := artifactProtos[rand.Intn(len(artifactProtos))]
		itemID := ub.AddStorageItem(proto, nil)
		return BoxRewardOutcome{Type: contentType, Value: itemID}
	}

	return BoxRewardOutcome{}
}
