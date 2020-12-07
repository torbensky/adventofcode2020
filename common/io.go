package common

import (
	"bufio"
	"io"
	"log"
)

// TokenFunc is a callback function used by ScanTokens
type TokenFunc func(token string)

// ScanLines fully scans an input stream, emitting lines as tokens
func ScanLines(reader io.Reader, fn TokenFunc) {
	ScanAllTokens(bufio.NewScanner(reader), fn)
}

// ScanSplit scans a stream, emitting one token at a time
//
// by default, each token is the contents of a single line (a line scanning function)
//
func ScanSplit(reader io.Reader, fn TokenFunc, splitFn bufio.SplitFunc) {
	scanner := bufio.NewScanner(reader)
	if splitFn != nil {
		scanner.Split(splitFn)
	}

	ScanAllTokens(scanner, fn)
}

// ReadStringLines reads all the newline separated lines into a string buffer
func ReadStringLines(reader io.Reader) []string {
	var lines []string
	ScanLines(reader, func(line string) {
		lines = append(lines, line)
	})

	return lines
}

// ScanAllTokens scans every token in the scanner, invoking the callback on each one
// stops when the end of reader is reached
//
func ScanAllTokens(scanner *bufio.Scanner, fn TokenFunc) {
	for scanner.Scan() {
		fn(scanner.Text())
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
	// End of stream, and no data/token left
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
