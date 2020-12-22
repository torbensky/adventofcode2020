package main

import (
	"fmt"
	"io"

	common "github.com/torbensky/adventofcode-common"
)

type card int

type deck []card

func (d deck) Sig() string {
	return fmt.Sprintf("%v", d)
}

func (d deck) copy() deck {
	c := make(deck, len(d))
	copy(c, d)
	return c
}

func (d deck) drawCard() (card, deck) {
	return d[0], d[1:]
}

func loadDecks(reader io.Reader) (deck, deck) {
	var deck1, deck2 deck

	deck1Start, deck2Start := false, false
	common.ScanLines(reader, func(line string) {
		if line == "" {
			return
		}
		if !deck1Start && line == "Player 1:" {
			deck1Start = true
			return
		}

		if !deck2Start && line == "Player 2:" {
			deck2Start = true
			return
		}

		if deck1Start && !deck2Start {
			deck1 = append(deck1, card(common.Atoi(line)))
			return
		}
		deck2 = append(deck2, card(common.Atoi(line)))
	})

	return deck1, deck2
}

func playRound(d1, d2 deck) (deck, deck) {
	c1, d1 := d1.drawCard()
	c2, d2 := d2.drawCard()

	if c1 > c2 {
		d1 = append(d1, []card{c1, c2}...)
	}

	if c2 > c1 {
		d2 = append(d2, []card{c2, c1}...)
	}

	return d1, d2
}

func main() {
	fmt.Printf("Part 1: %d\n", part1(common.OpenInputFile()))
	fmt.Printf("Part 2: %d\n", part2(common.OpenInputFile()))
}

func pickWinner(d1, d2 deck) deck {
	winner := d2
	if len(d1) > len(d2) {
		winner = d1
	}
	return winner
}

func part1(reader io.Reader) int {
	d1, d2 := loadDecks(reader)
	// fmt.Println(d1, d2)

	for len(d1) > 0 && len(d2) > 0 {
		d1, d2 = playRound(d1, d2)
	}
	// fmt.Println(d1, d2)
	winner := pickWinner(d1, d2)
	return calcDec(winner)
}

func calcDec(d deck) int {
	total := 0
	for i := len(d); i > 0; i-- {
		total += (len(d) + 1 - i) * int(d[i-1])
	}

	return total
}

type gameKey struct {
	d1 string
	d2 string
}

func doPlayRecursive(d1, d2 deck, game int) (deck, deck) {

	cache := make(map[gameKey]struct{})

	// round := 0
	for len(d1) > 0 && len(d2) > 0 {
		// round++
		// fmt.Println("game", game, "round", round)
		// fmt.Println("\tdeck1", d1)
		// fmt.Println("\tdeck2", d2)

		// Check if we have already seen this game sequence
		key := gameKey{d1: d1.Sig(), d2: d2.Sig()}
		// fmt.Printf("CACHE %v\n", key)
		if _, ok := cache[key]; ok {
			// fmt.Printf("GAME IN CACHE %v\n", key)
			// We have, player 1 (deck 1) auto wins
			return d1, deck{}
		}

		c1, nd1 := d1.drawCard()
		c2, nd2 := d2.drawCard()
		d1, d2 = nd1, nd2
		// fmt.Printf("\tdrew %d %d\n", c1, c2)

		// Check if we must begin a new game of recursive combat
		if int(c1) <= len(d1) && int(c2) <= len(d2) {
			_, nd2 := doPlayRecursive(d1[0:c1].copy(), d2[0:c2].copy(), game+1)
			if len(nd2) == 0 {
				// Player 1 won, their card is first
				d1 = append(d1, []card{c1, c2}...)
			} else {
				// Player 2 won, their card is first
				d2 = append(d2, []card{c2, c1}...)
			}
			continue
		}

		if c1 > c2 {
			d1 = append(d1, []card{c1, c2}...)
		}

		if c2 > c1 {
			d2 = append(d2, []card{c2, c1}...)
		}

		// Remember this game sequence
		cache[key] = struct{}{}
	}

	return d1, d2
}

func playRecursive(d1, d2 deck) (deck, deck) {
	return doPlayRecursive(d1, d2, 1)
}

func part2(reader io.Reader) int {
	d1, d2 := loadDecks(reader)
	// playRecursive()
	// fmt.Println(d1)
	// fmt.Println(d2)

	d1, d2 = playRecursive(d1, d2)
	winner := pickWinner(d1, d2)
	// fmt.Println("winner", winner)
	result := calcDec(winner)

	return result
}
