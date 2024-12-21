// Day 20: Race Condition
// https://adventofcode.com/2024/day/20

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/albertorestifo/dijkstra"
)

type (
	Point struct {
		x, y int
	}

	Vec struct {
		x, y int
	}

	Tile int

	Maze struct {
		width  int
		height int
		tiles  map[Point]Tile
		start  Point
		end    Point
	}
)

const (
	Wall Tile = iota
	Space
	Start
	End
)

var (
	Up    = Vec{0, -1}
	Right = Vec{1, 0}
	Down  = Vec{0, 1}
	Left  = Vec{-1, 0}
)

func (p Point) Add(v Vec) Point {
	return Point{p.x + v.x, p.y + v.y}
}

func (p Point) String() string {
	return fmt.Sprintf("%d,%d", p.x, p.y)
}

func (v Vec) Clockwise() Vec {
	return Vec{-v.y, v.x}
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

	maze := newMazeFromInput(string(inputContents))

	solve(maze)
}

func solve(maze Maze) {
	baseCost := dijkstraCost(maze)

	fmt.Println("Base Cost:", baseCost)

	cheatables := cheatableWalls(maze)
	cheatableSavings := make(map[int]int)

	for _, cheatable := range cheatables {
		maze.tiles[cheatable] = Space
		cost := dijkstraCost(maze)

		cheatableSavings[baseCost-cost]++

		// Set the cheatable wall back to a plain wall for the next iteration
		maze.tiles[cheatable] = Wall
	}

	part1 := 0

	for savings, count := range cheatableSavings {
		if savings >= 100 {
			part1 += count
		}
	}

	fmt.Println("Part 1:", part1)
}

func dijkstraCost(maze Maze) int {
	graph := dijkstra.Graph{}

	for y := 0; y < maze.height; y++ {
		for x := 0; x < maze.width; x++ {
			point := Point{x, y}

			if tile, ok := maze.tiles[point]; ok {
				if tile != Wall {
					neighbors := neighbors(maze, point)
					neighborMap := make(map[string]int)

					for _, neighbor := range neighbors {
						neighborMap[neighbor.String()] = 1
					}

					graph[point.String()] = neighborMap
				}
			}
		}
	}

	_, cost, _ := graph.Path(maze.start.String(), maze.end.String()) // skipping error handling

	return cost
}

func neighbors(maze Maze, point Point) []Point {
	neighbors := make([]Point, 0)
	dir := Up

	for i := 0; i < 4; i++ {
		target := point.Add(dir)

		if _, ok := maze.tiles[target]; ok {
			neighbors = append(neighbors, target)
		}

		dir = dir.Clockwise()
	}

	return neighbors
}

// Find every point that is cheatable
// A cheatable point is a single wall tile between two walkable paths
func cheatableWalls(maze Maze) []Point {
	cheatables := make(map[Point]struct{})

	for y := 0; y < maze.height; y++ {
		for x := 0; x < maze.width; x++ {
			point := Point{x, y}

			if tile, ok := maze.tiles[point]; ok {
				if tile != Wall {
					dir := Up

					for i := 0; i < 4; i++ {
						target := point.Add(dir)
						beyondTarget := target.Add(dir)

						targetTile, targetOk := maze.tiles[target]
						beyondTargetTile, beyondTargetOk := maze.tiles[beyondTarget]

						if targetOk && beyondTargetOk {
							if targetTile == Wall && beyondTargetTile != Wall {
								cheatables[target] = struct{}{}
							}
						}

						dir = dir.Clockwise()
					}

				}
			}
		}
	}

	cheatablePoints := make([]Point, len(cheatables))
	i := 0
	for key, _ := range cheatables {
		cheatablePoints[i] = key
		i++
	}

	return cheatablePoints
}

func newMazeFromInput(input string) Maze {
	lines := strings.Split(input, "\n")
	height := len(lines)
	width := len(lines[0])
	tiles := make(map[Point]Tile)
	start := Point{}
	end := Point{}

	for y, line := range lines {
		for x, ch := range line {
			point := Point{x, y}
			var tile Tile

			switch ch {
			case '#':
				tile = Wall
			case '.':
				tile = Space
			case 'S':
				tile = Start
				start = point
			case 'E':
				tile = End
				end = point
			}

			tiles[point] = tile
		}
	}

	return Maze{
		width,
		height,
		tiles,
		start,
		end,
	}
}
