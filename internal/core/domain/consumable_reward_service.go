package domain

import (
	"fmt"
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
		return nil, fmt.Errorf("box not found")
	}

	if box.Prototype.Category != StorageCategoryConsumable || box.Prototype.ConsumableData == nil {
		return nil, fmt.Errorf("item is not a consumable box")
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
		val := 100 + rand.Intn(400) // Placeholder logic
		ub.CreditLoot(PriceModel{Credits: val})
		return BoxRewardOutcome{Type: contentType, Value: val}

	case ConsumableContentsIron:
		val := 50 + rand.Intn(200)
		ub.CreditLoot(PriceModel{Iron: val})
		return BoxRewardOutcome{Type: contentType, Value: val}

	case ConsumableContentsTitanium:
		val := 25 + rand.Intn(100)
		ub.CreditLoot(PriceModel{Titanium: val})
		return BoxRewardOutcome{Type: contentType, Value: val}

	case ConsumableContentsAntimatter:
		val := 10 + rand.Intn(40)
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

	case ConsumableContentsMap:
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
