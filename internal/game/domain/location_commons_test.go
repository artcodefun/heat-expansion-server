package domain

import (
	"testing"
)

func TestAppropriateLocationDefense(t *testing.T) {
	cases := []struct {
		name        string
		stats       UserBaseStats
		locType     LocationType
		expectedMin float64
		expectedMax float64
	}{
		{
			name:        "fresh player, resourceful — floor applies",
			stats:       UserBaseStats{MaxSpace: DefaultMaxSpace},
			locType:     LocationTypeResourceful,
			expectedMin: baseResourcefulDefense,
			expectedMax: baseResourcefulDefense,
		},
		{
			name:        "fresh player, dangerous — floor applies",
			stats:       UserBaseStats{MaxSpace: DefaultMaxSpace},
			locType:     LocationTypeDangerous,
			expectedMin: baseDangerousDefense,
			expectedMax: baseDangerousDefense,
		},
		{
			name:        "space progression only, no army (MaxSpace=300)",
			stats:       UserBaseStats{MaxSpace: 300},
			locType:     LocationTypeResourceful,
			expectedMin: 30.0, // baseResourcefulDefense * (300/100)
			expectedMax: 30.0,
		},
		{
			name:        "army factor dominant over space (small base, strong army)",
			stats:       UserBaseStats{MaxSpace: DefaultMaxSpace, Attack: 1000, Defence: 1000},
			locType:     LocationTypeResourceful,
			expectedMin: 100.0, // baseResourcefulDefense * (1000 / spawnStartingPower)
			expectedMax: 100.0,
		},
		{
			name:    "space floor beats weak army (army deployed away)",
			stats:   UserBaseStats{MaxSpace: 500, Attack: 50, Defence: 50},
			locType: LocationTypeResourceful,
			// spaceFactor=5, armyFactor=0.5 → space wins → 10*5=50
			expectedMin: 50.0,
			expectedMax: 50.0,
		},
		{
			name:    "mid-game player (MaxSpace=500, decent army)",
			stats:   UserBaseStats{MaxSpace: 500, Attack: 5000, Defence: 4000},
			locType: LocationTypeResourceful,
			// armyFactor = 5000/100 = 50, spaceFactor = 5 → army wins → 10*50=500
			expectedMin: 500.0,
			expectedMax: 500.0,
		},
		{
			name:    "production player, army home (MaxSpace=2800, strong army)",
			stats:   UserBaseStats{MaxSpace: 2800, Attack: 50000, Defence: 40000},
			locType: LocationTypeResourceful,
			// armyFactor = 50000/100 = 500, spaceFactor = 28 → army wins → 10*500=5000
			expectedMin: 5000.0,
			expectedMax: 5000.0,
		},
		{
			name:    "production player, army deployed (MaxSpace=2800, no army home)",
			stats:   UserBaseStats{MaxSpace: 2800},
			locType: LocationTypeResourceful,
			// spaceFactor = 28, armyFactor = 0 → space floor → 10*28=280
			expectedMin: 280.0,
			expectedMax: 280.0,
		},
		{
			name:    "production player, dangerous location",
			stats:   UserBaseStats{MaxSpace: 2800, Attack: 50000, Defence: 40000},
			locType: LocationTypeDangerous,
			// armyFactor = 500 → 40*500=20000
			expectedMin: 20000.0,
			expectedMax: 20000.0,
		},
		{
			name:        "unknown location type returns zero",
			stats:       UserBaseStats{MaxSpace: DefaultMaxSpace, Attack: 9999, Defence: 9999},
			locType:     LocationTypeEmpty,
			expectedMin: 0,
			expectedMax: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := AppropriateLocationDefense(tc.stats, tc.locType)
			if got < tc.expectedMin || got > tc.expectedMax {
				t.Errorf("AppropriateLocationDefense = %.2f, want [%.2f, %.2f]", got, tc.expectedMin, tc.expectedMax)
			}
		})
	}
}

func TestFillDefenders_EmptyOnZeroTarget(t *testing.T) {
	proto := &ArmyItemPrototype{ID: 1, Faction: FactionMarauders, Defence: 10}
	var armies []ArmyStack
	var structures []DefenseStack

	FillDefenders(&armies, &structures, FactionMarauders, 0, []*ArmyItemPrototype{proto}, nil)

	if len(armies) != 0 || len(structures) != 0 {
		t.Errorf("expected no defenders for targetPower=0, got %d armies %d structures", len(armies), len(structures))
	}
}

func TestFillDefenders_EmptyOnNoFactionUnits(t *testing.T) {
	// Prototype belongs to a different faction
	proto := &ArmyItemPrototype{ID: 1, Faction: FactionFerrousSwarm, Defence: 10}
	var armies []ArmyStack
	var structures []DefenseStack

	FillDefenders(&armies, &structures, FactionMarauders, 100, []*ArmyItemPrototype{proto}, nil)

	if len(armies) != 0 || len(structures) != 0 {
		t.Errorf("expected no defenders when no units match faction")
	}
}

func TestFillDefenders_ReachesTargetPower(t *testing.T) {
	proto := &ArmyItemPrototype{ID: 1, Faction: FactionMarauders, Defence: 10}
	var armies []ArmyStack
	var structures []DefenseStack

	const target = 50.0
	FillDefenders(&armies, &structures, FactionMarauders, target, []*ArmyItemPrototype{proto}, nil)

	total := 0
	for _, s := range armies {
		total += s.Prototype.Defence * s.Count
	}
	if float64(total) < target {
		t.Errorf("combined defence power = %d, want >= %.0f", total, target)
	}
}

func TestFillDefenders_ExceedsOldIterationCap(t *testing.T) {
	// Each unit has defence=1, so reaching targetPower=200 requires >100 units —
	// more than the old maxIterations=100 cap would have allowed.
	proto := &ArmyItemPrototype{ID: 1, Faction: FactionMarauders, Defence: 1}
	var armies []ArmyStack
	var structures []DefenseStack

	const target = 200.0
	FillDefenders(&armies, &structures, FactionMarauders, target, []*ArmyItemPrototype{proto}, nil)

	total := 0
	for _, s := range armies {
		total += s.Prototype.Defence * s.Count
	}
	if float64(total) < target {
		t.Errorf("combined defence power = %d, want >= %.0f (old cap would have stopped at 100)", total, target)
	}
}
