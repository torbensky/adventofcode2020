package day20

import (
	"fmt"
	"io"
	"log"
	"math"
	"regexp"
	"strings"

	common "github.com/torbensky/adventofcode-common"
)

// Image is assumed to be a square image
type Image [][]byte

func (img Image) String() string {
	var sb strings.Builder
	for _, row := range img {
		sb.Write(row)
		sb.WriteString("\n")
	}
	return sb.String()
}

func (img Image) compare(other Image) bool {
	if len(img) != len(other) {
		return false
	}
	for y := 0; y < len(img); y++ {
		if len(img[y]) != len(other[y]) {
			return false
		}
		for x := 0; x < len(img[y]); x++ {
			if img[y][x] != other[y][x] {
				return false
			}
		}
	}
	return true
}

type coord2d struct {
	x int
	y int
}

/*

SeaMonsterFilter can filter an area of an image to detect a sea monster!

A sea monster looks like this:

         WIDTH=20

|---------------------|
|                  # |
|#    ##    ##    ###|    HEIGHT = 3
| #  #  #  #  #  #   |
|---------------------|

We store the 2D coordinates of the required pixelType (for our sea monster these are the '#')

The empty areas can be anything. So we just need to ensure that the required pixel types/positions are present

*/
type SeaMonsterFilter struct {
	// height of filter
	height int
	// width of filter
	width int
	// type of pixel we are looking for
	pixelType byte
	// set of pixels that we want to match the given type
	mask map[coord2d]struct{}
}

/*
 check applies the filter to a subset of the image. It assumes a coordinate system where top,right is 0,0:

(0,0)
   \             IMAGE
    +-----------------------------+
	|	(offsetX,offsetY)         |
	|		   \                  |
	|			+-------+         |
	|			|       |         |
	|			|       |         |
	|			|       |         |
	|			+-------+         |
	|                             |
	|			height = 5        |
	|			width = 9         |
	|                             |
    +-----------------------------+


 Returns false, -1 if there is no match.

 Returns (true, count) if there is a match. count is the number of pixels of the same type that were not included

 For example:

   - Assume a mask with 6 unique coordinates
   - Assume we are looking for the pixelType '#'
   - Assume there are total = 10 '#' pixels within the filter area

	If the filter matched:
		count = total - matched = 10 - 4
		return (true, count)

*/
func (f SeaMonsterFilter) check(img Image, offsetX, offsetY int) (bool, int) {

	// Detect whether the filter matches
	for c := range f.mask {
		if img[c.y+offsetY][c.x+offsetX] != f.pixelType {
			return false, -1
		}
	}

	// Count the total of the pixel type that the filter is for so we can tell how many weren't part of the mask
	totalOfType := 0
	for y := 0; y < f.height; y++ {
		for _, b := range img[y+offsetY] {
			if b == f.pixelType {
				totalOfType++
			}
		}
	}

	return true, totalOfType - len(f.mask)
}

// DetectAny finds whether or not there is a single sea monster in this image
func (f SeaMonsterFilter) DetectAny(img Image) bool {
	for offsetY := 0; offsetY < len(img)-f.height; offsetY++ {
		for offsetX := 0; offsetX < len(img[offsetY])-f.width; offsetX++ {
			matched, _ := f.check(img, offsetX, offsetY)
			if matched {
				return true
			}
		}
	}

	return false
}

// Scan finds all sea monsters and counts the number of pixels that are part of sea monster squares
// but not part of the sea monster
func (f SeaMonsterFilter) Scan(img Image) int {

	totalPixelsNotIncluded := 0
	for offsetY := 0; offsetY < len(img)-f.height; offsetY++ {
		for offsetX := 0; offsetX < len(img[offsetY])-f.width; offsetX++ {
			matched, count := f.check(img, offsetX, offsetY)
			if matched {
				totalPixelsNotIncluded += count
			}
		}
	}

	return totalPixelsNotIncluded
}

func emptyImage(size int) Image {
	img := make(Image, size)
	for y := 0; y < size; y++ {
		img[y] = make([]byte, size)
	}

	return img
}

func (img Image) FlipY() {
	for _, row := range img {
		for i := 0; i < len(row)/2; i++ {
			// swap pieces
			row[i], row[len(row)-1-i] = row[len(row)-1-i], row[i]
		}
	}
}

func (img Image) FlipX() {
	for i := 0; i < (len(img) / 2); i++ {
		img[i], img[len(img)-1-i] = img[len(img)-1-i], img[i]
	}
}

func (img Image) Rotate90() {
	size := len(img)

	// in-place matrix rotation
	for i := 0; i < size/2; i++ {
		for j := i; j < size-i-1; j++ {

			// rotate in batches of 4

			temp := img[i][j]

			img[i][j] = img[size-1-j][i]

			img[size-1-j][i] = img[size-1-i][size-1-j]

			img[size-1-i][size-1-j] = img[j][size-1-i]

			img[j][size-1-i] = temp
		}
	}
}

// TileSet represents a set of unique tiles ID->Tile
type TileSet map[int]Tile

// LoadTiles loads tiles from an io source
func LoadTiles(reader io.Reader) TileSet {
	tiles := make(TileSet)

	common.ScanSplit(reader, func(tile string) {

		tile = strings.TrimSpace(tile)
		if tile == "" {
			return
		}

		t := newTile(strings.Split(tile, "\n"))
		tiles[t.ID] = t
	}, common.SplitRecordsFunc)

	return tiles
}

/*

          TOP

	edge direction -->

           0
      -----------
	  |			|
	  |			|
 ^ 3  |			|  1  |
 |	  |			|     |
 |	  |			|     v
  	  -----------
	       2

		 BOTTOM

	<-- edge direction
*/

// Tile represents a tile of the challenge image puzzle
type Tile struct {
	ID    int
	Edges []Edge
	Image Image
}

func (t Tile) String() string {
	var sb strings.Builder

	le, re := t.Edges[leftEdge].String(), t.Edges[rightEdge].String()

	sb.WriteString(t.Edges[topEdge].String())
	sb.WriteByte('\n')

	for i := 0; i < len(t.Image); i++ {
		sb.WriteByte(le[len(le)-2-i])

		sb.Write(t.Image[i])
		sb.WriteByte(re[i+1])

		sb.WriteByte('\n')
	}

	sb.WriteString(reverse(t.Edges[bottomEdge].String()))
	sb.WriteByte('\n')

	return sb.String()
}

// FlipY flips the tile over the y-axis
func (t Tile) FlipY() {
	t.Edges[topEdge] = t.Edges[topEdge].flip()
	t.Edges[bottomEdge] = t.Edges[bottomEdge].flip()

	// swap left/right edges
	t.Edges[leftEdge], t.Edges[rightEdge] = t.Edges[rightEdge].flip(), t.Edges[leftEdge].flip()

	t.Image.FlipY()
}

// FlipX flips the tile over the x-axis
func (t Tile) FlipX() {
	// swap top/bottom edges
	t.Edges[topEdge], t.Edges[bottomEdge] = t.Edges[bottomEdge].flip(), t.Edges[topEdge].flip()

	t.Edges[leftEdge] = t.Edges[leftEdge].flip()
	t.Edges[rightEdge] = t.Edges[rightEdge].flip()

	t.Image.FlipX()
}

// Rotate90N performs a rotation N times
func (t Tile) Rotate90N(rotations int) {
	fmt.Printf("rotate requested %d times\n", rotations)
	if rotations < 0 {
		rotations += 4
	}
	rotations %= 4
	fmt.Printf("will rotate %d times\n", rotations)
	for i := 0; i < rotations; i++ {
		t.Rotate90()
	}
}

// Rotate90 rotates the tile by 90 degrees
func (t Tile) Rotate90() {
	t.Edges[topEdge], t.Edges[rightEdge], t.Edges[bottomEdge], t.Edges[leftEdge] = t.Edges[leftEdge], t.Edges[topEdge], t.Edges[rightEdge], t.Edges[bottomEdge]
	t.Image.Rotate90()
}

var tileIDRegex = regexp.MustCompile(`Tile (\d+)`)

func newTile(rows []string) Tile {

	matches := tileIDRegex.FindStringSubmatch(rows[0])
	rows = rows[1:]
	id := common.Atoi(matches[1])

	tileSize := len(rows)

	var le, re strings.Builder // left and right edges
	image := make(Image, tileSize-2)
	for rowNum, row := range rows {
		// Capture edge data
		le.WriteByte(row[0])
		re.WriteByte(row[len(row)-1])

		// Capture inner data
		if rowNum > 0 && rowNum < len(rows)-1 {
			image[rowNum-1] = []byte(row[1 : tileSize-1])
		}
	}

	return Tile{
		ID:    id,
		Image: image,
		Edges: []Edge{
			newEdge(rows[0]),                    // top
			newEdge(re.String()),                // right
			newEdge(reverse(rows[len(rows)-1])), // bottom
			newEdge(reverse(le.String())),       // left
		},
	}
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

type Edge struct {
	current int
	flipped int
}

func (e Edge) min() int {
	return min(e.current, e.flipped)
}

func (e Edge) flip() Edge {
	return Edge{current: e.flipped, flipped: e.current}
}

func (e Edge) String() string {
	var sb strings.Builder

	b := 1
	// NOTE: this is hard-coded for now to expect the length 10 edges of the challenge
	for i := 9; i >= 0; i-- {
		if b<<i&e.current != 0 {
			sb.WriteByte('#')
		} else {
			sb.WriteByte('.')
		}
	}

	return sb.String()
}

func newEdge(line string) Edge {
	unflipped := 0
	flipped := 0

	for i, c := range line {
		switch c {
		case '#':
			unflipped = setBit(int(unflipped), uint(len(line)-1-i))
			flipped = setBit(int(flipped), uint(i))
		case '.':
			// no-op
		default:
			panic("invalid char encountered")
		}
	}

	return Edge{flipped: flipped, current: unflipped}
}

func setBit(n int, i uint) int {
	n |= (1 << i)
	return n
}

func (e Edge) Match(o Edge) bool {
	return 0 == (o.current ^ e.current)
}

type Arrangement [][]Tile

func (a Arrangement) String() string {

	var sb strings.Builder
	for y := range a {
		for _, tile := range a[y] {
			sb.WriteString(fmt.Sprintf("%d\t", tile.ID))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
	// // TODO: don't use hard-coded tile size
	// tileSize := 10
	// puzzleSize := len(a) * tileSize // assume square
	// puzzle := make([][]byte, puzzleSize)
	// for i := 0; i < puzzleSize; i++ {
	// 	puzzle[i] = make([]byte, puzzleSize)
	// }

	// writeTile := func(t Tile, x, y int) {
	// 	for ty, row := range strings.Split(t.String(), "\n") {
	// 		for tx := range row {
	// 			fmt.Printf("write char %d,%d\n", x+tx, y+ty)
	// 			puzzle[y+ty][x+tx] = row[tx]
	// 		}
	// 	}
	// }

	// for y := range a {
	// 	for x, tile := range a[y] {
	// 		fmt.Printf("write tile %d,%d\n", x, y)
	// 		writeTile(tile, x*10, y*10)
	// 	}
	// }

	// var sb strings.Builder
	// for _, line := range puzzle {
	// 	sb.Write(line)
	// 	sb.WriteByte('\n')
	// }
	// return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (ts TileSet) getMinEdges() map[int]int {
	minEdges := make(map[int]int)
	for _, t := range ts {
		for _, e := range t.Edges {
			minEdges[e.min()]++
		}
	}
	return minEdges
}

type tileGroups struct {
	CornerTiles    TileSet
	PerimeterTiles TileSet
	InteriorTiles  TileSet
}

func (ts TileSet) pickOne() *Tile {
	for _, t := range ts {
		return &t
	}
	return nil
}

// GetTileGroups assumes a square NxN group of tiles
func (ts TileSet) GetTileGroups() tileGroups {
	result := tileGroups{
		CornerTiles:    make(TileSet),
		PerimeterTiles: make(TileSet),
		InteriorTiles:  make(TileSet),
	}

	if len(ts) == 0 {
		return result
	}

	if len(ts) == 1 {
		t := ts.pickOne()

		// In this case the tile is everything
		result.CornerTiles[t.ID] = *t
		result.InteriorTiles[t.ID] = *t
		result.PerimeterTiles[t.ID] = *t
		return result
	}

	uniqueMinEdges := ts.getMinEdges()
	// Find corner tiles
	for _, t := range ts {
		exteriorCount := 0
		for _, e := range t.Edges {
			minE := min(e.current, e.flipped)
			if uniqueMinEdges[minE] == 1 {
				exteriorCount++
			}
		}
		switch exteriorCount {
		case 1:
			result.PerimeterTiles[t.ID] = t
		case 2:
			result.CornerTiles[t.ID] = t
			result.PerimeterTiles[t.ID] = t
		case 0:
			result.InteriorTiles[t.ID] = t
		default:
			panic(fmt.Sprintf("unexpected edge matches %d", exteriorCount))
		}
	}

	return result
}

func abs(val int) int {
	if val < 0 {
		return -val
	}
	return val
}

func (t Tile) alignTo(other Tile, edge edgeType) {

	// Find the two edges that match
	for i := topEdge; i <= leftEdge; i++ {
		for j := topEdge; j <= leftEdge; j++ {
			if t.Edges[i].min() == other.Edges[j].min() {
				// flip, if necessary
				oppEdge := (j + 2) % 4
				fmt.Printf("tile %d edge match %d=>%d, rotating...\n", t.ID, i, j)
				breaker := 0
				for t.Edges[oppEdge].min() != other.Edges[j].min() {
					fmt.Println("ROTATING")
					t.Rotate90()
					breaker++
					if breaker > 10 {
						break
					}
				}

				fmt.Printf("tile %d looking at flip decision i=%d,j=%d\n", t.ID, i, j)
				if other.Edges[j].current != t.Edges[oppEdge].flipped {
					log.Printf("tile %d needs flipping (i=%d,j=%d)\n", t.ID, i, j)
					switch j {
					case topEdge, bottomEdge:
						t.FlipY()
					case leftEdge, rightEdge:
						t.FlipX()
					}
				}

				return
			}
		}
	}
}

// FindMatchTile finds a tile that matches the specified edge
func (ts TileSet) FindMatchTile(edge Edge) *Tile {
	for _, tile := range ts {
		for _, e := range tile.Edges {
			if e.min() == edge.min() {
				return &tile
			}
		}
	}

	return nil
}

// func (ts TileSet) doSolve(a Arrangement, x, y edgeType, iter int) {
// 	// handle last tile
// 	if len(ts) == 1 {
// 		for _, t := range ts {
// 			a[iter][iter] = t
// 			return
// 		}
// 	}
// 	lastCorner := a[iter][iter]

// 	// NOTE: assumes square puzzle
// 	for x := iter; x < len(a); x++ {
// 		a[iter][x]
// 	}

// 	// column
// 	for y := iter; y < len(a); y++ {

// 	}

// 	// expand out row

// 	// expand out column

// 	// Solve next iteration
// }

func (ts TileSet) Solve() Arrangement {
	var directions []edgeType
	groups := ts.GetTileGroups()

	// Pick a random corner
	corner := groups.CornerTiles.pickOne()
	delete(groups.PerimeterTiles, corner.ID)

	// Find the 2 edge types that we build off of
	for i := topEdge; i < 4; i++ {

		// We will use the edge type if we can find a match
		if groups.PerimeterTiles.FindMatchTile(corner.Edges[i]) != nil {
			directions = append(directions, i)
		}
	}

	if len(directions) != 2 {
		panic(fmt.Sprintf("A corner tile should only have 2 edge directions, found %d\n", len(directions)))
	}

	// Initialize a new, empty arrangement
	a := emptyArrangement(ts)
	fmt.Println("corner", corner.ID)
	a[0][0] = *corner
	delete(ts, corner.ID)
	delete(groups.PerimeterTiles, corner.ID)

	// NOTE: assumes square puzzle
	xDir, yDir := directions[0], directions[1]
	for iter := 0; iter < len(a)-1; iter++ {

		fmt.Println("---ACROSS---")
		for x := iter + 1; x < len(a); x++ {
			last := a[iter][x-1]
			fmt.Println("last", last.ID)
			next := groups.PerimeterTiles.FindMatchTile(last.Edges[xDir])
			if next == nil {
				fmt.Println(a)
				fmt.Println()
				for id := range groups.PerimeterTiles {
					fmt.Printf("%d ", id)
				}
				fmt.Println()
				panic("err")
			}
			fmt.Println("next", next.ID)
			delete(ts, next.ID)
			delete(groups.PerimeterTiles, next.ID)
			a[iter][x] = *next
			a[iter][x].alignTo(last, xDir)
		}

		fmt.Println("---DOWN---")
		for y := iter + 1; y < len(a); y++ {
			last := a[y-1][iter]
			fmt.Println("last", last.ID)
			next := groups.PerimeterTiles.FindMatchTile(last.Edges[yDir])
			fmt.Println("next", next.ID)
			delete(ts, next.ID)
			delete(groups.PerimeterTiles, next.ID)
			a[y][iter] = *next
			a[y][iter].alignTo(last, yDir)
		}
		fmt.Println(a)

		// Next group is a new, smaller square
		groups = ts.GetTileGroups()

		// Set next corner piece
		last := a[iter+1][iter]
		next := groups.CornerTiles.FindMatchTile(last.Edges[xDir])
		a[iter+1][iter+1] = *next
		a[iter+1][iter+1].alignTo(last, xDir)

		delete(ts, next.ID)
		delete(groups.PerimeterTiles, next.ID)
	}

	return a
}

func emptyArrangement(tiles map[int]Tile) Arrangement {
	size := int(math.Sqrt(float64(len(tiles))))
	if size*size != len(tiles) {
		panic("invalid tile arrangement")
	}

	result := make(Arrangement, size)
	for i := 0; i < size; i++ {
		result[i] = make([]Tile, size)
	}
	// i := 0
	// for _, tile := range tiles {
	// 	y := i / size
	// 	if result[y] == nil {
	// 		result[y] = make([]Tile, size)
	// 	}

	// 	x := i % size
	// 	result[y][x] = tile
	// 	i++
	// }

	return result
}

type edgeType int

const (
	topEdge edgeType = iota
	rightEdge
	bottomEdge
	leftEdge
)

// func (a Arrangement) validate() bool {

// 	for y := 1; y < len(a); y += 2 {
// 		for x := 1; x < len(a[y]); x += 2 {
// 			if !a.matchNeighbors(y, x) {
// 				return false
// 			}
// 		}
// 	}

// 	return true
// }

// func (a Arrangement) matchNeighbors(y, x int) bool {

// 	center := a[y][x]

// 	// match top
// 	if y-1 >= 0 && !center.Edges[topEdge].Match(a[y-1][x].Edges[bottomEdge]) {
// 		return false
// 	}

// 	// match right
// 	if len(a[y]) > (x+1) && !center.Edges[rightEdge].Match(a[y][x+1].Edges[leftEdge]) {
// 		return false
// 	}

// 	// match bottom
// 	if len(a) > (y+1) && !center.Edges[bottomEdge].Match(a[y+1][x].Edges[topEdge]) {
// 		return false
// 	}

// 	// match left
// 	return x-1 <= 0 || center.Edges[leftEdge].Match(a[y][x-1].Edges[rightEdge])
// }
