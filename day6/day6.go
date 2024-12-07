// Day 6: Guard Gallivant

package main

import (
	"fmt"
	"os"
)

type Point struct {
	x, y int
}

func (p Point) Add(v Vec) Point {
	return Point{x: p.x + v.x, y: p.y + v.y}
}

type Vec struct {
	x, y int
}

func (v Vec) Rotate() Vec {
	return Vec{x: -v.y, y: v.x}
}

type Visit struct {
	point     Point
	direction Vec
}

type Map struct {
	obstacles         map[Point]bool
	width, height     int
	startingPoint     Point
	startingDirection Vec
}

func newMap() Map {
	return Map{obstacles: make(map[Point]bool), width: 0, height: 0}
}

func (m *Map) PlaceObstable(p Point) {
	m.obstacles[p] = true
}

func (m *Map) RemoveObstacle(p Point) {
	delete(m.obstacles, p)
}

func (m *Map) HasObstacle(p Point) bool {
	_, ok := m.obstacles[p]

	return ok
}

func (m *Map) InBounds(p Point) bool {
	if p.x < 0 || p.y < 0 {
		return false
	}

	if p.x >= m.width || p.y >= m.height {
		return false
	}

	return true
}

func (m *Map) SetStartingPoint(p Point, v Vec) {
	m.startingPoint = p
	m.startingDirection = v
}

func (m *Map) SetDimensions(width, height int) {
	m.width = width
	m.height = height
}

type Guard struct {
	point     Point
	direction Vec
}

func newGuard(m *Map) Guard {
	return Guard{point: m.startingPoint, direction: m.startingDirection}
}

func main() {
	if cap(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [inputfile]\n", os.Args[0])
		os.Exit(-1)
	}

	m := newMap()

	readInputIntoMap(&m, os.Args[1])

	guard := newGuard(&m)

	part1(&guard, &m)
	part2(&guard, &m)
}

// Part 1 finds all the distinct places on the map the guard visited
func part1(guard *Guard, m *Map) {
	visitedPoints := make(map[Point]bool)

	// Loop guard movement until the guard attempts to leave the map
	for {
		next, _ := nextVisit(guard, m)

		if !m.InBounds(next.point) {
			break
		}

		// Add visit to list and move guard
		visitedPoints[next.point] = true

		guard.point = next.point
		guard.direction = next.direction
	}

	fmt.Println("Part 1:", len(visitedPoints))
}

// Part 2 finds all spots where an additional obstruction could be placed to cause the guard to go in an infinite loop
func part2(guard *Guard, m *Map) {
	numInfiniteLoops := 0

	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			p := Point{x, y}

			// Does an obstruction already exist here?
			if m.HasObstacle(p) {
				continue
			}

			// Obstruction can't be placed at starting block
			if p == m.startingPoint {
				continue
			}

			// Place an obstable in the test spot
			m.PlaceObstable(p)

			visits := make(map[Visit]bool)

			// Reset guard position
			guard.point = m.startingPoint
			guard.direction = m.startingDirection

			for {
				next, ok := nextVisit(guard, m)

				if !ok {
					break
				}

				// If guard already visited this space in this direction, then we've hit an infinite loop
				if _, ok := visits[next]; ok {
					numInfiniteLoops++
					break
				}

				// If the guard left the map, then we're also done
				if !m.InBounds(next.point) {
					break
				}

				// Add visit to list and move guard
				visits[Visit{next.point, next.direction}] = true

				guard.point = next.point
				guard.direction = next.direction
			}

			// Remove the obstacle
			m.RemoveObstacle(p)
		}
	}

	fmt.Println("Part 2:", numInfiniteLoops)
}

// Determine how the guard will move next
func nextVisit(g *Guard, m *Map) (Visit, bool) {
	dir := g.direction

	for {
		point := g.point.Add(dir)

		if !m.HasObstacle(point) {
			return Visit{point, dir}, true
		}

		// Obstable was in the way, rotate 90 degrees
		dir = dir.Rotate()
	}
}

func readInputIntoMap(m *Map, filename string) {
	contentBytes, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	y := 0
	x := 0

	for _, b := range contentBytes {
		switch b {
		case '\n':
			y++
			x = 0

		case '#':
			m.PlaceObstable(Point{x, y})
			x++

		case '^':
			m.SetStartingPoint(Point{x, y}, Vec{0, -1})
			x++

		default:
			x++
		}
	}

	if x != 0 {
		y++
	}

	m.SetDimensions(x, y)
}

func drawMap(m *Map, guard *Guard) string {
	var str string

	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			p := Point{x, y}

			if m.HasObstacle(p) {
				str += "#"
			} else if guard.point == p {
				switch guard.direction {
				case Vec{0, -1}:
					str += "^"

				case Vec{1, 0}:
					str += ">"

				case Vec{0, 1}:
					str += "v"

				case Vec{-1, 0}:
					str += "<"
				}
			} else {
				str += "."
			}

		}

		str += "\n"
	}

	return str
}
