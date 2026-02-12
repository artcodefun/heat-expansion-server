package dtos

import (
	"strings"

	"github.com/artcodefun/heat-expansion-api/internal/game/core/cqrs/readmodels"
)

type BuildCategory string

type BuildStatus string

// BuildCategory enum values
const (
	Control      BuildCategory = "CONTROL"
	Resources    BuildCategory = "RESOURCES"
	Defense      BuildCategory = "DEFENSE"
	Military     BuildCategory = "MILITARY"
	Intelligence BuildCategory = "INTELLIGENCE"
)

type ControlSubtype string

const (
	ControlSubtypeRepairCenter    ControlSubtype = "REPAIR_CENTER"
	ControlSubtypeCryptographyLab ControlSubtype = "CRYPTOGRAPHY_LAB"
	ControlSubtypeArtifactLab     ControlSubtype = "ARTIFACT_LAB"
	ControlSubtypeTradingTerminal ControlSubtype = "TRADING_TERMINAL"
	ControlSubtypeMailingTerminal ControlSubtype = "MAILING_TERMINAL"
)

type IntelligenceSubtype string

const (
	IntelligenceSubtypeScanner         IntelligenceSubtype = "SCANNER"
	IntelligenceSubtypeRadar           IntelligenceSubtype = "RADAR"
	IntelligenceSubtypeCloaking        IntelligenceSubtype = "CLOAKING"
	IntelligenceSubtypeScanInterceptor IntelligenceSubtype = "SCAN_INTERCEPTOR"
)

// BuildStatus enum values
const (
	BuildNew          BuildStatus = "NEW"
	BuildPending      BuildStatus = "PENDING"
	BuildInProduction BuildStatus = "IN_PRODUCTION"
	BuildPresent      BuildStatus = "PRESENT"
)

type BuildItemPrototypeDTO struct {
	ID               int           `json:"id"`
	Name             string        `json:"name"`
	Category         BuildCategory `json:"category"`
	Faction          Faction       `json:"faction"`
	ShortDescription string        `json:"short_description"`
	FullDescription  string        `json:"full_description"`
	Price            PriceModelDTO `json:"price"`
	Space            int           `json:"space"`
	ImageURL         string        `json:"image_url"`
	ProductionTime   int           `json:"production_time"`

	// Category-specific fields
	ControlData      *ControlBuildingDataDTO      `json:"control_data,omitempty"`
	ResourcesData    *ResourcesBuildingDataDTO    `json:"resources_data,omitempty"`
	DefenseData      *DefenseBuildingDataDTO      `json:"defense_data,omitempty"`
	MilitaryData     *MilitaryBuildingDataDTO     `json:"military_data,omitempty"`
	IntelligenceData *IntelligenceBuildingDataDTO `json:"intelligence_data,omitempty"`
}

type ControlBuildingDataDTO struct {
	Subtype ControlSubtype `json:"subtype"`
}

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

type DefenseBuildingDataDTO struct {
	DefenceBonus int `json:"defence_bonus"`
}

type MilitaryBuildingDataDTO struct {
	UnlockArmyCategory ArmyCategory `json:"unlock_army_category"`
}

type IntelligenceBuildingDataDTO struct {
	Subtype         IntelligenceSubtype `json:"subtype"`
	StealthStrength int                 `json:"stealth_strength"`
	ScanRange       int                 `json:"scan_range"`
	ScanCooldown    int                 `json:"scan_cooldown"`
}

type BuildItemNewDTO struct {
	BuildItemPrototypeDTO
}

type BuildItemPendingDTO struct {
	BaseOwnedItemDTO
	Prototype BuildItemPrototypeDTO `json:"prototype"`
}

type BuildItemInProductionDTO struct {
	BaseOwnedItemDTO
	Prototype         BuildItemPrototypeDTO `json:"prototype"`
	StartDate         int                   `json:"start_date"`
	CompletionDate    int                   `json:"completion_date"`
	CrystalsSkipPrice int                   `json:"crystals_skip_price"`
}

type BuildItemPresentDTO struct {
	BaseOwnedItemDTO
	Prototype BuildItemPrototypeDTO `json:"prototype"`
	Refund    PriceModelDTO         `json:"refund"`
}

func mapBuildItemPrototype(proto readmodels.BuildItemPrototype) BuildItemPrototypeDTO {
	dto := BuildItemPrototypeDTO{
		ID:               proto.ID,
		Name:             proto.Name,
		Category:         BuildCategory(proto.Category),
		Faction:          Faction(proto.Faction),
		ShortDescription: proto.ShortDescription,
		FullDescription:  proto.FullDescription,
		Price:            PriceModelFromReadModel(proto.Price),
		Space:            proto.Space,
		ImageURL:         proto.ImageURL,
		ProductionTime:   int(proto.ProductionTime),
	}

	if proto.ControlData != nil {
		dto.ControlData = &ControlBuildingDataDTO{
			Subtype: ControlSubtype(proto.ControlData.Subtype),
		}
	}
	if proto.ResourcesData != nil {
		dto.ResourcesData = &ResourcesBuildingDataDTO{
			CreditsProduction:    proto.ResourcesData.CreditsProduction,
			IronProduction:       proto.ResourcesData.IronProduction,
			TitaniumProduction:   proto.ResourcesData.TitaniumProduction,
			AntimatterProduction: proto.ResourcesData.AntimatterProduction,
			CreditsCapacity:      proto.ResourcesData.CreditsCapacity,
			IronCapacity:         proto.ResourcesData.IronCapacity,
			TitaniumCapacity:     proto.ResourcesData.TitaniumCapacity,
			AntimatterCapacity:   proto.ResourcesData.AntimatterCapacity,
		}
	}
	if proto.DefenseData != nil {
		dto.DefenseData = &DefenseBuildingDataDTO{
			DefenceBonus: proto.DefenseData.DefenceBonus,
		}
	}
	if proto.MilitaryData != nil {
		dto.MilitaryData = &MilitaryBuildingDataDTO{
			UnlockArmyCategory: ArmyCategory(proto.MilitaryData.UnlockArmyCategory),
		}
	}
	if proto.IntelligenceData != nil {
		dto.IntelligenceData = &IntelligenceBuildingDataDTO{
			Subtype:         IntelligenceSubtype(proto.IntelligenceData.Subtype),
			StealthStrength: proto.IntelligenceData.StealthStrength,
			ScanRange:       proto.IntelligenceData.ScanRange,
			ScanCooldown:    int(proto.IntelligenceData.ScanCooldown),
		}
	}

	return dto
}

func BuildItemsNewFromReadModels(items []*readmodels.BuildItemNew) []BuildItemNewDTO {
	out := make([]BuildItemNewDTO, 0, len(items))
	for _, item := range items {
		out = append(out, BuildItemNewDTO{BuildItemPrototypeDTO: mapBuildItemPrototype(item.Prototype)})
	}
	return out
}

func BuildItemsPendingFromReadModels(items []*readmodels.BuildItemPending) []BuildItemPendingDTO {
	out := make([]BuildItemPendingDTO, 0, len(items))
	for _, item := range items {
		out = append(out, BuildItemPendingDTO{
			BaseOwnedItemDTO: BaseOwnedItemDTOFromReadModel(item.BaseOwnedItem),
			Prototype:        mapBuildItemPrototype(item.Prototype),
		})
	}
	return out
}

func BuildItemsInProductionFromReadModels(items []*readmodels.BuildItemInProduction) []BuildItemInProductionDTO {
	out := make([]BuildItemInProductionDTO, 0, len(items))
	for _, item := range items {
		out = append(out, BuildItemInProductionDTO{
			BaseOwnedItemDTO:  BaseOwnedItemDTOFromReadModel(item.BaseOwnedItem),
			Prototype:         mapBuildItemPrototype(item.Prototype),
			StartDate:         int(item.StartDate),
			CompletionDate:    int(item.CompletionDate),
			CrystalsSkipPrice: item.CrystalsSkipPrice,
		})
	}
	return out
}

func BuildItemsPresentFromReadModels(items []*readmodels.BuildItemPresent) []BuildItemPresentDTO {
	out := make([]BuildItemPresentDTO, 0, len(items))
	for _, item := range items {
		out = append(out, BuildItemPresentDTO{
			BaseOwnedItemDTO: BaseOwnedItemDTOFromReadModel(item.BaseOwnedItem),
			Prototype:        mapBuildItemPrototype(item.Prototype),
			Refund:           PriceModelFromReadModel(item.Refund),
		})
	}
	return out
}

// BuildCategoryFromDTO normalizes a DTO category string to the read-model type.
func BuildCategoryFromDTO(value string) readmodels.BuildCategory {
	return readmodels.BuildCategory(strings.ToUpper(strings.TrimSpace(value)))
}
