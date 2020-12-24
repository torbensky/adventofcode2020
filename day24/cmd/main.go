package main

import (
	"fmt"
	"io"

	common "github.com/torbensky/adventofcode-common"
)

type hexPos struct {
	x int
	y int
	z int
}

type hexDirection int

const (
	northWest hexDirection = iota
	northEast
	east
	southEast
	southWest
	west
)

func (h hexPos) Neighbor(dir hexDirection) hexPos {
	switch dir {
	case northWest:
		return hexPos{
			x: h.x,
			y: h.y + 1,
			z: h.z - 1,
		}
	case northEast:
		return hexPos{
			x: h.x + 1,
			y: h.y,
			z: h.z - 1,
		}
	case east:
		return hexPos{
			x: h.x + 1,
			y: h.y - 1,
			z: h.z,
		}
	case southEast:
		return hexPos{
			x: h.x,
			y: h.y - 1,
			z: h.z + 1,
		}
	case southWest:
		return hexPos{
			x: h.x - 1,
			y: h.y,
			z: h.z + 1,
		}
	case west:
		return hexPos{
			x: h.x - 1,
			y: h.y + 1,
			z: h.z,
		}
	}

	panic("unexpected direction")
}

func (h hexPos) NorthEast() hexPos {
	return h.Neighbor(northEast)
}
func (h hexPos) East() hexPos {
	return h.Neighbor(east)
}
func (h hexPos) SouthEast() hexPos {
	return h.Neighbor(southEast)
}
func (h hexPos) SouthWest() hexPos {
	return h.Neighbor(southWest)
}
func (h hexPos) West() hexPos {
	return h.Neighbor(west)
}
func (h hexPos) NorthWest() hexPos {
	return h.Neighbor(northWest)
}

type hexGrid map[hexPos]bool

func (hg hexGrid) countNeighbors(pos hexPos, wantBlack bool, limit int) int {
	total := 0
	for _, cur := range pos.buildRing(1) {
		black, ok := hg[cur]
		if !ok {
			black = false
		}

		if black == wantBlack {
			total++
		}

		// short-circuit out if we hit the limit
		if total >= limit {
			return total
		}
	}
	return total
}

func (h hexPos) buildRing(radius int) []hexPos {
	cur := h
	for i := 0; i < radius; i++ {
		cur = cur.NorthWest()
	}

	ring := []hexPos{}
	for _, dir := range []hexDirection{southWest, southEast, east, northEast, northWest, west} {
		for j := 0; j < radius; j++ {
			cur = cur.Neighbor(dir)
			ring = append(ring, cur)
		}
	}

	return ring
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func (hg hexGrid) cycleDay() hexGrid {
	newGrid := make(hexGrid)

	seen := make(map[hexPos]struct{})

	checkHex := func(pos hexPos, isBlack bool) {
		if _, ok := seen[pos]; ok {
			return
		}
		if isBlack {
			// Any black tile with zero or more than 2 black tiles immediately adjacent to it is flipped to white.
			if count := hg.countNeighbors(pos, true, 6); count == 0 || count > 2 {
				newGrid[pos] = false
			} else {
				newGrid[pos] = true
			}
		} else {
			// Any white tile with exactly 2 black tiles immediately adjacent to it is flipped to black.
			if hg.countNeighbors(pos, true, 3) != 2 {
				newGrid[pos] = false
			} else {
				newGrid[pos] = true
			}
		}
	}

	processNeighbors := func(pos hexPos) {
		for _, cur := range pos.buildRing(1) {
			black, ok := hg[cur]
			if !ok {
				black = false
			}

			checkHex(cur, black)
			seen[cur] = struct{}{}
		}
	}

	// Find the largest ring we need to search
	for pos, isBlack := range hg {
		checkHex(pos, isBlack)
		seen[pos] = struct{}{}
		processNeighbors(pos)
	}

	return newGrid
}

func (hg hexGrid) countBlack() int {
	var total int
	for _, isBlack := range hg {
		if isBlack {
			total++
		}
	}
	return total
}

func (hg hexGrid) FollowInstructions(instructions string) {
	cur := hexPos{0, 0, 0}
	for i := 0; i < len(instructions); i++ {
		switch instructions[i] {
		case 'e':
			cur = cur.East()
		case 'w':
			cur = cur.West()
		case 'n':
			switch instructions[i+1] {
			case 'e':
				cur = cur.NorthEast()
				i = i + 1
			case 'w':
				cur = cur.NorthWest()
				i = i + 1
			default:
				panic("unexpected state")
			}
		case 's':
			switch instructions[i+1] {
			case 'e':
				cur = cur.SouthEast()
				i = i + 1
			case 'w':
				cur = cur.SouthWest()
				i = i + 1
			default:
				panic("unexpected state")
			}
		default:
			panic("unexpected input encountered")
		}
	}
	black, ok := hg[cur]
	if !ok {
		hg[cur] = true
	} else {
		hg[cur] = !black
	}
}

func newGrid() hexGrid {
	return hexGrid{}
}

func main() {
	fmt.Printf("Part 1: %d\n", part1(common.OpenInputFile()))
	fmt.Printf("Part 2: %d\n", part2(common.OpenInputFile()))
}

func part1(reader io.Reader) int {
	hg := newGrid()
	common.ScanLines(reader, func(line string) {
		hg.FollowInstructions(line)
	})

	return hg.countBlack()
}

func part2(reader io.Reader) int {
	hg := newGrid()
	common.ScanLines(reader, func(line string) {
		hg.FollowInstructions(line)
	})

	for i := 0; i < 100; i++ {
		hg = hg.cycleDay()
	}

	return hg.countBlack()
}
