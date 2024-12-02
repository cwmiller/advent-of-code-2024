// Day 1: Historian Hysteria

package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func main() {
	if cap(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s [part 1/2] [inputfile]\n", os.Args[0])
		os.Exit(-1)
	}

	inputFilename := os.Args[2]

	contentBytes, err := os.ReadFile(inputFilename)

	if err != nil {
		panic(err)
	}

	contents := string(contentBytes)

	lefts := make([]int, strings.Count(contents, "\n")+1)
	rights := make([]int, strings.Count(contents, "\n")+1)

	regex, _ := regexp.Compile(`(\d+)   (\d+)`)

	matches := regex.FindAllStringSubmatch(string(contents), -1)

	for i, line := range matches {
		left, _ := strconv.Atoi(line[1])
		right, _ := strconv.Atoi(line[2])

		lefts[i] = left
		rights[i] = right
	}

	part := os.Args[1]

	if part == "1" {
		part1(lefts, rights)
	} else {
		part2(lefts, rights)
	}
}

// Part 1 finds the distance differences between the left and right list
func part1(lefts []int, rights []int) {
	slices.Sort(lefts)
	slices.Sort(rights)

	sumOfDiffs := 0

	for i, left := range lefts {
		right := rights[i]
		diff := left - right

		// If negative, make it positive
		if diff < 0 {
			diff *= -1
		}

		sumOfDiffs += diff

		fmt.Printf("%d - %d = %d\n", left, right, diff)
	}

	fmt.Println("Sum of Diffs: ", sumOfDiffs)
}

// Part 2 finds number of occurrences a number from left appears in right
func part2(lefts []int, rights []int) {
	similarityScore := 0

	for _, left := range lefts {
		// Find how many times `left` appears in the rights list
		rightCnt := 0

		for _, right := range rights {
			if right == left {
				rightCnt++
			}
		}

		fmt.Printf("%d = %d\n", left, rightCnt)

		similarityScore += (left * rightCnt)
	}

	fmt.Println("Sum of Similiary Score: ", similarityScore)
}
