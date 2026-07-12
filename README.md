# h3-go

[![CI](https://github.com/dimchansky/h3-go/actions/workflows/ci.yml/badge.svg)](https://github.com/dimchansky/h3-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/dimchansky/h3-go.svg)](https://pkg.go.dev/github.com/dimchansky/h3-go)
[![Go Version](https://img.shields.io/github/go-mod/go-version/dimchansky/h3-go)](go.mod)
[![Latest Tag](https://img.shields.io/github/v/tag/dimchansky/h3-go?label=latest)](https://github.com/dimchansky/h3-go/tags)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE)

A **pure-Go** implementation of [Uber's H3](https://h3geo.org), the hexagonal
hierarchical geospatial indexing system — behaviorally equivalent to
**H3 C v4.4.0**, with no cgo, no external dependencies, and no `unsafe`
(both guarantees are enforced in CI).

## Why this library

The official Go option, [uber/h3-go](https://github.com/uber/h3-go), is a cgo
binding around the C library. That is a solid choice, but cgo brings real
costs: a C toolchain everywhere you build, harder cross-compilation, cgo call
overhead, and C memory outside the Go runtime's view. This project removes
the C dependency without giving up C-level correctness:

- **Pure, safe Go.** `CGO_ENABLED=0` builds work everywhere Go runs;
  cross-compiling is `GOOS=... go build`. Production code contains no
  `unsafe` — a hard invariant checked by a CI gate (`make check-unsafe`)
  across every build mode.
- **Complete.** All **78/78** public functions of H3 C v4.4.0 are covered —
  indexing, hierarchy, traversal, directed edges, vertexes, regions/polyfill,
  compaction, and measurement — not just the subset wrapped by bindings.
- **Correct by construction.** The implementation is a function-by-function
  port of the C sources, and a 227-file parity suite compiles the *original*
  upstream C and compares Go vs C behavior in-process (see
  [Testing](#testing-and-c-parity)). A separate differential suite
  cross-checks against the official cgo binding.
- **Allocation-aware by design.** Every collection API has a zero-allocation
  `Append*` form and, where it fits, a streaming `iter.Seq` form — measured
  by allocation assertions that run in CI, not just benchmarks.

## Install

```sh
go get github.com/dimchansky/h3-go
```

Requires Go ≥ 1.24 (the two most recent Go releases are tested in CI).
Nothing else: no C toolchain, no environment setup.

Install the upstream-compatible **`h3`** command-line utility separately:

```sh
go install github.com/dimchansky/h3-go/cmd/h3@latest
h3 latLngToCell -r 9 --lat 37.7759 --lng -122.4180
```

It implements all 63 H3 C v4.4.0 commands, including JSON/WKT/newline output
and stdin/file batch workflows. See the [CLI compatibility contract](docs/cli-compatibility.md).

## Usage

```go
import h3 "github.com/dimchansky/h3-go"
```

**Indexing** — coordinates are typed (`Angle`), so degree/radian mix-ups
don't compile:

```go
cell, err := h3.LatLngToCell(h3.LatLngDegs(37.7759, -122.4180), 9)
fmt.Println(cell)            // 8928308280fffff
center, _ := cell.LatLng()
fmt.Println(center.Lat.Deg()) // 37.7767...
```

**Typed indexes** — `Cell`, `DirectedEdge`, and `Vertex` are distinct
`uint64` types with validated parsing and JSON-ready text marshaling:

```go
c, err := h3.ParseCell("8928308280fffff") // mode-validated
_ = c.IsValid()                           // true
data, _ := json.Marshal(c)                // "8928308280fffff"
```

**Hierarchy and traversal**:

```go
parent, _ := cell.Parent(5)
children, _ := parent.Children(9)
disk, _ := cell.GridDisk(2)          // cell + everything within 2 rings
dist, _ := cell.GridDistance(other)
```

**Errors** are sentinels mirroring the C error codes:

```go
if _, err := cell.GridDiskUnsafe(2); errors.Is(err, h3.ErrPentagon) {
    // fall back to the safe variant
}
```

**Allocation control** — reuse buffers or stream, nothing is forced:

```go
buf := make([]h3.Cell, 0, 64)
disk, _ = cell.AppendGridDisk(buf[:0], 2)     // 0 allocs when capacity suffices

for child := range cell.ChildrenSeq(12) {     // 0 allocs, nothing materialized
    _ = child
}
```

**Regions**:

```go
polygon := h3.GeoPolygon{GeoLoop: h3.GeoLoop{
    h3.LatLngDegs(37.813, -122.408),
    h3.LatLngDegs(37.782, -122.386),
    h3.LatLngDegs(37.708, -122.390),
    h3.LatLngDegs(37.708, -122.507),
    h3.LatLngDegs(37.784, -122.511),
}}
cells, _ := h3.PolygonToCells(polygon, 7)
outline, _ := h3.CellsToMultiPolygon(cells)
compact, _ := h3.CompactCells(cells)
```

More runnable examples are in
[example_test.go](example_test.go) and on
[pkg.go.dev](https://pkg.go.dev/github.com/dimchansky/h3-go).

## Testing and C parity

Correctness is enforced in layers, all running in CI:

1. **cgo parity suite** (opt-in, `//go:build cgo && c2go`): 227 test files
   compile the pristine upstream C sources (downloaded to `testref/`, never
   vendored) and compare every ported function against the original C
   implementation in-process — exact values, exact error codes.
2. **Ported upstream tests**: the H3 project's own `testXxx.c` suites,
   translated to Go (tracked in [docs/ported-c-tests.md](docs/ported-c-tests.md)).
3. **Public API tests** with known vectors, pentagon edge cases, allocation
   assertions (`testing.AllocsPerRun`), and three fuzz targets.
4. **Differential tests** against the official uber/h3-go cgo binding
   (`interop/uberdiff`, separate module).
5. **Structural gates**: `make check-unsafe` (no `unsafe` reachable from any
   normal build), `make check-api` (every C public function ported and
   publicly represented), and a golden API-surface lock
   ([docs/api-surface.txt](docs/api-surface.txt)).

The pure-Go test suite (`make test`) needs nothing but Go. Only the parity
suite requires a C toolchain and network access, and it is strictly opt-in.

## Performance

Indicative numbers from `make bench` (Apple M-series, darwin/arm64, Go 1.25;
measure on your own hardware):

| Operation | Time | Allocations |
|---|---|---|
| `LatLngToCell` (res 9) | ~570 ns | 0 |
| `Cell.LatLng` | ~350 ns | 0 |
| `Cell.Boundary` | ~1.0 µs | 0 |
| `AppendGridDisk` (k=2, warm buffer) | ~270 ns | 0 |
| `ParseCell` | ~29 ns | 0 |
| `AppendPolygonToCells` (1253 cells) | ~0.94 ms | 3 (algorithm-internal, as in C) |

There is no cgo call overhead anywhere, and the hot paths are
zero-allocation by design (asserted in tests). No claims are made against
other libraries' throughput; run `make test-uberdiff` and `make bench` to
compare on your workload.

## Status and versioning

- **Maturity**: v0.x. The API is complete and exercised, but pre-1.0 —
  breaking changes remain possible between minor versions and are documented
  in the [CHANGELOG](CHANGELOG.md).
- **Upstream compatibility**: behaviorally equivalent to **H3 C v4.4.0**
  (`VersionMajor/Minor/Patch` report the target release). Upstream releases
  are adopted through a documented, tooling-assisted
  [sync workflow](docs/public-api-architecture.md#10-upstream-synchronization-workflow);
  the 4.3.0 → 4.4.0 sync was performed with it.
- **Go support**: the version in [go.mod](go.mod) up to the latest stable
  release (CI tests both ends).

## Project layout

- **Root package `h3`** — public API files (`indexing.go`, `traversal.go`, …)
  over a one-file-per-C-function ported implementation
  (`<cfile>_<function>.go`, each with a `// Ported from H3 C: <file>::<name>`
  attribution). Public operations carry `H3 C API: <name>` doc lines.
- **Parity harness** — `*_cgo.go` + `h3lib_*_c2go.c` behind
  `//go:build cgo && c2go`; excluded from every normal build.
- **`cmd/h3` + `internal/cli`** — dependency-free, pure-Go implementation of
  the upstream `h3` executable, with an injectable runner and 170 adapted
  upstream scenarios.
- `tools/apiinventory` — generates [docs/c-api-inventory.csv](docs/c-api-inventory.csv)
  (the full C↔Go function mapping) and enforces completeness.

Deeper documentation:

- [docs/public-api-architecture.md](docs/public-api-architecture.md) — design,
  decision records, measurements, upstream-sync workflow.
- [docs/DEVIATIONS.md](docs/DEVIATIONS.md) — intentional differences from C.
- [docs/FUTURE_WORK.md](docs/FUTURE_WORK.md) — backlog: deferred features,
  profiling-gated ideas, rejected designs, pre-v1.0.0 checklist.
- [CONTRIBUTING.md](CONTRIBUTING.md) — development workflow and porting rules.
- [docs/ci-policy.md](docs/ci-policy.md) — what CI runs when, and why the
  expensive suites are tiered.

## Development

```sh
make test              # pure-Go tests (CGO_ENABLED=0)
go test -race ./...    # race detector
make lint              # gofmt + go vet + golangci-lint + smrcptr
make check-unsafe      # gate: no unsafe reachable from any normal build
make bench             # benchmarks with allocation stats
go test -fuzz FuzzParseCell -fuzztime 30s .   # fuzzing

# Optional: C parity validation (needs a C toolchain + network)
make -C testref h3-source   # download upstream H3 sources
make test-c2go              # run the 227-file parity suite vs original C
make check-api              # every C public function ported & represented
make test-uberdiff          # differential vs the official cgo binding
make test-cli-diff          # all CLI scenarios vs the upstream C executable
make check-cli-inventory    # command/flag/test/fixture/source drift gate
```

## License and attribution

Apache-2.0. This project reimplements algorithms from
[Uber's H3](https://github.com/uber/h3); see [NOTICE](NOTICE) for
attribution.
