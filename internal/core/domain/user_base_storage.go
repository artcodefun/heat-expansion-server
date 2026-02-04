package domain

import (
	"fmt"

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

	// Check current active buffs count
	activeCount := 0
	for _, item := range ub.StorageItemsPresent {
		if item.IsActive && item.Prototype.BuffData != nil {
			activeCount++
		}
	}
	if activeCount >= ub.Stats.MaxActiveBuffs {
		return fmt.Errorf("maximum number of active buffs (%d) reached", ub.Stats.MaxActiveBuffs)
	}

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
		return fmt.Errorf("maximum number of simultaneous intel decryptions (%d) reached", ub.Stats.MaxActiveDecryptions)
	}

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

	// Check current active restorations count
	activeRestorations := 0
	for _, item := range ub.StorageItemsPresent {
		if item.Prototype.DamagedData != nil && item.ExpiresAt != nil && *item.ExpiresAt > now {
			activeRestorations++
		}
	}
	if activeRestorations >= ub.Stats.MaxActiveRestorations {
		return fmt.Errorf("maximum number of simultaneous restorations (%d) reached", ub.Stats.MaxActiveRestorations)
	}

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
			if ub.Stats.Space+unitProto.Space > ub.Stats.MaxSpace {
				return fmt.Errorf("not enough space to restore unit: required %d, available %d", unitProto.Space, ub.Stats.MaxSpace-ub.Stats.Space)
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
	if !ub.hasControlSubtype(ControlSubtypeArtifactLab) {
		return fmt.Errorf("artifact laboratory required to activate artifacts")
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
		return fmt.Errorf("maximum number of active artifacts (%d) reached", ub.Stats.MaxActiveArtifacts)
	}

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
