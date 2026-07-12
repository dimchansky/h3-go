// Command h3 is a pure-Go, drop-in replacement for the h3 command-line
// utility that ships with the upstream H3 C library: the same commands,
// flags, output formats, and exit codes (docs/cli-compatibility.md is the
// authoritative contract).
//
// main is deliberately trivial. All parsing, dispatch, and I/O live in
// internal/cli, whose Run takes injected arguments and streams and returns
// the exit code, so every CLI behavior is testable in-process without
// spawning a binary.
package main

import (
	"os"

	"github.com/dimchansky/h3-go/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
