# Comparison with uber/h3-go (the official cgo binding)

[uber/h3-go](https://github.com/uber/h3-go) is the official Go binding for
H3: it vendors the H3 C sources and calls them through cgo. This library is
a pure-Go reimplementation of the same C code. Both are correct, maintained
ways to use H3 from Go; they differ in build model, API shape, allocation
control, and per-operation performance. This document is the evidence-based
comparison behind the README's ["Why this library"](../README.md#why-this-library)
section: exact versions, a function-by-function coverage matrix, behavioral
differences, and trade-offs in both directions.

Companion documents:

- [Migration guide from uber/h3-go](migration-from-uber-h3-go.md) — type
  and call-site mappings with before/after code.
- [Benchmark results and methodology](benchmarks/README.md) — reproducible
  performance and memory comparisons ([interop/uberbench](../interop/uberbench/README.md)
  is the harness).
- [interop/uberdiff](../interop/uberdiff/README.md) — the differential
  correctness suite that cross-checks the two libraries.

## Versions compared

| Component | Version | Notes |
|---|---|---|
| This library | current commit | behavioral target **H3 C v4.4.0** (`VersionMajor/Minor/Patch`) |
| uber/h3-go | **v4.4.1** (pinned in [interop/uberbench](../interop/uberbench/go.mod) and [interop/uberdiff](../interop/uberdiff/go.mod)) | vendors **H3 C v4.4.1** |
| H3 C | v4.4.0 vs v4.4.1 | v4.4.1 changed only the `VERSION` metadata file — **no behavioral difference** ([release notes](https://github.com/uber/h3/releases/tag/v4.4.1)) |

Why uber/h3-go v4.4.1 and not v4.4.0: both binding releases vendor the same
H3 C v4.4.1 sources, and v4.4.1 additionally contains binding-level
performance fixes (stack-allocated coordinate conversion, zero-copy edge
slices, one-call `GridDisksUnsafe`). Pinning v4.4.1 therefore compares
against the binding's **best** release for the H3 4.4 line — the fair
choice for benchmarks — while keeping the underlying C behavior identical
to this library's v4.4.0 parity target.

Version skew to be aware of (checked 2026-07-12):

- The **latest stable uber/h3-go is v4.5.0**, which vendors **H3 C v4.5.0**
  — one minor release ahead of this library. H3 4.5.0 adds
  `reverseDirectedEdge` (exposed as `DirectedEdge.Reverse`), bidirectional
  `gridPathCells`, and stricter error reporting in
  `cellsToLinkedMultiPolygon`; the binding tracked it within days. This is
  the structural trade-off: the official binding adopts new C releases by
  re-vendoring, while this library ports them function by function
  (see the [sync workflow](public-api-architecture.md#10-upstream-synchronization-workflow)).
  Comparisons in this repository are **always at the common H3 4.4 level**;
  numbers and matrices here do not mix H3 versions.
- Uber's repository also contains **`x/h3go`, an experimental pure-Go H3
  implementation** (package doc: "an implementation of the H3 library
  entirely in Go"). As of 2026-07-12 it exists only on the `master` branch
  — it is in no tagged release, is not importable as a versioned module,
  and covers a subset of the API (indexing, hierarchy, sets, boundaries).
  It signals that Uber, too, sees value in a cgo-free H3, but it is not
  yet something to build on or benchmark against; this comparison will be
  extended to it if it ships (tracked in [FUTURE_WORK.md](FUTURE_WORK.md)).

The benchmark environment pins (Go version, C compiler, hardware, exact
commit, module versions, date) are recorded per result set in
[docs/benchmarks/](benchmarks/README.md).

## Coverage summary

Both libraries cover the H3 C v4.4.0 public API almost completely; the
differences are at the edges and in API shape:

- **This library**: all 78 public C functions are covered — 75 by exported
  API, 3 absorbed where Go makes them meaningless (`describeH3Error` →
  sentinel error values, `destroyLinkedMultiPolygon` → garbage collection,
  `maxFaceCount` → internal sizing). Enforced by `make check-api` against
  the upstream headers.
- **uber/h3-go v4.4.1**: 53 functions directly available, 10 available
  under a different Go shape, 13 absorbed (sizing helpers, unit-conversion
  constants, memory destructors), and 2 with no equivalent:
  single-origin `gridDiskUnsafe` and `constructCell`.
- **Beyond the C API**, each library adds its own conveniences:
  - *uber/h3-go only*: `ImmediateParent`/`ImmediateChildren` sugar,
    `IndexFromString`/`IndexToString` on raw `uint64`, a generic `Index`
    constraint, an optional result cap on `PolygonToCellsExperimental`
    (v4.5.0 adds `DirectedEdge.Reverse`).
  - *this library only*: a zero-allocation `Append*` form for every
    collection API, streaming iterators (`Cell.ChildrenSeq`, `CellsAtRes`,
    `PolygonToCellsExperimentalSeq`), exported sizing functions
    (`MaxGridDiskSize`, `MaxGridRingSize`, `MaxPolygonToCellsSize`,
    `MaxPolygonToCellsSizeExperimental`, `UncompactCellsSize`,
    `Cell.NumChildren`, `Cell.GridPathLen`), the typed `Angle` coordinate
    representation, validated parsing, and the pure-Go
    [`h3` CLI](../cmd/h3).

## Function-by-function matrix

The tables below are generated from
[comparison-uber-h3-go.csv](comparison-uber-h3-go.csv) by
[tools/ubercompare](../tools/README.md#ubercompare); `make check-ubercompare`
fails CI if the tables, the CSV, this library's locked API surface, and the
C API inventory drift apart. The binding-side symbols are verified against
the pinned uber/h3-go release by `TestMappingSymbolsExist` in
[interop/uberbench](../interop/uberbench/README.md).

<!-- BEGIN GENERATED: ubercompare (edit docs/comparison-uber-h3-go.csv and run `make gen-ubercompare`) -->

Legend — **uber/h3-go status**: `available` (same operation, directly
comparable shape), `different-shape` (same operation, different Go
types/structure), `absorbed` (not exposed because Go semantics make it
unnecessary — sizing helpers, memory destructors, unit converters),
`missing` (no equivalent). **Migration**: `mechanical` (rename/re-type at
the call site), `adaptation` (surrounding code must change shape),
`n/a` (nothing to migrate). A long dash (—) marks an intentionally
absent API.

### Indexing

| H3 C function | This library | uber/h3-go v4 | Status | Semantics | Allocation | Migration |
|---|---|---|---|---|---|---|
| `latLngToCell` | `LatLngToCell`; `LatLng.Cell` | `LatLngToCell`; `LatLng.Cell` | available | LatLng fields are typed Angle here vs float64 degrees there; construct with LatLngDegs | none vs Go-side struct conversion per call | mechanical |
| `cellToLatLng` | `Cell.LatLng` | `CellToLatLng`; `Cell.LatLng` | available | read coordinates via .Lat.Deg()/.Lng.Deg() here vs raw float64 degrees there | none in either | mechanical |
| `cellToBoundary` | `Cell.Boundary` | `CellToBoundary`; `Cell.Boundary` | different-shape | CellBoundary is a fixed-size value with Len/At/Verts here vs a []LatLng slice there | 0 allocs here vs 1 slice alloc there | adaptation |

### Index inspection

| H3 C function | This library | uber/h3-go v4 | Status | Semantics | Allocation | Migration |
|---|---|---|---|---|---|---|
| `getResolution` | `Cell.Resolution`; `DirectedEdge.Resolution`; `Vertex.Resolution` | `Cell.Resolution`; `DirectedEdge.Resolution`; `Vertex.Resolution` | available | identical (no error; works on any bits) | none | mechanical |
| `getBaseCellNumber` | `Cell.BaseCellNumber` | `BaseCellNumber`; `Cell.BaseCellNumber` | available | identical | none | mechanical |
| `getIndexDigit` | `Cell.IndexDigit` | `Cell.IndexDigit`; `DirectedEdge.IndexDigit`; `Vertex.IndexDigit` | available | binding also exposes it on edges/vertexes | none | mechanical |
| `stringToH3` | `ParseCell`; `ParseDirectedEdge`; `ParseVertex`; `Cell.UnmarshalText` | `CellFromString`; `DirectedEdgeFromString`; `VertexFromString`; `IndexFromString`; `Cell.UnmarshalText` | different-shape | parsing validates and returns errors here; the binding's *FromString swallow syntax errors and return index 0 (only UnmarshalText validates), so correct binding callers add IsValid | none in either | adaptation |
| `h3ToString` | `Cell.String`; `DirectedEdge.String`; `Vertex.String`; `Cell.MarshalText` | `CellToString`; `IndexToString`; `Cell.String`; `Cell.MarshalText` | available | same canonical lowercase-hex text in both | 1 string alloc in either | mechanical |
| `isValidCell` | `Cell.IsValid` | `Cell.IsValid` | available | identical | none | mechanical |
| `isValidIndex` | `IsValidIndex` | `IsValidIndex` | different-shape | takes a raw uint64 here vs a generic constrained by Cell/DirectedEdge/Vertex there | none | mechanical |
| `isResClassIII` | `Cell.IsResClassIII` | `Cell.IsResClassIII` | available | identical | none | mechanical |
| `isPentagon` | `Cell.IsPentagon` | `Cell.IsPentagon` | available | identical | none | mechanical |
| `maxFaceCount` | — (sizes IcosahedronFaces internally) | — (sizes IcosahedronFaces internally) | absorbed | not needed under Go slices | n/a | n/a |
| `getIcosahedronFaces` | `Cell.IcosahedronFaces` | `Cell.IcosahedronFaces` | available | identical ([]int; -1 slots pruned in both) | 1 slice alloc in either | mechanical |
| `constructCell` | `ConstructCell` | — | missing | no binding equivalent (H3 4.4.0 addition) | none here | n/a |

### Grid traversal

| H3 C function | This library | uber/h3-go v4 | Status | Semantics | Allocation | Migration |
|---|---|---|---|---|---|---|
| `maxGridDiskSize` | `MaxGridDiskSize` | — (internal sizing) | absorbed | exposed here for Append* buffer pre-sizing | n/a | n/a |
| `gridDisk` | `Cell.GridDisk`; `Cell.AppendGridDisk` | `GridDisk`; `Cell.GridDisk` | available | unordered and null-pruned in both | convenience allocates in both; Append form is 0-alloc with a warm buffer here | mechanical |
| `gridDiskUnsafe` | `Cell.GridDiskUnsafe`; `Cell.AppendGridDiskUnsafe` | — | missing | single-origin unsafe (ring-ordered fast path; ErrPentagon on distortion) not exposed by the binding — its GridDisksUnsafe covers only the batch form | Append form is 0-alloc here | n/a |
| `gridDiskDistances` | `Cell.GridDiskDistances`; `Cell.AppendGridDiskDistances` | `GridDiskDistances`; `Cell.GridDiskDistances` | different-shape | flat ([]Cell, []int32) here vs distance-indexed [][]Cell rings there; binding keeps H3_NULL slots inside rings on pentagon-affected disks, this library prunes them | 2 allocs vs k+2 allocs; Append form 0-alloc here | adaptation |
| `gridDiskDistancesSafe` | `Cell.GridDiskDistancesSafe` | `GridDiskDistancesSafe`; `Cell.GridDiskDistancesSafe` | different-shape | same shape difference as gridDiskDistances | allocating in both | adaptation |
| `gridDiskDistancesUnsafe` | `Cell.GridDiskDistancesUnsafe` | `GridDiskDistancesUnsafe`; `Cell.GridDiskDistancesUnsafe` | different-shape | same shape difference; ErrPentagon semantics identical | allocating in both | adaptation |
| `gridDisksUnsafe` | `GridDisksUnsafe` | `GridDisksUnsafe` | different-shape | one flat fixed-stride buffer (exact C layout; unpruned) here vs pruned per-origin [][]Cell there | 1 alloc here vs origins+1 allocs there | adaptation |
| `maxGridRingSize` | `MaxGridRingSize` | — (internal sizing) | absorbed | exposed here for Append* buffer pre-sizing | n/a | n/a |
| `gridRing` | `Cell.GridRing`; `Cell.AppendGridRing` | `GridRing`; `Cell.GridRing` | available | unordered and null-pruned in both | Append form 0-alloc here | mechanical |
| `gridRingUnsafe` | `Cell.GridRingUnsafe`; `Cell.AppendGridRingUnsafe` | `GridRingUnsafe`; `Cell.GridRingUnsafe` | available | identical (CCW ring order; ErrPentagon) | Append form 0-alloc here | mechanical |
| `gridDistance` | `Cell.GridDistance` | `GridDistance`; `Cell.GridDistance` | available | identical | none | mechanical |
| `gridPathCellsSize` | `Cell.GridPathLen` | — (internal sizing) | absorbed | exposed here for Append* buffer pre-sizing | n/a | n/a |
| `gridPathCells` | `Cell.GridPath`; `Cell.AppendGridPath` | `GridPath`; `Cell.GridPath` | available | identical (inclusive; same path) | Append form 0-alloc here | mechanical |
| `cellToLocalIj` | `CellToLocalIJ` | `CellToLocalIJ` | available | CoordIJ fields are int32 here (C overflow parity) vs int there | none | mechanical |
| `localIjToCell` | `LocalIJToCell` | `LocalIJToCell` | available | same CoordIJ note | none | mechanical |

### Hierarchy and compaction

| H3 C function | This library | uber/h3-go v4 | Status | Semantics | Allocation | Migration |
|---|---|---|---|---|---|---|
| `cellToParent` | `Cell.Parent` | `Cell.Parent`; `Cell.ImmediateParent` | available | no ImmediateParent sugar here: use c.Parent(c.Resolution()-1) | none | mechanical |
| `cellToChildrenSize` | `Cell.NumChildren` | — (internal sizing) | absorbed | exposed here for pre-sizing and planning | n/a | n/a |
| `cellToChildren` | `Cell.Children`; `Cell.AppendChildren`; `Cell.ChildrenSeq` | `Cell.Children`; `Cell.ImmediateChildren` | available | canonical order in both; iterator form streams without materializing here | Append form 0-alloc and Seq form allocation-free here vs slice per call there | mechanical |
| `cellToCenterChild` | `Cell.CenterChild` | `Cell.CenterChild` | available | identical | none | mechanical |
| `cellToChildPos` | `Cell.ChildPos` | `CellToChildPos`; `Cell.ChildPos` | available | int64 position here vs int there | none | mechanical |
| `childPosToCell` | `Cell.ChildAtPos` | `ChildPosToCell`; `Cell.ChildPosToCell` | available | named ChildAtPos here; position is int64 here vs int there | none | mechanical |
| `compactCells` | `CompactCells`; `AppendCompactCells` | `CompactCells` | available | identical (input must be unique same-res cells) | Append form reuses the result buffer here (algorithm scratch still allocates) | mechanical |
| `uncompactCellsSize` | `UncompactCellsSize` | — (internal sizing) | absorbed | exposed here for Append* buffer pre-sizing | n/a | n/a |
| `uncompactCells` | `UncompactCells`; `AppendUncompactCells` | `UncompactCells` | available | identical | Append form 0-alloc here | mechanical |

### Directed edges

| H3 C function | This library | uber/h3-go v4 | Status | Semantics | Allocation | Migration |
|---|---|---|---|---|---|---|
| `areNeighborCells` | `Cell.IsNeighbor` | `Cell.IsNeighbor` | available | identical | none | mechanical |
| `cellsToDirectedEdge` | `Cell.DirectedEdgeTo` | `Cell.DirectedEdge` | available | named DirectedEdgeTo here | none | mechanical |
| `isValidDirectedEdge` | `DirectedEdge.IsValid` | `DirectedEdge.IsValid` | available | identical | none | mechanical |
| `getDirectedEdgeOrigin` | `DirectedEdge.Origin` | `DirectedEdge.Origin` | available | identical | none | mechanical |
| `getDirectedEdgeDestination` | `DirectedEdge.Destination` | `DirectedEdge.Destination` | available | identical | none | mechanical |
| `directedEdgeToCells` | `DirectedEdge.Cells` | `DirectedEdge.Cells` | different-shape | two return values (origin, destination) here vs a 2-element []Cell there | none here vs 1 slice alloc there | mechanical |
| `originToDirectedEdges` | `Cell.DirectedEdges` | `Cell.DirectedEdges` | available | pentagon's deleted slot pruned in both (5 edges) | 1 slice alloc in either | mechanical |
| `directedEdgeToBoundary` | `DirectedEdge.Boundary` | `DirectedEdge.Boundary` | different-shape | CellBoundary value here vs []LatLng there (same cellToBoundary note) | 0 allocs here vs 1 there | adaptation |

### Vertexes

| H3 C function | This library | uber/h3-go v4 | Status | Semantics | Allocation | Migration |
|---|---|---|---|---|---|---|
| `cellToVertex` | `Cell.Vertex` | `Cell.Vertex`; `CellToVertex` | available | identical | none | mechanical |
| `cellToVertexes` | `Cell.Vertexes` | `Cell.Vertexes`; `CellToVertexes` | available | pentagon slot pruned in both (5 vertexes) | 1 slice alloc in either | mechanical |
| `vertexToLatLng` | `Vertex.LatLng` | `Vertex.LatLng`; `VertexToLatLng` | available | typed Angle coordinates here vs degrees there | none | mechanical |
| `isValidVertex` | `Vertex.IsValid` | `Vertex.IsValid`; `IsValidVertex` | available | identical | none | mechanical |

### Regions (polyfill and multi-polygon)

| H3 C function | This library | uber/h3-go v4 | Status | Semantics | Allocation | Migration |
|---|---|---|---|---|---|---|
| `maxPolygonToCellsSize` | `MaxPolygonToCellsSize` | — (internal sizing) | absorbed | exposed here for Append* buffer pre-sizing | n/a | n/a |
| `polygonToCells` | `PolygonToCells`; `AppendPolygonToCells` | `PolygonToCells`; `GeoPolygon.Cells` | available | center-containment in both; GeoPolygon vertices typed Angle here vs degrees there | Append form reuses result buffer here (3 algorithm-internal allocs remain, as in C); binding also mallocs the polygon on the C heap per call | mechanical |
| `maxPolygonToCellsSizeExperimental` | `MaxPolygonToCellsSizeExperimental` | — (internal sizing) | absorbed | exposed here for Append* buffer pre-sizing | n/a | n/a |
| `polygonToCellsExperimental` | `PolygonToCellsExperimental`; `AppendPolygonToCellsExperimental`; `PolygonToCellsExperimentalSeq` | `PolygonToCellsExperimental` | available | binding adds an optional variadic result cap; the Seq form here streams cells without materializing | Append and Seq forms here vs slice + C-heap polygon there | mechanical |
| `cellsToLinkedMultiPolygon` | `CellsToMultiPolygon` | `CellsToMultiPolygon` | available | same []GeoPolygon shape and loop order; vertices typed Angle here vs degrees there | Go-native structures here vs C linked list walked and freed there | mechanical |
| `destroyLinkedMultiPolygon` | — (garbage collected) | — (garbage collected) | absorbed | meaningless under GC in both | n/a | n/a |

### Measurement

| H3 C function | This library | uber/h3-go v4 | Status | Semantics | Allocation | Migration |
|---|---|---|---|---|---|---|
| `cellAreaRads2` | `Cell.AreaRads2` | `CellAreaRads2` | available | method here vs function there | none | mechanical |
| `cellAreaKm2` | `Cell.AreaKm2` | `CellAreaKm2` | available | method here vs function there | none | mechanical |
| `cellAreaM2` | `Cell.AreaM2` | `CellAreaM2` | available | method here vs function there | none | mechanical |
| `edgeLengthRads` | `DirectedEdge.LengthRads` | `EdgeLengthRads` | available | method here vs function there | none | mechanical |
| `edgeLengthKm` | `DirectedEdge.LengthKm` | `EdgeLengthKm` | available | method here vs function there | none | mechanical |
| `edgeLengthM` | `DirectedEdge.LengthM` | `EdgeLengthM` | available | method here vs function there | none | mechanical |
| `getHexagonAreaAvgKm2` | `HexagonAreaAvgKm2` | `HexagonAreaAvgKm2` | available | identical | none | mechanical |
| `getHexagonAreaAvgM2` | `HexagonAreaAvgM2` | `HexagonAreaAvgM2` | available | identical | none | mechanical |
| `getHexagonEdgeLengthAvgKm` | `HexagonEdgeLengthAvgKm` | `HexagonEdgeLengthAvgKm` | available | identical | none | mechanical |
| `getHexagonEdgeLengthAvgM` | `HexagonEdgeLengthAvgM` | `HexagonEdgeLengthAvgM` | available | identical | none | mechanical |
| `greatCircleDistanceRads` | `GreatCircleDistanceRads` | `GreatCircleDistanceRads` | available | no error in either; LatLng construction differs as above | none | mechanical |
| `greatCircleDistanceKm` | `GreatCircleDistanceKm` | `GreatCircleDistanceKm` | available | no error in either | none | mechanical |
| `greatCircleDistanceM` | `GreatCircleDistanceM` | `GreatCircleDistanceM` | available | no error in either | none | mechanical |

### Constants, conversions, and error description

| H3 C function | This library | uber/h3-go v4 | Status | Semantics | Allocation | Migration |
|---|---|---|---|---|---|---|
| `describeH3Error` | — (sentinel error messages) | — (sentinel error messages) | absorbed | both map H3 error codes to package-level Err* values matched with errors.Is; messages carry the C text; note the binding's ErrRsolutionMismatch spelling vs ErrResolutionMismatch here | n/a | mechanical |
| `degsToRads` | `Deg`; `Angle.Rad` | `DegsToRads` (constant) | different-shape | typed Angle constructors here vs multiply-by-constant there | none | mechanical |
| `radsToDegs` | `Rad`; `Angle.Deg` | `RadsToDegs` (constant) | different-shape | typed Angle accessors here vs multiply-by-constant there | none | mechanical |
| `getNumCells` | `NumCells` | `NumCells` | different-shape | (int64 with error) here vs plain int there (binding panics on out-of-range res) | none | adaptation |
| `res0CellCount` | `NumRes0Cells` | — (len of Res0Cells) | absorbed | a compile-time constant here | n/a | mechanical |
| `getRes0Cells` | `Res0Cells` | `Res0Cells` | available | cannot fail: no error return here vs ([]Cell. error) there | 1 slice alloc in either | mechanical |
| `pentagonCount` | `NumPentagons` | `NumPentagons` | available | constants in both | none | mechanical |
| `getPentagons` | `Pentagons` | `Pentagons` | available | identical | 1 slice alloc in either | mechanical |

<!-- END GENERATED: ubercompare -->

## Behavioral differences worth knowing

Everything below is asserted by the equivalence tests in
[interop/uberbench](../interop/uberbench/README.md) or the differential
suite in [interop/uberdiff](../interop/uberdiff/README.md); this is the
short list a migrating user actually hits. The full migration treatment is
in the [migration guide](migration-from-uber-h3-go.md).

- **Coordinate representation.** The binding's `LatLng` is
  `{Lat, Lng float64}` in **degrees**, converted to radians at every cgo
  boundary crossing. This library's `LatLng` holds typed `Angle` values
  (radians internally); degrees enter and exit through `LatLngDegs`,
  `Deg`, and `.Deg()`. Same math, but a struct literal with bare float64
  degrees does not compile here — which is the point.
- **Parsing.** `CellFromString`/`IndexFromString` in the binding swallow
  syntax errors and return index `0`; only `UnmarshalText` validates.
  `ParseCell` here both reports syntax errors and validates the index
  (`ErrCellInvalid`). Migrating code can usually delete a manual
  `IsValid()` check.
- **Result pruning.** Both libraries prune `H3_NULL` slots from disks,
  rings, polyfills, and pentagon edge/vertex lists. One asymmetry:
  `GridDiskDistances*` in the binding returns distance-indexed
  `[][]Cell` rings that **retain** zero slots on pentagon-affected disks
  (they collect in ring 0); this library returns a flat
  `([]Cell, []int32)` pair with nulls pruned.
- **Numeric widths.** Cell counts and child positions are `int64` here
  (`NumCells`, `Cell.ChildPos`, `Cell.NumChildren`) where the binding uses
  `int`; `CoordIJ` fields are `int32` here (preserving C overflow
  behavior) vs `int` there. The binding's `NumCells` panics on an
  out-of-range resolution; here it returns `ErrResolutionDomain`.
- **Errors.** Both expose sentinel `Err*` values matched with `errors.Is`,
  mapped from the same C error codes (the binding spells one
  `ErrRsolutionMismatch`; here it is `ErrResolutionMismatch`). The
  binding's `UnmarshalText` returns ad-hoc, non-sentinel errors.
- **Index types.** `Cell`/`DirectedEdge`/`Vertex` are `uint64` here and
  `int64` in the binding. Valid H3 indexes never set the high bit, so
  values convert loss-free in both directions.
- **Measurement functions** (`cellArea*`, `edgeLength*`,
  `greatCircleDistance*`) agree to ~1e-9 relative, not bit-exactly: C
  compilers contract multiply-adds into FMAs differently than the Go
  compiler, and the formulas amplify last-ulp differences. Cell indexes,
  by contrast, match exactly on every tested input.

## Trade-offs in both directions

Reasons to use **uber/h3-go**:

- **Fastest adoption of new H3 releases** — re-vendoring the C sources is
  cheaper than porting, so the binding usually tracks upstream minors
  first (it is on H3 4.5.0 today; this library is on 4.4.0).
- **Exact C execution** — it runs the reference implementation itself. If
  your requirement is "the same binary code as the C library", a binding
  is that by definition. (This library's answer is behavioral: a 227-file
  parity suite compares every ported function against the compiled C
  original.)
- **Maturity** — v4 has been stable and widely deployed for years; this
  library is v0.x with a [documented path to v1](FUTURE_WORK.md).
- **Per-operation speed in some workloads** — the C core is heavily
  optimized; for several operations the cgo boundary cost is smaller than
  the C-vs-Go codegen difference. See
  [the benchmark results](benchmarks/README.md) for which operations those
  are on the tested machines.

Reasons to use **this library**:

- **No cgo anywhere**: `CGO_ENABLED=0` builds, trivial cross-compilation
  (the CLI ships for six OS/arch targets from one runner), static binaries
  and scratch containers, no C toolchain in CI, race detector without cgo
  caveats, no cgo call overhead.
- **Allocation control**: `Append*` forms and iterators make hot paths
  zero-allocation (asserted in CI, measured in the benchmarks); the
  binding allocates per call and accepts no caller buffers.
- **All memory is Go memory**: visible to the profiler and the GC, no C
  heap the runtime cannot see (quantified in the
  [memory results](benchmarks/README.md)).
- **Typed API**: `Angle` makes degree/radian confusion uncompilable;
  parsing validates; sizing functions are exposed.
- **Complete at its target version**: 78/78 public functions plus the
  upstream-compatible CLI.

## Keeping this comparison honest

Maintainer checklist — run through it after any of: a new uber/h3-go
release, a new H3 C release, a public API change in this library, or new
benchmark runs:

1. `make check-ubercompare` — fails if the matrix, the generated tables,
   the C API inventory, or this library's API surface drift. After editing
   [comparison-uber-h3-go.csv](comparison-uber-h3-go.csv), regenerate with
   `make gen-ubercompare`.
2. On an uber/h3-go release: bump the pin in **both**
   [interop/uberdiff/go.mod](../interop/uberdiff/go.mod) and
   [interop/uberbench/go.mod](../interop/uberbench/go.mod); check which H3
   C version the new binding vendors (its `H3_VERSION` file); run
   `make test-uberdiff test-uberbench` (equivalence failures mean upstream
   behavior changed — reconcile before publishing numbers); update the
   matrix rows and this document's version tables; refresh the uber-only
   API list above.
3. On an H3 C release adopted by this library: `make api-inventory` first
   (CI's `check-api` gate), then revisit every matrix row whose C function
   changed, and re-run the full benchmark suite.
4. Re-running benchmarks: `make bench-uber` on each reference machine
   (see [benchmarks/README.md](benchmarks/README.md) for machine notes and
   the CI workflow) — never mix machines in one table, and refresh the
   README's curated numbers from the new `benchstat` output.
5. Verify every number quoted in the README appears in (or is directly
   computable from) a committed artifact under `docs/benchmarks/`.
