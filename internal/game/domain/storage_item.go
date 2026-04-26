package domain

import "math"

// StorageCategory represents the category of a storage item.
type StorageCategory string

type StorageStatus string

const (
	StorageStatusPresent  StorageStatus = "PRESENT"
	StorageStatusDeployed StorageStatus = "DEPLOYED"
)

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
	Name             TranslationKey
	Category         StorageCategory
	EstimatedWorth   int // Rough worth in credits
	ShortDescription TranslationKey
	FullDescription  TranslationKey
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
)

// BuffStorageData defines the properties for temporary stat enhancements.
type BuffStorageData struct {
	Type            BuffType
	Value           float32 // Multiplier or flat value depending on Type
	DurationSeconds int64   // Total active life of the buff once triggered
}

// HiddenLocationType identifies what kind of entity an intel item reveals on the map.
type HiddenLocationType string

const (
	HiddenLocationTypeResourceful HiddenLocationType = "RESOURCEFUL"
	HiddenLocationTypeDangerous   HiddenLocationType = "DANGEROUS" // e.g., NPC hostiles
	HiddenLocationTypeUserBase    HiddenLocationType = "USERBASE"  // Reveals a player base
)

// IntelStorageData defines properties for items that uncover hidden map nodes.
type IntelStorageData struct {
	Type              HiddenLocationType
	DecryptionSeconds int64 // Time required to "crack" the intel and reveal the location
}

// DamagedStorageData defines the requirements to restore a non-functional item (usually an army unit).
type DamagedStorageData struct {
	RestorePrice       PriceModel // Resources required for the repair
	RestorationSeconds int64      // Time required to restore the unit
	OriginalUnitID     int        // The ID of the ArmyUnitPrototype this will transform into
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
)

// ArtifactStorageData defines the properties for permanent passive items.
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
	ConsumableContentsIntel      ConsumableBoxContents = "INTEL"
	ConsumableContentsDamaged    ConsumableBoxContents = "DAMAGED"
	ConsumableContentsArtifact   ConsumableBoxContents = "ARTIFACT"
)

// ConsumableStorageData: properties for consumable items used in control buildings or special effects
type ConsumableStorageData struct {
	Type        ConsumableType
	BoxContents []ConsumableBoxContents // List of possible reward types
	BoxSize     int                     // Number of rolls or items inside
}

// StorageItemPresent represents a present storage item.
type StorageItemPresent struct {
	BaseOwnedItem
	Prototype StorageItemPrototype
	ExpiresAt *int64 // Unix timestamp for when the item (not necessarily the buff) disappears
	IsActive  bool   // Whether the item is currently active (e.g., toggled artifact or active buff)
}

// StorageItemDeployed represents a storage item reserved/deployed into an operation.
type StorageItemDeployed struct {
	BaseOwnedItem
	Prototype     StorageItemPrototype
	OperationKind OperationKind
	OperationID   int
	ExpiresAt     *int64
	IsActive      bool
}

// StorageItemSnap captures a snapshot of an active storage item (buff or artifact)
// to decouple resolution from potential prototype changes.
type StorageItemSnap struct {
	PrototypeID int
	Category    StorageCategory
	Buff        *BuffStorageData
	Artifact    *ArtifactStorageData
}

// StorageItemFromPresent converts a present storage item to a snapshot.
func StorageItemFromPresent(p StorageItemPresent) StorageItemSnap {
	return StorageItemSnap{
		PrototypeID: p.Prototype.ID,
		Category:    p.Prototype.Category,
		Buff:        p.Prototype.BuffData,
		Artifact:    p.Prototype.ArtifactData,
	}
}

// StorageItemsFromPresent maps present storage items to snapshots.
func StorageItemsFromPresent(items []StorageItemPresent) []StorageItemSnap {
	if len(items) == 0 {
		return nil
	}
	out := make([]StorageItemSnap, 0, len(items))
	for _, it := range items {
		if it.IsActive {
			out = append(out, StorageItemFromPresent(it))
		}
	}
	return out
}

// ProductionModifiers aggregates production-related multipliers.
type ProductionModifiers struct {
	CreditsProdMul  float64
	IronProdMul     float64
	TitaniumProdMul float64
}

// MilitaryModifiers aggregates combat-related multipliers.
type MilitaryModifiers struct {
	AttackMul   float64
	DefenceMul  float64
	StealthMul  float64
	CapacityMul float64
	SpeedMul    float64
}

// BaseModifiers aggregates all currently-active buffs/artifacts as multipliers.
type BaseModifiers struct {
	ProductionModifiers
	MilitaryModifiers
}

func IdentityBaseModifiers() BaseModifiers {
	return BaseModifiers{
		ProductionModifiers: ProductionModifiers{
			CreditsProdMul:  1,
			IronProdMul:     1,
			TitaniumProdMul: 1,
		},
		MilitaryModifiers: MilitaryModifiers{
			AttackMul:   1,
			DefenceMul:  1,
			StealthMul:  1,
			CapacityMul: 1,
			SpeedMul:    1,
		},
	}
}

func ModifiersFromSnaps(snaps []StorageItemSnap) BaseModifiers {
	m := IdentityBaseModifiers()
	for _, s := range snaps {
		if s.Buff != nil {
			m.ApplyBuff(s.Buff.Type, float64(s.Buff.Value))
		}
		if s.Artifact != nil {
			m.ApplyArtifact(s.Artifact.Type, float64(s.Artifact.Value))
		}
	}
	return m
}

func MilitaryModifiersFromSnaps(snaps []StorageItemSnap) MilitaryModifiers {
	return ModifiersFromSnaps(snaps).MilitaryModifiers
}

func mulInt(v int, mul float64) int {
	if v <= 0 {
		return 0
	}
	if mul <= 0 {
		return 0
	}
	return int(math.Round(float64(v) * mul))
}

func (m *BaseModifiers) ApplyBuff(t BuffType, v float64) {
	if v <= 0 {
		return
	}
	switch t {
	case BuffTypeCreditsProduction:
		m.CreditsProdMul *= v
	case BuffTypeIronProduction:
		m.IronProdMul *= v
	case BuffTypeTitaniumProduction:
		m.TitaniumProdMul *= v
	case BuffTypeAttackIncrease:
		m.AttackMul *= v
	case BuffTypeDefenceIncrease:
		m.DefenceMul *= v
	case BuffTypeStealthIncrease:
		m.StealthMul *= v
	case BuffTypeCapacityIncrease:
		m.CapacityMul *= v
	case BuffTypeSpeedIncrease:
		m.SpeedMul *= v
	}
}

func (m *BaseModifiers) ApplyArtifact(t ArtifactEffectType, v float64) {
	if v <= 0 {
		return
	}
	switch t {
	case ArtifactEffectTypeCreditsProduction:
		m.CreditsProdMul *= v
	case ArtifactEffectTypeIronProduction:
		m.IronProdMul *= v
	case ArtifactEffectTypeTitaniumProduction:
		m.TitaniumProdMul *= v
	case ArtifactEffectTypeAttackIncrease:
		m.AttackMul *= v
	case ArtifactEffectTypeDefenceIncrease:
		m.DefenceMul *= v
	case ArtifactEffectTypeStealthIncrease:
		m.StealthMul *= v
	case ArtifactEffectTypeCapacityIncrease:
		m.CapacityMul *= v
	case ArtifactEffectTypeSpeedIncrease:
		m.SpeedMul *= v
	}
}
