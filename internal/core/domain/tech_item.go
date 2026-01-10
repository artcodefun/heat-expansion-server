package domain

// TechCategory represents the category of a tech item.
type TechCategory string

const (
	TechCategoryArmy     TechCategory = "ARMY"
	TechCategoryBuild    TechCategory = "BUILD"
	TechCategoryBase     TechCategory = "BASE"
	TechCategoryPolitics TechCategory = "POLITICS"
)

// TechStatus represents the status of a tech item.
type TechStatus string

const (
	TechStatusInProgress TechStatus = "IN_PROGRESS"
	TechStatusDone       TechStatus = "DONE"
)

// TechItemPrototype is the base struct for tech item prototypes.
type TechItemPrototype struct {
	ID                 int
	Name               string
	Category           TechCategory
	UnlockTechnologyID *int // nil: available by default; non-nil: unlocked by this technology
	ShortDescription   string
	FullDescription    string
	Price              PriceModel
	ResearchTime       int64            // how many seconds it takes to research
	Improvement        *TechImprovement // optional numeric improvement offered by this technology
	ImageURL           string
}

// TechItemInProgress represents a tech item in progress with task details.
type TechItemInProgress struct {
	BaseOwnedItem
	Prototype         TechItemPrototype
	StartDate         int64
	CompletionDate    int64
	CrystalsSkipPrice int
}

// TechItemDone represents a completed tech item.
type TechItemDone struct {
	BaseOwnedItem
	Prototype    TechItemPrototype
	ResearchedAt int64 // Unix timestamp when research was completed
	Level        int   // current level of the technology (how many times it was researched)
}

// ImprovementType represents the type of improvement a technology can provide.
type ImprovementType string

const (
	ImprovementTypeSpaceCapacity           ImprovementType = "SPACE_CAPACITY"
	ImprovementTypeOperationsCount         ImprovementType = "OPERATIONS_COUNT"
	ImprovementTypeActiveBuffsCount        ImprovementType = "ACTIVE_BUFFS_COUNT"
	ImprovementTypeActiveArtifactsCount    ImprovementType = "ACTIVE_ARTIFACTS_COUNT"
	ImprovementTypeActiveRestorationsCount ImprovementType = "ACTIVE_RESTORATIONS_COUNT"
	ImprovementTypeBuildingProductionCount ImprovementType = "BUILDING_PRODUCTION_COUNT"
)

// TechImprovement describes a single numeric benefit that scales with technology level.
type TechImprovement struct {
	Type     ImprovementType
	Value    int
	MaxLevel *int // if nil, this improvement (and thus the tech) can be upgraded infinitely
}
