package main

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"

	common "github.com/torbensky/adventofcode-common"
)

func main() {
	fmt.Printf("Part 1: %d\n", part1(common.OpenInputFile()))
	fmt.Printf("Part 2: %d\n", part2(common.OpenInputFile()))
}

func part1(reader io.Reader) int {
	var mask string
	mem := make(map[int]int)
	common.ScanLines(reader, func(line string) {
		if strings.Contains(line, "mask") {
			mask = strings.Fields(line)[2]
		} else {
			ma := processMem(line)
			newVal := applyMask(ma.value, mask)
			mem[ma.location] = newVal
		}
	})

	sum := 0
	for _, val := range mem {
		sum += val
	}

	return sum
}

func applyMask(val int, mask string) int {
	for i, c := range mask {
		switch c {
		case 'X':
			// ignore
		case '1':
			val = setBit(val, uint(35-i))
		case '0':
			val = clearBit(val, uint(35-i))
		}
	}
	return val
}

func clearBit(n int, i uint) int {
	n &^= 1 << i
	return n
}

func setBit(n int, i uint) int {
	n |= (1 << i)
	return n
}

type memAssign struct {
	location int
	value    int
}

var memAssignRegex = regexp.MustCompile(`mem\[(\d+)\]\s*=\s*(\d+)`)

func processMem(line string) memAssign {
	results := memAssignRegex.FindStringSubmatch(line)
	if len(results) < 3 {
		log.Fatalf("unexpected format for mem assign %s: %v\n", line, results)
	}

	l, err := strconv.Atoi(results[1])
	common.MustNotError(err)
	v, err := strconv.Atoi(results[2])
	common.MustNotError(err)
	return memAssign{
		location: l,
		value:    v,
	}
}

func part2(reader io.Reader) int {
	var mask string
	mem := make(map[int]int)
	common.ScanLines(reader, func(line string) {
		if strings.Contains(line, "mask") {
			mask = strings.Fields(line)[2]
		} else {
			ma := processMem(line)
			memUpdates := applyMask2(ma, mask)
			for _, ma := range memUpdates {
				mem[ma.location] = ma.value
			}
		}
	})

	sum := 0
	for _, val := range mem {
		sum += val
	}

	return sum
}

func applyMask2(val memAssign, mask string) []memAssign {
	var results []memAssign

	locations := []int{val.location}
	for i, c := range mask {
		switch c {
		case 'X':
			var newLocations []int
			for j := 0; j < len(locations); j++ {
				newLocations = append(newLocations, clearBit(locations[j], uint(35-i)))
				locations[j] = setBit(locations[j], uint(35-i))
			}
			locations = append(locations, newLocations...)
		case '1':
			for j := 0; j < len(locations); j++ {
				locations[j] = setBit(locations[j], uint(35-i))
			}
		case '0':
			// do nothing
		}
	}

	for _, l := range locations {
		results = append(results, memAssign{
			location: l,
			value:    val.value,
		})
	}

	return results
}
