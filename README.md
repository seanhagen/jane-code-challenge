Jane Coding Challenge
=====================

My CLI app for the Jane coding challenge.

# Requirements

The minimum requirements are Go (at least v1.17), and
[mage](https://magefile.org/). Mage is a replacement for
[make](https://www.gnu.org/software/make/manual/make.html), which is fine, but
magefiles let us do so much more. Also Makefiles are a pain to read.

There are a few other tools that are required, but running `mage` will check
for the required tools, and give you the URL with instructions on how to
install the tool if it's missing.

# Getting A Binary 

There are two ways you can get a binary:

1. Download the appropriate binary [from the latest
   release](https://github.com/seanhagen/jane-code-challenge/releases/latest).
2. Use `mage` to build a binary on your computer -- see below.

The binaries on the release page are built using
[gox](https://github.com/mitchellh/gox), which makes cross-compilation for Go
a snap.

# Mage Targets

Running `mage -l` will output a list of targets, but the three main targets are
the following:

### test:run

This will run all the "basic" tests -- ie, no integration tests, anything that
requires setup & teardown, or anything that needs to reach out to a service.

### build:current 

This will build a binary for _your_ system. The binary can be found at
`output/rankings` (unless you change any of the settings in `magefile.go` ).

### build:all 

This target uses `gox` to build a binary for multiple platforms at once. The
variable `binaryTypes` in `magefile.go` controls what binaries are built.

## Other Useful Targets

* `build:debug` creates a debug binary that still has all of the symbols and
  whatnot required to get the most out of using gdb/delve/etc.
* `release:versionTag` uses the TAG envionment variable to update the version
  tag (both the contents of the `VERSION` file and the actual git tag)
* `release:create` creates a draft release on the GitHub repository

# Running The Application

You can build a binary and use that, or you can use `go run main.go`. The
following instructions will assume you're using a binary.

The binary has a built in help command, which also will be output if you run the
binary with no arguments. It should look something like this:

```
Rankings is a CLI application for reading in the results of soccer
matches and outputs the top 3 teams for each day based on the match
results.

Usage:
  rankings [command]

Available Commands:
  help        Help about any command
  parse       Read and parse match data to produce rankings

Flags:
  -h, --help   help for rankings

Use "rankings [command] --help" for more information about a command.
```

The `parse` command only takes a single argument: the path to a text file
containing the result data for soccer matches. 

The output will be sent to `stdout`.

# Notes

Some notes on how things could be improved, or potential pitfalls.

### Weird Bugs

#### Setting `ldflags`

Something that apparently doesn't translate 1:1 from a Makefile is setting the
`ldflags` flag when calling `go build`. For an app like this it's not _super_
important, but for larger apps I see it as pretty important.

Using `-ldflags` you can [embed information at build
time](https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications)
into the binary that's created. This can be super handy for backend services,
because you can embed the version & commit sha ( and other handy info ) into the
binary so that it can be used in multiple ways:

* as a header value in HTTP responses
* included as part of error messages sent to services like
  [HoneyBadger](https://www.honeybadger.io/)
* included as part of the tracing messages
* output from [a route](https://playground.bibabots.com/v1/info) to ensure that
  all deployed versions are up-to-date
  
There are probably a few other ways compile-time information can be useful,
these are just a few of the ways I've used such information in the past.

However, something about the way that Go passes arguments to commands doesn't
play well with the format of the `ldflags` string. It's not something I want to
spend the next two days debugging, so I removed it from `magefile.go`.

### Other Build Targets 

There are other potentially useful build targets that could be created, or ways
the current build targets could be improved.

#### release:create 

This could read the git history to create a CHANGELOG file and put it in the
`output` directory. It could also ask for a title or body for the release, so
that the step of having to go to GitHub and edit the draft release isn't needed.

#### test:watch 

There are packages that enable Go to watch files, using one of those a Mage
target could be written that watches all Go source files (including newly
created ones) and then runs tests when there are changes.

#### test:integration or test:all

This application doesn't have any, but an application that needs "bigger" tests
(ie, long running, integration, etc) could use a separate Mage target to run the
full suite of tests.

#### test:ci 

Using tools like `go-junit-report`, `gocover-cobertura`, and `golangci-lint`
Mage could output files for a CI system to show test code coverage and other
useful outputs.

