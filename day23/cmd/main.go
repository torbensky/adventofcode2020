package main

import (
	"fmt"
	"strconv"
	"strings"
)

// const input = "389125467"

const input = "463528179"

func main() {
	fmt.Println("Part 1:", part1())
	fmt.Println("Part 2", part2())
}

func part1() string {
	ring := parseInput(input)
	ring.playGame(100)
	return ring.part1String()
}

func part2() string {
	ring := parseInput(input)
	fmt.Println("extending...")
	ring.extendTo(1000000)
	fmt.Println("playing...")
	ring.playGame(10000000)
	return ring.part2String()
}

type cupList struct {
	head *cup
	tail *cup
}

func (list cupList) String() string {
	var sb strings.Builder
	iter := list.head
	for {
		sb.WriteString(iter.String())
		iter = iter.next

		if iter == nil {
			return sb.String()
		}
	}
}

func (list cupList) has(v int) bool {
	// Search for the current wanted item
	iter := list.head
	for {
		if iter.value == v {
			return true
		}

		iter = iter.next
		if iter == nil {
			return false
		}
	}
}

type cupRing struct {
	min     int          // min cup value seen
	max     int          // max cup value seen
	current *cup         // current cup in play
	index   map[int]*cup // support O(1) lookup of nodes in the ring
}

func (ring cupRing) last() *cup {
	iter := ring.current
	for iter.next != ring.current {
		iter = iter.next
	}

	return iter
}

func (ring *cupRing) extendTo(n int) {
	last := ring.last()
	val := ring.max + 1
	for val <= n {
		last.next = &cup{
			value: val,
		}
		ring.index[val] = last.next
		last = last.next
		val++
	}
	ring.max = n

	// make it a ring again
	last.next = ring.current
}

func (ring cupRing) takeN(n int) cupList {

	taken := ring.current.next
	list := cupList{head: taken, tail: taken}
	for i := 1; i < n; i++ {
		list.tail = list.tail.next
	}

	// New ring
	ring.current.next = list.tail.next
	list.tail.next = nil

	return list
}

func (ring cupRing) String() string {
	var sb strings.Builder
	iter := ring.current
	for {
		sb.WriteString(iter.String())
		iter = iter.next

		if iter == ring.current {
			return sb.String()
		}
	}
}

type cup struct {
	value int
	next  *cup
}

func (c cup) String() string {
	return strconv.Itoa(c.value)
}

func parseInput(input string) cupRing {
	val := int(input[0]) - '0'
	next := &cup{}
	ring := cupRing{current: &cup{
		value: val,
		next:  next,
	}, min: val, max: val, index: make(map[int]*cup)}
	ring.index[val] = ring.current

	for i := 1; i < len(input); i++ {

		val := int(input[i]) - '0'
		next.value = val
		ring.index[val] = next

		if val > ring.max {
			ring.max = val
		}

		if val < ring.min {
			ring.min = val
		}

		if i < len(input)-1 {
			next.next = &cup{}
			next = next.next
		}
	}

	// close up the ring
	next.next = ring.current

	return ring
}

func (ring cupRing) has(v int) bool {
	// Search for the current wanted item
	iter := ring.current
	for {
		if iter.value == v {
			return true
		}

		iter = iter.next
		if iter == ring.current {
			return false
		}
	}
}

func (ring cupRing) selectDestination(exclude cupList) int {
	want := ring.current.value - 1

	for {
		if want < ring.min {
			want = ring.max
		}

		// We know the ring has it if the excluded list doesn't
		if !exclude.has(want) {
			return want
		}

		want--
	}
}

// Return the product of the two numbers that come after 1
func (ring cupRing) part2String() string {
	for ring.current.value != 1 {
		ring.current = ring.current.next
	}
	n1 := ring.current.next
	n2 := n1.next
	return strconv.Itoa(n1.value * n2.value)
}

// Return the list of numbers that comes after "1"
func (ring cupRing) part1String() string {
	for ring.current.value != 1 {
		ring.current = ring.current.next
	}
	return ring.String()[1:]
}

func (ring cupRing) insertAfter(val int, other cupList) {
	insertAfter := ring.index[val]
	// // assume the value has to be in the list
	// for insertAfter.value != val {
	// 	insertAfter = insertAfter.next
	// }

	// glue in the new list
	oldNext := insertAfter.next
	insertAfter.next = other.head
	other.tail.next = oldNext
}

func (ring *cupRing) playRound() {
	threeCups := ring.takeN(3)
	dest := ring.selectDestination(threeCups)
	// fmt.Println("three", threeCups, "dest", dest, "remaining", ring)
	ring.insertAfter(dest, threeCups)
	ring.current = ring.current.next
}

func (ring *cupRing) playGame(numRounds int) {
	for i := 0; i < numRounds; i++ {
		ring.playRound()
		// fmt.Println(ring)select
	}
}
