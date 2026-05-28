package domain

import "slices"

// CreationSource identifies where a prototype is allowed to originate from.
type CreationSource string

const (
	CreationSourcePlayerBase    CreationSource = "PLAYER_BASE"
	CreationSourceBlackMarket   CreationSource = "BLACK_MARKET"
	CreationSourceNPCLocation   CreationSource = "NPC_LOCATION"
	CreationSourceConsumableBox CreationSource = "CONSUMABLE_BOX"
)

func FilterArmyItemPrototypesByCreationSource(prototypes []*ArmyItemPrototype, source CreationSource) []*ArmyItemPrototype {
	filtered := make([]*ArmyItemPrototype, 0, len(prototypes))
	for _, prototype := range prototypes {
		if slices.Contains(prototype.CreationSources, source) {
			filtered = append(filtered, prototype)
		}
	}
	return filtered
}

func FilterBuildItemPrototypesByCreationSource(prototypes []*BuildItemPrototype, source CreationSource) []*BuildItemPrototype {
	filtered := make([]*BuildItemPrototype, 0, len(prototypes))
	for _, prototype := range prototypes {
		if slices.Contains(prototype.CreationSources, source) {
			filtered = append(filtered, prototype)
		}
	}
	return filtered
}

func FilterStorageItemPrototypesByCreationSource(prototypes []*StorageItemPrototype, source CreationSource) []*StorageItemPrototype {
	filtered := make([]*StorageItemPrototype, 0, len(prototypes))
	for _, prototype := range prototypes {
		if slices.Contains(prototype.CreationSources, source) {
			filtered = append(filtered, prototype)
		}
	}
	return filtered
}
