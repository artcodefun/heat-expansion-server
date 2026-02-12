package dtos

import "github.com/artcodefun/heat-expansion-api/internal/game/core/domain"

type BaseStatsDTO struct {
	Credits               float64 `json:"credits"`
	CreditsCapacity       int     `json:"credits_capacity"`
	CreditsProduction     float64 `json:"credits_production"`
	Iron                  float64 `json:"iron"`
	IronCapacity          int     `json:"iron_capacity"`
	IronProduction        float64 `json:"iron_production"`
	Titanium              float64 `json:"titanium"`
	TitaniumCapacity      int     `json:"titanium_capacity"`
	TitaniumProduction    float64 `json:"titanium_production"`
	Antimatter            float64 `json:"antimatter"`
	AntimatterCapacity    int     `json:"antimatter_capacity"`
	AntimatterProduction  float64 `json:"antimatter_production"`
	Defence               int     `json:"defence"`
	Attack                int     `json:"attack"`
	Space                 int     `json:"space"`
	MaxSpace              int     `json:"max_space"`
	MaxOperations         int     `json:"max_operations"`
	MaxActiveBuffs        int     `json:"max_active_buffs"`
	MaxActiveArtifacts    int     `json:"max_active_artifacts"`
	MaxBuildingProduction int     `json:"max_building_production"`
	MaxActiveRestorations int     `json:"max_active_restorations"`
	MaxActiveDecryptions  int     `json:"max_active_decryptions"`
}

func BaseStatsDTOFromDomain(s domain.UserBaseStats) BaseStatsDTO {
	return BaseStatsDTO{
		Credits:               s.Credits,
		CreditsCapacity:       s.CreditsCapacity,
		CreditsProduction:     s.CreditsProduction,
		Iron:                  s.Iron,
		IronCapacity:          s.IronCapacity,
		IronProduction:        s.IronProduction,
		Titanium:              s.Titanium,
		TitaniumCapacity:      s.TitaniumCapacity,
		TitaniumProduction:    s.TitaniumProduction,
		Antimatter:            s.Antimatter,
		AntimatterCapacity:    s.AntimatterCapacity,
		AntimatterProduction:  s.AntimatterProduction,
		Defence:               s.Defence,
		Attack:                s.Attack,
		Space:                 s.Space,
		MaxSpace:              s.MaxSpace,
		MaxOperations:         s.MaxOperations,
		MaxActiveBuffs:        s.MaxActiveBuffs,
		MaxActiveArtifacts:    s.MaxActiveArtifacts,
		MaxBuildingProduction: s.MaxBuildingProduction,
		MaxActiveRestorations: s.MaxActiveRestorations,
		MaxActiveDecryptions:  s.MaxActiveDecryptions,
	}
}

func BaseStatsFromDTO(d BaseStatsDTO, calcTs int64) domain.UserBaseStats {
	return domain.UserBaseStats{
		Credits:               d.Credits,
		CreditsCapacity:       d.CreditsCapacity,
		CreditsProduction:     d.CreditsProduction,
		Iron:                  d.Iron,
		IronCapacity:          d.IronCapacity,
		IronProduction:        d.IronProduction,
		Titanium:              d.Titanium,
		TitaniumCapacity:      d.TitaniumCapacity,
		TitaniumProduction:    d.TitaniumProduction,
		Antimatter:            d.Antimatter,
		AntimatterCapacity:    d.AntimatterCapacity,
		AntimatterProduction:  d.AntimatterProduction,
		Defence:               d.Defence,
		Attack:                d.Attack,
		Space:                 d.Space,
		MaxSpace:              d.MaxSpace,
		MaxOperations:         d.MaxOperations,
		MaxActiveBuffs:        d.MaxActiveBuffs,
		MaxActiveArtifacts:    d.MaxActiveArtifacts,
		MaxBuildingProduction: d.MaxBuildingProduction,
		MaxActiveRestorations: d.MaxActiveRestorations,
		MaxActiveDecryptions:  d.MaxActiveDecryptions,
		CalculationTimestamp:  calcTs,
	}
}
