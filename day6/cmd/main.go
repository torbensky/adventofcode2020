package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	common "github.com/torbensky/adventofcode-common"
)

// Scans the questions file and returns counts
func scanQuestionsFile(path string) (int, int) {
	totalUnique, totalEveryone := 0, 0
	countQuestions := func(token string) {
		u, e := parseGroup(token)
		totalUnique += u
		totalEveryone += e
	}
	file := common.OpenInputFile()
	defer file.Close()
	common.ScanSplit(file, countQuestions, common.SplitRecordsFunc)

	return totalUnique, totalEveryone
}

func parseGroup(group string) (int, int) {
	// whitespace should separate each person
	peoplesAnswers := strings.Fields(group)
	numPeople := len(peoplesAnswers)

	// Find unique questions
	uniqueQuestions := make(map[rune]int) // count of people per question
	for i := 0; i < numPeople; i++ {
		for _, c := range peoplesAnswers[i] {
			uniqueQuestions[c]++
		}
	}

	// Find questions that everyone had
	totalEveryone := 0
	for _, count := range uniqueQuestions {
		if count == numPeople {
			totalEveryone++
		}
	}

	return len(uniqueQuestions), totalEveryone
}

func main() {

	// Validate program usage
	if len(os.Args) != 2 {
		log.Fatal("This command accepts only one argument: the path to the input file")
	}
	totalUnique, totalEveryone := scanQuestionsFile(os.Args[1])

	fmt.Printf("Part 1: %d\n", totalUnique)
	fmt.Printf("Part 2: %d\n", totalEveryone)
}
