package main

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	common "github.com/torbensky/adventofcode-common"
)

type operator string

// types of operations
const (
	nopOp = operator("nop")
	accOp = operator("acc")
	jmpOp = operator("jmp")
)

// Represents an instruction
type instruction struct {
	op  operator // the type of operation
	arg int      // the argument for the operation
}

func main() {
	program := loadProgram(common.OpenInputFile())

	_, acc := executeProgram(program)
	fmt.Printf("Part 1 - Result %d\n\n", acc)
	fmt.Printf("Part 2 - Result %d\n", fixProgram(program))
}

// Fixes the program according to Part 2
func fixProgram(instructions []instruction) int {
	for i, inst := range instructions {
		var new, old operator
		switch inst.op {
		case jmpOp:
			new, old = nopOp, jmpOp
		case nopOp:
			new, old = jmpOp, nopOp
		default:
			continue
		}

		completes, acc := executePatch(instructions, i, new, old)
		if completes {
			return acc
		}
	}

	return -1
}

// Patches the program and executes that, returning the result
func executePatch(instructions []instruction, patchIdx int, newOp, oldOp operator) (bool, int) {
	instructions[patchIdx].op = newOp
	looped, acc := executeProgram(instructions)
	instructions[patchIdx].op = oldOp
	return looped, acc
}

// executes a program, halting if an infinite loop is detected
// returns true/false depending on whether a loop was found and the value left in the accumulator
func executeProgram(instructions []instruction) (bool, int) {
	acc := 0
	i := 0
	executed := make(map[int]struct{})
	for i < len(instructions) {

		// Check if we already executed this line
		if _, ok := executed[i]; ok {
			return false, acc // yup - loop alert!
		}
		executed[i] = struct{}{} // remember we executed this line

		// Execute the instruction
		switch instructions[i].op {
		case jmpOp:
			i += instructions[i].arg
		case accOp:
			acc += instructions[i].arg
			i++
		case nopOp:
			i++
		default:
			log.Fatalf("unknown instruction %s encountered on line %d", instructions[i].op, i)
		}
	}

	return true, acc
}

// Loads a program from some data stream
func loadProgram(reader io.Reader) []instruction {

	var instructions []instruction

	parseLine := func(line string) {

		fields := strings.Fields(line)
		instruction := instruction{
			// NOTE: in the real world, probably should validate this input. But meh for this :P
			op: operator(fields[0]),
		}

		// parse out the argument
		numStr := fields[1]
		val, err := strconv.Atoi(numStr)
		common.MustNotError(err)
		instruction.arg = val

		instructions = append(instructions, instruction)
	}
	common.ScanLines(reader, parseLine)

	return instructions
}
