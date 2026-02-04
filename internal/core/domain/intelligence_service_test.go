package domain

import (
	"testing"
)

func TestIntelligenceService_ResolveScanVisibility(t *testing.T) {
	service := NewIntelligenceService()

	t.Run("no defender always visible", func(t *testing.T) {
		if !service.ResolveScanVisibility(10, nil) {
			t.Error("expected true when defender is nil")
		}
	})

	t.Run("scan strength >= cloaking", func(t *testing.T) {
		defender := &UserBaseModel{}
		defender.BuildingsPresent = []BuildItemPresent{
			{Prototype: BuildItemPrototype{
				IntelligenceData: &IntelligenceBuildingData{
					Subtype:         IntelligenceSubtypeCloaking,
					StealthStrength: 50,
				},
			}},
		}

		if !service.ResolveScanVisibility(50, defender) {
			t.Error("expected true when scan strength equals cloaking")
		}
		if !service.ResolveScanVisibility(60, defender) {
			t.Error("expected true when scan strength exceeds cloaking")
		}
		if service.ResolveScanVisibility(40, defender) {
			t.Error("expected false when scan strength is less than cloaking")
		}
	})
}

func TestIntelligenceService_ResolveRadarDetection(t *testing.T) {
	service := NewIntelligenceService()

	t.Run("no radar building", func(t *testing.T) {
		base := &UserBaseModel{}
		op := &MilitaryOperation{Units: []MilitaryUnitSnap{}}
		if service.ResolveRadarDetection(base, op) {
			t.Error("expected false when no radar building present")
		}
	})

	t.Run("non-stealthy operation", func(t *testing.T) {
		base := &UserBaseModel{
			BuildingsPresent: []BuildItemPresent{
				{Prototype: BuildItemPrototype{
					IntelligenceData: &IntelligenceBuildingData{
						Subtype: IntelligenceSubtypeRadar,
					},
				}},
			},
		}
		op := &MilitaryOperation{Units: []MilitaryUnitSnap{{Stealth: 0, Count: 1}}}
		if !service.ResolveRadarDetection(base, op) {
			t.Error("expected true for non-stealthy operation with radar")
		}
	})

	t.Run("stealthy operation contest", func(t *testing.T) {
		base := &UserBaseModel{
			BuildingsPresent: []BuildItemPresent{
				{Prototype: BuildItemPrototype{
					IntelligenceData: &IntelligenceBuildingData{
						Subtype:         IntelligenceSubtypeRadar,
						StealthStrength: 100,
					},
				}},
			},
		}

		// radar 100 > stealth 50 -> detected
		op1 := &MilitaryOperation{
			Units:          []MilitaryUnitSnap{{Stealth: 50, Count: 1}},
			TotalModifiers: MilitaryModifiersFromSnaps(nil),
		}
		if !service.ResolveRadarDetection(base, op1) {
			t.Error("expected true when radar strength > op stealth")
		}

		// radar 100 <= stealth 100 -> not detected
		op2 := &MilitaryOperation{
			Units:          []MilitaryUnitSnap{{Stealth: 100, Count: 1}},
			TotalModifiers: MilitaryModifiersFromSnaps(nil),
		}
		if service.ResolveRadarDetection(base, op2) {
			t.Error("expected false when radar strength <= op stealth")
		}
	})
}

func TestIntelligenceService_TriangulateScanSource(t *testing.T) {
	service := NewIntelligenceService()
	attackerCoords := Vector2i{X: 100, Y: 100}
	defenderCoords := Vector2i{X: 120, Y: 120} // Distance is sqrt(20^2 + 20^2) = ~28.28

	t.Run("no interceptor", func(t *testing.T) {
		defender := &UserBaseModel{Coordinates: defenderCoords}
		info := service.TriangulateScanSource(attackerCoords, defender, true)

		if info.PossibleSource != nil {
			t.Error("expected nil PossibleSource without interceptor")
		}
		if info.UncertaintyRadius != 0 {
			t.Errorf("expected 0 uncertainty, got %d", info.UncertaintyRadius)
		}
	})

	t.Run("with interceptor", func(t *testing.T) {
		defender := &UserBaseModel{
			Coordinates: defenderCoords,
			BuildingsPresent: []BuildItemPresent{
				{Prototype: BuildItemPrototype{
					IntelligenceData: &IntelligenceBuildingData{
						Subtype:         IntelligenceSubtypeScanInterceptor,
						StealthStrength: 10,
					},
				}},
			},
		}
		// dist is ~28, intercept is 10, radius should be ~18
		info := service.TriangulateScanSource(attackerCoords, defender, true)

		if info.PossibleSource == nil {
			t.Fatal("expected non-nil PossibleSource")
		}
		if info.UncertaintyRadius <= 0 {
			t.Errorf("expected positive uncertainty radius, got %d", info.UncertaintyRadius)
		}

		// Verify source is within radius of attacker
		distToAttacker := info.PossibleSource.DistanceTo(attackerCoords)
		// Small buffer for potential float precision/rounding in distance calculation vs radius
		if distToAttacker > float64(info.UncertaintyRadius)+1.0 {
			t.Errorf("estimated source %v too far from attacker %v (dist %f, radius %d)",
				info.PossibleSource, attackerCoords, distToAttacker, info.UncertaintyRadius)
		}
	})

	t.Run("powerful interceptor pinpointing", func(t *testing.T) {
		defender := &UserBaseModel{
			Coordinates: defenderCoords,
			BuildingsPresent: []BuildItemPresent{
				{Prototype: BuildItemPrototype{
					IntelligenceData: &IntelligenceBuildingData{
						Subtype:         IntelligenceSubtypeScanInterceptor,
						StealthStrength: 100, // much larger than distance
					},
				}},
			},
		}

		info := service.TriangulateScanSource(attackerCoords, defender, true)
		if info.UncertaintyRadius != 0 {
			t.Errorf("expected 0 uncertainty with powerful interceptor, got %d", info.UncertaintyRadius)
		}
		if info.PossibleSource == nil || *info.PossibleSource != attackerCoords {
			t.Errorf("expected pinpointed source %v, got %v", attackerCoords, info.PossibleSource)
		}
	})
}
