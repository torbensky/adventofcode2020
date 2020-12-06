package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// splitRecords splits on two consecutive, empty lines
// ignores "\r" carriage returns (so "\n\r\n" or even "\n\r\r\r\r...\n" will delimit tokens)
func splitRecords(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// End of file, and no data/token left
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// Check for token delim
	index := 0
	consecutiveNewLines := 0
	for ; index < len(data); index++ {
		switch data[index] {
		case '\n':
			consecutiveNewLines++
		case '\r':
			// ignore
		default:
			consecutiveNewLines = 0
		}

		// found token delim
		if consecutiveNewLines == 2 {
			// Note: may contain "\r" somewhere
			return index + 1, data[:index-1], nil
		}
	}

	// End of file, remaining data should be a token
	if atEOF {
		return len(data), data, nil
	}

	// Need MOAR
	return 0, nil, nil
}

// Scans the questions file and returns counts
func scanQuestionsFile(path string) (int, int) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(splitRecords)

	totalUnique, totalEveryone := 0, 0
	for scanner.Scan() {
		u, e := parseGroup(scanner.Text())
		totalUnique += u
		totalEveryone += e
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

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
