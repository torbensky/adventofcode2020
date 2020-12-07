package common

import (
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

// OpenInputFile opens the default input file specified by process args
func OpenInputFile() *os.File {
	return openFile(GetInputFilePath())
}
