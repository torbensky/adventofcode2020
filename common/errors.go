package common

import "log"

// MustNotError accepts an error. If it's not nil, it will fatally log and exit the program.
func MustNotError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
