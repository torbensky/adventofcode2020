package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/torbensky/adventofcode2020/common"
)

var lineRegex = regexp.MustCompile(`^(\d+)-(\d+) (\w): (\w+)$`)

func main() {

	validPart1Count := 0
	validPart2Count := 0
	// Scan passwords file and check validation
	common.ScanFile(common.GetInputFilePath(), func(line string) bool {
		// Parse out the data

		policy, password, err := processLine(line)
		if err != nil {
			log.Fatal(err)
		}

		// Count valid passwords

		if validatePart1(password, policy) {
			validPart1Count++
		}

		if validatePart2(password, policy) {
			validPart2Count++
		}
		return true
	}, nil)

	// RESULTS!

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

// Parses a line of the challenge input and returns structured data
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

// Validates the password according to part 2 requirements
// password must contain the designated character at either position v1 OR v2
// NOTE: positions are NOT zero-indexed
func validatePart2(password string, policy *passwordPolicy) bool {
	c1 := []rune(password)[policy.v1-1] // extract char at position v1
	c2 := []rune(password)[policy.v2-1] // extract char at position v2

	// Character must occur at EITHER position v1 OR v2
	return c1 == policy.char && c2 != policy.char || c1 != policy.char && c2 == policy.char
}

// Validates the password according to part 1 requirements
// password must contain between v1-v2 occurences of the designated character
func validatePart1(password string, policy *passwordPolicy) bool {
	count := strings.Count(password, string(policy.char))
	if count <= policy.v2 && count >= policy.v1 {
		return true
	}

	return false
}
