package mappers

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/gen"
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
