package cli

import (
	"errors"
	"fmt"
	"io"
	"runtime/debug"
	"strings"

	h3 "github.com/dimchansky/h3-go"
)

// environment bundles the streams for one invocation. Runners never touch
// os.Std* directly, so tests can substitute buffers and a fake stdin.
type environment struct {
	in     io.Reader
	out    io.Writer
	errOut io.Writer
}

// trackedWriter latches the first write error and suppresses all further
// writes. This is the CLI's only output-failure handling: there is no
// SIGPIPE machinery; Run inspects the latched error after a command
// succeeds and turns it into exit code 1 (e.g. stdout closed mid-stream by
// a downstream `head`).
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

// The write helpers discard errors on purpose: every writer that reaches a
// runner is a *trackedWriter, which already latches the first failure.
func writeText(w io.Writer, args ...any)             { _, _ = fmt.Fprint(w, args...) }
func writef(w io.Writer, format string, args ...any) { _, _ = fmt.Fprintf(w, format, args...) }
func writeln(w io.Writer, args ...any)               { _, _ = fmt.Fprintln(w, args...) }

// command describes one subcommand: its upstream-matching name, the
// description shown in help listings, the option specs handed to
// parseOptions, and the runner invoked with the parsed arguments. A runner
// returns nil on success, an h3 sentinel error to exit with that H3 code,
// or a *commandError when it has already written its own diagnostic.
type command struct {
	name        string
	description string
	options     []optionSpec
	run         func(environment, parsedArgs) error
}

// buildVersion is stamped by release builds via
// -ldflags "-X github.com/dimchansky/h3-go/internal/cli.buildVersion=<tag>"
// (.github/workflows/release-builds.yml); versionText falls back to module
// build info for go-install'ed binaries.
var buildVersion = "devel"

// Run executes the upstream-compatible h3 command and returns its process
// exit code. It is stateless: arguments and streams are injected, nothing
// global is mutated, so tests invoke it repeatedly in-process.
//
// The exit-code policy mirrors the C CLI (docs/cli-compatibility.md):
// 0 on success, the numeric H3 error code (1–19) on operation failure, 1
// for no/unknown command or an output-write failure — and, deliberately, 0
// for help output and for recognized-command parser errors, because
// upstream's argument parser reports those without failing the process.
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
		// !execute covers both --help and parser errors; each already
		// printed its output and both exit 0 per the upstream contract.
		parsed, execute := parseOptions("h3 "+cmd.name, cmd.description, args[1:], cmd.options, out, errOut)
		if !execute {
			return 0
		}
		err := cmd.run(environment{in: stdin, out: out, errOut: errOut}, parsed)
		if err == nil && out.err != nil {
			return 1 // command succeeded but stdout failed (e.g. broken pipe)
		}
		code := errorCode(err)
		// commandError means the runner already wrote a C-matching
		// diagnostic; anything else gets the generic upstream error line.
		var quiet *commandError
		if code != 0 && !errors.As(err, &quiet) {
			writef(errOut, "Error %d: %s\n", code, errorDescription(code))
		}
		return code
	}
	writeln(out, "Please use h3 --help to see how to use this command.")
	return 1
}

// versionText resolves the build metadata shown by --version: the
// ldflags-injected buildVersion when present, otherwise the module version
// recorded by `go install`, otherwise "devel".
func versionText() string {
	if buildVersion != "devel" {
		return buildVersion
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return buildVersion
}

// commandError carries an explicit exit code for failures whose diagnostic
// the runner has already written (via failDirect). Run recognizes it and
// suppresses the generic "Error N: ..." stderr line that H3 sentinel errors
// would otherwise produce.
type commandError struct{ code int }

func (e *commandError) Error() string { return fmt.Sprintf("command exited with status %d", e.code) }

// failDirect writes an exact upstream diagnostic (wording matters — the
// scenario suite matches it) and returns the quiet exit-1 error.
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

// codeErrors is index-aligned with the H3 C error codes: position N holds
// the sentinel for code N (position 0 is success). The CLI exits with the
// numeric code, so the order must match h3api.h exactly — do not reorder.
var codeErrors = []error{
	nil, h3.ErrFailed, h3.ErrDomain, h3.ErrLatLngDomain, h3.ErrResolutionDomain,
	h3.ErrCellInvalid, h3.ErrDirectedEdgeInvalid, h3.ErrUndirectedEdgeInvalid,
	h3.ErrVertexInvalid, h3.ErrPentagon, h3.ErrDuplicateInput, h3.ErrNotNeighbors,
	h3.ErrResolutionMismatch, h3.ErrMemoryAlloc, h3.ErrMemoryBounds,
	h3.ErrOptionInvalid, h3.ErrIndexInvalid, h3.ErrBaseCellDomain,
	h3.ErrDigitDomain, h3.ErrDeletedDigit,
}

// errorCode maps a runner error to the process exit code: an explicit
// commandError code, the matching H3 error code, or 1 for anything else.
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

// errorDescription returns the upstream describeH3Error text for a code.
// The Go sentinels carry an "h3: " prefix that C output does not have, so
// it is stripped here.
func errorDescription(code int) string {
	if code <= 0 || code >= len(codeErrors) {
		return "Invalid error code"
	}
	return strings.TrimPrefix(codeErrors[code].Error(), "h3: ")
}
