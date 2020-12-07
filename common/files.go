package common

import (
	"bufio"
	"log"
	"os"
)

func openFile(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

// GetInputFilePath reads the path to the input file from the process args
//
// It expects the process to be invoked with a single argument (that being the input file path)
//
// It fatally exits if args do not match expectations
//
func GetInputFilePath() string {

	// Validate program usage

	if len(os.Args) != 2 {
		log.Fatal("This command accepts only one argument: the path to the input file")
	}

	return os.Args[1]
}

// ScanLines is a shorter overload of ScanFile
//
// It scans each line in a file
//
func ScanLines(path string, fn StoppableTokenFunc) {
	ScanFile(path, fn, nil)
}

// ScanFile scans a file, emitting one token at a time
//
// by default, each token is the contents of a single line (a line scanning function)
//
func ScanFile(path string, fn StoppableTokenFunc, splitFn bufio.SplitFunc) {
	file := openFile(path)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if splitFn != nil {
		scanner.Split(splitFn)
	}

	ScanTokens(scanner, fn)
}

// ReadStringLines reads all the newline separated lines into a string buffer
func ReadStringLines(path string) []string {
	var lines []string
	ScanLines(path, AllTokensFunc(func(line string) {
		lines = append(lines, line)
	}))

	return lines
}
