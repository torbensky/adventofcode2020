package main

import (
	"fmt"
	"io"
	"sort"
	"strconv"

	common "github.com/torbensky/adventofcode-common"
)

func main() {
	fmt.Printf("Part 1: %d\n", part1(common.OpenInputFile()))
	fmt.Printf("Part 2: %d\n", part2(common.OpenInputFile()))
}

func part1(reader io.Reader) int {
	adapters := loadData(reader)

	currentJoltage := 0
	num1Diff := 0
	num3Diff := 1 // the last jump to our device is always a 3 so start at 1
	for _, a := range adapters {
		diff := a - currentJoltage
		currentJoltage = a
		if diff == 1 {
			num1Diff++
		}
		if diff == 3 {
			num3Diff++
		}
	}

	return num1Diff * num3Diff
}

func part2(reader io.Reader) int {
	adapters := loadData(reader)
	return countPermutations(adapters, 0)
}

// return sorted list
func loadData(reader io.Reader) []int {
	var adapters []int
	common.ScanLines(reader, func(line string) {
		num, err := strconv.Atoi(line)
		common.MustNotError(err)
		adapters = append(adapters, num)
	})
	sort.Ints(adapters)
	return adapters
}

var cache = make(map[int]int)

func countPermutations(a []int, last int) int {
	if len(a) == 0 {
		return 1
	}

	if v, ok := cache[last]; ok {
		return v
	}

	count := 0
	max := len(a)
	if max > 3 {
		max = 3
	}

	for i := 0; i < max; i++ {
		diff := a[i] - last
		if diff > 3 {
			break
		}
		count += countPermutations(a[i+1:], a[i])
	}

	cache[last] = count
	return count
}
