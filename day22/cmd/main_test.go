package main

import (
	"bufio"
	"os"
	"testing"
)

const part1Answer = 306
const part2Answer = 291

func TestPart1(t *testing.T) {
	t.Parallel()
	reader := bufio.NewReader(openTestInput(t))
	got := part1(reader)
	want := part1Answer
	if want != got {
		t.Errorf("expected %d got %d\n", want, got)
	}
}

func TestPart2(t *testing.T) {
	t.Parallel()
	reader := bufio.NewReader(openTestInput(t))
	got := part2(reader)
	want := part2Answer
	if want != got {
		t.Errorf("expected %d got %d\n", want, got)
	}
}

func openTestInput(t *testing.T) *os.File {
	file, err := os.Open("../test-input.txt")
	if err != nil {
		t.Fatalf("unable to open test data: %v\n", err)
	}

	return file
}
