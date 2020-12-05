package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

// represents a numeric range [lower,upper]
type numRange struct {
	lower int
	upper int
}

// splits a numeric rane into two halves
func splitRange(r numRange) (numRange, numRange) {
	newRangeSize := (r.upper - r.lower) / 2

	return numRange{
			lower: r.lower,
			upper: r.lower + newRangeSize,
		}, numRange{
			lower: r.lower + newRangeSize + 1,
			upper: r.upper,
		}
}

// A decoded boarding pass with row/column for seating
type boardingPass struct {
	row    int
	column int
	seatId int
}

func loadData(path string) []*boardingPass {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var passes []*boardingPass
	for scanner.Scan() {
		rowRange := numRange{lower: 0, upper: 127}
		colRange := numRange{lower: 0, upper: 7}

		line := scanner.Text()
		for i, c := range line {
			if i < 7 {
				lower, upper := splitRange(rowRange)
				switch c {
				case 'F':
					rowRange = lower
				case 'B':
					rowRange = upper
				}
			} else {
				lower, upper := splitRange(colRange)
				switch c {
				case 'L':
					colRange = lower
				case 'R':
					colRange = upper
				}
			}
		}

		passes = append(passes, &boardingPass{
			row:    rowRange.lower,
			column: colRange.lower,
			seatId: rowRange.lower*8 + colRange.lower,
		})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return passes
}

func main() {

	// Validate program usage
	if len(os.Args) != 2 {
		log.Fatal("This command accepts only one argument: the path to the input file")
	}
	passes := loadData(os.Args[1])

	var seatIds []int
	for _, p := range passes {
		id := p.seatId
		seatIds = append(seatIds, id)
	}
	sort.Ints(seatIds)
	last := seatIds[0]
	for _, id := range seatIds[1:] {
		if id-last > 1 {
			break
		}
		last = id
	}

	fmt.Println("Part 1")
	fmt.Printf("Highest seat ID is: %d\n", seatIds[len(seatIds)-1])

	fmt.Println("Part 2")
	fmt.Printf("My seat is: %d\n", last+1)
}

func mustNotError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
