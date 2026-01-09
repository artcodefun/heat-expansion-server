package readmodels

// StorageCategory represents the category of a storage item.
type StorageCategory string

const (
	StorageCategoryBuff       StorageCategory = "BUFF"
	StorageCategoryIntel      StorageCategory = "INTEL"
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
	IntelData      *IntelStorageData
	DamagedData    *DamagedStorageData
	ArtifactData   *ArtifactStorageData
	ConsumableData *ConsumableStorageData
}

// BuffType identifies the specific stat or resource a buff affects.
type BuffType string

const (
	BuffTypeCreditsProduction  BuffType = "CREDITS_PRODUCTION"
	BuffTypeIronProduction     BuffType = "IRON_PRODUCTION"
	BuffTypeTitaniumProduction BuffType = "TITANIUM_PRODUCTION"

	BuffTypeAttackIncrease   BuffType = "ATTACK_INCREASE"
	BuffTypeDefenceIncrease  BuffType = "DEFENCE_INCREASE"
	BuffTypeStealthIncrease  BuffType = "STEALTH_INCREASE"
	BuffTypeCapacityIncrease BuffType = "CAPACITY_INCREASE"
	BuffTypeSpeedIncrease    BuffType = "SPEED_INCREASE"

	BuffTypePricesDecrease BuffType = "PRICES_DECREASE"
)

// BuffStorageData: bonuses for buffs
type BuffStorageData struct {
	Type            BuffType
	Value           float32
	DurationSeconds int64
}

// HiddenLocationType identifies what kind of entity an intel item reveals on the map.
type HiddenLocationType string

const (
	HiddenLocationTypeResourcefulCredits  HiddenLocationType = "RESOURCEFUL_CREDITS"
	HiddenLocationTypeResourcefulIron     HiddenLocationType = "RESOURCEFUL_IRON"
	HiddenLocationTypeResourcefulTitanium HiddenLocationType = "RESOURCEFUL_TITANIUM"
	HiddenLocationTypeDangerous           HiddenLocationType = "DANGEROUS" // e.g., NPC hostiles
	HiddenLocationTypeUserBase            HiddenLocationType = "USERBASE"  // Reveals a player base
)

// IntelStorageData: properties for intel items
type IntelStorageData struct {
	Type              HiddenLocationType
	DecryptionSeconds int64
}

// DamagedStorageData: properties for damaged items (e.g., damaged army units)
type DamagedStorageData struct {
	RestorePrice   PriceModel // Cost to restore the item
	OriginalUnitID int        // Reference to the original army unit prototype
}

// ArtifactEffectType identifies the permanent passive bonus provided by an artifact.
type ArtifactEffectType string

const (
	ArtifactEffectTypeCreditsProduction  ArtifactEffectType = "CREDITS_PRODUCTION"
	ArtifactEffectTypeIronProduction     ArtifactEffectType = "IRON_PRODUCTION"
	ArtifactEffectTypeTitaniumProduction ArtifactEffectType = "TITANIUM_PRODUCTION"

	ArtifactEffectTypeAttackIncrease   ArtifactEffectType = "ATTACK_INCREASE"
	ArtifactEffectTypeDefenceIncrease  ArtifactEffectType = "DEFENCE_INCREASE"
	ArtifactEffectTypeStealthIncrease  ArtifactEffectType = "STEALTH_INCREASE"
	ArtifactEffectTypeCapacityIncrease ArtifactEffectType = "CAPACITY_INCREASE"
	ArtifactEffectTypeSpeedIncrease    ArtifactEffectType = "SPEED_INCREASE"

	ArtifactEffectTypePricesDecrease ArtifactEffectType = "PRICES_DECREASE"
)

// ArtifactStorageData: properties for artifacts
type ArtifactStorageData struct {
	Type  ArtifactEffectType
	Value float32
}

// ConsumableType identifies the type of consumable
type ConsumableType string

const (
	ConsumableTypeBox         ConsumableType = "BOX"          // A loot box containing multiple rewards
	ConsumableTypeWarpCapsule ConsumableType = "WARP_CAPSULE" // An item with a specific gameplay function
)

// ConsumableBoxContents defines the pool of possible rewards within a loot box.
type ConsumableBoxContents string

const (
	ConsumableContentsCredits    ConsumableBoxContents = "CREDITS"
	ConsumableContentsIron       ConsumableBoxContents = "IRON"
	ConsumableContentsTitanium   ConsumableBoxContents = "TITANIUM"
	ConsumableContentsAntimatter ConsumableBoxContents = "ANTIMATTER"
	ConsumableContentsCrystals   ConsumableBoxContents = "CRYSTALS"
	ConsumableContentsBuff       ConsumableBoxContents = "BUFF"
	ConsumableContentsMap        ConsumableBoxContents = "MAP"
	ConsumableContentsDamaged    ConsumableBoxContents = "DAMAGED"
	ConsumableContentsArtifact   ConsumableBoxContents = "ARTIFACT"
)

// ConsumableStorageData: properties for consumable items used in control buildings or special effects
type ConsumableStorageData struct {
	Type        ConsumableType
	BoxContents []ConsumableBoxContents
	BoxSize     int
}

// StorageItemPresent represents a present storage item.
type StorageItemPresent struct {
	BaseOwnedItem
	Prototype StorageItemPrototype
	ExpiresAt *int64
	IsActive  bool
}
