package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/torbensky/adventofcode2020/common"
)

// types of operations
const (
	nopOp = "nop"
	accOp = "acc"
	jmpOp = "jmp"
)

// Represents an instruction
type instruction struct {
	op  string // the type of operation
	arg int    // the argument for the operation
}

func main() {
	program := loadProgram(common.OpenInputFile())

	_, acc := executeProgram(program)
	fmt.Printf("Part 1 - Result %d\n\n", acc)
	fmt.Printf("Part 2 - Total %d\n", fixProgram(program))
}

func fixProgram(instructions []*instruction) int {
	for i, inst := range instructions {
		var completes bool
		var acc int
		switch inst.op {
		case jmpOp:
			completes, acc = testFixedProgram(instructions, i, nopOp, jmpOp)
		case nopOp:
			completes, acc = testFixedProgram(instructions, i, jmpOp, nopOp)
		default:
			continue
		}

		if completes {
			return acc
		}
	}

	return -1
}

func testFixedProgram(instructions []*instruction, fixedIdx int, newOp, oldOp string) (bool, int) {
	instructions[fixedIdx].op = newOp
	looped, acc := executeProgram(instructions)
	instructions[fixedIdx].op = oldOp
	return looped, acc
}

// executes a program, halting if an infinite loop is detected
// returns true/false depending on loop and the value left in the accumulator
func executeProgram(instructions []*instruction) (bool, int) {
	acc := 0
	i := 0
	executed := make(map[int]struct{})
	for i < len(instructions) {

		// Check if we already executed this line
		if _, ok := executed[i]; ok {
			return false, acc // yup - loop alert!
		}
		executed[i] = struct{}{} // remember we executed this line

		// Execute the instructin
		switch instructions[i].op {
		case jmpOp:
			i += instructions[i].arg
			continue
		case accOp:
			acc += instructions[i].arg
			fallthrough //
		default:
			// no-op
			i++
		}
	}

	return true, acc
}

func loadProgram(reader io.Reader) []*instruction {
	var instructions []*instruction
	parseLine := func(line string) {
		fields := strings.Fields(line)
		instruction := &instruction{
			op: fields[0],
		}

		numStr := fields[1]
		if numStr[0] == '+' {
			numStr = numStr[1:]
		}
		val, err := strconv.Atoi(numStr)
		common.MustNotError(err)
		instruction.arg = val
		instructions = append(instructions, instruction)
	}
	common.ScanLines(reader, parseLine)
	return instructions
}
