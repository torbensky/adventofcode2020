package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestEvaluateExpression(t *testing.T) {

	for _, c := range []struct {
		expr   string
		result int
	}{
		{"1 + 2", 3},
		{"2 * (2 + 1)", 6},
		{"2 * (2 * 3)", 12},
		{"2 * (2 * 3 + 2 * 3)", 48},
		{"2 * 3 + (4 * 5)", 26},
		{"5 + (8 * 3 + 9 + 3 * 4 * 3)", 437},
		{"5 * 9 * (7 * 3 * 3 + 9 * 3 + (8 + 6 * 4))", 12240},
		{"((2 + 4 * 9) * (6 + 9 * 8 + 6) + 6) + 2 + 4 * 2", 13632},
	} {
		fmt.Printf("expression '%s'\n", c.expr)
		got := evaluateExpr(c.expr)
		if c.result != got {
			t.Fatalf("%s wanted %d got %d", c.expr, c.result, got)
		}
	}

}

func TestPart2(t *testing.T) {
	t.Parallel()
	reader := bufio.NewReader(openTestInput(t))
	got := part2(reader)
	want := 231235959382961
	if want != got {
		t.Errorf("expected %d got %d\n", want, got)
	}
}

func TestPart1(t *testing.T) {
	t.Parallel()
	reader := bufio.NewReader(openTestInput(t))
	got := part1(reader)
	want := 8929569623593
	if want != got {
		t.Errorf("expected %d got %d\n", want, got)
	}
}

func openTestInput(t *testing.T) *os.File {
	file, err := os.Open("../input.txt")
	if err != nil {
		t.Fatalf("unable to open test data: %v\n", err)
	}

	return file
}
