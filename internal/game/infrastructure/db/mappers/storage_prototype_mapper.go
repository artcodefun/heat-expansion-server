package mappers

import (
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
)

func StoragePrototypeFromDB(p gen.StorageItemPrototype) *domain.StorageItemPrototype {
	var buff *domain.BuffStorageData
	if p.BuffData.Valid {
		var dto dtos.BuffStorageDataDTO
		unmarshalIfValid(p.BuffData, &dto)
		buff = dtos.BuffStorageDataFromDTO(&dto)
	}
	var intel *domain.IntelStorageData
	if p.IntelData.Valid {
		var dto dtos.IntelStorageDataDTO
		unmarshalIfValid(p.IntelData, &dto)
		intel = dtos.IntelStorageDataFromDTO(&dto)
	}
	var dmg *domain.DamagedStorageData
	if p.DamagedData.Valid {
		var dto dtos.DamagedStorageDataDTO
		unmarshalIfValid(p.DamagedData, &dto)
		dmg = dtos.DamagedStorageDataFromDTO(&dto)
	}
	var art *domain.ArtifactStorageData
	if p.ArtifactData.Valid {
		var dto dtos.ArtifactStorageDataDTO
		unmarshalIfValid(p.ArtifactData, &dto)
		art = dtos.ArtifactStorageDataFromDTO(&dto)
	}
	var cons *domain.ConsumableStorageData
	if p.ConsumableData.Valid {
		var dto dtos.ConsumableStorageDataDTO
		unmarshalIfValid(p.ConsumableData, &dto)
		cons = dtos.ConsumableStorageDataFromDTO(&dto)
	}

	proto := &domain.StorageItemPrototype{
		ID:               int(p.ID),
		Name:             p.Name,
		Category:         domain.StorageCategory(p.Category),
		CreationSources:  creationSourcesFromJSON(p.CreationSources),
		EstimatedWorth:   int(p.EstimatedWorth),
		ShortDescription: nullStringToString(&p.ShortDescription.String, p.ShortDescription.Valid),
		FullDescription:  nullStringToString(&p.FullDescription.String, p.FullDescription.Valid),
		ImageURL:         nullStringToString(&p.ImageUrl.String, p.ImageUrl.Valid),
		BuffData:         buff,
		IntelData:        intel,
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
