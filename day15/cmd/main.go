package main

import (
	"fmt"
)

var testInput = []int{}

func main() {
	fmt.Printf("Part 1: %d\n", part1(testInput))
	fmt.Printf("Part 2: %d\n", part2(testInput))
}

func part1(vals []int) int {
	return playSayGame(vals, 2020)
}

func part2(vals []int) int {
	return playSayGame(vals, 30000000)
}

func playSayGame(vals []int, numTurns int) int {
	initialLength := len(vals)
	spoken := make(map[int]int)
	for i, s := range vals[:initialLength-1] {
		spoken[s] = i + 1
	}
	last := vals[initialLength-1]
	for turn := initialLength; turn < numTurns; turn++ {
		next := 0
		if val, ok := spoken[last]; ok {
			next = turn - val
		}
		spoken[last] = turn
		last = next
	}
	return last
}
