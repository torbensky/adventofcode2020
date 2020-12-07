package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/torbensky/adventofcode2020/common"
)

func main() {

	// Scan in numeric data
	var viableValues []int
	common.ScanFile(common.GetInputFilePath(), func(line string) bool {
		val, err := strconv.Atoi(line)
		if err != nil {
			log.Printf("unable to parse integer - skipping line '%s'\n", line)
			return true
		}

		// Don't consider values > 2020
		if val > 2020 {
			return true
		}
		viableValues = append(viableValues, val)
		return true
	}, nil)

	sort.Ints(viableValues) // sort required for solution algorithms

	// Find the pair that sums to 2020
	val1, val2, err := solvePart1(viableValues)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Part 1 Solution: %d * %d = %d\n", val1, val2, val1*val2)

	val1, val2, val3, err := solvePart2(viableValues)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Part 1 Solution: %d * %d * %d = %d\n", val1, val2, val3, val1*val2*val3)
}

// solvePart2 finds 3 numbers in a sorted list that sum to 2020
// NOTE: expects an already sorted integer list to work properly
func solvePart2(values []int) (int, int, int, error) {
	var sum int
	for idx, val1 := range values {
		for i := idx; i < len(values); i++ {
			for j := i; j < len(values); j++ {
				sum = val1 + values[i] + values[j]
				if sum == 2020 {
					return val1, values[i], values[j], nil
				}
				if sum > 2020 {
					break
				}
			}
		}
	}

	return 0, 0, 0, fmt.Errorf("no solution")
}

// solvePart1 finds 2 numbers in a sorted list that sum to 2020
// NOTE: expects an already sorted integer list to work properly
func solvePart1(values []int) (int, int, error) {
	for idx, val1 := range values {
		for i := idx; i < len(values); i++ {
			if val1+values[i] == 2020 {
				return val1, values[i], nil
			}

			if val1+values[i] > 2020 {
				break
			}
		}
	}

	return 0, 0, fmt.Errorf("no solution")
}
