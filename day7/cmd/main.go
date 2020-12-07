package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/torbensky/adventofcode2020/common"
)

type bagRule struct {
	bagType  string         // type of outer bag
	contains map[string]int // types of inner bags + count required
}

func main() {
	bagRules := loadRules(common.GetInputFilePath())

	fmt.Println()
	fmt.Printf("Part 1 - Total %d\n\n", calcPart1(bagRules))
	fmt.Printf("Part 2 - Total %d\n\n", countAllInnerBags(bagRules, "shiny gold"))
}

func loadRules(filePath string) map[string]*bagRule {
	bagRules := make(map[string]*bagRule)
	parseRuleLine := func(line string) {
		words := strings.Fields(line)
		outerBag := strings.Join(words[0:2], " ")

		newRule := &bagRule{
			bagType:  outerBag,
			contains: make(map[string]int),
		}

		for i := 4; i < len(words); i += 4 {
			// Check for "no other bags"
			if words[i] == "no" {
				break
			}

			// Find how many bags are required
			numBags, err := strconv.Atoi(words[i])
			if err != nil {
				log.Fatal("can't process bag count")
			}

			innerBag := strings.Join(words[i+1:i+3], " ")
			newRule.contains[innerBag] = numBags
		}

		// Only add rules that are not dead-ends
		if len(newRule.contains) > 0 {
			bagRules[outerBag] = newRule
		}
	}

	common.ScanLines(filePath, common.AllTokensFunc(parseRuleLine))

	return bagRules
}

func calcPart1(bagRules map[string]*bagRule) int {
	// Count the total number of ways to have shiny gold bags
	totalWithShinyGold := 0
	for bt := range bagRules {
		if canContain(bagRules, bt, "shiny gold") {
			totalWithShinyGold++
		}
	}
	return totalWithShinyGold
}

// Checks whether a given bag is allowed to contain the target bag
func canContain(rules map[string]*bagRule, outerBag, targetBag string) bool {
	// base condition, can we go further?
	if rules[outerBag] == nil {
		return false // no
	}

	// Did we find it?
	if _, ok := rules[outerBag].contains[targetBag]; ok {
		return true
	}

	// Maybe an inner bag allows...
	for bt := range rules[outerBag].contains {
		if canContain(rules, bt, targetBag) {
			return true
		}
	}

	// Nope, no inner bags contain it either
	return false
}

// Count the number of inner bags that must be within the given bag type
func countAllInnerBags(rules map[string]*bagRule, bagType string) int {
	// Base condition
	if rules[bagType] == nil {
		return 0
	}

	total := 0
	for bt, count := range rules[bagType].contains {
		total += count + count*countAllInnerBags(rules, bt)
	}

	return total
}
