package main

import (
	"testing"
)

func TestParse(t *testing.T) {
	ring := parseInput(input)

	got := ring.String()
	if got != input {
		t.Fatalf("input parse failed - %q != %q\n", input, got)
	}
}
