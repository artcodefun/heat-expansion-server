package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestTradePayload_New_Valid(t *testing.T) {
	itemID := uuid.New()
	p, err := NewTradePayload(
		PriceModel{Credits: 100},
		[]TradeStorageItemSnap{{ItemID: itemID, PrototypeID: 10, Category: StorageCategoryArtifact}},
		[]TradeArmyItemSnap{{PrototypeID: 20, Count: 3}},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Resources.Credits != 100 || len(p.Storage) != 1 || len(p.Army) != 1 {
		t.Fatalf("unexpected payload: %#v", p)
	}
}

func TestTradePayload_New_RejectsEmpty(t *testing.T) {
	_, err := NewTradePayload(PriceModel{}, nil, nil)
	if err == nil {
		t.Fatalf("expected empty payload error")
	}
}

func TestTradePayload_New_RejectsDuplicateStorageItem(t *testing.T) {
	id := uuid.New()
	_, err := NewTradePayload(
		PriceModel{Credits: 1},
		[]TradeStorageItemSnap{
			{ItemID: id, PrototypeID: 1, Category: StorageCategoryBuff},
			{ItemID: id, PrototypeID: 1, Category: StorageCategoryBuff},
		},
		nil,
	)
	if err == nil {
		t.Fatalf("expected duplicate storage item error")
	}
}

func TestTradePayload_CapacityHelpers(t *testing.T) {
	p, err := NewTradePayload(
		PriceModel{Credits: 100, Iron: 10},
		nil,
		[]TradeArmyItemSnap{{PrototypeID: 1, Count: 2, Capacity: 7}},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p.ProvidedArmyCapacity() != 14 {
		t.Fatalf("expected provided army capacity 14, got %v", p.ProvidedArmyCapacity())
	}

	expectedRequired := p.Resources.CreditsWorth() / WorthCapacityMultiplier
	if p.RequiredResourceCapacity() != expectedRequired {
		t.Fatalf("expected required capacity %v, got %v", expectedRequired, p.RequiredResourceCapacity())
	}
}
