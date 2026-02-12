package dtos

import (
	"strings"

	"github.com/artcodefun/heat-expansion-api/internal/game/application/cqrs/readmodels"
)

type StorageCategory string

// StorageCategory enum for DTOs
const (
	StorageCategoryBuff       StorageCategory = "BUFF"
	StorageCategoryIntel      StorageCategory = "INTEL"
	StorageCategoryDamaged    StorageCategory = "DAMAGED"
	StorageCategoryArtifact   StorageCategory = "ARTIFACT"
	StorageCategoryConsumable StorageCategory = "CONSUMABLE"
)

type BuffType string

const (
	BuffCreditsProduction  BuffType = "CREDITS_PRODUCTION"
	BuffIronProduction     BuffType = "IRON_PRODUCTION"
	BuffTitaniumProduction BuffType = "TITANIUM_PRODUCTION"
	BuffAttackIncrease     BuffType = "ATTACK_INCREASE"
	BuffDefenceIncrease    BuffType = "DEFENCE_INCREASE"
	BuffStealthIncrease    BuffType = "STEALTH_INCREASE"
	BuffCapacityIncrease   BuffType = "CAPACITY_INCREASE"
	BuffSpeedIncrease      BuffType = "SPEED_INCREASE"
)

type HiddenLocationType string

const (
	HiddenLocationResourceful HiddenLocationType = "RESOURCEFUL"
	HiddenLocationDangerous   HiddenLocationType = "DANGEROUS"
	HiddenLocationUserBase    HiddenLocationType = "USERBASE"
)

type ArtifactEffectType string

const (
	ArtifactCreditsProduction  ArtifactEffectType = "CREDITS_PRODUCTION"
	ArtifactIronProduction     ArtifactEffectType = "IRON_PRODUCTION"
	ArtifactTitaniumProduction ArtifactEffectType = "TITANIUM_PRODUCTION"
	ArtifactAttackIncrease     ArtifactEffectType = "ATTACK_INCREASE"
	ArtifactDefenceIncrease    ArtifactEffectType = "DEFENCE_INCREASE"
	ArtifactStealthIncrease    ArtifactEffectType = "STEALTH_INCREASE"
	ArtifactCapacityIncrease   ArtifactEffectType = "CAPACITY_INCREASE"
	ArtifactSpeedIncrease      ArtifactEffectType = "SPEED_INCREASE"
)

type ConsumableType string

const (
	ConsumableBox         ConsumableType = "BOX"
	ConsumableWarpCapsule ConsumableType = "WARP_CAPSULE"
)

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

type StorageItemPrototypeDTO struct {
	ID               int             `json:"id"`
	Name             string          `json:"name"`
	Category         StorageCategory `json:"category"`
	ShortDescription string          `json:"short_description"`
	FullDescription  string          `json:"full_description"`
	ImageURL         string          `json:"image_url"`

	// Category-specific fields
	BuffData       *BuffStorageDataDTO       `json:"buff_data,omitempty"`
	IntelData      *IntelStorageDataDTO      `json:"intel_data,omitempty"`
	DamagedData    *DamagedStorageDataDTO    `json:"damaged_data,omitempty"`
	ArtifactData   *ArtifactStorageDataDTO   `json:"artifact_data,omitempty"`
	ConsumableData *ConsumableStorageDataDTO `json:"consumable_data,omitempty"`
}

type BuffStorageDataDTO struct {
	Type            BuffType `json:"type"`
	Value           float32  `json:"value"`
	DurationSeconds int64    `json:"duration_seconds"`
}

type IntelStorageDataDTO struct {
	Type              HiddenLocationType `json:"type"`
	DecryptionSeconds int64              `json:"decryption_seconds"`
}

type DamagedStorageDataDTO struct {
	RestorePrice   PriceModelDTO `json:"restore_price"`
	OriginalUnitID int           `json:"original_unit_id"`
}

type ArtifactStorageDataDTO struct {
	Type  ArtifactEffectType `json:"type"`
	Value float32            `json:"value"`
}

type ConsumableStorageDataDTO struct {
	Type        ConsumableType          `json:"type"`
	BoxContents []ConsumableBoxContents `json:"box_contents"`
	BoxSize     int                     `json:"box_size"`
}

type StorageItemPresentDTO struct {
	BaseOwnedItemDTO
	Prototype StorageItemPrototypeDTO `json:"prototype"`
	ExpiresAt *int64                  `json:"expires_at,omitempty"`
	IsActive  bool                    `json:"is_active"`
}

func StorageCategoryFromDTO(c string) readmodels.StorageCategory {
	return readmodels.StorageCategory(strings.ToUpper(strings.TrimSpace(c)))
}

func mapStorageItemPrototype(proto readmodels.StorageItemPrototype) StorageItemPrototypeDTO {
	var buff *BuffStorageDataDTO
	if proto.BuffData != nil {
		buff = &BuffStorageDataDTO{
			Type:            BuffType(proto.BuffData.Type),
			Value:           proto.BuffData.Value,
			DurationSeconds: proto.BuffData.DurationSeconds,
		}
	}
	var intel *IntelStorageDataDTO
	if proto.IntelData != nil {
		intel = &IntelStorageDataDTO{
			Type:              HiddenLocationType(proto.IntelData.Type),
			DecryptionSeconds: proto.IntelData.DecryptionSeconds,
		}
	}
	var dmg *DamagedStorageDataDTO
	if proto.DamagedData != nil {
		dmg = &DamagedStorageDataDTO{
			RestorePrice:   PriceModelFromReadModel(proto.DamagedData.RestorePrice),
			OriginalUnitID: proto.DamagedData.OriginalUnitID,
		}
	}
	var art *ArtifactStorageDataDTO
	if proto.ArtifactData != nil {
		art = &ArtifactStorageDataDTO{
			Type:  ArtifactEffectType(proto.ArtifactData.Type),
			Value: proto.ArtifactData.Value,
		}
	}
	var cons *ConsumableStorageDataDTO
	if proto.ConsumableData != nil {
		contents := make([]ConsumableBoxContents, len(proto.ConsumableData.BoxContents))
		for i, c := range proto.ConsumableData.BoxContents {
			contents[i] = ConsumableBoxContents(c)
		}
		cons = &ConsumableStorageDataDTO{
			Type:        ConsumableType(proto.ConsumableData.Type),
			BoxContents: contents,
			BoxSize:     proto.ConsumableData.BoxSize,
		}
	}

	return StorageItemPrototypeDTO{
		ID:               proto.ID,
		Name:             proto.Name,
		Category:         StorageCategory(proto.Category),
		ShortDescription: proto.ShortDescription,
		FullDescription:  proto.FullDescription,
		ImageURL:         proto.ImageURL,
		BuffData:         buff,
		IntelData:        intel,
		DamagedData:      dmg,
		ArtifactData:     art,
		ConsumableData:   cons,
	}
}

func StorageItemsPresentFromReadModels(items []*readmodels.StorageItemPresent) []StorageItemPresentDTO {
	out := make([]StorageItemPresentDTO, 0, len(items))
	for _, item := range items {
		out = append(out, StorageItemPresentDTO{
			BaseOwnedItemDTO: BaseOwnedItemDTOFromReadModel(item.BaseOwnedItem),
			Prototype:        mapStorageItemPrototype(item.Prototype),
			ExpiresAt:        item.ExpiresAt,
			IsActive:         item.IsActive,
		})
	}
	return out
}
