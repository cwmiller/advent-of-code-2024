// Day 16: Reindeer Maze
// https://adventofcode.com/2024/day/16

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/albertorestifo/dijkstra"
)

type Point struct {
	x, y int
}

func (p Point) Add(v Vec) Point {
	return Point{p.x + v.x, p.y + v.y}
}

type Vec struct {
	x, y int
}

func (v Vec) Clockwise() Vec {
	return Vec{-v.y, v.x}
}

func (v Vec) CounterClockwise() Vec {
	return Vec{v.y, -v.x}
}

var Up = Vec{0, -1}
var Right = Vec{1, 0}
var Down = Vec{0, 1}
var Left = Vec{-1, 0}

type Maze struct {
	width, height int
	tiles         map[Point]Tile
	startTile     Tile
	endTile       Tile
}

func (m *Maze) AddTile(kind TileKind, pos Point) {
	tile := Tile{kind, pos}
	m.tiles[pos] = tile

	if kind == Start {
		m.startTile = tile
	}

	if kind == End {
		m.endTile = tile
	}
}

func (m *Maze) GetTile(pos Point) (Tile, bool) {
	tile, ok := m.tiles[pos]

	return tile, ok
}

type Neighbor struct {
	tile Tile
	dir  Vec
}

func (m *Maze) GetNeighbors(tile Tile) []Neighbor {
	neighbors := make([]Neighbor, 0)

	dir := Up
	for i := 0; i < 4; i++ {
		neighborPos := tile.pos.Add(dir)
		neighbor := m.tiles[neighborPos]

		if neighbor.kind != Wall {
			neighbors = append(neighbors, Neighbor{neighbor, dir})
		}

		dir = dir.Clockwise()
	}

	return neighbors
}

func (m *Maze) String() string {
	var str string

	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			switch m.tiles[Point{x, y}].kind {
			case Wall:
				str += "#"
			case Path:
				str += "."
			case Start:
				str += "S"
			case End:
				str += "E"
			}
		}

		str += "\n"
	}

	return str
}

func NewMaze(width, height int) *Maze {
	return &Maze{
		width:  width,
		height: height,
		tiles:  make(map[Point]Tile),
	}
}

type TileKind int

const (
	Wall TileKind = iota
	Path
	Start
	End
)

type Tile struct {
	kind TileKind
	pos  Point
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

func newMazeFromInput(input string) *Maze {
	lines := strings.Split(input, "\n")

	height := len(lines)
	width := len(lines[0])
	maze := NewMaze(width, height)

	for y, line := range lines {
		for x, ch := range line {
			var kind TileKind

			switch ch {
			case '#':
				kind = Wall
			case '.':
				kind = Path
			case 'S':
				kind = Start
			case 'E':
				kind = End
			}

			maze.AddTile(kind, Point{x, y})
		}
	}

	return maze
}

type Node struct {
	pos Point
	dir Vec
}

func (n Node) String() string {
	return fmt.Sprintf("%d,%d-%d,%d", n.pos.x, n.pos.y, n.dir.x, n.dir.y)
}

func solve(maze *Maze) {
	// Create graph
	graph := dijkstra.Graph{}

	for y := 0; y < maze.height; y++ {
		for x := 0; x < maze.width; x++ {
			if tile, ok := maze.GetTile(Point{x, y}); ok {
				if tile.kind != Wall {
					dir := Up

					for i := 0; i < 4; i++ {
						node := Node{tile.pos, dir}
						neighborMap := make(map[string]int)

						neighbors := maze.GetNeighbors(tile)

						for _, neighbor := range neighbors {
							neighborNode := Node{neighbor.tile.pos, neighbor.dir}
							dist := 1

							// Add another 1000 distance if we have to rotate
							if dir != neighbor.dir {
								dist += 1000
							}

							neighborMap[neighborNode.String()] = dist
						}

						graph[node.String()] = neighborMap

						dir = dir.Clockwise()
					}
				}
			}
		}
	}

	startNode := Node{maze.startTile.pos, Right}
	endNode := Node{maze.endTile.pos, Up}

	_, cost, _ := graph.Path(startNode.String(), endNode.String()) // skipping error handling

	fmt.Println("Cost:", cost)
}
