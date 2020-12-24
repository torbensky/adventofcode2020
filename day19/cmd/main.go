package main

import (
	"fmt"
	"io"
	"strings"

	common "github.com/torbensky/adventofcode-common"
)

var debug = false

type elementKind int

const (
	otherRule elementKind = iota
	literal
)

type element string

func (e element) kind() elementKind {
	switch e {
	case "a", "b":
		return literal
	default:
		return otherRule
	}
}

func (e element) literal() byte {
	return e[0]
}

func (e element) ruleNum() int {
	return common.Atoi(string(e))
}

type group []element

func (g group) String() string {
	var sb strings.Builder

	for _, e := range g {
		sb.WriteString(fmt.Sprintf(" %s ", e))
	}

	return sb.String()
}

type rule struct {
	anyOf []group
}

func (r rule) String() string {
	var sb strings.Builder

	for i, g := range r.anyOf {
		sb.WriteString(g.String())
		if i != len(r.anyOf)-1 {
			sb.WriteString(" | ")
		}
	}

	return sb.String()
}

func parseRuleGroup(data string) group {
	parts := strings.Fields(strings.TrimSpace(strings.ReplaceAll(data, "\"", "")))
	g := make(group, len(parts))

	for i, p := range parts {
		g[i] = element(p)
	}

	return g
}

func parseRuleLine(line string) (int, rule) {
	ruleNumAndData := strings.Split(strings.TrimSpace(line), ":")
	ruleNum := common.Atoi(ruleNumAndData[0])
	result := rule{}

	groups := strings.Split(strings.TrimSpace(ruleNumAndData[1]), "|")
	result.anyOf = make([]group, len(groups))

	for i, g := range groups {
		result.anyOf[i] = parseRuleGroup(g)
	}

	return ruleNum, result
}

func matchGroup(line string, pos int, g group, rules map[int]rule, depth int) []int {

	if pos == len(line) {
		if debug {
			fmt.Printf("\n%sEOL for group %q (%q at pos %d)\n", strings.Repeat("\t", depth), g, line, pos)
		}
		return nil
	}

	if debug {
		if depth > 50 {
			panic("too deep")
		}
		fmt.Printf("\n%smatching %q to group %q (%q at pos %d)\n", strings.Repeat("\t", depth), line[pos:], g, line, pos)
	}

	// base case
	if len(g) == 1 {
		switch g[0].kind() {
		case otherRule:
			return matchRule2(line, pos, rules, g[0].ruleNum(), depth+1)
		case literal:
			if g[0].literal() == line[pos] {
				return []int{pos + 1}
			}

			return nil
		}
	}

	possible := matchGroup(line, pos, g[0:1], rules, depth+1)

	var nextPossible []int
	for _, p := range possible {
		next := matchGroup(line, p, g[1:], rules, depth+1)
		if next != nil {
			nextPossible = append(nextPossible, next...)
		}
	}

	return nextPossible
}

func matchRule2(line string, pos int, rules map[int]rule, ruleNum, depth int) []int {
	rule := rules[ruleNum]
	if debug {
		if depth > 50 {
			panic("too deep")
		}
		fmt.Printf("\n%smatching %q for rule %d: %q\n", strings.Repeat("\t", depth), line, ruleNum, rule)
	}

	var possible []int
	for _, g := range rule.anyOf {
		np := matchGroup(line, pos, g, rules, depth)
		possible = append(possible, np...)
	}

	return possible
}

func matchRule(line string, rules map[int]rule, ruleNum int) bool {
	possibleMatches := matchRule2(line, 0, rules, ruleNum, 0)

	if debug {
		fmt.Println("possible matches were: ", possibleMatches)
	}

	for _, p := range possibleMatches {
		if p == len(line) {
			return true
		}
	}

	return false
}

func main() {
	fmt.Printf("Part 1: %d\n", part1(common.OpenInputFile()))
	fmt.Printf("Part 2: %d\n", part2(common.OpenInputFile()))
}

func part1(reader io.Reader) int {
	matchCount := 0
	doneRules := false
	rules := make(map[int]rule)
	common.ScanLines(reader, func(line string) {
		if line == "" {
			if !doneRules {
				doneRules = true
			}
			return
		}

		if !doneRules {
			ruleNum, rule := parseRuleLine(line)
			rules[ruleNum] = rule
			return
		}

		if matchRule(line, rules, 0) {
			matchCount++
		}
	})
	return matchCount
}

func replace11(n int) string {
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		sb.WriteString(strings.Repeat("42 ", i))
		sb.WriteString(strings.Repeat("31 ", i))

		if i < n {
			sb.WriteString("| ")
		}
	}

	return sb.String()
}

func replace8(n int) string {
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		var sb strings.Builder
		for i := 1; i <= n; i++ {
			sb.WriteString(strings.Repeat("42 ", i))

			if i < n {
				sb.WriteString("| ")
			}
		}

		return sb.String()
	}

	return sb.String()
}

func part2(reader io.Reader) int {
	matchCount := 0
	doneRules := false
	rules := make(map[int]rule)
	common.ScanLines(reader, func(line string) {
		if line == "" {
			if !doneRules {
				doneRules = true
			}
			return
		}

		if strings.HasPrefix(line, "8:") {
			fmt.Println("replacing 8")
			line = "8: " + replace8(20)
		}

		if strings.HasPrefix(line, "11:") {
			fmt.Println("replacing 11")
			line = "11: " + replace11(20)
		}

		if !doneRules {
			ruleNum, rule := parseRuleLine(line)
			rules[ruleNum] = rule
			return
		}

		if matchRule(line, rules, 0) {
			matchCount++
		}
	})
	return matchCount
}
