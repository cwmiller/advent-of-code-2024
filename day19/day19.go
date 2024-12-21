// Day 19: Linen Layout
// https://adventofcode.com/2024/day/19

package main

import (
	"fmt"
	"os"
	"strings"
)

type Input struct {
	patterns []string
	designs  []string
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

	input := parseInput(string(inputContents))

	solve(input)
}

// Solves Part 1 and Part 2
// Part 1 returns how many designs can be made with the given patterns
// Part 2 returns how many different combinations of patterns are possible for all designs
func solve(input Input) {
	successful := 0
	totalCombinations := 0
	cache := make(map[string]int)

	for _, design := range input.designs {
		combinations := solutions(design, input.patterns, cache)

		if combinations > 0 {
			successful++
			totalCombinations += combinations
		}
	}

	fmt.Println("Part 1:", successful)
	fmt.Println("Part 2:", totalCombinations)
}

// Returns the number of ways a design can be made from the given patterns
func solutions(design string, patterns []string, cache map[string]int) int {
	if cached, ok := cache[design]; ok {
		return cached
	}

	// Successful end result, the string has been fully matched
	if design == "" {
		return 1
	}

	// Iterate through each pattern, checking if it matches the start of the design
	// If matched, the pattern is trimmed off the start of the design and re-processed
	total := 0
	for _, pattern := range patterns {
		if strings.HasPrefix(design, pattern) {
			next, _ := strings.CutPrefix(design, pattern)

			total += solutions(next, patterns, cache)
		}
	}

	cache[design] = total

	return total
}

// Retrieve patterns and designs from input text file
func parseInput(input string) Input {
	lines := strings.Split(input, "\n")

	patterns := strings.Fields(strings.Replace(lines[0], ",", " ", -1))
	designs := lines[2:]

	return Input{
		patterns,
		designs,
	}
}
