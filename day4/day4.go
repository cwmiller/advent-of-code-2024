// Day 4: Ceres Search

package main

import (
	"errors"
	"fmt"
	"os"
)

type Point struct {
	x, y int
}

type Vec struct {
	x, y int
}

type Grid [][]byte

func (g *Grid) Get(point Point) (byte, error) {
	if point.y < 0 || len(*g) <= point.y || point.x < 0 || len((*g)[point.y]) <= point.x {
		return 0, errors.New("invalid coordinates")
	}

	return (*g)[point.y][point.x], nil
}

func (g *Grid) AddRow(values []byte) {
	*g = append(*g, values)
}

// Returns a len-length word starting from the given point
func (g *Grid) Word(start Point, dir Vec, len int) string {
	//x, y := start.x, start.y
	point := Point{start.x, start.y}
	word := ""

	for i := 0; i < len; i++ {
		val, err := g.Get(point)

		if err == nil {
			word += string(val)
		}

		point.x += dir.x
		point.y += dir.y

	}

	return word
}

// Returns a slice of len long words that originate at the given point
func (g *Grid) AllWords(point Point, len int) []string {
	words := make([]string, 8)

	// Grab all words going clockwise from the starting point
	words[0] = g.Word(point, Vec{0, -1}, len)
	words[1] = g.Word(point, Vec{1, -1}, len)
	words[2] = g.Word(point, Vec{1, 0}, len)
	words[3] = g.Word(point, Vec{1, 1}, len)
	words[4] = g.Word(point, Vec{0, 1}, len)
	words[5] = g.Word(point, Vec{-1, 1}, len)
	words[6] = g.Word(point, Vec{-1, 0}, len)
	words[7] = g.Word(point, Vec{-1, -1}, len)

	return words
}

// Returns a slice of len long words that are in an X pattern from the center point given
func (g *Grid) XWords(center Point, len int) []string {
	words := make([]string, 4)

	// Top-left
	words[0] = g.Word(Point{center.x - (len / 2), center.y - (len / 2)}, Vec{1, 1}, len)
	// Top-right
	words[1] = g.Word(Point{center.x + (len / 2), center.y - (len / 2)}, Vec{-1, 1}, len)
	// Bottom-left
	words[2] = g.Word(Point{center.x - (len / 2), center.y + (len / 2)}, Vec{1, -1}, len)
	// Bottom-right
	words[3] = g.Word(Point{center.x + (len / 2), center.y + (len / 2)}, Vec{-1, -1}, len)

	return words
}

func main() {
	if cap(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [inputfile]\n", os.Args[0])
		os.Exit(-1)
	}

	grid := readFileIntoGrid(os.Args[1])

	fmt.Println("XMAS:", xmasCount(grid))
	fmt.Println("X-MAS:", crossMasCount(grid))
}

func readFileIntoGrid(filename string) *Grid {
	content, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	grid := make(Grid, 0)
	row := make([]byte, 0)

	for _, val := range content {
		if val == '\n' {
			grid.AddRow(row)
			row = make([]byte, 0)
		} else {
			row = append(row, val)
		}
	}

	if len(row) > 0 {
		grid.AddRow(row)
	}

	return &grid
}

// Get number of instances "XMAS" appears in the grid
func xmasCount(grid *Grid) int {
	starterPoints := make([]Point, 0)

	for y := range *grid {
		for x := range (*grid)[y] {
			point := Point{x, y}

			val, err := grid.Get(point)

			if err == nil && val == 'X' {
				starterPoints = append(starterPoints, point)
			}
		}
	}

	count := 0

	for _, point := range starterPoints {
		words := grid.AllWords(point, 4)

		for _, word := range words {
			if word == "XMAS" {
				count++
			}
		}
	}

	return count
}

// Get number of instances an X format that includes two "MAS" occurs
func crossMasCount(grid *Grid) int {
	centerPoints := make([]Point, 0)

	for y := range *grid {
		for x := range (*grid)[y] {
			point := Point{x, y}

			val, err := grid.Get(point)

			if err == nil && val == 'A' {
				centerPoints = append(centerPoints, point)
			}
		}
	}

	found := 0

	for _, point := range centerPoints {
		words := grid.XWords(point, 3)
		wordCount := 0

		for _, word := range words {
			if word == "MAS" {
				wordCount++
			}
		}

		if wordCount >= 2 {
			found++
		}

	}

	return found
}
