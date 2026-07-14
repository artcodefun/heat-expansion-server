package dtos

import "github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs/readmodels"

// GetPrototypeURI binds the :id parameter for prototype endpoints.
type GetPrototypeURI struct {
	ID int64 `uri:"id" binding:"required"`
}

// PriceModelDTO is the JSON representation of a resource price.
type PriceModelDTO struct {
	Credits    int64 `json:"credits"`
	Iron       int64 `json:"iron"`
	Titanium   int64 `json:"titanium"`
	Antimatter int64 `json:"antimatter"`
}

func priceToDTO(p readmodels.PriceModel) PriceModelDTO {
	return PriceModelDTO{Credits: p.Credits, Iron: p.Iron, Titanium: p.Titanium, Antimatter: p.Antimatter}
}

func priceFromDTO(d PriceModelDTO) readmodels.PriceModel {
	return readmodels.PriceModel{Credits: d.Credits, Iron: d.Iron, Titanium: d.Titanium, Antimatter: d.Antimatter}
}

// ── Army ──────────────────────────────────────────────────────────────────────

type ArmyPrototypeDTO struct {
	ID                 int64         `json:"id"`
	Name               string        `json:"name"               binding:"required"`
	Category           string        `json:"category"           binding:"required"`
	CreationSources    []string      `json:"creation_sources"`
	Faction            string        `json:"faction"            binding:"required"`
	UnlockTechnologyID *int64        `json:"unlock_technology_id,omitempty"`
	ShortDescription   string        `json:"short_description"`
	FullDescription    string        `json:"full_description"`
	Price              PriceModelDTO `json:"price"`
	ProductionTime     int64         `json:"production_time"`
	Space              int32         `json:"space"`
	ImageURL           string        `json:"image_url"`
	Attack             int32         `json:"attack"`
	Defence            int32         `json:"defence"`
	Capacity           int32         `json:"capacity"`
	Stealth            int32         `json:"stealth"`
	Speed              int32         `json:"speed"`
}

func ArmyPrototypeDTOFromModel(m *readmodels.ArmyPrototype) ArmyPrototypeDTO {
	return ArmyPrototypeDTO{
		ID:                 m.ID,
		Name:               m.Name,
		Category:           m.Category,
		CreationSources:    m.CreationSources,
		Faction:            m.Faction,
		UnlockTechnologyID: m.UnlockTechnologyID,
		ShortDescription:   m.ShortDescription,
		FullDescription:    m.FullDescription,
		Price:              priceToDTO(m.Price),
		ProductionTime:     m.ProductionTime,
		Space:              m.Space,
		ImageURL:           m.ImageURL,
		Attack:             m.Attack,
		Defence:            m.Defence,
		Capacity:           m.Capacity,
		Stealth:            m.Stealth,
		Speed:              m.Speed,
	}
}

func ArmyPrototypeDTOToModel(d ArmyPrototypeDTO) *readmodels.ArmyPrototype {
	return &readmodels.ArmyPrototype{
		ID:                 d.ID,
		Name:               d.Name,
		Category:           d.Category,
		CreationSources:    d.CreationSources,
		Faction:            d.Faction,
		UnlockTechnologyID: d.UnlockTechnologyID,
		ShortDescription:   d.ShortDescription,
		FullDescription:    d.FullDescription,
		Price:              priceFromDTO(d.Price),
		ProductionTime:     d.ProductionTime,
		Space:              d.Space,
		ImageURL:           d.ImageURL,
		Attack:             d.Attack,
		Defence:            d.Defence,
		Capacity:           d.Capacity,
		Stealth:            d.Stealth,
		Speed:              d.Speed,
	}
}

// ── Build ──────────────────────────────────────────────────────────────────────

type BuildControlDataDTO struct {
	Subtype string `json:"subtype"`
}

type BuildResourcesDataDTO struct {
	CreditsProduction    float64 `json:"credits_production"`
	IronProduction       float64 `json:"iron_production"`
	TitaniumProduction   float64 `json:"titanium_production"`
	AntimatterProduction float64 `json:"antimatter_production"`
	CreditsCapacity      int32   `json:"credits_capacity"`
	IronCapacity         int32   `json:"iron_capacity"`
	TitaniumCapacity     int32   `json:"titanium_capacity"`
	AntimatterCapacity   int32   `json:"antimatter_capacity"`
}

type BuildDefenseDataDTO struct {
	DefenceBonus int32 `json:"defence_bonus"`
}

type BuildMilitaryDataDTO struct {
	UnlockArmyCategory string `json:"unlock_army_category"`
}

type BuildIntelligenceDataDTO struct {
	Subtype         string `json:"subtype"`
	StealthStrength int32  `json:"stealth_strength"`
	ScanRange       int32  `json:"scan_range"`
	ScanCooldown    int64  `json:"scan_cooldown"`
}

type BuildPrototypeDTO struct {
	ID                 int64                     `json:"id"`
	Name               string                    `json:"name"               binding:"required"`
	Category           string                    `json:"category"           binding:"required"`
	CreationSources    []string                  `json:"creation_sources"`
	Faction            string                    `json:"faction"            binding:"required"`
	UnlockTechnologyID *int64                    `json:"unlock_technology_id,omitempty"`
	ShortDescription   string                    `json:"short_description"`
	FullDescription    string                    `json:"full_description"`
	Price              PriceModelDTO             `json:"price"`
	ProductionTime     int64                     `json:"production_time"`
	Space              int32                     `json:"space"`
	ImageURL           string                    `json:"image_url"`
	ControlData        *BuildControlDataDTO      `json:"control_data,omitempty"`
	ResourcesData      *BuildResourcesDataDTO    `json:"resources_data,omitempty"`
	DefenseData        *BuildDefenseDataDTO      `json:"defense_data,omitempty"`
	MilitaryData       *BuildMilitaryDataDTO     `json:"military_data,omitempty"`
	IntelligenceData   *BuildIntelligenceDataDTO `json:"intelligence_data,omitempty"`
}

func BuildPrototypeDTOFromModel(m *readmodels.BuildPrototype) BuildPrototypeDTO {
	d := BuildPrototypeDTO{
		ID:                 m.ID,
		Name:               m.Name,
		Category:           m.Category,
		CreationSources:    m.CreationSources,
		Faction:            m.Faction,
		UnlockTechnologyID: m.UnlockTechnologyID,
		ShortDescription:   m.ShortDescription,
		FullDescription:    m.FullDescription,
		Price:              priceToDTO(m.Price),
		ProductionTime:     m.ProductionTime,
		Space:              m.Space,
		ImageURL:           m.ImageURL,
	}
	if m.ControlData != nil {
		d.ControlData = &BuildControlDataDTO{Subtype: m.ControlData.Subtype}
	}
	if m.ResourcesData != nil {
		r := m.ResourcesData
		d.ResourcesData = &BuildResourcesDataDTO{
			CreditsProduction: r.CreditsProduction, IronProduction: r.IronProduction,
			TitaniumProduction: r.TitaniumProduction, AntimatterProduction: r.AntimatterProduction,
			CreditsCapacity: r.CreditsCapacity, IronCapacity: r.IronCapacity,
			TitaniumCapacity: r.TitaniumCapacity, AntimatterCapacity: r.AntimatterCapacity,
		}
	}
	if m.DefenseData != nil {
		d.DefenseData = &BuildDefenseDataDTO{DefenceBonus: m.DefenseData.DefenceBonus}
	}
	if m.MilitaryData != nil {
		d.MilitaryData = &BuildMilitaryDataDTO{UnlockArmyCategory: m.MilitaryData.UnlockArmyCategory}
	}
	if m.IntelligenceData != nil {
		i := m.IntelligenceData
		d.IntelligenceData = &BuildIntelligenceDataDTO{Subtype: i.Subtype, StealthStrength: i.StealthStrength, ScanRange: i.ScanRange, ScanCooldown: i.ScanCooldown}
	}
	return d
}

func BuildPrototypeDTOToModel(d BuildPrototypeDTO) *readmodels.BuildPrototype {
	m := &readmodels.BuildPrototype{
		ID:                 d.ID,
		Name:               d.Name,
		Category:           d.Category,
		CreationSources:    d.CreationSources,
		Faction:            d.Faction,
		UnlockTechnologyID: d.UnlockTechnologyID,
		ShortDescription:   d.ShortDescription,
		FullDescription:    d.FullDescription,
		Price:              priceFromDTO(d.Price),
		ProductionTime:     d.ProductionTime,
		Space:              d.Space,
		ImageURL:           d.ImageURL,
	}
	if d.ControlData != nil {
		m.ControlData = &readmodels.BuildControlData{Subtype: d.ControlData.Subtype}
	}
	if d.ResourcesData != nil {
		r := d.ResourcesData
		m.ResourcesData = &readmodels.BuildResourcesData{
			CreditsProduction: r.CreditsProduction, IronProduction: r.IronProduction,
			TitaniumProduction: r.TitaniumProduction, AntimatterProduction: r.AntimatterProduction,
			CreditsCapacity: r.CreditsCapacity, IronCapacity: r.IronCapacity,
			TitaniumCapacity: r.TitaniumCapacity, AntimatterCapacity: r.AntimatterCapacity,
		}
	}
	if d.DefenseData != nil {
		m.DefenseData = &readmodels.BuildDefenseData{DefenceBonus: d.DefenseData.DefenceBonus}
	}
	if d.MilitaryData != nil {
		m.MilitaryData = &readmodels.BuildMilitaryData{UnlockArmyCategory: d.MilitaryData.UnlockArmyCategory}
	}
	if d.IntelligenceData != nil {
		i := d.IntelligenceData
		m.IntelligenceData = &readmodels.BuildIntelligenceData{Subtype: i.Subtype, StealthStrength: i.StealthStrength, ScanRange: i.ScanRange, ScanCooldown: i.ScanCooldown}
	}
	return m
}

// ── Storage ───────────────────────────────────────────────────────────────────

type StorageBuffDataDTO struct {
	Type            string  `json:"type"`
	Value           float32 `json:"value"`
	DurationSeconds int64   `json:"duration_seconds"`
}

type StorageIntelDataDTO struct {
	Type              string `json:"type"`
	DecryptionSeconds int64  `json:"decryption_seconds"`
}

type StorageDamagedDataDTO struct {
	RestorePrice       PriceModelDTO `json:"restore_price"`
	RestorationSeconds int64         `json:"restoration_seconds"`
	OriginalUnitID     int64         `json:"original_unit_id"`
}

type StorageArtifactDataDTO struct {
	Type  string  `json:"type"`
	Value float32 `json:"value"`
}

type StorageConsumableDataDTO struct {
	Type        string   `json:"type"`
	BoxContents []string `json:"box_contents"`
	BoxSize     int32    `json:"box_size"`
}

type StoragePrototypeDTO struct {
	ID               int64                     `json:"id"`
	Name             string                    `json:"name"     binding:"required"`
	Category         string                    `json:"category" binding:"required"`
	CreationSources  []string                  `json:"creation_sources"`
	EstimatedWorth   int32                     `json:"estimated_worth"`
	ShortDescription string                    `json:"short_description"`
	FullDescription  string                    `json:"full_description"`
	ImageURL         string                    `json:"image_url"`
	BuffData         *StorageBuffDataDTO       `json:"buff_data,omitempty"`
	IntelData        *StorageIntelDataDTO      `json:"intel_data,omitempty"`
	DamagedData      *StorageDamagedDataDTO    `json:"damaged_data,omitempty"`
	ArtifactData     *StorageArtifactDataDTO   `json:"artifact_data,omitempty"`
	ConsumableData   *StorageConsumableDataDTO `json:"consumable_data,omitempty"`
}

func StoragePrototypeDTOFromModel(m *readmodels.StoragePrototype) StoragePrototypeDTO {
	d := StoragePrototypeDTO{
		ID: m.ID, Name: m.Name, Category: m.Category, CreationSources: m.CreationSources,
		EstimatedWorth: m.EstimatedWorth, ShortDescription: m.ShortDescription,
		FullDescription: m.FullDescription, ImageURL: m.ImageURL,
	}
	if m.BuffData != nil {
		d.BuffData = &StorageBuffDataDTO{Type: m.BuffData.Type, Value: m.BuffData.Value, DurationSeconds: m.BuffData.DurationSeconds}
	}
	if m.IntelData != nil {
		d.IntelData = &StorageIntelDataDTO{Type: m.IntelData.Type, DecryptionSeconds: m.IntelData.DecryptionSeconds}
	}
	if m.DamagedData != nil {
		d.DamagedData = &StorageDamagedDataDTO{RestorePrice: priceToDTO(m.DamagedData.RestorePrice), RestorationSeconds: m.DamagedData.RestorationSeconds, OriginalUnitID: m.DamagedData.OriginalUnitID}
	}
	if m.ArtifactData != nil {
		d.ArtifactData = &StorageArtifactDataDTO{Type: m.ArtifactData.Type, Value: m.ArtifactData.Value}
	}
	if m.ConsumableData != nil {
		d.ConsumableData = &StorageConsumableDataDTO{Type: m.ConsumableData.Type, BoxContents: m.ConsumableData.BoxContents, BoxSize: m.ConsumableData.BoxSize}
	}
	return d
}

func StoragePrototypeDTOToModel(d StoragePrototypeDTO) *readmodels.StoragePrototype {
	m := &readmodels.StoragePrototype{
		ID: d.ID, Name: d.Name, Category: d.Category, CreationSources: d.CreationSources,
		EstimatedWorth: d.EstimatedWorth, ShortDescription: d.ShortDescription,
		FullDescription: d.FullDescription, ImageURL: d.ImageURL,
	}
	if d.BuffData != nil {
		m.BuffData = &readmodels.StorageBuffData{Type: d.BuffData.Type, Value: d.BuffData.Value, DurationSeconds: d.BuffData.DurationSeconds}
	}
	if d.IntelData != nil {
		m.IntelData = &readmodels.StorageIntelData{Type: d.IntelData.Type, DecryptionSeconds: d.IntelData.DecryptionSeconds}
	}
	if d.DamagedData != nil {
		m.DamagedData = &readmodels.StorageDamagedData{RestorePrice: priceFromDTO(d.DamagedData.RestorePrice), RestorationSeconds: d.DamagedData.RestorationSeconds, OriginalUnitID: d.DamagedData.OriginalUnitID}
	}
	if d.ArtifactData != nil {
		m.ArtifactData = &readmodels.StorageArtifactData{Type: d.ArtifactData.Type, Value: d.ArtifactData.Value}
	}
	if d.ConsumableData != nil {
		m.ConsumableData = &readmodels.StorageConsumableData{Type: d.ConsumableData.Type, BoxContents: d.ConsumableData.BoxContents, BoxSize: d.ConsumableData.BoxSize}
	}
	return m
}

// ── Tech ──────────────────────────────────────────────────────────────────────

type TechImprovementDTO struct {
	Type     string `json:"type"`
	Value    int32  `json:"value"`
	MaxLevel int32  `json:"max_level"`
}

type TechPrototypeDTO struct {
	ID                 int64               `json:"id"`
	Name               string              `json:"name"     binding:"required"`
	Category           string              `json:"category" binding:"required"`
	UnlockTechnologyID int64               `json:"unlock_technology_id"`
	ShortDescription   string              `json:"short_description"`
	FullDescription    string              `json:"full_description"`
	Price              PriceModelDTO       `json:"price"`
	ResearchTime       int64               `json:"research_time"`
	ImageURL           string              `json:"image_url"`
	Improvement        *TechImprovementDTO `json:"improvement,omitempty"`
}

func TechPrototypeDTOFromModel(m *readmodels.TechPrototype) TechPrototypeDTO {
	d := TechPrototypeDTO{
		ID: m.ID, Name: m.Name, Category: m.Category, UnlockTechnologyID: m.UnlockTechnologyID,
		ShortDescription: m.ShortDescription, FullDescription: m.FullDescription,
		Price: priceToDTO(m.Price), ResearchTime: m.ResearchTime, ImageURL: m.ImageURL,
	}
	if m.Improvement != nil {
		d.Improvement = &TechImprovementDTO{Type: m.Improvement.Type, Value: m.Improvement.Value, MaxLevel: m.Improvement.MaxLevel}
	}
	return d
}

func TechPrototypeDTOToModel(d TechPrototypeDTO) *readmodels.TechPrototype {
	m := &readmodels.TechPrototype{
		ID: d.ID, Name: d.Name, Category: d.Category, UnlockTechnologyID: d.UnlockTechnologyID,
		ShortDescription: d.ShortDescription, FullDescription: d.FullDescription,
		Price: priceFromDTO(d.Price), ResearchTime: d.ResearchTime, ImageURL: d.ImageURL,
	}
	if d.Improvement != nil {
		m.Improvement = &readmodels.TechImprovement{Type: d.Improvement.Type, Value: d.Improvement.Value, MaxLevel: d.Improvement.MaxLevel}
	}
	return m
}
