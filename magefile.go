//go:build mage

// Mostly concerned with building the 'rankings' binary, including cross-platform compliation.
// Also has some dev-ops tasks related to pushing new releases to GitHub.
package main

/**
 * File: magefile.go
 * Date: 2021-11-19 18:10:04
 * Creator: Sean Patrick Hagen <sean.hagen@gmail.com>
 */

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/go-git/go-git/v5"
	ver "github.com/hashicorp/go-version"
	"github.com/logrusorgru/aurora"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	// output directory, where all build output will go
	outputDir = "output"
	// go compile command
	goCmd = "go"
)

var (
	// Aliases gives us some nicer ways to call common stages
	Aliases = map[string]interface{}{
		"build": Build.Current,
	}

	// file containing current version of the project
	versionFile = "VERSION"

	// name of the binary to generate
	binName = "ratings"

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
	requiredCommands  = []string{"upx", "gox", "ghr", "goimports"}
	commandMinVersion = []*ver.Version{
		ver.Must(ver.NewVersion("3.96")),
		nil,
		ver.Must(ver.NewVersion("0.14.0")),
		nil,
	}
	commandPages = []string{
		"https://github.com/upx/upx/releases/latest",
		"https://github.com/mitchellh/gox#installation",
		"https://github.com/tcnksm/ghr#install",
		"https://pkg.go.dev/golang.org/x/tools/cmd/goimports#pkg-overview",
	}
	commandPaths = []string{}

	// utility, etc -- do not edit these!
	goodMsg  = aurora.Green("GOOD")
	errorMsg = aurora.Red("ERROR")

	// needed when uploading artifacts to a repo that isn't owned directly by the user
	ghrUser = "seanhagen"

	ldFlagsBase = ""
	version     = ""
	repo        = ""
	build       = ""
	branch      = ""

	rootPath = ""
	goFiles  = []string{}
)

func init() {
	setExecutable()
	setRootPath()
}

/*
 * Stages
 * Each exported function is callable by mage
 */

// Util is the namespace for utility stages
type Util mg.Namespace

// Fmt formats the code using gofmt
func (Util) Fmt() error {
	mg.Deps(initVars, checkForCommands)
	fmt.Printf("Formatting all go source files...")
	args := append([]string{"-w"}, goFiles...)
	err := sh.Run("goimports", args...)
	return printResult(err)
}

// Test is the namespace for stages that run tests
type Test mg.Namespace

// Run ...
func (Test) Run() {
	mg.SerialDeps(
		initVars,
		downloadDeps,
		checkForCommands,
		checkCommandVersions,
		makeOutputDir,
		mg.F(runTests, true),
	)
}

// Build is the namespace for stages that output binaries
type Build mg.Namespace

// Build generates a compact binary for your current system
func (Build) Current() {
	mg.SerialDeps(
		initVars,
		downloadDeps,
		checkForCommands,
		checkCommandVersions,
		makeOutputDir,
		// probably some other steps (running tests, etc)
		mg.F(runTests, false),
		mg.F(buildForCurrent, false),
		upxAllBinaries,
	)
}

// Debug generates a debug-ready binary for your current system
func (Build) Debug() {
	mg.SerialDeps(
		initVars,
		downloadDeps,
		checkForCommands,
		checkCommandVersions,
		makeOutputDir,
		// probably some other steps (running tests, etc)
		mg.F(runTests, false),
		mg.F(buildForCurrent, true),
	)
}

// All cross-compiles the application for multiple operating systems and architectures
func (Build) All() {
	mg.SerialDeps(
		initVars,
		downloadDeps,
		checkForCommands,
		checkCommandVersions,
		clean,
		makeOutputDir,
		// probably some other steps (running tests, etc)
		mg.F(runTests, false),
		buildForAllPlatforms,
		upxAllBinaries,
	)
}

// Release is for tasks related to pushing new versions & assets
type Release mg.Namespace

// VersionTag uses the value of the TAG environment variable to tag the current commit
// and push the tag. The tag must be a valid semantic version.
func (Release) VersionTag() error {
	mg.Deps(initVars)
	var err error
	tag := strings.TrimSpace(os.Getenv("TAG"))
	if tag == "" {
		return errors.New("TAG environment variable is required and can't be blank")
	}
	if _, err = ver.NewVersion(tag); err != nil {
		return err
	}

	err = os.WriteFile(versionFile, []byte(tag), 0644)
	if err != nil {
		return err
	}

	if err = sh.RunV("git", "add", versionFile); err != nil {
		return err
	}
	if err = sh.RunV("git", "commit", "-m", fmt.Sprintf("Version bumped to %v", tag)); err != nil {
		return err
	}
	if err = sh.RunV("git", "push", "origin", branch); err != nil {
		return err
	}

	if err = sh.RunV("git", "tag", "-a", "$TAG", "-m", fmt.Sprintf("updating version to %v", tag)); err != nil {
		return err
	}
	if err = sh.RunV("git", "push", "origin", "$TAG"); err != nil {
		return err
	}

	defer func() {
		if err != nil {
			sh.RunV("git", "tag", "--delete", "$TAG")
			sh.RunV("git", "push", "--delete", "origin", "$TAG")
		}
	}()
	return nil
}

// Create uses ghx to upload cross-platform binaries as a
// release on GitHub. Uses the VERSION file to determine
// release name. The release is created as a draft.
func (Release) Create() {
	mg.SerialDeps(
		checkToken,
		Build.All,
		mg.F(uploadRelease, false, false),
	)
}

// Replace uses ghx to replace the release that matches the
// version set in VERSION.
func (Release) Replace() {
	mg.SerialDeps(
		checkToken,
		Build.All,
		mg.F(uploadRelease, true, false),
	)
}

// Delete uses ghx to delete the release matching the version set in VERSION
func (Release) Delete() {
	mg.SerialDeps(
		checkToken,
		Build.All,
		mg.F(uploadRelease, false, true),
	)
}

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
		should := commandMinVersion[i]
		if should != nil {
			fmt.Printf("Checking that command '%v' is version %v or later...", v, should.String())
			var have *ver.Version
			var err error

			switch v {
			case "upx":
				have, err = getUPXVersion()
			case "golangci-lint":
				have, err = getGoLangCIVersion()
			case "ghr":
				have, err = getGhrVersion()
			}

			if err != nil {
				fmt.Printf(" %v\n", errorMsg)
				return err
			}

			if have.LessThan(should) {
				return fmt.Errorf(
					"command '%v' reports version '%v', need at least version '%v'",
					v,
					have.String(),
					should.String(),
				)
			}
			fmt.Printf(" have version %v -- %v\n", have.String(), aurora.Green("GOOD!"))
		} else {
			fmt.Printf("Command '%v' doesn't output version, assuming that it's new enough.\n", v)
		}
	}
	return nil
}

func getUPXVersion() (*ver.Version, error) {
	out, err := runVersionCmd("upx", []string{"-V"})
	if err != nil {
		return nil, err
	}

	//get: upx 3.96-git-db7ba31ca8ce+
	lines := strings.SplitN(out, "\n", 2)
	v := strings.Replace(lines[0], "upx ", "", -1)
	bits := strings.SplitN(v, "-", 2)
	v = bits[0]
	// transform to: 3.96

	return ver.NewVersion(v)
}

func getGoLangCIVersion() (*ver.Version, error) {
	out, err := runVersionCmd("golangci-lint", []string{"--version"})
	if err != nil {
		return nil, err
	}

	// get: golangci-lint has version 1.43.0 built from 861262b7 on 2021-11-03
	lines := strings.SplitN(out, "\n", 2)
	v := lines[0]
	bits := strings.Split(v, " ")
	for i, x := range bits {
		if x == "version" {
			v = bits[i+1]
			break
		}
	}
	// transform to: 1.43.0
	return ver.NewVersion(v)
}

func getGhrVersion() (*ver.Version, error) {
	out, err := runVersionCmd("ghr", []string{"--version"})
	if err != nil {
		return nil, err
	}
	// get: ghr version v0.14.0
	v := strings.Replace(out, "ghr version v", "", 1)
	v = strings.TrimSpace(v)
	// transform to: 0.14.0

	return ver.NewSemver(v)
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
	fmt.Printf("Starting cross-compilation...")
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
	if !run {
		fmt.Printf(" %v\n", errorMsg)
		return fmt.Errorf("command not run")
	}
	return printResult(err)
}

func buildForCurrent(debug bool) error {
	mg.Deps(initVars)
	// TODO: run go build, output to output directory

	env := getEnv()
	args := []string{
		"build",
		"-a", fmt.Sprintf(`-ldflags='%v'`, ldFlagsBase),
	}

	out := binaryOut
	switch debug {
	case true:
		out = fmt.Sprintf("%v_debug", binaryOut)
		args = append(args, "-o", out, "-gcflags='all=-N -l'", "-race")
	case false:
		args = append(args, "-o", binaryOut)
	}

	args = append(args, ".")

	if debug {
		env["CGO_ENABLED"] = "1"
		fmt.Printf("Creating debug binary at '%v'...", out)
	} else {
		fmt.Printf("Creating binary at '%v'...", out)
	}

	err := sh.RunWith(env, goCmd, args...)
	return printResult(err)
}

func upxAllBinaries() error {
	fmt.Printf("Running UPX on all generated binaries\n")
	files, err := ioutil.ReadDir(outputDir)
	if err != nil {
		fmt.Printf(" %v\n", errorMsg)
		return err
	}

	for _, f := range files {
		fn := fmt.Sprintf("%v/%v", outputDir, f.Name())
		fmt.Printf("\tUPX on '%v'...", fn)
		err := sh.Run("upx", "-q", "-9", fn)
		if x := printResult(err); x != nil {
			return x
		}
	}
	return nil
}

func runTests(verbose bool) error {
	if verbose {
		fmt.Printf("Running tests in verbose mode:\n")
		return runAndStreamOutput("go", "test", "-v", "-p", "1", "-timeout", "20m", "./...")
	}
	fmt.Printf("Running tests...")
	err := sh.Run("go", "test", "-p", "1", "-timeout", "20m", "./...")
	return printResult(err)
}

func downloadDeps() error {
	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}
	return nil
}

func clean() error {
	fmt.Printf("Removing output directory...")
	err := sh.Run("rm", "-rf", outputDir)
	return printResult(err)
}

/*
 * Helpers!
 * These are just utility functions.
 */

func printResult(err error) error {
	if err != nil {
		fmt.Printf(" %v\n", errorMsg)
	}
	fmt.Printf(" %v\n", goodMsg)
	return err
}

func getIOBuffers() (*bytes.Buffer, *bytes.Buffer) {
	return bytes.NewBufferString(""), bytes.NewBufferString("")
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
	return map[string]string{
		"CGO_ENABLED": "0", // static binaries for great justice
	}
}

func runAndStreamOutput(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)

	c.Env = os.Environ()
	c.Dir = rootPath

	fmt.Printf("%s\n\n", c.String())

	stdout, _ := c.StdoutPipe()
	errbuf := bytes.Buffer{}
	c.Stderr = &errbuf
	c.Start()

	reader := bufio.NewReader(stdout)
	line, err := reader.ReadString('\n')
	for err == nil {
		fmt.Print(line)
		line, err = reader.ReadString('\n')
	}

	if err := c.Wait(); err != nil {
		fmt.Printf(errbuf.String())
		// fmt.Printf("Error: %s\n", err)
		// os.Exit(1)
		return err
	}
	return nil
}

// Some variables have external dependencies (like git) which may not always be available.
func initVars() {
	mg.Deps(setGoFiles, setVersion, setLdFlags)
}

func setRootPath() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting pwd: %s\n", err)
		os.Exit(1)
	}
	if err := os.Setenv("RANKING_ROOTPATH", pwd); err != nil {
		fmt.Printf("Error setting root path: %s\n", err)
		os.Exit(1)
	}
	rootPath = pwd
}

func setGoFiles() {
	// GOFILES := $(shell find . -name "*.go" -type f ! -path "*/bindata.go")
	files, err := sh.Output("find", ".", "-name", "*.go", "-type", "f")

	if err != nil {
		fmt.Printf("Error getting go files: %s\n", err)
		os.Exit(1)
	}
	for _, f := range strings.Split(string(files), "\n") {
		if strings.HasSuffix(f, ".go") {
			goFiles = append(goFiles, rootPath+strings.TrimLeft(f, "."))
		}
	}
}

func setExecutable() {
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	binaryOut = fmt.Sprintf("%v/%v", outputDir, binName)
}

func setVersion() {
	bits, err := ioutil.ReadFile("./VERSION")
	if err != nil {
		version = "unknown-version"
	}
	version = strings.TrimSpace(string(bits))
}

func setLdFlags() {
	mg.Deps(setVersion)

	repo, ver, branch, build := "unknown", "unknown", "unknown", "unknown"
	base := "-X main.Repo=%v -X main.Version=%v -X main.Branch=%v -X main.Build=%v"

	defer func() {
		ldFlagsBase = fmt.Sprintf(base, repo, ver, branch, build)
	}()

	r, err := git.PlainOpen(".")
	if err != nil {
		return
	}

	repo = getRepo(r)
	branch = getBranch(r)
	build = getBuild(r)
	return
}

func checkToken() error {
	t, ok := os.LookupEnv("GITHUB_TOKEN")
	t = strings.TrimSpace(t)
	if !ok {
		return fmt.Errorf("environment variable GITHUB_TOKEN must be set to deploy a release")
	}
	if t == "" {
		return fmt.Errorf("GITHUB_TOKEN must not be empty")
	}
	return nil
}

func uploadRelease(replace, delete bool) error {
	if replace && delete {
		return fmt.Errorf("can't replace AND delete a release")
	}

	fmt.Printf("Uploading binaries to GitHub...")
	// rm $(BUILDDIR)/$(BINARY)
	_, err := os.Stat(binaryOut)
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf(" %v\n", errorMsg)
		return err
	}

	if os.IsExist(err) {
		err := sh.Run("rm", binaryOut)
		if err != nil {
			fmt.Printf(" %v\n", errorMsg)
			return err
		}
	}

	args := []string{
		"-u",
		ghrUser,
	}

	if replace {
		args = append(args, "-replace")
	}

	if delete {
		args = append(args, "-delete")
	}

	if !replace && !delete {
		args = append(args, "-draft")
	}

	args = append(args, fmt.Sprintf("v%v", version), outputDir)

	// ghr v$(VERSION) $(BUILDDIR)
	err = sh.Run("ghr", args...)
	return printResult(err)
}
