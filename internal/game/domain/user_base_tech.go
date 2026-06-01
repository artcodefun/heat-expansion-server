package domain

import (
	"github.com/google/uuid"
)

// Returns all technology prototypes the user can research based on unlocks and prerequisites
func (ub *UserBaseModel) AvailableTechnologies(allPrototypes []*TechItemPrototype) []*TechItemPrototype {
	available := []*TechItemPrototype{}
	for _, proto := range allPrototypes {
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

		// Check level limits
		currentLevel := ub.GetTechLevel(proto.ID)
		if proto.Improvement == nil {
			// Capability unlock tech: if level > 0, cannot research again
			if currentLevel > 0 {
				continue
			}
		} else if proto.Improvement.MaxLevel != nil && currentLevel >= *proto.Improvement.MaxLevel {
			// Scalable tech: check MaxLevel limit if present
			continue
		}

		// Check unlock condition (if any)
		if proto.UnlockTechnologyID != nil && !ub.HasTech(*proto.UnlockTechnologyID) {
			continue
		}
		available = append(available, proto)
	}
	return available
}

// HasTech returns true if the technology is researched at least at level 1.
func (ub *UserBaseModel) HasTech(techID int) bool {
	for _, t := range ub.TechnologiesDone {
		if t.Prototype.ID == techID {
			return true
		}
	}
	return false
}

// GetTechLevel returns the current researched level of a technology.
func (ub *UserBaseModel) GetTechLevel(techID int) int {
	for _, t := range ub.TechnologiesDone {
		if t.Prototype.ID == techID {
			return t.Level
		}
	}
	return 0
}

// StartTechResearch queues a technology for research
func (ub *UserBaseModel) StartTechResearch(proto *TechItemPrototype) error {
	ub.recalculateStats()
	defer ub.recalculateStats()
	// Ensure this prototype is actually available for this base
	if len(ub.AvailableTechnologies([]*TechItemPrototype{proto})) == 0 {
		return NewError("error.domain.tech.not_available_for_research", nil)
	}

	// Scale cost and duration by 0.5× per level (level 0→1 = 1×, 1→2 = 1.5×, 2→3 = 2×, …)
	multiplier := 1.0 + float64(ub.GetTechLevel(proto.ID))*0.5
	effectivePrice := proto.Price.MultiplyFloat(multiplier)
	effectiveResearchTime := int64(float64(proto.ResearchTime) * multiplier)

	// Validate resources
	if err := ub.Stats.CheckResources(effectivePrice); err != nil {
		return err
	}
	// Subtract price
	ub.Stats.SubtractResources(effectivePrice)
	// Add to in-progress
	now := NowUnix()
	completionDate := now + effectiveResearchTime
	crystalsSkipPrice := max(1, int(effectiveResearchTime/60))
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
			found := false
			for i := range ub.TechnologiesDone {
				if ub.TechnologiesDone[i].Prototype.ID == tech.Prototype.ID {
					ub.TechnologiesDone[i].Level++
					ub.TechnologiesDone[i].ResearchedAt = tech.CompletionDate
					found = true
					break
				}
			}
			if !found {
				done := TechItemDone{
					BaseOwnedItem: NewBaseOwnedItem(ub.ID),
					Prototype:     tech.Prototype,
					ResearchedAt:  tech.CompletionDate,
					Level:         1,
				}
				ub.TechnologiesDone = append(ub.TechnologiesDone, done)
			}
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
		return NewError("error.domain.tech.in_progress_not_found", H{"item_id": techItemID})
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
