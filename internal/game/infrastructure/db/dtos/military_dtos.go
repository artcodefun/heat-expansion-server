package dtos

import "github.com/artcodefun/heat-expansion-api/internal/game/core/domain"

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

func MilitaryUnitDTOFromDomain(u domain.MilitaryUnitSnap) MilitaryUnitDTO {
	return MilitaryUnitDTO{PrototypeID: u.PrototypeID, Category: string(u.Category), Attack: u.Attack, Defence: u.Defence, Capacity: u.Capacity, Stealth: u.Stealth, Speed: u.Speed, Count: u.Count}
}
func MilitaryUnitFromDTO(d MilitaryUnitDTO) domain.MilitaryUnitSnap {
	return domain.MilitaryUnitSnap{PrototypeID: d.PrototypeID, Category: domain.ArmyCategory(d.Category), Attack: d.Attack, Defence: d.Defence, Capacity: d.Capacity, Stealth: d.Stealth, Speed: d.Speed, Count: d.Count}
}

func DefenseStructureDTOFromDomain(s domain.DefenseStructureSnap) DefenseStructureDTO {
	return DefenseStructureDTO{PrototypeID: s.PrototypeID, Defence: s.Defence, Count: s.Count}
}
func DefenseStructureFromDTO(d DefenseStructureDTO) domain.DefenseStructureSnap {
	return domain.DefenseStructureSnap{PrototypeID: d.PrototypeID, Defence: d.Defence, Count: d.Count}
}

type MilitaryModifiersDTO struct {
	AttackMul   float64 `json:"attack_mul"`
	DefenceMul  float64 `json:"defence_mul"`
	StealthMul  float64 `json:"stealth_mul"`
	CapacityMul float64 `json:"capacity_mul"`
	SpeedMul    float64 `json:"speed_mul"`
}

func MilitaryModifiersDTOFromDomain(m domain.MilitaryModifiers) MilitaryModifiersDTO {
	return MilitaryModifiersDTO{
		AttackMul:   m.AttackMul,
		DefenceMul:  m.DefenceMul,
		StealthMul:  m.StealthMul,
		CapacityMul: m.CapacityMul,
		SpeedMul:    m.SpeedMul,
	}
}

func MilitaryModifiersFromDTO(d MilitaryModifiersDTO) domain.MilitaryModifiers {
	return domain.MilitaryModifiers{
		AttackMul:   d.AttackMul,
		DefenceMul:  d.DefenceMul,
		StealthMul:  d.StealthMul,
		CapacityMul: d.CapacityMul,
		SpeedMul:    d.SpeedMul,
	}
}

type StorageItemSnapDTO struct {
	PrototypeID  int                     `json:"prototype_id"`
	Category     string                  `json:"category"`
	BuffData     *BuffStorageDataDTO     `json:"buff_data,omitempty"`
	ArtifactData *ArtifactStorageDataDTO `json:"artifact_data,omitempty"`
}

func StorageItemSnapDTOFromDomain(s domain.StorageItemSnap) StorageItemSnapDTO {
	return StorageItemSnapDTO{
		PrototypeID:  s.PrototypeID,
		Category:     string(s.Category),
		BuffData:     BuffStorageDataDTOFromDomain(s.Buff),
		ArtifactData: ArtifactStorageDataDTOFromDomain(s.Artifact),
	}
}

func StorageItemSnapFromDTO(d StorageItemSnapDTO) domain.StorageItemSnap {
	return domain.StorageItemSnap{
		PrototypeID: d.PrototypeID,
		Category:    domain.StorageCategory(d.Category),
		Buff:        BuffStorageDataFromDTO(d.BuffData),
		Artifact:    ArtifactStorageDataFromDTO(d.ArtifactData),
	}
}

type TrophyDTO struct {
	PrototypeID int `json:"prototype_id"`
}

func TrophyDTOFromDomain(t domain.TrophyStorageItem) TrophyDTO {
	return TrophyDTO{PrototypeID: t.PrototypeID}
}

func TrophyFromDTO(d TrophyDTO) domain.TrophyStorageItem {
	return domain.TrophyStorageItem{PrototypeID: d.PrototypeID}
}
