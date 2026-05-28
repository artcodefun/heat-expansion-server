package domain

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestStorage_StartIntelDecryption_RespectsMaxDecryptions(t *testing.T) {
	base := newBaseWithDefaults(35)
	base.BuildingsPresent = append(base.BuildingsPresent, BuildItemPresent{
		Prototype: BuildItemPrototype{
			Category: BuildCategoryControl,
			ControlData: &ControlBuildingData{
				Subtype: ControlSubtypeCryptographyLab,
			},
		},
	})
	base.recalculateStats() // MaxActiveDecryptions = 1

	proto := StorageItemPrototype{
		ID:        500,
		IntelData: &IntelStorageData{Type: HiddenLocationTypeUserBase, DecryptionSeconds: 100},
	}

	item1ID := base.AddStorageItem(proto, nil)
	item2ID := base.AddStorageItem(proto, nil)

	// first one should start
	if err := base.StartIntelDecryptionByID(item1ID); err != nil {
		t.Fatalf("failed to start first decryption: %v", err)
	}

	// second one should fail
	if err := base.StartIntelDecryptionByID(item2ID); err == nil {
		t.Errorf("expected error for exceeding MaxActiveDecryptions (1), got nil")
	} else if !strings.HasPrefix(err.Error(), "error.domain.storage.max_decryptions_reached") {
		t.Errorf("unexpected error message: %v", err)
	}

	// Increase limit via tech
	techProto := TechItemPrototype{
		ID: 600,
		Improvement: &TechImprovement{
			Type:  ImprovementTypeActiveDecryptionsCount,
			Value: 1,
		},
	}
	base.TechnologiesDone = append(base.TechnologiesDone, TechItemDone{
		Prototype: techProto,
		Level:     1,
	})
	base.recalculateStats() // MaxActiveDecryptions = 2

	// Now it should work
	if err := base.StartIntelDecryptionByID(item2ID); err != nil {
		t.Errorf("failed to start second decryption after tech upgrade: %v", err)
	}
}

func TestStorage_BuffActivateAndExpire(t *testing.T) {
	SetTestNow(t, 20_000)
	base := newBaseWithDefaults(4)
	// add a buff storage item
	buff := StorageItemPresent{
		BaseOwnedItem: NewBaseOwnedItem(base.ID),
		Prototype: StorageItemPrototype{
			ID:       300,
			Name:     "Space Booster",
			Category: StorageCategoryBuff,
			BuffData: &BuffStorageData{DurationSeconds: 100},
		},
	}
	base.StorageItemsPresent = []StorageItemPresent{buff}

	// Activate
	if err := base.ActivateBuffByID(buff.ID); err != nil {
		t.Fatalf("ActivateBuffByID error: %v", err)
	}
	events := base.PullEvents()
	if len(events) == 0 {
		t.Fatalf("expected BuffActivatedEvent")
	}
	if _, ok := events[0].(BuffActivatedEvent); !ok {
		t.Fatalf("expected BuffActivatedEvent, got %T", events[0])
	}
	if base.StorageItemsPresent[0].ExpiresAt == nil || *base.StorageItemsPresent[0].ExpiresAt != 20_100 {
		t.Fatalf("expected ExpiresAt=20100")
	}

	// advance time past expiration and delete
	SetTestNow(t, 20_200)
	deleted := base.DeleteExpiredBuffs()
	if deleted != 1 {
		t.Fatalf("expected 1 expired buff deleted, got %d", deleted)
	}
}

func TestStorage_MoveIntelDecryptionQueue_CompletesDecryption(t *testing.T) {
	SetTestNow(t, 30_000)
	base := newBaseWithDefaults(10)
	// mock crypto lab
	base.BuildingsPresent = append(base.BuildingsPresent, BuildItemPresent{
		Prototype: BuildItemPrototype{
			Category: BuildCategoryControl,
			ControlData: &ControlBuildingData{
				Subtype: ControlSubtypeCryptographyLab,
			},
		},
	})
	base.recalculateStats() // MaxActiveDecryptions = 1

	proto := StorageItemPrototype{
		ID:        501,
		IntelData: &IntelStorageData{Type: HiddenLocationTypeUserBase, DecryptionSeconds: 100},
	}

	// Add item and start manually
	itemID := base.AddStorageItem(proto, nil)
	idx := -1
	for i, it := range base.StorageItemsPresent {
		if it.ID == itemID {
			idx = i
			break
		}
	}
	expiresAt := int64(30_100)
	base.StorageItemsPresent[idx].ExpiresAt = &expiresAt
	base.StorageItemsPresent[idx].IsActive = true

	// Before completion
	_, err := base.DecryptIntelItemByID(itemID)
	if err == nil {
		t.Fatalf("expected error: intel not ready yet")
	}

	// After completion
	SetTestNow(t, 30_101)
	_, err = base.DecryptIntelItemByID(itemID)
	if err != nil {
		t.Fatalf("failed to decrypt: %v", err)
	}

	if len(base.StorageItemsPresent) != 0 {
		t.Fatalf("expected intel item to be removed after decryption")
	}

	// Should emit event
	events := base.PullEvents()
	found := false
	for _, e := range events {
		if _, ok := e.(IntelDecryptionFinishedEvent); ok {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected IntelDecryptionFinishedEvent")
	}
}

func TestStorage_DeleteStorageItemByID(t *testing.T) {
	base := newBaseWithDefaults(1)
	itemID := base.AddStorageItem(StorageItemPrototype{ID: 100}, nil)
	if len(base.StorageItemsPresent) != 1 {
		t.Fatal("item not added")
	}

	if err := base.DeletePresentStorageItemByID(itemID); err != nil {
		t.Fatalf("failed to delete item: %v", err)
	}
	if len(base.StorageItemsPresent) != 0 {
		t.Fatal("item not removed")
	}
}

func TestStorage_ReceiveStorageItem_AddsItem(t *testing.T) {
	base := newBaseWithDefaults(2)
	proto := StorageItemPrototype{
		ID:              200,
		CreationSources: []CreationSource{CreationSourcePlayerBase},
		Category:        StorageCategoryBuff,
	}

	if err := base.ReceiveStorageItem(proto); err != nil {
		t.Fatalf("ReceiveStorageItem error: %v", err)
	}
	if len(base.StorageItemsPresent) != 1 {
		t.Fatalf("expected one storage item, got %+v", base.StorageItemsPresent)
	}
}

func TestStorage_AddTradeDeployedStorageItems_NormalizesBaseOwnership(t *testing.T) {
	base := newBaseWithDefaults(42)
	otherBaseID := base.ID + 1000
	itemID := uuid.Must(uuid.NewV7())
	proto := StorageItemPrototype{ID: 901, Category: StorageCategoryArtifact}

	base.AddTradeDeployedStorageItems([]StorageItemDeployed{{
		BaseOwnedItem: BaseOwnedItem{ID: itemID, UserBaseID: otherBaseID},
		Prototype:     proto,
		OperationKind: OperationKindTrade,
		OperationID:   77,
	}}, 77)

	if len(base.StorageItemsDeployed) != 1 {
		t.Fatalf("expected one deployed storage item, got %+v", base.StorageItemsDeployed)
	}
	if base.StorageItemsDeployed[0].BaseOwnedItem.UserBaseID != base.ID {
		t.Fatalf("expected deployed storage item to be normalized to current base ID, got %d", base.StorageItemsDeployed[0].BaseOwnedItem.UserBaseID)
	}
	if base.StorageItemsDeployed[0].ID == itemID {
		t.Fatalf("expected deployed storage item to have new ID, not original %s", itemID)
	}

	base.ReturnAllDeployedStorageFromOperation(OperationKindTrade, 77)
	if len(base.StorageItemsPresent) != 1 {
		t.Fatalf("expected storage item to return to present, got %+v", base.StorageItemsPresent)
	}
	if base.StorageItemsPresent[0].BaseOwnedItem.UserBaseID != base.ID {
		t.Fatalf("expected restored storage item to be normalized to current base ID, got %d", base.StorageItemsPresent[0].BaseOwnedItem.UserBaseID)
	}
}

func TestStorage_ActivateArtifact_RequiresArtifactLab(t *testing.T) {
	base := newBaseWithDefaults(40)
	artifact := StorageItemPresent{
		BaseOwnedItem: NewBaseOwnedItem(base.ID),
		Prototype: StorageItemPrototype{
			ID:           800,
			Name:         "Ancient Relic",
			Category:     StorageCategoryArtifact,
			ArtifactData: &ArtifactStorageData{Type: "CombatBoost", Value: 0.1},
		},
	}
	base.StorageItemsPresent = []StorageItemPresent{artifact}

	// Try to activate without lab - should fail
	err := base.ActivateArtifactByID(artifact.ID)
	if err == nil {
		t.Errorf("expected error when activating artifact without Artifact Lab, got nil")
	} else if err.Error() != "error.domain.storage.artifact_lab_required" {
		t.Errorf("unexpected error message: %v", err)
	}

	// Add Artifact Lab
	base.BuildingsPresent = append(base.BuildingsPresent, BuildItemPresent{
		Prototype: BuildItemPrototype{
			Category: BuildCategoryControl,
			ControlData: &ControlBuildingData{
				Subtype: ControlSubtypeArtifactLab,
			},
		},
	})

	// Now activation should succeed
	if err := base.ActivateArtifactByID(artifact.ID); err != nil {
		t.Fatalf("failed to activate artifact after adding lab: %v", err)
	}

	if !base.StorageItemsPresent[0].IsActive {
		t.Errorf("expected artifact to be active")
	}
}

func TestStorage_ActivateBuffTwice_ErrorsAndDoesNotDuplicate(t *testing.T) {
	SetTestNow(t, 21_000)
	base := newBaseWithDefaults(6)
	buff := StorageItemPresent{
		BaseOwnedItem: NewBaseOwnedItem(base.ID),
		Prototype: StorageItemPrototype{
			ID:       400,
			Name:     "Space Booster",
			Category: StorageCategoryBuff,
			BuffData: &BuffStorageData{DurationSeconds: 50},
		},
	}
	base.StorageItemsPresent = []StorageItemPresent{buff}

	// first activation succeeds
	if err := base.ActivateBuffByID(buff.ID); err != nil {
		t.Fatalf("first ActivateBuffByID error: %v", err)
	}
	firstExpiresAt := base.StorageItemsPresent[0].ExpiresAt
	if firstExpiresAt == nil || *firstExpiresAt != 21_050 {
		t.Fatalf("expected ExpiresAt to be set on first activation, got %+v", firstExpiresAt)
	}
	base.PullEvents() // clear

	// second activation should return error and not change ExpiresAt or emit events
	if err := base.ActivateBuffByID(buff.ID); err == nil {
		t.Fatalf("expected error on second ActivateBuffByID for same buff")
	}
	secondExpiresAt := base.StorageItemsPresent[0].ExpiresAt
	if secondExpiresAt == nil || *secondExpiresAt != *firstExpiresAt {
		t.Fatalf("expected ExpiresAt to remain unchanged on second activation, got %+v", secondExpiresAt)
	}
	if events := base.PullEvents(); len(events) != 0 {
		t.Fatalf("expected no additional events on second activation, got %v", events)
	}
}
