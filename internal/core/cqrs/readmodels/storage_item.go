package readmodels

// StorageCategory represents the category of a storage item.
type StorageCategory string

const (
	StorageCategoryBuff       StorageCategory = "BUFF"
	StorageCategoryMap        StorageCategory = "MAP"
	StorageCategoryDamaged    StorageCategory = "DAMAGED"
	StorageCategoryArtifact   StorageCategory = "ARTIFACT"
	StorageCategoryConsumable StorageCategory = "CONSUMABLE"
)

// StorageItemPrototype is the base struct for storage item prototypes.
type StorageItemPrototype struct {
	ID               int
	Name             string
	Category         StorageCategory
	ShortDescription string
	FullDescription  string
	ImageURL         string

	// Category-specific fields
	BuffData       *BuffStorageData
	MapData        *MapStorageData
	DamagedData    *DamagedStorageData
	ArtifactData   *ArtifactStorageData
	ConsumableData *ConsumableStorageData
}

// BuffStorageData: bonuses for buffs
type BuffStorageData struct {
	SpaceCapacityBonus int
	AttackBonus        int
	DefenceBonus       int
	DurationSeconds    int64
	ActivatedAt        *int64
}

// MapStorageData: properties for map items
type MapStorageData struct {
	RevealedArea string // e.g., region or sector identifier
	ScanRange    int    // How much area is revealed
	// Add more as needed
}

// DamagedStorageData: properties for damaged items (e.g., damaged army units)
type DamagedStorageData struct {
	RestorePrice   PriceModel // Cost to restore the item
	OriginalUnitID int        // Reference to the original army unit prototype
	DamageLevel    int        // Severity of damage
	// Add more as needed
}

// ArtifactStorageData: properties for artifacts
type ArtifactStorageData struct {
	PassiveEffect string // Description of constant effect
	Rarity        string // e.g., "Common", "Rare", "Legendary"
	Lore          string // Optional: flavor text or story
	// Add more as needed
}

// ConsumableStorageData: properties for consumable items used in control buildings or special effects
type ConsumableStorageData struct {
	EffectType   string      // e.g., "RESTORE_UNIT", "BOOST_PRODUCTION", etc.
	Uses         int         // Number of uses (if applicable)
	RestorePrice *PriceModel // Optional: cost to use/restore if relevant
	// Add more as needed
}

// StorageItemPresent represents a present storage item.
type StorageItemPresent struct {
	BaseOwnedItem
	Prototype StorageItemPrototype
	Refund    PriceModel
}
