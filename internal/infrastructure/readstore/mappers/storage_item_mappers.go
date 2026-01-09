package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/gen"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

func StoragePrototypeFromPresentRow(r gen.ListPresentStorageItemsRow) readmodels.StorageItemPrototype {
	return readmodels.StorageItemPrototype{
		ID:               int(r.ProtoID),
		Name:             r.Name,
		Category:         readmodels.StorageCategory(r.Category),
		ShortDescription: nullString(r.ShortDescription),
		FullDescription:  nullString(r.FullDescription),
		ImageURL:         nullString(r.ImageUrl),
		BuffData:         buffStorageDataFromJSON(r.BuffData),
		MapData:          mapStorageDataFromJSON(r.MapData),
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
	refund := readmodels.PriceModel{Credits: jd.Refund.Credits, Iron: jd.Refund.Iron, Titanium: jd.Refund.Titanium, Antimatter: jd.Refund.Antimatter}
	return readmodels.StorageItemPresent{
		BaseOwnedItem: readmodels.BaseOwnedItem{ID: uuid.UUID(r.ID), UserBaseID: int(r.BaseID)},
		Prototype:     StoragePrototypeFromPresentRow(r),
		Refund:        refund,
		ActivatedAt:   jd.ActivatedAt,
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
		SpaceCapacityBonus: d.SpaceCapacityBonus,
		AttackBonus:        d.AttackBonus,
		DefenceBonus:       d.DefenceBonus,
		DurationSeconds:    d.DurationSeconds,
	}
}

func mapStorageDataFromJSON(nm pqtype.NullRawMessage) *readmodels.MapStorageData {
	if !nm.Valid {
		return nil
	}
	var d dtos.MapStorageDataDTO
	if err := json.Unmarshal(nm.RawMessage, &d); err != nil {
		return nil
	}
	return &readmodels.MapStorageData{
		RevealedArea: d.RevealedArea,
		ScanRange:    d.ScanRange,
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
		OriginalUnitID: d.OriginalUnitID,
		DamageLevel:    d.DamageLevel,
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
		PassiveEffect: d.PassiveEffect,
		Rarity:        d.Rarity,
		Lore:          d.Lore,
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
	var price *readmodels.PriceModel
	if d.RestorePrice != nil {
		price = &readmodels.PriceModel{
			Credits:    d.RestorePrice.Credits,
			Iron:       d.RestorePrice.Iron,
			Titanium:   d.RestorePrice.Titanium,
			Antimatter: d.RestorePrice.Antimatter,
		}
	}
	return &readmodels.ConsumableStorageData{
		EffectType:   d.EffectType,
		Uses:         d.Uses,
		RestorePrice: price,
	}
}
