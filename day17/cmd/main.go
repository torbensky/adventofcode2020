package main

import (
	"fmt"
	"io"

	common "github.com/torbensky/adventofcode-common"
)

type coord3d struct {
	X int
	Y int
	Z int
}

func (c coord3d) text() string {
	return fmt.Sprintf("[x:%d,y:%d,z:%d]", c.X, c.Y, c.Z)
}

func add(c1, c2 coord3d) coord3d {
	return coord3d{
		X: c1.X + c2.X,
		Y: c1.Y + c2.Y,
		Z: c1.Z + c2.Z,
	}
}

func make3dBaseVectors() []coord3d {
	components := [3]int{-1, 0, 1}
	var baseVectors []coord3d
	for _, x := range components {
		for _, y := range components {
			for _, z := range components {
				if x == 0 && y == 0 && z == 0 {
					continue
				}

				baseVectors = append(baseVectors, coord3d{X: x, Y: y, Z: z})
			}
		}
	}
	return baseVectors
}

var vectors3d = make3dBaseVectors()

type space3d map[coord3d]struct{}

// counts up to "limit" number of neighbors in space
func (s space3d) countNeighbors(coord coord3d, limit int) int {
	active := 0
	for _, vec := range vectors3d {
		neighborPos := add(coord, vec)

		if _, ok := s[neighborPos]; ok {
			active++

			// stop early if we hit the limit
			if active >= limit {
				return active
			}
		}
	}

	return active
}

func (s space3d) cycle() space3d {
	// Make a copy so we can switch all nodes "at once"
	nextSpace := make(space3d)

	// The only things we need to know about in the space are the active nodes
	for cube := range s {

		// Check if the active cube remains active
		if count := s.countNeighbors(cube, 4); count == 2 || count == 3 {
			nextSpace[cube] = struct{}{}
		}

		// Next, we can check every neighbour of the active cube to see if it is inactive
		for _, vec := range vectors3d {

			neighbor := add(cube, vec)

			// inactive nodes aren't in the space
			if _, ok := s[neighbor]; !ok {
				// inactive neighbors will activate if there are exactly 3 active neighbors
				if count := s.countNeighbors(neighbor, 4); count == 3 {
					nextSpace[neighbor] = struct{}{}
				}
			}
		}
	}

	return nextSpace
}

func main() {
	fmt.Printf("Part 1: %d\n", part1(common.OpenInputFile()))
	fmt.Printf("Part 2: %d\n", part2(common.OpenInputFile()))
}

func load3dSpace(reader io.Reader) space3d {
	s := space3d{}
	y := 0
	common.ScanLines(reader, func(line string) {

		for x, b := range line {
			// Only need to store active nodes, assume all other coords inactive
			if b == '#' {
				s[coord3d{X: x, Y: y, Z: 0}] = struct{}{}
			}
		}

		y++
	})

	return s
}

func part1(reader io.Reader) int {
	space := load3dSpace(reader)
	for i := 0; i < 6; i++ {
		space = space.cycle()
	}
	// All nodes mapped are active
	return len(space)
}

func part2(reader io.Reader) int {
	space := load4dSpace(reader)
	for i := 0; i < 6; i++ {
		space = space.cycle()
	}
	// All nodes mapped are active
	return len(space)
}

// TODO: not sure how best to share a lot of this common code in Go without generics... too tired to bother

type coord4d struct {
	X int
	Y int
	Z int
	W int
}

func (c coord4d) text() string {
	return fmt.Sprintf("[x:%d,y:%d,z:%d,w:%d]", c.X, c.Y, c.Z, c.W)
}

func add4d(c1, c2 coord4d) coord4d {
	return coord4d{
		X: c1.X + c2.X,
		Y: c1.Y + c2.Y,
		Z: c1.Z + c2.Z,
		W: c1.W + c2.W,
	}
}

func make4dBaseVectors() []coord4d {
	components := [3]int{-1, 0, 1}
	var baseVectors []coord4d
	for _, x := range components {
		for _, y := range components {
			for _, z := range components {
				for _, w := range components {
					if x == 0 && y == 0 && z == 0 && w == 0 {
						continue
					}

					baseVectors = append(baseVectors, coord4d{X: x, Y: y, Z: z, W: w})
				}
			}
		}
	}
	return baseVectors
}

var vectors4d = make4dBaseVectors()

type space4d map[coord4d]struct{}

// counts up to "limit" number of neighbors in space
func (s space4d) countNeighbors(coord coord4d, limit int) int {
	active := 0
	for _, vec := range vectors4d {
		neighborPos := add4d(coord, vec)

		if _, ok := s[neighborPos]; ok {
			active++

			// stop early if we hit the limit
			if active >= limit {
				return active
			}
		}
	}

	return active
}

func (s space4d) cycle() space4d {
	// Make a copy so we can switch all nodes "at once"
	nextSpace := make(space4d)

	// The only things we need to know about in the space are the active nodes
	for cube := range s {

		// Check if the active cube remains active
		if count := s.countNeighbors(cube, 4); count == 2 || count == 3 {
			nextSpace[cube] = struct{}{}
		}

		// Next, we can check every neighbour of the active cube to see if it is inactive
		for _, vec := range vectors4d {

			neighbor := add4d(cube, vec)

			// inactive nodes aren't in the space
			if _, ok := s[neighbor]; !ok {
				// inactive neighbors will activate if there are exactly 3 active neighbors
				if count := s.countNeighbors(neighbor, 4); count == 3 {
					nextSpace[neighbor] = struct{}{}
				}
			}
		}
	}

	return nextSpace
}

func load4dSpace(reader io.Reader) space4d {
	s := space4d{}
	y := 0
	common.ScanLines(reader, func(line string) {

		for x, b := range line {
			// Only need to store active nodes, assume all other coords inactive
			if b == '#' {
				s[coord4d{X: x, Y: y, Z: 0}] = struct{}{}
			}
		}

		y++
	})

	return s
}
