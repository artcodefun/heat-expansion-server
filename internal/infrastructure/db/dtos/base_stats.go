package dtos

import "github.com/artcodefun/heat-expansion-api/internal/core/domain"

type BaseStatsDTO struct {
	Credits              int     `json:"credits"`
	CreditsCapacity      int     `json:"credits_capacity"`
	CreditsProduction    float64 `json:"credits_production"`
	Iron                 int     `json:"iron"`
	IronCapacity         int     `json:"iron_capacity"`
	IronProduction       float64 `json:"iron_production"`
	Titanium             int     `json:"titanium"`
	TitaniumCapacity     int     `json:"titanium_capacity"`
	TitaniumProduction   float64 `json:"titanium_production"`
	Antimatter           int     `json:"antimatter"`
	AntimatterCapacity   int     `json:"antimatter_capacity"`
	AntimatterProduction float64 `json:"antimatter_production"`
	Defence              int     `json:"defence"`
	Attack               int     `json:"attack"`
	Space                int     `json:"space"`
	SpaceCapacity        int     `json:"space_capacity"`
}

func BaseStatsDTOFromDomain(s domain.UserBaseStats) BaseStatsDTO {
	return BaseStatsDTO{
		Credits:              s.Credits,
		CreditsCapacity:      s.CreditsCapacity,
		CreditsProduction:    s.CreditsProduction,
		Iron:                 s.Iron,
		IronCapacity:         s.IronCapacity,
		IronProduction:       s.IronProduction,
		Titanium:             s.Titanium,
		TitaniumCapacity:     s.TitaniumCapacity,
		TitaniumProduction:   s.TitaniumProduction,
		Antimatter:           s.Antimatter,
		AntimatterCapacity:   s.AntimatterCapacity,
		AntimatterProduction: s.AntimatterProduction,
		Defence:              s.Defence,
		Attack:               s.Attack,
		Space:                s.Space,
		SpaceCapacity:        s.SpaceCapacity,
	}
}

func BaseStatsFromDTO(d BaseStatsDTO, calcTs int64) domain.UserBaseStats {
	return domain.UserBaseStats{
		Credits:              d.Credits,
		CreditsCapacity:      d.CreditsCapacity,
		CreditsProduction:    d.CreditsProduction,
		Iron:                 d.Iron,
		IronCapacity:         d.IronCapacity,
		IronProduction:       d.IronProduction,
		Titanium:             d.Titanium,
		TitaniumCapacity:     d.TitaniumCapacity,
		TitaniumProduction:   d.TitaniumProduction,
		Antimatter:           d.Antimatter,
		AntimatterCapacity:   d.AntimatterCapacity,
		AntimatterProduction: d.AntimatterProduction,
		Defence:              d.Defence,
		Attack:               d.Attack,
		Space:                d.Space,
		SpaceCapacity:        d.SpaceCapacity,
		CalculationTimestamp: calcTs,
	}
}
