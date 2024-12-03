// Day 3: Mull It Over

package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	if cap(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [inputfile]\n", os.Args[0])
		os.Exit(-1)
	}

	contents := readInputfile(os.Args[1])

	// The sum of the results of all mul() instructions in the input
	results := 0

	// For part 2, only the results of mul() instructions following do() instructions count
	conditionalResults := 0

	// mul instructions are in the format "mul(X,Y) where and Y are 1-3 digit numbers
	mulRx, _ := regexp.Compile(`mul\((\d{1,3}),(\d{1,3})\)`)

	// Besides the mul() instruction, there's also do() and don't() commands that enable & disable the mul instruction for Part 2
	allCommandsRx, _ := regexp.Compile(`do\(\)|don't\(\)|mul\((\d{1,3}),(\d{1,3})\)`)

	// Find the start & end indexes of all commands within the input
	commands := allCommandsRx.FindAllStringIndex(string(contents), -1)

	// As we iterate through all commands found, mulEnabled will be toggled as it finds do() and don't() commands
	// mul() commands only apply in Part 2 if a do() command proceeded it
	mulEnabled := true

	for _, indexes := range commands {
		startIdx, endIdx := indexes[0], indexes[1]
		command := contents[startIdx:endIdx]

		switch command {
		case "do()":
			mulEnabled = true

		case "don't()":
			mulEnabled = false

		default:
			// Must be mul() command, parse it to get the operands
			mul := mulRx.FindStringSubmatch(command)
			x, _ := strconv.Atoi(mul[1])
			y, _ := strconv.Atoi(mul[2])

			result := x * y
			results += result

			if mulEnabled {
				conditionalResults += result
			}
		}
	}

	fmt.Println("Part 1 result:", results)
	fmt.Println("Part 2 result:", conditionalResults)
}

func readInputfile(filename string) string {
	contentBytes, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	return string(contentBytes)
}
