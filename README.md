Jane Coding Challenge
=====================

My CLI app for the Jane coding challenge.

# Requirements

All that's required is Go (at least v1.17), and [mage](https://magefile.org/).

Mage is a replacement for
[make](https://www.gnu.org/software/make/manual/make.html), which is fine, but
magefiles let us do so much more. Also Makefiles are a pain to read.

# Testing

Run `mage testBasic` to run just the basic tests.

Right now this CLI app doesn't have any integration or other fancy/funky tests,
so that's the only stage setup for testing right now.

# Building

There are a few ways to build the binaries for this app.

The default output directory is `output`, which is generated as part of the
build process.

### For Your System

Run `mage build`, and the binary can be found at `output/ratings`. The binary
will have been run through `upx`, so it should be pretty tiny.

### Debug Binary

Run `mage buildDebug` to generate a binary with all debug strings and whatnot
kept intact rather than being stripped from the binary.

### Cross-Platform

Using [gox](https://github.com/mitchellh/gox), we can easily compile for
multiple platforms at once. Run `mage buildAll` and multiple binaries will be
generated in the `output` directory.

To control what operating systems and architectures get binaries built for them,
modify the `binaryTypes` variable in `magefile.go`.

# Running 

To run the binary like so: `rankings parse <file>`, where `<file>` is a path to
a text file containing the result data for soccer matches.

# Etc

A spot for other miscellaneous notes.
