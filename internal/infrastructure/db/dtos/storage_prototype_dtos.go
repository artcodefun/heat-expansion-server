package dtos

import "github.com/artcodefun/heat-expansion-api/internal/core/domain"

// BuffStorageDataDTO represents JSON shape for buff storage items in prototypes.
type BuffStorageDataDTO struct {
	Type            string  `json:"type"`
	Value           float32 `json:"value"`
	DurationSeconds int64   `json:"duration_seconds"`
}

func BuffStorageDataDTOFromDomain(d *domain.BuffStorageData) *BuffStorageDataDTO {
	if d == nil {
		return nil
	}
	return &BuffStorageDataDTO{
		Type:            string(d.Type),
		Value:           d.Value,
		DurationSeconds: d.DurationSeconds,
	}
}

func BuffStorageDataFromDTO(d *BuffStorageDataDTO) *domain.BuffStorageData {
	if d == nil {
		return nil
	}
	return &domain.BuffStorageData{
		Type:            domain.BuffType(d.Type),
		Value:           d.Value,
		DurationSeconds: d.DurationSeconds,
	}
}

// IntelStorageDataDTO represents JSON shape for intel storage items in prototypes.
type IntelStorageDataDTO struct {
	Type              string `json:"type"`
	DecryptionSeconds int64  `json:"decryption_seconds"`
}

func IntelStorageDataDTOFromDomain(d *domain.IntelStorageData) *IntelStorageDataDTO {
	if d == nil {
		return nil
	}
	return &IntelStorageDataDTO{
		Type:              string(d.Type),
		DecryptionSeconds: d.DecryptionSeconds,
	}
}

func IntelStorageDataFromDTO(d *IntelStorageDataDTO) *domain.IntelStorageData {
	if d == nil {
		return nil
	}
	return &domain.IntelStorageData{
		Type:              domain.HiddenLocationType(d.Type),
		DecryptionSeconds: d.DecryptionSeconds,
	}
}

// DamagedStorageDataDTO represents JSON shape for damaged storage items in prototypes.
type DamagedStorageDataDTO struct {
	RestorePrice       PriceDTO `json:"restore_price"`
	RestorationSeconds int64    `json:"restoration_seconds"`
	OriginalUnitID     int      `json:"original_unit_id"`
}

func DamagedStorageDataDTOFromDomain(d *domain.DamagedStorageData) *DamagedStorageDataDTO {
	if d == nil {
		return nil
	}
	return &DamagedStorageDataDTO{
		RestorePrice:       PriceDTOFromDomain(d.RestorePrice),
		RestorationSeconds: d.RestorationSeconds,
		OriginalUnitID:     d.OriginalUnitID,
	}
}

func DamagedStorageDataFromDTO(d *DamagedStorageDataDTO) *domain.DamagedStorageData {
	if d == nil {
		return nil
	}
	return &domain.DamagedStorageData{
		RestorePrice:       PriceFromDTO(d.RestorePrice),
		RestorationSeconds: d.RestorationSeconds,
		OriginalUnitID:     d.OriginalUnitID,
	}
}

// ArtifactStorageDataDTO represents JSON shape for artifact storage items in prototypes.
type ArtifactStorageDataDTO struct {
	Type  string  `json:"type"`
	Value float32 `json:"value"`
}

func ArtifactStorageDataDTOFromDomain(d *domain.ArtifactStorageData) *ArtifactStorageDataDTO {
	if d == nil {
		return nil
	}
	return &ArtifactStorageDataDTO{
		Type:  string(d.Type),
		Value: d.Value,
	}
}

func ArtifactStorageDataFromDTO(d *ArtifactStorageDataDTO) *domain.ArtifactStorageData {
	if d == nil {
		return nil
	}
	return &domain.ArtifactStorageData{
		Type:  domain.ArtifactEffectType(d.Type),
		Value: d.Value,
	}
}

// ConsumableStorageDataDTO represents JSON shape for consumable storage items in prototypes.
type ConsumableStorageDataDTO struct {
	Type        string   `json:"type"`
	BoxContents []string `json:"box_contents"`
	BoxSize     int      `json:"box_size"`
}

func ConsumableStorageDataDTOFromDomain(d *domain.ConsumableStorageData) *ConsumableStorageDataDTO {
	if d == nil {
		return nil
	}
	contents := make([]string, len(d.BoxContents))
	for i, c := range d.BoxContents {
		contents[i] = string(c)
	}
	return &ConsumableStorageDataDTO{
		Type:        string(d.Type),
		BoxContents: contents,
		BoxSize:     d.BoxSize,
	}
}

func ConsumableStorageDataFromDTO(d *ConsumableStorageDataDTO) *domain.ConsumableStorageData {
	if d == nil {
		return nil
	}
	contents := make([]domain.ConsumableBoxContents, len(d.BoxContents))
	for i, c := range d.BoxContents {
		contents[i] = domain.ConsumableBoxContents(c)
	}
	return &domain.ConsumableStorageData{
		Type:        domain.ConsumableType(d.Type),
		BoxContents: contents,
		BoxSize:     d.BoxSize,
	}
}
