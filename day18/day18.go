// Day 18: RAM Run
// https://adventofcode.com/2024/day/18

package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/albertorestifo/dijkstra"
	"github.com/gbin/goncurses"
)

type Point struct {
	x, y int
}

func (p Point) Add(v Vec) Point {
	return Point{p.x + v.x, p.y + v.y}
}

// Convert point to a string representation
// Used to register the point with the dijkstra library
func (p Point) String() string {
	return fmt.Sprintf("%d,%d", p.x, p.y)
}

var Up = Vec{0, -1}
var Right = Vec{1, 0}
var Down = Vec{0, 1}
var Left = Vec{-1, 0}

type Vec struct {
	x, y int
}

func (v Vec) Clockwise() Vec {
	return Vec{-v.y, v.x}
}

// Holds the bytes that have dropped in
type Ram map[Point]bool

// Holds all the byte points from the input file
type Input []Point

func main() {
	if cap(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s width height [input-file]\n", os.Args[0])
		os.Exit(-1)
	}

	width, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	height, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	inputContents, err := os.ReadFile(os.Args[3])
	if err != nil {
		panic(err)
	}

	corruptions := parseInput(string(inputContents))
	ram := make(Ram)

	solve(corruptions, ram, width, height)
}

// Drops each byte one by one from input into ram and finds the best path from start to end
// Each iteration is displayed showing the found path
// Stops once it reaches a byte that blocks any access to the end point
func solve(input Input, ram Ram, width, height int) {
	stdscr, _ := goncurses.Init()
	defer goncurses.End()

	// Part 1 requires finding the cost of the 1024th byte
	kbCost := "Pending"

	for i, n := range input {
		stdscr.Clear()

		ram[n] = true

		paths, cost, _ := pathfind(ram, width, height)

		if i == 1023 {
			kbCost = strconv.Itoa(cost)
		}

		stdscr.Println("Byte:", i+1, " Point:", n, " Cost:", cost)
		stdscr.Println("1024 Cost:", kbCost)
		stdscr.Println(strings.Repeat("-", width))

		// Display map
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				pt := Point{x, y}

				// Points within the found path to the end are displayed as an O
				// Corrupted spaces a #
				// Walkable spaces a .
				if slices.Index(paths, pt.String()) > -1 {
					stdscr.Print("O")
				} else {

					if ram[pt] {
						stdscr.Print("#")
					} else {
						stdscr.Print(".")
					}
				}
			}

			stdscr.Print("\n")
		}

		stdscr.Refresh()

		if cost == 0 {
			stdscr.GetChar()
		}
	}
}

// Convert input file of points to a slice
func parseInput(input string) []Point {
	points := make([]Point, 0)
	rx := regexp.MustCompile(`(\d+),(\d+)`)

	matches := rx.FindAllStringSubmatch(input, -1)
	for _, match := range matches {
		x, _ := strconv.Atoi(match[1])
		y, _ := strconv.Atoi(match[2])

		points = append(points, Point{x, y})
	}

	return points
}

// Perform dijkstra of the RAM to find the best path around the corrupted bytes
func pathfind(ram Ram, width, height int) ([]string, int, error) {
	graph := dijkstra.Graph{}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pt := Point{x, y}

			if ram[pt] {
				continue
			}

			neighbors := make(map[string]int)

			dir := Up
			for i := 0; i < 4; i++ {
				dir = dir.Clockwise()
				neighbor := pt.Add(dir)

				// Don't add any neighbor that's outside the RAM grid
				// Or is corrupted by a byte
				if neighbor.x < 0 || neighbor.y < 0 {
					continue
				}

				if neighbor.x >= width || neighbor.y >= height {
					continue
				}

				if ram[neighbor] {
					continue
				}

				// There's no difference in cost between any neighbor
				neighbors[neighbor.String()] = 1
			}

			graph[pt.String()] = neighbors
		}
	}

	start := Point{0, 0}
	end := Point{width - 1, height - 1}

	path, cost, err := graph.Path(start.String(), end.String())

	return path, cost, err
}
