package domain

import (
	"fmt"

	"github.com/google/uuid"
)

// UserBaseModel represents a military base in a sector.
type UserBaseModel struct {
	EventProducer
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
}

// NewUserBaseModel constructs a fresh user base aggregate, ensuring all slices are initialized
// and stats reflect the current state (even when empty).
func NewUserBaseModel(baseID int, userID int, coords Vector2i) *UserBaseModel {
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

// Helper for checking unlocked tech
func hasTech(techs []TechItemDone, techID int) bool {
	for _, t := range techs {
		if t.Prototype.ID == techID {
			return true
		}
	}
	return false
}

// Domain logic for building creation

// Returns all building prototypes the user can create based on unlocked technologies
func (ub *UserBaseModel) AvailableBuildings(allPrototypes []*BuildItemPrototype) []*BuildItemPrototype {
	available := []*BuildItemPrototype{}
	for _, proto := range allPrototypes {
		if proto.UnlockTechnologyID == nil || hasTech(ub.TechnologiesDone, *proto.UnlockTechnologyID) {
			available = append(available, proto)
		}
	}
	return available
}

// GetProductionCompletionTime returns the completion time for a building in production by item ID.
func (ub *UserBaseModel) GetProductionCompletionTime(id uuid.UUID) (int64, bool) {
	for _, prod := range ub.BuildingsInProduction {
		if prod.ID == id {
			return prod.CompletionDate, true
		}
	}
	return 0, false
}

// Queues a new building for production
func (ub *UserBaseModel) AddToBuildQueue(proto *BuildItemPrototype) error {
	defer ub.recalculateStats()

	// Ensure this prototype is actually available for this base
	if len(ub.AvailableBuildings([]*BuildItemPrototype{proto})) == 0 {
		return fmt.Errorf("this building is not available for production")
	}

	// Calculate total space after adding this building
	totalSpace := ub.Stats.Space + proto.Space
	if totalSpace > ub.Stats.SpaceCapacity {
		return fmt.Errorf("not enough space to queue building: required %d, available %d", totalSpace, ub.Stats.SpaceCapacity)
	}

	// Validate resources (example: credits, iron, titanium, antimatter)
	if proto.Price.Credits > ub.Stats.Credits {
		return fmt.Errorf("not enough credits")
	}
	if proto.Price.Iron > ub.Stats.Iron {
		return fmt.Errorf("not enough iron")
	}
	if proto.Price.Titanium > ub.Stats.Titanium {
		return fmt.Errorf("not enough titanium")
	}
	if proto.Price.Antimatter > ub.Stats.Antimatter {
		return fmt.Errorf("not enough antimatter")
	}
	// Subtract price from resources
	ub.Stats.Credits -= proto.Price.Credits
	ub.Stats.Iron -= proto.Price.Iron
	ub.Stats.Titanium -= proto.Price.Titanium
	ub.Stats.Antimatter -= proto.Price.Antimatter

	// Always add to pending
	pendingItem := BuildItemPending{
		BaseOwnedItem: NewBaseOwnedItem(ub.ID),
		Prototype:     *proto,
	}
	ub.BuildingsPending = append(ub.BuildingsPending, pendingItem)
	// Optionally emit event for building added to pending
	// ub.AddEvent(NewBuildingProductionPendingEvent(ub.ID, proto.ID))

	// Immediately process the queue
	ub.MoveBuildQueue()

	return nil
}

// Moves finished production items to present and starts next pending item
func (ub *UserBaseModel) MoveBuildQueue() {
	defer ub.recalculateStats()

	now := NowUnix()
	var remainingInProduction []BuildItemInProduction
	for _, prod := range ub.BuildingsInProduction {
		if prod.CompletionDate <= now {
			// Move to present
			present := BuildItemPresent{
				BaseOwnedItem: NewBaseOwnedItem(ub.ID),
				Prototype:     prod.Prototype,
				Refund:        prod.Prototype.Price.Divide(10),
			}
			ub.BuildingsPresent = append(ub.BuildingsPresent, present)
			// Emit event for building production finished
			ub.AddEvent(NewBuildingProductionFinishedEvent(ub.ID, prod.ID, present.ID))
		} else {
			remainingInProduction = append(remainingInProduction, prod)
		}
	}
	ub.BuildingsInProduction = remainingInProduction

	// If no items in production, start next pending
	if len(ub.BuildingsInProduction) == 0 && len(ub.BuildingsPending) > 0 {
		next := ub.BuildingsPending[0]
		ub.BuildingsPending = ub.BuildingsPending[1:]
		startDate := now
		completionDate := startDate + next.Prototype.ProductionTime
		crystalsSkipPrice := int(next.Prototype.ProductionTime / 60)
		newProd := BuildItemInProduction{
			BaseOwnedItem:     NewBaseOwnedItem(ub.ID),
			Prototype:         next.Prototype,
			StartDate:         startDate,
			CompletionDate:    completionDate,
			CrystalsSkipPrice: crystalsSkipPrice,
		}
		ub.BuildingsInProduction = append(ub.BuildingsInProduction, newProd)

		// Emit event for building production started
		ub.AddEvent(NewBuildingProductionStartedEvent(ub.ID, newProd.ID, completionDate))
	}
}

// CancelPendingBuildingByID removes a pending building by item ID and refunds its price.
func (ub *UserBaseModel) CancelPendingBuildingByID(itemID uuid.UUID) error {
	defer ub.recalculateStats()

	idx := -1
	for i, item := range ub.BuildingsPending {
		if item.ID == itemID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("pending building with ID %s not found", itemID)
	}
	item := ub.BuildingsPending[idx]
	// Refund resources
	ub.Stats.Credits += item.Prototype.Price.Credits
	ub.Stats.Iron += item.Prototype.Price.Iron
	ub.Stats.Titanium += item.Prototype.Price.Titanium
	ub.Stats.Antimatter += item.Prototype.Price.Antimatter
	// Remove from pending
	ub.BuildingsPending = append(ub.BuildingsPending[:idx], ub.BuildingsPending[idx+1:]...)
	// Optionally emit event for cancellation
	ub.AddEvent(NewBuildingProductionCancelledEvent(ub.ID, item.ID))
	return nil
}

// SpeedUpBuildingProduction finishes building production immediately for the given item ID.
func (ub *UserBaseModel) SpeedUpBuildingProduction(buildingItemID uuid.UUID) error {
	idx := -1
	for i, item := range ub.BuildingsInProduction {
		if item.BaseOwnedItem.ID == buildingItemID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("in-production building with ID %s not found", buildingItemID)
	}
	// Set completion date to now
	ub.BuildingsInProduction[idx].CompletionDate = NowUnix()
	// Capture the in-production item ID before moving the queue (it may be removed)
	spedUpItemID := ub.BuildingsInProduction[idx].ID
	ub.MoveBuildQueue()
	ub.AddEvent(NewBuildingProductionSpeedupEvent(ub.ID, spedUpItemID))

	return nil
}

// DeletePresentBuildingByID removes a present building by item ID and emits an event.
func (ub *UserBaseModel) DeletePresentBuildingByID(itemID uuid.UUID) error {
	idx := -1
	var item BuildItemPresent
	for i, b := range ub.BuildingsPresent {
		if b.ID == itemID {
			idx = i
			item = b
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("present building with ID %s not found", itemID)
	}
	// Refund resources to base from item's Refund field
	ub.Stats.Credits += item.Refund.Credits
	ub.Stats.Iron += item.Refund.Iron
	ub.Stats.Titanium += item.Refund.Titanium
	ub.Stats.Antimatter += item.Refund.Antimatter
	// Remove from present
	ub.BuildingsPresent = append(ub.BuildingsPresent[:idx], ub.BuildingsPresent[idx+1:]...)
	// Emit event for deletion
	ub.AddEvent(NewBuildingPresentDeletedEvent(ub.ID, itemID))
	ub.recalculateStats()
	return nil
}

// Helper to count present military buildings by army category
func (ub *UserBaseModel) countMilitaryBuildingsForCategory(category ArmyCategory) int {
	count := 0
	for _, b := range ub.BuildingsPresent {
		if b.Prototype.MilitaryData != nil && ArmyCategory(b.Prototype.MilitaryData.UnlockArmyCategory) == category {
			count++
		}
	}
	return count
}

func (ub *UserBaseModel) hasControlSubtype(subtype ControlSubtype) bool {
	for _, b := range ub.BuildingsPresent {
		if b.Prototype.Category == BuildCategoryControl &&
			b.Prototype.ControlData != nil &&
			b.Prototype.ControlData.Subtype == subtype {
			return true
		}
	}
	return false
}

// Returns all army prototypes the user can create based on unlocked technologies and present military buildings
func (ub *UserBaseModel) AvailableArmies(allPrototypes []*ArmyItemPrototype) []*ArmyItemPrototype {
	available := []*ArmyItemPrototype{}
	for _, proto := range allPrototypes {
		// Check tech unlock
		if proto.UnlockTechnologyID != nil && !hasTech(ub.TechnologiesDone, *proto.UnlockTechnologyID) {
			continue
		}
		// Check for present military building of matching category
		if ub.countMilitaryBuildingsForCategory(proto.Category) == 0 {
			continue
		}
		available = append(available, proto)
	}
	return available
}

// Queues a new army item for production (batch with count)
func (ub *UserBaseModel) QueueArmy(proto *ArmyItemPrototype, count int) error {
	defer ub.recalculateStats()

	if count < 1 {
		return fmt.Errorf("count must be at least 1")
	}

	// Ensure this prototype is actually available for this base
	if len(ub.AvailableArmies([]*ArmyItemPrototype{proto})) == 0 {
		return fmt.Errorf("this army item is not available for production")
	}

	// Validate available space (armies in queue and production should reserve space
	// just like buildings do).
	requiredSpace := proto.Space * count
	totalSpace := ub.Stats.Space + requiredSpace
	if totalSpace > ub.Stats.SpaceCapacity {
		return fmt.Errorf("not enough space to queue army: required %d, available %d", totalSpace, ub.Stats.SpaceCapacity)
	}

	// Validate resources
	totalPrice := PriceModel{
		Credits:    proto.Price.Credits * count,
		Iron:       proto.Price.Iron * count,
		Titanium:   proto.Price.Titanium * count,
		Antimatter: proto.Price.Antimatter * count,
	}
	if totalPrice.Credits > ub.Stats.Credits {
		return fmt.Errorf("not enough credits")
	}
	if totalPrice.Iron > ub.Stats.Iron {
		return fmt.Errorf("not enough iron")
	}
	if totalPrice.Titanium > ub.Stats.Titanium {
		return fmt.Errorf("not enough titanium")
	}
	if totalPrice.Antimatter > ub.Stats.Antimatter {
		return fmt.Errorf("not enough antimatter")
	}
	// Subtract price
	ub.Stats.Credits -= totalPrice.Credits
	ub.Stats.Iron -= totalPrice.Iron
	ub.Stats.Titanium -= totalPrice.Titanium
	ub.Stats.Antimatter -= totalPrice.Antimatter
	// Add to pending (merge with existing batch if same prototype)
	found := false
	for i, p := range ub.ArmiesPending {
		if p.Prototype.ID == proto.ID {
			ub.ArmiesPending[i].Count += count
			found = true
			break
		}
	}
	if !found {
		pending := ArmyItemPending{
			BaseOwnedItem: NewBaseOwnedItem(ub.ID),
			Prototype:     *proto,
			Count:         count,
		}
		ub.ArmiesPending = append(ub.ArmiesPending, pending)
	}
	ub.MoveArmyQueue()
	return nil
}

// Moves finished army items to present and starts next pending batch
func (ub *UserBaseModel) MoveArmyQueue() {
	defer ub.recalculateStats()

	now := NowUnix()
	var remainingInProduction []ArmyItemInProduction
	for _, prod := range ub.ArmiesInProduction {
		if prod.CompletionDate <= now {
			// Move to present: increment count if already present
			found := false
			for i, ap := range ub.ArmiesPresent {
				if ap.Prototype.ID == prod.Prototype.ID {
					ub.ArmiesPresent[i].Count++
					found = true
					break
				}
			}
			if !found {
				present := ArmyItemPresent{
					BaseOwnedItem: NewBaseOwnedItem(ub.ID),
					Prototype:     prod.Prototype,
					Count:         1,
					Refund:        prod.Prototype.Price.Divide(10),
				}
				ub.ArmiesPresent = append(ub.ArmiesPresent, present)
			}
			// Record event for army production finished
			ub.AddEvent(NewArmyProductionFinishedEvent(ub.ID, prod.ID))
		} else {
			remainingInProduction = append(remainingInProduction, prod)
		}
	}
	ub.ArmiesInProduction = remainingInProduction

	// New logic: fill available slots for each category
	categorySlots := map[ArmyCategory]int{}
	for _, b := range ub.BuildingsPresent {
		if b.Prototype.MilitaryData != nil {
			cat := ArmyCategory(b.Prototype.MilitaryData.UnlockArmyCategory)
			categorySlots[cat]++
		}
	}
	inProductionCount := map[ArmyCategory]int{}
	for _, prod := range ub.ArmiesInProduction {
		inProductionCount[prod.Prototype.Category]++
	}
	var newPending []ArmyItemPending
	for _, pending := range ub.ArmiesPending {
		cat := pending.Prototype.Category
		slots := categorySlots[cat]
		inProd := inProductionCount[cat]
		if inProd < slots {
			startDate := now
			completionDate := startDate + pending.Prototype.ProductionTime
			crystalsSkipPrice := int(pending.Prototype.ProductionTime / 60)
			newProd := ArmyItemInProduction{
				BaseOwnedItem:     NewBaseOwnedItem(ub.ID),
				Prototype:         pending.Prototype,
				StartDate:         startDate,
				CompletionDate:    completionDate,
				CrystalsSkipPrice: crystalsSkipPrice,
			}
			ub.ArmiesInProduction = append(ub.ArmiesInProduction, newProd)
			inProductionCount[cat]++
			if pending.Count > 1 {
				pending.Count--
				newPending = append(newPending, pending)
			}
			// Record event for army production started
			ub.AddEvent(NewArmyProductionStartedEvent(ub.ID, pending.ID, completionDate))
		} else {
			newPending = append(newPending, pending)
		}
	}
	ub.ArmiesPending = newPending
}

// CancelPendingArmyByID removes a pending army item by item ID and refunds its price for a given count.
func (ub *UserBaseModel) CancelPendingArmyByID(itemID uuid.UUID, count int) error {
	defer ub.recalculateStats()

	idx := -1
	for i, item := range ub.ArmiesPending {
		if item.ID == itemID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("pending army item with ID %s not found", itemID)
	}
	item := ub.ArmiesPending[idx]
	if count < 1 || count > item.Count {
		return fmt.Errorf("invalid cancel count: %d (pending count: %d)", count, item.Count)
	}
	// Refund resources for the canceled amount
	ub.Stats.Credits += item.Prototype.Price.Credits * count
	ub.Stats.Iron += item.Prototype.Price.Iron * count
	ub.Stats.Titanium += item.Prototype.Price.Titanium * count
	ub.Stats.Antimatter += item.Prototype.Price.Antimatter * count
	if count == item.Count {
		// Remove from pending
		ub.ArmiesPending = append(ub.ArmiesPending[:idx], ub.ArmiesPending[idx+1:]...)
	} else {
		// Decrement batch count
		ub.ArmiesPending[idx].Count -= count
	}
	// Record event for army production cancellation
	ub.AddEvent(NewArmyProductionCancelledEvent(ub.ID, item.ID, count))
	return nil
}

// SpeedUpArmyProduction finishes army production immediately for the given item ID.
func (ub *UserBaseModel) SpeedUpArmyProduction(armyItemID uuid.UUID) error {
	idx := -1
	for i, item := range ub.ArmiesInProduction {
		if item.BaseOwnedItem.ID == armyItemID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("in-production army item with ID %s not found", armyItemID)
	}
	// Set completion date to now
	ub.ArmiesInProduction[idx].CompletionDate = NowUnix()
	// Capture the in-production item ID before moving the queue (it may be removed)
	spedUpItemID := ub.ArmiesInProduction[idx].ID
	ub.MoveArmyQueue()
	// Record event for army production speedup
	ub.AddEvent(NewArmyProductionSpeedupEvent(ub.ID, spedUpItemID))
	return nil
}

// DeletePresentArmyByID removes a present army item by item ID and refunds resources for a given count.
func (ub *UserBaseModel) DeletePresentArmyByID(itemID uuid.UUID, count int) error {
	idx := -1
	var item ArmyItemPresent
	for i, a := range ub.ArmiesPresent {
		if a.ID == itemID {
			idx = i
			item = a
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("present army item with ID %s not found", itemID)
	}
	if count < 1 || count > item.Count {
		return fmt.Errorf("invalid delete count: %d (present count: %d)", count, item.Count)
	}
	// Refund resources for the deleted amount
	ub.Stats.Credits += item.Refund.Credits * count
	ub.Stats.Iron += item.Refund.Iron * count
	ub.Stats.Titanium += item.Refund.Titanium * count
	ub.Stats.Antimatter += item.Refund.Antimatter * count
	if count == item.Count {
		// Remove from present
		ub.ArmiesPresent = append(ub.ArmiesPresent[:idx], ub.ArmiesPresent[idx+1:]...)
	} else {
		// Decrement batch count
		ub.ArmiesPresent[idx].Count -= count
	}
	// Record event for army present deletion
	ub.AddEvent(NewArmyPresentDeletedEvent(ub.ID, itemID, count))
	ub.recalculateStats()
	return nil
}

// Returns all technology prototypes the user can research based on unlocks and prerequisites
func (ub *UserBaseModel) AvailableTechnologies(allPrototypes []*TechItemPrototype) []*TechItemPrototype {
	available := []*TechItemPrototype{}
	for _, proto := range allPrototypes {
		// Already researched?
		if hasTech(ub.TechnologiesDone, proto.ID) {
			continue
		}
		// Already in progress?
		alreadyInProgress := false
		for _, t := range ub.TechnologiesInProgress {
			if t.Prototype.ID == proto.ID {
				alreadyInProgress = true
				break
			}
		}
		if alreadyInProgress {
			continue
		}
		// Check unlock condition (if any)
		if proto.UnlockTechnologyID != nil && !hasTech(ub.TechnologiesDone, *proto.UnlockTechnologyID) {
			continue
		}
		available = append(available, proto)
	}
	return available
}

// StartTechResearch queues a technology for research
func (ub *UserBaseModel) StartTechResearch(proto *TechItemPrototype) error {
	defer ub.recalculateStats()
	// Ensure this prototype is actually available for this base
	if len(ub.AvailableTechnologies([]*TechItemPrototype{proto})) == 0 {
		return fmt.Errorf("this technology is not available for research")
	}

	// Validate resources
	if proto.Price.Credits > ub.Stats.Credits {
		return fmt.Errorf("not enough credits")
	}
	if proto.Price.Iron > ub.Stats.Iron {
		return fmt.Errorf("not enough iron")
	}
	if proto.Price.Titanium > ub.Stats.Titanium {
		return fmt.Errorf("not enough titanium")
	}
	if proto.Price.Antimatter > ub.Stats.Antimatter {
		return fmt.Errorf("not enough antimatter")
	}
	// Subtract price
	ub.Stats.Credits -= proto.Price.Credits
	ub.Stats.Iron -= proto.Price.Iron
	ub.Stats.Titanium -= proto.Price.Titanium
	ub.Stats.Antimatter -= proto.Price.Antimatter
	// Add to in-progress
	now := NowUnix()
	completionDate := now + proto.ResearchTime
	crystalsSkipPrice := int(proto.ResearchTime / 60)
	inProgress := TechItemInProgress{
		BaseOwnedItem:     NewBaseOwnedItem(ub.ID),
		Prototype:         *proto,
		StartDate:         now,
		CompletionDate:    completionDate,
		CrystalsSkipPrice: crystalsSkipPrice,
	}
	ub.TechnologiesInProgress = append(ub.TechnologiesInProgress, inProgress)
	// Emit event for tech research started
	ub.AddEvent(NewTechResearchStartedEvent(ub.ID, inProgress.BaseOwnedItem.ID, proto.ID, completionDate))
	return nil
}

// MoveTechQueue moves finished techs to done and starts next in-progress (if any)
func (ub *UserBaseModel) MoveTechQueue() {
	defer ub.recalculateStats()
	now := NowUnix()
	var remainingInProgress []TechItemInProgress
	for _, tech := range ub.TechnologiesInProgress {
		if tech.CompletionDate <= now {
			// Move to done
			done := TechItemDone{
				BaseOwnedItem: NewBaseOwnedItem(ub.ID),
				Prototype:     tech.Prototype,
				ResearchedAt:  tech.CompletionDate,
			}
			ub.TechnologiesDone = append(ub.TechnologiesDone, done)
			// Emit event for tech research finished
			ub.AddEvent(NewTechResearchFinishedEvent(ub.ID, tech.BaseOwnedItem.ID, tech.Prototype.ID))
		} else {
			remainingInProgress = append(remainingInProgress, tech)
		}
	}
	ub.TechnologiesInProgress = remainingInProgress
}

// SpeedUpTechResearch finishes tech research immediately for the given item ID
func (ub *UserBaseModel) SpeedUpTechResearch(techItemID uuid.UUID) error {
	idx := -1
	for i, item := range ub.TechnologiesInProgress {
		if item.BaseOwnedItem.ID == techItemID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("in-progress tech with ID %s not found", techItemID)
	}
	// Set completion date to now
	ub.TechnologiesInProgress[idx].CompletionDate = NowUnix()
	// Capture IDs before moving the queue (the entry may be removed)
	spedUpItemID := ub.TechnologiesInProgress[idx].BaseOwnedItem.ID
	spedUpProtoID := ub.TechnologiesInProgress[idx].Prototype.ID
	ub.MoveTechQueue()
	// Emit event for tech research speedup
	ub.AddEvent(NewTechResearchSpeedupEvent(ub.ID, spedUpItemID, spedUpProtoID))
	return nil
}

// DeletePresentStorageItemByID removes a present storage item by item ID.
func (ub *UserBaseModel) DeletePresentStorageItemByID(itemID uuid.UUID) error {
	idx := -1
	for i, s := range ub.StorageItemsPresent {
		if s.ID == itemID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("present storage item with ID %s not found", itemID)
	}

	// Remove from present
	ub.StorageItemsPresent = append(ub.StorageItemsPresent[:idx], ub.StorageItemsPresent[idx+1:]...)
	// Emit event for deletion
	ub.AddEvent(NewStorageItemPresentDeletedEvent(ub.ID, itemID))
	ub.recalculateStats()
	return nil
}

// AddStorageItem adds a new storage item to the base.
func (ub *UserBaseModel) AddStorageItem(proto StorageItemPrototype, expiresAt *int64) uuid.UUID {
	item := StorageItemPresent{
		BaseOwnedItem: NewBaseOwnedItem(ub.ID),
		Prototype:     proto,
		ExpiresAt:     expiresAt,
		IsActive:      false,
	}
	ub.StorageItemsPresent = append(ub.StorageItemsPresent, item)
	return item.ID
}

// ActivateBuffByID activates a buff storage item by item ID, sets ExpiresAt, emits event, and returns error if not found or already activated
func (ub *UserBaseModel) ActivateBuffByID(itemID uuid.UUID) error {
	defer ub.recalculateStats()
	now := NowUnix()
	for i, item := range ub.StorageItemsPresent {
		if item.ID == itemID && item.Prototype.BuffData != nil {
			// Check if already activated
			if item.ExpiresAt != nil {
				return fmt.Errorf("buff item with ID %s is already activated", itemID)
			}
			// Set ExpiresAt
			expiresAt := now + item.Prototype.BuffData.DurationSeconds
			ub.StorageItemsPresent[i].ExpiresAt = &expiresAt
			ub.StorageItemsPresent[i].IsActive = true

			// Emit event for buff activation
			ub.AddEvent(NewBuffActivatedEvent(ub.ID, itemID))
			return nil
		}
	}
	return fmt.Errorf("buff storage item with ID %s not found or not a buff", itemID)
}

// DeleteExpiredBuffs removes expired buffs from storage.
func (ub *UserBaseModel) DeleteExpiredBuffs() int {
	now := NowUnix()
	var remaining []StorageItemPresent
	processed := 0
	for _, item := range ub.StorageItemsPresent {
		if item.ExpiresAt != nil && *item.ExpiresAt <= now && item.Prototype.BuffData != nil {
			processed++
			continue
		}
		remaining = append(remaining, item)
	}
	ub.StorageItemsPresent = remaining
	if processed > 0 {
		ub.recalculateStats()
	}
	return processed
}

// DecryptIntelItemByID completes the decryption process for a specific intel item.
func (ub *UserBaseModel) DecryptIntelItemByID(itemID uuid.UUID) (HiddenLocationType, error) {
	if !ub.hasControlSubtype(ControlSubtypeCryptographyLab) {
		return "", fmt.Errorf("cryptography lab required to decrypt intel")
	}
	now := NowUnix()
	idx := -1
	for i, item := range ub.StorageItemsPresent {
		if item.ID == itemID && item.Prototype.IntelData != nil {
			if item.ExpiresAt == nil || *item.ExpiresAt > now {
				return "", fmt.Errorf("intel item %s is not ready for decryption completion", itemID)
			}
			idx = i
			break
		}
	}
	if idx == -1 {
		return "", fmt.Errorf("ready intel item %s not found", itemID)
	}

	item := ub.StorageItemsPresent[idx]
	intelType := item.Prototype.IntelData.Type
	// Emit event
	ub.AddEvent(NewIntelDecryptionFinishedEvent(ub.ID, item.ID, intelType))

	// Remove from storage
	ub.StorageItemsPresent = append(ub.StorageItemsPresent[:idx], ub.StorageItemsPresent[idx+1:]...)
	return intelType, nil
}

// RestoreDamagedItemByID completes the restoration process for a specific damaged item.
func (ub *UserBaseModel) RestoreDamagedItemByID(itemID uuid.UUID, armyProtos []*ArmyItemPrototype) error {
	if !ub.hasControlSubtype(ControlSubtypeRepairCenter) {
		return fmt.Errorf("repair center required to restore damaged units")
	}
	defer ub.recalculateStats()
	now := NowUnix()
	idx := -1
	for i, item := range ub.StorageItemsPresent {
		if item.ID == itemID && item.Prototype.DamagedData != nil {
			if item.ExpiresAt == nil || *item.ExpiresAt > now {
				return fmt.Errorf("damaged item %s is not ready for restoration completion", itemID)
			}
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("ready damaged item %s not found", itemID)
	}

	item := ub.StorageItemsPresent[idx]
	data := item.Prototype.DamagedData
	var unitProto *ArmyItemPrototype
	for _, p := range armyProtos {
		if p.ID == data.OriginalUnitID {
			unitProto = p
			break
		}
	}

	if unitProto == nil {
		return fmt.Errorf("original unit prototype %d not found", data.OriginalUnitID)
	}

	// Add to present armies
	found := false
	for j, p := range ub.ArmiesPresent {
		if p.Prototype.ID == unitProto.ID {
			ub.ArmiesPresent[j].Count++
			found = true
			break
		}
	}
	if !found {
		ub.ArmiesPresent = append(ub.ArmiesPresent, ArmyItemPresent{
			BaseOwnedItem: NewBaseOwnedItem(ub.ID),
			Prototype:     *unitProto,
			Count:         1,
			Refund:        unitProto.Price.Divide(2),
		})
	}

	// Emit event
	ub.AddEvent(NewDamagedItemRestoredEvent(ub.ID, item.ID))

	// Remove from storage
	ub.StorageItemsPresent = append(ub.StorageItemsPresent[:idx], ub.StorageItemsPresent[idx+1:]...)
	return nil
}

// StartIntelDecryptionByID starts the decryption process for an intel item.
func (ub *UserBaseModel) StartIntelDecryptionByID(itemID uuid.UUID) error {
	if !ub.hasControlSubtype(ControlSubtypeCryptographyLab) {
		return fmt.Errorf("cryptography lab required to start intel decryption")
	}
	now := NowUnix()
	for i, item := range ub.StorageItemsPresent {
		if item.ID == itemID && item.Prototype.IntelData != nil {
			if item.ExpiresAt != nil {
				return fmt.Errorf("intel item with ID %s is already being decrypted", itemID)
			}
			expiresAt := now + item.Prototype.IntelData.DecryptionSeconds
			ub.StorageItemsPresent[i].ExpiresAt = &expiresAt
			ub.StorageItemsPresent[i].IsActive = true

			ub.AddEvent(NewIntelDecryptionStartedEvent(ub.ID, itemID))
			return nil
		}
	}
	return fmt.Errorf("intel storage item with ID %s not found or not an intel item", itemID)
}

// StartDamagedItemRestorationByID starts the restoration process for a damaged item.
func (ub *UserBaseModel) StartDamagedItemRestorationByID(itemID uuid.UUID, armyProtos []*ArmyItemPrototype) error {
	if !ub.hasControlSubtype(ControlSubtypeRepairCenter) {
		return fmt.Errorf("repair center required to start restoration")
	}
	defer ub.recalculateStats()
	now := NowUnix()

	for i, item := range ub.StorageItemsPresent {
		if item.ID == itemID && item.Prototype.DamagedData != nil {
			if item.ExpiresAt != nil {
				return fmt.Errorf("damaged item with ID %s is already being restored", itemID)
			}
			data := item.Prototype.DamagedData

			// Find original unit prototype
			var unitProto *ArmyItemPrototype
			for _, p := range armyProtos {
				if p.ID == data.OriginalUnitID {
					unitProto = p
					break
				}
			}
			if unitProto == nil {
				return fmt.Errorf("original unit prototype %d not found", data.OriginalUnitID)
			}

			// Validate space
			if ub.Stats.Space+unitProto.Space > ub.Stats.SpaceCapacity {
				return fmt.Errorf("not enough space to restore unit: required %d, available %d", unitProto.Space, ub.Stats.SpaceCapacity-ub.Stats.Space)
			}

			// Validate resources
			if err := ub.Stats.CheckResources(data.RestorePrice); err != nil {
				return err
			}

			// Deduct price
			ub.Stats.SubtractResources(data.RestorePrice)

			// Start restoration
			expiresAt := now + data.RestorationSeconds
			ub.StorageItemsPresent[i].ExpiresAt = &expiresAt
			ub.StorageItemsPresent[i].IsActive = true

			ub.AddEvent(NewDamagedItemRestorationStartedEvent(ub.ID, itemID))
			return nil
		}
	}
	return fmt.Errorf("damaged storage item with ID %s not found or not a damaged item", itemID)
}

// ActivateArtifactByID enables the bonus of an artifact.
func (ub *UserBaseModel) ActivateArtifactByID(itemID uuid.UUID) error {
	defer ub.recalculateStats()
	for i, item := range ub.StorageItemsPresent {
		if item.ID == itemID && item.Prototype.ArtifactData != nil {
			if item.IsActive {
				return fmt.Errorf("artifact with ID %s is already active", itemID)
			}
			ub.StorageItemsPresent[i].IsActive = true
			ub.AddEvent(NewArtifactActivatedEvent(ub.ID, itemID))
			return nil
		}
	}
	return fmt.Errorf("artifact storage item with ID %s not found or not an artifact", itemID)
}

// DeactivateArtifactByID disables the bonus of an artifact.
func (ub *UserBaseModel) DeactivateArtifactByID(itemID uuid.UUID) error {
	defer ub.recalculateStats()
	for i, item := range ub.StorageItemsPresent {
		if item.ID == itemID && item.Prototype.ArtifactData != nil {
			if !item.IsActive {
				return fmt.Errorf("artifact with ID %s is not active", itemID)
			}
			ub.StorageItemsPresent[i].IsActive = false
			ub.AddEvent(NewArtifactDeactivatedEvent(ub.ID, itemID))
			return nil
		}
	}
	return fmt.Errorf("artifact storage item with ID %s not found or not an artifact", itemID)
}

// ArmyDeploymentRequest represents a request to allocate a number of units
// from a present army stack identified by its item ID.
type ArmyDeploymentRequest struct {
	PresentItemID uuid.UUID
	Count         int
}

// DeploymentReadyItem describes an army stack that is ready to be deployed; it ties a present item ID to its prototype and the requested count.
type DeploymentReadyItem struct {
	PresentItemID uuid.UUID
	Prototype     ArmyItemPrototype
	Count         int
}

// GetReadyToDeployArmy validates each deployment request against the current ArmiesPresent inventory
// and returns a list of DeploymentReadyItems that can be safely used to build operation units
// before the actual allocation mutates the base state.
func (ub *UserBaseModel) GetReadyToDeployArmy(requests []ArmyDeploymentRequest) ([]DeploymentReadyItem, error) {
	if len(requests) == 0 {
		return nil, fmt.Errorf("no units provided for deployment")
	}
	readyToDeploy := []DeploymentReadyItem{}
	for _, request := range requests {

		presentItemID, count := request.PresentItemID, request.Count

		idx := -1
		for i, p := range ub.ArmiesPresent {
			if p.ID == request.PresentItemID {
				idx = i
				break
			}
		}
		if idx == -1 {
			return []DeploymentReadyItem{}, fmt.Errorf("present army item %s not found", presentItemID)
		}
		p := ub.ArmiesPresent[idx]
		if count < 1 || count > p.Count {
			return []DeploymentReadyItem{}, fmt.Errorf("invalid count %d (available %d)", count, p.Count)
		}

		readyToDeploy = append(readyToDeploy, DeploymentReadyItem{
			PresentItemID: p.ID,
			Count:         count,
			Prototype:     p.Prototype,
		})
	}

	return readyToDeploy, nil
}

// AllocateArmyToOperation removes 'count' from a present stack, records them as deployed,
// and returns a deployed chunk describing what was allocated. Use conversion helpers
// to build OperationUnits for the military operation aggregate if needed.
func (ub *UserBaseModel) AllocateArmyToOperation(request ArmyDeploymentRequest, operationID int) (ArmyItemDeployed, error) {
	defer ub.recalculateStats()

	presentItemID, count := request.PresentItemID, request.Count

	idx := -1
	for i, p := range ub.ArmiesPresent {
		if p.ID == presentItemID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return ArmyItemDeployed{}, fmt.Errorf("present army item %s not found", presentItemID)
	}
	p := ub.ArmiesPresent[idx]
	if count < 1 || count > p.Count {
		return ArmyItemDeployed{}, fmt.Errorf("invalid count %d (available %d)", count, p.Count)
	}

	// Decrement/remove from present
	if count == p.Count {
		ub.ArmiesPresent = append(ub.ArmiesPresent[:idx], ub.ArmiesPresent[idx+1:]...)
	} else {
		ub.ArmiesPresent[idx].Count -= count
	}

	// Merge into deployed list by (operationID, prototypeID)
	merged := false
	for i, d := range ub.ArmiesDeployed {
		if d.OperationID == operationID && d.Prototype.ID == p.Prototype.ID {
			ub.ArmiesDeployed[i].Count += count
			merged = true
			break
		}
	}
	if !merged {
		ub.ArmiesDeployed = append(ub.ArmiesDeployed, ArmyItemDeployed{
			BaseOwnedItem: NewBaseOwnedItem(ub.ID),
			Prototype:     p.Prototype,
			OperationID:   operationID,
			Count:         count,
		})
	}

	deployedChunk := ArmyItemDeployed{
		BaseOwnedItem: NewBaseOwnedItem(ub.ID),
		Prototype:     p.Prototype,
		OperationID:   operationID,
		Count:         count,
	}
	// Optionally emit an allocation event
	// ub.AddEvent(NewArmyAllocatedToOperationEvent(ub.ID, operationID, p.ID, p.Prototype.ID, count))
	return deployedChunk, nil
}

// (deprecated) ReturnArmyFromOperation was replaced by ReturnAllDeployedFromOperation.

// CleanupDeployedForOperation removes any remaining deployed entries for an operation (e.g., on cancel/fail).
func (ub *UserBaseModel) CleanupDeployedForOperation(operationID int) {
	if len(ub.ArmiesDeployed) == 0 {
		return
	}
	out := ub.ArmiesDeployed[:0]
	for _, d := range ub.ArmiesDeployed {
		if d.OperationID != operationID {
			out = append(out, d)
		}
	}
	ub.ArmiesDeployed = out
	ub.recalculateStats()
}

// --- Aggregate methods used by military operations ---

// ReturnAllDeployedFromOperation merges all deployed units for the given operation
// back into ArmiesPresent and removes the deployed entries.
func (ub *UserBaseModel) ReturnAllDeployedFromOperation(operationID int) {
	// Merge deployed counts into present by prototype
	for _, d := range ub.ArmiesDeployed {
		if d.OperationID != operationID || d.Count <= 0 {
			continue
		}
		merged := false
		for i := range ub.ArmiesPresent {
			if ub.ArmiesPresent[i].Prototype.ID == d.Prototype.ID {
				ub.ArmiesPresent[i].Count += d.Count
				merged = true
				break
			}
		}
		if !merged {
			ub.ArmiesPresent = append(ub.ArmiesPresent, ArmyItemPresent{
				BaseOwnedItem: NewBaseOwnedItem(ub.ID),
				Prototype:     d.Prototype,
				Count:         d.Count,
				Refund:        d.Prototype.Price.Divide(10),
			})
		}
	}
	// Remove deployed entries for this operation and recalc stats
	ub.CleanupDeployedForOperation(operationID)
}

// TrimDeployedToSurvivors reduces deployed counts for this operation down to the
// survivors by prototype, removing any stacks that were completely destroyed.
func (ub *UserBaseModel) TrimDeployedToSurvivors(operationID int, survivors []MilitaryUnitSnap) {
	defer ub.recalculateStats()
	if len(ub.ArmiesDeployed) == 0 {
		return
	}
	remainByProto := map[int]int{}
	for _, u := range survivors {
		if u.Count > 0 {
			remainByProto[u.PrototypeID] += u.Count
		}
	}

	newDeployed := make([]ArmyItemDeployed, 0, len(ub.ArmiesDeployed))
	for _, d := range ub.ArmiesDeployed {
		if d.OperationID != operationID {
			newDeployed = append(newDeployed, d)
			continue
		}

		remain := remainByProto[d.Prototype.ID]
		if remain > 0 {
			if d.Count > remain {
				d.Count = remain
			}
			remainByProto[d.Prototype.ID] -= d.Count
			newDeployed = append(newDeployed, d)
		}
	}
	ub.ArmiesDeployed = newDeployed
}

// ApplyDefenderArmyRemaining sets ArmiesPresent counts to the provided remaining defenders
// by prototype ID, removing any stacks that were completely destroyed.
func (ub *UserBaseModel) ApplyDefenderArmyRemaining(remaining []MilitaryUnitSnap) {
	defer ub.recalculateStats()
	remainByProto := map[int]int{}
	for _, u := range remaining {
		if u.Count > 0 {
			remainByProto[u.PrototypeID] += u.Count
		}
	}

	newArmies := make([]ArmyItemPresent, 0, len(ub.ArmiesPresent))
	for _, p := range ub.ArmiesPresent {
		if newCount, ok := remainByProto[p.Prototype.ID]; ok && newCount > 0 {
			p.Count = newCount
			newArmies = append(newArmies, p)
		}
	}
	ub.ArmiesPresent = newArmies
}

// ApplyRemainingDefensiveStructures adjusts defensive BuildingsPresent to match the
// remaining structures (by PrototypeID). Non-defensive buildings are left untouched.
func (ub *UserBaseModel) ApplyRemainingDefensiveStructures(remaining []DefenseStructureSnap) {
	defer ub.recalculateStats()
	if len(ub.BuildingsPresent) == 0 {
		return
	}
	// Count how many instances of each defensive structure prototype should remain.
	keepCountByID := map[int]int{}
	for _, s := range remaining {
		if s.Count > 0 && s.PrototypeID != 0 {
			keepCountByID[s.PrototypeID] += s.Count
		}
	}
	filtered := make([]BuildItemPresent, 0, len(ub.BuildingsPresent))
	for _, b := range ub.BuildingsPresent {
		// Keep non-defensive buildings as-is.
		if b.Prototype.DefenseData == nil {
			filtered = append(filtered, b)
			continue
		}
		// For defensive buildings, keep up to the specified remaining count per prototype ID.
		id := b.Prototype.ID
		if keepCountByID[id] > 0 {
			filtered = append(filtered, b)
			keepCountByID[id]--
		}
		// Else: this defensive building instance was destroyed; drop it.
	}
	ub.BuildingsPresent = filtered
}

// DeductLoot subtracts the provided loot from the base's resources, clamped at zero.
func (ub *UserBaseModel) DeductLoot(loot PriceModel) {
	if loot.Credits > 0 {
		ub.Stats.Credits = maxInt(ub.Stats.Credits-loot.Credits, 0)
	}
	if loot.Iron > 0 {
		ub.Stats.Iron = maxInt(ub.Stats.Iron-loot.Iron, 0)
	}
	if loot.Titanium > 0 {
		ub.Stats.Titanium = maxInt(ub.Stats.Titanium-loot.Titanium, 0)
	}
	if loot.Antimatter > 0 {
		ub.Stats.Antimatter = maxInt(ub.Stats.Antimatter-loot.Antimatter, 0)
	}
}

// CreditLoot adds the provided loot to the base's resources, clamped by capacities.
func (ub *UserBaseModel) CreditLoot(loot PriceModel) {
	if loot.Credits > 0 {
		ub.Stats.Credits = min(ub.Stats.Credits+loot.Credits, ub.Stats.CreditsCapacity)
	}
	if loot.Iron > 0 {
		ub.Stats.Iron = min(ub.Stats.Iron+loot.Iron, ub.Stats.IronCapacity)
	}
	if loot.Titanium > 0 {
		ub.Stats.Titanium = min(ub.Stats.Titanium+loot.Titanium, ub.Stats.TitaniumCapacity)
	}
	if loot.Antimatter > 0 {
		ub.Stats.Antimatter = min(ub.Stats.Antimatter+loot.Antimatter, ub.Stats.AntimatterCapacity)
	}
}

// AddTrophies adds the provided trophies to the base's storage.
// It requires the prototypes to be resolved by the caller.
func (ub *UserBaseModel) AddTrophies(trophies []TrophyStorageItem, protos map[int]StorageItemPrototype) {
	for _, t := range trophies {
		if p, ok := protos[t.PrototypeID]; ok {
			ub.StorageItemsPresent = append(ub.StorageItemsPresent, StorageItemPresent{
				BaseOwnedItem: NewBaseOwnedItem(ub.ID),
				Prototype:     p,
			})
		}
	}
}

func (ub *UserBaseModel) TotalRadarStealthStrength() int {
	total := 0
	for _, b := range ub.BuildingsPresent {
		if b.Prototype.IntelligenceData != nil && b.Prototype.IntelligenceData.Subtype == IntelligenceSubtypeRadar {
			total += b.Prototype.IntelligenceData.StealthStrength
		}
	}
	return total
}

// ActiveModifiers returns the currently-active multipliers.
// Buffs must be IsActive and non-expired; artifacts must be IsActive.
func (ub *UserBaseModel) ActiveModifiers() BaseModifiers {
	m := IdentityBaseModifiers()
	now := NowUnix()

	for _, it := range ub.StorageItemsPresent {
		if !it.IsActive {
			continue
		}

		// Timed buffs
		if it.Prototype.BuffData != nil {
			if it.ExpiresAt != nil && *it.ExpiresAt <= now {
				continue
			}
			m.ApplyBuff(it.Prototype.BuffData.Type, float64(it.Prototype.BuffData.Value))
			continue
		}

		// Permanent artifacts (toggleable by IsActive)
		if it.Prototype.ArtifactData != nil {
			m.ApplyArtifact(it.Prototype.ArtifactData.Type, float64(it.Prototype.ArtifactData.Value))
			continue
		}
	}

	return m
}

// Default capacities and stats for UserBaseStats
const (
	DefaultCreditsCapacity    = 10000
	DefaultIronCapacity       = 5000
	DefaultTitaniumCapacity   = 2500
	DefaultAntimatterCapacity = 1000
	DefaultDefence            = 100
	DefaultAttack             = 0
	DefaultSpaceCapacity      = 50
)

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

func (s *UserBaseStats) CheckResources(price PriceModel) error {
	if price.Credits > s.Credits {
		return fmt.Errorf("insufficient credits")
	}
	if price.Iron > s.Iron {
		return fmt.Errorf("insufficient iron")
	}
	if price.Titanium > s.Titanium {
		return fmt.Errorf("insufficient titanium")
	}
	if price.Antimatter > s.Antimatter {
		return fmt.Errorf("insufficient antimatter")
	}
	return nil
}

func (s *UserBaseStats) SubtractResources(price PriceModel) {
	s.Credits -= price.Credits
	s.Iron -= price.Iron
	s.Titanium -= price.Titanium
	s.Antimatter -= price.Antimatter
}

// RecalculateStats updates the UserBaseStats based on present items and default constants.
func (ub *UserBaseModel) recalculateStats() {
	stats := UserBaseStats{}
	// Set default capacities
	stats.CreditsCapacity = DefaultCreditsCapacity
	stats.IronCapacity = DefaultIronCapacity
	stats.TitaniumCapacity = DefaultTitaniumCapacity
	stats.AntimatterCapacity = DefaultAntimatterCapacity
	stats.Defence = DefaultDefence
	stats.Attack = DefaultAttack
	stats.SpaceCapacity = DefaultSpaceCapacity

	// Aggregate bonuses from present buildings
	for _, b := range ub.BuildingsPresent {
		proto := b.Prototype
		// Resources buildings
		if proto.ResourcesData != nil {
			stats.CreditsCapacity += proto.ResourcesData.CreditsCapacity
			stats.IronCapacity += proto.ResourcesData.IronCapacity
			stats.TitaniumCapacity += proto.ResourcesData.TitaniumCapacity
			stats.AntimatterCapacity += proto.ResourcesData.AntimatterCapacity
			stats.CreditsProduction += proto.ResourcesData.CreditsProduction
			stats.IronProduction += proto.ResourcesData.IronProduction
			stats.TitaniumProduction += proto.ResourcesData.TitaniumProduction
			stats.AntimatterProduction += proto.ResourcesData.AntimatterProduction
		}
		// Defense buildings
		if proto.DefenseData != nil {
			stats.Defence += proto.DefenseData.DefenceBonus
		}
		// Space is always added
		stats.Space += proto.Space
	}

	// Include space from buildings in production
	for _, b := range ub.BuildingsInProduction {
		stats.Space += b.Prototype.Space
	}

	// Include space from buildings pending
	for _, b := range ub.BuildingsPending {
		stats.Space += b.Prototype.Space
	}
	// Include space from armies present
	for _, a := range ub.ArmiesPresent {
		stats.Space += a.Prototype.Space * a.Count
	}

	// Include space from armies deployed (still occupy capacity)
	for _, d := range ub.ArmiesDeployed {
		stats.Space += d.Prototype.Space * d.Count
	}

	// Include space from armies in production
	for _, a := range ub.ArmiesInProduction {
		stats.Space += a.Prototype.Space
	}

	// Include space from armies pending
	for _, a := range ub.ArmiesPending {
		stats.Space += a.Prototype.Space * a.Count
	}

	// Aggregate power from present armies
	for _, a := range ub.ArmiesPresent {
		stats.Defence += a.Prototype.Defence
		stats.Attack += a.Prototype.Attack
	}

	// Apply researched technology effects (additive for space/attack/defence, multiplicative % for production)
	// EffectTypeResourceBonus is treated as a percentage boost across all resource production rates.
	for _, tech := range ub.TechnologiesDone {
		for _, eff := range tech.Prototype.Effects {
			switch eff.EffectType {
			case EffectTypeSpaceBonus:
				stats.SpaceCapacity += eff.Value
			case EffectTypeDefenceBonus:
				stats.Defence += eff.Value
			case EffectTypeAttackBonus:
				stats.Attack += eff.Value
			case EffectTypeResourceBonus:
				mult := 1 + float64(eff.Value)/100.0
				if mult < 0 { // guard against negative values producing negative production
					mult = 0
				}
				stats.CreditsProduction *= mult
				stats.IronProduction *= mult
				stats.TitaniumProduction *= mult
				stats.AntimatterProduction *= mult
			}
		}
	}

	// Apply modifiers from storage items (buffs and artifacts)
	mods := ub.ActiveModifiers()
	stats.CreditsProduction *= mods.CreditsProdMul
	stats.IronProduction *= mods.IronProdMul
	stats.TitaniumProduction *= mods.TitaniumProdMul
	// (Antimatter production doesn't have a specific multiplier in the current BuffTypes/ArtifactTypes)

	stats.Attack = mulInt(stats.Attack, mods.AttackMul)
	stats.Defence = mulInt(stats.Defence, mods.DefenceMul)

	// Calculate current resources based on previous value, production rate, and elapsed time
	prevStats := ub.Stats
	now := NowUnix()
	delta := now - prevStats.CalculationTimestamp

	if delta > 0 {
		stats.Credits = prevStats.Credits + int(stats.CreditsProduction*float64(delta))
		if stats.Credits > stats.CreditsCapacity {
			stats.Credits = stats.CreditsCapacity
		}

		stats.Iron = prevStats.Iron + int(stats.IronProduction*float64(delta))
		if stats.Iron > stats.IronCapacity {
			stats.Iron = stats.IronCapacity
		}

		stats.Titanium = prevStats.Titanium + int(stats.TitaniumProduction*float64(delta))
		if stats.Titanium > stats.TitaniumCapacity {
			stats.Titanium = stats.TitaniumCapacity
		}

		stats.Antimatter = prevStats.Antimatter + int(stats.AntimatterProduction*float64(delta))
		if stats.Antimatter > stats.AntimatterCapacity {
			stats.Antimatter = stats.AntimatterCapacity
		}
	} else {
		stats.Credits = prevStats.Credits
		stats.Iron = prevStats.Iron
		stats.Titanium = prevStats.Titanium
		stats.Antimatter = prevStats.Antimatter
	}

	stats.CalculationTimestamp = now
	ub.Stats = stats
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
