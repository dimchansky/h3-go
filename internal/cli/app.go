// Package cli implements the upstream-compatible h3 command-line utility.
package cli

import (
	"errors"
	"fmt"
	"io"
	"runtime/debug"
	"strings"

	h3 "github.com/dimchansky/h3-go"
)

type environment struct {
	in     io.Reader
	out    io.Writer
	errOut io.Writer
}

type trackedWriter struct {
	w   io.Writer
	err error
}

func (w *trackedWriter) Write(p []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}
	n, err := w.w.Write(p)
	w.err = err
	return n, err
}

func writeText(w io.Writer, args ...any)             { _, _ = fmt.Fprint(w, args...) }
func writef(w io.Writer, format string, args ...any) { _, _ = fmt.Fprintf(w, format, args...) }
func writeln(w io.Writer, args ...any)               { _, _ = fmt.Fprintln(w, args...) }

type command struct {
	name        string
	description string
	options     []optionSpec
	run         func(environment, parsedArgs) error
}

var buildVersion = "devel"

// Run executes the upstream-compatible h3 command and returns its process exit code.
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	out := &trackedWriter{w: stdout}
	errOut := &trackedWriter{w: stderr}
	if len(args) == 0 {
		writeln(out, "Please use h3 --help to see how to use this command.")
		return 1
	}
	commands := commandIndex()
	if strings.EqualFold(args[0], "--help") || strings.EqualFold(args[0], "-h") {
		printGeneralHelp(out, commands)
		if out.err != nil {
			return 1
		}
		return 0
	}
	if strings.EqualFold(args[0], "--version") {
		writef(out, "h3 4.4.0 (%s)\n", versionText())
		if out.err != nil {
			return 1
		}
		return 0
	}
	for _, cmd := range commands {
		if !strings.EqualFold(args[0], cmd.name) {
			continue
		}
		parsed, execute := parseOptions("h3 "+cmd.name, cmd.description, args[1:], cmd.options, out, errOut)
		if !execute {
			return 0
		}
		err := cmd.run(environment{in: stdin, out: out, errOut: errOut}, parsed)
		if err == nil && out.err != nil {
			return 1
		}
		code := errorCode(err)
		var quiet *commandError
		if code != 0 && !errors.As(err, &quiet) {
			writef(errOut, "Error %d: %s\n", code, errorDescription(code))
		}
		return code
	}
	writeln(out, "Please use h3 --help to see how to use this command.")
	return 1
}

func versionText() string {
	if buildVersion != "devel" {
		return buildVersion
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return buildVersion
}

type commandError struct{ code int }

func (e *commandError) Error() string { return fmt.Sprintf("command exited with status %d", e.code) }

func failDirect(w io.Writer, message string) error {
	writeText(w, message)
	return &commandError{code: 1}
}

func printGeneralHelp(w io.Writer, commands []command) {
	writeln(w, "h3: Please use one of the subcommands listed to perform an H3 calculation. Use h3 <SUBCOMMAND> --help for details on the usage of any subcommand.")
	writeln(w, "H3 4.4.0")
	writeln(w)
	writeln(w, "\t-h, --help\tShow this help message.")
	for _, cmd := range commands {
		writef(w, "\t%s\t%s\n", cmd.name, cmd.description)
	}
}

var codeErrors = []error{
	nil, h3.ErrFailed, h3.ErrDomain, h3.ErrLatLngDomain, h3.ErrResolutionDomain,
	h3.ErrCellInvalid, h3.ErrDirectedEdgeInvalid, h3.ErrUndirectedEdgeInvalid,
	h3.ErrVertexInvalid, h3.ErrPentagon, h3.ErrDuplicateInput, h3.ErrNotNeighbors,
	h3.ErrResolutionMismatch, h3.ErrMemoryAlloc, h3.ErrMemoryBounds,
	h3.ErrOptionInvalid, h3.ErrIndexInvalid, h3.ErrBaseCellDomain,
	h3.ErrDigitDomain, h3.ErrDeletedDigit,
}

func errorCode(err error) int {
	if err == nil {
		return 0
	}
	var commandErr *commandError
	if errors.As(err, &commandErr) {
		return commandErr.code
	}
	for code := 1; code < len(codeErrors); code++ {
		if errors.Is(err, codeErrors[code]) {
			return code
		}
	}
	return 1
}

func errorDescription(code int) string {
	if code <= 0 || code >= len(codeErrors) {
		return "Invalid error code"
	}
	return strings.TrimPrefix(codeErrors[code].Error(), "h3: ")
}
