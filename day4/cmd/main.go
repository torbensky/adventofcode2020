package main

import (
	"bufio"
	"log"
	"os"
)

// Buffer to store each line of a file
var fileLines []string

// Loads entire file into memory
func loadFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var line string
	for scanner.Scan() {
		line = scanner.Text()
		fileLines = append(fileLines, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	// Validate program usage
	if len(os.Args) != 2 {
		log.Fatal("This command accepts only one argument: the path to the input file")
	}
	loadFile(os.Args[1])

	// Part 1

	// TODO

	// Part 2

	// TODO
}
