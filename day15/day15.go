// Day 15: Warehouse Woes
// https://adventofcode.com/2024/day/15

package main

import (
	"fmt"
	"os"
	"regexp"
)

type point struct {
	x, y int
}

func (pt point) Add(v vec) point {
	return point{pt.x + v.x, pt.y + v.y}
}

type vec struct {
	x, y int
}

var Up = vec{0, -1}
var Right = vec{1, 0}
var Down = vec{0, 1}
var Left = vec{-1, 0}

const (
	Wall = iota
	Empty
	Box

	// In Part 2, boxes are 2 spaces wide
	BoxLeft
	BoxRight
)

type warehouse map[point]int

func newWarehouse() warehouse {
	return make(map[point]int)
}

// Populate walls and boxes in a warehouse from puzzle input
// For Part 2, the doubleWide argument is set which causes all walls and boxes to be twice as wide
func (wh warehouse) populate(input string, doubleWide bool) point {
	y := 0
	x := 0

	var robotPos point

	for _, c := range input {
		if c == '\n' {
			// Empty line detected, end of map
			if x == 0 {
				break
			}

			x = 0
			y++

			continue
		}

		var singleWideKind int
		var leftKind int
		var rightKind int

		switch c {
		case '#':
			singleWideKind = Wall
			leftKind = Wall
			rightKind = Wall
		case '@':
			//kind = Robot
			// @ indicates the robot's starting position
			// We'll track his movement outside the map
			robotPos = point{x, y}

			singleWideKind = Empty
			leftKind = Empty
			rightKind = Empty
		case 'O':
			singleWideKind = Box
			leftKind = BoxLeft
			rightKind = BoxRight

		case '.':
			singleWideKind = Empty
			leftKind = Empty
			rightKind = Empty
		}

		if doubleWide {
			wh[point{x, y}] = leftKind
			wh[point{x + 1, y}] = rightKind

			x += 2
		} else {
			wh[point{x, y}] = singleWideKind
			x++
		}
	}

	return robotPos
}

// Determine if the robot can move to the point given
// If the `adjust` parameter is set, boxes will be moved in the warehouse. Else it just determines that they CAN be moved
func (wh warehouse) canMove(pt point, dir vec, adjust bool) bool {
	if target, ok := wh[pt]; ok {
		switch target {
		case Wall:
			return false
		case Empty:
			return true
		case Box:
			// Since there's a box in the way, see if the box can be moved to the next space
			// If so, then this box can move to that spot
			nextPoint := pt.Add(dir)

			if wh.canMove(nextPoint, dir, adjust) {
				if adjust {
					wh[nextPoint] = Box
					wh[pt] = Empty
				}

				return true
			}

		// For double-wide boxes we need to check if both sides of the box can move
		// If moving horizontally, then the same logic works as single-wide boxes. It will move both sides of the boxes.
		// But if moving vertically, then we need to check that both sides of the box can move
		case BoxLeft, BoxRight:
			if dir == Left || dir == Right {
				// Handle horizontal direction, where the single-wide box logic works
				nextPoint := pt.Add(dir)
				if wh.canMove(nextPoint, dir, adjust) {
					if adjust {
						wh[nextPoint] = target
						wh[pt] = Empty
					}

					return true
				}

			} else {
				// Handle vertical direction where both sides of the boxes can move one or two boxes in the way
				var otherSideDir vec
				var otherSideKind int

				if target == BoxLeft {
					otherSideDir = Right
					otherSideKind = BoxRight
				} else {
					otherSideDir = Left
					otherSideKind = BoxLeft
				}

				// Find both points where the box will move if moved in the direction given
				thisSideNextPoint := pt.Add(dir)
				otherSideNextPoint := pt.Add(otherSideDir).Add(dir)

				if wh.canMove(thisSideNextPoint, dir, adjust) && wh.canMove(otherSideNextPoint, dir, adjust) {
					if adjust {
						wh[thisSideNextPoint] = target
						wh[otherSideNextPoint] = otherSideKind
						wh[pt] = Empty
						wh[pt.Add(otherSideDir)] = Empty
					}

					return true
				}
			}
		}
	}

	return false
}

// Tries to make space at the given point for the robot to move into
// Returns if the space is now open for the robot to move to
func (wh warehouse) tryMove(pt point, dir vec) bool {
	// First check if the robot can move to the space
	if wh.canMove(pt, dir, false) {
		// After verifying that it can, commit any box movements to allow the robot to move
		return wh.canMove(pt, dir, true)
	}

	return false
}

// Returns the points of all boxes in the warehouse
func (wh warehouse) allBoxPoints() []point {
	pts := make([]point, 0)

	y := 0
	x := 0

	for {
		x = 0

		if _, ok := wh[point{x, y}]; !ok {
			break
		}

		for {
			pt := point{x, y}
			kind, ok := wh[pt]
			if !ok {
				break
			}

			if kind == Box || kind == BoxLeft {
				pts = append(pts, pt)
			}

			x++
		}

		y++
	}

	return pts
}

// Generate a displayable map of the warehouse
func (wh warehouse) String() string {
	var str string
	y := 0
	x := 0

	for {
		x = 0

		if _, ok := wh[point{x, y}]; !ok {
			break
		}

		for {
			pt := point{x, y}
			kind, ok := wh[pt]
			if !ok {
				break
			}

			switch kind {
			case Empty:
				str += "."
			case Wall:
				str += "#"
			case Box:
				str += "O"
			case BoxLeft:
				str += "["
			case BoxRight:
				str += "]"
			}

			x++
		}

		str += "\n"

		y++
	}

	return str
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

	movements := movementsFromInput(string(inputContents))

	// Part 1
	{
		wh := newWarehouse()
		robotPos := wh.populate(string(inputContents), false)
		solve(wh, robotPos, movements)
	}

	// Part 2
	{
		wh := newWarehouse()
		robotPos := wh.populate(string(inputContents), true)

		solve(wh, robotPos, movements)
	}
}

// Performs all robot movements in a warehouse and returns the GPS Sum of all boxes after they're moved
func solve(wh warehouse, robot point, movements []vec) {
	for _, movement := range movements {
		targetPos := robot.Add(movement)

		if wh.tryMove(targetPos, movement) {
			robot = targetPos
		}
	}

	// Calculate the sum of all box GPS coordinates
	// GPS coordinates of a box are (Y * 100) + X
	gpsSum := 0
	boxPoints := wh.allBoxPoints()

	for _, boxPoint := range boxPoints {
		gpsSum += (100 * boxPoint.y) + boxPoint.x
	}

	fmt.Println("GPS Sum:", gpsSum)
}

// Parse input file contents to retrieve all movement commands for the robot
func movementsFromInput(input string) []vec {
	movements := make([]vec, 0)

	rx := regexp.MustCompile(`[<^>v]+`)
	matches := rx.FindAllString(input, -1)

	for _, match := range matches {
		for _, c := range match {
			var dir vec

			switch c {
			case '^':
				dir = Up
			case '>':
				dir = Right
			case 'v':
				dir = Down
			case '<':
				dir = Left
			}

			movements = append(movements, dir)
		}
	}

	return movements
}
