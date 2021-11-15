package main

import (
	"fmt"
)

/**
 * File: main.go
 * Date: 2021-11-15 14:15:29
 * Creator: Sean Patrick Hagen <sean.hagen@gmail.com>
 */

/**
 * Requirements:
 *  - input is text (file name as input argument)
 *  - output is text (output to stdout)
 *  - detect when new match day starts (teams only play once per match day)
 *  - assign points to each team (3 for win, tie/draw is 1, loss is 0)
 *  - if more than 1 team are tied for points, should have same rank and print
 *    in alphabetical order
 */

func main() {
	fmt.Println("hello world!")

	// attempt to open file
	//  - error for file unopenable

	// for each line in file:
	//  1. parse line: get each team and their score
	//  2. has either team played before this day?
	//  2-yes: start a new match day, set as current match day
	//  3. parse match results ( assign points )
	//  4. calculations?
	//  5. ???
	//  6. profit, probably

	// after file is fully parsed, output each match day

}
