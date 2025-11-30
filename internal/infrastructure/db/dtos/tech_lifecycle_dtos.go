package dtos

import "github.com/artcodefun/heat-expansion-api/internal/core/domain"

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
}

func TechDoneDTOFromDomain(t domain.TechItemDone) TechDoneDTO {
	return TechDoneDTO{ResearchedAt: t.ResearchedAt}
}
func TechDoneFromDTO(d TechDoneDTO, owned domain.BaseOwnedItem, proto domain.TechItemPrototype) domain.TechItemDone {
	return domain.TechItemDone{BaseOwnedItem: owned, Prototype: proto, ResearchedAt: d.ResearchedAt}
}

type TechnologyEffectDTO struct {
	Type  string `json:"type"`
	Value int    `json:"value"`
}

func TechnologyEffectDTOFromDomain(e domain.TechnologyEffect) TechnologyEffectDTO {
	return TechnologyEffectDTO{Type: string(e.EffectType), Value: e.Value}
}
func TechnologyEffectFromDTO(d TechnologyEffectDTO) domain.TechnologyEffect {
	return domain.TechnologyEffect{EffectType: domain.EffectType(d.Type), Value: d.Value}
}
