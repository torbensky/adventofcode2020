package day20

import (
	"strings"
	"testing"

	common "github.com/torbensky/adventofcode-common"
)

// func TestSolve(t *testing.T) {
// 	f := common.OpenFile("./test-input.txt")
// 	defer f.Close()
// 	tiles := LoadTiles(f)
// 	a := tiles.Solve()
// 	fmt.Println(a)
// 	t.Fail()
// }

func TestArrangement(t *testing.T) {
	f := common.OpenFile("./test-input.txt")
	defer f.Close()
	tiles := LoadTiles(f)
	a := emptyArrangement(tiles)

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

func TestTileSet(t *testing.T) {
	f := common.OpenFile("./test-input.txt")
	defer f.Close()

	tiles := LoadTiles(f)
	if len(tiles) != 9 {
		t.Fatalf("expected 9 tiles, got %d\n", len(tiles))
	}

	groups := tiles.GetTileGroups()
	if len(groups.CornerTiles) != 4 {
		t.Fatalf("corner tiles wrong: wanted %d got %d", 4, len(groups.CornerTiles))
	}

	if len(groups.PerimeterTiles) != 8 {
		t.Fatalf("perimiter tiles wrong: wanted %d got %d", 8, len(groups.PerimeterTiles))
	}

	if len(groups.InteriorTiles) != 1 {
		t.Fatalf("interior tiles wrong: wanted %d got %d", 1, len(groups.PerimeterTiles))
	}

	if _, ok := groups.InteriorTiles[1427]; !ok {
		t.Fatalf("Interior is missing tile %d\n", 1427)
	}

	for _, tn := range []int{1951, 3079, 2971, 1171} {
		if _, ok := groups.CornerTiles[tn]; !ok {
			t.Fatalf("Corners is missing tile %d\n", tn)
		}
	}

	for _, tn := range []int{1951, 3079, 2971, 1171, 2311, 2473, 1489, 2729} {
		if _, ok := groups.PerimeterTiles[tn]; !ok {
			t.Fatalf("Perimiter is missing tile %d\n", tn)
		}
	}

	for _, c := range groups.CornerTiles {
		delete(groups.PerimeterTiles, c.ID)
		count := 0
		for i := 0; i < 4; i++ {
			matched := groups.PerimeterTiles.FindMatchTile(c.Edges[i])
			if matched != nil {
				count++
			}
		}
		if count != 2 {
			t.Fatalf("incorrect match count %d for corner tile %s\n", count, c)
		}
	}

	t2729 := tiles[2729]
	delete(tiles, 2729)
	matched := tiles.FindMatchTile(t2729.Edges[rightEdge])
	if matched.ID != 1427 {
		t.Fatalf("tile did not match: wanted %d, got %d", 1427, matched.ID)
	}

	delete(tiles, 1427)
	matched = tiles.FindMatchTile(t2729.Edges[rightEdge])
	if matched != nil {
		t.Fatal("expected no match")
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

func TestTileFlipY(t *testing.T) {
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
	t2311.FlipY()
	wanted := newEdge(".#..#.##..")
	if !t2311.Edges[topEdge].Match(wanted) {
		t.Fatalf("top edge is wrong %s != %s", t2311.Edges[topEdge], wanted)
	}

	wanted = newEdge(".#####..#.")
	if !t2311.Edges[rightEdge].Match(wanted) {
		t.Fatalf("right edge is wrong %s != %s", t2311.Edges[rightEdge], wanted)
	}

	wanted = newEdge("###..###..")
	if !t2311.Edges[bottomEdge].Match(wanted) {
		t.Fatalf("bottom edge is wrong %s != %s", t2311.Edges[bottomEdge], wanted)
	}

	wanted = newEdge("...#.##..#")
	if !t2311.Edges[leftEdge].Match(wanted) {
		t.Fatalf("left edge is wrong %s != %s", t2311.Edges[leftEdge], wanted)
	}

	t2311.FlipY()
	wanted = newEdge(reverse(".#..#.##.."))
	if !t2311.Edges[topEdge].Match(wanted) {
		t.Fatalf("top edge is wrong %s != %s", t2311.Edges[topEdge], wanted)
	}

	wanted = newEdge("...#.##..#")
	if !t2311.Edges[rightEdge].Match(wanted) {
		t.Fatalf("right edge is wrong %s != %s", t2311.Edges[rightEdge], wanted)
	}

	wanted = newEdge(reverse("###..###.."))
	if !t2311.Edges[bottomEdge].Match(wanted) {
		t.Fatalf("bottom edge is wrong %s != %s", t2311.Edges[bottomEdge], wanted)
	}

	wanted = newEdge(".#####..#.")
	if !t2311.Edges[leftEdge].Match(wanted) {
		t.Fatalf("left edge is wrong %s != %s", t2311.Edges[leftEdge], wanted)
	}
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func TestTileFlipX(t *testing.T) {
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
	t2311.FlipX()
	wanted := newEdge("..###..###")
	if !t2311.Edges[topEdge].Match(wanted) {
		t.Fatalf("top edge is wrong %s != %s", t2311.Edges[topEdge], wanted)
	}

	wanted = newEdge("#..##.#...")
	if !t2311.Edges[rightEdge].Match(wanted) {
		t.Fatalf("right edge is wrong %s != %s", t2311.Edges[rightEdge], wanted)
	}

	wanted = newEdge("..##.#..#.")
	if !t2311.Edges[bottomEdge].Match(wanted) {
		t.Fatalf("bottom edge is wrong %s != %s", t2311.Edges[bottomEdge], wanted)
	}

	wanted = newEdge(".#..#####.")
	if !t2311.Edges[leftEdge].Match(wanted) {
		t.Fatalf("left edge is wrong %s != %s", t2311.Edges[leftEdge], wanted)
	}

	t2311.FlipX()

	wanted = Edge{current: 0b0011010010}
	if !t2311.Edges[topEdge].Match(wanted) {
		t.Fatalf("top edge is wrong %s != %s", t2311.Edges[topEdge], wanted)
	}

	wanted = Edge{current: 0b0001011001}
	if !t2311.Edges[rightEdge].Match(wanted) {
		t.Fatalf("right edge is wrong %s != %s", t2311.Edges[rightEdge], wanted)
	}

	wanted = Edge{current: 0b0011100111}
	if !t2311.Edges[bottomEdge].Match(wanted) {
		t.Fatalf("bottom edge is wrong %s != %s", t2311.Edges[bottomEdge], wanted)
	}

	wanted = Edge{current: 0b0111110010}
	if !t2311.Edges[leftEdge].Match(wanted) {
		t.Fatalf("left edge is wrong %s != %s", t2311.Edges[leftEdge], wanted)
	}
}

func TestTileRotate(t *testing.T) {
	theTile := newTile(strings.Split(`Tile 2311:
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

	t.Log("BEFORE\n", theTile)
	theTile.Rotate90()
	t.Log("AFTER\n", theTile)

	want := newEdge(".#..#####.")
	if !theTile.Edges[topEdge].Match(want) {
		t.Fatalf("rotate failed for top: %s != %s\n", want, theTile.Edges[topEdge])
	}
	want = newEdge("..##.#..#.")
	if !theTile.Edges[rightEdge].Match(want) {
		t.Fatalf("rotate failed for right: %s != %s\n", want, theTile.Edges[rightEdge])
	}
	want = newEdge("#..##.#...")
	if !theTile.Edges[bottomEdge].Match(want) {
		t.Fatalf("rotate failed for bottom: %s != %s\n", want, theTile.Edges[bottomEdge])
	}
	want = newEdge("..###..###")
	if !theTile.Edges[leftEdge].Match(want) {
		t.Fatalf("rotate failed for left: %s != %s\n", want, theTile.Edges[leftEdge])
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

	// Test tile String()

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
