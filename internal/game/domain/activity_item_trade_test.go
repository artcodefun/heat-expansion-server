package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewActivityFromTradeOperation(t *testing.T) {
	SetTestNow(t, 1_000)

	userID := uuid.New()
	item := NewActivityFromTradeOperation(userID, 12, 99)

	if item.Kind != ActivityKindTrade {
		t.Fatalf("expected trade activity kind, got %s", item.Kind)
	}
	if item.BaseID != 12 {
		t.Fatalf("expected base 12, got %d", item.BaseID)
	}
	if item.Trade == nil || item.Trade.OpID != 99 {
		t.Fatalf("expected trade payload with op 99, got %+v", item.Trade)
	}
	if item.CreatedAt != 1_000 {
		t.Fatalf("expected created at 1000, got %d", item.CreatedAt)
	}
}

func TestNewTradeAlert(t *testing.T) {
	SetTestNow(t, 1_000)

	op := &TradeOperation{
		SenderUserID:   uuid.New(),
		SenderBaseID:   12,
		ReceiverUserID: uuid.New(),
		ReceiverBaseID: 34,
	}
	alert := NewTradeAlert(op, false, TradeAlertKindCreated)
	if alert.Kind != AlertKindTrade {
		t.Fatalf("expected trade alert, got %s", alert.Kind)
	}
	if alert.Title != "alert.trade.created.title" {
		t.Fatalf("unexpected title key: %s", alert.Title)
	}
	if alert.Content != "alert.trade.created.content" {
		t.Fatalf("unexpected content key: %s", alert.Content)
	}
	if alert.UserID != op.ReceiverUserID {
		t.Fatalf("expected receiver alert user, got %s", alert.UserID)
	}
	if alert.BaseID == nil || *alert.BaseID != op.ReceiverBaseID {
		t.Fatalf("expected receiver base %d, got %+v", op.ReceiverBaseID, alert.BaseID)
	}
}
