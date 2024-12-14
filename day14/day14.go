// Day 14: Restroom Redoubt
// https://adventofcode.com/2024/day/14

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"regexp"
	"strconv"
)

type xy struct {
	x, y int
}

type point xy
type vec xy

func (p point) Add(v vec) point {
	return point{p.x + v.x, p.y + v.y}
}

type robot struct {
	pos point
	vel vec
}

type area struct {
	width, height int
	robots        []robot
}

// Move all the robots in the area one time
// Robots move according to their velocity, and can teleport to the other side of the area if they go out of bounds
func (a *area) moveRobots() {
	for i, robot := range a.robots {
		newPos := robot.pos.Add(robot.vel)

		// If position goes outside area, wrap it around
		if newPos.y < 0 {
			newPos.y += a.height
		}

		if newPos.y >= a.height {
			newPos.y = newPos.y % a.height
		}

		if newPos.x < 0 {
			newPos.x += a.width
		}

		if newPos.x >= a.width {
			newPos.x = newPos.x % a.width
		}

		a.robots[i].pos = newPos
	}
}

// Groups robots by the quardrant they're in
func (a *area) robotsByQuadrant() [4][]robot {
	quads := [4][]robot{
		make([]robot, 0),
		make([]robot, 0),
		make([]robot, 0),
		make([]robot, 0),
	}

	// Robots along the middle axis are not in a quadrant
	midX := a.width / 2
	midY := a.height / 2

	for _, robot := range a.robots {
		x, y := robot.pos.x, robot.pos.y

		// Skip robots in the center lines
		if x == midX || y == midY {
			continue
		}

		quad := 0

		if y > midY {
			quad = 2
		}

		if x > midX {
			quad++
		}

		quads[quad] = append(quads[quad], robot)
	}

	return quads
}

// Return a map of the area, showing the number of robots in each position
func (a *area) String() string {
	var str string

	for y := 0; y < a.height; y++ {
		for x := 0; x < a.width; x++ {
			// Find number of robots at this point
			pos := point{x, y}
			count := 0

			for _, robot := range a.robots {
				if robot.pos == pos {
					count++
				}
			}

			if count > 0 {
				str += strconv.Itoa(count)
			} else {
				str += "."
			}
		}

		str += "\n"
	}

	return str
}

// Create an image of the area
// Unoccupied plots will be a black pixel, while plots with at least one robot will be a white pixel
func (a *area) ToImage() *image.RGBA {
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{a.width, a.height}})

	for y := 0; y < a.height; y++ {
		for x := 0; x < a.width; x++ {
			// Find number of robots at this point
			pos := point{x, y}
			count := 0

			for _, robot := range a.robots {
				if robot.pos == pos {
					count++
					break
				}
			}

			// If there was at least one robot, put a white pixel
			if count > 0 {
				img.Set(x, y, color.White)
			} else {
				img.Set(x, y, color.Black)
			}
		}
	}

	return img
}

// Create a new area with the given dimensions
func newArea(width, height int) *area {
	return &area{
		width,
		height,
		make([]robot, 0),
	}
}

func main() {
	if cap(os.Args) < 6 {
		fmt.Fprintf(os.Stderr, "Usage: %s width height max-seconds [input-file] [output-folder]\n", os.Args[0])
		os.Exit(-1)
	}

	var width, height, maxSeconds int
	var err error

	if width, err = strconv.Atoi(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid width\n")
		os.Exit(-1)
	}

	if height, err = strconv.Atoi(os.Args[2]); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid height\n")
		os.Exit(-1)
	}

	if maxSeconds, err = strconv.Atoi(os.Args[3]); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid max-seconds\n")
		os.Exit(-1)
	}

	area := newArea(width, height)
	loadRobotsFromInput(area, os.Args[4])

	generateFiles(area, maxSeconds, os.Args[5])
}

// Generate an image for every tick up to `ticks`
// Each image will be saved to disk with the tick and the calculated safety factor
// For Part 1, we need the safety factor of the 100th tick
// For Part 2, we have to generate images of X ticks and look for one where a christmas tree appears
func generateFiles(area *area, ticks int, outputFolder string) {
	for i := 0; i < ticks; i++ {
		area.moveRobots()

		// Find the number of robots in each quadrant
		// This is used to calculate the safety factor
		robotsByQuadrant := area.robotsByQuadrant()
		safetyFactor := 1

		for _, robots := range robotsByQuadrant {
			if len(robots) > 0 {
				safetyFactor *= len(robots)
			}
		}

		img := area.ToImage()

		// Output filename is the number of ticks elapsed and the safety factor
		f, err := os.Create(fmt.Sprintf("%s/%d - %d.png", outputFolder, i+1, safetyFactor))

		if err != nil {
			panic(err)
		}

		png.Encode(f, img)
	}
}

// Input file contains each robot patrolling the bathroom
// Each line is a robot and
func loadRobotsFromInput(area *area, filename string) {
	bytes, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	rx := regexp.MustCompile(`p=(\d+),(\d+) v=(\-?\d+),(\-?\d+)`)
	matches := rx.FindAllStringSubmatch(string(bytes), -1)

	for _, match := range matches {
		area.robots = append(area.robots, robot{
			point{mustAtoi(match[1]), mustAtoi(match[2])},
			vec{mustAtoi(match[3]), mustAtoi(match[4])},
		})
	}
}

func mustAtoi(str string) int {
	n, err := strconv.Atoi(str)

	if err != nil {
		panic(err)
	}

	return n
}
