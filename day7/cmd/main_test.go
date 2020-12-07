package main

import (
	"strings"
	"testing"
)

const (
	example1 = `light red bags contain 1 bright white bag, 2 muted yellow bags.
	dark orange bags contain 3 bright white bags, 4 muted yellow bags.
	bright white bags contain 1 shiny gold bag.
	muted yellow bags contain 2 shiny gold bags, 9 faded blue bags.
	shiny gold bags contain 1 dark olive bag, 2 vibrant plum bags.
	dark olive bags contain 3 faded blue bags, 4 dotted black bags.
	vibrant plum bags contain 5 faded blue bags, 6 dotted black bags.
	faded blue bags contain no other bags.
	dotted black bags contain no other bags.`
	example2 = `shiny gold bags contain 2 dark red bags.
	dark red bags contain 2 dark orange bags.
	dark orange bags contain 2 dark yellow bags.
	dark yellow bags contain 2 dark green bags.
	dark green bags contain 2 dark blue bags.
	dark blue bags contain 2 dark violet bags.
	dark violet bags contain no other bags.`
)

var conditions = []struct {
	data  string
	part1 int
	part2 int
}{
	{data: example1, part1: 4, part2: 32},
	{data: example2, part1: 0, part2: 126},
}

func TestPart1(t *testing.T) {
	t.Parallel()
	for i, cond := range conditions {
		reader := strings.NewReader(cond.data)
		rules := loadRules(reader)
		want := cond.part1
		got := calcPart1(rules)
		if want != got {
			t.Errorf("Example %d: expected %d got %d\n", i+1, want, got)
		}
	}
}

func TestPart2(t *testing.T) {
	t.Parallel()
	for i, cond := range conditions {
		reader := strings.NewReader(cond.data)
		rules := loadRules(reader)
		want := cond.part2
		got := countAllInnerBags(rules, "shiny gold")
		if want != got {
			t.Errorf("Example %d: expected %d got %d\n", i+1, want, got)
		}
	}
}
