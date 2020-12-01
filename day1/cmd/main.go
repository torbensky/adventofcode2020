package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {

	// Validate program usage

	if len(os.Args) != 2 {
		log.Fatal("This command accepts only one argument: the path to the input file")
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Scan in numeric data

	scanner := bufio.NewScanner(file)
	var viableValues []int
	for scanner.Scan() {
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Printf("unable to parse integer - skipping line '%s'\n", scanner.Text())
			continue
		}

		// Don't consider values > 2020
		if val > 2020 {
			continue
		}
		viableValues = append(viableValues, val)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Find the pair that sums to 2020
	val1, val2, err := solve(viableValues)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Solution: %d * %d = %d\n", val1, val2, val1*val2)
}

func solve(values []int) (int, int, error) {
	sort.Ints(values)
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
