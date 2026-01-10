package dtos

import "github.com/artcodefun/heat-expansion-api/internal/core/domain"

type SpyResultDTO struct {
	Outcome           string            `json:"outcome"`
	AttackerRemaining []MilitaryUnitDTO `json:"attacker_remaining"`
	DefenderRemaining []MilitaryUnitDTO `json:"defender_remaining"`
	DefendersBefore   []MilitaryUnitDTO `json:"defenders_before"`
}

type AttackResultDTO struct {
	Outcome             string                `json:"outcome"`
	AttackerRemaining   []MilitaryUnitDTO     `json:"attacker_remaining"`
	DefenderRemaining   []MilitaryUnitDTO     `json:"defender_remaining"`
	RemainingStructures []DefenseStructureDTO `json:"remaining_structures"`
	Loot                PriceDTO              `json:"loot"`
	Trophies            []TrophyDTO           `json:"trophies"`
	DefendersBefore     []MilitaryUnitDTO     `json:"defenders_before"`
	StructuresBefore    []DefenseStructureDTO `json:"structures_before"`
}

func SpyResultDTOFromDomain(s *domain.SpyResult) *SpyResultDTO {
	if s == nil {
		return nil
	}
	out := &SpyResultDTO{Outcome: string(s.Outcome)}
	if len(s.AttackerRemaining) > 0 {
		out.AttackerRemaining = make([]MilitaryUnitDTO, 0, len(s.AttackerRemaining))
		for _, u := range s.AttackerRemaining {
			out.AttackerRemaining = append(out.AttackerRemaining, MilitaryUnitDTOFromDomain(u))
		}
	}
	if len(s.DefenderRemaining) > 0 {
		out.DefenderRemaining = make([]MilitaryUnitDTO, 0, len(s.DefenderRemaining))
		for _, u := range s.DefenderRemaining {
			out.DefenderRemaining = append(out.DefenderRemaining, MilitaryUnitDTOFromDomain(u))
		}
	}
	if len(s.DefendersBefore) > 0 {
		out.DefendersBefore = make([]MilitaryUnitDTO, 0, len(s.DefendersBefore))
		for _, u := range s.DefendersBefore {
			out.DefendersBefore = append(out.DefendersBefore, MilitaryUnitDTOFromDomain(u))
		}
	}
	return out
}
func SpyResultFromDTO(d *SpyResultDTO) *domain.SpyResult {
	if d == nil {
		return nil
	}
	out := &domain.SpyResult{Outcome: domain.SpyOutcome(d.Outcome)}
	if len(d.AttackerRemaining) > 0 {
		out.AttackerRemaining = make([]domain.MilitaryUnitSnap, 0, len(d.AttackerRemaining))
		for _, u := range d.AttackerRemaining {
			out.AttackerRemaining = append(out.AttackerRemaining, MilitaryUnitFromDTO(u))
		}
	}
	if len(d.DefenderRemaining) > 0 {
		out.DefenderRemaining = make([]domain.MilitaryUnitSnap, 0, len(d.DefenderRemaining))
		for _, u := range d.DefenderRemaining {
			out.DefenderRemaining = append(out.DefenderRemaining, MilitaryUnitFromDTO(u))
		}
	}
	if len(d.DefendersBefore) > 0 {
		out.DefendersBefore = make([]domain.MilitaryUnitSnap, 0, len(d.DefendersBefore))
		for _, u := range d.DefendersBefore {
			out.DefendersBefore = append(out.DefendersBefore, MilitaryUnitFromDTO(u))
		}
	}
	return out
}

func AttackResultDTOFromDomain(a *domain.AttackResult) *AttackResultDTO {
	if a == nil {
		return nil
	}
	out := &AttackResultDTO{Outcome: string(a.Outcome), Loot: PriceDTOFromDomain(a.Loot)}
	if len(a.AttackerRemaining) > 0 {
		out.AttackerRemaining = make([]MilitaryUnitDTO, 0, len(a.AttackerRemaining))
		for _, u := range a.AttackerRemaining {
			out.AttackerRemaining = append(out.AttackerRemaining, MilitaryUnitDTOFromDomain(u))
		}
	}
	if len(a.DefenderRemaining) > 0 {
		out.DefenderRemaining = make([]MilitaryUnitDTO, 0, len(a.DefenderRemaining))
		for _, u := range a.DefenderRemaining {
			out.DefenderRemaining = append(out.DefenderRemaining, MilitaryUnitDTOFromDomain(u))
		}
	}
	if len(a.RemainingStructures) > 0 {
		out.RemainingStructures = make([]DefenseStructureDTO, 0, len(a.RemainingStructures))
		for _, s := range a.RemainingStructures {
			out.RemainingStructures = append(out.RemainingStructures, DefenseStructureDTOFromDomain(s))
		}
	}
	if len(a.Trophies) > 0 {
		out.Trophies = make([]TrophyDTO, 0, len(a.Trophies))
		for _, t := range a.Trophies {
			out.Trophies = append(out.Trophies, TrophyDTOFromDomain(t))
		}
	}
	if len(a.DefendersBefore) > 0 {
		out.DefendersBefore = make([]MilitaryUnitDTO, 0, len(a.DefendersBefore))
		for _, u := range a.DefendersBefore {
			out.DefendersBefore = append(out.DefendersBefore, MilitaryUnitDTOFromDomain(u))
		}
	}
	if len(a.StructuresBefore) > 0 {
		out.StructuresBefore = make([]DefenseStructureDTO, 0, len(a.StructuresBefore))
		for _, s := range a.StructuresBefore {
			out.StructuresBefore = append(out.StructuresBefore, DefenseStructureDTOFromDomain(s))
		}
	}
	return out
}
func AttackResultFromDTO(d *AttackResultDTO) *domain.AttackResult {
	if d == nil {
		return nil
	}
	out := &domain.AttackResult{Outcome: domain.AttackOutcome(d.Outcome), Loot: PriceFromDTO(d.Loot)}
	if len(d.AttackerRemaining) > 0 {
		out.AttackerRemaining = make([]domain.MilitaryUnitSnap, 0, len(d.AttackerRemaining))
		for _, u := range d.AttackerRemaining {
			out.AttackerRemaining = append(out.AttackerRemaining, MilitaryUnitFromDTO(u))
		}
	}
	if len(d.DefenderRemaining) > 0 {
		out.DefenderRemaining = make([]domain.MilitaryUnitSnap, 0, len(d.DefenderRemaining))
		for _, u := range d.DefenderRemaining {
			out.DefenderRemaining = append(out.DefenderRemaining, MilitaryUnitFromDTO(u))
		}
	}
	if len(d.RemainingStructures) > 0 {
		out.RemainingStructures = make([]domain.DefenseStructureSnap, 0, len(d.RemainingStructures))
		for _, s := range d.RemainingStructures {
			out.RemainingStructures = append(out.RemainingStructures, DefenseStructureFromDTO(s))
		}
	}
	if len(d.Trophies) > 0 {
		out.Trophies = make([]domain.TrophyStorageItem, 0, len(d.Trophies))
		for _, t := range d.Trophies {
			out.Trophies = append(out.Trophies, TrophyFromDTO(t))
		}
	}
	if len(d.DefendersBefore) > 0 {
		out.DefendersBefore = make([]domain.MilitaryUnitSnap, 0, len(d.DefendersBefore))
		for _, u := range d.DefendersBefore {
			out.DefendersBefore = append(out.DefendersBefore, MilitaryUnitFromDTO(u))
		}
	}
	if len(d.StructuresBefore) > 0 {
		out.StructuresBefore = make([]domain.DefenseStructureSnap, 0, len(d.StructuresBefore))
		for _, s := range d.StructuresBefore {
			out.StructuresBefore = append(out.StructuresBefore, DefenseStructureFromDTO(s))
		}
	}
	return out
}
