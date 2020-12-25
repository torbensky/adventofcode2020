package main

import "fmt"

const pk1 = 3469259
const pk2 = 13170438

func main() {
	part1()
}

func part1() {
	result := 1

	// Search for the loop sizes that produce the public keys that we found
	var pk1Loops, pk2Loops int
	for i := 1; ; i++ {
		result *= 7
		result %= 20201227
		if result == pk1 {
			pk1Loops = i
		}
		if result == pk2 {
			pk2Loops = i
		}
		if pk1Loops != 0 && pk2Loops != 0 {
			break
		}
	}

	fmt.Println("pk1Loops", pk1Loops, "pk2loops", pk2Loops)

	ek1 := transform(pk1, pk2Loops)
	ek2 := transform(pk1, pk2Loops)
	fmt.Println(ek1, ek2)
}

func transform(subject, loopSize int) int {
	result := 1

	for i := 0; i < loopSize; i++ {
		result *= subject
		result %= 20201227
	}

	return result
}
