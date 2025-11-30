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
	ResearchTime       int64 // how many seconds it takes to research
	ImageURL           string
	Effects            []TechnologyEffect
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
}

// EffectType represents the type of effect a technology can have.
type EffectType string

const (
	EffectTypeSpaceBonus    EffectType = "SPACE_BONUS"
	EffectTypeDefenceBonus  EffectType = "DEFENCE_BONUS"
	EffectTypeAttackBonus   EffectType = "ATTACK_BONUS"
	EffectTypeResourceBonus EffectType = "RESOURCE_BONUS"
	// Add more as needed
)

// TechnologyEffect describes an effect a technology has on a base or items.
type TechnologyEffect struct {
	EffectType EffectType
	Value      int
}
