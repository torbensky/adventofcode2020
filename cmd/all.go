package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	common "github.com/torbensky/adventofcode-common"
)

func main() {
	baseDir := common.GetInputFilePath()
	files, err := ioutil.ReadDir(baseDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() && strings.HasPrefix(f.Name(), "day") {

			fmt.Println("==========================================================")
			fmt.Printf("= Running %s...\n", f.Name())
			fmt.Println("==========================================================")
			fmt.Println()

			mainFile := filepath.Join(baseDir, f.Name(), "cmd", "main.go")
			inputFile := filepath.Join(baseDir, f.Name(), "input.txt")
			cmd := exec.Command("go", "run", mainFile, inputFile)

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			// Execute the command
			if err := cmd.Run(); err != nil {
				log.Panic(err)
			}
		}
	}
}
