// Day 2: Red-Nosed Reports

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if cap(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [inputfile]\n", os.Args[0])
		os.Exit(-1)
	}

	inputFilename := os.Args[1]

	contentBytes, err := os.ReadFile(inputFilename)

	if err != nil {
		panic(err)
	}

	contents := string(contentBytes)

	// Data will be split into a multi-dimensional array
	levels := make([][]int, strings.Count(contents, "\n")+1)

	// Parse each line of input
	// Each line is a level with reports separated by a space
	lines := strings.Split(contents, "\n")

	for i, line := range lines {
		reportStrs := strings.Split(line, " ")
		level := make([]int, cap(reportStrs))

		for j, reportStr := range reportStrs {
			level[j], _ = strconv.Atoi(reportStr)
		}

		levels[i] = level
	}

	run(levels)
}

func run(levels [][]int) {
	totalPart1Safe := 0
	totalPart2Safe := 0

	for levelIdx, level := range levels {
		part1Safe := safe(level)
		part2Safe := part1Safe

		if !part1Safe {
			part2Safe = dampenedSafe(level)
		}

		fmt.Println("Level ", levelIdx, " = ", part1Safe, part2Safe)

		if part1Safe {
			totalPart1Safe++
		}

		if part2Safe {
			totalPart2Safe++
		}
	}

	fmt.Println("Total Safe (Part 1): ", totalPart1Safe)
	fmt.Println("Total Safe (Part 2): ", totalPart2Safe)
}

// Check a level for safety
func safe(level []int) bool {
	return safeWithDampener(level, -1)
}

// Allow a report to be skipped
// Part 2 allows for a single report to be omitted when checking for safeness
func safeWithDampener(level []int, dampenIdx int) bool {
	// We must keep track of the last report to determine the difference between reports and if the values are incrementing or decrementing
	var lastReport *int
	var lastDirection *string

	safe := true

	// Loop through each level and determine if its safe or not
	for reportIdx, report := range level {
		// Allow a report to be skipped
		if dampenIdx == reportIdx {
			continue
		}

		// If this is the first report, then we are always safe
		if lastReport != nil {
			diff := *lastReport - report
			absDiff := diff

			if absDiff < 0 {
				absDiff *= -1
			}

			// Difference must be between 1 and 3 (inclusively)
			if absDiff < 1 || absDiff > 3 {
				safe = false
				break
			}

			// Difference is safe, but we must also be consistently ascending or descening in value
			direction := "ascending"
			if report < *lastReport {
				direction = "descending"
			}

			if lastDirection != nil {
				if direction != *lastDirection {
					safe = false
					break
				}
			}

			lastDirection = &direction
		}

		lastReport = &report
	}

	return safe
}

// Dampened mode (part 2) allows a single report to be excluded from the safe detection
func dampenedSafe(level []int) bool {
	for i := range level {
		if safeWithDampener(level, i) {
			return true
		}
	}

	return false
}
