package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

type NumRange struct {
	lower int
	upper int
}

func splitRange(r NumRange) (NumRange, NumRange) {
	newRangeSize := (r.upper - r.lower) / 2

	return NumRange{
			lower: r.lower,
			upper: r.lower + newRangeSize,
		}, NumRange{
			lower: r.lower + newRangeSize + 1,
			upper: r.upper,
		}
}

type boardingPass struct {
	row    int
	column int
}

var passes []*boardingPass

func loadData(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		rowRange := NumRange{lower: 0, upper: 127}
		colrange := NumRange{lower: 0, upper: 7}
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
				lower, upper := splitRange(colrange)
				switch c {
				case 'L':
					colrange = lower
				case 'R':
					colrange = upper
				}
			}
		}
		passes = append(passes, &boardingPass{
			row:    rowRange.lower,
			column: colrange.lower,
		})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func main() {

	// Validate program usage
	if len(os.Args) != 2 {
		log.Fatal("This command accepts only one argument: the path to the input file")
	}
	loadData(os.Args[1])

	var seatIds []int
	for _, p := range passes {
		id := p.row*8 + p.column
		seatIds = append(seatIds, id)
	}
	sort.Ints(seatIds)
	last := seatIds[0]
	for _, id := range seatIds[1:] {
		fmt.Printf("id %d, last%d \n", id, last)
		if id-last > 1 {
			fmt.Printf("missing seat might be %d\n", id)
		}
		last = id
	}

	// fmt.Printf("highest seat id is %d\n", highest)
}

func mustNotError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
