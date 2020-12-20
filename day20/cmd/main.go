package main

import (
	"fmt"
	"io"

	common "github.com/torbensky/adventofcode-common"
	"github.com/torbensky/adventofcode2020/day20"
)

func main() {
	fmt.Printf("Part 1: %d\n", part1(common.OpenInputFile()))
	fmt.Printf("Part 2: %d\n", part2(common.OpenInputFile()))
}

func part1(reader io.Reader) int {
	tiles := day20.LoadTiles(reader)

	groups := tiles.GetTileGroups()

	result := 1
	for _, t := range groups.CornerTiles {
		result *= t.ID
	}

	return result
}

func part2(reader io.Reader) int {
	common.ScanLines(reader, func(line string) {
		// fmt.Println(line)
	})
	// TODO: implement me
	return -1
}
