package domain

import (
	"math/rand"
	"testing"
)

func TestBasePlacementService_EmptyWorldReturnsCenter(t *testing.T) {
	svc := NewBasePlacementServiceWithConfig(3, rand.New(rand.NewSource(1)))

	x, y := svc.FindFreeChunkForBase(nil)

	if x != 0 || y != 0 {
		t.Fatalf("expected first base at center (0,0), got (%d,%d)", x, y)
	}
}

func TestBasePlacementService_RespectsMinEuclideanDistance(t *testing.T) {
	svc := NewBasePlacementServiceWithConfig(3, rand.New(rand.NewSource(1)))
	occupied := []Vector2i{{X: 0, Y: 0}}

	x, y := svc.FindFreeChunkForBase(occupied)

	dx := x - occupied[0].X
	dy := y - occupied[0].Y
	if dx*dx+dy*dy < 9 {
		t.Fatalf("expected distance >= 3 from occupied base, got (%d,%d)", x, y)
	}
}

func TestBasePlacementService_PrefersCenterWhenPossible(t *testing.T) {
	svc := NewBasePlacementServiceWithConfig(3, rand.New(rand.NewSource(1)))
	occupied := []Vector2i{{X: 50, Y: 50}, {X: -50, Y: -50}}

	x, y := svc.FindFreeChunkForBase(occupied)

	if x != 0 || y != 0 {
		t.Fatalf("expected center-biased placement at (0,0), got (%d,%d)", x, y)
	}
}

func TestBasePlacementService_ExpandsOutsideDenseCenter(t *testing.T) {
	svc := NewBasePlacementServiceWithConfig(1, rand.New(rand.NewSource(1)))
	occupiedSet := make(map[[2]int]struct{})
	occupied := make([]Vector2i, 0)
	for x := -2; x <= 2; x++ {
		for y := -2; y <= 2; y++ {
			occupied = append(occupied, Vector2i{X: x, Y: y})
			occupiedSet[[2]int{x, y}] = struct{}{}
		}
	}

	x, y := svc.FindFreeChunkForBase(occupied)

	if _, exists := occupiedSet[[2]int{x, y}]; exists {
		t.Fatalf("expected a free coordinate, got occupied (%d,%d)", x, y)
	}
	if max(x, -x) <= 2 && max(y, -y) <= 2 {
		t.Fatalf("expected placement outside dense center square, got (%d,%d)", x, y)
	}
}

func TestBasePlacementService_RandomizesEqualCandidates(t *testing.T) {
	occupied := []Vector2i{{X: 0, Y: 0}}
	results := make(map[[2]int]struct{})

	for seed := int64(1); seed <= 16; seed++ {
		svc := NewBasePlacementServiceWithConfig(3, rand.New(rand.NewSource(seed)))
		x, y := svc.FindFreeChunkForBase(occupied)
		results[[2]int{x, y}] = struct{}{}
	}

	if len(results) < 2 {
		t.Fatalf("expected randomized tie-breaking to produce multiple coordinates, got %d", len(results))
	}
}

func TestBasePlacementService_SimulateWorldPopulation25(t *testing.T) {
	svc := NewBasePlacementServiceWithConfig(3, rand.New(rand.NewSource(42)))
	occupied := make([]Vector2i, 0, 25)
	seen := make(map[[2]int]struct{}, 25)

	for i := 0; i < 25; i++ {
		x, y := svc.FindFreeChunkForBase(occupied)
		key := [2]int{x, y}
		if _, exists := seen[key]; exists {
			t.Fatalf("duplicate coordinate generated at step %d: (%d,%d)", i+1, x, y)
		}
		seen[key] = struct{}{}
		occupied = append(occupied, Vector2i{X: x, Y: y})
		t.Logf("%02d: (%d,%d)", i+1, x, y)
	}
}
