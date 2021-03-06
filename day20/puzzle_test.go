package day20

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	common "github.com/torbensky/adventofcode-common"
)

func TestSolve(t *testing.T) {
	tiles := loadTestTiles(t)
	a := tiles.Solve()
	fmt.Println(a)
	t.Fail()
}

// func TestArrangement(t *testing.T) {
// 	tiles := loadTestTiles(t)
// 	a := emptyArrangement(tiles)

// 	if len(a) != 3 {
// 		t.Fatal("arrangement should have height of 3")
// 	}

// 	if len(a[0]) != 3 {
// 		t.Fatal("arrangement should have width of 3")
// 	}
// }

func TestImageFlip(t *testing.T) {
	img := Image{
		[]byte{'1', '2', '3'},
		[]byte{'4', '5', '6'},
		[]byte{'7', '8', '9'},
	}
	img.FlipX()
	imgF := Image{
		[]byte{'7', '8', '9'},
		[]byte{'4', '5', '6'},
		[]byte{'1', '2', '3'},
	}
	if !img.compare(imgF) {
		t.Logf("expected:\n%s", imgF)
		t.Logf("got:\n%s", img)
		t.Fatal("img flipX failed")
	}

	img.FlipY()
	imgF = Image{
		[]byte{'9', '8', '7'},
		[]byte{'6', '5', '4'},
		[]byte{'3', '2', '1'},
	}
	if !img.compare(imgF) {
		t.Logf("expected:\n%s", imgF)
		t.Logf("got:\n%s", img)
		t.Fatal("img flipY failed")
	}
}

func TestImageRotate(t *testing.T) {
	img := Image{
		[]byte{'1', '2', '3'},
		[]byte{'4', '5', '6'},
		[]byte{'7', '8', '9'},
	}
	img.Rotate90()
	imgR := Image{
		[]byte{'7', '4', '1'},
		[]byte{'8', '5', '2'},
		[]byte{'9', '6', '3'},
	}
	if !img.compare(imgR) {
		t.Logf("expected:\n%s", imgR)
		t.Logf("got:\n%s", img)
		t.Fatal("img rotate90 failed")
	}

	img.Rotate90()
	imgR = Image{
		[]byte{'9', '8', '7'},
		[]byte{'6', '5', '4'},
		[]byte{'3', '2', '1'},
	}
	if !img.compare(imgR) {
		t.Logf("expected:\n%s", imgR)
		t.Logf("got:\n%s", img)
		t.Fatal("img rotate90 failed")
	}

	img.Rotate90()
	img.Rotate90()
	imgR = Image{
		[]byte{'1', '2', '3'},
		[]byte{'4', '5', '6'},
		[]byte{'7', '8', '9'},
	}
	if !img.compare(imgR) {
		t.Logf("expected:\n%s", imgR)
		t.Logf("got:\n%s", img)
		t.Fatal("img rotate90 failed")
	}
}

func TestAlignTiles(t *testing.T) {
	tiles := loadTestTiles(t)
	/*

		1951    2311    3079
		2729    1427    2473
		2971    1489    1171

	*/
	tiles[1427].FlipX()
	t.Logf("matching to center tile 1427...\n%s", tiles[1427])

	// tiles[1427].alignTo(tiles[2311], topEdge)
	// if tiles[2311].Edges[topEdge] != tiles[1427].Edges[bottomEdge] {
	// 	t.Fatal("t2311 top != t1427 bottom")
	// }

	tiles[2311].alignTo(tiles[1427], topEdge)
	if !tiles[2311].Edges[bottomEdge].Match(tiles[1427].Edges[topEdge].flip()) {
		t.Logf("2311:\n%s", tiles[2311])
		t.Fatal("2311 bottom failed to align to top of 1427")
	}

	tiles[2729].alignTo(tiles[1427], leftEdge)
	if !tiles[2729].Edges[rightEdge].Match(tiles[1427].Edges[leftEdge].flip()) {
		t.Logf("2729:\n%s", tiles[2729])
		t.Fatal("2729 right failed to align to left of 1427")
	}

	tiles[1489].alignTo(tiles[1427], bottomEdge)
	fmt.Printf("align 1489:\n%s", tiles[1489])
	if !tiles[1489].Edges[topEdge].Match(tiles[1427].Edges[bottomEdge].flip()) {
		t.Logf("1489:\n%s", tiles[1489])

		t.Fatal("1489 top failed to align to bottom of 1427")
	}

	tiles[2473].alignTo(tiles[1427], rightEdge)
	if !tiles[2473].Edges[leftEdge].Match(tiles[1427].Edges[rightEdge].flip()) {
		t.Logf("2473:\n%s", tiles[2473])
		t.Fatal("2473 left failed to align to right of 1427")
	}

	tiles[3079].alignTo(tiles[2311], rightEdge)
	if tiles[2311].Edges[rightEdge] != tiles[3079].Edges[leftEdge].flip() {
		t.Log("aligned to:\n", tiles[2311].String())
		t.Fatal("tiles[2311] right != tiles[3079] left\n")
	}

}

func TestTileSet(t *testing.T) {
	tiles := loadTestTiles(t)
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

}

func TestTileMatch(t *testing.T) {
	tiles := loadTestTiles(t)

	var deleted []Tile
	del := func(t Tile) {
		deleted = append(deleted, t)
		delete(tiles, t.ID)
	}
	restore := func() {
		for _, t := range deleted {
			tiles[t.ID] = t
		}
		deleted = nil
	}

	for _, c := range []struct {
		id      int
		matches []int
	}{
		{id: 1951, matches: []int{2311, 2729}},
		{id: 2729, matches: []int{1951, 2971, 1427}},
		{id: 2971, matches: []int{2729, 1489}},

		{id: 2311, matches: []int{1951, 3079}},
		{id: 1427, matches: []int{2311, 2473, 1489, 2729}},
		{id: 1489, matches: []int{2971, 1427, 1171}},

		{id: 3079, matches: []int{2311, 2473}},
		{id: 2473, matches: []int{1427, 3079, 1171}},
		{id: 1171, matches: []int{2473, 1489}},
	} {
		tile := tiles[c.id]
		del(tile)

		matched := []int{}
		for i := 0; i < 4; i++ {
			e := tile.Edges[i]
			m := tiles.FindMatchTile(e)
			if m != nil {
				matched = append(matched, m.ID)
			}
		}

		for _, m := range c.matches {
			found := false
			for _, mid := range matched {
				if mid == m {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("tile %d missing match %d\n", tile.ID, m)
			}
		}

		restore()
	}
}

func TestTileSetMatchTile(t *testing.T) {
	tiles := loadTestTiles(t)
	groups := tiles.GetTileGroups()

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

	mTile := tiles[2729]
	delete(tiles, 2729)
	matched := tiles.FindMatchTile(mTile.Edges[rightEdge])
	if matched.ID != 1427 {
		t.Fatalf("tile did not match: wanted %d, got %d", 1427, matched.ID)
	}

	mTile = tiles[1427]
	delete(tiles, 1427)
	matched = tiles.FindMatchTile(mTile.Edges[leftEdge])
	if matched != nil {
		t.Fatal("expected no match")
	}
}

func TestLoadTiles(t *testing.T) {
	loadTestTiles(t)
}

func TestNewTile(t *testing.T) {
	tiles := loadTestTiles(t)
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

	for et := range t2311.Edges {
		if t2311.Edges[et] != tiles[2311].Edges[et] {
			t.Fatal("Tile 2311 loaded improperly\n")
		}
	}
}

func TestTileFlipY(t *testing.T) {
	tiles := loadTestTiles(t)
	tiles[2311].FlipY()
	wanted := newEdge(".#..#.##..")
	if !tiles[2311].Edges[topEdge].Match(wanted) {
		t.Fatalf("top edge is wrong %s != %s", tiles[2311].Edges[topEdge], wanted)
	}

	wanted = newEdge(".#####..#.")
	if !tiles[2311].Edges[rightEdge].Match(wanted) {
		t.Fatalf("right edge is wrong %s != %s", tiles[2311].Edges[rightEdge], wanted)
	}

	wanted = newEdge("..###..###")
	if !tiles[2311].Edges[bottomEdge].Match(wanted) {
		t.Fatalf("bottom edge is wrong %s != %s", tiles[2311].Edges[bottomEdge], wanted)
	}

	wanted = newEdge("#..##.#...")
	if !tiles[2311].Edges[leftEdge].Match(wanted) {
		t.Fatalf("left edge is wrong %s != %s", tiles[2311].Edges[leftEdge], wanted)
	}

	tiles[2311].FlipY()
	wanted = newEdge("..##.#..#.")
	if !tiles[2311].Edges[topEdge].Match(wanted) {
		t.Fatalf("top edge is wrong %s != %s", tiles[2311].Edges[topEdge], wanted)
	}

	wanted = newEdge("...#.##..#")
	if !tiles[2311].Edges[rightEdge].Match(wanted) {
		t.Fatalf("right edge is wrong %s != %s", tiles[2311].Edges[rightEdge], wanted)
	}

	wanted = newEdge("###..###..")
	if !tiles[2311].Edges[bottomEdge].Match(wanted) {
		t.Fatalf("bottom edge is wrong %s != %s", tiles[2311].Edges[bottomEdge], wanted)
	}

	wanted = newEdge(".#..#####.")
	if !tiles[2311].Edges[leftEdge].Match(wanted) {
		t.Fatalf("left edge is wrong %s != %s", tiles[2311].Edges[leftEdge], wanted)
	}
}

func TestTileFlipX(t *testing.T) {
	tiles := loadTestTiles(t)
	tiles[2311].FlipX()
	wanted := newEdge("..###..###")
	if !tiles[2311].Edges[topEdge].Match(wanted) {
		t.Fatalf("top edge is wrong %s != %s", tiles[2311].Edges[topEdge], wanted)
	}

	wanted = newEdge("#..##.#...")
	if !tiles[2311].Edges[rightEdge].Match(wanted) {
		t.Fatalf("right edge is wrong %s != %s", tiles[2311].Edges[rightEdge], wanted)
	}

	wanted = newEdge(".#..#.##..")
	if !tiles[2311].Edges[bottomEdge].Match(wanted) {
		t.Fatalf("bottom edge is wrong %s != %s", tiles[2311].Edges[bottomEdge], wanted)
	}

	wanted = newEdge(".#####..#.")
	if !tiles[2311].Edges[leftEdge].Match(wanted) {
		t.Fatalf("left edge is wrong %s != %s", tiles[2311].Edges[leftEdge], wanted)
	}

	tiles[2311].FlipX()

	wanted = Edge{current: 0b0011010010}
	if !tiles[2311].Edges[topEdge].Match(wanted) {
		t.Fatalf("top edge is wrong %s != %s", tiles[2311].Edges[topEdge], wanted)
	}

	wanted = Edge{current: 0b0001011001}
	if !tiles[2311].Edges[rightEdge].Match(wanted) {
		t.Fatalf("right edge is wrong %s != %s", tiles[2311].Edges[rightEdge], wanted)
	}

	wanted = Edge{current: 0b1110011100}
	if !tiles[2311].Edges[bottomEdge].Match(wanted) {
		t.Fatalf("bottom edge is wrong %s != %s", tiles[2311].Edges[bottomEdge], wanted)
	}

	wanted = Edge{current: 0b0100111110}
	if !tiles[2311].Edges[leftEdge].Match(wanted) {
		t.Fatalf("left edge is wrong %s != %s", tiles[2311].Edges[leftEdge], wanted)
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
	want = newEdge("...#.##..#")
	if !theTile.Edges[bottomEdge].Match(want) {
		t.Fatalf("rotate failed for bottom: %s != %s\n", want, theTile.Edges[bottomEdge])
	}
	want = newEdge("###..###..")
	if !theTile.Edges[leftEdge].Match(want) {
		t.Fatalf("rotate failed for left: %s != %s\n", want, theTile.Edges[leftEdge])
	}

	// 	theTile = newTile(strings.Split(`Tile 1111:

	// `, "\n"))
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
		t.Fatal("tiles don't match")
	}

	wanted := Edge{current: 0b0011010010}
	if !t1.Edges[topEdge].Match(wanted) {
		t.Fatalf("top edge is wrong %s != %s", t1.Edges[topEdge], wanted)
	}

	wanted = Edge{current: 0b0001011001}
	if !t1.Edges[rightEdge].Match(wanted) {
		t.Fatalf("right edge is wrong %s != %s", t1.Edges[rightEdge], wanted)
	}

	wanted = Edge{current: 0b1110011100}
	if !t1.Edges[bottomEdge].Match(wanted) {
		t.Fatalf("bottom edge is wrong %s != %s", t1.Edges[bottomEdge], wanted)
	}

	wanted = Edge{current: 0b0100111110}
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

	e := newEdge("##########")
	if e.current != 0b1111111111 || e.flipped != 0b1111111111 {
		t.Fatalf("newEdge failed to load properly: %s,%s", strconv.FormatInt(int64(e.current), 2), strconv.FormatInt(int64(e.flipped), 2))
	}

	e = newEdge("..........")
	if e.current != 0 || e.flipped != 0 {
		t.Fatalf("newEdge failed to load properly: %s,%s", strconv.FormatInt(int64(e.current), 2), strconv.FormatInt(int64(e.flipped), 2))
	}
}

func loadTestTiles(t *testing.T) TileSet {
	f := common.OpenFile("./test-input.txt")
	defer f.Close()
	tiles := LoadTiles(f)
	if len(tiles) != 9 {
		t.Fatalf("expected 9 tiles, got %d\n", len(tiles))
	}
	return tiles
}

func TestTileString(t *testing.T) {
	tiles := loadTestTiles(t)
	for id := range tiles {
		tileData := getTileData(fmt.Sprint(id))
		if tileData != tiles[id].String() {
			t.Fatal("failed match")
		}

		tiles[id].Rotate90N(4)
		if tileData != tiles[id].String() {
			t.Fatal("failed match")
		}

		tiles[id].Rotate90()
		if tileData != tiles[id].String() {
			t.Fatal("failed match")
		}
	}
}

func getTileData(tileID string) string {
	f := common.OpenFile("./test-input.txt")
	defer f.Close()

	var result strings.Builder
	found := false
	done := false
	common.ScanLines(f, func(line string) {
		if done {
			return
		}

		if strings.Contains(line, tileID) {
			found = true
			return
		}

		if found {
			if line == "" {
				done = true
			} else {
				result.WriteString(line)
				result.WriteByte('\n')
			}
		}
	})
	return result.String()
}
