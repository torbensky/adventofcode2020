package main

import (
	"fmt"
	"io"
	"strings"

	common "github.com/torbensky/adventofcode-common"
)

type cell = byte
type coord []int

const (
	floor    cell = '.'
	empty    cell = 'L'
	occupied cell = '#'
)

var vectors = []coord{
	{1, 1},
	{1, -1},
	{-1, 1},
	{-1, -1},
	{1, 0},
	{-1, 0},
	{0, 1},
	{0, -1},
}

func main() {
	fmt.Printf("Part 1: %d\n", part1(common.OpenInputFile()))
	fmt.Printf("Part 2: %d\n", part2(common.OpenInputFile()))
}

func loadGrid(reader io.Reader) [][]cell {
	var grid [][]cell
	common.ScanLines(reader, func(line string) {
		// TODO: use a byte reader
		grid = append(grid, []cell(line))
	})
	return grid
}

func part1(reader io.Reader) int {
	grid := loadGrid(reader)
	for {
		stable := simulate(grid)
		if stable {
			break
		}
	}

	return countOccupied(grid)
}

func simulate(grid [][]cell) bool {
	gridChanged := false
	var swaps []coord
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			switch grid[y][x] {
			case floor:
				// nothing
			case empty:
				if countAdjacent(grid, y, x, occupied) == 0 {
					gridChanged = true
					swaps = append(swaps, []int{y, x})
				}
			case occupied:
				if countAdjacent(grid, y, x, occupied) >= 4 {
					gridChanged = true
					swaps = append(swaps, []int{y, x})
				}
			}
		}
	}

	for _, coord := range swaps {
		swapOccupancy(grid, coord[0], coord[1])
	}

	return !gridChanged
}

func countAdjacent(grid [][]cell, y, x int, cellType cell) int {
	matches := 0

	for _, v := range vectors {
		if !isValidCoord(grid, y+v[0], x+v[1]) {
			continue
		}

		c := grid[y+v[0]][x+v[1]]
		if c == cellType {
			matches++
		}
	}

	return matches
}

func countOccupied(grid [][]cell) int {
	count := 0
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			if grid[i][j] == occupied {
				count++
			}
		}
	}

	return count
}

func part2(reader io.Reader) int {
	grid := loadGrid(reader)
	for {
		stable := simulatePart2(grid)
		if stable {
			break
		}
	}
	return countOccupied(grid)
}

func simulatePart2(grid [][]cell) bool {
	gridChanged := false
	var swaps []coord
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			switch grid[y][x] {
			case floor:
				// nothing
			case empty:
				if countAdjacentPart2(grid, y, x, occupied) == 0 {
					gridChanged = true
					swaps = append(swaps, []int{y, x})
				}
			case occupied:
				if countAdjacentPart2(grid, y, x, occupied) >= 5 {
					gridChanged = true
					swaps = append(swaps, []int{y, x})
				}
			}
		}
	}

	for _, coord := range swaps {
		swapOccupancy(grid, coord[0], coord[1])
	}

	return !gridChanged
}

func countAdjacentPart2(grid [][]cell, y, x int, cellType cell) int {
	matches := 0

	for _, v := range vectors {
		if !isValidCoord(grid, y+v[0], x+v[1]) {
			continue
		}

		if doesVectorHit(grid, y, x, v, cellType) {
			matches++
		}
	}

	return matches
}

func doesVectorHit(grid [][]cell, y, x int, vector []int, cellType cell) bool {
	dy, dx := y, x
	for {

		// follow the vector!
		dy, dx = dy+vector[0], dx+vector[1]

		// Did vector reach end of the room?
		if !isValidCoord(grid, dy, dx) {
			return false
		}

		gt := grid[dy][dx]

		// Did it hit the thing we are looking for?
		if gt == cellType {
			return true
		}

		// Did it hit some other non-empty cell?
		if gt != floor {
			return false
		}
	}
}

func swapOccupancy(grid [][]cell, y, x int) {
	if grid[y][x] == empty {
		grid[y][x] = occupied
	} else {
		grid[y][x] = empty
	}
}

func isValidCoord(grid [][]cell, y, x int) bool {
	if y < 0 || x < 0 {
		return false
	}
	if y >= len(grid) || x >= len(grid[y]) {
		return false
	}

	return true
}

func printGrid(grid [][]cell) {
	fmt.Println(strings.Repeat("=", len(grid[0])))
	fmt.Println()
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			fmt.Printf(string(grid[y][x]))
		}
		fmt.Println()
	}
	fmt.Println()
	fmt.Println(strings.Repeat("=", len(grid[0])))
}
