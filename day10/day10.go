// Day 10: Hoof It
// https://adventofcode.com/2024/day/10

package main

import (
	"fmt"
	"os"
	"strconv"
)

// Trails begin with a height of 0 and increment by one until reaching a 9
const (
	TrailHead = 0
	TrailTail = 9
)

// A x,y coordinate on the topography map
type point struct {
	x, y int
}

// A vector direction
type vec struct {
	x, y int
}

// Topography map contains all x,y positions of the map and the height of each position
type topographyMap struct {
	width, height int
	positions     map[point]int
}

func main() {
	if cap(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [inputfile]\n", os.Args[0])
		os.Exit(-1)
	}

	tmap := newMapFromInput(os.Args[1])

	// Score is the sum of all trail scores
	score := 0

	// Rating is the sum of all trail ratings
	rating := 0

	// Loop over all positions in the typography map looking for trailheads (height of 0)
	// Calculate the score and rating for that trailhead and add to the sums
	for y := 0; y < tmap.height; y++ {
		for x := 0; x < tmap.width; x++ {
			pt := point{x, y}
			pos := tmap.positions[pt]

			if pos == TrailHead {
				score += trailheadScore(tmap, pt)
				rating += completableTrails(tmap, pt)
			}
		}
	}

	fmt.Println("Part 1 Score:", score)
	fmt.Println("Part 2 Rating:", rating)
}

// Parse input text file and generate a topography map from it
// Example:
// 89010123
// 78121874
// 87430965
// 96549874
// 45678903
// 32019012
// 01329801
// 10456732
func newMapFromInput(filename string) topographyMap {
	bytes, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	y := 0
	x := 0

	positions := make(map[point]int)

	for _, b := range bytes {
		switch b {
		case '\n':
			y++
			x = 0
		default:
			positions[point{x, y}], _ = strconv.Atoi(string(b))
			x++
		}
	}

	if x != 0 {
		y++
	}

	return topographyMap{
		width:     x,
		height:    y,
		positions: positions,
	}
}

// Trailhead score is the unique count of trailtails that are accesible from a single trailhead
func trailheadScore(tmap topographyMap, pt point) int {
	// Find all tails that can be reached starting from this point
	tails := completableTails(tmap, pt)

	// Filter the list down to just unique points
	founds := make(map[point]bool, 0)
	uniqueTails := make([]point, 0)

	for _, tailPt := range tails {
		if !founds[tailPt] {
			uniqueTails = append(uniqueTails, tailPt)
			founds[tailPt] = true
		}
	}

	return len(uniqueTails)
}

// Returns the number of unique trails from the given point that ultimately reach a tail
// This is known as the trail rating
func completableTrails(tmap topographyMap, pt point) int {
	height := tmap.positions[pt]
	count := 0

	// Rotate around current position looking for heights 1 greater than the current point, or 9 which is the trail tail
	dir := vec{0, -1}

	for i := 0; i < 4; i++ {
		candidatePt := point{pt.x + dir.x, pt.y + dir.y}

		if candidate, ok := tmap.positions[candidatePt]; ok {
			if candidate == TrailTail && height == TrailTail-1 {
				// A tail has been found! Add to the count
				count++
			} else if candidate == height+1 {
				// A node has been found with a height 1 greater than our own, which means we can walk there!
				// Do a recursive call to walk there and keep following the trail!
				count += completableTrails(tmap, candidatePt)
			}
		}

		dir = vec{-dir.y, dir.x}
	}

	return count
}

// Returns a list of tails that can be reached from the given point
// Duplicate tails can be returned if a trail forks and both paths end up at the same tail
func completableTails(tmap topographyMap, pt point) []point {
	height := tmap.positions[pt]
	tails := make([]point, 0)

	// Rotate around current position looking for heights 1 greater than the current point, or 9 which is the trail tail
	dir := vec{0, -1}

	for i := 0; i < 4; i++ {
		candidatePt := point{pt.x + dir.x, pt.y + dir.y}

		if candidate, ok := tmap.positions[candidatePt]; ok {
			if candidate == TrailTail && height == TrailTail-1 {
				// A tail has been found! Add to the list
				tails = append(tails, candidatePt)
			} else if candidate == height+1 {
				// A node has been found with a height 1 greater than our own, which means we can walk there!
				// Do a recursive call to walk there and keep following the trail!
				tails = append(tails, completableTails(tmap, candidatePt)...)
			}
		}

		dir = vec{-dir.y, dir.x}
	}

	return tails
}
