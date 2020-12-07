package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/torbensky/adventofcode2020/common"
)

// Scans the questions file and returns counts
func scanQuestionsFile(path string) (int, int) {
	totalUnique, totalEveryone := 0, 0
	common.ScanFile(common.GetInputFilePath(), func(token string) bool {
		u, e := parseGroup(token)
		totalUnique += u
		totalEveryone += e
		return true
	}, common.SplitRecordsFunc)

	return totalUnique, totalEveryone
}

func parseGroup(data string) (int, int) {
	// whitespace should separate each person
	peoplesAnswers := strings.Fields(data)
	numPeople := len(peoplesAnswers)

	// Find unique questions
	uniqueQuestions := map[rune]struct{}{}
	for i := 0; i < numPeople; i++ {
		for _, c := range peoplesAnswers[i] {
			uniqueQuestions[c] = struct{}{}
		}
	}

	// Find questions that everyone had
	totalEveryone := 0
	for question := range uniqueQuestions {
		count := 0
		for i := 0; i < numPeople; i++ {
			if strings.ContainsRune(peoplesAnswers[i], question) {
				count++
			}
		}

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
