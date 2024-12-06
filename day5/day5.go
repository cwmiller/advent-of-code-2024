// Day 5: Print Queue

package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Represents a page order rule where X comes before Y
type OrderRule struct {
	x, y int
}

// Collection of rules to apply to a page update
type RuleBook struct {
	rules []OrderRule
}

func newRuleBook() *RuleBook {
	return &RuleBook{}
}

// Add a rule to the rule book collection
func (b *RuleBook) AddRule(x, y int) error {
	// Ensure a rule doesn't already exist for these pages
	_, err := b.GetRule(x, y)

	if err == nil {
		return errors.New("order rule already exists")
	}

	b.rules = append(b.rules, OrderRule{x, y})

	return nil
}

// Retrieve a rule from the book given the two numbers
// Numbers can be in either order
func (b *RuleBook) GetRule(page1, page2 int) (OrderRule, error) {
	for _, rule := range b.rules {
		if (rule.x == page1 && rule.y == page2) || (rule.x == page2 && rule.y == page1) {
			return rule, nil
		}
	}

	return OrderRule{}, errors.New("rule not found")
}

// A page update contains a series of page numbers that are expected to satisfy the rule book
type PageUpdate []int

// Determine if a page update containing a sequence of page numbers passes the ordering rules given
func (u *PageUpdate) PassesRules(ruleBook *RuleBook) bool {
	for i := 0; i < len(*u); i++ {
		for j := i + 1; j < len(*u); j++ {
			x := (*u)[i]
			y := (*u)[j]

			rule, err := ruleBook.GetRule(x, y)

			if err == nil {
				// Does rule match the order the number pair is in?
				if rule.x != x && rule.y != y {
					return false
				}
			}
		}
	}

	return true
}

// Reorder a page update to satisfy the rule book
// Returns new PageUpdate of fixed values
func (u *PageUpdate) Fix(ruleBook *RuleBook) PageUpdate {
	fixed := *u

	for i := 0; i < len(fixed); i++ {
		for j := i + 1; j < len(fixed); j++ {
			x := fixed[i]
			y := fixed[j]

			rule, err := ruleBook.GetRule(x, y)

			if err == nil {
				// Does rule match the order the number pair is in?
				// If not, swap them!
				if rule.x != x && rule.y != y {
					fixed[i] = y
					fixed[j] = x
				}
			}
		}
	}

	return fixed
}

// Collection of page updates
type PageUpdateCollection struct {
	pageUpdates []PageUpdate
}

func newPageUpdateCollection() *PageUpdateCollection {
	return &PageUpdateCollection{
		pageUpdates: make([]PageUpdate, 0),
	}
}

// Add a page update to the collection
func (c *PageUpdateCollection) Add(update PageUpdate) {
	c.pageUpdates = append(c.pageUpdates, update)
}

// Return all page updates that satisfy the rule book
func (c *PageUpdateCollection) FindAllPassingPages(ruleBook *RuleBook) []PageUpdate {
	passingUpdates := make([]PageUpdate, 0)

	for _, update := range c.pageUpdates {
		if update.PassesRules(ruleBook) {
			passingUpdates = append(passingUpdates, update)
		}
	}

	return passingUpdates
}

// Return all page updates that DO NOT satisfy the rule book
func (c *PageUpdateCollection) FindAllFailingPages(ruleBook *RuleBook) []PageUpdate {
	failingUpdates := make([]PageUpdate, 0)

	for _, update := range c.pageUpdates {
		if !update.PassesRules(ruleBook) {
			failingUpdates = append(failingUpdates, update)
		}
	}

	return failingUpdates
}

func main() {
	if cap(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [inputfile]\n", os.Args[0])
		os.Exit(-1)
	}

	input := readInputfile(os.Args[1])
	ruleBook := newRuleBook()
	pageCollection := newPageUpdateCollection()

	populateRuleBook(ruleBook, input)
	populatePageUpdatesCollection(pageCollection, input)

	part1(pageCollection, ruleBook)
	part2(pageCollection, ruleBook)

}

// Part 1 returns the sum of the middle page numbers in all the passing page updates
func part1(pageCollection *PageUpdateCollection, ruleBook *RuleBook) {
	part1Result := 0

	passingUpdates := pageCollection.FindAllPassingPages(ruleBook)

	for _, update := range passingUpdates {
		middleIdx := len(update) / 2

		part1Result += update[middleIdx]
	}

	fmt.Println("Part 1:", part1Result)
}

// Part 2 reorders the failing page updates and returns the sum of the middle page numbers
func part2(pageCollection *PageUpdateCollection, ruleBook *RuleBook) {
	part2Result := 0

	failingUpdates := pageCollection.FindAllFailingPages(ruleBook)

	for _, update := range failingUpdates {
		fixedUpdate := update.Fix(ruleBook)
		middleIdx := len(fixedUpdate) / 2

		part2Result += fixedUpdate[middleIdx]

	}

	fmt.Println("Part 2:", part2Result)
}

// Read contents from input file into a string
func readInputfile(filename string) string {
	contentBytes, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	return string(contentBytes)
}

// Find all order rule within the input file's contents and add them to the rule book
// Rules are in the format X|Y where X must come before Y in the page order
func populateRuleBook(book *RuleBook, input string) {
	rx := regexp.MustCompile(`(\d+)\|(\d+)`)

	matches := rx.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		x, _ := strconv.Atoi(match[1])
		y, _ := strconv.Atoi(match[2])

		book.AddRule(x, y)
	}
}

// Find all page updates within the input file's contents and add them to the collection
// Page updates are a variable length string of numbers concatenated by commas
func populatePageUpdatesCollection(collection *PageUpdateCollection, input string) {
	rx := regexp.MustCompile(`.*,.*`)

	matches := rx.FindAllString(input, -1)

	for _, match := range matches {
		pagesStr := strings.Split(match, ",")
		pages := make([]int, len(pagesStr))

		for i, pageStr := range pagesStr {
			pages[i], _ = strconv.Atoi(pageStr)
		}

		collection.Add(pages)
	}
}
