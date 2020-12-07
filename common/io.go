package common

import (
	"bufio"
	"log"
)

// StoppableTokenFunc is a callback function used by ScanTokens
//
// return true to continue, false to stop processing further tokens
//
type StoppableTokenFunc func(token string) bool

// AllTokensFunc is an implementation of the StoppableTokenFunc that never wants to stop
//
func AllTokensFunc(fn func(token string)) StoppableTokenFunc {
	return func(token string) bool {
		fn(token)
		return true
	}
}

// ScanTokens scans every token in the scanner, invoking the callback on each one
// stops either when the end of file is reached or when the callback returns false
//
func ScanTokens(scanner *bufio.Scanner, fn StoppableTokenFunc) {
	for scanner.Scan() {
		if !fn(scanner.Text()) {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// SplitRecordsFunc splits on two consecutive, empty lines
// It is an implementation of https://golang.org/pkg/bufio/#SplitFunc
//
// Note: ignores "\r" carriage returns (so "\n\r\n" or even "\n\r\r\r\r...\n" will delimit tokens)
//
func SplitRecordsFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// End of file, and no data/token left
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// Check for token delim
	index := 0
	consecutiveNewLines := 0
	for ; index < len(data); index++ {
		switch data[index] {
		case '\n':
			consecutiveNewLines++
		case '\r':
			// ignore
		default:
			consecutiveNewLines = 0
		}

		// found token delim
		if consecutiveNewLines == 2 {
			// Note: may contain "\r" somewhere
			return index + 1, data[:index-1], nil
		}
	}

	// End of file, remaining data should be a token
	if atEOF {
		return len(data), data, nil
	}

	// Need MOAR
	return 0, nil, nil
}
