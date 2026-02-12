package dtos

import "github.com/artcodefun/heat-expansion-api/internal/game/core/domain"

type LocationResourceStatsDTO struct {
	Credits    int `json:"credits"`
	Iron       int `json:"iron"`
	Titanium   int `json:"titanium"`
	Antimatter int `json:"antimatter"`
}

func LocationResourceStatsDTOFromDomain(s domain.LocationResourceStats) LocationResourceStatsDTO {
	return LocationResourceStatsDTO{Credits: s.Credits, Iron: s.Iron, Titanium: s.Titanium, Antimatter: s.Antimatter}
}
func LocationResourceStatsFromDTO(d LocationResourceStatsDTO, calcTs int64) domain.LocationResourceStats {
	return domain.LocationResourceStats{Credits: d.Credits, Iron: d.Iron, Titanium: d.Titanium, Antimatter: d.Antimatter, CalculationTimestamp: calcTs}
}
