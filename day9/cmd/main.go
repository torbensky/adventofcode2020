package main

import (
	"fmt"
	"io"
	"strconv"

	common "github.com/torbensky/adventofcode-common"
)

func main() {
	nums := loadNumbers(common.OpenInputFile())

	fmt.Printf("Part 1 - Result %d\n\n", part1(nums, 25))
	fmt.Printf("Part 2 - Result %d\n", part2(nums, 14144619))
}

func part2(nums []int, value int) int {
	for i := 0; i < len(nums); i++ {
		sum := 0
		smallest := nums[i]
		largest := nums[i]
		for j := i; j < len(nums); j++ {
			if smallest > nums[j] {
				smallest = nums[j]
			}
			if largest < nums[j] {
				largest = nums[j]
			}
			sum += nums[j]
			if sum == value {
				return smallest + largest
			}
			if sum > value {
				break
			}
		}
	}
	return -1
}

func part1(nums []int, preamble int) int {
	for i := preamble; i < len(nums); i++ {
		if !hasPastNThatSum(nums[i-preamble:i], preamble, nums[i]) {
			return nums[i]
		}
	}

	return -1
}

func hasPastNThatSum(nums []int, n int, target int) bool {
	for j := 0; j < n; j++ {
		for q := 0; q < n; q++ {
			// Don't allow adding to self
			if q == j {
				continue
			}

			if target == nums[j]+nums[q] {
				return true
			}
		}
	}
	return false
}

func loadNumbers(reader io.Reader) []int {
	var nums []int
	parseLine := func(line string) {
		val, err := strconv.Atoi(line)
		common.MustNotError(err)
		nums = append(nums, val)
	}
	common.ScanLines(reader, parseLine)

	return nums
}
