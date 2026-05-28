package domain

import (
	"slices"

	"github.com/google/uuid"
)

// Returns all building prototypes the user can create based on unlocked technologies
func (ub *UserBaseModel) AvailableBuildings(allPrototypes []*BuildItemPrototype) []*BuildItemPrototype {
	available := []*BuildItemPrototype{}
	for _, proto := range allPrototypes {
		if !slices.Contains(proto.CreationSources, CreationSourcePlayerBase) {
			continue
		}
		// Players can only build EXO_COALITION buildings
		if proto.Faction != FactionExoCoalition {
			continue
		}
		if proto.UnlockTechnologyID == nil || ub.HasTech(*proto.UnlockTechnologyID) {
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
	ub.recalculateStats()
	defer ub.recalculateStats()

	// Ensure this prototype is actually available for this base
	if len(ub.AvailableBuildings([]*BuildItemPrototype{proto})) == 0 {
		return NewError("error.domain.building.not_available_for_production", nil)
	}

	// Calculate total space after adding this building
	totalSpace := ub.Stats.Space + proto.Space
	if totalSpace > ub.Stats.MaxSpace {
		return NewError("error.domain.building.not_enough_space", H{
			"required":  totalSpace,
			"available": ub.Stats.MaxSpace,
		})
	}

	// Validate resources (example: credits, iron, titanium, antimatter)
	if err := ub.Stats.CheckResources(proto.Price); err != nil {
		return err
	}
	// Subtract price from resources
	ub.Stats.SubtractResources(proto.Price)

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

// ReceiveBuilding instantly adds a completed building to the base.
func (ub *UserBaseModel) ReceiveBuilding(proto BuildItemPrototype) error {
	ub.recalculateStats()
	defer ub.recalculateStats()

	totalSpace := ub.Stats.Space + proto.Space
	if totalSpace > ub.Stats.MaxSpace {
		return NewError("error.domain.building.not_enough_space", H{"required": totalSpace, "available": ub.Stats.MaxSpace})
	}

	ub.BuildingsPresent = append(ub.BuildingsPresent, BuildItemPresent{
		BaseOwnedItem: NewBaseOwnedItem(ub.ID),
		Prototype:     proto,
		Refund:        proto.Price.Divide(10),
	})
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

	// Start next pending items up to MaxBuildingProduction limit
	for len(ub.BuildingsInProduction) < ub.Stats.MaxBuildingProduction && len(ub.BuildingsPending) > 0 {
		next := ub.BuildingsPending[0]
		ub.BuildingsPending = ub.BuildingsPending[1:]
		startDate := now
		completionDate := startDate + next.Prototype.ProductionTime
		crystalsSkipPrice := max(1, int(next.Prototype.ProductionTime/60))
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
		return NewError("error.domain.building.pending_not_found", H{"item_id": itemID})
	}
	item := ub.BuildingsPending[idx]
	// Refund resources
	ub.CreditLoot(item.Prototype.Price)
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
		return NewError("error.domain.building.in_production_not_found", H{"item_id": buildingItemID})
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
		return NewError("error.domain.building.present_not_found", H{"item_id": itemID})
	}
	// Refund resources to base from item's Refund field
	ub.CreditLoot(item.Refund)
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

func (ub *UserBaseModel) TotalRadarStealthStrength() int {
	return ub.totalStealthStrengthForSubtype(IntelligenceSubtypeRadar)
}

func (ub *UserBaseModel) TotalCloakingStealthStrength() int {
	return ub.totalStealthStrengthForSubtype(IntelligenceSubtypeCloaking)
}

func (ub *UserBaseModel) TotalInterceptionStealthStrength() int {
	return ub.totalStealthStrengthForSubtype(IntelligenceSubtypeScanInterceptor)
}

func (ub *UserBaseModel) totalStealthStrengthForSubtype(subtype IntelligenceSubtype) int {
	total := 0
	for _, b := range ub.BuildingsPresent {
		if b.Prototype.IntelligenceData != nil && b.Prototype.IntelligenceData.Subtype == subtype {
			total += b.Prototype.IntelligenceData.StealthStrength
		}
	}
	return total
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

// EnsureStartingBuildingsPresent adds basic production buildings if they are missing.
func (ub *UserBaseModel) EnsureStartingBuildingsPresent(allPrototypes []*BuildItemPrototype) {
	defer ub.recalculateStats()

	var basicCredits, basicIron, basicInfantryBuilding *BuildItemPrototype
	for _, p := range allPrototypes {
		if p.Faction != FactionExoCoalition {
			continue
		}

		if p.Category == BuildCategoryResources && p.ResourcesData != nil {
			if p.ResourcesData.CreditsProduction > 0 && (basicCredits == nil || p.ResourcesData.CreditsProduction < basicCredits.ResourcesData.CreditsProduction) {
				basicCredits = p
			}
			if p.ResourcesData.IronProduction > 0 && (basicIron == nil || p.ResourcesData.IronProduction < basicIron.ResourcesData.IronProduction) {
				basicIron = p
			}
		}

		if p.Category == BuildCategoryMilitary && p.MilitaryData != nil && p.MilitaryData.UnlockArmyCategory == ArmyCategoryInfantry {
			if basicInfantryBuilding == nil || p.Price.CreditsWorth() < basicInfantryBuilding.Price.CreditsWorth() {
				basicInfantryBuilding = p
			}
		}
	}

	hasPrototype := func(id int) bool {
		for _, b := range ub.BuildingsPresent {
			if b.Prototype.ID == id {
				return true
			}
		}
		return false
	}

	if basicCredits != nil && !hasPrototype(basicCredits.ID) {
		ub.BuildingsPresent = append(ub.BuildingsPresent, BuildItemPresent{
			BaseOwnedItem: NewBaseOwnedItem(ub.ID),
			Prototype:     *basicCredits,
			Refund:        basicCredits.Price.Divide(10),
		})
	}
	if basicIron != nil && !hasPrototype(basicIron.ID) {
		ub.BuildingsPresent = append(ub.BuildingsPresent, BuildItemPresent{
			BaseOwnedItem: NewBaseOwnedItem(ub.ID),
			Prototype:     *basicIron,
			Refund:        basicIron.Price.Divide(10),
		})
	}
	if basicInfantryBuilding != nil && !hasPrototype(basicInfantryBuilding.ID) {
		ub.BuildingsPresent = append(ub.BuildingsPresent, BuildItemPresent{
			BaseOwnedItem: NewBaseOwnedItem(ub.ID),
			Prototype:     *basicInfantryBuilding,
			Refund:        basicInfantryBuilding.Price.Divide(10),
		})
	}
}
