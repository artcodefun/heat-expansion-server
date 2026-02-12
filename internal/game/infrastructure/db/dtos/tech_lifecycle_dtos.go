package dtos

import "github.com/artcodefun/heat-expansion-api/internal/game/core/domain"

type TechInProgressDTO struct {
	StartDate         int64 `json:"start_date"`
	CompletionDate    int64 `json:"completion_date"`
	CrystalsSkipPrice int   `json:"crystals_skip_price"`
}

func TechInProgressDTOFromDomain(t domain.TechItemInProgress) TechInProgressDTO {
	return TechInProgressDTO{StartDate: t.StartDate, CompletionDate: t.CompletionDate, CrystalsSkipPrice: t.CrystalsSkipPrice}
}
func TechInProgressFromDTO(d TechInProgressDTO, owned domain.BaseOwnedItem, proto domain.TechItemPrototype) domain.TechItemInProgress {
	return domain.TechItemInProgress{BaseOwnedItem: owned, Prototype: proto, StartDate: d.StartDate, CompletionDate: d.CompletionDate, CrystalsSkipPrice: d.CrystalsSkipPrice}
}

type TechDoneDTO struct {
	ResearchedAt int64 `json:"researched_at"`
	Level        int   `json:"level"`
}

func TechDoneDTOFromDomain(t domain.TechItemDone) TechDoneDTO {
	return TechDoneDTO{ResearchedAt: t.ResearchedAt, Level: t.Level}
}
func TechDoneFromDTO(d TechDoneDTO, owned domain.BaseOwnedItem, proto domain.TechItemPrototype) domain.TechItemDone {
	return domain.TechItemDone{BaseOwnedItem: owned, Prototype: proto, ResearchedAt: d.ResearchedAt, Level: d.Level}
}
