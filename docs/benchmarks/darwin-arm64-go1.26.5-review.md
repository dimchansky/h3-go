# Darwin/arm64 Go 1.26.5 benchmark refresh review

This is the internal review of the July 2026 local benchmark refresh. It
compares the new Go 1.26.5 artifact, measured at source commit `f207783`,
with the previous Go 1.25.4 artifact from Git commit `1b7583c` (measured at
source commit `a64705f`). The benchmark names, deterministic inputs,
`-count=10 -benchtime=1s -benchmem` flags, pinned uber/h3-go v4.4.1 module,
and pinned benchstat version are the same.

The source commits are not otherwise identical: intervening public API
additions did not change the benchmark bodies or the existing operations
they call, but they can still affect code layout. The measurements therefore
show how the published result sets differ; they do not isolate the Go
toolchain as the only cause. The previous raw files remain available in Git
history and are not retained as a second headline dataset.

## Timing classification

The classification uses the paired benchstat comparison of the two committed
raw files, plus three separate focused checks of hierarchy, collection, and
cheap scalar calls. The focused checks used three samples at 500 ms each and
were review-only; they are not mixed into the published artifact.

| Classification | Review |
|---|---|
| Broadly unchanged within benchmark noise | Most geometry and polygon operations moved by less than about 5%. Several sub-microsecond pure-Go scalar calls were 2–4% slower, while many binding rows were a few percent faster; these small code-layout-sensitive shifts do not support a broad conclusion. |
| Materially improved | Pure-Go children traversal improved by roughly 5–10% across depths and modes, `CompactCells` by about 12%, `UncompactCells` by about 10%, and grid-disk variants by roughly 4–7%. The binding also improved materially on several cheap cgo-boundary calls, so this is not a pure-Go-only movement. |
| Materially regressed | No broad or reproducible pure-Go regression was established. |
| Changed allocation behavior | No benchmark changed `allocs/op`. Reported `B/op` is effectively unchanged; the warm uncompact row moved from 10 to 9 amortized B/op while remaining at zero allocations, which is not an allocation-count change. |
| Inconclusive | Depth-5 children is statistically tied in the official run. Small scalar changes and cross-run peak-RSS differences remain sensitive to runtime layout and high-water-mark noise. |

Across all benchmark variants, the cross-version benchstat time geomean is
about 4% lower in the new run. This aggregate mixes pure, warm, cold, and
uber modes and is not a claim that Go 1.26.5 universally improves
performance.

## Changed within-run conclusions

- Children at depth 3 changed from a small binding advantage to a small
  pure-Go advantage. At depth 5, the old binding advantage became a
  statistical tie.
- `UncompactCells` changed from a statistical tie to an approximately 11%
  pure-Go advantage.
- The pure-Go service-workload advantage remains, but is smaller: the binding
  is about 20% slower in the new run rather than about 31% slower previously.
- The main platform-specific conclusions remain: on this M1 Max,
  `LatLngToCell`, `CellToBoundary`, and polygon fill favor the binding, while
  cheap accessors favor pure Go; compact and multi-polygon conversion still
  favor the binding on both published environments.

## Process-memory comparison

The semantic checksums still match for every pure/binding workload pair.
Retained Go heap and cumulative allocation totals are broadly stable. The
new process run records higher peak RSS for some pure workloads—notably
polyfill and compact—and a higher Go runtime `Sys` baseline for several
short-lived scenarios. Conversely, some binding and pure scenarios changed
only slightly. Because peak RSS is a single process high-water mark that
includes runtime and allocator slack, these cross-run RSS movements are
classified as inconclusive rather than as an implementation memory
regression. The within-run findings remain useful: the binding's scalar
workload performs about six million mallocs, and its compact and retained
polyfill scenarios use less peak RSS than the pure-Go scenarios in this run.

The measurement ran on AC power with Codex UI and normal desktop services
active. No concurrent repository work ran during timing, but the machine was
not an isolated lab host; small changes should be treated accordingly.
