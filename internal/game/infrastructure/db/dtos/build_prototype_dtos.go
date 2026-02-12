package dtos

import "github.com/artcodefun/heat-expansion-api/internal/game/core/domain"

// ControlBuildingDataDTO represents JSON shape for control building data in prototypes.
type ControlBuildingDataDTO struct {
	Subtype string `json:"subtype"`
}

func ControlBuildingDataDTOFromDomain(d *domain.ControlBuildingData) *ControlBuildingDataDTO {
	if d == nil {
		return nil
	}
	return &ControlBuildingDataDTO{Subtype: string(d.Subtype)}
}

func ControlBuildingDataFromDTO(d *ControlBuildingDataDTO) *domain.ControlBuildingData {
	if d == nil {
		return nil
	}
	return &domain.ControlBuildingData{Subtype: domain.ControlSubtype(d.Subtype)}
}

// ResourcesBuildingDataDTO represents JSON shape for resource-building data.
type ResourcesBuildingDataDTO struct {
	CreditsProduction    float64 `json:"credits_production"`
	IronProduction       float64 `json:"iron_production"`
	TitaniumProduction   float64 `json:"titanium_production"`
	AntimatterProduction float64 `json:"antimatter_production"`
	CreditsCapacity      int     `json:"credits_capacity"`
	IronCapacity         int     `json:"iron_capacity"`
	TitaniumCapacity     int     `json:"titanium_capacity"`
	AntimatterCapacity   int     `json:"antimatter_capacity"`
}

func ResourcesBuildingDataDTOFromDomain(d *domain.ResourcesBuildingData) *ResourcesBuildingDataDTO {
	if d == nil {
		return nil
	}
	return &ResourcesBuildingDataDTO{
		CreditsProduction:    d.CreditsProduction,
		IronProduction:       d.IronProduction,
		TitaniumProduction:   d.TitaniumProduction,
		AntimatterProduction: d.AntimatterProduction,
		CreditsCapacity:      d.CreditsCapacity,
		IronCapacity:         d.IronCapacity,
		TitaniumCapacity:     d.TitaniumCapacity,
		AntimatterCapacity:   d.AntimatterCapacity,
	}
}

func ResourcesBuildingDataFromDTO(d *ResourcesBuildingDataDTO) *domain.ResourcesBuildingData {
	if d == nil {
		return nil
	}
	return &domain.ResourcesBuildingData{
		CreditsProduction:    d.CreditsProduction,
		IronProduction:       d.IronProduction,
		TitaniumProduction:   d.TitaniumProduction,
		AntimatterProduction: d.AntimatterProduction,
		CreditsCapacity:      d.CreditsCapacity,
		IronCapacity:         d.IronCapacity,
		TitaniumCapacity:     d.TitaniumCapacity,
		AntimatterCapacity:   d.AntimatterCapacity,
	}
}

// DefenseBuildingDataDTO represents JSON shape for defense-building data.
type DefenseBuildingDataDTO struct {
	DefenceBonus int `json:"defence_bonus"`
}

func DefenseBuildingDataDTOFromDomain(d *domain.DefenseBuildingData) *DefenseBuildingDataDTO {
	if d == nil {
		return nil
	}
	return &DefenseBuildingDataDTO{
		DefenceBonus: d.DefenceBonus,
	}
}

func DefenseBuildingDataFromDTO(d *DefenseBuildingDataDTO) *domain.DefenseBuildingData {
	if d == nil {
		return nil
	}
	return &domain.DefenseBuildingData{
		DefenceBonus: d.DefenceBonus,
	}
}

// MilitaryBuildingDataDTO represents JSON shape for military-building data.
type MilitaryBuildingDataDTO struct {
	UnlockArmyCategory string `json:"unlock_army_category"`
}

func MilitaryBuildingDataDTOFromDomain(d *domain.MilitaryBuildingData) *MilitaryBuildingDataDTO {
	if d == nil {
		return nil
	}
	return &MilitaryBuildingDataDTO{UnlockArmyCategory: string(d.UnlockArmyCategory)}
}

func MilitaryBuildingDataFromDTO(d *MilitaryBuildingDataDTO) *domain.MilitaryBuildingData {
	if d == nil {
		return nil
	}
	return &domain.MilitaryBuildingData{UnlockArmyCategory: domain.ArmyCategory(d.UnlockArmyCategory)}
}

// IntelligenceBuildingDataDTO represents JSON shape for intelligence-building data.
type IntelligenceBuildingDataDTO struct {
	Subtype         string `json:"subtype"`
	StealthStrength int    `json:"stealth_strength"`
	ScanRange       int    `json:"scan_range"`
	ScanCooldown    int64  `json:"scan_cooldown"`
}

func IntelligenceBuildingDataDTOFromDomain(d *domain.IntelligenceBuildingData) *IntelligenceBuildingDataDTO {
	if d == nil {
		return nil
	}
	return &IntelligenceBuildingDataDTO{
		Subtype:         string(d.Subtype),
		StealthStrength: d.StealthStrength,
		ScanRange:       d.ScanRange,
		ScanCooldown:    d.ScanCooldown,
	}
}

func IntelligenceBuildingDataFromDTO(d *IntelligenceBuildingDataDTO) *domain.IntelligenceBuildingData {
	if d == nil {
		return nil
	}
	return &domain.IntelligenceBuildingData{
		Subtype:         domain.IntelligenceSubtype(d.Subtype),
		StealthStrength: d.StealthStrength,
		ScanRange:       d.ScanRange,
		ScanCooldown:    d.ScanCooldown,
	}
}
