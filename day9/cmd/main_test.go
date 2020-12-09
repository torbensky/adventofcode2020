package main

import (
	"strings"
	"testing"
)

const (
	example1 = `35
20
15
25
47
40
62
55
65
95
102
117
150
182
127
219
299
277
309
576`
)

var conditions = []struct {
	data  string
	part1 int
	part2 int
}{
	{data: example1, part1: 127, part2: 62},
}

func TestPart1(t *testing.T) {
	t.Parallel()
	for i, cond := range conditions {
		reader := strings.NewReader(cond.data)
		nums := loadNumbers(reader)
		got := part1(nums, 5)
		want := cond.part1
		if want != got {
			t.Errorf("Example %d: expected %d got %d\n", i+1, want, got)
		}
	}
}

func TestPart2(t *testing.T) {
	t.Parallel()
	for i, cond := range conditions {
		reader := strings.NewReader(cond.data)
		nums := loadNumbers(reader)
		got := part2(nums, 127)
		want := cond.part2
		if want != got {
			t.Errorf("Example %d: expected %d got %d\n", i+1, want, got)
		}
	}
}
