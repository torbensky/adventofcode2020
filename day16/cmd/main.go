package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	common "github.com/torbensky/adventofcode-common"
)

func main() {
	fmt.Printf("Part 1: %d\n", part1(common.OpenInputFile()))
	fmt.Printf("Part 2: %d\n", part2(common.OpenInputFile()))
}

type fieldSet map[string]struct{}

// returns a set that is the intersection of sets a and b
func intersect(a, b fieldSet) fieldSet {
	result := make(fieldSet)

	for k := range a {
		if _, ok := b[k]; ok {
			result[k] = struct{}{}
		}
	}

	for k := range b {
		if _, ok := a[k]; ok {
			result[k] = struct{}{}
		}
	}

	return result
}

type ticketSchema map[string][2]fieldRange

func (ts ticketSchema) print() {
	for k, v := range ts {
		fmt.Printf("%s:\n\t[", k)
		for _, r := range v {
			fmt.Printf(" (%d-%d) ", r.start, r.end)
		}
		fmt.Println("]")
	}
}

func (ts ticketSchema) findFields(t ticket, column int) []string {
	var fields []string

	for label := range t.rules[column] {
		fields = append(fields, label)
	}

	return fields
}

func loadTicketData(reader io.Reader) (ticketSchema, []ticket, ticket, int) {

	schema := make(ticketSchema)
	var yourTicket ticket
	var validTickets []ticket

	scanningErrors := 0
	scanMode := 0
	common.ScanLines(reader, func(line string) {

		if line == "" {
			scanMode++
			return
		}

		switch scanMode {
		case 0:
			// Field declaration: "class: 1-3 or 5-7"
			name, ranges := parseFieldLine(line)
			schema[name] = ranges
		case 1:
			// Your ticket header "your ticket:"
			scanMode++
		case 2:
			// Your ticket data "7,1,14"
			// parse data for our ticket
			tikt, _, _ := readTicketData(schema, line)
			// "your ticket" should always be valid
			yourTicket = tikt
			// our ticket  == valid
			// validTickets = append(validTickets, tikt)
		case 3:
			// Nearby tickets header "nearby tickets:"
			scanMode++
		case 4:
			// scan nearby ticket data to the end of the file
			tikt, valid, errors := readTicketData(schema, line)
			if valid {
				validTickets = append(validTickets, tikt)
			} else {
				scanningErrors += errors
			}
		default:
			log.Fatal("unhandled scan mode encountered")
		}
	})

	return schema, validTickets, yourTicket, scanningErrors
}

func part1(reader io.Reader) int {
	_, _, _, errors := loadTicketData(reader)
	return errors
}

func findFieldMatches(fields ticketSchema, val int) fieldSet {
	matches := make(fieldSet)
	for label, ranges := range fields {
		for _, r := range ranges {
			if val >= r.start && val <= r.end {
				// fmt.Printf("valid match found\n")
				// valid field match
				matches[label] = struct{}{}
			}
		}
	}

	// no fields match
	return matches
}

func parseFieldLine(line string) (string, [2]fieldRange) {
	fields := strings.Fields(line)
	var ranges []fieldRange
	for _, f := range fields {
		if strings.Contains(f, "-") {
			ranges = append(ranges, parseRange(f))
		}
	}

	if len(ranges) != 2 {
		log.Fatalf("unexpected fields data: %s\n", line)
	}

	label := strings.Split(line, ":")[0]

	return label, [2]fieldRange{ranges[0], ranges[1]}
}

type fieldRange struct {
	start int
	end   int
}

func parseRange(rangeStr string) fieldRange {
	result := strings.Split(rangeStr, "-")
	if len(result) != 2 {
		log.Fatal("parseRange failed")
	}

	rStart := common.Atoi(result[0])
	rEnd := common.Atoi(result[1])

	return fieldRange{start: rStart, end: rEnd}
}

type ticket struct {
	values []int
	rules  columnRules
}

func (t ticket) print() {
	for col, val := range t.values {
		fmt.Printf("%d:%d ", col, val)
	}
	fmt.Println()
}

// Attempts to parse the ticket from a line of data
//
// Returns the invalid value if the ticket is found to be invalid
//
// The ticket should be ignored if valid is false
//
func readTicketData(schema ticketSchema, line string) (ticket, bool, int) {
	columns := strings.Split(line, ",")
	t := ticket{rules: make(columnRules), values: make([]int, len(columns))}
	for fieldPos, field := range columns {
		val := common.Atoi(field)
		matches := findFieldMatches(schema, val)
		if len(matches) == 0 {
			return t, false, val
		}
		t.values[fieldPos] = val
		t.rules[fieldPos] = matches
	}

	if len(t.rules) == len(columns) {
		return t, true, -1
	}

	panic("somehow ended up with no invalidations but not enough matches")
}

type columnRules map[int]fieldSet

func (cr columnRules) countFields(col int) int {
	return len(cr[col])
}

func (cr columnRules) deleteField(field string) {
	for col := 0; col < len(cr); col++ {
		delete(cr[col], field)
	}
}

func (cr columnRules) identifyNext() (int, string) {
	for col, fields := range cr {
		if len(fields) == 1 {
			for field := range fields {
				return col, field
			}
		}
	}

	panic("bad state - unable to identify further fields")
}

func newColumnRules(tickets []ticket) columnRules {
	rules := make(columnRules)

	for _, t := range tickets {
		for col, fields := range t.rules {

			if _, ok := rules[col]; !ok {
				rules[col] = fields
				continue
			}

			rules[col] = intersect(rules[col], fields)
		}
	}

	return rules
}

func part2(reader io.Reader) int {
	schema, tickets, yourTicket, _ := loadTicketData(reader)

	identified := identifyFields(schema, tickets)
	total := 1
	for col, field := range identified {
		if strings.HasPrefix(field, "departure") {
			total *= yourTicket.values[col]
		}
	}

	return total
}

func identifyFields(schema ticketSchema, tickets []ticket) map[int]string {
	cr := newColumnRules(tickets)
	identified := make(map[int]string)
	for {

		col, field := cr.identifyNext()
		identified[col] = field
		cr.deleteField(field)

		if len(identified) == len(schema) {
			break
		}
	}

	return identified
}
