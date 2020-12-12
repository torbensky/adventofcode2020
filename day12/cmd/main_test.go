package main

import (
	"bufio"
	"os"
	"testing"
)

const part1Answer = 25
const part2Answer = 286

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

func TestRotate90(t *testing.T) {
	for _, testRot := range []struct {
		initial  direction
		left     bool
		expected direction
	}{
		// Left
		{north, true, west},
		{west, true, south},
		{south, true, east},
		{east, true, north},
		// Right
		{north, false, east},
		{west, false, north},
		{south, false, west},
		{east, false, south},
	} {
		got := rotate90(testRot.initial, testRot.left)
		if got != testRot.expected {
			t.Errorf("initial %c, rotate left=%t wanted %c, got %c\n", headingNames[testRot.initial], testRot.left, headingNames[testRot.expected], headingNames[got])
		}
	}
}

func openTestInput(t *testing.T) *os.File {
	file, err := os.Open("../test-input.txt")
	if err != nil {
		t.Fatalf("unable to open test data: %v\n", err)
	}

	return file
}
