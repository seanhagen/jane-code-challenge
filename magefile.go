// // +build mage

package main

/**
 * File: magefile.go
 * Date: 2021-11-19 18:10:04
 * Creator: Sean Patrick Hagen <sean.hagen@gmail.com>
 */

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const OutputDir = "output"

var (
	binaryOut = fmt.Sprintf("%v/ratings")
	// oses

	requiredCommands  = []string{"upx", "gox", "golangci-lint"}
	commandMinVersion = []string{"3.96"}
	commandPages      = []string{
		"https://github.com/upx/upx/releases/latest",
		"https://github.com/mitchellh/gox#installation",
		"https://golangci-lint.run/usage/install/#local-installation",
	}

	commandPaths = []string{}
)

// gox for cross compile
// upx to make binaries smaller
// golangci for checking that everything is good

func checkForCommands() error {
	for i, v := range requiredCommands {
		fmt.Printf("Checking for '%v' command...", v)

		out, err := sh.Output("which", v)
		if err != nil {
			fmt.Printf(" %v\n", aurora.Red("ERROR"))
			return fmt.Errorf("unable to find command %v, refer to %v on how to install it", v, commandPages[i])
		}

		commandPaths = append(commandPaths, out)

		fmt.Printf(" %v\n", aurora.Green("GOOD!"))
	}

	return nil
}

func makeOutputDir() error {
	stat, err := os.Stat(OutputDir)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("unable to check for output directory: %w", err)
	}

	if !stat.IsDir() {
		return fmt.Errorf("can't create output directory, a file named '%v' already exists", OutputDir)
	}

	if os.IsNotExist(err) {
		if err := sh.Run("mkdir", OutputDir); err != nil {
			return err
		}
	}

	return nil
}

func buildPlatforms() error {
	mg.SerialDeps(checkForCommands, makeOutputDir)
	return nil
}

// Build ...
func Build() {
	mg.SerialDeps(buildPlatforms)
}

// Install ... Runs go mod download and then installs the binary.
func Install() error {
	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}

	return nil //sh.Run("go", "install", "./...")
}
