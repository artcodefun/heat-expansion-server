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

	return &domain.StorageItemPrototype{
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
}

func StoragePrototypesFromDB(src []gen.StorageItemPrototype) []*domain.StorageItemPrototype {
	dst := make([]*domain.StorageItemPrototype, len(src))
	for i, p := range src {
		dst[i] = StoragePrototypeFromDB(p)
	}
	return dst
}

// StoragePrototypeToCreateParams maps a domain prototype to the sqlc insert params.
// The category-specific data blocks are serialized via their DTO shapes; a nil
// pointer becomes SQL NULL.
func StoragePrototypeToCreateParams(p *domain.StorageItemPrototype) gen.CreateStoragePrototypeParams {
	return gen.CreateStoragePrototypeParams{
		ID:               int64(p.ID),
		Name:             string(p.Name),
		Category:         string(p.Category),
		EstimatedWorth:   int32(p.EstimatedWorth),
		ShortDescription: stringToNullString(string(p.ShortDescription)),
		FullDescription:  stringToNullString(string(p.FullDescription)),
		ImageUrl:         stringToNullString(p.ImageURL),
		BuffData:         toNullRawMessage(dtos.BuffStorageDataDTOFromDomain(p.BuffData)),
		IntelData:        toNullRawMessage(dtos.IntelStorageDataDTOFromDomain(p.IntelData)),
		DamagedData:      toNullRawMessage(dtos.DamagedStorageDataDTOFromDomain(p.DamagedData)),
		ArtifactData:     toNullRawMessage(dtos.ArtifactStorageDataDTOFromDomain(p.ArtifactData)),
		ConsumableData:   toNullRawMessage(dtos.ConsumableStorageDataDTOFromDomain(p.ConsumableData)),
		CreationSources:  creationSourcesToJSON(p.CreationSources),
	}
}

// StoragePrototypeToUpdateParams maps a domain prototype to the sqlc update params,
// keyed by p.ID.
func StoragePrototypeToUpdateParams(p *domain.StorageItemPrototype) gen.UpdateStoragePrototypeParams {
	return gen.UpdateStoragePrototypeParams{
		ID:               int64(p.ID),
		Name:             string(p.Name),
		Category:         string(p.Category),
		EstimatedWorth:   int32(p.EstimatedWorth),
		ShortDescription: stringToNullString(string(p.ShortDescription)),
		FullDescription:  stringToNullString(string(p.FullDescription)),
		ImageUrl:         stringToNullString(p.ImageURL),
		BuffData:         toNullRawMessage(dtos.BuffStorageDataDTOFromDomain(p.BuffData)),
		IntelData:        toNullRawMessage(dtos.IntelStorageDataDTOFromDomain(p.IntelData)),
		DamagedData:      toNullRawMessage(dtos.DamagedStorageDataDTOFromDomain(p.DamagedData)),
		ArtifactData:     toNullRawMessage(dtos.ArtifactStorageDataDTOFromDomain(p.ArtifactData)),
		ConsumableData:   toNullRawMessage(dtos.ConsumableStorageDataDTOFromDomain(p.ConsumableData)),
		CreationSources:  creationSourcesToJSON(p.CreationSources),
	}
}
