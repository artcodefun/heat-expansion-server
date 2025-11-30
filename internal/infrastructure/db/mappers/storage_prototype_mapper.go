package mappers

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/gen"
)

func StoragePrototypeFromDB(p gen.StorageItemPrototype) *domain.StorageItemPrototype {
	var buff *domain.BuffStorageData
	if p.BuffData.Valid {
		var tmp domain.BuffStorageData
		unmarshalIfValid(p.BuffData, &tmp)
		buff = &tmp
	}
	var mp *domain.MapStorageData
	if p.MapData.Valid {
		var tmp domain.MapStorageData
		unmarshalIfValid(p.MapData, &tmp)
		mp = &tmp
	}
	var dmg *domain.DamagedStorageData
	if p.DamagedData.Valid {
		var tmp domain.DamagedStorageData
		unmarshalIfValid(p.DamagedData, &tmp)
		dmg = &tmp
	}
	var art *domain.ArtifactStorageData
	if p.ArtifactData.Valid {
		var tmp domain.ArtifactStorageData
		unmarshalIfValid(p.ArtifactData, &tmp)
		art = &tmp
	}
	var cons *domain.ConsumableStorageData
	if p.ConsumableData.Valid {
		var tmp domain.ConsumableStorageData
		unmarshalIfValid(p.ConsumableData, &tmp)
		cons = &tmp
	}

	proto := &domain.StorageItemPrototype{
		ID:               int(p.ID),
		Name:             p.Name,
		Category:         domain.StorageCategory(p.Category),
		ShortDescription: nullStringToString(&p.ShortDescription.String, p.ShortDescription.Valid),
		FullDescription:  nullStringToString(&p.FullDescription.String, p.FullDescription.Valid),
		ImageURL:         nullStringToString(&p.ImageUrl.String, p.ImageUrl.Valid),
		BuffData:         buff,
		MapData:          mp,
		DamagedData:      dmg,
		ArtifactData:     art,
		ConsumableData:   cons,
	}
	return proto
}

func StoragePrototypesFromDB(src []gen.StorageItemPrototype) []*domain.StorageItemPrototype {
	dst := make([]*domain.StorageItemPrototype, len(src))
	for i, p := range src {
		dst[i] = StoragePrototypeFromDB(p)
	}
	return dst
}
