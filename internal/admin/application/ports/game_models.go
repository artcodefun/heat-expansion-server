package ports

// PriceModel is the resource price used by game prototypes.
type PriceModel struct {
	Credits    int64 `json:"credits"`
	Iron       int64 `json:"iron"`
	Titanium   int64 `json:"titanium"`
	Antimatter int64 `json:"antimatter"`
}

// ArmyPrototype is the admin model for an army item prototype.
type ArmyPrototype struct {
	ID                 int64      `json:"id"`
	Name               string     `json:"name"`
	Category           string     `json:"category"`
	CreationSources    []string   `json:"creation_sources"`
	Faction            string     `json:"faction"`
	UnlockTechnologyID *int64     `json:"unlock_technology_id,omitempty"`
	ShortDescription   string     `json:"short_description"`
	FullDescription    string     `json:"full_description"`
	Price              PriceModel `json:"price"`
	ProductionTime     int64      `json:"production_time"`
	Space              int32      `json:"space"`
	ImageURL           string     `json:"image_url"`
	Attack             int32      `json:"attack"`
	Defence            int32      `json:"defence"`
	Capacity           int32      `json:"capacity"`
	Stealth            int32      `json:"stealth"`
	Speed              int32      `json:"speed"`
}

// BuildControlData carries category data for control buildings.
type BuildControlData struct {
	Subtype string `json:"subtype"`
}

// BuildResourcesData carries category data for resource buildings.
type BuildResourcesData struct {
	CreditsProduction    float64 `json:"credits_production"`
	IronProduction       float64 `json:"iron_production"`
	TitaniumProduction   float64 `json:"titanium_production"`
	AntimatterProduction float64 `json:"antimatter_production"`
	CreditsCapacity      int32   `json:"credits_capacity"`
	IronCapacity         int32   `json:"iron_capacity"`
	TitaniumCapacity     int32   `json:"titanium_capacity"`
	AntimatterCapacity   int32   `json:"antimatter_capacity"`
}

// BuildDefenseData carries category data for defense buildings.
type BuildDefenseData struct {
	DefenceBonus int32 `json:"defence_bonus"`
}

// BuildMilitaryData carries category data for military buildings.
type BuildMilitaryData struct {
	UnlockArmyCategory string `json:"unlock_army_category"`
}

// BuildIntelligenceData carries category data for intelligence buildings.
type BuildIntelligenceData struct {
	Subtype         string `json:"subtype"`
	StealthStrength int32  `json:"stealth_strength"`
	ScanRange       int32  `json:"scan_range"`
	ScanCooldown    int64  `json:"scan_cooldown"`
}

// BuildPrototype is the admin model for a build item prototype.
// Exactly one of the category data fields is non-nil.
type BuildPrototype struct {
	ID                 int64      `json:"id"`
	Name               string     `json:"name"`
	Category           string     `json:"category"`
	CreationSources    []string   `json:"creation_sources"`
	Faction            string     `json:"faction"`
	UnlockTechnologyID *int64     `json:"unlock_technology_id,omitempty"`
	ShortDescription   string     `json:"short_description"`
	FullDescription    string     `json:"full_description"`
	Price              PriceModel `json:"price"`
	ProductionTime     int64      `json:"production_time"`
	Space              int32      `json:"space"`
	ImageURL           string     `json:"image_url"`

	ControlData      *BuildControlData      `json:"control_data,omitempty"`
	ResourcesData    *BuildResourcesData    `json:"resources_data,omitempty"`
	DefenseData      *BuildDefenseData      `json:"defense_data,omitempty"`
	MilitaryData     *BuildMilitaryData     `json:"military_data,omitempty"`
	IntelligenceData *BuildIntelligenceData `json:"intelligence_data,omitempty"`
}

// StorageBuffData carries category data for buff storage items.
type StorageBuffData struct {
	Type            string  `json:"type"`
	Value           float32 `json:"value"`
	DurationSeconds int64   `json:"duration_seconds"`
}

// StorageIntelData carries category data for intel storage items.
type StorageIntelData struct {
	Type              string `json:"type"`
	DecryptionSeconds int64  `json:"decryption_seconds"`
}

// StorageDamagedData carries category data for damaged-unit storage items.
type StorageDamagedData struct {
	RestorePrice       PriceModel `json:"restore_price"`
	RestorationSeconds int64      `json:"restoration_seconds"`
	OriginalUnitID     int64      `json:"original_unit_id"`
}

// StorageArtifactData carries category data for artifact storage items.
type StorageArtifactData struct {
	Type  string  `json:"type"`
	Value float32 `json:"value"`
}

// StorageConsumableData carries category data for consumable storage items.
type StorageConsumableData struct {
	Type        string   `json:"type"`
	BoxContents []string `json:"box_contents"`
	BoxSize     int32    `json:"box_size"`
}

// StoragePrototype is the admin model for a storage item prototype.
// Exactly one of the category data fields is non-nil.
type StoragePrototype struct {
	ID               int64    `json:"id"`
	Name             string   `json:"name"`
	Category         string   `json:"category"`
	CreationSources  []string `json:"creation_sources"`
	EstimatedWorth   int32    `json:"estimated_worth"`
	ShortDescription string   `json:"short_description"`
	FullDescription  string   `json:"full_description"`
	ImageURL         string   `json:"image_url"`

	BuffData       *StorageBuffData       `json:"buff_data,omitempty"`
	IntelData      *StorageIntelData      `json:"intel_data,omitempty"`
	DamagedData    *StorageDamagedData    `json:"damaged_data,omitempty"`
	ArtifactData   *StorageArtifactData   `json:"artifact_data,omitempty"`
	ConsumableData *StorageConsumableData `json:"consumable_data,omitempty"`
}

// TechImprovement is the optional numeric improvement offered by a technology.
type TechImprovement struct {
	Type     string `json:"type"`
	Value    int32  `json:"value"`
	MaxLevel int32  `json:"max_level"` // 0 = no cap
}

// TechPrototype is the admin model for a tech item prototype.
type TechPrototype struct {
	ID                 int64            `json:"id"`
	Name               string           `json:"name"`
	Category           string           `json:"category"`
	UnlockTechnologyID int64            `json:"unlock_technology_id"` // 0 = no prerequisite
	ShortDescription   string           `json:"short_description"`
	FullDescription    string           `json:"full_description"`
	Price              PriceModel       `json:"price"`
	ResearchTime       int64            `json:"research_time"`
	ImageURL           string           `json:"image_url"`
	Improvement        *TechImprovement `json:"improvement,omitempty"`
}

// Translation is the admin model for a single localised string entry.
type Translation struct {
	Key    string `json:"key"`
	Locale string `json:"locale"`
	Value  string `json:"value"`
}
