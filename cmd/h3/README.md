# `h3` — the H3 command-line utility in pure Go

This directory builds **`h3`**, a drop-in, dependency-free replacement for
the command-line utility that ships with the upstream
[H3 C library](https://github.com/uber/h3). All 63 commands of H3 C v4.4.0
are implemented with the same flags and aliases, input sources, output
formats, and exit codes. The authoritative compatibility contract — including
the deliberate compatibility quirks — is
[docs/cli-compatibility.md](../../docs/cli-compatibility.md).

## Install

```sh
go install github.com/dimchansky/h3-go/cmd/h3@latest
```

No C toolchain is needed (`CGO_ENABLED=0` works). Every `v*` tag also builds
reproducible archives — linux, macOS, and Windows on amd64 and arm64, with a
`SHA256SUMS` file — as artifacts of the
[release-builds workflow](../../.github/workflows/release-builds.yml).
From a checkout: `make build-cli`.

## Usage

`h3 --help` lists every command; `h3 <command> --help` shows its flags.

```sh
h3 latLngToCell -r 9 --lat 37.7759 --lng -122.4180
# "8928308280fffff"

h3 cellToBoundary -c 8928308280fffff -f wkt
# POLYGON((-122.4171997184 37.7751977829, ...))

h3 gridDisk -c 8928308280fffff -k 1
# [ "8928308280fffff", "8928308280bffff", ... ]

# Cell-set and polygon commands accept inline values, files, and stdin
# (`-i --`) exactly where the upstream CLI does:
h3 compactCells -i cells.txt -f newline
printf '[[37.775, -122.418], [40.689, -74.044]]' | h3 greatCircleDistanceKm -i --
# 4126.3699216676
```

Output formats follow upstream: `json` (default), plus `wkt`, `newline`, or
`numeric` where the corresponding upstream command supports them. Coordinate
and polygon inputs are JSON arrays in latitude/longitude order.

## Exit codes

Matching the upstream CLI: `0` on success; an H3 operation failure exits with
the numeric H3 error code (1–19) and a classification on stderr; help output
and recognized-command argument errors exit `0` (that is the C parser's
contract); no command or an unknown command exits `1`.

## Version

```sh
h3 --version
# h3 4.4.0 (v0.2.0)
```

The first number is the H3 C compatibility version. The parenthesized value
is this module's version: release builds inject the git tag via
`-ldflags "-X github.com/dimchansky/h3-go/internal/cli.buildVersion=<tag>"`
(see [.github/workflows/release-builds.yml](../../.github/workflows/release-builds.yml));
`go install`ed binaries fall back to Go module build metadata. `--version` is
the one deliberate addition over the upstream command set.

## Where the code lives

[main.go](main.go) is intentionally minimal: it hands `os.Args`, the standard
streams, and the exit code over to [`internal/cli`](../../internal/cli),
which implements parsing, command dispatch, and rendering on top of the
public `h3` package API. That split keeps the whole CLI testable in-process.

Test layers (see [docs/cli-compatibility.md](../../docs/cli-compatibility.md)):

```sh
make test-cli           # all 170 upstream CLI scenarios, in-process
make test-cli-process   # real binary: pipes, stderr, exit statuses
make test-cli-diff      # differential vs the compiled upstream C h3 binary
make check-cli-inventory # drift gate over commands/tests/fixtures/sources
```
