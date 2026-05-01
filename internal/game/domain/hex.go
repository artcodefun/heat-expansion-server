package domain

// HexRing returns all hex coordinates at exactly dist steps from center,
// using pointy-top even-row offset coordinates.
func HexRing(center Vector2i, dist int) []Vector2i {
	if dist <= 0 {
		return []Vector2i{center}
	}
	centerCube := offsetToCube(center)
	directions := []hexCube{
		{1, -1, 0}, {1, 0, -1}, {0, 1, -1},
		{-1, 1, 0}, {-1, 0, 1}, {0, -1, 1},
	}
	start := hexCube{
		q: centerCube.q + directions[4].q*dist,
		r: centerCube.r + directions[4].r*dist,
		s: centerCube.s + directions[4].s*dist,
	}
	results := make([]Vector2i, 0, 6*dist)
	curr := start
	for i := 0; i < 6; i++ {
		for j := 0; j < dist; j++ {
			results = append(results, cubeToOffset(curr))
			curr = hexCube{
				q: curr.q + directions[i].q,
				r: curr.r + directions[i].r,
				s: curr.s + directions[i].s,
			}
		}
	}
	return results
}

type hexCube struct {
	q, r, s int
}

func offsetToCube(v Vector2i) hexCube {
	q := v.X - (v.Y-(v.Y&1))/2
	r := v.Y
	return hexCube{q: q, r: r, s: -q - r}
}

func cubeToOffset(c hexCube) Vector2i {
	x := c.q + (c.r-(c.r&1))/2
	y := c.r
	return Vector2i{X: x, Y: y}
}
