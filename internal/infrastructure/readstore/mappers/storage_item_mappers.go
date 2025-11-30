package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/gen"
	"github.com/google/uuid"
)

func StoragePrototypeFromPresentRow(r gen.ListPresentStorageItemsRow) readmodels.StorageItemPrototype {
	return readmodels.StorageItemPrototype{
		ID:               int(r.ProtoID),
		Name:             r.Name,
		Category:         readmodels.StorageCategory(r.Category),
		ShortDescription: nullString(r.ShortDescription),
		FullDescription:  nullString(r.FullDescription),
		ImageURL:         nullString(r.ImageUrl),
		BuffData:         jsonToNullRaw[readmodels.BuffStorageData](r.BuffData),
		MapData:          jsonToNullRaw[readmodels.MapStorageData](r.MapData),
		DamagedData:      jsonToNullRaw[readmodels.DamagedStorageData](r.DamagedData),
		ArtifactData:     jsonToNullRaw[readmodels.ArtifactStorageData](r.ArtifactData),
		ConsumableData:   jsonToNullRaw[readmodels.ConsumableStorageData](r.ConsumableData),
	}
}

func StorageItemPresentFromRow(r gen.ListPresentStorageItemsRow) readmodels.StorageItemPresent {
	var jd dtos.StoragePresentDTO
	if r.PresentData.Valid {
		_ = json.Unmarshal(r.PresentData.RawMessage, &jd)
	}
	refund := readmodels.PriceModel{Credits: jd.Refund.Credits, Iron: jd.Refund.Iron, Titanium: jd.Refund.Titanium, Antimatter: jd.Refund.Antimatter}
	return readmodels.StorageItemPresent{BaseOwnedItem: readmodels.BaseOwnedItem{ID: uuid.UUID(r.ID), UserBaseID: int(r.BaseID)}, Prototype: StoragePrototypeFromPresentRow(r), Refund: refund}
}
