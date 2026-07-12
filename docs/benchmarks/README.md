# Benchmark results: this library vs uber/h3-go

This directory holds the committed, reproducible artifacts of the
comparative benchmark suite in
[interop/uberbench](../../interop/uberbench/README.md), which measures this
pure-Go library against the official cgo binding
[uber/h3-go](https://github.com/uber/h3-go) on semantically equivalent
operations over identical deterministic inputs. The interpretation lives in
[comparison-uber-h3-go.md](../comparison-uber-h3-go.md); curated headline
numbers appear in the [README](../../README.md#performance) and must always
be traceable to an artifact here.

## What is in each result set

One directory per machine/OS — ratios are **never** comparable across
directories, only within one:

| Directory | Environment |
|---|---|
| [darwin-arm64](darwin-arm64/metadata.txt) | Apple M1 Max laptop (dedicated, mains power, minimal load) |
| `linux-amd64` | GitHub Actions `ubuntu-latest` shared runner (see noise caveats) — populated from the [benchmarks workflow](../../.github/workflows/benchmarks.yml) artifact |

Each contains:

- `metadata.txt` — full environment pin: repository commit, uber/h3-go
  module version, vendored H3 C version, Go and C compiler versions,
  CGO settings, CPU model, core count, benchmark flags, date.
- `bench-raw.txt` — raw `go test -bench` output, all repetitions
  (`-count=10 -benchmem` unless metadata says otherwise).
- `benchstat.txt` — [benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat)
  summary (median ± confidence interval per implementation column, with
  A/B ratios and p-values). **This is the generated table README numbers
  come from.**
- `benchstat.csv` — the same, machine-readable.
- `memory.tsv` — the process-level memory matrix from
  [cmd/memprobe](../../interop/uberbench/cmd/memprobe/main.go).

## Methodology

- **Equivalence before timing.** The suite refuses to benchmark a pairing
  whose results differ: `run.sh` runs the equivalence tests first
  (identical cell sets, coordinates within 1e-9°, measurements within
  1e-8 relative). Where APIs return differently shaped results
  (`GridDiskDistances`, `GridDisksUnsafe`, `CellToBoundary`), both the
  benchmark source and the comparison doc say so explicitly.
- **Identical inputs, deterministic datasets.** Fixed-seed global
  coordinates (poles and antimeridian included), the upstream SF test
  polygon (with-hole variants in the equivalence tests), pentagon sets,
  precomputed neighbor/path pairs. Both implementations cycle the same
  arrays with the same index arithmetic.
- **Statistics.** Each benchmark runs `-count=10` (10 independent
  samples); tables report benchstat medians with confidence intervals,
  and benchstat marks statistically insignificant deltas (p ≥ 0.05) as
  `~`. Outliers are visible in the raw file; nothing is hand-averaged.
- **Usage modes are labeled, not blended.** `impl=pure` vs `impl=uber`
  compares the two convenience (allocating) APIs — the apples-to-apples
  headline. `impl=pure-warm` is this library's buffer-reuse `Append*`
  path, which has **no binding equivalent** and is reported separately;
  `impl=pure-cold` (Append with nil buffer) exists to show it costs the
  same as the convenience form. The binding's most efficient usage *is*
  its convenience API (it accepts no caller buffers); its batch API
  (`GridDisksUnsafe`) is benchmarked where it exists.
- **Single-goroutine measurements** (standard `go test -bench`), default
  `GOMAXPROCS` recorded in metadata.

## Memory: what `B/op` can and cannot see

Go's benchmark allocation metrics observe **only the Go heap**. The
binding does real work there too (output buffers are Go-allocated), so
`B/op`/`allocs/op` comparisons are meaningful for results — but the
binding *additionally* allocates on the **C heap** (polygon inputs,
linked multi-polygon structures), which `B/op` reports as zero, and C
stack/register traffic at every cgo crossing. Therefore:

- `bench-raw.txt` / `benchstat.txt` report `B/op`/`allocs/op` — read
  them as *Go-heap* numbers, a lower bound on the binding's true
  footprint for polygon workloads.
- `memory.tsv` reports whole-process **peak RSS** (`getrusage`) with one
  process per (implementation, workload) pair, plus retained Go heap
  after a final GC, cumulative allocation counts, and GC cycles. Peak
  RSS is implementation-agnostic: it includes Go heap, Go runtime, and
  any C heap. Columns are documented in the
  [memprobe source](../../interop/uberbench/cmd/memprobe/main.go).
- Memory categories separated by the workloads: results the caller keeps
  (`retained-polyfill`, and every `…; last result retained` workload),
  temporary algorithm memory (peak RSS minus retained heap),
  and steady-state churn (`scalar-1m` retains nothing). Caller-provided
  reusable buffers exist only in this library (`Append*`) and are covered
  by the `pure-warm` benchmark rows, not by memprobe.

RSS is an upper-bound style metric (high-water mark, includes allocator
slack and the runtime); treat small differences as noise and large,
reproducible differences as signal. No sub-MB precision claims are made.

## Noise and machine-specificity caveats

- Numbers are valid **only for the recorded machine, OS, compilers, and
  library versions**; ratios move between microarchitectures (and between
  C compilers — the binding's speed depends on how its vendored C was
  compiled).
- The darwin-arm64 set comes from a laptop: macOS performs its own
  scheduling/frequency management (no user-controllable governor);
  the run notes in `metadata.txt` state power and load conditions.
- The linux-amd64 set comes from a **shared** GitHub Actions runner:
  noisy neighbors widen confidence intervals; treat small deltas there
  as `~` even when benchstat resolves them. This is also why benchmarks
  are a manual/scheduled workflow
  ([benchmarks.yml](../../.github/workflows/benchmarks.yml)), never a
  per-push CI gate.
- Sub-microsecond benchmarks are sensitive to code layout and branch
  predictors; prefer the batch rows (`BatchLatLngToCell`,
  `ServiceWorkload`) when estimating end-to-end impact.

## Reading the benchstat tables

`benchstat.txt` pivots implementations into columns (`-col /impl`). For
example, in the `sec/op` table the `pure` and `uber` columns are the two
convenience APIs and the final rows give geomeans. `vs base` percentages
compare against the first column (`pure`); negative means the other
column is faster. The `B/op` and `allocs/op` tables follow (Go heap
only — see above).

## Reproducing and refreshing

```sh
make bench-uber                  # full local run -> docs/benchmarks/<goos>-<goarch>/
gh workflow run benchmarks       # Linux amd64 run in CI; download the artifact
```

Refresh procedure and cross-checks are step 4 of the
[maintainer checklist](../comparison-uber-h3-go.md#keeping-this-comparison-honest).
Commit the regenerated directory wholesale (metadata + raw + summaries
must stay from the same run); then re-verify any README numbers against
the new `benchstat.txt`.

## Selected results

See the per-environment `benchstat.txt` for the full tables; the
[README's Performance section](../../README.md#performance) curates a
representative subset with commentary, including the pairings where the
binding is faster. The one-paragraph darwin-arm64 story (2026-07-12 run):

> The cgo boundary costs the binding roughly 30–45 ns and 1–2 Go
> allocations per crossing, so cheap calls are lopsidedly faster in pure
> Go (`Resolution` ~50×, `Parent` ~10×, `IsValidCell` ~9×) and the mixed
> service workload runs ~24% faster. The binding's C core is honestly
> faster on compute-bound kernels (`compactCells` −59%,
> `cellsToMultiPolygon` −44%, polyfill −23%, `gridPath` −32%,
> `latLngToCell` −13%). The `pure-cold` column is statistically
> indistinguishable from `pure` everywhere (p ≥ 0.05), confirming the
> convenience API and `Append*`-with-nil are the same path; `pure-warm`
> rows allocate zero.

A worked example of why `memory.tsv` exists:
`CellsToMultiPolygon` shows the binding at 4 Go allocs / ~2 KiB per call
vs 1187 allocs / ~57 KiB here — but the binding builds the entire linked
multi-polygon on the **C heap** first, which `B/op` cannot see. The
process-level view (`multipolygon-sf9` row) shows near-identical peak RSS
and retained heap for both. Read Go-heap numbers as *visibility*, not
totality; read RSS for the whole-process truth.
