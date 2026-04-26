package domain

import (
	"testing"

	"github.com/google/uuid"
)

func hasTradeReturnArrivedEvent(events []DomainEvent) bool {
	for _, e := range events {
		if _, ok := e.(TradeOperationReturnArrivedEvent); ok {
			return true
		}
	}
	return false
}

func hasTradeCreatedEventWithID(events []DomainEvent, operationID int) bool {
	for _, e := range events {
		ev, ok := e.(TradeOperationCreatedEvent)
		if ok && ev.OperationID == operationID {
			return true
		}
	}
	return false
}

func makeTradePayloadForTests(t *testing.T) TradePayload {
	t.Helper()
	p, err := NewTradePayload(PriceModel{Credits: 50}, nil, []TradeArmyItemSnap{{PrototypeID: 1, Count: 1, Capacity: 2}})
	if err != nil {
		t.Fatalf("payload setup failed: %v", err)
	}
	return p
}

func makeTransportForTests() []MilitaryUnitSnap {
	return []MilitaryUnitSnap{{PrototypeID: 1, Count: 1, Speed: 100, Capacity: 10, Attack: 1, Defence: 1}}
}

func TestTradeOperationLifecycle_HappyPath(t *testing.T) {
	SetTestNow(t, 1_000)
	offered := makeTradePayloadForTests(t)
	requested := makeTradePayloadForTests(t)

	op, err := NewTradeOperation(
		uuid.New(),
		1,
		uuid.New(),
		2,
		Vector2i{X: 1, Y: 1},
		Vector2i{X: 4, Y: 5},
		offered,
		requested,
		makeTransportForTests(),
		nil,
	)
	if err != nil {
		t.Fatalf("new operation failed: %v", err)
	}
	if op.Phase != TradePhasePending {
		t.Fatalf("expected pending phase, got %s", op.Phase)
	}
	if op.Result != TradeResultUnknown {
		t.Fatalf("expected unknown result on creation, got %s", op.Result)
	}
	if op.ExpiresAt != 1_000+24*60*60 {
		t.Fatalf("expected expires_at to be now+24h, got %d", op.ExpiresAt)
	}

	if err := op.Accept(); err != nil {
		t.Fatalf("accept failed: %v", err)
	}
	if op.Phase != TradePhaseOutbound {
		t.Fatalf("expected outbound, got %s", op.Phase)
	}

	if err := op.OnArrive(); err != nil {
		t.Fatalf("on arrive failed: %v", err)
	}
	if op.Phase != TradePhaseArrived {
		t.Fatalf("expected arrived, got %s", op.Phase)
	}

	if err := op.StartReturn(); err != nil {
		t.Fatalf("start return failed: %v", err)
	}
	if op.Phase != TradePhaseReturning {
		t.Fatalf("expected returning, got %s", op.Phase)
	}

	if err := op.OnReturnArrive(); err != nil {
		t.Fatalf("on return arrive failed: %v", err)
	}
	if op.Phase != TradePhaseCompleted {
		t.Fatalf("expected completed, got %s", op.Phase)
	}
	if op.Result != TradeResultSuccess {
		t.Fatalf("expected success result, got %s", op.Result)
	}
}

func TestTradeOperation_RequiresUnitsOnBothLegs(t *testing.T) {
	SetTestNow(t, 900)

	sender := uuid.New()
	receiver := uuid.New()

	resourceOnly, err := NewTradePayload(PriceModel{Credits: 50}, nil, nil)
	if err != nil {
		t.Fatalf("resource-only payload should be valid: %v", err)
	}
	armyOnly, err := NewTradePayload(PriceModel{}, nil, []TradeArmyItemSnap{{PrototypeID: 10, Count: 1, Capacity: 0}})
	if err != nil {
		t.Fatalf("army-only payload should be valid: %v", err)
	}

	_, err = NewTradeOperation(
		sender,
		1,
		receiver,
		2,
		Vector2i{X: 1, Y: 1},
		Vector2i{X: 4, Y: 5},
		resourceOnly,
		resourceOnly,
		nil,
		nil,
	)
	if err == nil {
		t.Fatalf("expected empty convoy trade to be rejected")
	}

	_, err = NewTradeOperation(
		sender,
		1,
		receiver,
		2,
		Vector2i{X: 1, Y: 1},
		Vector2i{X: 4, Y: 5},
		armyOnly,
		resourceOnly,
		nil,
		nil,
	)
	if err == nil {
		t.Fatalf("expected return leg without units to be rejected")
	}

	_, err = NewTradeOperation(
		sender,
		1,
		receiver,
		2,
		Vector2i{X: 1, Y: 1},
		Vector2i{X: 4, Y: 5},
		resourceOnly,
		armyOnly,
		nil,
		nil,
	)
	if err == nil {
		t.Fatalf("expected outbound leg without units to be rejected")
	}
}

func TestTradeOperation_EmitCreatedEvent_UsesAssignedID(t *testing.T) {
	SetTestNow(t, 1_500)
	offered := makeTradePayloadForTests(t)
	requested := makeTradePayloadForTests(t)

	op, err := NewTradeOperation(
		uuid.New(),
		1,
		uuid.New(),
		2,
		Vector2i{X: 1, Y: 1},
		Vector2i{X: 4, Y: 5},
		offered,
		requested,
		makeTransportForTests(),
		nil,
	)
	if err != nil {
		t.Fatalf("new operation failed: %v", err)
	}

	op.ID = 42
	op.EmitCreatedEvent()
	if !hasTradeCreatedEventWithID(op.PullEvents(), 42) {
		t.Fatalf("expected created event with assigned operation id")
	}
}

func TestTradeOperation_ExpireIfPending_Idempotent(t *testing.T) {
	SetTestNow(t, 2_000)
	offered := makeTradePayloadForTests(t)
	requested := makeTradePayloadForTests(t)

	op, err := NewTradeOperation(
		uuid.New(),
		1,
		uuid.New(),
		2,
		Vector2i{X: 1, Y: 1},
		Vector2i{X: 2, Y: 3},
		offered,
		requested,
		makeTransportForTests(),
		nil,
	)
	if err != nil {
		t.Fatalf("new operation failed: %v", err)
	}

	if !op.ExpireIfPending() {
		t.Fatalf("expected expire transition true")
	}
	if op.Phase != TradePhaseCompleted {
		t.Fatalf("expected completed phase after expiration, got %s", op.Phase)
	}
	if op.Result != TradeResultExpired {
		t.Fatalf("expected expired result after expiration, got %s", op.Result)
	}
	if op.ExpireIfPending() {
		t.Fatalf("expected second expiration attempt to be no-op")
	}
}

func TestTradeOperation_GuardsInvalidPhase(t *testing.T) {
	SetTestNow(t, 3_000)
	offered := makeTradePayloadForTests(t)
	requested := makeTradePayloadForTests(t)

	op, err := NewTradeOperation(
		uuid.New(),
		1,
		uuid.New(),
		2,
		Vector2i{X: 1, Y: 1},
		Vector2i{X: 2, Y: 2},
		offered,
		requested,
		makeTransportForTests(),
		nil,
	)
	if err != nil {
		t.Fatalf("new operation failed: %v", err)
	}

	if err := op.OnArrive(); err == nil {
		t.Fatalf("expected on arrive to fail from pending")
	}
	if err := op.OnReturnArrive(); err == nil {
		t.Fatalf("expected on return arrive to fail from pending")
	}
	if err := op.CancelByInitiator(); err != nil {
		t.Fatalf("cancel in pending should succeed: %v", err)
	}
	if op.Phase != TradePhaseCompleted {
		t.Fatalf("expected completed phase after pending cancel, got %s", op.Phase)
	}
	if op.Result != TradeResultCanceled {
		t.Fatalf("expected canceled result after pending cancel, got %s", op.Result)
	}
	if err := op.Accept(); err == nil {
		t.Fatalf("expected accept to fail after cancellation")
	}
}

func TestTradeOperation_CancelOutbound_StartsReturnAndFinishesCancelled(t *testing.T) {
	SetTestNow(t, 6_000)
	offered := makeTradePayloadForTests(t)
	requested := makeTradePayloadForTests(t)

	op, err := NewTradeOperation(
		uuid.New(),
		1,
		uuid.New(),
		2,
		Vector2i{X: 1, Y: 1},
		Vector2i{X: 4, Y: 5},
		offered,
		requested,
		makeTransportForTests(),
		nil,
	)
	if err != nil {
		t.Fatalf("new operation failed: %v", err)
	}

	if err := op.Accept(); err != nil {
		t.Fatalf("accept failed: %v", err)
	}

	SetTestNow(t, 6_010)
	if err := op.CancelByInitiator(); err != nil {
		t.Fatalf("cancel in outbound should succeed: %v", err)
	}
	if op.Phase != TradePhaseReturning {
		t.Fatalf("expected returning after outbound cancel, got %s", op.Phase)
	}
	if op.Result != TradeResultCanceled {
		t.Fatalf("expected canceled result after outbound cancel, got %s", op.Result)
	}
	if op.ReturnArriveAt <= op.ReturnDepartAt {
		t.Fatalf("expected positive return travel after outbound cancel")
	}

	SetTestNow(t, op.ReturnArriveAt)
	if err := op.OnReturnArrive(); err != nil {
		t.Fatalf("on return arrive failed: %v", err)
	}
	if op.Phase != TradePhaseCompleted {
		t.Fatalf("expected completed terminal phase after outbound cancel return, got %s", op.Phase)
	}
	if op.Result != TradeResultCanceled {
		t.Fatalf("expected canceled terminal result after outbound cancel return, got %s", op.Result)
	}
}

func TestTradeOperation_Decline_CompletesWithDeclinedResult(t *testing.T) {
	SetTestNow(t, 7_000)
	offered := makeTradePayloadForTests(t)
	requested := makeTradePayloadForTests(t)

	op, err := NewTradeOperation(
		uuid.New(),
		1,
		uuid.New(),
		2,
		Vector2i{X: 1, Y: 1},
		Vector2i{X: 3, Y: 3},
		offered,
		requested,
		makeTransportForTests(),
		nil,
	)
	if err != nil {
		t.Fatalf("new operation failed: %v", err)
	}

	if err := op.Decline(); err != nil {
		t.Fatalf("decline failed: %v", err)
	}
	if op.Phase != TradePhaseCompleted {
		t.Fatalf("expected completed phase after decline, got %s", op.Phase)
	}
	if op.Result != TradeResultDeclined {
		t.Fatalf("expected declined result after decline, got %s", op.Result)
	}
	if !hasTradeReturnArrivedEvent(op.PullEvents()) {
		t.Fatalf("expected return-arrived event after decline")
	}
}

func TestTradeOperation_ExpireIfPending_EmitsReturnArrivedEvent(t *testing.T) {
	SetTestNow(t, 8_000)
	offered := makeTradePayloadForTests(t)
	requested := makeTradePayloadForTests(t)

	op, err := NewTradeOperation(
		uuid.New(),
		1,
		uuid.New(),
		2,
		Vector2i{X: 1, Y: 1},
		Vector2i{X: 2, Y: 2},
		offered,
		requested,
		makeTransportForTests(),
		nil,
	)
	if err != nil {
		t.Fatalf("new operation failed: %v", err)
	}

	if !op.ExpireIfPending() {
		t.Fatalf("expected expire transition")
	}
	if !hasTradeReturnArrivedEvent(op.PullEvents()) {
		t.Fatalf("expected return-arrived event after expiration")
	}
}

func TestTradeOperation_UpdatePhaseBasedOnTime_OutboundToArrived(t *testing.T) {
	SetTestNow(t, 8_500)
	offered := makeTradePayloadForTests(t)
	requested := makeTradePayloadForTests(t)

	op, err := NewTradeOperation(
		uuid.New(),
		1,
		uuid.New(),
		2,
		Vector2i{X: 1, Y: 1},
		Vector2i{X: 3, Y: 3},
		offered,
		requested,
		makeTransportForTests(),
		nil,
	)
	if err != nil {
		t.Fatalf("new operation failed: %v", err)
	}

	if err := op.Accept(); err != nil {
		t.Fatalf("accept failed: %v", err)
	}

	SetTestNow(t, op.OutboundArriveAt)
	op.UpdatePhaseBasedOnTime()
	if op.Phase != TradePhaseArrived {
		t.Fatalf("expected arrived phase after outbound arrival time, got %s", op.Phase)
	}
}

func TestTradeOperation_UpdatePhaseBasedOnTime_ReturningToCompleted(t *testing.T) {
	SetTestNow(t, 9_000)
	offered := makeTradePayloadForTests(t)
	requested := makeTradePayloadForTests(t)

	op, err := NewTradeOperation(
		uuid.New(),
		1,
		uuid.New(),
		2,
		Vector2i{X: 1, Y: 1},
		Vector2i{X: 3, Y: 3},
		offered,
		requested,
		makeTransportForTests(),
		nil,
	)
	if err != nil {
		t.Fatalf("new operation failed: %v", err)
	}

	if err := op.Accept(); err != nil {
		t.Fatalf("accept failed: %v", err)
	}
	if err := op.OnArrive(); err != nil {
		t.Fatalf("on arrive failed: %v", err)
	}
	if err := op.StartReturn(); err != nil {
		t.Fatalf("start return failed: %v", err)
	}

	SetTestNow(t, op.ReturnArriveAt)
	op.UpdatePhaseBasedOnTime()
	if op.Phase != TradePhaseCompleted {
		t.Fatalf("expected completed phase after return arrival time, got %s", op.Phase)
	}
	if op.Result != TradeResultSuccess {
		t.Fatalf("expected success result on completion, got %s", op.Result)
	}
}

func TestTradeOperation_New_ValidatesOutboundLegCapacityWithPayloadAndTransport(t *testing.T) {
	SetTestNow(t, 4_000)

	// outbound required = 20 (200 credits), provided = 5 (payload army) + 10 (transport) = 15
	offered, err := NewTradePayload(PriceModel{Credits: 200}, nil, []TradeArmyItemSnap{{PrototypeID: 1, Count: 1, Capacity: 5}})
	if err != nil {
		t.Fatalf("offered payload failed: %v", err)
	}
	requested, err := NewTradePayload(PriceModel{Credits: 0}, nil, []TradeArmyItemSnap{{PrototypeID: 4, Count: 1, Capacity: 1}})
	if err != nil {
		t.Fatalf("requested payload failed: %v", err)
	}

	_, err = NewTradeOperation(
		uuid.New(),
		1,
		uuid.New(),
		2,
		Vector2i{X: 1, Y: 1},
		Vector2i{X: 2, Y: 2},
		offered,
		requested,
		[]MilitaryUnitSnap{{PrototypeID: 2, Count: 1, Capacity: 10, Speed: 100}},
		nil,
	)
	if err == nil {
		t.Fatalf("expected outbound capacity validation error")
	}
}

func TestTradeOperation_New_ValidatesReturnLegCapacityWithPayloadAndTransport(t *testing.T) {
	SetTestNow(t, 5_000)

	offered, err := NewTradePayload(PriceModel{Credits: 0}, nil, []TradeArmyItemSnap{{PrototypeID: 5, Count: 1, Capacity: 1}})
	if err != nil {
		t.Fatalf("offered payload failed: %v", err)
	}
	// return required = 15 (150 credits), provided = 2 (payload army) + 10 (transport) = 12
	requested, err := NewTradePayload(PriceModel{Credits: 150}, nil, []TradeArmyItemSnap{{PrototypeID: 3, Count: 1, Capacity: 2}})
	if err != nil {
		t.Fatalf("requested payload failed: %v", err)
	}

	_, err = NewTradeOperation(
		uuid.New(),
		1,
		uuid.New(),
		2,
		Vector2i{X: 1, Y: 1},
		Vector2i{X: 2, Y: 2},
		offered,
		requested,
		[]MilitaryUnitSnap{{PrototypeID: 2, Count: 1, Capacity: 10, Speed: 100}},
		nil,
	)
	if err == nil {
		t.Fatalf("expected return capacity validation error")
	}
}
