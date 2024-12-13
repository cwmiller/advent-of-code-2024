// Day 12: Garden Groups
// https://adventofcode.com/2024/day/12

package main

import (
	"fmt"
	"os"
)

type point struct {
	x, y int
}

func (p point) Add(v vec) point {
	return point{p.x + v.x, p.y + v.y}
}

type vec struct {
	x, y int
}

func (v vec) Rotate() vec {
	return vec{-v.y, v.x}
}

var Up = vec{0, -1}
var Right = vec{1, 0}
var Down = vec{0, 1}
var Left = vec{-1, 0}

var UpLeft = vec{-1, -1}
var UpRight = vec{1, -1}
var DownRight = vec{1, 1}
var DownLeft = vec{-1, -1}

// Puzzle input is a 2-d map of characters
// Each character represents a different type of plant
type farmMap struct {
	width, height int
	plants        map[point]rune
}

// Track the number of plots a plant takes up as well as its edges and sides
type region struct {
	area  int
	edges int
	sides int
}

func main() {
	if cap(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [inputfile]\n", os.Args[0])
		os.Exit(-1)
	}

	farm := newMapFromInput(os.Args[1])

	regions := findRegions(farm)
	totalEdgePrice := 0
	totalSidePrice := 0

	for _, region := range regions {
		totalEdgePrice += region.area * region.edges
		totalSidePrice += region.area * region.sides
	}

	fmt.Println("Part 1 Edge Price:", totalEdgePrice)
	fmt.Println("Part 2 Side Price:", totalSidePrice)
}

func newMapFromInput(filename string) farmMap {
	bytes, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	y := 0
	x := 0

	plants := make(map[point]rune)

	for _, b := range bytes {
		switch b {
		case '\n':
			y++
			x = 0
		default:
			plants[point{x, y}] = rune(b)
			x++
		}
	}

	if x != 0 {
		y++
	}

	return farmMap{
		width:  x,
		height: y,
		plants: plants,
	}
}

func findRegions(farm farmMap) []region {
	mapped := make(map[point]bool)
	regions := make([]region, 0)

	for y := 0; y < farm.height; y++ {
		for x := 0; x < farm.width; x++ {
			pt := point{x, y}

			if _, ok := mapped[pt]; !ok {
				region := walkRegion(farm, pt, mapped)
				regions = append(regions, region)
			}
		}
	}

	return regions
}

func walkRegion(farm farmMap, pt point, mappings map[point]bool) region {
	mappings[pt] = true
	plantType := farm.plants[pt]

	area := 1
	edges := 0
	sides := countSides(farm, pt)

	// Check all directions around this plot to see if the same plant is adjacent
	// Hitting the edge of the farm or a different plant type will count as an edge
	// If an adjacent plot contains the same plant type, then add its area and edges to our own
	dir := Up

	for i := 0; i < 4; i++ {
		adjacentPt := pt.Add(dir)

		// Check if point is outside bounds or if it's a different plant
		adjacentPlantType, inBounds := farm.plants[adjacentPt]

		if !inBounds || adjacentPlantType != plantType {
			edges++
		} else {
			if !mappings[adjacentPt] {
				adjacentRegion := walkRegion(farm, adjacentPt, mappings)

				area += adjacentRegion.area
				edges += adjacentRegion.edges
				sides += adjacentRegion.sides
			}
		}

		dir = dir.Rotate()
	}

	return region{area, edges, sides}
}

// Sides are continuous edges along line of plots of the same plant
func countSides(farm farmMap, pt point) int {
	plantType := farm.plants[pt]

	// Sides is equal to the number of corners
	// Rotate around the given point checking for corners
	dir1 := Left
	dir2 := Up
	diag := UpLeft

	sides := 0

	for i := 0; i < 4; i++ {
		dir1PlantType := farm.plants[pt.Add(dir1)]
		dir2PlantType := farm.plants[pt.Add(dir2)]
		diagPlantType := farm.plants[pt.Add(diag)]

		// Plot has different plants on both corner sides
		if dir1PlantType != plantType && dir2PlantType != plantType {
			sides++
			// Plot has the same plants on both corner sides, but something different diagonally
		} else if dir1PlantType == plantType && dir2PlantType == plantType && diagPlantType != plantType {
			sides++
		}

		dir1 = dir1.Rotate()
		dir2 = dir2.Rotate()
		diag = diag.Rotate()
	}

	return sides
}
