# h3-go

[![CI](https://github.com/dimchansky/h3-go/actions/workflows/ci.yml/badge.svg)](https://github.com/dimchansky/h3-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/dimchansky/h3-go.svg)](https://pkg.go.dev/github.com/dimchansky/h3-go)
[![Go Version](https://img.shields.io/github/go-mod/go-version/dimchansky/h3-go)](go.mod)
[![Latest Tag](https://img.shields.io/github/v/tag/dimchansky/h3-go?label=latest)](https://github.com/dimchansky/h3-go/tags)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE)

A **pure-Go** implementation of [Uber's H3](https://h3geo.org), the hexagonal
hierarchical geospatial indexing system — behaviorally equivalent to
**H3 C v4.4.0**, with no cgo, no external dependencies, and no `unsafe`
(both guarantees are enforced in CI). It ships two things:

- the **`h3` library** — a typed, allocation-aware Go API covering all 78
  public functions of H3 C v4.4.0;
- the **`h3` command-line utility** — a drop-in, pure-Go replacement for the
  upstream C `h3` executable (all 63 commands).

This is a reimplementation, not a cgo binding: the C library is ported to Go
function by function, and every ported function is traceable to — and
parity-tested against — its C original.

## Why this library

The official Go option, [uber/h3-go](https://github.com/uber/h3-go), is a
cgo binding around the C library — a solid, mature choice. If you are using
it (or choosing between the two), the practical differences are these; each
one is backed by the
[full comparison](docs/comparison-uber-h3-go.md) and by
[measured, equivalence-gated benchmarks](docs/benchmarks/README.md) rather
than assumption:

- **No cgo, anywhere.** `CGO_ENABLED=0` builds work everywhere Go runs:
  cross-compiling is `GOOS=... go build`, static binaries and scratch
  containers just work, CI needs no C toolchain. Production code also
  contains no `unsafe` — a hard invariant checked by a CI gate
  (`make check-unsafe`) across every build mode.
- **Same behavior, verified.** A 227-file parity suite compiles the
  *original* upstream C and compares Go vs C behavior in-process (see
  [Correctness](#correctness-and-testing)); differential and equivalence
  suites cross-check against uber/h3-go itself on identical inputs
  ([interop/uberdiff](interop/uberdiff), [interop/uberbench](interop/uberbench)).
- **More H3 4.4 functionality exposed to Go.** This library covers all
  **78/78** public H3 C v4.4.0 functions. At that common H3 version, the
  binding has no equivalent for `ConstructCell` or the single-origin
  `Cell.GridDiskUnsafe`. Beyond the C API, this library adds streaming
  iterators (`Cell.ChildrenSeq`, `CellsAtRes`,
  `PolygonToCellsExperimentalSeq`), caller-owned `Append*` forms, exported
  sizing helpers, grouped distance rings, and the pure-Go `h3` CLI. See the
  [versioned capability comparison](docs/comparison-uber-h3-go.md#coverage-summary).
- **Allocation control the binding cannot offer.** The main variable-size
  hierarchy, traversal, compaction, and polyfill paths have
  caller-owned-buffer alternatives. Several reach zero allocations once
  warm; algorithms that need scratch space still allocate internally.
  The binding accepts no caller buffers. All memory here is Go memory,
  visible to the profiler and the GC; the binding also allocates on the C
  heap, which Go tooling cannot see (quantified in the
  [memory results](docs/benchmarks/README.md#memory-what-bop-can-and-cannot-see)).
- **A typed API that catches mistakes at compile time.** Coordinates are
  `Angle`-typed (degree/radian mix-ups don't compile), parsing validates
  indexes instead of silently returning 0, error values are `errors.Is`
  sentinels mirroring the C error codes.

And the trade-offs, equally explicit:

- **This is v0.x**; the API is complete and heavily exercised, but pre-1.0
  breaking changes remain possible ([status](#status-and-versioning)). The
  binding's v4 line has years of production maturity.
- **The binding tracks new H3 releases sooner** (it vendors C sources; it is
  on H3 4.5.0 today, this library on 4.4.0 — see the
  [version notes](docs/comparison-uber-h3-go.md#versions-compared)).
- **Pure Go is not automatically faster.** Removing the cgo boundary wins
  many workloads; the binding's optimized C core wins others — the
  [Performance](#performance) section shows both honestly.
- **Migration is a real (if mostly mechanical) change**: different
  coordinate construction, a few reshaped results. The
  [migration and upgrade guide](docs/migration-from-uber-h3-go.md) maps every
  API, highlights the additional capabilities available after switching,
  and keeps its example verified by a test.

If you need exact lockstep with the newest C release or maximum ecosystem
maturity, the official binding remains a good choice. If you want H3
behavior with Go's build story, memory model, and type system, that is what
this library is for.

## Install

```sh
go get github.com/dimchansky/h3-go
```

Requires Go ≥ 1.24 (the two most recent Go releases are tested in CI).
Nothing else: no C toolchain, no environment setup.

## Library quick start

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

## The `h3` command-line utility

A pure-Go, dependency-free replacement for the upstream C `h3` executable —
same commands, same flags, same output formats, same exit codes:

```sh
go install github.com/dimchansky/h3-go/cmd/h3@latest
```

```sh
h3 latLngToCell -r 9 --lat 37.7759 --lng -122.4180
# "8928308280fffff"

h3 cellToBoundary -c 8928308280fffff -f wkt
# POLYGON((-122.4171997184 37.7751977829, -122.4161283578 37.7768804484, ...))

# stdin and file batch workflows, exactly where upstream supports them:
printf '[[37.775, -122.418], [40.689, -74.044]]' | h3 greatCircleDistanceKm -i --
# 4126.3699216676
h3 compactCells -i cells.txt -f newline
```

All 63 H3 C v4.4.0 commands are implemented and locked by the 170 upstream
CLI test scenarios plus differential runs against the compiled C binary.
See [cmd/h3](cmd/h3) for the CLI README and
[docs/cli-compatibility.md](docs/cli-compatibility.md) for the exact
compatibility contract. Every `v*` tag builds reproducible archives for six
OS/architecture targets
([release-builds workflow](.github/workflows/release-builds.yml)).

## Correctness and testing

Correctness is enforced in layers; each answers a different question and none
substitutes for another:

1. **cgo parity suite** (opt-in, `//go:build cgo && c2go`): 227 test files
   compile the pristine upstream C sources (downloaded to `testref/`, never
   vendored) and compare every ported function against the original C
   implementation in-process — exact values, exact error codes.
2. **Ported upstream tests**: the H3 project's own test suites, translated to
   Go and tracked case-by-case in a reviewed inventory
   ([docs/ported-c-tests.md](docs/ported-c-tests.md)).
3. **Public API tests** with known vectors, pentagon edge cases, allocation
   assertions (`testing.AllocsPerRun`), and seven fuzz targets covering
   parsing, round-trips, and all upstream fuzzer input domains.
4. **Large fixture suites**: the three input-driven upstream programs replayed
   over 526,546 golden coordinate/boundary records (needs `testref/`; runs
   nightly and on every release tag in CI).
5. **CLI compatibility tests**: all 170 upstream CLI scenarios in-process,
   process-level pipe/exit-status tests, and opt-in differential execution
   against the compiled upstream `h3` binary.
6. **Differential tests** against the official uber/h3-go cgo binding
   ([interop/uberdiff](interop/uberdiff), separate module), plus the
   benchmark suite's own equivalence gates over every benchmarked operation
   pairing ([interop/uberbench](interop/uberbench)).
7. **Structural gates**: `make check-unsafe` (no `unsafe` reachable from any
   normal build), `make check-api` (every C public function ported and
   publicly represented), a golden API-surface lock
   ([docs/api-surface.txt](docs/api-surface.txt)), and drift gates over the
   test/CLI inventories.

The pure-Go test suite (`make test`) needs nothing but Go. Only the parity
and differential suites require a C toolchain and network access, and they
are strictly opt-in. [docs/ci-policy.md](docs/ci-policy.md) explains which
layer runs when in CI.

## Performance

Performance is operation- and platform-specific. This pure-Go implementation
avoids cgo-call overhead and offers reusable-buffer APIs; the official C core
is faster for some heavier geometry operations. The comparison against
**uber/h3-go v4.4.1** covers 33 scenarios with identical deterministic inputs,
semantic-equivalence gates, 10 repetitions, and `benchstat` summaries.

<!-- BEGIN GENERATED: benchdocs README (run `make gen-benchdocs`) -->

The complete report covers **all 33 scenarios**; this scorecard counts the
statistically faster implementation for each execution-time result.

| Environment | this library faster | uber/h3-go faster | no clear difference |
|---|---:|---:|---:|
| Apple M1 Max · `darwin/arm64` | **12** | **20** | 1 |
| GitHub Actions AMD EPYC 7763 · `linux/amd64` | **26** | **5** | 2 |

**[Explore every scenario, timing, allocation, process-memory, and reusable-buffer result →](docs/benchmarks/results.md)**

Several outcomes reverse between the two environments. That is a measured
platform difference, not evidence that either implementation is universally
faster; absolute timings must not be compared across machines.
<!-- END GENERATED: benchdocs README -->

Caveats that matter: Go's `B/op` cannot see the binding's C-heap allocations,
and Linux/amd64 results come from a shared CI runner with wider noise. The
[full results](docs/benchmarks/results.md) include every scenario, confidence
intervals, allocation comparisons, all process-memory workloads,
reusable-buffer modes, provenance, and interpretation guidance. See the
[methodology and reproduction guide](docs/benchmarks/README.md) to reproduce
the suite.

## Status and versioning

- **Maturity**: v0.x. The API is complete and exercised, but pre-1.0 —
  breaking changes remain possible between minor versions and are documented
  in the [CHANGELOG](CHANGELOG.md). The remaining pre-v1.0.0 checklist lives
  in [docs/FUTURE_WORK.md](docs/FUTURE_WORK.md).
- **Upstream compatibility**: behaviorally equivalent to **H3 C v4.4.0**
  (`VersionMajor/Minor/Patch` report the target release). Upstream releases
  are adopted through a documented, tooling-assisted
  [sync workflow](docs/public-api-architecture.md#10-upstream-synchronization-workflow);
  the 4.3.0 → 4.4.0 sync was performed with it.
- **Go support**: the version in [go.mod](go.mod) up to the latest stable
  release (CI tests both ends).

## Repository map

| Path | What it is |
|---|---|
| `*.go` (root) | The `h3` library — public API files (`cell.go`, `traversal.go`, …) layered over a one-file-per-C-function port (each with a `// Ported from H3 C: <file>::<name>` attribution) |
| `*_cgo.go`, `h3lib_*_c2go.c` | Opt-in C parity harness behind `//go:build cgo && c2go`; excluded from every normal build |
| [`cmd/h3`](cmd/h3) | The `h3` executable — a minimal `main` that delegates to `internal/cli` |
| [`internal/cli`](internal/cli) | CLI implementation: upstream-compatible parser, command registry, output encoders, exit-code mapping; consumes only the public `h3` API |
| [`interop/uberdiff`](interop/uberdiff) | Separate Go module that differentially tests this library against the official uber/h3-go cgo binding (nightly CI) |
| [`interop/uberbench`](interop/uberbench) | Separate Go module benchmarking this library against uber/h3-go — equivalence-gated benchmarks plus process-level memory probes; results in [docs/benchmarks](docs/benchmarks/README.md) |
| [`testref`](testref) | Scaffolding that downloads pristine upstream H3 C sources for the parity suite and gates — never vendored |
| [`tools`](tools) | Maintenance commands: API/test/CLI inventories, upstream symbol diff, docs link check ([tools/README.md](tools/README.md)) |
| [`docs`](docs) | Design records, compatibility contracts, and generated inventories ([docs/README.md](docs/README.md)) |
| [`.github/workflows`](.github/workflows) | Tiered CI: fast checks + C parity on every code change, heavy suites nightly and on tags ([docs/ci-policy.md](docs/ci-policy.md)) |

## Documentation

Pick your path:

- **Library user** → [quick start](#library-quick-start) →
  [pkg.go.dev reference](https://pkg.go.dev/github.com/dimchansky/h3-go) →
  [intentional deviations from C](docs/DEVIATIONS.md).
- **CLI user** → [`h3` utility](#the-h3-command-line-utility) →
  [cmd/h3 README](cmd/h3/README.md) →
  [compatibility contract](docs/cli-compatibility.md).
- **Contributor** → [CONTRIBUTING.md](CONTRIBUTING.md) →
  [architecture & decision records](docs/public-api-architecture.md) →
  [CI policy](docs/ci-policy.md) → [lint policy](docs/lint-policy.md).
  Coding agents and new maintainers: start with the
  [quick reference](AGENTS.md).
- **Upstream-sync maintainer** → [CONTRIBUTING.md](CONTRIBUTING.md) →
  [upstream test equivalence](docs/ported-c-tests.md) →
  [maintenance tools](tools/README.md) →
  [example sync record](docs/sync/4.3.0-to-4.4.0.md).

The full annotated index of every document and generated inventory is
[docs/README.md](docs/README.md).

## Development

```sh
make test              # pure-Go tests (CGO_ENABLED=0)
go test -race ./...    # race detector
make lint              # gofmt + go vet + golangci-lint + smrcptr
make check-unsafe      # gate: no unsafe reachable from any normal build
make check-docs        # Markdown link/anchor checker
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

See [CONTRIBUTING.md](CONTRIBUTING.md) for the ground rules (they are
CI-enforced) and the upstream porting workflow.

## License and attribution

Apache-2.0. This project reimplements algorithms from
[Uber's H3](https://github.com/uber/h3); see [NOTICE](NOTICE) for
attribution.
