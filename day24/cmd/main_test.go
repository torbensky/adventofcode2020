package main

import (
	"bufio"
	"os"
	"testing"

	common "github.com/torbensky/adventofcode-common"
)

func TestHexRing(t *testing.T) {
	ring := hexPos{}.buildRing(1)
	if len(ring) != 6 {
		t.Errorf("wanted 6, got %d", len(ring))
	}
	want := hexPos{x: -1, y: 1, z: 0}
	got := ring[0]
	if got != want {
		t.Errorf("wanted %v got %v\n", want, got)
	}

	want = hexPos{x: 0, y: -1, z: 1}
	got = ring[2]
	if got != want {
		t.Errorf("wanted %v got %v\n", want, got)
	}

	want = hexPos{x: 0, y: 1, z: -1}
	got = ring[5]
	if got != want {
		t.Errorf("wanted %v got %v\n", want, got)
	}

	ring = hexPos{}.buildRing(2)
	if len(ring) != 12 {
		t.Errorf("wanted 12, got %d", len(ring))
	}
	want = hexPos{x: -1, y: 2, z: -1}
	got = ring[0]
	if got != want {
		t.Errorf("wanted %v got %v\n", want, got)
	}

	want = hexPos{x: 0, y: -2, z: 2}
	got = ring[5]
	if got != want {
		t.Errorf("wanted %v got %v\n", want, got)
	}

	want = hexPos{x: 0, y: 2, z: -2}
	got = ring[11]
	if got != want {
		t.Errorf("wanted %v got %v\n", want, got)
	}
}

func TestHexGrid(t *testing.T) {
	g := newGrid()

	want := 0
	got := g.countBlack()
	if want != got {
		t.Errorf("wanted %d got %d\n", want, got)
	}

	g.FollowInstructions("ne")
	want = 1
	got = g.countBlack()
	if want != got {
		t.Errorf("wanted %d got %d\n", want, got)
	}

	g.FollowInstructions("ne")
	want = 0
	got = g.countBlack()
	if want != got {
		t.Errorf("wanted %d got %d\n", want, got)
	}

	g.FollowInstructions("neeeeee")
	want = 1
	got = g.countBlack()
	if want != got {
		t.Errorf("wanted %d got %d\n", want, got)
	}

	g.FollowInstructions("neeeeee")
	want = 0
	got = g.countBlack()
	if want != got {
		t.Errorf("wanted %d got %d\n", want, got)
	}

	g.FollowInstructions("nwswseenenwsw")
	want = 1
	got = g.countBlack()
	if want != got {
		t.Errorf("wanted %d got %d\n", want, got)
	}
}

func TestPart1(t *testing.T) {
	t.Parallel()
	f := openTestInput(t)
	defer f.Close()
	reader := bufio.NewReader(f)
	hg := newGrid()
	common.ScanLines(reader, func(line string) {
		hg.FollowInstructions(line)
	})

	got := hg.countBlack()
	want := 10
	if want != got {
		t.Errorf("expected %d got %d\n", want, got)
	}
}

func TestPart2(t *testing.T) {
	t.Parallel()
	f := openTestInput(t)
	defer f.Close()
	reader := bufio.NewReader(f)
	hg := newGrid()
	common.ScanLines(reader, func(line string) {
		hg.FollowInstructions(line)
	})

	got := hg.countBlack()
	want := 10
	if want != got {
		t.Errorf("expected %d got %d\n", want, got)
	}

	for _, want := range []int{15, 12, 25, 14, 23, 28, 41, 37} {
		hg = hg.cycleDay()
		got = hg.countBlack()
		if want != got {
			t.Errorf("expected %d got %d\n", want, got)
		}
	}

	for i := 0; i < 92; i++ {
		hg = hg.cycleDay()
	}
	got = hg.countBlack()
	want = 2208
	if want != got {
		t.Errorf("expected %d got %d\n", want, got)
	}
}

func TestHexPos(t *testing.T) {
	cur := hexPos{}

	if cur.x != 0 || cur.y != 0 || cur.z != 0 {
		t.Error("initial hex is wrong")
	}

	cur = hexPos{}.NorthWest()
	if cur.x != 0 || cur.y != 1 || cur.z != -1 {
		t.Errorf("NorthWest is wrong %v\n", cur)
	}

	cur = hexPos{}.West()
	if cur.x != -1 || cur.y != 1 || cur.z != 0 {
		t.Errorf("West is wrong: %v\n", cur)
	}

	cur = hexPos{}.SouthWest()
	if cur.x != -1 || cur.y != 0 || cur.z != 1 {
		t.Errorf("SouthWest is wrong: %v\n", cur)
	}

	cur = hexPos{}.SouthEast()
	if cur.x != 0 || cur.y != -1 || cur.z != 1 {
		t.Errorf("SouthEast is wrong: %v\n", cur)
	}

	cur = hexPos{}.East()
	if cur.x != 1 || cur.y != -1 || cur.z != 0 {
		t.Errorf("East is wrong: %v\n", cur)
	}

	cur = hexPos{}.NorthEast()
	if cur.x != 1 || cur.y != 0 || cur.z != -1 {
		t.Errorf("NorthEast is wrong: %v\n", cur)
	}
}

func openTestInput(t *testing.T) *os.File {
	file, err := os.Open("../test-input.txt")
	if err != nil {
		t.Fatalf("unable to open test data: %v\n", err)
	}

	return file
}
