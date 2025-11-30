package dtos

import "github.com/artcodefun/heat-expansion-api/internal/core/domain"

type MilitaryUnitDTO struct {
	PrototypeID int    `json:"prototype_id"`
	Category    string `json:"category"`
	Attack      int    `json:"attack"`
	Defence     int    `json:"defence"`
	Capacity    int    `json:"capacity"`
	Stealth     int    `json:"stealth"`
	Speed       int    `json:"speed"`
	Count       int    `json:"count"`
}

type DefenseStructureDTO struct {
	PrototypeID int `json:"prototype_id"`
	Defence     int `json:"defence"`
	Count       int `json:"count"`
}

func MilitaryUnitDTOFromDomain(u domain.MilitaryUnit) MilitaryUnitDTO {
	return MilitaryUnitDTO{PrototypeID: u.PrototypeID, Category: string(u.Category), Attack: u.Attack, Defence: u.Defence, Capacity: u.Capacity, Stealth: u.Stealth, Speed: u.Speed, Count: u.Count}
}
func MilitaryUnitFromDTO(d MilitaryUnitDTO) domain.MilitaryUnit {
	return domain.MilitaryUnit{PrototypeID: d.PrototypeID, Category: domain.ArmyCategory(d.Category), Attack: d.Attack, Defence: d.Defence, Capacity: d.Capacity, Stealth: d.Stealth, Speed: d.Speed, Count: d.Count}
}

func DefenseStructureDTOFromDomain(s domain.DefenseStructure) DefenseStructureDTO {
	return DefenseStructureDTO{PrototypeID: s.PrototypeID, Defence: s.Defence, Count: s.Count}
}
func DefenseStructureFromDTO(d DefenseStructureDTO) domain.DefenseStructure {
	return domain.DefenseStructure{PrototypeID: d.PrototypeID, Defence: d.Defence, Count: d.Count}
}
