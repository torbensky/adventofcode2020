package day20

import (
	"strings"
	"testing"

	common "github.com/torbensky/adventofcode-common"
)

func TestArrangement(t *testing.T) {
	f := common.OpenFile("./test-input.txt")
	defer f.Close()
	tiles := LoadTiles(f)
	a := newArrangement(tiles)

	if len(a) != 3 {
		t.Fatal("arrangement should have height of 3")
	}

	if len(a[0]) != 3 {
		t.Fatal("arrangement should have width of 3")
	}
}

func TestAlignTiles(t *testing.T) {
	t2311 := newTile(strings.Split(`Tile 2311:
..##.#..#.
##..#.....
#...##..#.
####.#...#
##.##.###.
##...#.###
.#.#.#..##
..#....#..
###...#.#.
..###..###`, "\n"))
	t1427 := newTile(strings.Split(`Tile 1427:
###.##.#..
.#..#.##..
.#.##.#..#
#.#.#.##.#
....#...##
...##..##.
...#.#####
.#.####.#.
..#..###.#
..##.#..#.`, "\n"))

	t1427.alignTo(t2311, topEdge)
	if t2311.Edges[topEdge] != t1427.Edges[bottomEdge] {
		t.Fatal("t2311 top != t1427 bottom")
	}

	t3079 := newTile(strings.Split(`Tile 3079:
#.#.#####.
.#..######
..#.......
######....
####.#..#.
.#...#.##.
#.#####.##
..#.###...
..#.......
..#.###...`, "\n"))

	t.Log("pre align:\n", t3079.String())
	t3079.alignTo(t2311, rightEdge)
	if t2311.Edges[rightEdge] != t3079.Edges[leftEdge] {
		t.Log("post-align:\n", t3079.String())
		t.Log("aligned to:\n", t2311.String())
		t.Fatal("t2311 right != t1427 left\n")
	}
}

func TestLoadTiles(t *testing.T) {
	f := common.OpenFile("./test-input.txt")
	defer f.Close()

	tiles := LoadTiles(f)
	if len(tiles) != 9 {
		t.Fatalf("expected 9 tiles, got %d\n", len(tiles))
	}

	t2311 := newTile(strings.Split(`Tile 2311:
..##.#..#.
##..#.....
#...##..#.
####.#...#
##.##.###.
##...#.###
.#.#.#..##
..#....#..
###...#.#.
..###..###`, "\n"))
	for i := range t2311.Edges {
		if t2311.Edges[i] != tiles[2311].Edges[i] {
			t.Fatal("Tile 2311 loaded improperly\n")
		}
	}
}

func TestTile(t *testing.T) {
	t1 := newTile(strings.Split(`Tile 2311:
..##.#..#.
##..#.....
#...##..#.
####.#...#
##.##.###.
##...#.###
.#.#.#..##
..#....#..
###...#.#.
..###..###`, "\n"))

	if t1.ID != 2311 {
		t.Fatalf("tile id wanted %d got %d\n", 2311, t1.ID)
	}

	wanted := Edge{current: 0b0011010010}
	if !t1.Edges[topEdge].Match(wanted) {
		t.Fatalf("top edge is wrong %s != %s", t1.Edges[topEdge], wanted)
	}

	wanted = Edge{current: 0b0001011001}
	if !t1.Edges[rightEdge].Match(wanted) {
		t.Fatalf("right edge is wrong %s != %s", t1.Edges[rightEdge], wanted)
	}

	wanted = Edge{current: 0b0011100111}
	if !t1.Edges[bottomEdge].Match(wanted) {
		t.Fatalf("bottom edge is wrong %s != %s", t1.Edges[bottomEdge], wanted)
	}

	wanted = Edge{current: 0b0111110010}
	if !t1.Edges[leftEdge].Match(wanted) {
		t.Fatalf("left edge is wrong %s != %s", t1.Edges[leftEdge], wanted)
	}

	tileStr := `..##.#..#.
##..#.....
#...##..#.
####.#...#
##.##.###.
##...#.###
.#.#.#..##
..#....#..
###...#.#.
..###..###
`
	if tileStr != t1.String() {
		t.Log("wanted\n", tileStr)
		t.Log("got\n", t1.String())
		t.Fatal("tile edges don't match")
	}

	t.Log("BEFORE\n", t1)
	t1.Rotate90()
	if t1.Edges[topEdge].String() != ".#..#####." {
		t.Log("AFTER\n", t1)
		t.Fatalf("rotate failed: %s != %s\n", t1.Edges[topEdge], ".#..#####.")
	}
}

func TestEdge(t *testing.T) {
	e1 := newEdge("#........#")
	e2 := newEdge("#........#")

	if e1.String() != "#........#" {
		t.Fatalf("edge String() doesn't work: %s != %s\n", e1.String(), "#........#")
	}

	if !e1.Match(e2) || !e2.Match(e1) {
		t.Fatalf("missing match %s == %s\n", e1, e2)
	}

	e3 := newEdge(".#.......#")
	if e1.Match(e3) || e2.Match(e3) {
		t.Fatal("unexpected match!")
	}

	if e3.String() != ".#.......#" {
		t.Fatalf("edge String() doesn't work: %s != %s\n", e3.String(), ".#.......#")
	}

	e4 := newEdge("..#....#..")
	want := Edge{current: 0b0010000100}
	if !e4.Match(want) {
		t.Fatalf("wanted %s got %s", want, e4)
	}

	// test some rotations
	if !e1.Match(e1.flip()) {
		t.Fatalf("symmetrical rotation failed")
	}
	if e3.Match(e3.flip()) {
		t.Fatalf("unsymmetrical rotation failed")
	}

	want = Edge{current: 0b1000000010}
	if !e3.flip().Match(want) {
		t.Fatalf("wanted %s got %s", want, e3)
	}
}
