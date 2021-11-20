//go:build mage

package main

/**
 * File: magefile.go
 * Date: 2021-11-19 18:10:04
 * Creator: Sean Patrick Hagen <sean.hagen@gmail.com>
 */

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/go-version"
	"github.com/logrusorgru/aurora"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// name of the binary to generate
const binName = "ratings"

// output directory, where all build output will go
const outputDir = "output"

var (
	// the location of the binary when building for the local platform in the Build stage
	binaryOut = fmt.Sprintf("%v/%v", outputDir, binName)

	// other OS-es to build for, used in the BuildAll stage
	// see output of `go tool dist list` to see valid values here
	binaryTypes = []struct {
		name  string
		archs []string
	}{
		{"linux", []string{"amd64", "arm", "arm64"}},
		{"darwin", []string{"amd64", "arm64"}},
		{"windows", []string{"amd64", "arm64"}},
	}

	// dev/ops stuff
	// gox for cross compile
	// upx to make binaries smaller
	// golangci for checking that everything is good
	requiredCommands  = []string{"upx", "gox", "golangci-lint"}
	commandMinVersion = []string{"3.96", "-", "1.43.0"}
	commandPages      = []string{
		"https://github.com/upx/upx/releases/latest",
		"https://github.com/mitchellh/gox#installation",
		"https://golangci-lint.run/usage/install/#local-installation",
	}
	commandPaths = []string{}

	// utility, etc -- do not edit these!
	goodMsg  = aurora.Green("GOOD")
	errorMsg = aurora.Red("ERROR")

	ldFlagsBase = getLdFlagBase()
)

/*
 * Stages
 * Each exported function is callable by mage
 */

// Build ...
func Build() {
	mg.SerialDeps(
		checkForCommands,
		checkCommandVersions,
		makeOutputDir,
		// probably some other steps (running tests, etc)
		runTests,
		buildForCurrent,
		upxAllBinaries,
	)
}

// BuildAll ...
func BuildAll() {
	mg.SerialDeps(
		checkForCommands,
		checkCommandVersions,
		makeOutputDir,
		// probably some other steps (running tests, etc)
		// runTests,
		buildForAllPlatforms,
		upxAllBinaries,
	)
}

// RunTests ...
func RunTests() {
	mg.SerialDeps(
		checkForCommands,
		checkCommandVersions,
		makeOutputDir,
		runTests,
	)
}

// // Install ... Runs go mod download and then installs the binary.
// func Install() error {
// 	return nil //sh.Run("go", "install", "./...")
// }

/*
 * Sub-commands!
 * These are parts of build steps, broken out into easier to understand sub-functions
 */

func checkForCommands() error {
	env := getEnv()

	for i, v := range requiredCommands {
		fmt.Printf("Checking for '%v' command...", v)

		ioOut, ioErr := getIOBuffers()
		ran, err := sh.Exec(env, ioOut, ioErr, "which", v)
		if !ran && err == nil {
			fmt.Printf(" %v\n", errorMsg)
			return fmt.Errorf("command didn't run, but no error, not sure what happened there...")
		} else {
			if err != nil {
				fmt.Printf(" %v\n", errorMsg)
				return fmt.Errorf("unable to find command %v, refer to %v on how to install it", v, commandPages[i])
			}
			commandPaths = append(commandPaths,
				strings.TrimSpace(ioOut.String()),
			)
		}
		fmt.Printf(" %v!\n", goodMsg)
	}

	return nil
}

func checkCommandVersions() error {
	for i, v := range requiredCommands {
		ver := commandMinVersion[i]
		if ver != "-" {
			fmt.Printf("Checking that command '%v' is version %v or later...", v, ver)

			should, err := version.NewVersion(ver)
			if err != nil {
				return fmt.Errorf("unable to parse required version '%v', error: %w", ver, err)
			}
			var have *version.Version

			switch v {
			case "upx":
				have, err = getUPXVersion()
			case "golangci-lint":
				have, err = getGoLangCIVersion()
			}

			if have.LessThan(should) {
				return fmt.Errorf(
					"command '%v' reports version '%v', need at least version '%v'",
					v,
					have.String(),
					should.String(),
				)
			} else {
				fmt.Printf(" have version %v -- %v\n", have.String(), aurora.Green("GOOD!"))
			}
		} else {
			fmt.Printf("Command '%v' doesn't output version, assuming that it's new enough.\n", v)
		}
	}
	return nil
}

func getUPXVersion() (*version.Version, error) {
	out, err := runVersionCmd("upx", []string{"-V"})
	if err != nil {
		return nil, err
	}

	//get: upx 3.96-git-db7ba31ca8ce+
	lines := strings.SplitN(out, "\n", 2)
	ver := strings.Replace(lines[0], "upx ", "", -1)
	bits := strings.SplitN(ver, "-", 2)
	ver = bits[0]
	// transform to: 3.96

	return version.NewVersion(ver)
}

func getGoLangCIVersion() (*version.Version, error) {
	out, err := runVersionCmd("golangci-lint", []string{"--version"})
	if err != nil {
		return nil, err
	}

	// get: golangci-lint has version 1.43.0 built from 861262b7 on 2021-11-03
	lines := strings.SplitN(out, "\n", 2)
	ver := lines[0]
	bits := strings.Split(ver, " ")
	for i, v := range bits {
		if v == "version" {
			ver = bits[i+1]
			break
		}
	}
	// transform to: 1.43.0
	return version.NewVersion(ver)
}

// returns stdout if there was no error, stderr if there was an error ( along with the go error )
func runVersionCmd(cmd string, args []string) (string, error) {
	env := getEnv()
	ioOut, ioErr := getIOBuffers()

	ran, err := sh.Exec(env, ioOut, ioErr, cmd, args...)
	if !ran {
		return ioErr.String(), fmt.Errorf("command did not run: %w", err)
	}

	if err != nil {
		return ioErr.String(), fmt.Errorf("unable to run command: %w", err)
	}

	return ioOut.String(), nil
}

func makeOutputDir() error {
	fmt.Printf("Generating output directory...")
	stat, err := os.Stat(outputDir)
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf(" %v\n", errorMsg)
		return fmt.Errorf("unable to check for output directory: %w", err)
	}

	if err != nil && os.IsNotExist(err) {
		if err := sh.Run("mkdir", outputDir); err != nil {
			fmt.Printf(" %v\n", errorMsg)
			return err
		}
		fmt.Printf(" %v\n", goodMsg)
		return nil
	}

	if stat.IsDir() {
		fmt.Printf(" %v\n", goodMsg)
		return nil
	}

	if !stat.IsDir() {
		fmt.Printf(" %v\n", errorMsg)
		return fmt.Errorf("can't create output directory, a file named '%v' already exists", outputDir)
	}

	return fmt.Errorf("should have returned earlier if everything was okay...")
}

func buildForAllPlatforms() error {
	env := getEnv()
	ioOut, ioErr := getIOBuffers()
	osa := ""
	for _, v := range binaryTypes {
		for _, vv := range v.archs {
			if osa == "" {
				osa = fmt.Sprintf("%v/%v", v.name, vv)
			} else {
				osa = fmt.Sprintf("%v %v/%v", osa, v.name, vv)
			}
		}
	}

	out := fmt.Sprintf("%v/{{.OS}}_{{.Arch}}_%v", outputDir, binName)
	args := []string{
		"-ldflags", ldFlagsBase,
		"-osarch", osa,
		"-output", out,
	}

	run, err := sh.Exec(env, ioOut, ioErr, "gox", args...)
	// TODO: handle errors, output a list of the built binaries
	spew.Dump(run, err, ioOut.String(), ioErr.String())
	return nil
}

func buildForCurrent() error {
	// TODO: run go build, output to output directory
	return fmt.Errorf("not yet")
}

func upxAllBinaries() error {
	// TODO: run upx on each binary in output directory
	return fmt.Errorf("can't upx yet")
}

func runGoInstall() error {
	// TODO: install all the things
	return fmt.Errorf("can't install yet")
}

func runTests() error {
	// TODO: run all the tests
	return fmt.Errorf("nope, can't do that yet")
}

func downloadDeps() error {
	// TODO: do this
	// 	if err := sh.Run("go", "mod", "download"); err != nil {
	// 		return err
	// 	}
	return fmt.Errorf("no deps yet")
}

/*
 * Helpers!
 * These are just utility functions.
 */

func getIOBuffers() (*bytes.Buffer, *bytes.Buffer) {
	return bytes.NewBufferString(""), bytes.NewBufferString("")
}

func getLdFlagBase() (ldBase string) {
	repo, ver, branch, build := "unknown", "unknown", "unknown", "unknown"
	base := "-X main.Repo=%v -X main.Version=%v -X main.Branch=%v -X main.Build=%v -s -w"

	defer func() {
		ldBase = fmt.Sprintf(base, repo, ver, branch, build)
	}()

	// wd, err := os.Getwd()
	// if err != nil {
	// 	return
	// }

	r, err := git.PlainOpen(".")
	if err != nil {
		return
	}

	repo = getRepo(r)
	ver = getVersion()
	branch = getBranch(r)
	build = getBuild(r)
	return
}

func getRepo(r *git.Repository) string {
	conf, err := r.Config()
	if err != nil {
		return "unknown-repo"
	}

	origin, ok := conf.Remotes["origin"]
	if !ok {
		fmt.Printf("%v: no remote named %v, can't get repo url!\n", errorMsg, aurora.Bold("origin"))
		return "unknown-repo"
	}

	if len(origin.URLs) <= 0 {
		fmt.Printf("%v: no URLs set in origin!\n", errorMsg)
		return "unknown-repo"
	}

	return origin.URLs[0]
}

func getVersion() string {
	bits, err := ioutil.ReadFile("./VERSION")
	if err != nil {
		return "unknown-version"
	}

	return string(bits)
}

func getBranch(r *git.Repository) string {
	h, err := r.Head()
	if err != nil {
		return "unknown-repo"
	}
	return h.Name().Short()
}

func getBuild(r *git.Repository) string {
	h, err := r.Head()
	if err != nil {
		return "unknown-build"
	}
	return h.Hash().String()

}

func getEnv() map[string]string {
	return map[string]string{}
}
