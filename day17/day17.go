// Day 17: Chronospatial Computer
// https://adventofcode.com/2024/day/17

package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/cwmiller/advent-of-code-2024/day17/cpu"
)

type monitor struct {
	vals []string
	i    int
}

func (mon *monitor) output(val int) {
	mon.vals[mon.i] = strconv.Itoa(val)

	mon.i++
}

func (mon *monitor) value() string {
	return strings.Join(mon.vals[:mon.i], ",")
}

func (mon *monitor) reset() {
	mon.i = 0
}

func newMonitor() *monitor {
	return &monitor{
		vals: make([]string, 16),
	}
}

func main() {
	if cap(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [input-file]\n", os.Args[0])
		os.Exit(-1)
	}

	inputContents, err := os.ReadFile(os.Args[1])

	if err != nil {
		panic(err)
	}

	solvePart1(string(inputContents))
	solvePart2(string(inputContents))
}

func solvePart1(input string) {
	mon := newMonitor()
	parsedInput := parseInput(input)
	cpu := cpu.NewCpu(parsedInput.rom, mon.output)

	cpu.SetA(parsedInput.a)
	cpu.SetB(parsedInput.b)
	cpu.SetC(parsedInput.c)

	for !cpu.Halted() {
		cpu.Step()
	}

	fmt.Println("Output:", mon.value())
}

func solvePart2(input string) {
}

type FromInput struct {
	a, b, c int
	program string
	rom     cpu.Rom
}

func parseInput(input string) FromInput {
	rx := regexp.MustCompile(`(Register \D|Program): ([0-9,]+)`)
	matches := rx.FindAllStringSubmatch(input, -1)

	var a, b, c int
	var program string
	rom := make(cpu.Rom, 0)

	for _, match := range matches {
		switch match[1] {
		case "Register A":
			a, _ = strconv.Atoi(match[2])
		case "Register B":
			b, _ = strconv.Atoi(match[2])
		case "Register C":
			c, _ = strconv.Atoi(match[2])
		case "Program":
			program = match[2]
			vals := strings.Split(match[2], ",")

			for _, val := range vals {
				val, _ := strconv.Atoi(val)
				rom = append(rom, val)
			}
		}
	}

	return FromInput{
		a,
		b,
		c,
		program,
		rom,
	}
}

/*
func newCpuFromInput(input string, output cpu.OutputHandler) (*cpu.Cpu, {
	rx := regexp.MustCompile(`(Register \D|Program): ([0-9,]+)`)
	matches := rx.FindAllStringSubmatch(input, -1)

	var a, b, c int
	rom := make(cpu.Rom, 0)

	for _, match := range matches {
		switch match[1] {
		case "Register A":
			a, _ = strconv.Atoi(match[2])
		case "Register B":
			b, _ = strconv.Atoi(match[2])
		case "Register C":
			c, _ = strconv.Atoi(match[2])
		case "Program":
			vals := strings.Split(match[2], ",")

			for _, val := range vals {
				val, _ := strconv.Atoi(val)
				rom = append(rom, val)
			}
		}
	}

	return cpu.NewCpu(rom, output)
}
*/
