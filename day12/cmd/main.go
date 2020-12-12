package main

import (
	"fmt"
	"io"
	"strconv"

	"github.com/torbensky/adventofcode2020/common"
)

type direction byte

const (
	north direction = 'N'
	east  direction = 'E'
	south direction = 'S'
	west  direction = 'W'
)

type action byte

const (
	flyNorth    action = 'N'
	flyEast     action = 'E'
	flySouth    action = 'S'
	flyWest     action = 'W'
	rotateLeft  action = 'L'
	rotateRight action = 'R'
	moveForward action = 'F'
)

type coord struct {
	n int
	e int
}

func main() {
	fmt.Printf("Part 1: %d\n", part1(common.OpenInputFile()))
	fmt.Printf("Part 2: %d\n", part2(common.OpenInputFile()))
}

func part1(reader io.Reader) int {
	vessel := newShip(false)
	common.ScanLines(reader, func(line string) {
		action := action(line[0])
		val, err := strconv.Atoi(string(line[1:]))
		common.MustNotError(err)
		vessel.fly(part1Pilot, action, val)
	})

	return abs(vessel.position.e) + abs(vessel.position.n)
}

func part2(reader io.Reader) int {
	vessel := newShip(false)
	common.ScanLines(reader, func(line string) {
		a := action(line[0])
		val, err := strconv.Atoi(string(line[1:]))
		common.MustNotError(err)
		vessel.fly(part2Pilot, a, val)
	})

	return abs(vessel.position.e) + abs(vessel.position.n)
}

type pilot func(s *ship, a action, val int)

func part1Pilot(s *ship, a action, val int) {
	switch a {
	case flyNorth:
		s.move(val, 0)
	case flySouth:
		s.move(-val, 0)
	case flyEast:
		s.move(0, val)
	case flyWest:
		s.move(0, -val)
	case rotateLeft:
		s.rotateHeading(true, val)
	case rotateRight:
		s.rotateHeading(false, val)
	case moveForward:
		s.moveForward(val)
	}
}

func part2Pilot(s *ship, a action, val int) {
	switch a {
	case flyNorth:
		s.moveWaypoint(val, 0)
	case flySouth:
		s.moveWaypoint(-val, 0)
	case flyEast:
		s.moveWaypoint(0, val)
	case flyWest:
		s.moveWaypoint(0, -val)
	case rotateLeft:
		// rotate waypoint left
		s.rotateWaypoint(true, val)
	case rotateRight:
		// rotate waypoint east
		s.rotateWaypoint(false, val)
	case moveForward:
		// move ship waypoint amount
		s.moveToWaypoint(val)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type ship struct {
	debug    bool
	heading  direction
	position coord
	waypoint coord
}

func newShip(debug bool) ship {
	return ship{
		debug:   debug,
		heading: east,
		position: coord{
			n: 0,
			e: 0,
		},
		waypoint: coord{
			n: 1,
			e: 10,
		},
	}
}

func (s *ship) move(n, e int) {
	s.position.n += n
	s.position.e += e
}

func (s *ship) moveForward(amount int) {
	switch s.heading {
	case north:
		s.position.n += amount
	case south:
		s.position.n -= amount
	case east:
		s.position.e += amount
	case west:
		s.position.e -= amount
	}
}

func (s *ship) rotateHeading(left bool, degrees int) {
	numRotations := (degrees % 360) / 90
	for i := 0; i < numRotations; i++ {
		s.heading = rotate90(s.heading, left)
	}
}

func rotate90(heading direction, left bool) direction {
	switch heading {
	case north:
		if left {
			return west
		}
		return east
	case south:
		if left {
			return east
		}
		return west
	case west:
		if left {
			return south
		}
		return north
	case east:
		if left {
			return north
		}
		return south
	}
	panic("unhandled rotation case")
}

func (s *ship) print() {
	fmt.Printf("heading=%s north=%d east=%d waypoint[n=%d,e=%d]\n", string(s.heading), s.position.n, s.position.e, s.waypoint.n, s.waypoint.e)
}

func (s *ship) rotateWaypoint(left bool, degrees int) {
	numRotations := (degrees % 360) / 90
	for i := 0; i < numRotations; i++ {
		if left {
			s.waypoint.e, s.waypoint.n = -s.waypoint.n, s.waypoint.e
		} else {
			s.waypoint.e, s.waypoint.n = s.waypoint.n, -s.waypoint.e
		}
	}
}

func (s *ship) fly(p pilot, a action, val int) {
	if s.debug {
		fmt.Printf("%s%d\n", string(a), val)
	}
	p(s, a, val)
	if s.debug {
		s.print()
	}
}

func (s *ship) moveWaypoint(n, e int) {
	s.waypoint.n += n
	s.waypoint.e += e
}

func (s *ship) moveToWaypoint(units int) {
	s.position.e += s.waypoint.e * units
	s.position.n += s.waypoint.n * units
}
