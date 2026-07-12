# uberbench — comparative benchmarks vs the official cgo binding

This module benchmarks the pure-Go H3 library in this repository against
the official cgo binding [uber/h3-go](https://github.com/uber/h3-go) on
semantically equivalent operations over identical, deterministic datasets,
and measures process-level memory that Go benchmark metrics cannot see.

The results, with full methodology and environment pins, live in
[docs/benchmarks](../../docs/benchmarks/README.md); the interpretation and
the functionality comparison live in
[docs/comparison-uber-h3-go.md](../../docs/comparison-uber-h3-go.md).

## Why a separate module

The root library has **zero dependencies and zero cgo** — a hard, CI-gated
invariant. Benchmarking against the binding requires importing it (and
therefore cgo, a C toolchain, and first-run network access). Keeping that
in a separate module with a `replace` directive means nothing here is ever
part of a consumer's build. The same reasoning as
[interop/uberdiff](../uberdiff/README.md), which cross-checks correctness
rather than performance.

## What is measured

| Category | Benchmarks |
|---|---|
| Scalar / cgo-boundary-sensitive | `LatLngToCell`, `CellToLatLng`, `CellToBoundary`, `CellToParent`, `GridDistance`, `IsNeighbor`, `CellArea`, `ParseCell`, `CellToString`, `IsValidCell`, `Resolution` |
| Fixed-size geometry | `DirectedEdges`, `DirectedEdgeBoundary`, `Vertexes`, `VertexLatLng` |
| Collections | `Children` (depths 1/3/5), `GridDisk` (k=1/5/20), `GridDiskDistances`, `GridRing`, `GridPath`, `GridDisksUnsafe`, `Compact`, `Uncompact`, `PolygonToCells` (SF polygon at res 7/9/11), `CellsToMultiPolygon` |
| Batch workloads | 10,000-point `LatLngToCell` batch; a service-style enrich pipeline (index point → k=1 disk → coarser parent) |

Each benchmark carries an `/impl=` sub-name so
[benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat) can pivot
implementations into columns:

- `impl=pure` — this library's convenience (allocating) API;
- `impl=pure-cold` — this library's `Append*` form with a nil buffer
  (shows convenience and `Append*` are the same path);
- `impl=pure-warm` — this library's `Append*` form with a reused buffer
  (the zero-allocation path);
- `impl=uber` — the binding. Its APIs allocate per call and accept no
  caller buffers, so this is also its most efficient supported usage;
  where the binding has a batch-amortizing API (`GridDisksUnsafe`) it is
  benchmarked as well.

Comparing `pure` vs `uber` is the apples-to-apples convenience-API
comparison. `pure-warm` has **no equivalent in the binding** — it is
reported for users who control their buffers, not as the headline.

## Equivalence gating

`equivalence_test.go` verifies — on the exact benchmark datasets — that
both libraries return semantically equivalent results for every
benchmarked pairing: identical cell/edge/vertex sets, coordinates within
1e-9°, measurements within 1e-8 relative (see the file's comment for why
bit-exactness is not expected for trigonometry-heavy functions). `run.sh`
runs these tests before benchmarking and aborts on any disagreement.
Where the two APIs return differently *shaped* results (e.g.
`GridDiskDistances`: flat pair here, distance-indexed rings there), the
benchmark source notes it and the equivalence test normalizes the shapes
before comparing.

## Process-level memory: cmd/memprobe

`B/op` / `allocs/op` only observe the **Go heap**. The binding also
allocates on the **C heap** (polygon inputs, linked multi-polygon
structures), which Go accounting cannot see. `cmd/memprobe` therefore runs
one (implementation, workload) pair per process and reports peak RSS
(`getrusage`), retained Go heap after a final GC, cumulative allocation
counters, and a work checksum. The workload matrix covers large polyfills,
deep children expansions, a 2-million-cell uncompact/compact, a
million-call scalar loop, multi-polygon assembly, and a retention
scenario. `TestEquivalenceMemWorkloads` pins that both implementations
compute identical checksums.

## Running

```sh
# from the repository root
make bench-uber        # full suite -> docs/benchmarks/<goos>-<goarch>/
make test-uberbench    # just the equivalence tests

# or directly, with knobs
cd interop/uberbench
COUNT=10 BENCHTIME=1s ./run.sh
```

Requires cgo, a C toolchain, and network access on first run (to fetch
the binding). Results land in `docs/benchmarks/<goos>-<goarch>/` with a
`metadata.txt` pinning commit, module versions, toolchains, and hardware.

In CI, the suite runs in the manual/scheduled
[benchmarks workflow](../../.github/workflows/benchmarks.yml) — never as a
per-push gate: shared runners are too noisy for pass/fail thresholds, and
benchmark regressions are a thing to investigate, not a thing to block
unrelated PRs on.

## What the numbers do and do not show

They quantify, per operation and environment, the cost or benefit of the
cgo boundary, the two libraries' allocation behavior, and buffer-reuse
headroom. They do **not** prove either library "faster" in general:
results differ by operation, workload shape, and machine — the binding
wins some pairings honestly (its C core is heavily optimized) — and
shared-runner numbers carry extra noise. Correctness equivalence is
established here only to the extent the gate tests assert; the exhaustive
correctness story is the root module's C parity suite and
[interop/uberdiff](../uberdiff/README.md).

## Updating for a new uber/h3-go or H3 release

1. Bump `github.com/uber/h3-go/v4` in `go.mod` (`go get -u` + `go mod tidy`).
2. Check which H3 C version the new binding vendors (`H3_VERSION` file in
   the module) against this library's target (`h3.VersionMajor/Minor/Patch`);
   if they diverge, say so in `docs/benchmarks/README.md` and
   `docs/comparison-uber-h3-go.md` — do not present cross-H3-version
   numbers as an implementation comparison.
3. `make test-uberbench` — equivalence failures here mean upstream
   behavior changed; reconcile before benchmarking.
4. Re-run `make bench-uber` on the reference machines and commit the new
   artifacts; refresh any numbers quoted in the README and
   `docs/comparison-uber-h3-go.md` (see the
   [maintainer checklist](../../docs/comparison-uber-h3-go.md#keeping-this-comparison-honest)).
