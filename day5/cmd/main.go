package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/torbensky/adventofcode2020/common"
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
	seatID int
}

type boardingPassList []boardingPass

func (s boardingPassList) Len() int {
	return len(s)
}
func (s boardingPassList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s boardingPassList) Less(i, j int) bool {
	return s[i].seatID < s[j].seatID
}

func loadData(path string) boardingPassList {

	var passes boardingPassList
	common.ScanLines(common.GetInputFilePath(), func(line string) bool {
		rowRange := numRange{lower: 0, upper: 127}
		colRange := numRange{lower: 0, upper: 7}

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

		passes = append(passes, boardingPass{
			row:    rowRange.lower,
			column: colRange.lower,
			seatID: rowRange.lower*8 + colRange.lower,
		})

		return true
	})

	return passes
}

func main() {

	// Validate program usage
	if len(os.Args) != 2 {
		log.Fatal("This command accepts only one argument: the path to the input file")
	}
	passes := loadData(os.Args[1])

	// Sort passes so they are ordered increasing by ID
	sort.Sort(passes)

	// Iterate over each seat until we find a gap in the ID's
	lastSeatID := passes[0].seatID
	for _, p := range passes[1:] {
		if p.seatID-lastSeatID > 1 {
			break
		}
		lastSeatID = p.seatID
	}

	fmt.Println("Part 1")
	fmt.Printf("Highest seat ID is: %d\n", passes[len(passes)-1].seatID)

	fmt.Println("Part 2")
	fmt.Printf("My seat is: %d\n", lastSeatID+1)
}

func mustNotError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
