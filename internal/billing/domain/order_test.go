package domain

import (
	"testing"

	"github.com/google/uuid"
)

func newTestPendingOrder() *PurchaseOrder {
	return NewPendingOrder(
		uuid.New(),
		uuid.New(),
		100,
		9900,
		"RUB",
		PaymentProviderYooKassa,
	)
}

func TestNewPendingOrder(t *testing.T) {
	SetTestNow(t, 1_000)

	userID := uuid.New()
	packageID := uuid.New()
	order := NewPendingOrder(userID, packageID, 100, 9900, "RUB", PaymentProviderYooKassa)

	if order.ID == uuid.Nil {
		t.Fatalf("expected a generated order ID, got nil")
	}
	if order.UserID != userID || order.PackageID != packageID {
		t.Fatalf("unexpected ids: user=%s package=%s", order.UserID, order.PackageID)
	}
	if order.Status != OrderStatusPending {
		t.Fatalf("expected status PENDING, got %s", order.Status)
	}
	if order.Crystals != 100 || order.AmountMinorUnits != 9900 || order.Currency != "RUB" {
		t.Fatalf("unexpected order amounts: %+v", order)
	}
	if order.CreatedAt != 1_000 || order.UpdatedAt != 1_000 {
		t.Fatalf("expected timestamps to use NowUnix (1000), got created=%d updated=%d", order.CreatedAt, order.UpdatedAt)
	}
	if events := order.PullEvents(); len(events) != 0 {
		t.Fatalf("expected no events on creation, got %d", len(events))
	}
}

func TestPurchaseOrder_MarkPaid(t *testing.T) {
	SetTestNow(t, 2_000)

	order := newTestPendingOrder()
	if err := order.MarkPaid(); err != nil {
		t.Fatalf("MarkPaid from PENDING returned error: %v", err)
	}
	if order.Status != OrderStatusPaid {
		t.Fatalf("expected status PAID, got %s", order.Status)
	}
	if order.UpdatedAt != 2_000 {
		t.Fatalf("expected UpdatedAt to use NowUnix (2000), got %d", order.UpdatedAt)
	}

	events := order.PullEvents()
	if len(events) != 1 {
		t.Fatalf("expected exactly one event, got %d", len(events))
	}
	paid, ok := events[0].(OrderPaidEvent)
	if !ok {
		t.Fatalf("expected OrderPaidEvent, got %T", events[0])
	}
	if paid.OrderID != order.ID || paid.UserID != order.UserID || paid.PackageID != order.PackageID || paid.Crystals != order.Crystals {
		t.Fatalf("OrderPaidEvent fields mismatch: %+v vs order %+v", paid, order)
	}
}

func TestPurchaseOrder_MarkPaid_Idempotent(t *testing.T) {
	order := newTestPendingOrder()
	if err := order.MarkPaid(); err != nil {
		t.Fatalf("first MarkPaid returned error: %v", err)
	}
	_ = order.PullEvents() // drain the event from the first transition

	if err := order.MarkPaid(); err != nil {
		t.Fatalf("idempotent MarkPaid returned error: %v", err)
	}
	if order.Status != OrderStatusPaid {
		t.Fatalf("expected status to remain PAID, got %s", order.Status)
	}
	if events := order.PullEvents(); len(events) != 0 {
		t.Fatalf("expected no second event on idempotent MarkPaid, got %d", len(events))
	}
}

func TestPurchaseOrder_MarkPaid_InvalidFromFailed(t *testing.T) {
	order := newTestPendingOrder()
	if err := order.MarkFailed(); err != nil {
		t.Fatalf("MarkFailed setup returned error: %v", err)
	}
	_ = order.PullEvents()

	if err := order.MarkPaid(); err == nil {
		t.Fatalf("expected error transitioning FAILED -> PAID, got nil")
	}
	if order.Status != OrderStatusFailed {
		t.Fatalf("expected status to remain FAILED, got %s", order.Status)
	}
	if events := order.PullEvents(); len(events) != 0 {
		t.Fatalf("expected no event on invalid transition, got %d", len(events))
	}
}

func TestPurchaseOrder_MarkFailed(t *testing.T) {
	SetTestNow(t, 3_000)

	order := newTestPendingOrder()
	if err := order.MarkFailed(); err != nil {
		t.Fatalf("MarkFailed from PENDING returned error: %v", err)
	}
	if order.Status != OrderStatusFailed {
		t.Fatalf("expected status FAILED, got %s", order.Status)
	}
	if order.UpdatedAt != 3_000 {
		t.Fatalf("expected UpdatedAt to use NowUnix (3000), got %d", order.UpdatedAt)
	}

	events := order.PullEvents()
	if len(events) != 1 {
		t.Fatalf("expected exactly one event, got %d", len(events))
	}
	failed, ok := events[0].(OrderFailedEvent)
	if !ok {
		t.Fatalf("expected OrderFailedEvent, got %T", events[0])
	}
	if failed.OrderID != order.ID {
		t.Fatalf("OrderFailedEvent OrderID mismatch: got %s want %s", failed.OrderID, order.ID)
	}
}

func TestPurchaseOrder_MarkFailed_Idempotent(t *testing.T) {
	order := newTestPendingOrder()
	if err := order.MarkFailed(); err != nil {
		t.Fatalf("first MarkFailed returned error: %v", err)
	}
	_ = order.PullEvents()

	if err := order.MarkFailed(); err != nil {
		t.Fatalf("idempotent MarkFailed returned error: %v", err)
	}
	if events := order.PullEvents(); len(events) != 0 {
		t.Fatalf("expected no second event on idempotent MarkFailed, got %d", len(events))
	}
}

func TestPurchaseOrder_MarkFailed_InvalidFromPaid(t *testing.T) {
	order := newTestPendingOrder()
	if err := order.MarkPaid(); err != nil {
		t.Fatalf("MarkPaid setup returned error: %v", err)
	}
	_ = order.PullEvents()

	if err := order.MarkFailed(); err == nil {
		t.Fatalf("expected error transitioning PAID -> FAILED, got nil")
	}
	if order.Status != OrderStatusPaid {
		t.Fatalf("expected status to remain PAID, got %s", order.Status)
	}
	if events := order.PullEvents(); len(events) != 0 {
		t.Fatalf("expected no event on invalid transition, got %d", len(events))
	}
}
