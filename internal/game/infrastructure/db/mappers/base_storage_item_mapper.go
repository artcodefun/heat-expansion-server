package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
)

func HydrateStorageItems(base *domain.UserBaseModel, rows []gen.BaseStorageItem, proto map[int]*domain.StorageItemPrototype) {
	for _, r := range rows {
		p := proto[int(r.PrototypeID)]
		if p == nil {
			continue
		}
		owned := domain.BaseOwnedItem{ID: r.ID, UserBaseID: base.ID}
		// Only PRESENT supported currently
		var d dtos.StoragePresentDTO
		unmarshalIfValid(r.PresentData, &d)
		base.StorageItemsPresent = append(base.StorageItemsPresent, dtos.StoragePresentFromDTO(d, owned, *p))
	}
}

// DehydrateStorageItems converts present storage items into insert params
// for the base_storage_items table.
func DehydrateStorageItems(base *domain.UserBaseModel) []gen.InsertBaseStorageItemParams {
	now := domain.NowUnix()
	out := make([]gen.InsertBaseStorageItemParams, 0, len(base.StorageItemsPresent))

	for _, it := range base.StorageItemsPresent {
		presentRaw := BuildStoragePresentRaw(it)
		stateJSON, _ := json.Marshal(map[string]any{}) // placeholder empty state
		out = append(out, gen.InsertBaseStorageItemParams{
			ID:          it.ID,
			BaseID:      int64(base.ID),
			PrototypeID: int64(it.Prototype.ID),
			Status:      "PRESENT",
			PresentData: presentRaw,
			State:       stateJSON,
			CreatedAt:   now,
		})
	}
	return out
}
