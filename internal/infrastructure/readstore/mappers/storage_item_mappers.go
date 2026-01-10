package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/gen"
	"github.com/sqlc-dev/pqtype"
)

func StoragePrototypeFromModel(r gen.StorageItemPrototype) readmodels.StorageItemPrototype {
	return readmodels.StorageItemPrototype{
		ID:               int(r.ID),
		Name:             r.Name,
		Category:         readmodels.StorageCategory(r.Category),
		ShortDescription: nullString(r.ShortDescription),
		FullDescription:  nullString(r.FullDescription),
		ImageURL:         nullString(r.ImageUrl),
		BuffData:         buffStorageDataFromJSON(r.BuffData),
		IntelData:        intelStorageDataFromJSON(r.IntelData),
		DamagedData:      damagedStorageDataFromJSON(r.DamagedData),
		ArtifactData:     artifactStorageDataFromJSON(r.ArtifactData),
		ConsumableData:   consumableStorageDataFromJSON(r.ConsumableData),
	}
}

func StoragePrototypeFromPresentRow(r gen.ListPresentStorageItemsRow) readmodels.StorageItemPrototype {
	return readmodels.StorageItemPrototype{
		ID:               int(r.ProtoID),
		Name:             r.Name,
		Category:         readmodels.StorageCategory(r.Category),
		ShortDescription: nullString(r.ShortDescription),
		FullDescription:  nullString(r.FullDescription),
		ImageURL:         nullString(r.ImageUrl),
		BuffData:         buffStorageDataFromJSON(r.BuffData),
		IntelData:        intelStorageDataFromJSON(r.IntelData),
		DamagedData:      damagedStorageDataFromJSON(r.DamagedData),
		ArtifactData:     artifactStorageDataFromJSON(r.ArtifactData),
		ConsumableData:   consumableStorageDataFromJSON(r.ConsumableData),
	}
}

func StorageItemPresentFromRow(r gen.ListPresentStorageItemsRow) readmodels.StorageItemPresent {
	var jd dtos.StoragePresentDTO
	if r.PresentData.Valid {
		_ = json.Unmarshal(r.PresentData.RawMessage, &jd)
	}
	return readmodels.StorageItemPresent{
		BaseOwnedItem: readmodels.BaseOwnedItem{ID: r.ID, UserBaseID: int(r.BaseID)},
		Prototype:     StoragePrototypeFromPresentRow(r),
		ExpiresAt:     jd.ExpiresAt,
		IsActive:      jd.IsActive,
	}
}

// Storage prototype detail helpers: JSONB (DTO shape) -> readmodels.* data

func buffStorageDataFromJSON(nm pqtype.NullRawMessage) *readmodels.BuffStorageData {
	if !nm.Valid {
		return nil
	}
	var d dtos.BuffStorageDataDTO
	if err := json.Unmarshal(nm.RawMessage, &d); err != nil {
		return nil
	}
	return &readmodels.BuffStorageData{
		Type:            readmodels.BuffType(d.Type),
		Value:           d.Value,
		DurationSeconds: d.DurationSeconds,
	}
}

func intelStorageDataFromJSON(nm pqtype.NullRawMessage) *readmodels.IntelStorageData {
	if !nm.Valid {
		return nil
	}
	var d dtos.IntelStorageDataDTO
	if err := json.Unmarshal(nm.RawMessage, &d); err != nil {
		return nil
	}
	return &readmodels.IntelStorageData{
		Type:              readmodels.HiddenLocationType(d.Type),
		DecryptionSeconds: d.DecryptionSeconds,
	}
}

func damagedStorageDataFromJSON(nm pqtype.NullRawMessage) *readmodels.DamagedStorageData {
	if !nm.Valid {
		return nil
	}
	var d dtos.DamagedStorageDataDTO
	if err := json.Unmarshal(nm.RawMessage, &d); err != nil {
		return nil
	}
	return &readmodels.DamagedStorageData{
		RestorePrice: readmodels.PriceModel{
			Credits:    d.RestorePrice.Credits,
			Iron:       d.RestorePrice.Iron,
			Titanium:   d.RestorePrice.Titanium,
			Antimatter: d.RestorePrice.Antimatter,
		},
		RestorationSeconds: d.RestorationSeconds,
		OriginalUnitID:     d.OriginalUnitID,
	}
}

func artifactStorageDataFromJSON(nm pqtype.NullRawMessage) *readmodels.ArtifactStorageData {
	if !nm.Valid {
		return nil
	}
	var d dtos.ArtifactStorageDataDTO
	if err := json.Unmarshal(nm.RawMessage, &d); err != nil {
		return nil
	}
	return &readmodels.ArtifactStorageData{
		Type:  readmodels.ArtifactEffectType(d.Type),
		Value: d.Value,
	}
}

func consumableStorageDataFromJSON(nm pqtype.NullRawMessage) *readmodels.ConsumableStorageData {
	if !nm.Valid {
		return nil
	}
	var d dtos.ConsumableStorageDataDTO
	if err := json.Unmarshal(nm.RawMessage, &d); err != nil {
		return nil
	}
	contents := make([]readmodels.ConsumableBoxContents, len(d.BoxContents))
	for i, c := range d.BoxContents {
		contents[i] = readmodels.ConsumableBoxContents(c)
	}
	return &readmodels.ConsumableStorageData{
		Type:        readmodels.ConsumableType(d.Type),
		BoxContents: contents,
		BoxSize:     d.BoxSize,
	}
}
