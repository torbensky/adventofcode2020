package main

import (
	"fmt"
	"io"

	"github.com/torbensky/adventofcode2020/common"
)

func main() {
	fmt.Printf("Part 1: %d\n", part1(common.OpenInputFile()))
	fmt.Printf("Part 2: %d\n", part2(common.OpenInputFile()))
}

func part1(reader io.Reader) int {
	common.ScanLines(reader, func(line string) {
		fmt.Println(line)
	})
	// TODO: implement me
	return -1
}

func part2(reader io.Reader) int {
	common.ScanLines(reader, func(line string) {
		fmt.Println(line)
	})
	// TODO: implement me
	return -1
}
