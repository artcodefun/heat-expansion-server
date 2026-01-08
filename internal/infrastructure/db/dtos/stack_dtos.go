package dtos

import "github.com/artcodefun/heat-expansion-api/internal/core/domain"

type ArmyStackDTO struct {
	PrototypeID int `json:"prototype_id"`
	Count       int `json:"count"`
}

type DefenseStackDTO struct {
	PrototypeID int `json:"prototype_id"`
	Count       int `json:"count"`
}

func ArmyStackFromDTO(d ArmyStackDTO, p domain.ArmyItemPrototype) domain.ArmyStack {
	return domain.ArmyStack{
		Prototype: p,
		Count:     d.Count,
	}
}

func ArmyStackDTOFromDomain(s domain.ArmyStack) ArmyStackDTO {
	return ArmyStackDTO{
		PrototypeID: s.Prototype.ID,
		Count:       s.Count,
	}
}

func DefenseStackFromDTO(d DefenseStackDTO, p domain.BuildItemPrototype) domain.DefenseStack {
	return domain.DefenseStack{
		Prototype: p,
		Count:     d.Count,
	}
}

func DefenseStackDTOFromDomain(s domain.DefenseStack) DefenseStackDTO {
	return DefenseStackDTO{
		PrototypeID: s.Prototype.ID,
		Count:       s.Count,
	}
}
