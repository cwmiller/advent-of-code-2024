// Day 13: Claw Contraption
// https://adventofcode.com/2024/day/13

package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type button struct {
	x, y int64
}

type prize struct {
	x, y int64
}

type machine struct {
	a     button
	b     button
	prize prize
}

func main() {
	if cap(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [inputfile]\n", os.Args[0])
		os.Exit(-1)
	}

	machines := machinesFromInput(os.Args[1])
	var total1Tokens int64
	var total2Tokens int64

	for _, machine := range machines {
		total1Tokens += solve(machine, 0)
		total2Tokens += solve(machine, 10000000000000)
	}

	fmt.Println("Part 1 total tokens:", total1Tokens)
	fmt.Println("Part 2 total tokens:", total2Tokens)
}

func machinesFromInput(filename string) []machine {
	bytes, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	rx := regexp.MustCompile(`(?m:Button A: X\+(\d+), Y\+(\d+)\nButton B: X\+(\d+), Y\+(\d+)\nPrize: X=(\d+), Y=(\d+))`)

	// Machines are separated by newlines
	machineContents := strings.Split(string(bytes), "\n\n")

	machines := make([]machine, len(machineContents))

	for i, machineContent := range machineContents {
		matches := rx.FindStringSubmatch(machineContent)

		machines[i] = machine{
			a:     button{atoi(matches[1]), atoi(matches[2])},
			b:     button{atoi(matches[3]), atoi(matches[4])},
			prize: prize{atoi(matches[5]), atoi(matches[6])},
		}
	}

	return machines
}

// Calculates the best solution of button presses
// Returns the number of tokens required to press the buttons needed to get to the prize
// A button = 3 tokens, B button = 1 token
// For Part 2, the prize X and Y get incremented by 10000000000000
func solve(machine machine, prizeIncrement int64) int64 {
	prizeX := machine.prize.x + prizeIncrement
	prizeY := machine.prize.y + prizeIncrement

	// Cramer's Rule
	// AX + BX = PX
	// AY + BY = PY
	coefficient := (machine.a.x * machine.b.y) - (machine.b.x * machine.a.y)
	x := (prizeX * machine.b.y) - (prizeY * machine.b.x)
	y := (prizeY * machine.a.x) - (prizeX * machine.a.y)

	a := x / coefficient
	b := y / coefficient

	// Test if the result actually matches the prize position
	// If a failure, return 0 tokens
	testX := a*machine.a.x + b*machine.b.x
	testY := a*machine.a.y + b*machine.b.y

	if prizeX != testX || prizeY != testY {
		return 0
	}

	return a*3 + b
}

func atoi(str string) int64 {
	n, err := strconv.ParseInt(str, 10, 64)

	if err != nil {
		panic(err)
	}

	return n
}
