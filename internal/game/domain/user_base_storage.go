package domain

import (
	"github.com/google/uuid"
)

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
		return NewError("error.domain.storage.present_not_found", H{"item_id": itemID})
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

	// Check current active buffs count
	activeCount := 0
	for _, item := range ub.StorageItemsPresent {
		if item.IsActive && item.Prototype.BuffData != nil {
			activeCount++
		}
	}
	if activeCount >= ub.Stats.MaxActiveBuffs {
		return NewError("error.domain.storage.max_buffs_reached", H{"max": ub.Stats.MaxActiveBuffs})
	}

	for i, item := range ub.StorageItemsPresent {
		if item.ID == itemID && item.Prototype.BuffData != nil {
			// Check if already activated
			if item.ExpiresAt != nil {
				return NewError("error.domain.storage.buff_already_active", H{"item_id": itemID})
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
	return NewError("error.domain.storage.not_a_buff", H{"item_id": itemID})
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
		return "", NewError("error.domain.storage.cryptography_lab_required", nil)
	}
	now := NowUnix()
	idx := -1
	for i, item := range ub.StorageItemsPresent {
		if item.ID == itemID && item.Prototype.IntelData != nil {
			if item.ExpiresAt == nil || *item.ExpiresAt > now {
				return "", NewError("error.domain.storage.intel_not_ready", H{"item_id": itemID})
			}
			idx = i
			break
		}
	}
	if idx == -1 {
		return "", NewError("error.domain.storage.intel_not_found", H{"item_id": itemID})
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
		return NewError("error.domain.storage.repair_center_required", nil)
	}
	defer ub.recalculateStats()
	now := NowUnix()
	idx := -1
	for i, item := range ub.StorageItemsPresent {
		if item.ID == itemID && item.Prototype.DamagedData != nil {
			if item.ExpiresAt == nil || *item.ExpiresAt > now {
				return NewError("error.domain.storage.damaged_not_ready", H{"item_id": itemID})
			}
			idx = i
			break
		}
	}
	if idx == -1 {
		return NewError("error.domain.storage.damaged_not_found", H{"item_id": itemID})
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
		return NewError("error.domain.storage.unit_prototype_not_found", H{"proto_id": data.OriginalUnitID})
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
		return NewError("error.domain.storage.cryptography_lab_required", nil)
	}
	defer ub.recalculateStats()
	now := NowUnix()

	// Check current active decryptions count
	activeCount := 0
	for _, item := range ub.StorageItemsPresent {
		if item.IsActive && item.Prototype.IntelData != nil && item.ExpiresAt != nil && *item.ExpiresAt > now {
			activeCount++
		}
	}
	if activeCount >= ub.Stats.MaxActiveDecryptions {
		return NewError("error.domain.storage.max_decryptions_reached", H{"max": ub.Stats.MaxActiveDecryptions})
	}

	for i, item := range ub.StorageItemsPresent {
		if item.ID == itemID && item.Prototype.IntelData != nil {
			if item.ExpiresAt != nil {
				return NewError("error.domain.storage.intel_already_decrypting", H{"item_id": itemID})
			}
			expiresAt := now + item.Prototype.IntelData.DecryptionSeconds
			ub.StorageItemsPresent[i].ExpiresAt = &expiresAt
			ub.StorageItemsPresent[i].IsActive = true

			ub.AddEvent(NewIntelDecryptionStartedEvent(ub.ID, itemID))
			return nil
		}
	}
	return NewError("error.domain.storage.not_intel", H{"item_id": itemID})
}

// StartDamagedItemRestorationByID starts the restoration process for a damaged item.
func (ub *UserBaseModel) StartDamagedItemRestorationByID(itemID uuid.UUID, armyProtos []*ArmyItemPrototype) error {
	if !ub.hasControlSubtype(ControlSubtypeRepairCenter) {
		return NewError("error.domain.storage.repair_center_required", nil)
	}
	defer ub.recalculateStats()
	now := NowUnix()

	// Check current active restorations count
	activeRestorations := 0
	for _, item := range ub.StorageItemsPresent {
		if item.Prototype.DamagedData != nil && item.ExpiresAt != nil && *item.ExpiresAt > now {
			activeRestorations++
		}
	}
	if activeRestorations >= ub.Stats.MaxActiveRestorations {
		return NewError("error.domain.storage.max_restorations_reached", H{"max": ub.Stats.MaxActiveRestorations})
	}

	for i, item := range ub.StorageItemsPresent {
		if item.ID == itemID && item.Prototype.DamagedData != nil {
			if item.ExpiresAt != nil {
				return NewError("error.domain.storage.damaged_already_restoring", H{"item_id": itemID})
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
				return NewError("error.domain.storage.unit_prototype_not_found", H{"proto_id": data.OriginalUnitID})
			}

			// Validate space
			if ub.Stats.Space+unitProto.Space > ub.Stats.MaxSpace {
				return NewError("error.domain.storage.not_enough_space", H{"required": unitProto.Space, "available": ub.Stats.MaxSpace - ub.Stats.Space})
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
	return NewError("error.domain.storage.not_damaged", H{"item_id": itemID})
}

// ActivateArtifactByID enables the bonus of an artifact.
func (ub *UserBaseModel) ActivateArtifactByID(itemID uuid.UUID) error {
	if !ub.hasControlSubtype(ControlSubtypeArtifactLab) {
		return NewError("error.domain.storage.artifact_lab_required", nil)
	}
	defer ub.recalculateStats()

	// Check current active artifacts count
	activeCount := 0
	for _, item := range ub.StorageItemsPresent {
		if item.IsActive && item.Prototype.ArtifactData != nil {
			activeCount++
		}
	}
	if activeCount >= ub.Stats.MaxActiveArtifacts {
		return NewError("error.domain.storage.max_artifacts_reached", H{"max": ub.Stats.MaxActiveArtifacts})
	}

	for i, item := range ub.StorageItemsPresent {
		if item.ID == itemID && item.Prototype.ArtifactData != nil {
			if item.IsActive {
				return NewError("error.domain.storage.artifact_already_active", H{"item_id": itemID})
			}
			ub.StorageItemsPresent[i].IsActive = true
			ub.AddEvent(NewArtifactActivatedEvent(ub.ID, itemID))
			return nil
		}
	}
	return NewError("error.domain.storage.not_artifact", H{"item_id": itemID})
}

// DeactivateArtifactByID disables the bonus of an artifact.
func (ub *UserBaseModel) DeactivateArtifactByID(itemID uuid.UUID) error {
	defer ub.recalculateStats()
	for i, item := range ub.StorageItemsPresent {
		if item.ID == itemID && item.Prototype.ArtifactData != nil {
			if !item.IsActive {
				return NewError("error.domain.storage.artifact_not_active", H{"item_id": itemID})
			}
			ub.StorageItemsPresent[i].IsActive = false
			ub.AddEvent(NewArtifactDeactivatedEvent(ub.ID, itemID))
			return nil
		}
	}
	return NewError("error.domain.storage.not_artifact", H{"item_id": itemID})
}

// ActiveModifiers returns the currently-active multipliers.
// Buffs must be IsActive and non-expired; artifacts must be IsActive.
func (ub *UserBaseModel) ActiveModifiers() BaseModifiers {
	return ModifiersFromSnaps(ub.ActiveStorageSnaps())
}

// ActiveStorageSnaps returns snapshots of all currently active storage items (buffs/artifacts).
func (ub *UserBaseModel) ActiveStorageSnaps() []StorageItemSnap {
	var active []StorageItemSnap
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
			active = append(active, StorageItemFromPresent(it))
			continue
		}

		// Permanent artifacts (toggleable by IsActive)
		if it.Prototype.ArtifactData != nil {
			active = append(active, StorageItemFromPresent(it))
			continue
		}
	}

	return active
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
