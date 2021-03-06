package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	common "github.com/torbensky/adventofcode-common"
)

type passport struct {
	data map[string]string
}

// Loads a list of passports from a file
func loadPassportsData(path string) []passport {
	var passports []passport
	addPassport := func(token string) {
		passport := parsePassport(token)
		passports = append(passports, passport)
		return
	}
	file := common.OpenInputFile()
	defer file.Close()
	common.ScanSplit(file, addPassport, common.SplitRecordsFunc)

	return passports
}

// Parses a passport from a chunk of text
func parsePassport(raw string) passport {
	parsed := passport{
		data: make(map[string]string),
	}
	pairs := strings.Fields(raw)
	for _, pair := range pairs {
		passKeyVal := strings.Split(pair, ":")
		parsed.data[passKeyVal[0]] = passKeyVal[1]
	}

	return parsed
}

// Required passport fields
var requiredPassportFields = []string{"byr", "ecl", "eyr", "hcl", "hgt", "iyr", "pid"}

// Checks whether the passport has all required fields
func (p *passport) hasRequiredFields() bool {
	for _, field := range requiredPassportFields {
		if _, ok := p.data[field]; !ok {
			return false
		}
	}

	return true
}

var heightRegex = regexp.MustCompile(`^(\d+)((cm)|(in))$`)
var pidRegex = regexp.MustCompile(`^(\d{9})$`)
var hclRegex = regexp.MustCompile(`^#[0-9a-f]{6}$`)

func (p *passport) isValid() bool {
	// years
	for _, yearValidation := range []struct {
		field string
		min   int
		max   int
	}{
		{"byr", 1920, 2002},
		{"iyr", 2010, 2020},
		{"eyr", 2020, 2030},
	} {
		year, err := strconv.Atoi(p.data[yearValidation.field])
		common.MustNotError(err)
		if year < yearValidation.min || year > yearValidation.max {
			logInvalid(p, yearValidation.field)
			return false
		}
	}

	// height
	heightMatch := heightRegex.FindStringSubmatch(p.data["hgt"])
	if len(heightMatch) != 5 {
		logInvalid(p, "hgt")
		return false
	}

	height, err := strconv.Atoi(heightMatch[1])
	common.MustNotError(err)
	if heightMatch[2] == "cm" {
		if height < 150 || height > 193 {
			logInvalid(p, "hgt")
			return false
		}
	} else {
		if height < 59 || height > 76 {
			logInvalid(p, "hgt")
			return false
		}
	}

	// Eye colour
	switch p.data["ecl"] {
	case "amb", "blu", "brn", "gry", "grn", "hzl", "oth":
		// valid
	default:
		logInvalid(p, "ecl")
		return false
	}

	// Passport ID
	if !pidRegex.MatchString(p.data["pid"]) {
		logInvalid(p, "pid")
		return false
	}

	// Hair colour
	if !hclRegex.MatchString(p.data["hcl"]) {
		logInvalid(p, "hcl")
		return false
	}

	return true
}

func main() {

	// Validate program usage
	if len(os.Args) != 2 {
		log.Fatal("This command accepts only one argument: the path to the input file")
	}
	passports := loadPassportsData(os.Args[1])

	part1Count := 0
	part2Count := 0
	for _, p := range passports {
		if !p.hasRequiredFields() {
			continue
		}
		part1Count++

		if p.isValid() {
			part2Count++
		}
	}
	fmt.Printf("Part 1 valid passport count %d\n", part1Count)
	fmt.Printf("Part 2 valid passport count %d\n", part2Count)
}

// Set to empty string to disable invalid logging
var debugField = ""

// helper function to debug
// no-op when debug disabled or field doesn't match debug field
func logInvalid(p *passport, field string) {
	if debugField == "" {
		return
	}

	if field != debugField {
		return
	}

	fmt.Println("invalid passport")
	fmt.Println("=================================================================================================================")
	fmt.Printf("\tfield:\t%s\t%s\n\n", field, p.data[field])
	fmt.Printf("\tdata:\t%v\n", p.data)
	fmt.Println("=================================================================================================================")
}
