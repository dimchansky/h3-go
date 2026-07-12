// Package cli implements the h3 command-line utility as a compatible
// pure-Go replacement for the CLI that ships with the upstream H3 C library
// (src/apps/filters/h3.c). The compatibility contract — commands, flags,
// formats, exit codes, and the reviewed inventories that lock them — is
// docs/cli-compatibility.md.
//
// The package is internal on purpose: the supported surface is the behavior
// of the h3 executable (built from cmd/h3), not a Go API. Everything here is
// built on the public h3 package alone — no cgo, no unsafe, no third-party
// dependencies.
//
// # Structure
//
// [Run] is the whole entry point. It receives the argument vector and the
// three standard streams and returns the process exit code; it keeps no
// state between calls, which is what lets the full upstream scenario suite
// execute in-process against buffers instead of spawned binaries.
//
// Commands are declared as [command] descriptors (name, description, option
// specs, run function) grouped by area in the commands_*.go files and
// assembled by commandIndex. Argument parsing (parser.go) reimplements the
// upstream args.c contract. Input decoding and output rendering shared by
// several commands live in io.go; command-specific rendering stays next to
// its run function.
//
// # Compatibility rules that shape the code
//
// Behavior follows the C CLI even where a Go program would naturally differ.
// The load-bearing quirks, each locked by the 170-scenario suite and the
// differential tests against the compiled C binary:
//
//   - Recognized-command parser errors (missing/duplicate/unknown argument)
//     print help to stderr and exit 0, because upstream's argument parser
//     reports failure without setting a process error. Only "no command" or
//     an unrecognized command exits 1.
//   - H3 operation failures exit with the numeric H3 error code (1–19) and
//     print "Error N: <description>" to stderr.
//   - Hex cell arguments are parsed permissively (invalid input becomes
//     index 0 and fails validation later), matching C's sscanf behavior.
//   - Formatting mirrors upstream printf calls, including the %.10f
//     coordinate precision, the %.6f multipolygon precision, and
//     gridDistance printing in hex.
//   - Command registration order is observable in `h3 --help` and therefore
//     fixed.
//
// Help whitespace and diagnostic wording are deliberately not byte-locked;
// see docs/cli-compatibility.md for the full policy, including the additive
// --version command.
//
// # Testing
//
// Four layers, from fast to heavyweight: TestUpstreamCLICompatibility runs
// all upstream scenarios in-process; TestParserAndTopLevelExitContract locks
// the exit-code policy; TestBinaryProcessContract builds the real binary and
// checks pipes, stderr, and exit statuses; TestDifferentialWithCCLI (opt-in
// via H3_CLI_C_BINARY, `make test-cli-diff`) replays every scenario against
// the compiled upstream C executable and compares outputs.
package cli
