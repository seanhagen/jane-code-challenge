package main

/**
 * File: main.go
 * Date: 2021-11-15 14:15:29
 * Creator: Sean Patrick Hagen <sean.hagen@gmail.com>
 */

import "github.com/seanhagen/jane/cmd"

var (
	// Version is set by the build process, contains semantic version
	Version string
	// Build is set by the build process, contains sha tag of build
	Build string
	// Repo is set by the build process, contains the repo where the code for this binary was built from
	Repo string
	// Branch is set by the build process, contains the branch of the repo the binary was built from
	Branch string
)

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
	cmd.Execute()
}
