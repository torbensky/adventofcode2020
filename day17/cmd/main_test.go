package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"testing"
)

const part1Answer = 112
const part2Answer = 848

const exampleData = `.#.
..#
###`

func TestSpaceCycles(t *testing.T) {
	space := load3dSpace(strings.NewReader(exampleData))

	if len(space) != 5 {
		t.Fatalf("cycle 0: want 5 got %d\n", len(space))
	}

	for i := 0; i < 6; i++ {
		space = space.cycle()

		if i == 0 && len(space) != 11 {
			t.Fatalf("cycle 1: want 11 got %d\n", len(space))
		}

		if i == 5 && len(space) != 112 {
			t.Fatalf("cycle 6: want 112 got %d\n", len(space))
		}
	}
}

func TestCountNeighbors(t *testing.T) {
	space := load3dSpace(strings.NewReader(exampleData))
	for _, val := range []struct {
		coord coord3d
		n     int
	}{
		{coord: coord3d{X: 0, Y: 0, Z: 0}, n: 1},
		{coord: coord3d{X: 0, Y: 1, Z: 0}, n: 3},
		{coord: coord3d{X: 1, Y: 1, Z: 0}, n: 5},
		{coord: coord3d{X: 1, Y: 2, Z: 0}, n: 3},
		{coord: coord3d{X: 3, Y: 2, Z: 0}, n: 2},
		{coord: coord3d{X: 3, Y: 3, Z: 0}, n: 1},
		{coord: coord3d{X: 2, Y: 3, Z: 0}, n: 2},
		{coord: coord3d{X: 2, Y: 4, Z: 0}, n: 0},
	} {
		if count := space.countNeighbors(val.coord, 100); count != val.n {
			log.Fatalf("countNeighbors(%s): wanted %d got %d\n", val.coord.text(), val.n, count)
		}
	}
}

func TestMake3dBaseVectors(t *testing.T) {
	allVecs := make3dBaseVectors()
	if len(allVecs) != 26 {
		t.Fatal("not enough vectors")
	}
}

func TestLoadData(t *testing.T) {
	space := load3dSpace(strings.NewReader(exampleData))

	// The example data should have five active nodes
	if len(space) != 5 {
		t.Fatalf("want 5 got %d\n", len(space))
	}

	for _, coord := range []coord3d{
		{X: 1, Y: 0, Z: 0},
		{X: 2, Y: 1, Z: 0},
		{X: 0, Y: 2, Z: 0},
		{X: 1, Y: 2, Z: 0},
		{X: 2, Y: 2, Z: 0},
	} {
		if _, ok := space[coord]; !ok {
			log.Fatalf("missing expected active node: %v\n", coord)
		}
	}

}

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
