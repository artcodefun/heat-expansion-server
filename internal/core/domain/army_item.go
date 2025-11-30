package domain

// ArmyCategory represents the category of an army item.
type ArmyCategory string

const (
	ArmyCategoryInfantry  ArmyCategory = "INFANTRY"
	ArmyCategoryArmored   ArmyCategory = "ARMORED"
	ArmyCategoryArtillery ArmyCategory = "ARTILLERY"
	ArmyCategoryAviation  ArmyCategory = "AVIATION"
	ArmyCategorySpy       ArmyCategory = "SPY"
	ArmyCategorySpecial   ArmyCategory = "SPECIAL"
)

// ArmyStatus represents the status of an army item.
type ArmyStatus string

const (
	ArmyStatusPending      ArmyStatus = "PENDING"
	ArmyStatusInProduction ArmyStatus = "IN_PRODUCTION"
	ArmyStatusPresent      ArmyStatus = "PRESENT"
	ArmyStatusDeployed     ArmyStatus = "DEPLOYED"
)

// ArmyItemPrototype is the base struct for army item prototypes.
type ArmyItemPrototype struct {
	ID                 int
	Name               string
	Category           ArmyCategory
	UnlockTechnologyID *int // nil: available by default; non-nil: unlocked by this technology
	ShortDescription   string
	FullDescription    string
	Price              PriceModel
	ProductionTime     int64 // how many seconds it takes to create
	Space              int
	ImageURL           string
	Attack             int
	Defence            int
	Capacity           int
	Stealth            int
	Speed              int
}

// ArmyItemPending represents a pending army item with count.
type ArmyItemPending struct {
	BaseOwnedItem
	Prototype ArmyItemPrototype
	Count     int
}

// ArmyItemInProduction represents an army item in production with task details.
type ArmyItemInProduction struct {
	BaseOwnedItem
	Prototype         ArmyItemPrototype
	StartDate         int64
	CompletionDate    int64
	CrystalsSkipPrice int
}

// ArmyItemPresent represents a present army item with count and refund.
type ArmyItemPresent struct {
	BaseOwnedItem
	Prototype ArmyItemPrototype
	Count     int
	Refund    PriceModel
}

// ArmyItemDeployed represents units allocated to a military operation and currently away from the base.
type ArmyItemDeployed struct {
	BaseOwnedItem
	Prototype   ArmyItemPrototype
	OperationID int // owning military operation id
	Count       int
}
