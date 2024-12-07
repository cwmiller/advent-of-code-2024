// Day 7: Bridge Repair

package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Allowable operands
const (
	Addition = iota
	Multiplication
	Concatenate
)

// An operation is an addition, multiplication, or concatenation of a operand against a previous result
type Operation struct {
	operator int
	operand  int64
}

// Perform operation against accumulator
func (o Operation) Result(accumulator int64) int64 {
	switch o.operator {
	case Addition:
		return accumulator + o.operand

	case Multiplication:

		return accumulator * o.operand
	case Concatenate:
		res, _ := strconv.ParseInt(fmt.Sprintf("%d%d", accumulator, o.operand), 10, 64)
		return res

	default:
		return 0
	}
}

// A calibration holds an expected result along with every possible series of operations to achieve that result
// A calibration is considered valid if ANY of the possibilities are valid
type Calibration struct {
	result     int64
	operations [][]Operation
}

// Reduces operations and returns if the result matches the expected result
func (c Calibration) IsValid() bool {
	for _, test := range c.operations {
		// Reduce operations
		var accumulator int64 = 0

		for _, op := range test {
			accumulator = op.Result(accumulator)
		}

		// Does result match?
		if accumulator == c.result {
			return true
		}
	}

	return false
}

// Represents an input consisting of a result and all numbers that make up that result
type Input struct {
	result  int64
	numbers []int64
}

// Returns a collection of every possible sequence of operations that can be made with the series of numbers given
// Concatenate is only available in part 2
func buildPossibilities(collection [][]Operation, remainingNumbers []int64, includeConcatenate bool) [][]Operation {
	additionalOperations := make([][]Operation, 0)
	number := remainingNumbers[0]

	// Collection always starts with a single operation: Adding the first number
	if len(collection) == 0 {
		collection = make([][]Operation, 1)
		collection[0] = []Operation{{Addition, number}}

		remainingNumbers = remainingNumbers[1:]
		number = remainingNumbers[0]
	}

	// Cycles the current collection
	// Append an Addition operation for the next number in sequence to ALL existing operations
	// Also create another series of operations where the next number in sequence is Multiplied instead of Added
	// For Part 2, a Concatenate operation is also created
	for i, existing := range collection {
		collection[i] = append(collection[i], Operation{Addition, number})

		// Copy the existing sequence in order to add a new one where multiplication is performed instead of add
		multiplication := make([]Operation, len(existing))
		copy(multiplication, existing)

		multiplication = append(multiplication, Operation{Multiplication, number})
		additionalOperations = append(additionalOperations, multiplication)

		// Part 2 includes a concatenate operation that must also be created
		if includeConcatenate {
			concatenate := make([]Operation, len(existing))
			copy(concatenate, existing)

			concatenate = append(concatenate, Operation{Concatenate, number})
			additionalOperations = append(additionalOperations, concatenate)
		}
	}

	// Merge the two collections together
	collection = append(collection, additionalOperations...)

	// Pop off the number we just added operations for
	// If we're left with no numbers, then we are done
	// Else continue on
	remainingNumbers = remainingNumbers[1:]

	if len(remainingNumbers) == 0 {
		return collection
	}

	return buildPossibilities(collection, remainingNumbers, includeConcatenate)

}

func newCalibration(input Input, includeConcatenate bool) Calibration {
	collection := buildPossibilities(make([][]Operation, 0), input.numbers, includeConcatenate)

	return Calibration{input.result, collection}
}

func main() {
	if cap(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [inputfile]\n", os.Args[0])
		os.Exit(-1)
	}

	inputs := readInput(os.Args[1])

	part1(inputs)
	part2(inputs)
}

// Read input file into a series of Inputs
func readInput(filename string) []Input {
	inputs := make([]Input, 0)

	contentBytes, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	contents := string(contentBytes)

	// Each line is in the format RESULT: X Y Z..
	rx := regexp.MustCompile(`(\d+): ([\d+| ]+)`)

	matches := rx.FindAllStringSubmatch(contents, -1)

	for _, line := range matches {
		result, _ := strconv.ParseInt(line[1], 10, 64)

		numbersStr := strings.Split(line[2], " ")
		numbers := make([]int64, len(numbersStr))

		for i, num := range numbersStr {
			numbers[i], _ = strconv.ParseInt(num, 10, 64)
		}

		inputs = append(inputs, Input{result, numbers})
	}

	return inputs
}

// Part 1 tests all inputs filling in operator gaps with * or + to valid valid calibrations
// Answer is the sum of all valid calibrations' results
func part1(inputs []Input) {
	var result int64 = 0

	// Create a Calibration consisting of the input's result and ALL possible equations
	for _, input := range inputs {
		calibration := newCalibration(input, false)

		if calibration.IsValid() {

			result += calibration.result
		}
	}

	fmt.Println("Part 1:", result)
}

// Part 2 adds a concatenation operation that concatenates the operand to the previous value
func part2(inputs []Input) {
	var result int64 = 0

	// Create a Calibration consisting of the input's result and ALL possible equations
	for _, input := range inputs {
		calibration := newCalibration(input, true)

		if calibration.IsValid() {
			result += calibration.result
		}
	}

	fmt.Println("Part 2:", result)
}
