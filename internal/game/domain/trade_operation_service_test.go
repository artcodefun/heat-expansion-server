package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestBuildTradeOperationForCreation_Success(t *testing.T) {
	SetTestNow(t, 15_000)

	sender := newBaseWithDefaults(10)
	sender.UserID = uuid.New()
	sender.Coordinates = Vector2i{X: 1, Y: 1}
	receiver := newBaseWithDefaults(11)
	receiver.UserID = uuid.New()
	receiver.Coordinates = Vector2i{X: 4, Y: 5}

	transportProto := ArmyItemPrototype{ID: 201, Category: ArmyCategoryInfantry, Capacity: 15, Speed: 90}
	offeredProto := ArmyItemPrototype{ID: 202, Category: ArmyCategoryInfantry, Capacity: 8, Speed: 95}
	requestedProto := ArmyItemPrototype{ID: 203, Category: ArmyCategoryInfantry, Capacity: 5, Speed: 100}

	sender.ArmiesPresent = []ArmyItemPresent{
		{BaseOwnedItem: NewBaseOwnedItem(sender.ID), Prototype: transportProto, Count: 4},
		{BaseOwnedItem: NewBaseOwnedItem(sender.ID), Prototype: offeredProto, Count: 3},
	}
	receiver.ArmiesPresent = []ArmyItemPresent{
		{BaseOwnedItem: NewBaseOwnedItem(receiver.ID), Prototype: requestedProto, Count: 2},
	}

	offeredStorageID := uuid.New()
	requestedStorageID := uuid.New()
	sender.StorageItemsPresent = []StorageItemPresent{{
		BaseOwnedItem: BaseOwnedItem{ID: offeredStorageID, UserBaseID: sender.ID},
		Prototype:     StorageItemPrototype{ID: 301, Category: StorageCategoryArtifact},
	}}
	receiver.StorageItemsPresent = []StorageItemPresent{{
		BaseOwnedItem: BaseOwnedItem{ID: requestedStorageID, UserBaseID: receiver.ID},
		Prototype:     StorageItemPrototype{ID: 302, Category: StorageCategoryConsumable},
	}}

	op, err := BuildTradeOperationForCreation(
		sender,
		receiver,
		PriceModel{Credits: 30},
		[]ArmyDeploymentRequest{{PresentItemID: sender.ArmiesPresent[1].ID, Count: 2}},
		[]uuid.UUID{offeredStorageID},
		PriceModel{Credits: 20},
		[]ArmyDeploymentRequest{{PresentItemID: receiver.ArmiesPresent[0].ID, Count: 1}},
		[]uuid.UUID{requestedStorageID},
		[]ArmyDeploymentRequest{{PresentItemID: sender.ArmiesPresent[0].ID, Count: 2}},
	)
	if err != nil {
		t.Fatalf("build trade operation failed: %v", err)
	}

	if op.SourceCoordinates != sender.Coordinates {
		t.Fatalf("expected source coordinates from sender base")
	}
	if op.TargetCoordinates != receiver.Coordinates {
		t.Fatalf("expected target coordinates from receiver base")
	}
	if len(op.TransportUnits) != 1 || op.TransportUnits[0].Count != 2 {
		t.Fatalf("expected transport units from transport requests")
	}
	if len(op.OfferedPayload.Storage) != 1 || op.OfferedPayload.Storage[0].ItemID != offeredStorageID {
		t.Fatalf("expected offered storage payload from sender item ids")
	}
	if len(op.RequestedPayload.Storage) != 1 || op.RequestedPayload.Storage[0].ItemID != requestedStorageID {
		t.Fatalf("expected requested storage payload from receiver item ids")
	}
}

func TestBuildTradeOperationForCreation_NilParticipants(t *testing.T) {
	_, err := BuildTradeOperationForCreation(
		nil,
		nil,
		PriceModel{Credits: 1},
		nil,
		nil,
		PriceModel{Credits: 1},
		nil,
		nil,
		nil,
	)
	if err == nil {
		t.Fatalf("expected participant mismatch error")
	}
}

func TestTradeOperationService_ValidateTradeDeploymentAllocationAvailability_AllowsSharedStackWithinAvailable(t *testing.T) {
	base := newBaseWithDefaults(1)
	shared := uuid.New()
	base.ArmiesPresent = []ArmyItemPresent{{
		BaseOwnedItem: BaseOwnedItem{ID: shared, UserBaseID: base.ID},
		Prototype:     ArmyItemPrototype{ID: 100, Category: ArmyCategoryInfantry, Capacity: 5, Speed: 100},
		Count:         5,
	}}

	err := validateTradeDeploymentAllocationAvailability(
		base,
		[]ArmyDeploymentRequest{{PresentItemID: shared, Count: 2}},
		[]ArmyDeploymentRequest{{PresentItemID: shared, Count: 1}},
	)
	if err != nil {
		t.Fatalf("expected shared stack allocation to be allowed, got: %v", err)
	}
}

func TestTradeOperationService_ValidateTradeDeploymentAllocationAvailability_RejectsWhenCombinedExceedsAvailable(t *testing.T) {
	base := newBaseWithDefaults(1)
	shared := uuid.New()
	base.ArmiesPresent = []ArmyItemPresent{{
		BaseOwnedItem: BaseOwnedItem{ID: shared, UserBaseID: base.ID},
		Prototype:     ArmyItemPrototype{ID: 100, Category: ArmyCategoryInfantry, Capacity: 5, Speed: 100},
		Count:         5,
	}}

	err := validateTradeDeploymentAllocationAvailability(
		base,
		[]ArmyDeploymentRequest{{PresentItemID: shared, Count: 4}},
		[]ArmyDeploymentRequest{{PresentItemID: shared, Count: 2}},
	)
	if err == nil {
		t.Fatalf("expected invalid deployment count error")
	}
}

func TestTradeOperationService_ValidatePayloadAvailability_StorageInactiveRequired(t *testing.T) {
	base := newBaseWithDefaults(100)
	itemID := uuid.New()
	base.StorageItemsPresent = []StorageItemPresent{{
		BaseOwnedItem: BaseOwnedItem{ID: itemID, UserBaseID: base.ID},
		Prototype:     StorageItemPrototype{ID: 10, Category: StorageCategoryArtifact},
		IsActive:      true,
	}}

	payload, err := NewTradePayload(
		PriceModel{Credits: 1},
		[]TradeStorageItemSnap{{ItemID: itemID, PrototypeID: 10, Category: StorageCategoryArtifact}},
		nil,
	)
	if err != nil {
		t.Fatalf("payload creation failed: %v", err)
	}

	service := &TradeOperationService{}
	err = service.validatePayloadAvailability(base, payload)
	if err == nil {
		t.Fatalf("expected not-tradeable storage error")
	}
}

func TestTradeOperationService_CommitSenderForTradeCreation(t *testing.T) {
	SetTestNow(t, 20_000)

	sender := newBaseWithDefaults(101)
	sender.UserID = uuid.New()
	receiver := newBaseWithDefaults(102)
	receiver.UserID = uuid.New()

	transportProto := ArmyItemPrototype{ID: 11, Category: ArmyCategoryInfantry, Capacity: 10, Speed: 100}
	offeredProto := ArmyItemPrototype{ID: 12, Category: ArmyCategoryInfantry, Capacity: 5, Speed: 90}

	sender.ArmiesPresent = []ArmyItemPresent{
		{BaseOwnedItem: NewBaseOwnedItem(sender.ID), Prototype: transportProto, Count: 3},
		{BaseOwnedItem: NewBaseOwnedItem(sender.ID), Prototype: offeredProto, Count: 4},
	}
	storageID := uuid.New()
	sender.StorageItemsPresent = []StorageItemPresent{
		{BaseOwnedItem: BaseOwnedItem{ID: storageID, UserBaseID: sender.ID}, Prototype: StorageItemPrototype{ID: 41, Category: StorageCategoryArtifact}},
	}

	offered, err := NewTradePayload(
		PriceModel{Credits: 100},
		[]TradeStorageItemSnap{{ItemID: storageID, PrototypeID: 41, Category: StorageCategoryArtifact}},
		[]TradeArmyItemSnap{{PrototypeID: offeredProto.ID, Count: 2}},
	)
	if err != nil {
		t.Fatalf("offered payload failed: %v", err)
	}
	requested, err := NewTradePayload(PriceModel{Credits: 50}, nil, nil)
	if err != nil {
		t.Fatalf("requested payload failed: %v", err)
	}
	op, err := NewTradeOperation(sender.UserID, sender.ID, receiver.UserID, receiver.ID, Vector2i{X: 1, Y: 1}, Vector2i{X: 2, Y: 3}, offered, requested, []MilitaryUnitSnap{{PrototypeID: transportProto.ID, Count: 1, Capacity: 10, Speed: 100}}, nil)
	if err != nil {
		t.Fatalf("operation creation failed: %v", err)
	}
	op.ID = 777

	service := NewTradeOperationService(sender, receiver, op)
	err = service.CommitSenderForTradeCreation()
	if err != nil {
		t.Fatalf("commit sender failed: %v", err)
	}

	if sender.Stats.Credits >= 1000 {
		t.Fatalf("expected sender credits to decrease")
	}
	if len(sender.StorageItemsPresent) != 0 {
		t.Fatalf("expected offered storage item removed")
	}
	if len(sender.StorageItemsDeployed) != 1 || sender.StorageItemsDeployed[0].ID != storageID {
		t.Fatalf("expected offered storage item deployed for operation")
	}
	if len(sender.ArmiesDeployed) == 0 {
		t.Fatalf("expected deployed armies after sender commit")
	}
}

func TestTradeOperationService_AcceptAndCommitReceiver(t *testing.T) {
	SetTestNow(t, 30_000)

	sender := newBaseWithDefaults(201)
	sender.UserID = uuid.New()
	receiver := newBaseWithDefaults(202)
	receiver.UserID = uuid.New()

	reqArmyProto := ArmyItemPrototype{ID: 21, Category: ArmyCategoryInfantry, Capacity: 7, Speed: 90}
	receiver.ArmiesPresent = []ArmyItemPresent{{BaseOwnedItem: NewBaseOwnedItem(receiver.ID), Prototype: reqArmyProto, Count: 3}}
	reqStorageID := uuid.New()
	receiver.StorageItemsPresent = []StorageItemPresent{{BaseOwnedItem: BaseOwnedItem{ID: reqStorageID, UserBaseID: receiver.ID}, Prototype: StorageItemPrototype{ID: 51, Category: StorageCategoryConsumable}}}

	offered, err := NewTradePayload(PriceModel{Credits: 10}, nil, nil)
	if err != nil {
		t.Fatalf("offered payload failed: %v", err)
	}
	requested, err := NewTradePayload(
		PriceModel{Credits: 20},
		[]TradeStorageItemSnap{{ItemID: reqStorageID, PrototypeID: 51, Category: StorageCategoryConsumable}},
		[]TradeArmyItemSnap{{PrototypeID: reqArmyProto.ID, Count: 2}},
	)
	if err != nil {
		t.Fatalf("requested payload failed: %v", err)
	}

	op, err := NewTradeOperation(sender.UserID, sender.ID, receiver.UserID, receiver.ID, Vector2i{X: 1, Y: 1}, Vector2i{X: 3, Y: 4}, offered, requested, []MilitaryUnitSnap{{PrototypeID: 99, Count: 1, Capacity: 20, Speed: 120}}, nil)
	if err != nil {
		t.Fatalf("operation creation failed: %v", err)
	}
	op.ID = 888

	service := NewTradeOperationService(sender, receiver, op)
	if err := service.AcceptAndCommitReceiver(); err != nil {
		t.Fatalf("commit receiver failed: %v", err)
	}

	if op.Phase != TradePhaseOutbound {
		t.Fatalf("expected operation phase OUTBOUND, got %s", op.Phase)
	}
	if len(receiver.StorageItemsPresent) != 0 {
		t.Fatalf("expected requested storage item removed from receiver")
	}
	if len(receiver.StorageItemsDeployed) != 1 || receiver.StorageItemsDeployed[0].ID != reqStorageID {
		t.Fatalf("expected requested storage item deployed from receiver")
	}
	if len(receiver.ArmiesDeployed) == 0 {
		t.Fatalf("expected requested army deployed from receiver")
	}
}

func TestTradeOperationService_CancelAndReleaseReceiverIfCommitted(t *testing.T) {
	SetTestNow(t, 40_000)

	sender := newBaseWithDefaults(301)
	sender.UserID = uuid.New()
	receiver := newBaseWithDefaults(302)
	receiver.UserID = uuid.New()

	reqArmyProto := ArmyItemPrototype{ID: 31, Category: ArmyCategoryInfantry, Capacity: 6, Speed: 100}
	receiver.ArmiesPresent = []ArmyItemPresent{{BaseOwnedItem: NewBaseOwnedItem(receiver.ID), Prototype: reqArmyProto, Count: 3}}
	reqStorageID := uuid.New()
	receiver.StorageItemsPresent = []StorageItemPresent{{BaseOwnedItem: BaseOwnedItem{ID: reqStorageID, UserBaseID: receiver.ID}, Prototype: StorageItemPrototype{ID: 61, Category: StorageCategoryConsumable}}}

	receiverCreditsBefore := receiver.Stats.Credits

	offered, err := NewTradePayload(PriceModel{Credits: 10}, nil, nil)
	if err != nil {
		t.Fatalf("offered payload failed: %v", err)
	}
	requested, err := NewTradePayload(
		PriceModel{Credits: 20},
		[]TradeStorageItemSnap{{ItemID: reqStorageID, PrototypeID: 61, Category: StorageCategoryConsumable}},
		[]TradeArmyItemSnap{{PrototypeID: reqArmyProto.ID, Count: 2}},
	)
	if err != nil {
		t.Fatalf("requested payload failed: %v", err)
	}

	op, err := NewTradeOperation(sender.UserID, sender.ID, receiver.UserID, receiver.ID, Vector2i{X: 1, Y: 1}, Vector2i{X: 3, Y: 4}, offered, requested, []MilitaryUnitSnap{{PrototypeID: 99, Count: 1, Capacity: 20, Speed: 120}}, nil)
	if err != nil {
		t.Fatalf("operation creation failed: %v", err)
	}
	op.ID = 999

	service := NewTradeOperationService(sender, receiver, op)
	if err := service.AcceptAndCommitReceiver(); err != nil {
		t.Fatalf("commit receiver failed: %v", err)
	}
	if op.Phase != TradePhaseOutbound {
		t.Fatalf("expected outbound before cancel, got %s", op.Phase)
	}
	if len(receiver.StorageItemsPresent) != 0 {
		t.Fatalf("expected requested storage removed on accept")
	}

	SetTestNow(t, 40_010)
	if err := service.CancelAndReleaseReceiverIfCommitted(); err != nil {
		t.Fatalf("cancel + receiver rollback failed: %v", err)
	}

	if op.Result != TradeResultCanceled {
		t.Fatalf("expected canceled operation result, got %s", op.Result)
	}
	if op.Phase != TradePhaseReturning {
		t.Fatalf("expected returning phase after outbound cancel, got %s", op.Phase)
	}

	if receiver.Stats.Credits != receiverCreditsBefore {
		t.Fatalf("expected receiver credits to be restored, got %v want %v", receiver.Stats.Credits, receiverCreditsBefore)
	}
	if len(receiver.StorageItemsPresent) != 1 || receiver.StorageItemsPresent[0].ID != reqStorageID {
		t.Fatalf("expected requested storage to be restored")
	}
	if len(receiver.StorageItemsDeployed) != 0 {
		t.Fatalf("expected receiver deployed storage to be released")
	}

	if len(receiver.ArmiesDeployed) != 0 {
		t.Fatalf("expected receiver deployed army to be released")
	}
	if len(receiver.ArmiesPresent) != 1 || receiver.ArmiesPresent[0].Count != 3 {
		t.Fatalf("expected receiver present army restored to original count")
	}
}

func TestTradeOperationService_FinalizeSenderAfterReturn_Success(t *testing.T) {
	SetTestNow(t, 50_000)

	sender := newBaseWithDefaults(401)
	sender.UserID = uuid.New()
	receiver := newBaseWithDefaults(402)
	receiver.UserID = uuid.New()

	requestedStorageID := uuid.New()
	offered, err := NewTradePayload(PriceModel{Credits: 5}, nil, nil)
	if err != nil {
		t.Fatalf("offered payload failed: %v", err)
	}
	requested, err := NewTradePayload(
		PriceModel{Credits: 20},
		[]TradeStorageItemSnap{{ItemID: requestedStorageID, PrototypeID: 71, Category: StorageCategoryArtifact}},
		nil,
	)
	if err != nil {
		t.Fatalf("requested payload failed: %v", err)
	}

	op, err := NewTradeOperation(sender.UserID, sender.ID, receiver.UserID, receiver.ID, Vector2i{X: 1, Y: 1}, Vector2i{X: 2, Y: 3}, offered, requested, []MilitaryUnitSnap{{PrototypeID: 99, Count: 1, Capacity: 20, Speed: 120}}, nil)
	if err != nil {
		t.Fatalf("operation creation failed: %v", err)
	}
	op.ID = 1234
	op.Phase = TradePhaseCompleted
	op.Result = TradeResultSuccess

	sender.Stats.Credits = 0
	senderCreditsBefore := sender.Stats.Credits
	sender.ArmiesDeployed = []ArmyItemDeployed{{
		BaseOwnedItem: NewBaseOwnedItem(sender.ID),
		Prototype:     ArmyItemPrototype{ID: 81, Category: ArmyCategoryInfantry, Capacity: 5, Speed: 90},
		OperationKind: OperationKindTrade,
		OperationID:   op.ID,
		Count:         2,
	}}
	sender.StorageItemsDeployed = []StorageItemDeployed{{
		BaseOwnedItem: BaseOwnedItem{ID: requestedStorageID, UserBaseID: sender.ID},
		Prototype:     StorageItemPrototype{ID: 71, Category: StorageCategoryArtifact, ArtifactData: &ArtifactStorageData{Type: ArtifactEffectTypeAttackIncrease, Value: 0.1}},
		OperationKind: OperationKindTrade,
		OperationID:   op.ID,
	}}

	service := NewTradeOperationService(sender, receiver, op)
	if err := service.FinalizeSenderAfterReturn(); err != nil {
		t.Fatalf("finalize sender failed: %v", err)
	}

	if sender.Stats.Credits != senderCreditsBefore+20 {
		t.Fatalf("expected sender credits increased by requested payload")
	}
	if len(sender.StorageItemsPresent) != 1 || sender.StorageItemsPresent[0].ID != requestedStorageID {
		t.Fatalf("expected requested storage restored to sender")
	}
	if sender.StorageItemsPresent[0].Prototype.ArtifactData == nil {
		t.Fatalf("expected full requested storage prototype data preserved")
	}
	if len(sender.StorageItemsDeployed) != 0 {
		t.Fatalf("expected sender deployed storage released")
	}
	if len(sender.ArmiesDeployed) != 0 {
		t.Fatalf("expected sender deployed units released")
	}
}

func TestTradeOperationService_FinalizeSenderAfterReturn_RestoresOfferedOnDeclined(t *testing.T) {
	SetTestNow(t, 55_000)

	sender := newBaseWithDefaults(501)
	sender.UserID = uuid.New()
	receiver := newBaseWithDefaults(502)
	receiver.UserID = uuid.New()

	offeredStorageID := uuid.New()
	offered, err := NewTradePayload(
		PriceModel{Credits: 40},
		[]TradeStorageItemSnap{{ItemID: offeredStorageID, PrototypeID: 91, Category: StorageCategoryArtifact}},
		nil,
	)
	if err != nil {
		t.Fatalf("offered payload failed: %v", err)
	}
	requested, err := NewTradePayload(PriceModel{Credits: 10}, nil, nil)
	if err != nil {
		t.Fatalf("requested payload failed: %v", err)
	}

	op, err := NewTradeOperation(sender.UserID, sender.ID, receiver.UserID, receiver.ID, Vector2i{X: 1, Y: 1}, Vector2i{X: 2, Y: 3}, offered, requested, []MilitaryUnitSnap{{PrototypeID: 99, Count: 1, Capacity: 20, Speed: 120}}, nil)
	if err != nil {
		t.Fatalf("operation creation failed: %v", err)
	}
	op.ID = 2001
	op.Phase = TradePhaseCompleted
	op.Result = TradeResultDeclined

	sender.Stats.Credits = 0
	sender.ArmiesDeployed = []ArmyItemDeployed{{
		BaseOwnedItem: NewBaseOwnedItem(sender.ID),
		Prototype:     ArmyItemPrototype{ID: 82, Category: ArmyCategoryInfantry, Capacity: 5, Speed: 90},
		OperationKind: OperationKindTrade,
		OperationID:   op.ID,
		Count:         3,
	}}
	sender.StorageItemsDeployed = []StorageItemDeployed{{
		BaseOwnedItem: BaseOwnedItem{ID: offeredStorageID, UserBaseID: sender.ID},
		Prototype:     StorageItemPrototype{ID: 91, Category: StorageCategoryArtifact, ArtifactData: &ArtifactStorageData{Type: ArtifactEffectTypeDefenceIncrease, Value: 0.2}},
		OperationKind: OperationKindTrade,
		OperationID:   op.ID,
	}}

	service := NewTradeOperationService(sender, receiver, op)
	if err := service.FinalizeSenderAfterReturn(); err != nil {
		t.Fatalf("finalize sender on declined failed: %v", err)
	}

	if sender.Stats.Credits != 40 {
		t.Fatalf("expected sender offered credits restored, got %v", sender.Stats.Credits)
	}
	if len(sender.StorageItemsPresent) != 1 || sender.StorageItemsPresent[0].ID != offeredStorageID {
		t.Fatalf("expected offered storage restored to sender")
	}
	if len(sender.StorageItemsDeployed) != 0 {
		t.Fatalf("expected sender deployed storage released on declined")
	}
	if len(sender.ArmiesDeployed) != 0 {
		t.Fatalf("expected sender deployed units released on declined")
	}
}

func TestTradeOperationService_ProcessArrivalAndStartReturn(t *testing.T) {
	SetTestNow(t, 70_000)

	sender := newBaseWithDefaults(701)
	sender.UserID = uuid.New()
	sender.Coordinates = Vector2i{X: 1, Y: 1}
	receiver := newBaseWithDefaults(702)
	receiver.UserID = uuid.New()
	receiver.Coordinates = Vector2i{X: 5, Y: 5}

	transportProto := ArmyItemPrototype{ID: 201, Category: ArmyCategoryInfantry, Capacity: 20, Speed: 100}
	offeredProto := ArmyItemPrototype{ID: 202, Category: ArmyCategoryInfantry, Capacity: 5, Speed: 80}
	requestedProto := ArmyItemPrototype{ID: 203, Category: ArmyCategoryInfantry, Capacity: 5, Speed: 80}

	offeredStorageID := uuid.New()
	offeredPayload, err := NewTradePayload(
		PriceModel{Credits: 50},
		[]TradeStorageItemSnap{{ItemID: offeredStorageID, PrototypeID: 301, Category: StorageCategoryArtifact}},
		[]TradeArmyItemSnap{{PrototypeID: offeredProto.ID, Count: 2, Capacity: offeredProto.Capacity}},
	)
	if err != nil {
		t.Fatalf("offered payload failed: %v", err)
	}
	requestedPayload, err := NewTradePayload(PriceModel{Credits: 30}, nil, []TradeArmyItemSnap{{PrototypeID: requestedProto.ID, Count: 1, Capacity: requestedProto.Capacity}})
	if err != nil {
		t.Fatalf("requested payload failed: %v", err)
	}

	op, err := NewTradeOperation(
		sender.UserID, sender.ID, receiver.UserID, receiver.ID,
		sender.Coordinates, receiver.Coordinates,
		offeredPayload, requestedPayload,
		[]MilitaryUnitSnap{{PrototypeID: transportProto.ID, Count: 1, Capacity: transportProto.Capacity, Speed: transportProto.Speed}},
		nil,
	)
	if err != nil {
		t.Fatalf("operation creation failed: %v", err)
	}
	op.ID = 3001
	op.Phase = TradePhaseArrived

	// Sender deployed: 1 transport + 2 offered army
	sender.ArmiesDeployed = []ArmyItemDeployed{
		{BaseOwnedItem: NewBaseOwnedItem(sender.ID), Prototype: transportProto, OperationKind: OperationKindTrade, OperationID: op.ID, Count: 1},
		{BaseOwnedItem: NewBaseOwnedItem(sender.ID), Prototype: offeredProto, OperationKind: OperationKindTrade, OperationID: op.ID, Count: 2},
	}
	sender.StorageItemsDeployed = []StorageItemDeployed{
		{BaseOwnedItem: BaseOwnedItem{ID: offeredStorageID, UserBaseID: sender.ID}, Prototype: StorageItemPrototype{ID: 301, Category: StorageCategoryArtifact, ArtifactData: &ArtifactStorageData{Type: ArtifactEffectTypeAttackIncrease, Value: 0.15}}, OperationKind: OperationKindTrade, OperationID: op.ID},
	}
	// Receiver deployed: 1 requested army
	receiver.ArmiesDeployed = []ArmyItemDeployed{
		{BaseOwnedItem: NewBaseOwnedItem(receiver.ID), Prototype: requestedProto, OperationKind: OperationKindTrade, OperationID: op.ID, Count: 1},
	}
	receiver.Stats.Credits = 0
	receiverCreditsBefore := receiver.Stats.Credits

	svc := NewTradeOperationService(sender, receiver, op)
	if err := svc.ProcessArrivalAndStartReturn(); err != nil {
		t.Fatalf("ProcessArrivalAndStartReturn failed: %v", err)
	}

	// Operation should now be RETURNING with a scheduled return
	if op.Phase != TradePhaseReturning {
		t.Fatalf("expected operation phase RETURNING, got %s", op.Phase)
	}
	if op.ReturnArriveAt <= 0 {
		t.Fatalf("expected ReturnArriveAt to be set")
	}

	// Offered resources and storage should be credited to receiver
	if receiver.Stats.Credits != receiverCreditsBefore+50 {
		t.Fatalf("expected receiver credits increased by 50, got %v", receiver.Stats.Credits)
	}
	if len(receiver.StorageItemsPresent) != 1 {
		t.Fatalf("expected offered storage transferred to receiver, got %d items", len(receiver.StorageItemsPresent))
	}
	if receiver.StorageItemsPresent[0].Prototype.ArtifactData == nil {
		t.Fatalf("expected full offered storage prototype data transferred to receiver")
	}

	// Offered army should be in receiver's present (not sender's deployed)
	foundInReceiver := false
	for _, p := range receiver.ArmiesPresent {
		if p.Prototype.ID == offeredProto.ID && p.Count == 2 {
			foundInReceiver = true
			break
		}
	}
	if !foundInReceiver {
		t.Fatalf("expected offered army (2 units) in receiver's present armies")
	}

	// Sender should only have transport + requested army in deployed
	senderDeployedByProto := make(map[int]int)
	for _, d := range sender.ArmiesDeployed {
		if d.OperationKind == OperationKindTrade && d.OperationID == op.ID {
			senderDeployedByProto[d.Prototype.ID] += d.Count
		}
	}
	if senderDeployedByProto[transportProto.ID] != 1 {
		t.Fatalf("expected 1 transport unit still deployed on sender, got %d", senderDeployedByProto[transportProto.ID])
	}
	if senderDeployedByProto[requestedProto.ID] != 1 {
		t.Fatalf("expected 1 requested army unit in sender's deployed for return trip, got %d", senderDeployedByProto[requestedProto.ID])
	}
	if senderDeployedByProto[offeredProto.ID] != 0 {
		t.Fatalf("expected offered army removed from sender's deployed, got %d", senderDeployedByProto[offeredProto.ID])
	}

	// Receiver should have no deployed units remaining
	for _, d := range receiver.ArmiesDeployed {
		if d.OperationKind == OperationKindTrade && d.OperationID == op.ID {
			t.Fatalf("expected receiver to have no deployed units for this operation after swap")
		}
	}
	if len(receiver.StorageItemsDeployed) != 0 {
		t.Fatalf("expected receiver to have no deployed storage for this operation after swap")
	}
}

func TestTradeOperationService_ProcessArrivalAndStartReturn_WrongPhase(t *testing.T) {
	SetTestNow(t, 71_000)

	sender := newBaseWithDefaults(711)
	sender.UserID = uuid.New()
	receiver := newBaseWithDefaults(712)
	receiver.UserID = uuid.New()

	offered, err := NewTradePayload(PriceModel{Credits: 1}, nil, nil)
	if err != nil {
		t.Fatalf("offered payload failed: %v", err)
	}
	requested, err := NewTradePayload(PriceModel{Credits: 1}, nil, nil)
	if err != nil {
		t.Fatalf("requested payload failed: %v", err)
	}

	op, err := NewTradeOperation(
		sender.UserID, sender.ID, receiver.UserID, receiver.ID,
		Vector2i{X: 1, Y: 1}, Vector2i{X: 2, Y: 2},
		offered, requested,
		[]MilitaryUnitSnap{{PrototypeID: 99, Count: 1, Capacity: 20, Speed: 100}},
		nil,
	)
	if err != nil {
		t.Fatalf("operation creation failed: %v", err)
	}
	op.ID = 3002
	op.Phase = TradePhaseOutbound

	svc := NewTradeOperationService(sender, receiver, op)
	if err := svc.ProcessArrivalAndStartReturn(); err == nil {
		t.Fatalf("expected error for non-ARRIVED phase")
	}
}

func TestTradeOperationService_FinalizeSenderAfterReturn_RestoresOfferedOnExpired(t *testing.T) {
	SetTestNow(t, 56_000)

	sender := newBaseWithDefaults(601)
	sender.UserID = uuid.New()
	receiver := newBaseWithDefaults(602)
	receiver.UserID = uuid.New()

	offered, err := NewTradePayload(PriceModel{Credits: 15, Iron: 5}, nil, nil)
	if err != nil {
		t.Fatalf("offered payload failed: %v", err)
	}
	requested, err := NewTradePayload(PriceModel{Credits: 1}, nil, nil)
	if err != nil {
		t.Fatalf("requested payload failed: %v", err)
	}

	op, err := NewTradeOperation(sender.UserID, sender.ID, receiver.UserID, receiver.ID, Vector2i{X: 1, Y: 1}, Vector2i{X: 2, Y: 3}, offered, requested, []MilitaryUnitSnap{{PrototypeID: 99, Count: 1, Capacity: 20, Speed: 120}}, nil)
	if err != nil {
		t.Fatalf("operation creation failed: %v", err)
	}
	op.ID = 2002
	op.Phase = TradePhaseCompleted
	op.Result = TradeResultExpired

	sender.Stats.Credits = 0
	sender.Stats.Iron = 0

	service := NewTradeOperationService(sender, receiver, op)
	if err := service.FinalizeSenderAfterReturn(); err != nil {
		t.Fatalf("finalize sender on expired failed: %v", err)
	}

	if sender.Stats.Credits != 15 || sender.Stats.Iron != 5 {
		t.Fatalf("expected offered resources restored, got credits=%v iron=%v", sender.Stats.Credits, sender.Stats.Iron)
	}
}
