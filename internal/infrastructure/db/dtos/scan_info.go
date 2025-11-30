package dtos

import "github.com/artcodefun/heat-expansion-api/internal/core/domain"

type ScanInfoDTO struct {
	Credits    int `json:"credits"`
	Iron       int `json:"iron"`
	Titanium   int `json:"titanium"`
	Antimatter int `json:"antimatter"`
	Defence    int `json:"defence"`
	Attack     int `json:"attack"`
	Space      int `json:"space"`
}

func ScanInfoDTOFromDomain(s domain.ScanInfo) ScanInfoDTO {
	return ScanInfoDTO{Credits: s.Credits, Iron: s.Iron, Titanium: s.Titanium, Antimatter: s.Antimatter, Defence: s.Defence, Attack: s.Attack, Space: s.Space}
}
func ScanInfoFromDTO(d ScanInfoDTO) domain.ScanInfo {
	return domain.ScanInfo{Credits: d.Credits, Iron: d.Iron, Titanium: d.Titanium, Antimatter: d.Antimatter, Defence: d.Defence, Attack: d.Attack, Space: d.Space}
}
