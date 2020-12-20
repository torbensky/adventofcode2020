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

func (img Image) FlipX() {
	for _, row := range img {
		for i := 0; i < len(row)/2; i++ {
			// swap pieces
			row[i], row[len(row)-1-i] = row[len(row)-1-i], row[i]
		}
	}
}

func (img Image) FlipY() {
	for i := 0; i < (len(img) / 2); i++ {
		img[i], img[len(img)-1-i] = img[len(img)-1-i], img[i]
	}
}

func (img Image) Rotate90() {
	size := len(img)

	// in-place matrix rotation
	for x := 0; x < size/2; x++ {
		for y := x; y < size-x-1; y++ {

			// rotate in batches of 4

			temp := img[x][y]

			img[x][y] = img[y][size-1-x]
			img[y][size-1-x] = img[size-1-x][size-1-y]
			img[size-1-x][size-1-y] = img[size-1-y][x]
			img[size-1-y][x] = temp
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
		sb.WriteByte(le[i+1])
		sb.Write(t.Image[i])
		sb.WriteByte(re[i+1])
		sb.WriteByte('\n')
	}
	sb.WriteString(t.Edges[bottomEdge].String())
	sb.WriteByte('\n')

	return sb.String()
}

// FlipY flips the tile over the y-axis
func (t Tile) FlipY() {
	t.Edges[topEdge] = t.Edges[topEdge].flip()
	t.Edges[bottomEdge] = t.Edges[bottomEdge].flip()

	// swap left/right edges
	t.Edges[leftEdge], t.Edges[rightEdge] = t.Edges[rightEdge], t.Edges[leftEdge]

	t.Image.FlipY()
}

// FlipX flips the tile over the x-axis
func (t Tile) FlipX() {
	// swap top/bottom edges
	t.Edges[topEdge], t.Edges[bottomEdge] = t.Edges[bottomEdge], t.Edges[topEdge]

	t.Edges[leftEdge] = t.Edges[leftEdge].flip()
	t.Edges[rightEdge] = t.Edges[rightEdge].flip()

	t.Image.FlipX()
}

// Rotate90N performs a rotation N times
func (t Tile) Rotate90N(rotations int) {
	rotationsNeeded := (4 + rotations) % 4
	fmt.Printf("will rotate %d times\n", rotationsNeeded)
	for i := 0; i < rotationsNeeded; i++ {
		t.Rotate90()
	}
}

// Rotate90 rotates the tile by 90 degrees
func (t Tile) Rotate90() {

	for i, e := range t.Edges {
		fmt.Println(i, e)
	}

	t.Edges[topEdge], t.Edges[rightEdge], t.Edges[bottomEdge], t.Edges[leftEdge] = t.Edges[leftEdge].flip(), t.Edges[topEdge].flip(), t.Edges[rightEdge].flip(), t.Edges[bottomEdge].flip()
	fmt.Println()
	for i, e := range t.Edges {
		fmt.Println(i, e)
	}
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
			newEdge(rows[0]),           // top
			newEdge(re.String()),       // right
			newEdge(rows[len(rows)-1]), // bottom
			newEdge(le.String()),       // left
		},
	}
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
	CornerTiles     TileSet
	OneEdgeOutTiles TileSet
	InteriorTiles   TileSet
}

func (ts TileSet) GetTileGroups() tileGroups {
	result := tileGroups{
		CornerTiles:     make(TileSet),
		OneEdgeOutTiles: make(TileSet),
		InteriorTiles:   make(TileSet),
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
			result.OneEdgeOutTiles[t.ID] = t
		case 2:
			result.CornerTiles[t.ID] = t
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
			e1 := t.Edges[i]
			e2 := other.Edges[j]

			if e1.min() == e2.min() {
				// rotate, if necessary
				if i != j {
					log.Printf("tile %d needs rotations\n", t.ID)
					t.Rotate90N(2 + int(i-j))
				}

				// flip, if necessary
				if !e1.Match(e2) {
					log.Printf("tile %d needs flipping\n", t.ID)
					switch i {
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

func (ts TileSet) matchTile(edge Edge) Tile {
	for _, tile := range ts {
		for _, e := range tile.Edges {
			if e.min() == edge.min() {
				return tile
			}
		}
	}

	panic("no match")
}

func (ts TileSet) Solve() Arrangement {
	// Start with some random arrangement
	a := newArrangement(ts)

	// Find correct corners
	// groups := ts.GetTileGroups()

	// Find correct edges

	// Find next layer in

	// repeat until center

	return a
}

func newArrangement(tiles map[int]Tile) Arrangement {
	size := int(math.Sqrt(float64(len(tiles))))
	if size*size != len(tiles) {
		panic("invalid tile arrangement")
	}

	result := make(Arrangement, size)

	i := 0
	for _, tile := range tiles {
		y := i / size
		if result[y] == nil {
			result[y] = make([]Tile, size)
		}

		x := i % size
		result[y][x] = tile
		i++
	}

	return result
}

type edgeType int

const (
	topEdge edgeType = iota
	rightEdge
	bottomEdge
	leftEdge
)

/*

          TOP

           0
      -----------
	  |			|
	  |			|
   3  |			|  1
	  |			|
	  |			|
  	  -----------
	       2

	     BOTTOM
*/

func (a Arrangement) validate() bool {

	for y := 1; y < len(a); y += 2 {
		for x := 1; x < len(a[y]); x += 2 {
			if !a.matchNeighbors(y, x) {
				return false
			}
		}
	}

	return true
}

func (a Arrangement) matchNeighbors(y, x int) bool {

	center := a[y][x]

	// match top
	if y-1 >= 0 && !center.Edges[topEdge].Match(a[y-1][x].Edges[bottomEdge]) {
		return false
	}

	// match right
	if len(a[y]) > (x+1) && !center.Edges[rightEdge].Match(a[y][x+1].Edges[leftEdge]) {
		return false
	}

	// match bottom
	if len(a) > (y+1) && !center.Edges[bottomEdge].Match(a[y+1][x].Edges[topEdge]) {
		return false
	}

	// match left
	return x-1 <= 0 || center.Edges[leftEdge].Match(a[y][x-1].Edges[rightEdge])
}
