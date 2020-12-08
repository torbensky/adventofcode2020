package main

import (
	"strings"
	"testing"
)

const (
	example1 = `nop +0
	acc +1
	jmp +4
	acc +3
	jmp -3
	acc -99
	acc +1
	jmp -4
	acc +6`
)

var conditions = []struct {
	data  string
	part1 int
	part2 int
}{
	{data: example1, part1: 5, part2: 8},
}

func TestPart1(t *testing.T) {
	t.Parallel()
	for i, cond := range conditions {
		reader := strings.NewReader(cond.data)
		prog := loadProgram(reader)
		_, got := executeProgram(prog)
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
		prog := loadProgram(reader)
		got := fixProgram(prog)
		want := cond.part2
		if want != got {
			t.Errorf("Example %d: expected %d got %d\n", i+1, want, got)
		}
	}
}
