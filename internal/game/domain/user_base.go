package domain

import (
	"github.com/google/uuid"
)

// UserBaseModel represents a military base in a sector.
type UserBaseModel struct {
	EventProducer
	ID          int
	Coordinates Vector2i
	UserID      uuid.UUID
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
}

// NewUserBaseModel constructs a fresh user base aggregate, ensuring all slices are initialized
// and stats reflect the current state (even when empty).
func NewUserBaseModel(baseID int, userID uuid.UUID, coords Vector2i) *UserBaseModel {
	ub := &UserBaseModel{
		ID:                     baseID,
		Coordinates:            coords,
		UserID:                 userID,
		ArmiesPending:          []ArmyItemPending{},
		ArmiesPresent:          []ArmyItemPresent{},
		ArmiesInProduction:     []ArmyItemInProduction{},
		ArmiesDeployed:         []ArmyItemDeployed{},
		BuildingsPending:       []BuildItemPending{},
		BuildingsPresent:       []BuildItemPresent{},
		BuildingsInProduction:  []BuildItemInProduction{},
		TechnologiesInProgress: []TechItemInProgress{},
		TechnologiesDone:       []TechItemDone{},
		StorageItemsPresent:    []StorageItemPresent{},
	}
	ub.recalculateStats()
	return ub
}

// EmitCreated records a domain event indicating this base has been created.
func (ub *UserBaseModel) EmitCreated() {
	if ub == nil || ub.ID <= 0 {
		return
	}
	ub.AddEvent(NewUserBaseCreatedEvent(ub.ID, ub.UserID))
}

// BaseOwnedItem is embedded in all items that belong to a user base.
type BaseOwnedItem struct {
	ID         uuid.UUID
	UserBaseID int
}

func NewBaseOwnedItem(baseId int) BaseOwnedItem {
	return BaseOwnedItem{
		ID:         uuid.Must(uuid.NewV7()),
		UserBaseID: baseId,
	}
}
