package domain

import (
	"math"
	"math/rand"
)

type IntelligenceService struct{}

func NewIntelligenceService() *IntelligenceService {
	return &IntelligenceService{}
}

// ResolveScanVisibility determines if an attacker's scan strength is enough to see through a defender's cloaking.
func (s *IntelligenceService) ResolveScanVisibility(scanStrength int, defender *UserBaseModel) bool {
	if defender == nil {
		return true
	}
	return scanStrength >= defender.TotalCloakingStealthStrength()
}

// ResolveRadarDetection determines if a defender's radar strength is enough to see an incoming operation.
func (s *IntelligenceService) ResolveRadarDetection(base *UserBaseModel, op *MilitaryOperation) bool {
	// First, check if base has any functional radar
	hasRadar := false
	for _, b := range base.BuildingsPresent {
		if b.Prototype.IntelligenceData != nil && b.Prototype.IntelligenceData.Subtype == IntelligenceSubtypeRadar {
			hasRadar = true
			break
		}
	}
	if !hasRadar {
		return false
	}

	// If op is not stealthy, it's always detected by radar
	if op.TotalStealth() <= 0 {
		return true
	}

	// For stealthy ops, we need more radar strength than total op stealth
	return base.TotalRadarStealthStrength() > op.TotalStealth()
}

// TriangulateScanSource generates the detection info and estimated source coordinates for a defender.
func (s *IntelligenceService) TriangulateScanSource(attackerCoords Vector2i, defender *UserBaseModel, scanPenetrated bool) ScanInterceptInfo {
	interceptPower := defender.TotalInterceptionStealthStrength()

	info := ScanInterceptInfo{
		ScannedCoordinates:     defender.Coordinates,
		ScanPenetratedCloaking: scanPenetrated,
	}

	if interceptPower <= 0 {
		return info
	}

	dist := defender.Coordinates.DistanceTo(attackerCoords)
	// Uncertainty shrinks as interceptPower grows relative to the distance.
	radius := int(dist) - interceptPower
	if radius < 0 {
		radius = 0
	}

	estimatedSource := s.randomSectorInCircle(attackerCoords, radius)
	info.PossibleSource = &estimatedSource
	info.UncertaintyRadius = radius

	return info
}

func (s *IntelligenceService) randomSectorInCircle(center Vector2i, radius int) Vector2i {
	if radius <= 0 {
		return center
	}
	angle := rand.Float64() * 2 * math.Pi
	r := math.Sqrt(rand.Float64()) * float64(radius)
	return Vector2i{
		X: center.X + int(math.Round(r*math.Cos(angle))),
		Y: center.Y + int(math.Round(r*math.Sin(angle))),
	}
}
