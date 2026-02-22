package domain

import (
	"math/rand"
	"time"
)

// BasePlacementService provides logic for finding a free chunk for a new base.
type BasePlacementService struct {
	minBaseDistance int
	random          *rand.Rand
}

const defaultMinBaseDistance = 3

func NewBasePlacementService() *BasePlacementService {
	return NewBasePlacementServiceWithConfig(defaultMinBaseDistance, nil)
}

func NewBasePlacementServiceWithConfig(minBaseDistance int, random *rand.Rand) *BasePlacementService {
	if minBaseDistance < 0 {
		minBaseDistance = 0
	}
	if random == nil {
		random = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	return &BasePlacementService{minBaseDistance: minBaseDistance, random: random}
}

// FindFreeChunkForBase returns coordinates for a free sector for a new base.
// It searches outward from map center, enforces minimum Euclidean spacing from existing bases,
// and randomly breaks ties between equally-good candidates.
func (s *BasePlacementService) FindFreeChunkForBase(occupied []Vector2i) (int, int) {
	occupiedSet := make(map[[2]int]struct{}, len(occupied))
	bucketSize := s.minBaseDistance
	if bucketSize < 1 {
		bucketSize = 1
	}
	bucketRange := 0
	if s.minBaseDistance > 0 {
		bucketRange = (s.minBaseDistance + bucketSize - 1) / bucketSize
	}
	buckets := make(map[[2]int][]Vector2i)
	for _, c := range occupied {
		occupiedSet[[2]int{c.X, c.Y}] = struct{}{}
		bKey := [2]int{floorDiv(c.X, bucketSize), floorDiv(c.Y, bucketSize)}
		buckets[bKey] = append(buckets[bKey], c)
	}

	minDistSq := int64(s.minBaseDistance) * int64(s.minBaseDistance)
	for radius := 0; ; radius++ {
		bestDistSq := int64(-1)
		best := make([]Vector2i, 0, intMax(1, 8*radius))
		for x := -radius; x <= radius; x++ {
			for y := -radius; y <= radius; y++ {
				if intMax(intAbs(x), intAbs(y)) != radius {
					continue
				}
				if _, exists := occupiedSet[[2]int{x, y}]; exists {
					continue
				}
				candidate := Vector2i{X: x, Y: y}
				if minDistSq > 0 && !isFarEnoughFromBases(candidate, buckets, bucketSize, bucketRange, minDistSq) {
					continue
				}
				distSq := int64(x*x + y*y)
				if bestDistSq == -1 || distSq < bestDistSq {
					bestDistSq = distSq
					best = best[:0]
				}
				if distSq == bestDistSq {
					best = append(best, candidate)
				}
			}
		}
		if len(best) > 0 {
			picked := best[s.random.Intn(len(best))]
			return picked.X, picked.Y
		}
	}
}

func isFarEnoughFromBases(candidate Vector2i, buckets map[[2]int][]Vector2i, bucketSize int, bucketRange int, minDistSq int64) bool {
	bx := floorDiv(candidate.X, bucketSize)
	by := floorDiv(candidate.Y, bucketSize)
	for dx := -bucketRange; dx <= bucketRange; dx++ {
		for dy := -bucketRange; dy <= bucketRange; dy++ {
			for _, base := range buckets[[2]int{bx + dx, by + dy}] {
				dx64 := int64(candidate.X - base.X)
				dy64 := int64(candidate.Y - base.Y)
				if dx64*dx64+dy64*dy64 < minDistSq {
					return false
				}
			}
		}
	}
	return true
}

func floorDiv(a int, b int) int {
	if b <= 0 {
		return 0
	}
	if a >= 0 {
		return a / b
	}
	return -(((-a) + b - 1) / b)
}

func intMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func intAbs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
