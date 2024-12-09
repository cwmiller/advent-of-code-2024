// Day 8: Resonant Collinearity

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

func (p Point) Sub(v Vec) Point {
	return Point{x: p.x - v.x, y: p.y - v.y}
}

type Vec struct {
	x, y int
}

// Represents a pair of towers with the same frequency
type PointPair struct {
	a Point
	b Point
}

// Create a vector of the difference between the two points in a pair
func (pp PointPair) Diff() Vec {
	return Vec{
		pp.a.x - pp.b.x,
		pp.a.y - pp.b.y,
	}
}

type Map struct {
	width, height int
	// Keep track of every tower based on x,y coordinate
	towers map[Point]Frequency
	// Keep track of every point a frequency is at for quick lookup
	frequencies map[Frequency][]Point
	// All antinodes placed on the map
	antinodes map[Point]bool
}

func newMap() *Map {
	return &Map{
		width:       0,
		height:      0,
		towers:      make(map[Point]Frequency),
		frequencies: make(map[Frequency][]Point),
		antinodes:   make(map[Point]bool),
	}
}

// Add a tower with a frequency to the map
func (m *Map) AddTower(p Point, freq Frequency) {
	// Does tower already exist at point?
	_, ok := m.towers[p]

	if !ok {
		m.towers[p] = freq
		_, ok := m.frequencies[freq]

		if !ok {
			m.frequencies[freq] = make([]Point, 0)
		}

		m.frequencies[freq] = append(m.frequencies[freq], p)
	}
}

// Returns if the given point is within the height and width of the map
func (m *Map) InBounds(p Point) bool {
	if p.x < 0 || p.y < 0 {
		return false
	}

	if p.x >= m.width || p.y >= m.height {
		return false
	}

	return true
}

// Return the frequency (if one exists) at the given point
func (m *Map) GetTowerAtPoint(p Point) (Frequency, bool) {
	freq, ok := m.towers[p]

	return freq, ok
}

// Get a list of all points a tower exists for a frequency
func (m *Map) GetPointsForFrequency(freq Frequency) ([]Point, bool) {
	points, ok := m.frequencies[freq]

	if !ok {
		return nil, false
	}

	ret := make([]Point, len(points))
	copy(ret, points)

	return ret, true
}

// Add an antinode to the map
func (m *Map) AddAntiNode(p Point) bool {
	// An antinode cannot be placed outside the map
	if !m.InBounds(p) {
		return false
	}

	m.antinodes[p] = true

	return true
}

// Get a list of all unique frequencies that are plotted on the map
func (m *Map) Frequencies() []Frequency {
	frequencies := make([]Frequency, len(m.frequencies))

	i := 0
	for key := range m.frequencies {
		frequencies[i] = key
		i++
	}

	return frequencies
}

// Set the width and height of the map
func (m *Map) SetDimensions(width, height int) {
	m.width = width
	m.height = height
}

func (m Map) String() string {
	var str string
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			p := Point{x, y}

			frequency, tok := m.GetTowerAtPoint(p)
			_, aok := m.antinodes[p]

			if !tok && !aok {
				str += "."
			} else {
				if aok {
					str += "#"
				} else {
					str += string(frequency)
				}
			}
		}

		str += "\n"
	}

	return str
}

type Frequency string

func main() {
	if cap(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [inputfile]\n", os.Args[0])
		os.Exit(-1)
	}

	m := newMap()

	readInputIntoMap(os.Args[1], m)

	part1(m)

	fmt.Println("-------------------")

	part2(m)
}

// Part 1 finds pairs of towers with the same frequency and adds an antinode
// on both sides of the towers at the same distance as the towers are from each other
func part1(m *Map) {
	frequencies := m.Frequencies()

	for _, freq := range frequencies {
		pairs := freqPointPairs(m, freq)

		for _, pair := range pairs {
			anA := pair.a.Add(pair.Diff())
			anB := pair.b.Sub(pair.Diff())

			part1PlaceAntiNode(m, anA, freq)
			part1PlaceAntiNode(m, anB, freq)
		}
	}

	fmt.Println(m)

	fmt.Println("Part1:", len(m.antinodes))
}

// Part 1 has a restriction where antinodes cannot be placed on top of towers with the same frequency
func part1PlaceAntiNode(m *Map, p Point, freq Frequency) {
	// An antinode can't be placed if a tower of the same frequency is already there
	towerFreq, ok := m.GetTowerAtPoint(p)

	if !ok || freq != towerFreq {
		m.AddAntiNode(p)
	}
}

// Part 2 doesn't stop with the antinodes being placed just 1 on each side of the pair,
// but instead they keep being placed along the line until the end of the map is reached
// The antennas also become antinodes
func part2(m *Map) {
	frequencies := m.Frequencies()

	for _, freq := range frequencies {
		pairs := freqPointPairs(m, freq)

		for _, pair := range pairs {
			// Place antinodes on both towers
			m.AddAntiNode(pair.a)
			m.AddAntiNode(pair.b)

			// Add antinodes all along the line created by the pair until we go off the map

			for anA := pair.a.Add(pair.Diff()); m.InBounds(anA); anA = anA.Add(pair.Diff()) {
				m.AddAntiNode(anA)
			}

			for anB := pair.b.Sub(pair.Diff()); m.InBounds(anB); anB = anB.Sub(pair.Diff()) {
				m.AddAntiNode(anB)
			}
		}
	}

	fmt.Println(m)

	fmt.Println("Part2:", len(m.antinodes))
}

// Create a list of every possible pairing of towers with the same frequency
func freqPointPairs(m *Map, freq Frequency) []PointPair {
	pairs := make([]PointPair, 0)

	freqPoints, _ := m.GetPointsForFrequency(freq)

	for aIdx := 0; aIdx < len(freqPoints); aIdx++ {
		pointA := freqPoints[aIdx]

		// Create a pair with every other point with the same frequency
		for bIdx := 0; bIdx < len(freqPoints); bIdx++ {
			pointB := freqPoints[bIdx]

			// Skip Point A
			if pointA == pointB {
				continue
			}

			pairs = append(pairs, PointPair{pointA, pointB})
		}
	}

	return pairs
}

func readInputIntoMap(filename string, m *Map) {
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

		case '.':
			x++

		default:
			m.AddTower(Point{x, y}, Frequency(b))
			x++
		}
	}

	if x != 0 {
		y++
	}

	m.SetDimensions(x, y)
}
