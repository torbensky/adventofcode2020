package main

import (
	"fmt"
	"io"
	"math"
	"math/big"
	"strconv"
	"strings"

	common "github.com/torbensky/adventofcode-common"
)

func main() {
	fmt.Printf("Part 1: %d\n", part1(common.OpenInputFile()))
	fmt.Printf("Part 2: %d\n", part2(common.OpenInputFile()))
}

func part1(reader io.Reader) int {
	lines := common.ReadStringLines(reader)
	departAt, err := strconv.Atoi(lines[0])
	common.MustNotError(err)

	busses := strings.Split(lines[1], ",")
	var busNums []int
	for _, b := range busses {
		if b == "x" {
			continue
		}
		bn, err := strconv.Atoi(b)
		common.MustNotError(err)
		busNums = append(busNums, bn)
	}

	closestBus := -1
	closestTime := math.MaxInt64
	for _, bn := range busNums {
		time := departAt / bn
		time *= bn
		if time < departAt {
			time += bn
		}
		if time < closestTime {
			closestTime = time
			closestBus = bn
		}
	}

	return closestBus * (closestTime - departAt)
}

func part2(reader io.Reader) int64 {
	lines := common.ReadStringLines(reader)

	busses := strings.Split(lines[1], ",")
	var a []*big.Int
	var n []*big.Int
	for i, b := range busses {
		if b == "x" {
			continue
		}
		busNum, err := strconv.Atoi(b)
		common.MustNotError(err)
		a = append(a, big.NewInt(int64(busNum-i)))
		n = append(n, big.NewInt(int64(busNum)))
	}

	result, err := crt(a, n)
	common.MustNotError(err)
	return result.Int64()

	// for {
	// 	runBus(busList[0])

	// 	for i := 1; i < len(busList); i++ {
	// 		runBusToClosest(busList[i], busList[0].time)
	// 	}

	// 	scheduleWorks := true
	// 	for i := 1; i < len(busList); i++ {
	// 		baseTime := busList[0].time
	// 		timeDiff := busList[i].time - baseTime
	// 		if timeDiff != busList[i].order {
	// 			scheduleWorks = false
	// 			break
	// 		}
	// 	}

	// 	if scheduleWorks {
	// 		return busList[0].time
	// 	}
	// }
}

var one = big.NewInt(1)

// https://rosettacode.org/wiki/Chinese_remainder_theorem#Go
func crt(a, n []*big.Int) (*big.Int, error) {
	p := new(big.Int).Set(n[0])
	for _, n1 := range n[1:] {
		p.Mul(p, n1)
	}
	var x, q, s, z big.Int
	for i, n1 := range n {
		q.Div(p, n1)
		z.GCD(nil, &s, n1, &q)
		if z.Cmp(one) != 0 {
			return nil, fmt.Errorf("%d not coprime", n1)
		}
		x.Add(&x, s.Mul(a[i], s.Mul(&s, &q)))
	}
	return x.Mod(&x, p), nil
}

// func runBusses(busses []*bus) {
// 	for _, bus := range busses {
// 		bus.time += bus.num
// 	}
// }

// func runBus(bus *bus) {
// 	bus.time += bus.num
// }

// func runBusToClosest(bus *bus, time int) {
// 	newTime := time / bus.num
// 	newTime *= bus.num
// 	if newTime < time {
// 		newTime += bus.num
// 	}
// 	bus.time = newTime
// }
