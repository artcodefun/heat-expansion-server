package domain

// BasePlacementService provides logic for finding a free chunk for a new base.
type BasePlacementService struct{}

func NewBasePlacementService() *BasePlacementService {
	return &BasePlacementService{}
}

// FindFreeChunkForBase returns coordinates for a free chunk given occupied coordinates.
func (s *BasePlacementService) FindFreeChunkForBase(occupied []Vector2i) (int, int) {
	if len(occupied) == 0 {
		return 0, 0
	}
	minX, minY := occupied[0].X, occupied[0].Y
	maxX, maxY := minX, minY
	for _, c := range occupied {
		if c.X < minX {
			minX = c.X
		}
		if c.X > maxX {
			maxX = c.X
		}
		if c.Y < minY {
			minY = c.Y
		}
		if c.Y > maxY {
			maxY = c.Y
		}
	}
	chunkSize := 10
	chunkCounts := make(map[[2]int]int)
	for _, coord := range occupied {
		chunkX := (coord.X - minX) / chunkSize
		chunkY := (coord.Y - minY) / chunkSize
		chunkCounts[[2]int{chunkX, chunkY}]++
	}
	centerChunkX := (0 - minX) / chunkSize
	centerChunkY := (0 - minY) / chunkSize
	minCount := -1
	var targetChunk [2]int
	for cx := 0; cx <= (maxX-minX)/chunkSize; cx++ {
		for cy := 0; cy <= (maxY-minY)/chunkSize; cy++ {
			count := chunkCounts[[2]int{cx, cy}]
			centerDist := abs(cx-centerChunkX) + abs(cy-centerChunkY)
			if minCount == -1 || count < minCount || (count == minCount && centerDist < abs(targetChunk[0]-centerChunkX)+abs(targetChunk[1]-centerChunkY)) {
				minCount = count
				targetChunk = [2]int{cx, cy}
			}
		}
	}
	occupiedSet := make(map[[2]int]struct{}, len(occupied))
	for _, c := range occupied {
		occupiedSet[[2]int{c.X, c.Y}] = struct{}{}
	}
	for x := minX + targetChunk[0]*chunkSize; x < minX+(targetChunk[0]+1)*chunkSize; x++ {
		for y := minY + targetChunk[1]*chunkSize; y < minY+(targetChunk[1]+1)*chunkSize; y++ {
			if _, exists := occupiedSet[[2]int{x, y}]; !exists {
				return x, y
			}
		}
	}
	return 0, 0
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
