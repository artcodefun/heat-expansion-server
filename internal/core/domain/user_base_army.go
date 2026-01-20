package domain

import (
	"fmt"

	"github.com/google/uuid"
)

// Returns all army prototypes the user can create based on unlocked technologies and present military buildings
func (ub *UserBaseModel) AvailableArmies(allPrototypes []*ArmyItemPrototype) []*ArmyItemPrototype {
	available := []*ArmyItemPrototype{}
	for _, proto := range allPrototypes {
		// Check tech unlock
		if proto.UnlockTechnologyID != nil && !ub.HasTech(*proto.UnlockTechnologyID) {
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
	if totalSpace > ub.Stats.MaxSpace {
		return fmt.Errorf("not enough space to queue army: required %d, available %d", totalSpace, ub.Stats.MaxSpace)
	}

	// Validate resources
	totalPrice := PriceModel{
		Credits:    proto.Price.Credits * count,
		Iron:       proto.Price.Iron * count,
		Titanium:   proto.Price.Titanium * count,
		Antimatter: proto.Price.Antimatter * count,
	}
	if err := ub.Stats.CheckResources(totalPrice); err != nil {
		return err
	}
	// Subtract price
	ub.Stats.SubtractResources(totalPrice)
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
			crystalsSkipPrice := max(1, int(pending.Prototype.ProductionTime/60))
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
	refund := PriceModel{
		Credits:    item.Prototype.Price.Credits * count,
		Iron:       item.Prototype.Price.Iron * count,
		Titanium:   item.Prototype.Price.Titanium * count,
		Antimatter: item.Prototype.Price.Antimatter * count,
	}
	ub.CreditLoot(refund)
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
	refund := PriceModel{
		Credits:    item.Refund.Credits * count,
		Iron:       item.Refund.Iron * count,
		Titanium:   item.Refund.Titanium * count,
		Antimatter: item.Refund.Antimatter * count,
	}
	ub.CreditLoot(refund)
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

	// Check MaxOperations limit
	activeOps := make(map[int]bool)
	for _, d := range ub.ArmiesDeployed {
		activeOps[d.OperationID] = true
	}
	if !activeOps[operationID] && len(activeOps) >= ub.Stats.MaxOperations {
		return ArmyItemDeployed{}, fmt.Errorf("maximum number of operations (%d) reached", ub.Stats.MaxOperations)
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
