package domain

import (
	"testing"
)

func TestTech_StartAndSpeedUp_EmitsEvents(t *testing.T) {
	SetTestNow(t, 10_000)
	base := newBaseWithDefaults(3)
	tech := &TechItemPrototype{
		ID:           200,
		Name:         "Improved Mining",
		Category:     TechCategoryBuild,
		Price:        PriceModel{Credits: 250},
		ResearchTime: 90,
	}
	if err := base.StartTechResearch(tech); err != nil {
		t.Fatalf("StartTechResearch error: %v", err)
	}
	events := base.PullEvents()
	if len(events) == 0 {
		t.Fatalf("expected TechResearchStartedEvent on start")
	}
	_, ok := events[0].(TechResearchStartedEvent)
	if !ok {
		t.Fatalf("expected first event to be TechResearchStartedEvent, got %T", events[0])
	}
	if got := len(base.TechnologiesInProgress); got != 1 {
		t.Fatalf("expected 1 tech in progress, got %d", got)
	}
	// resources should be debited by tech price
	if base.Stats.Credits != 1000-250 {
		t.Fatalf("unexpected credits after StartTechResearch: %+v", base.Stats)
	}

	// speed up
	inProgID := base.TechnologiesInProgress[0].BaseOwnedItem.ID
	base.PullEvents()
	if err := base.SpeedUpTechResearch(inProgID); err != nil {
		t.Fatalf("SpeedUpTechResearch error: %v", err)
	}
	events = base.PullEvents()
	var gotFinished, gotSpeedup bool
	for _, e := range events {
		switch e.(type) {
		case TechResearchFinishedEvent:
			gotFinished = true
		case TechResearchSpeedupEvent:
			gotSpeedup = true
		}
	}
	if !gotFinished || !gotSpeedup {
		t.Fatalf("expected tech finished and speedup events, got finished=%v speedup=%v", gotFinished, gotSpeedup)
	}
	if len(base.TechnologiesInProgress) != 0 {
		t.Fatalf("expected no technologies in progress after speedup")
	}
	if len(base.TechnologiesDone) != 1 || base.TechnologiesDone[0].Prototype.ID != tech.ID {
		t.Fatalf("expected tech to be moved to done after speedup, got %+v", base.TechnologiesDone)
	}
}

func TestTech_HasTechAndGetTechLevel(t *testing.T) {
	base := newBaseWithDefaults(1)
	techID := 123
	if base.HasTech(techID) {
		t.Errorf("expected HasTech to be false initially")
	}
	if level := base.GetTechLevel(techID); level != 0 {
		t.Errorf("expected level 0, got %d", level)
	}

	base.TechnologiesDone = append(base.TechnologiesDone, TechItemDone{
		Prototype: TechItemPrototype{ID: techID},
		Level:     2,
	})

	if !base.HasTech(techID) {
		t.Errorf("expected HasTech to be true after adding tech")
	}
	if level := base.GetTechLevel(techID); level != 2 {
		t.Errorf("expected level 2, got %d", level)
	}
}

func TestTech_MoveTechQueue_CompletesResearch(t *testing.T) {
	SetTestNow(t, 10_000)
	base := newBaseWithDefaults(1)
	techProto := TechItemPrototype{
		ID:           300,
		ResearchTime: 100,
	}

	// Manually add to in-progress
	base.TechnologiesInProgress = append(base.TechnologiesInProgress, TechItemInProgress{
		BaseOwnedItem:  NewBaseOwnedItem(base.ID),
		Prototype:      techProto,
		StartDate:      10_000,
		CompletionDate: 10_100,
	})

	// Before time
	base.MoveTechQueue()
	if len(base.TechnologiesInProgress) != 1 {
		t.Fatalf("expected 1 in progress before completion time")
	}

	// After time
	SetTestNow(t, 10_101)
	base.MoveTechQueue()

	if len(base.TechnologiesInProgress) != 0 {
		t.Fatalf("expected 0 in progress after completion time")
	}
	if len(base.TechnologiesDone) != 1 || base.TechnologiesDone[0].Prototype.ID != 300 {
		t.Fatalf("expected tech 300 in done")
	}
}

func TestTech_StartTechResearch_NotAvailableWhenAlreadyInProgress(t *testing.T) {
	SetTestNow(t, 11_000)
	base := newBaseWithDefaults(5)
	tech := &TechItemPrototype{
		ID:           201,
		Name:         "Shielding",
		Category:     TechCategoryBuild,
		Price:        PriceModel{Credits: 100},
		ResearchTime: 30,
	}
	// mark tech as already in progress
	base.TechnologiesInProgress = []TechItemInProgress{{
		BaseOwnedItem:  NewBaseOwnedItem(base.ID),
		Prototype:      *tech,
		StartDate:      NowUnix(),
		CompletionDate: NowUnix() + 30,
	}}

	if err := base.StartTechResearch(tech); err == nil {
		t.Fatalf("expected error when starting research for a tech already in progress")
	}
	if len(base.TechnologiesInProgress) != 1 {
		t.Fatalf("expected in-progress list to remain unchanged")
	}
}

func TestTech_StartTechResearch_ScalesPriceAndTimeByHalfPerLevel(t *testing.T) {
	proto := &TechItemPrototype{
		ID:           999,
		Category:     TechCategoryBase,
		Price:        PriceModel{Credits: 1000, Iron: 200, Titanium: 40, Antimatter: 2},
		ResearchTime: 900,
		Improvement:  &TechImprovement{Type: ImprovementTypeSpaceCapacity, Value: 50},
	}

	cases := []struct {
		level            int
		wantCredits      float64
		wantIron         float64
		wantTitanium     float64
		wantAntimatter   float64
		wantResearchTime int64
	}{
		// level 0→1: multiplier 1.0
		{0, 1000, 200, 40, 2, 900},
		// level 1→2: multiplier 1.5 (floor truncation)
		{1, 1500, 300, 60, 3, 1350},
		// level 2→3: multiplier 2.0
		{2, 2000, 400, 80, 4, 1800},
	}

	for _, tc := range cases {
		SetTestNow(t, 10_000)
		base := newBaseWithDefaults(1)
		base.Stats.Credits = 10_000
		base.Stats.Iron = 10_000
		base.Stats.Titanium = 10_000
		base.Stats.Antimatter = 10
		if tc.level > 0 {
			base.TechnologiesDone = []TechItemDone{{
				BaseOwnedItem: NewBaseOwnedItem(base.ID),
				Prototype:     *proto,
				Level:         tc.level,
			}}
		}

		creditsBefore := base.Stats.Credits
		ironBefore := base.Stats.Iron
		titaniumBefore := base.Stats.Titanium
		antimatterBefore := base.Stats.Antimatter

		if err := base.StartTechResearch(proto); err != nil {
			t.Fatalf("level %d: StartTechResearch error: %v", tc.level, err)
		}

		if got := creditsBefore - base.Stats.Credits; got != tc.wantCredits {
			t.Errorf("level %d: credits deducted = %.0f, want %.0f", tc.level, got, tc.wantCredits)
		}
		if got := ironBefore - base.Stats.Iron; got != tc.wantIron {
			t.Errorf("level %d: iron deducted = %.0f, want %.0f", tc.level, got, tc.wantIron)
		}
		if got := titaniumBefore - base.Stats.Titanium; got != tc.wantTitanium {
			t.Errorf("level %d: titanium deducted = %.0f, want %.0f", tc.level, got, tc.wantTitanium)
		}
		if got := antimatterBefore - base.Stats.Antimatter; got != tc.wantAntimatter {
			t.Errorf("level %d: antimatter deducted = %.0f, want %.0f", tc.level, got, tc.wantAntimatter)
		}
		if got := base.TechnologiesInProgress[0].CompletionDate - base.TechnologiesInProgress[0].StartDate; got != tc.wantResearchTime {
			t.Errorf("level %d: research time = %d, want %d", tc.level, got, tc.wantResearchTime)
		}
	}
}

func TestTech_StartTechResearch_NotAvailableWhenAlreadyDone(t *testing.T) {
	SetTestNow(t, 11_000)
	base := newBaseWithDefaults(5)
	tech := &TechItemPrototype{
		ID:           201,
		Name:         "Shielding",
		Category:     TechCategoryBuild,
		Price:        PriceModel{Credits: 100},
		ResearchTime: 30,
	}
	// mark tech as already done (Level 1) so AvailableTechnologies returns empty
	base.TechnologiesDone = []TechItemDone{{
		BaseOwnedItem: NewBaseOwnedItem(base.ID),
		Prototype:     *tech,
		Level:         1,
	}}

	if err := base.StartTechResearch(tech); err == nil {
		t.Fatalf("expected error when starting research for an already done tech")
	}
	if len(base.TechnologiesInProgress) != 0 {
		t.Fatalf("expected no technologies in progress when start fails")
	}
	if events := base.PullEvents(); len(events) != 0 {
		t.Fatalf("expected no events when StartTechResearch fails, got %v", events)
	}
}
