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
	validPart1Count := 0
	validPart2Count := 0
	for scanner.Scan() {
		policy, password, err := processLine(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		if validatePart1(password, policy) {
			validPart1Count++
		}

		if validatePart2(password, policy) {
			validPart2Count++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Part1 - found %d valid passwords\n", validPart1Count)
	fmt.Printf("Part2 - found %d valid passwords\n", validPart2Count)
}

// encodes the information in the toboggan password policy database
//
// Example:
// 1-3 a: abcde
// v1 = 1
// v2 = 3
// char = "a"
type passwordPolicy struct {
	v1   int  // first numerical value of password policy
	v2   int  // second numerical value of password policy
	char rune // the char that must occur
}

func processLine(text string) (*passwordPolicy, string, error) {
	results := lineRegex.FindStringSubmatch(text)
	if len(results) < 5 {
		return nil, "", fmt.Errorf("line does not match expected format")
	}

	v1, err := strconv.Atoi(results[1])
	if err != nil {
		return nil, "", err
	}

	v2, err := strconv.Atoi(results[2])
	if err != nil {
		return nil, "", err
	}

	return &passwordPolicy{
		v1:   v1,
		v2:   v2,
		char: rune(results[3][0]),
	}, results[4], nil
}

func validatePart2(password string, policy *passwordPolicy) bool {
	c1 := []rune(password)[policy.v1-1]
	c2 := []rune(password)[policy.v2-1]

	return c1 == policy.char && c2 != policy.char || c1 != policy.char && c2 == policy.char
}

func validatePart1(password string, policy *passwordPolicy) bool {
	count := strings.Count(password, string(policy.char))
	if count <= policy.v2 && count >= policy.v1 {
		return true
	}

	return false
}
