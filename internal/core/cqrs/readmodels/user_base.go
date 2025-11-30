package readmodels

import (
	"github.com/google/uuid"
)

// UserBaseModel represents a military base in a sector.
type UserBaseModel struct {
	ID          int
	Coordinates Vector2i
	UserID      int
	LocationDetails

	ArmiesPending      []ArmyItemPending
	ArmiesPresent      []ArmyItemPresent
	ArmiesInProduction []ArmyItemInProduction
	ArmiesDeployed     []ArmyItemDeployed // deployed units grouped by operation

	BuildingsPending      []BuildItemPending
	BuildingsPresent      []BuildItemPresent
	BuildingsInProduction []BuildItemInProduction

	TechnologiesInProgress []TechItemInProgress
	TechnologiesDone       []TechItemDone

	StorageItemsPresent []StorageItemPresent

	Stats UserBaseStats
	// ...existing fields...
}

// UserBaseStats represents current properties of a base.
type UserBaseStats struct {
	Credits              int
	CreditsCapacity      int
	CreditsProduction    float64
	Iron                 int
	IronCapacity         int
	IronProduction       float64
	Titanium             int
	TitaniumCapacity     int
	TitaniumProduction   float64
	Antimatter           int
	AntimatterCapacity   int
	AntimatterProduction float64
	Defence              int
	Attack               int
	Space                int
	SpaceCapacity        int
	CalculationTimestamp int64 // Unix timestamp of last resource calculation
}

// BaseOwnedItem is embedded in all items that belong to a user base.
type BaseOwnedItem struct {
	ID         uuid.UUID
	UserBaseID int
}
