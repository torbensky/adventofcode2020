package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var lineRegex = regexp.MustCompile(`^(\d+)-(\d+) (\w): (\w+)$`)

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
	validCount := 0
	for scanner.Scan() {
		policy, password, err := processLine(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		if validate(password, policy) {
			validCount++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d valid passwords\n", validCount)

}

type passwordPolicy struct {
	min  int    // min occurrences of char
	max  int    // max occurrences of char
	char string // the char that must occur
}

func processLine(text string) (*passwordPolicy, string, error) {
	results := lineRegex.FindStringSubmatch(text)
	if len(results) < 5 {
		return nil, "", fmt.Errorf("line does not match expected format")
	}

	min, err := strconv.Atoi(results[1])
	if err != nil {
		return nil, "", err
	}

	max, err := strconv.Atoi(results[2])
	if err != nil {
		return nil, "", err
	}

	return &passwordPolicy{
		min:  min,
		max:  max,
		char: results[3],
	}, results[4], nil
}

func validate(password string, policy *passwordPolicy) bool {
	count := strings.Count(password, policy.char)
	if count <= policy.max && count >= policy.min {
		return true
	}

	return false
}
