package readmodels

// BuildCategory represents the category of a build item.
type BuildCategory string

const (
	BuildCategoryControl      BuildCategory = "CONTROL"
	BuildCategoryResources    BuildCategory = "RESOURCES"
	BuildCategoryDefense      BuildCategory = "DEFENSE"
	BuildCategoryMilitary     BuildCategory = "MILITARY"
	BuildCategoryIntelligence BuildCategory = "INTELLIGENCE"
)

// BuildStatus represents the status of a build item.
type BuildStatus string

const (
	BuildStatusPending      BuildStatus = "PENDING"
	BuildStatusInProduction BuildStatus = "IN_PRODUCTION"
	BuildStatusPresent      BuildStatus = "PRESENT"
)

// BuildItemPrototype is the base struct for build item prototypes.
type BuildItemPrototype struct {
	ID                 int
	Name               string
	Category           BuildCategory
	UnlockTechnologyID *int // nil: available by default; non-nil: unlocked by this technology
	ShortDescription   string
	FullDescription    string
	Price              PriceModel
	ProductionTime     int64 // how many seconds it takes to create
	Space              int
	ImageURL           string

	// Category-specific fields
	ControlData      *ControlBuildingData
	ResourcesData    *ResourcesBuildingData
	DefenseData      *DefenseBuildingData
	MilitaryData     *MilitaryBuildingData
	IntelligenceData *IntelligenceBuildingData
}

// Control buildings: e.g., repair center, traiding terminal
type ControlBuildingData struct {
	Subtype ControlSubtype
}

type ControlSubtype string

const (
	ControlSubtypeRepairCenter     ControlSubtype = "REPAIR_CENTER"
	ControlSubtypeTraidingTerminal ControlSubtype = "TRADING_TERMINAL"
	ControlSubtypeMailingTerminal  ControlSubtype = "MAILING_TERMINAL"
	// Add more as needed
)

// Resources buildings: e.g., resource production/capacity
type ResourcesBuildingData struct {
	CreditsProduction    float64
	IronProduction       float64
	TitaniumProduction   float64
	AntimatterProduction float64
	CreditsCapacity      int
	IronCapacity         int
	TitaniumCapacity     int
	AntimatterCapacity   int
}

// Defense buildings: e.g., defence bonus, shield strength
type DefenseBuildingData struct {
	DefenceBonus   int
	ShieldStrength int
	// Add more as needed
}

// Military buildings: e.g., unlock units, training speed
type MilitaryBuildingData struct {
	UnlockArmyCategory ArmyCategory
}

// Intelligence buildings: e.g., vision range, stealth bonus
type IntelligenceBuildingData struct {
	Subtype            IntelligenceSubtype
	StealthStrength    int
	TargetLocationType LocationType
	ScanRange          int
	// ScanCooldown defines how many seconds to wait between automatic scans
	ScanCooldown int64
	// Add more as needed
}

type IntelligenceSubtype string

const (
	IntelligenceSubtypeRadar    IntelligenceSubtype = "RADAR"
	IntelligenceSubtypeCloaking IntelligenceSubtype = "CLOAKING"
	// Add more as needed
)

// BuildItemPending represents a pending build item.
type BuildItemPending struct {
	BaseOwnedItem
	Prototype BuildItemPrototype
}

// BuildItemInProduction represents a build item in production with task details.
type BuildItemInProduction struct {
	BaseOwnedItem
	Prototype         BuildItemPrototype
	StartDate         int64
	CompletionDate    int64
	CrystalsSkipPrice int
}

// BuildItemPresent represents a present build item with building ID and refund.
type BuildItemPresent struct {
	BaseOwnedItem
	Prototype BuildItemPrototype
	Refund    PriceModel
}

// BuildItemNew represents a build item prototype available to be queued (no ownership yet).
type BuildItemNew struct {
	Prototype BuildItemPrototype
}
