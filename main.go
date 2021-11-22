package main

/**
 * File: main.go
 * Date: 2021-11-15 14:15:29
 * Creator: Sean Patrick Hagen <sean.hagen@gmail.com>
 */

var (
	// Version is set by the build process, contains semantic version
	Version string
	// Build is set by the build process, contains sha tag of build
	Build string
	// Repo is set by the build process, contains the repo where the code for this binary was built from
	Repo string
	// Branch is set by the build process, contains the branch of the repo the binary was built from
	Branch string
import (
	"github.com/seanhagen/jane-coding-challenge/cmd"
)

func main() {
	cmd.Execute()
}
