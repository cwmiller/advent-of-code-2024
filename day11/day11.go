// Day 11: Plutonian Pebbles
// https://adventofcode.com/2024/day/11

package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

// Track the results of each blink of a stone
// This is used in a map to reduce the number of duplicate computations needed
type cache map[cacheEntry]int

type cacheEntry struct {
	stone     int64
	iteration int
}

func main() {
	if cap(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [inputfile]\n", os.Args[0])
		os.Exit(-1)
	}

	stones := readInputFile(os.Args[1])

	part1Total := 0
	part2Total := 0
	cache := make(cache)

	// Iterate over each stone and calculate the number of stones remaining after that stone is blinked X times
	// Add this to the overall totals
	// The cache is shared amongst both parts, so all totals have already been calculated for the first 25 blinks when it starts part 2
	for _, stone := range stones {
		part1Total += blink(stone, 25, cache)
		part2Total += blink(stone, 75, cache)
	}

	fmt.Println("Part 1:", part1Total)
	fmt.Println("Part 2:", part2Total)
}

// Read input file and return a list of stones
// Input is just a series of numbers separated by spaces
// Each number represents a stone with a number engraved on it
func readInputFile(filename string) []int64 {
	inputBytes, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	rx := regexp.MustCompile(`(\d+)`)
	matches := rx.FindAllStringSubmatch(string(inputBytes), -1)

	stones := make([]int64, len(matches))

	for i, match := range matches {
		stones[i], _ = strconv.ParseInt(match[1], 10, 64)
	}

	return stones
}

// "Blink" a stone a number of times and return the number of stones left afterwards
// This is a recursive function
func blink(stone int64, blinks int, cache cache) int {
	if blinks == 0 {
		return 1
	}

	cacheEntry := cacheEntry{stone, blinks}
	var total int

	// If the result for this stone & blinks has already been calculated, return it
	if total, ok := cache[cacheEntry]; ok {
		return total
	}

	stoneStr := fmt.Sprintf("%d", stone)

	// If the stone is 0, then it becomes 1
	if stone == 0 {
		total += blink(1, blinks-1, cache)
		// If the stone is an even number of digits, then the number is cut in half and becomes two stones
	} else if len(stoneStr)%2 == 0 {
		// Split stone number in half
		halfIdx := len(stoneStr) / 2

		l, _ := strconv.ParseInt(stoneStr[:halfIdx], 10, 64)
		r, _ := strconv.ParseInt(stoneStr[halfIdx:], 10, 64)

		total += (blink(l, blinks-1, cache) + blink(r, blinks-1, cache))
		// If the previous rules don't apply, then the stone number is multiplied by 2024
	} else {
		total += blink(stone*2024, blinks-1, cache)
	}

	// Cache the result so we don't have to calculate it again in the future
	cache[cacheEntry] = total

	return total
}
