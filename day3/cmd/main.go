package main

import (
	"fmt"
	"log"

	common "github.com/torbensky/adventofcode-common"
)

// The rune that indicates a tree tile on the map
const tree = '#'

// Buffer to store each line of a file
var fileLines []string

// loads map data from the file and returns the map dimensions (width,height)
func loadMapData() (int, int) {
	file := common.OpenInputFile()
	defer file.Close()
	fileLines = common.ReadStringLines(file)
	return len(fileLines[0]), len(fileLines)
}

// gets the map tile at the specified coordinate
func getMapTile(x, y int) (rune, error) {
	// Validate y coordinate is in map bounds
	if y < 0 || y > len(fileLines) {
		return 0, fmt.Errorf("y coord outside of map bounds")
	}

	// Validate x coordinate is in map bounds
	if x < 0 || x > len(fileLines[y]) {
		return 0, fmt.Errorf("x coord is outside of map bounds")
	}

	// y is the line num, x is the rune
	return rune(fileLines[y][x]), nil
}

func main() {
	width, height := loadMapData()

	// Part 1

	fmt.Println("Part One")
	fmt.Printf("Total trees: %d\n", traverseSlope(3, 1, width, height))
	fmt.Println()

	// Part 2

	fmt.Println("Part Two")
	total := 1
	for _, val := range []struct {
		dx int
		dy int
	}{
		{dx: 1, dy: 1},
		{dx: 3, dy: 1},
		{dx: 5, dy: 1},
		{dx: 7, dy: 1},
		{dx: 1, dy: 2},
	} {
		next := traverseSlope(val.dx, val.dy, width, height)
		fmt.Printf("- Right %d, down %d = %d\n", val.dx, val.dy, next)
		total *= next
	}

	fmt.Printf("Total trees: %d\n", total)
}

// Fully traverses a "slope" across a map, counting the number of trees that are encountered
//
// A "slope" is defined by a horizontal (dx) and vertical (dy) distance that you move
//
// For example, dx=3,dy=1 means you move 3 tiles to the right, and 1 down (starting from the top left 0,0)
//
// "width" and "height" are the bounds of the map
//
// (0,0) is the top left, (width-1,height-1) is the bottom right
func traverseSlope(dx, dy, width, height int) int {
	treesEncountered := 0

	right := 0
	for down := dy; down < height; down += dy {
		right = (right + dx) % width // the slope wraps around horizontally (this had me stuck for a while!)
		tile, err := getMapTile(right, down)
		if err != nil {
			log.Fatal(err)
		}

		if tile == tree {
			treesEncountered++
		}
	}

	return treesEncountered
}
