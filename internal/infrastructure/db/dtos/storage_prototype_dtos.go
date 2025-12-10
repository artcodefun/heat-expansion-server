package dtos

import "github.com/artcodefun/heat-expansion-api/internal/core/domain"

// BuffStorageDataDTO represents JSON shape for buff storage items in prototypes.
type BuffStorageDataDTO struct {
	SpaceCapacityBonus int    `json:"space_capacity_bonus"`
	AttackBonus        int    `json:"attack_bonus"`
	DefenceBonus       int    `json:"defence_bonus"`
	DurationSeconds    int64  `json:"duration_seconds"`
	ActivatedAt        *int64 `json:"activated_at,omitempty"`
}

func BuffStorageDataDTOFromDomain(d *domain.BuffStorageData) *BuffStorageDataDTO {
	if d == nil {
		return nil
	}
	return &BuffStorageDataDTO{
		SpaceCapacityBonus: d.SpaceCapacityBonus,
		AttackBonus:        d.AttackBonus,
		DefenceBonus:       d.DefenceBonus,
		DurationSeconds:    d.DurationSeconds,
		ActivatedAt:        d.ActivatedAt,
	}
}

func BuffStorageDataFromDTO(d *BuffStorageDataDTO) *domain.BuffStorageData {
	if d == nil {
		return nil
	}
	return &domain.BuffStorageData{
		SpaceCapacityBonus: d.SpaceCapacityBonus,
		AttackBonus:        d.AttackBonus,
		DefenceBonus:       d.DefenceBonus,
		DurationSeconds:    d.DurationSeconds,
		ActivatedAt:        d.ActivatedAt,
	}
}

// MapStorageDataDTO represents JSON shape for map storage items in prototypes.
type MapStorageDataDTO struct {
	RevealedArea string `json:"revealed_area"`
	ScanRange    int    `json:"scan_range"`
}

func MapStorageDataDTOFromDomain(d *domain.MapStorageData) *MapStorageDataDTO {
	if d == nil {
		return nil
	}
	return &MapStorageDataDTO{
		RevealedArea: d.RevealedArea,
		ScanRange:    d.ScanRange,
	}
}

func MapStorageDataFromDTO(d *MapStorageDataDTO) *domain.MapStorageData {
	if d == nil {
		return nil
	}
	return &domain.MapStorageData{
		RevealedArea: d.RevealedArea,
		ScanRange:    d.ScanRange,
	}
}

// DamagedStorageDataDTO represents JSON shape for damaged storage items in prototypes.
type DamagedStorageDataDTO struct {
	RestorePrice   PriceDTO `json:"restore_price"`
	OriginalUnitID int      `json:"original_unit_id"`
	DamageLevel    int      `json:"damage_level"`
}

func DamagedStorageDataDTOFromDomain(d *domain.DamagedStorageData) *DamagedStorageDataDTO {
	if d == nil {
		return nil
	}
	return &DamagedStorageDataDTO{
		RestorePrice:   PriceDTOFromDomain(d.RestorePrice),
		OriginalUnitID: d.OriginalUnitID,
		DamageLevel:    d.DamageLevel,
	}
}

func DamagedStorageDataFromDTO(d *DamagedStorageDataDTO) *domain.DamagedStorageData {
	if d == nil {
		return nil
	}
	return &domain.DamagedStorageData{
		RestorePrice:   PriceFromDTO(d.RestorePrice),
		OriginalUnitID: d.OriginalUnitID,
		DamageLevel:    d.DamageLevel,
	}
}

// ArtifactStorageDataDTO represents JSON shape for artifact storage items in prototypes.
type ArtifactStorageDataDTO struct {
	PassiveEffect string `json:"passive_effect"`
	Rarity        string `json:"rarity"`
	Lore          string `json:"lore"`
}

func ArtifactStorageDataDTOFromDomain(d *domain.ArtifactStorageData) *ArtifactStorageDataDTO {
	if d == nil {
		return nil
	}
	return &ArtifactStorageDataDTO{
		PassiveEffect: d.PassiveEffect,
		Rarity:        d.Rarity,
		Lore:          d.Lore,
	}
}

func ArtifactStorageDataFromDTO(d *ArtifactStorageDataDTO) *domain.ArtifactStorageData {
	if d == nil {
		return nil
	}
	return &domain.ArtifactStorageData{
		PassiveEffect: d.PassiveEffect,
		Rarity:        d.Rarity,
		Lore:          d.Lore,
	}
}

// ConsumableStorageDataDTO represents JSON shape for consumable storage items in prototypes.
type ConsumableStorageDataDTO struct {
	EffectType   string    `json:"effect_type"`
	Uses         int       `json:"uses"`
	RestorePrice *PriceDTO `json:"restore_price,omitempty"`
}

func ConsumableStorageDataDTOFromDomain(d *domain.ConsumableStorageData) *ConsumableStorageDataDTO {
	if d == nil {
		return nil
	}
	var priceDTO *PriceDTO
	if d.RestorePrice != nil {
		v := PriceDTOFromDomain(*d.RestorePrice)
		priceDTO = &v
	}
	return &ConsumableStorageDataDTO{
		EffectType:   d.EffectType,
		Uses:         d.Uses,
		RestorePrice: priceDTO,
	}
}

func ConsumableStorageDataFromDTO(d *ConsumableStorageDataDTO) *domain.ConsumableStorageData {
	if d == nil {
		return nil
	}
	var price *domain.PriceModel
	if d.RestorePrice != nil {
		v := PriceFromDTO(*d.RestorePrice)
		price = &v
	}
	return &domain.ConsumableStorageData{
		EffectType:   d.EffectType,
		Uses:         d.Uses,
		RestorePrice: price,
	}
}
