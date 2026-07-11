# H3-Go: Public API Architecture

Status: **proposal — awaiting approval before implementation**
Scope: repository `github.com/dimchansky/h3-go` at commit `52d76be`, upstream reference H3 C **v4.3.0** (`testref/h3-4.3.0`).

This document proposes how to layer an idiomatic, strongly typed, allocation-aware public Go
API on top of the existing mechanically ported C implementation, without losing
function-level traceability to upstream H3 C. Every load-bearing claim below was verified
against the actual repository; the experiments backing the design are reproduced in
[Appendix A](#appendix-a-experimental-evidence).

---

## 0. Executive summary

**The port is functionally complete at the C level and functionally unusable at the Go
level.** All **75 of 75** public functions of H3 C 4.3.0 are ported and parity-tested
(see `docs/c-api-inventory.csv`, regenerable with `go run ./tools/apiinventory`), but every
one of them is an *unexported* Go function — the package currently exports **zero H3
operations**. What *is* exported (30 types, ~121 constants, 16 lookup tables) is mostly
accidental: C identifiers that happened to start with an uppercase letter.

The recommended architecture:

1. **Single root package `h3`, two file layers, no `internal/` split.** The mechanically
   ported C-shaped code stays where it is, unexported, one file per C function. The public
   API is added as new plainly named files (`cell.go`, `edge.go`, `region.go`, …) in the
   same package.
2. **`type Cell uint64` becomes the defined public type, and the internal `H3Index`
   becomes a type alias for it** (`type h3Index = Cell`). This single-line change makes
   `[]h3Index` and `[]Cell` the *same type*, so every hot slice-producing algorithm
   (gridDisk, children, compact, polygonToCells, …) is callable from the public API with
   **zero copies, zero unsafe, zero generics** — verified by compiling and running both
   full test suites with the alias in place. `DirectedEdge` and `Vertex` are separate
   defined `uint64` types; all edge/vertex collections in H3 are fixed-size (≤ 6), so the
   scalar conversions at the boundary cost nothing.
3. **Dual-form collection APIs**: a convenient allocating form
   (`Cell.GridDisk(k) ([]Cell, error)`) plus a zero-allocation `Append*` form
   (`Cell.AppendGridDisk(dst []Cell, k int)`), and Go 1.23 `iter.Seq` iterators over the
   already-ported C iterator structs (measured **0 allocs** for iterating 343 children).
4. **One mechanical unexport sweep** turns the accidentally exported C-style identifiers
   unexported (case change only, names otherwise preserved), validated by the 224-file cgo
   parity suite.
5. **No compatibility layer.** There are no tags, no releases, and no callable exported
   API today; nothing usable exists to preserve.
6. **Upstream sync stays function-level**: filename convention + `// Ported from H3 C:`
   attributions + the committed `tools/apiinventory` tool + a CI completeness gate.
7. **Hard invariant (DR-007): the production library is safe Go only.** No file selected
   by any normal build (`go build ./...`, `go test ./...`, `CGO_ENABLED=0 go test ./...`,
   `go test -race ./...`, any GOOS/GOARCH) imports `unsafe`; no public API requires or
   exposes an unsafe abstraction. `unsafe` exists only inside the opt-in cgo parity
   harness behind `cgo && c2go` tags, and a two-layer CI gate (§9.2) keeps it that way.
   Full audit in [Appendix B](#appendix-b-unsafe-audit).

One urgent unrelated finding: **CI is red** (last three runs failed): the `make fmt` gate
fails on 8 unformatted files at HEAD, and behind it `go test -race ./...` enables cgo,
which pulls in the `//go:build cgo` interop files, which `#include` C sources that are not
in the repository. Formatting plus fixing the build tags to `cgo && c2go` is Phase 0.

---

## 1. Current-state assessment

### 1.1 What the repository actually is

The layout differs substantially from what `README.md`, `AGENTS.md`, and `CLAUDE.md`
describe (those documents predate commit `efbef78 "move evrything to root package"` and
describe an `internal/`-package + oracle-CLI design that no longer exists; `AGENTS.md` and
`CLAUDE.md` are `.gitignore`d local files). The real state:

- **One flat root package `h3`** containing:
  - **274 pure-Go non-test files** — one C function per file, named
    `<cfile>_<function>.go` / `<cfile>__<function>.go`, each carrying a
    `// Ported from H3 C: <file>::<name>` attribution;
  - **19 `*_cgo.go` interop files** + **17 `h3lib_*_c2go.c` shims** that compile the
    original C sources for in-process parity testing;
  - **290 test files**: 224 `*_parity_test.go` (Go-vs-C comparisons), 31 ported upstream
    unit-test files (`testGridDisk_test.go`, …, tracked in `docs/ported-c-tests.md`), 35 other
    plain Go tests.
- `testref/` holds the downloaded upstream source (`make -C testref` / `make ref`); only
  the small oracle scaffolding is committed.
- `go.mod`: `go 1.22`, zero dependencies. No git tags; 332 commits; single author.
- `h3.go.backup` / `h3_test.go.backup`: an earlier ~96-function public API draft (commits
  `bdfb7f9`…`3013484`) that used an `unsafe`-based generic `castSlice[From, To ~uint64]`
  to bridge `[]Cell` ↔ `[]H3Index`. It was renamed to `.backup` in commit `91d349a` and no
  longer compiles (its helpers were deleted). This proposal supersedes it — and makes its
  `unsafe` bridge unnecessary.

### 1.2 Fidelity to C

Fidelity is excellent and is the repository's main asset:

- Signatures mirror C exactly: out-pointers and caller-sized buffers, `H3Error` returns.
  E.g. `gridDiskDistances(origin H3Index, k int32, out []H3Index, distances []int32) H3Error`
  ([algos__gridDiskDistances.go](../algos__gridDiskDistances.go)),
  `cellToVertexes(cell H3Index, vertexes *[6]H3Index) H3Error`,
  `polygonToCells(geoPolygon *GeoPolygon, res int32, flags uint32, out []H3Index) H3Error`.
- C `int` → `int32` everywhere (494 uses) to preserve 32-bit overflow semantics; C
  `int64_t` → `int64` (100 uses).
- Behavior is pinned by 224 cgo parity test files plus 31 ported upstream test files.

Two deliberate deviations already exist and matter for the API design:

1. **`Angle`** (`type Angle float64`, radians inside): `LatLng.Lat/Lng` and all four
   `BBox` fields are `Angle`, with `.Rad()` unwrapping at 17 arithmetic sites (e.g.
   [h3index__latLngToCell.go:11](../h3index__latLngToCell.go)). The type flows through the
   entire port transparently and is parity-tested. This is a *good* deviation: it removes
   the single most common H3 user bug (degrees-vs-radians), and — because the internal
   representation *is* the public representation — it costs nothing (see §5.4).
2. **`CellBoundary.Verts` is a growing `[]LatLng`**, whereas C uses an embedded fixed
   `LatLng verts[10]` array. The Go grow loop in
   [faceijk__faceIjkToCellBoundary.go:89](../faceijk__faceIjkToCellBoundary.go) reallocates
   on *every vertex*: measured **6 heap allocations per `cellToBoundary` call** where C
   performs zero (Appendix A, E4). This is a fidelity *and* performance regression to fix.

### 1.3 The export surface is accidental

Verified counts (full lists in §2.3):

| Bucket | Count | Examples |
|---|---|---|
| Exported functions | 14 | `Rad`, `Deg`, 11 `Angle` methods, `FLAG_GET_CONTAINMENT_MODE` |
| Exported types | 30 | `H3Index`, `LatLng`, … but also `CoordIJK`, `FaceIJK`, `Vec2d`, `VertexGraph`, `IterCellsPolygonCompact` |
| Exported constants | ~121 | `M_PI`, `MAX_H3_RES`, `H3_BC_MASK_NEGATIVE`, `E_PENTAGON`, `BUFF_SIZE` |
| Exported lookup tables | 16 | `DIRECTIONS`, `PENTAGON_ROTATIONS`, `NEW_DIGIT_II`, `RES0_BBOXES` |
| **Exported H3 operations** | **0** | — |

The inconsistency is symptomatic of case-driven exporting: rotation tables are public
while the equally internal `baseCellData`/`faceNeighbors` tables are private; the
`IterCellsChildren` struct is fully exported while its siblings are half-exported
(`IterCellsPolygonCompact` has public `Cell`/`Error` and private `res`, `flags`, `bboxes`,
`started` — [iterators_types.go](../iterators_types.go)).

### 1.4 Allocation / conversion hotspots

- `cellToBoundary`: 6 allocs/call (§1.2). Also `polygon__cellBoundaryCrossesGeoLoop.go`
  allocates a normalized boundary copy per call *inside the polyfill hot loop*.
- `polygonToCells` allocates its search/found hash arrays and a per-hex ring buffer
  internally (C does too — inherent to the algorithm, apart from the ring buffer which C
  keeps on the stack).
- `compactCells` allocates three working arrays per call (C mallocs equivalently).
- Linked-geo output allocates one node per vertex/loop/polygon (mirrors C malloc;
  inherent).
- `gridDisk`-family functions use the caller's `out` buffer as a **linear-probing hash
  set** ([algos__gridDiskDistancesInternal.go](../algos__gridDiskDistancesInternal.go)), so
  any buffer-reuse API **must zero the buffer before each call** — a real but cheap cost
  (measured: wrapper with `clear` + size computation adds ~58 ns to a 325 ns k=2 disk, 0
  allocs; Appendix A, E2).

### 1.5 Testing state

Strong: 75/75 public functions parity-tested against the real C objects (compiled from
pristine upstream sources via the `h3lib_*_c2go.c` shim scheme — no vendored C code), plus
31 ported upstream test files, plus exhaustive tests (bounded by partial base-cell
iteration rather than `testing.Short()`).

Gaps:

- **CI is red, for two stacked reasons.** (a) The first CI step, `make fmt`, fails:
  8 files at HEAD are not `gofmt -s` clean (`algos_cellsToLinkedMultiPolygon*.go`,
  `algos_cgo.go`, `algos_h3SetToVertexGraph_parity_test.go`, `baseCells_testBaseCells*`,
  `bbox_testBBoxInternal_test.go`, `polyfill_cellToBBox_test.go`). (b) Behind it, the
  test step would fail anyway: interop files are gated on `//go:build cgo` alone (205 of
  224 parity tests likewise), so *any* cgo-enabled build — including CI's
  `go test -race ./...` and gopls in a default IDE session — tries to compile
  `#include "mathExtensions.c"` without the include paths that only `make test-c2go`
  provides. Reproduced locally: `CGO_ENABLED=1 go test -race ./...` →
  `fatal error: 'mathExtensions.c' file not found`. The intended design (per
  `C2GO_README.md`, since folded into `CONTRIBUTING.md`) was `cgo && c2go`; only 14
  parity files and 3 interop files actually
  say that.
- No public API to test yet: no allocation assertions, no fuzz targets, no examples, no
  API-surface lock. §9 addresses all of these.
- Build-tag inconsistency: 205×`cgo` / 14×`cgo && c2go` / 5×`c2go` across parity files.

### 1.6 Maintenance risk against future upstream releases

Currently *low* at the implementation layer (that is the repo's whole design: one file per
C function, attribution comments, parity harness that can be pointed at a new version via
`make test-c2go H3VER=4.4.0`) — but three things undermine it:

1. Filename convention drift: 49 of the 75 C-*public* functions live in
   double-underscore files (e.g. `algos__gridDisk.go`), although `__` is documented to
   mean "C-internal `_`-prefixed helper"; `h3Index_getBaseCellNumber.go` breaks the
   lowercase-module-prefix convention. Mapping therefore must rely on the attribution
   comments, not filenames — the committed `tools/apiinventory` does exactly that.
2. Attribution format variance (`latLng.c::H3_EXPORT(edgeLengthM)`, `(static function)`
   suffixes, one comma-separated multi-constant attribution) breaks naive greps; the tool
   normalizes these.
3. The stale top-level docs actively mislead (they describe a package layout and an
   oracle protocol that do not exist).

---

## 2. API inventory

### 2.1 Mechanically generated inventory

`tools/apiinventory` (committed with this document) parses `h3api.h.in` for `H3_EXPORT`
declarations and the repo for attribution comments, and emits
[`docs/c-api-inventory.csv`](./c-api-inventory.csv) — 309 rows: **75 C-public functions
(75 ported, 0 missing)** + 234 C-internal helper attributions (197 funcs, 26 vars, 10
consts, 1 type). Regenerate with:

```sh
go run ./tools/apiinventory > docs/c-api-inventory.csv
```

C-internal ports per module: coordijk.c 28, h3Index.c 20, polyfill.c 18, utility.c 17,
faceijk.c 15, bbox.c 14, baseCells.c 13, algos.c 12, latLng.c 12, linkedGeo.c 12,
vertex.c 11, h3Index.h 10, polygon.c 9, vertexGraph.c 9, localij.c 8, iterators.c 7,
polygonAlgos.h 6, vec2d.c 3, vec3d.c 3, others 7.

Deprecations: none are marked deprecated in the 4.3.0 header.
`polygonToCellsExperimental` / `maxPolygonToCellsSizeExperimental` are flagged
*experimental, subject to change in minor versions*; `cellsToLinkedMultiPolygon` is
documented upstream as a "binding-only concept". The pre-4.0 `hexRange` family no longer
exists in 4.x.

### 2.2 C public API → proposed public Go API

Legend: *impl* = existing unexported Go port (unchanged); *public* = proposed new wrapper.
`Append*` forms and `iter.Seq` forms are listed where they exist; conventions in §4/§5.
"—" in *public* means intentionally not exposed, with the reason.

**Indexing & inspection (h3Index.c, baseCells.c)**

| C function | Go impl (exists) | Proposed public Go API |
|---|---|---|
| `latLngToCell` | `latLngToCell` | `LatLngToCell(ll LatLng, res int) (Cell, error)` |
| `cellToLatLng` | `cellToLatLng` | `Cell.LatLng() (LatLng, error)` |
| `cellToBoundary` | `cellToBoundary` | `Cell.Boundary() (CellBoundary, error)` |
| `getResolution` | `getResolution` | `Cell.Resolution() int` (also on `DirectedEdge`, `Vertex`) |
| `getBaseCellNumber` | `getBaseCellNumber` | `Cell.BaseCellNumber() int` |
| `stringToH3` | `stringToH3` | `ParseCell(s string) (Cell, error)`; `ParseDirectedEdge`; `ParseVertex` (mode-validated) |
| `h3ToString` | `h3ToString` | `Cell.String() string` (+ `DirectedEdge`, `Vertex`); `MarshalText`/`UnmarshalText` on all three |
| `isValidCell` | `isValidCell` | `Cell.IsValid() bool` |
| `isResClassIII` | `isResClassIII` | `Cell.IsResClassIII() bool` |
| `isPentagon` | `isPentagon` | `Cell.IsPentagon() bool` |
| `getIcosahedronFaces` | `getIcosahedronFaces` | `Cell.IcosahedronFaces() ([]int, error)` (≤ 5, `-1` holes pruned) |
| `maxFaceCount` | `maxFaceCount` | — sizing detail hidden inside `IcosahedronFaces` |
| `getPentagons` | `getPentagons` | `Pentagons(res int) ([]Cell, error)` |
| `pentagonCount` | `pentagonCount` | `const NumPentagons = 12` |
| `getRes0Cells` | `getRes0Cells` | `Res0Cells() []Cell` |
| `res0CellCount` | `res0CellCount` | `const NumRes0Cells = 122` |
| `describeH3Error` | `describeH3Error` | — surfaces as the `Error()` text of the sentinel errors (§4.5) |

**Hierarchy (h3Index.c)**

| C function | Go impl | Proposed public Go API |
|---|---|---|
| `cellToParent` | `cellToParent` | `Cell.Parent(res int) (Cell, error)` |
| `cellToCenterChild` | `cellToCenterChild` | `Cell.CenterChild(res int) (Cell, error)` |
| `cellToChildren` | `cellToChildren` | `Cell.Children(res int) ([]Cell, error)`; `Cell.AppendChildren(dst []Cell, res int) ([]Cell, error)`; `Cell.ChildrenSeq(res int) iter.Seq[Cell]` |
| `cellToChildrenSize` | `cellToChildrenSize` | `Cell.NumChildren(res int) (int64, error)` |
| `cellToChildPos` | `cellToChildPos` | `Cell.ChildPos(parentRes int) (int64, error)` |
| `childPosToCell` | `childPosToCell` | `Cell.ChildAtPos(pos int64, res int) (Cell, error)` (receiver = parent) |
| `compactCells` | `compactCells` | `CompactCells(cells []Cell) ([]Cell, error)`; `AppendCompactCells(dst, cells []Cell) ([]Cell, error)` |
| `uncompactCells` | `uncompactCells` | `UncompactCells(cells []Cell, res int) ([]Cell, error)`; `AppendUncompactCells(dst, cells []Cell, res int) ([]Cell, error)` |
| `uncompactCellsSize` | `uncompactCellsSize` | `UncompactCellsSize(cells []Cell, res int) (int64, error)` |

**Traversal (algos.c, localij.c)**

| C function | Go impl | Proposed public Go API |
|---|---|---|
| `gridDisk` | `gridDisk` | `Cell.GridDisk(k int) ([]Cell, error)`; `Cell.AppendGridDisk(dst []Cell, k int) ([]Cell, error)` (holes pruned) |
| `gridDiskDistances` | `gridDiskDistances` | `Cell.GridDiskDistances(k int) ([]Cell, []int32, error)`; `Append` form takes `(dst []Cell, dstDist []int32, k int)` |
| `gridDiskDistancesSafe` | `gridDiskDistancesSafe` | `Cell.GridDiskDistancesSafe(k int)` (same shape) |
| `gridDiskDistancesUnsafe` | `gridDiskDistancesUnsafe` | `Cell.GridDiskDistancesUnsafe(k int)` (same shape) |
| `gridDiskUnsafe` | `gridDiskUnsafe` | `Cell.GridDiskUnsafe(k int) ([]Cell, error)` (ordered; `ErrPentagon` on distortion) |
| `gridDisksUnsafe` | `gridDisksUnsafe` | `GridDisksUnsafe(origins []Cell, k int) ([]Cell, error)` |
| `gridRing` | `gridRing` | `Cell.GridRing(k int) ([]Cell, error)` + `Append` form (holes pruned) |
| `gridRingUnsafe` | `gridRingUnsafe` | `Cell.GridRingUnsafe(k int) ([]Cell, error)` + `Append` form |
| `maxGridDiskSize` | `maxGridDiskSize` | `MaxGridDiskSize(k int) (int64, error)` (for `Append*` pre-sizing) |
| `maxGridRingSize` | `maxGridRingSize` | `MaxGridRingSize(k int) (int64, error)` |
| `gridDistance` | `gridDistance` | `Cell.GridDistance(other Cell) (int, error)` |
| `gridPathCells` | `gridPathCells` | `Cell.GridPath(other Cell) ([]Cell, error)` + `Append` form |
| `gridPathCellsSize` | `gridPathCellsSize` | `Cell.GridPathLen(other Cell) (int, error)` |
| `cellToLocalIj` | `cellToLocalIj` | `CellToLocalIJ(origin, c Cell) (CoordIJ, error)` |
| `localIjToCell` | `localIjToCell` | `LocalIJToCell(origin Cell, ij CoordIJ) (Cell, error)` |

**Regions (algos.c, polyfill.c, linkedGeo.c)**

| C function | Go impl | Proposed public Go API |
|---|---|---|
| `polygonToCells` | `polygonToCells` | `PolygonToCells(p GeoPolygon, res int) ([]Cell, error)` + `Append` form |
| `maxPolygonToCellsSize` | `maxPolygonToCellsSize` | `MaxPolygonToCellsSize(p GeoPolygon, res int) (int64, error)` |
| `polygonToCellsExperimental` | `polygonToCellsExperimental` | `PolygonToCellsExperimental(p GeoPolygon, res int, mode ContainmentMode) ([]Cell, error)` + `Append` form; doc-flagged experimental |
| `maxPolygonToCellsSizeExperimental` | `maxPolygonToCellsSizeExperimental` | `MaxPolygonToCellsSizeExperimental(p GeoPolygon, res int, mode ContainmentMode) (int64, error)` |
| `cellsToLinkedMultiPolygon` | `cellsToLinkedMultiPolygon` | `CellsToMultiPolygon(cells []Cell) ([]GeoPolygon, error)` (linked list converted to slices) |
| `destroyLinkedMultiPolygon` | `destroyLinkedMultiPolygon` | — GC makes it meaningless in Go |

**Directed edges (directedEdge.c)**

| C function | Go impl | Proposed public Go API |
|---|---|---|
| `areNeighborCells` | `areNeighborCells` | `Cell.IsNeighbor(other Cell) (bool, error)` |
| `cellsToDirectedEdge` | `cellsToDirectedEdge` | `Cell.DirectedEdgeTo(dest Cell) (DirectedEdge, error)` |
| `isValidDirectedEdge` | `isValidDirectedEdge` | `DirectedEdge.IsValid() bool` |
| `getDirectedEdgeOrigin` | `getDirectedEdgeOrigin` | `DirectedEdge.Origin() (Cell, error)` |
| `getDirectedEdgeDestination` | `getDirectedEdgeDestination` | `DirectedEdge.Destination() (Cell, error)` |
| `directedEdgeToCells` | `directedEdgeToCells` | `DirectedEdge.Cells() (origin, destination Cell, err error)` |
| `originToDirectedEdges` | `originToDirectedEdges` | `Cell.DirectedEdges() ([]DirectedEdge, error)` (≤ 6; pentagon hole pruned) |
| `directedEdgeToBoundary` | `directedEdgeToBoundary` | `DirectedEdge.Boundary() (CellBoundary, error)` |

**Vertexes (vertex.c)**

| C function | Go impl | Proposed public Go API |
|---|---|---|
| `cellToVertex` | `cellToVertex` | `Cell.Vertex(vertexNum int) (Vertex, error)` |
| `cellToVertexes` | `cellToVertexes` | `Cell.Vertexes() ([]Vertex, error)` (5 for pentagons; hole pruned) |
| `vertexToLatLng` | `vertexToLatLng` | `Vertex.LatLng() (LatLng, error)` |
| `isValidVertex` | `isValidVertex` | `Vertex.IsValid() bool` |

**Metrics & misc (latLng.c)**

| C function | Go impl | Proposed public Go API |
|---|---|---|
| `greatCircleDistanceKm/M/Rads` | same names | `GreatCircleDistanceKm/M/Rads(a, b LatLng) float64` |
| `cellAreaKm2/M2/Rads2` | same names | `Cell.AreaKm2/AreaM2/AreaRads2() (float64, error)` |
| `edgeLengthKm/M/Rads` | same names | `DirectedEdge.LengthKm/LengthM/LengthRads() (float64, error)` |
| `getHexagonAreaAvgKm2/M2` | same names | `HexagonAreaAvgKm2/M2(res int) (float64, error)` |
| `getHexagonEdgeLengthAvgKm/M` | same names | `HexagonEdgeLengthAvgKm/M(res int) (float64, error)` |
| `getNumCells` | `getNumCells` | `NumCells(res int) (int64, error)` |
| `degsToRads` / `radsToDegs` | same names | — replaced by `Angle`: `Deg(x).Rad()` / `Rad(x).Deg()` |

**Go-only additions** (no C counterpart): `Angle`/`Deg`/`Rad` (already exist);
`NewLatLng(lat, lng Angle) LatLng` and `LatLngDegs(lat, lng float64) LatLng` conveniences;
`CellsAtRes(res int) iter.Seq[Cell]` (wraps the ported `iterInitRes`/`iterStepRes` from
iterators.c, which is not part of h3api.h.in but is documented upstream API);
`Cell.ChildrenSeq`; text marshaling; curated constants (`MaxResolution`, `NumBaseCells`,
`NumPentagons`, `NumRes0Cells`, `MaxCellBoundaryVerts`, `VersionMajor/Minor/Patch`).

**Intentionally never exposed**: all 234 C-internal helpers; the C iterator structs
(superseded by `iter.Seq`); `destroy*` functions; the `utility.c` print helpers (already
unexported; test-support only); `describeH3Error`, `degsToRads`, `radsToDegs`,
`maxFaceCount` (absorbed as noted above); `GeoMultiPolygon` (defined in the C header but
used by no C public function — confirmed).

### 2.3 Currently exported identifiers and their fate

- **Keep (already correct)**: `Angle`, `Deg`, `Rad` + methods; `LatLng`, `GeoLoop`,
  `GeoPolygon`, `CoordIJ`, `ContainmentMode` (constants renamed Go-style, §7).
- **Reshape**: `CellBoundary` becomes a fixed-array value type (§4.4); `H3Index` becomes
  the unexported alias `h3Index = Cell` (§4.2); `H3Error`/`E_*` become unexported, replaced
  publicly by sentinel errors (§4.5).
- **Unexport (mechanical case change, names otherwise preserved)**: internal types
  (`CoordIJK`, `FaceIJK`, `FaceOrientIJK`, `BBox`, `Vec2d`, `Vec3d`, `VertexGraph`,
  `VertexNode`, `BaseCellData`, `BaseCellRotation`, `PentagonDirectionFaces`, `Direction`,
  `Overage`, `LongitudeNormalization`, 4 `IterCells*` types, `LinkedLatLng`,
  `LinkedGeoLoop`, `LinkedGeoPolygon`), ~121 C-macro constants (`M_PI`, `MAX_H3_RES`,
  `H3_*_MASK*`, `CENTER_DIGIT`…, `BUFF_SIZE`, …), 16 lookup tables (`DIRECTIONS`,
  `PENTAGON_ROTATIONS*`, `NEW_DIGIT_*`, `NEW_ADJUSTMENT_*`, `UNIT_VECS`,
  `FAILED_DIRECTIONS`, `MAX_EDGE_LENGTH_RADS`, `NORTH/SOUTH_POLE_CELLS`,
  `VALID_RANGE_BBOX`, `RES0_BBOXES`), and `FLAG_GET_CONTAINMENT_MODE` /
  `FLAG_CONTAINMENT_MODE_MASK`. Naming rules and the one known collision (`DIRECTIONS` vs
  the existing unexported `directions`) are specified in §7.

---

## 3. Proposed package architecture

### 3.1 Options considered

**Option A — single root package, two file layers (recommended).** The ported C-shaped
code (unexported) and the public typed API live in one `package h3`. Public files get
plain topical names (`cell.go`, `edge.go`, `vertex.go`, `latlng.go`, `boundary.go`,
`hierarchy.go`, `traversal.go`, `region.go`, `localij.go`, `metrics.go`, `errors.go`,
`iter.go`, `doc.go`); ported files keep their `<cfile>_<function>.go` names.

**Option B — `internal/h3impl` for the ported code, thin public root.** Structurally
attractive but fails on a hard language constraint: the public `Cell` must be a defined
type *with methods*, and methods can only be declared in the package that defines the
type. If `Cell` is defined in the root, `internal/h3impl` cannot mention it (import
cycle), so the implementation must use `uint64`/its own `H3Index` — and then `[]Cell` and
`[]h3impl.H3Index` are **different types that cannot be converted without copying or
`unsafe`** (Go spec, Conversions: slice types are convertible only if they have identical
element types, modulo struct tags). The reverse trick — define `Cell` in the internal
package and alias it publicly (`type Cell = h3impl.Cell`) — compiles, but then all methods
live in an `internal` package: `pkg.go.dev` documentation for the root shows an alias
whose method set is documented nowhere reachable, and users cannot see the type's docs.
Option B therefore forces exactly the copies/`unsafe` this project wants to avoid, plus a
500-file move. Rejected (DR-001).

**Option C — multiple internal packages mirroring C modules** (`internal/coordijk`,
`internal/faceijk`, …). Everything wrong with B, plus qualifier churn inside the ported
bodies (`coordijk.IjkAdd(...)`), which *does* damage line-level traceability, plus
import-cycle traps (h3Index.c ↔ algos.c call each other freely). Rejected.

**Option D — generated compatibility layer** (code-gen the wrappers from an inventory).
The wrappers are not mechanical enough to generate: hole-pruning policy, dst-buffer
semantics, error mapping, and receiver choice vary per function. A generator would be more
code than the ~60 wrappers themselves. Rejected; generation is used only for *inventory
and verification* (tools/apiinventory), not for API code.

### 3.2 Why the single package holds up

- **Zero-cost bridging is only possible here.** With `type Cell uint64` and
  `type h3Index = Cell` in the same package, the entire ported layer operates on the
  public type without knowing it. Verified: the alias builds with **zero edits to any
  ported file**, and both the pure-Go suite and the full cgo parity suite pass unchanged
  (Appendix A, E1).
- **godoc hygiene does not require package separation** — it requires unexporting, which
  is needed anyway (§7). After the sweep, `go doc` shows only the intended API.
- **Testability**: parity tests are white-box (they call unexported functions) and must
  stay in-package; a package split would force an awkward export-for-test shim layer.
- **Scale**: ~600 files in one package is unusual but harmless — compilation units in Go
  are packages, and the compiler handles this size trivially (full build ≈ seconds).
  Discoverability inside the repo comes from the filename convention, not the package
  graph.
- The **`tools/apiinventory`** command is the only additional package (a `main`, not
  importable).

### 3.3 File-layer conventions

| Layer | Files | Naming | Visibility |
|---|---|---|---|
| Ported C code | `<cfile>_<func>.go`, `<cfile>__<func>.go`, `*_constants.go`, `*_types.go` | C names preserved (case-lowered where currently exported) | unexported |
| Public API | `cell.go`, `edge.go`, `vertex.go`, `latlng.go`, `boundary.go`, `hierarchy.go`, `traversal.go`, `region.go`, `localij.go`, `metrics.go`, `errors.go`, `iter.go`, `doc.go` | Go names | exported |
| Parity harness | `*_cgo.go`, `h3lib_*_c2go.c`, `*_parity_test.go` | as today | build-tagged `cgo && c2go` |

Every public wrapper carries a machine-readable back-reference in its doc comment:

```go
// GridDisk returns all cells within grid distance k of c, in no particular
// order. ...
//
// H3 C API: gridDisk.
func (c Cell) GridDisk(k int) ([]Cell, error)
```

`tools/apiinventory` will be extended (Phase 6) to parse `H3 C API:` lines and enforce
that every C-public function has either a public wrapper or an entry in the documented
omissions list (§2.2 "intentionally never exposed") — turning API completeness into a CI
gate.

---

## 4. Public type model

### 4.1 Overview

| Type | Declaration | Kind | Methods | Zero value | String/Text |
|---|---|---|---|---|---|
| `Cell` | `type Cell uint64` | defined | full API surface | invalid index (== `h3Null`); `IsValid()==false` | 15-digit lowercase hex; `MarshalText`/`UnmarshalText` |
| `DirectedEdge` | `type DirectedEdge uint64` | defined | edge ops | invalid | same |
| `Vertex` | `type Vertex uint64` | defined | vertex ops | invalid | same |
| `h3Index` | `type h3Index = Cell` | **alias**, unexported | — | — | — |
| `Angle` | `type Angle float64` (radians) | defined (exists) | `Deg`, `Rad`, trig, `String` ("…°") | 0 rad | no text marshaling (see §12-Q5) |
| `LatLng` | `struct { Lat, Lng Angle }` | defined (exists) | `Cell(res)`, `String` | 0°,0° (meaningful) | none (GeoJSON ambiguity, §12-Q5) |
| `CellBoundary` | `struct { verts [10]LatLng; n int32 }` | defined value type | `Verts() []LatLng`, `Len()`, `At(i)` | empty boundary | — |
| `GeoLoop` | `type GeoLoop []LatLng` | defined (exists) | maybe `IsClockwise()` later | nil = empty | — |
| `GeoPolygon` | `struct { GeoLoop GeoLoop; Holes []GeoLoop }` | defined (exists) | — | empty polygon | — |
| `CoordIJ` | `struct { I, J int32 }` | defined (exists) | — | origin | — |
| `ContainmentMode` | `type ContainmentMode int` | defined (exists) | `String()` | `ContainmentCenter` (=0, matches C) | — |

Explicitly *not* introduced: a public umbrella `Index` type. See §4.3.

### 4.2 `Cell` and the alias trick (the load-bearing decision)

The C layer uses `H3Index` for cells, edges, and vertexes alike, and every *large* (non
fixed-size) index collection in the entire 75-function API carries **cells only**
(gridDisk/ring/path, children, compact/uncompact, polygonToCells, res0, pentagons).
Edge and vertex collections are fixed-size ≤ 6 (`originToDirectedEdges`,
`cellToVertexes`) or scalar.

Therefore: make `Cell` the defined type the implementation *actually uses*, via a
single-line change to `h3api_types.go`:

```go
// Cell is an H3 cell (hexagon or pentagon) index.
type Cell uint64

// h3Index mirrors C H3Index. It is an alias of Cell so the mechanically
// ported code and the public API share one type: []h3Index IS []Cell.
type h3Index = Cell
```

Consequences, all verified experimentally (Appendix A, E1–E3):

- `[]h3Index` and `[]Cell` are the *same type* — dst buffers pass through with no
  conversion, no copy, no `unsafe`, no generics.
- Ported code needs **zero body edits** (only the `H3Index` → `h3Index` spelling, a
  word-boundary sed across 117 files, executed in the unexport sweep; the alias compiles
  under the old exported spelling too, which allows landing the type change and the rename
  as two separately verifiable commits).
- Scalar conversions `DirectedEdge(x)`, `Vertex(x)`, `uint64(c)` compile to nothing.
- Inside edge/vertex wrappers, a stack `[6]h3Index` is filled by the ported function and
  converted element-wise into the result type — no heap traffic.
- Semantic note: internal code occasionally holds an edge or vertex index in a variable
  typed `h3Index`(=`Cell`). That is exactly as (un)typed as the C original — the type
  distinction is a *public-boundary* guarantee, enforced where values enter (parsing,
  `UnmarshalText`, explicit conversion documented as unchecked) and by validating
  operations, not a whole-program invariant. This mirrors how the C library itself works.

`uint64` (not uber's `int64`) stays: it matches C `H3Index`, the port's bit arithmetic,
and avoids sign-extension traps in mask code. Interop with uber types is a documented
explicit conversion.

### 4.3 Why no public `Index` type

An `Index uint64` sitting "above" `Cell`/`DirectedEdge`/`Vertex` would need conversions at
every use (defined types don't auto-convert), duplicate every mode-generic operation, and
add a second way to hold an invalid value. The two genuine mode-generic needs are served
without it:

- *Parsing unknown input*: `ParseCell`/`ParseDirectedEdge`/`ParseVertex` validate the mode
  bits and return typed values (unlike C `stringToH3` and uber's `IndexFromString`, which
  accept anything and defer validation to the caller — a documented wart we do not copy).
- *Generic helpers*: an unexported constraint `type index interface { Cell | DirectedEdge | Vertex }`
  supports shared parse/format/marshal implementations (§6). If demand appears, a public
  `IsValidIndex[T index]` can be added later without breaking anything.

Validation policy: constructors/parsers validate; pure bit accessors (`Resolution`,
`BaseCellNumber`, `IsValid`, `String`) never error; algorithmic operations return the C
layer's domain errors. Raw conversions (`Cell(0x8f28...)`) are legal and unchecked, same
as C — documented with a pointer to `IsValid`.

### 4.4 `CellBoundary` as a value type

C: `typedef struct { int numVerts; LatLng verts[MAX_CELL_BNDRY_VERTS]; } CellBoundary;` —
a 168-byte stack value. The current port's `Verts []LatLng` slice deviates and costs 6
allocs/call (§1.2). Proposal: restore C fidelity *and* fix the hotspot:

```go
// CellBoundary is the boundary of a cell or edge, in ccw order.
// The zero value is an empty boundary. It is a value type; copying is cheap
// and no heap allocation is involved.
type CellBoundary struct {
    verts [MaxCellBoundaryVerts]LatLng // MaxCellBoundaryVerts = 10, C parity
    n     int32
}

func (b *CellBoundary) Len() int         { return int(b.n) }
func (b *CellBoundary) At(i int) LatLng  { return b.verts[i] }
func (b *CellBoundary) Verts() []LatLng  { return b.verts[:b.n] } // aliases b
```

`Cell.Boundary()` returns it by value; typical usage keeps everything on the stack.
Fields stay unexported so `n` cannot be corrupted; `Verts()` uses a pointer receiver so
the returned slice aliases the caller's variable rather than a temporary. The internal
ported functions (`_faceIjkToCellBoundary`, `_faceIjkPentToCellBoundary`,
`cellToBoundary`, `directedEdgeToBoundary`, `polygon__cellBoundaryCrossesGeoLoop`) are
updated to fill the fixed array exactly as the C code does — this *increases* C fidelity
and is re-verified by the existing parity tests, which compare values, not storage.

### 4.5 Errors

Public: 15 package-level sentinels mirroring `H3ErrorCodes` (the shape the `.backup` draft
and uber both use, minus uber's frozen `ErrRsolutionMismatch` typo):

```go
var (
    ErrFailed              = errors.New("h3: the operation failed")           // E_FAILED
    ErrDomain              = errors.New("h3: argument was outside of acceptable range")
    ErrLatLngDomain        = errors.New("h3: latitude or longitude arguments were outside of acceptable range")
    ErrResolutionDomain    = errors.New("h3: resolution argument was outside of acceptable range")
    ErrCellInvalid         = errors.New("h3: cell argument was not valid")
    ErrDirectedEdgeInvalid = errors.New("h3: directed edge argument was not valid")
    // … one per C code 1–15, message text taken from describeH3Error …
)

func toErr(e h3Error) error { // array lookup, no allocation, inlines
    if int(e) < len(errTable) { return errTable[e] }
    return ErrFailed
}
```

`errors.Is` works naturally. The C `H3Error`/`E_*` identifiers become unexported
(`h3Error`, `eSuccess`, …) — they remain the currency of the ported layer. Should users
need numeric codes for cross-language stability, a `Code(err) (int, bool)` accessor can be
added later without redesign (§12-Q7).

### 4.6 Formatting, parsing, marshaling

- `String()` on `Cell`/`DirectedEdge`/`Vertex`: lowercase hex without `0x` (C
  `h3ToString` behavior; matches uber and h3geo tooling). Implemented via
  `strconv.AppendUint` into a stack buffer — no fmt, no allocation beyond the string
  itself.
- `ParseCell` etc.: accept the same forms C `stringToH3` accepts (hex, with the mode then
  validated). Unlike uber's `IndexFromString`, parse errors are returned, not swallowed.
- `MarshalText`/`UnmarshalText` on all three index types (JSON interop follows for free;
  hex strings survive JavaScript number precision, the standard H3-in-JSON practice).
- `LatLng`/`Angle`: **no** JSON/text marshaling in v0 — any choice (radians? degrees?
  `[lng, lat]` GeoJSON order?) silently mis-integrates with someone; explicit accessors
  are safer until a `geojson` helper package exists (§12-Q5).

### 4.7 Compatibility with existing users

There are none that could exist: no tags, no releases, no callable API (the 14 exported
funcs are the `Angle` helpers). `H3Index`, `CellBoundary{NumVerts, Verts}`,
`LinkedGeoPolygon` etc. disappear or change shape without deprecation ceremony (DR-006).

---

## 5. Zero-copy API strategy

Cost classes used below: **(0)** truly free (type identity / compiles away), **(s)**
syntactic conversion, no runtime cost, **(k)** fixed small cost independent of n,
**(A)** allocation inherent to the requested result, **(X)** avoidable copy — none of
these remain in the design.

### 5.1 Scalars — (0)/(s)

`Cell(x)`, `uint64(c)`, `DirectedEdge(v)` are representation-preserving conversions;
wrappers like `Cell.Resolution() int { return int(getResolution(c)) }` fully inline
(Appendix A, E3). `int` ↔ `int32`/`int64` narrowing for `res`, `k`, `vertexNum` is (s)
with explicit range checks already performed by the C layer (`E_RES_DOMAIN` etc.).

### 5.2 Cell slices — (0) by type identity

The alias makes every `out []h3Index` parameter accept `[]Cell` directly. The dst-buffer
wrapper pattern (validated in Appendix A, E2: **cold 1 alloc — the result itself; warm 0
allocs**):

```go
func (c Cell) AppendGridDisk(dst []Cell, k int) ([]Cell, error) {
    var sz int64
    if err := maxGridDiskSize(int32(k), &sz); err != eSuccess {
        return dst, toErr(err)
    }
    n := len(dst)
    dst = slices.Grow(dst, int(sz))[:n+int(sz)]
    win := dst[n:]
    clear(win) // gridDisk uses the out buffer as a linear-probing hash set
    if err := gridDisk(c, int32(k), win); err != eSuccess {
        return dst[:n], toErr(err)
    }
    m := compactNonNull(win) // in-place; preserves relative order
    return dst[:n+m], nil
}

func (c Cell) GridDisk(k int) ([]Cell, error) { return c.AppendGridDisk(nil, k) }
```

Notes:

- `clear` on the window is mandatory (hash-set probing, §1.4) and costs ~10 ns at k=2;
  it exists in the C usage contract too ("array must be zeroed").
- **Hole pruning** (`compactNonNull`) is in-place — no second buffer. Pruning policy per
  function follows the C contract inventory (§2.1): prune `H3_NULL` for
  gridDisk/gridRing/polygonToCells outputs and pentagon holes in
  `DirectedEdges`/`Vertexes`/`IcosahedronFaces`; *trim* (not filter) for
  compact/uncompact. `GridDiskDistances*` prunes cells and distances in tandem.
- `Append*` semantics are true append (results after `len(dst)`), matching
  `strconv.AppendInt` conventions; pass `buf[:0]` to reuse.

### 5.3 Fixed-size outputs — (k), stack only

```go
func (c Cell) DirectedEdges() ([]DirectedEdge, error) {
    var raw [6]h3Index                     // stack
    if err := originToDirectedEdges(c, &raw); err != eSuccess { return nil, toErr(err) }
    out := make([]DirectedEdge, 0, 6)      // the requested result — (A)
    for _, e := range raw {
        if e != h3Null { out = append(out, DirectedEdge(e)) }
    }
    return out, nil
}
```

The 6-element element-wise conversion is (k); the single `make` is (A). An `Append` form
is unnecessary at this size but trivially added if profiling ever disagrees.

### 5.4 Coordinates, boundaries, polygons — (0) because public type == internal type

This is the second structural decision: the public geometry types (`LatLng` with `Angle`
fields, `GeoLoop`, `GeoPolygon`) **are** the types the ported code computes with. A
`GeoPolygon` with thousands of vertices enters `PolygonToCells` without any per-vertex
conversion or copy — the wrapper passes `&p` through. Had the public types used degrees
(as uber's do), every polygon input and boundary output would pay an O(n) convert-copy.
This is the concrete, measured reason to keep `Angle` (DR-003) — plus it makes the
degrees/radians bug class unrepresentable.

`CellBoundary` as a value type (§4.4) removes the 6-alloc hotspot; `Boundary()` is then
(k) — one 168-byte stack-to-stack copy.

### 5.5 Linked multipolygon output — (A) only

`CellsToMultiPolygon` walks the ported linked structure once to count, allocates the exact
`[]GeoPolygon`/`[]GeoLoop`/`[]LatLng` result, and fills it. The linked nodes themselves
are the algorithm's own working allocation (C mallocs identically); the conversion adds
only the result buffers the user asked for. No exposed linked types.

### 5.6 Iterators — (0) steady-state

Range-over-func wrappers over the ported iterator structs; measured **0 allocs** for a
full 343-child iteration including the closure (Appendix A, E5):

```go
func (c Cell) ChildrenSeq(res int) iter.Seq[Cell] {
    return func(yield func(Cell) bool) {
        var it iterCellsChildren
        iterInitParent(c, int32(res), &it)
        for ; it.h != h3Null; iterStepChild(&it) {
            if !yield(it.h) { return }
        }
    }
}
```

Invalid input yields an empty sequence (exactly the C null-iterator contract). For the
polygon iterator, which can fail mid-iteration, the shape is `iter.Seq2[Cell, error]` or a
final-`Err()` handle — deferred to Phase 5 with a decision in §12-Q6. This feature sets
the `go` directive floor at **1.23**; recommendation: `go 1.24` (adds `testing.B.Loop` for
the benchmark suite; every currently supported toolchain ≥ 1.25 satisfies it). Upgrading
further brings nothing this library needs.

### 5.7 Where allocation is inherent — and stays

`Children`, `GridDisk`, `PolygonToCells` convenience forms allocate exactly the result
(cold path, 1 alloc). `compactCells`'s three working arrays and `polygonToCells`'s hash
arrays are the algorithm's own cost (same in C); they could later be offered a scratch
buffer via an options struct if profiling demands it — explicitly out of scope now.

---

## 6. Generics assessment

Generics are needed in exactly **two small places**; everywhere else they would be net
harm. Assessment against actual code:

| Site | Verdict | Rationale |
|---|---|---|
| `castSlice[From, To ~uint64]` (`h3.go.backup:920`, `unsafe`) | **Delete — obsoleted** | The alias makes `[]Cell` ≡ `[]h3Index`; the only remaining cross-type collections are ≤ 6 elements (§5.3). The generic+unsafe bridge solved a problem the type system now solves. |
| Parse/format/marshal for `Cell`/`DirectedEdge`/`Vertex` | **Use generics (unexported)** | `func parseIndex[T ~uint64](s string, mode int) (T, error)` + shared `appendHex[T ~uint64]`. Removes 3× duplication of fiddly hex/validation code; instantiates to identical machine code for one gcshape (`uint64`); inlines like any small function; zero allocation change. Public generic surface: none initially. |
| Pentagon-hole pruning of `[6]h3Index` into `[]DirectedEdge` / `[]Vertex` | **Borderline — allow one tiny helper** | `pruneConvert[T ~uint64](raw []h3Index) []T`. Two call sites; fine either as generic or duplicated; the generic keeps the pruning policy in one place. |
| Ported algorithm layer (e.g. making `gridDisk` generic over `~uint64`) | **Reject** | Would change 117 files' signatures for zero benefit (alias already gives type identity), hurt line-level C traceability, and risk gcshape/inlining regressions in hot recursion (`_gridDiskDistancesInternal`). |
| uber-style public `Index` constraint (`interface{ Cell \| DirectedEdge \| Vertex }`) | **Defer** | No operation in the proposed surface needs it; add later compatibly if a `IsValidIndex[T]`-style need materializes. |

---

## 7. Naming strategy

### 7.1 Rules by identifier class

| Class | Rule | Examples |
|---|---|---|
| Ported C functions | Keep C name verbatim, unexported (as today); `_`-prefixed C statics keep the `_` prefix | `gridDisk`, `_faceIjkToH3` |
| Ported C types (currently exported) | Lowercase first letter; all-caps prefixes drop to lowercase as a word | `CoordIJK`→`coordIJK`, `FaceIJK`→`faceIJK`, `BBox`→`bbox`, `Vec2d`→`vec2d`, `IterCellsChildren`→`iterCellsChildren` |
| Ported C constants/macros | `SCREAMING_SNAKE` → lowerCamel: lowercase first segment, TitleCase the rest | `MAX_H3_RES`→`maxH3Res`, `M_PI`→`mPi`, `E_PENTAGON`→`ePentagon`, `H3_NULL`→`h3Null`, `NUM_HEX_VERTS`→`numHexVerts` |
| Ported C tables | Same rule; **on collision with an existing unexported name, prefix with the C module** | `PENTAGON_ROTATIONS`→`pentagonRotations`; known collision: `DIRECTIONS`(algos.c)→`algosDirections` because unexported `directions` already exists in the vertex module |
| Public functions/methods | Go-style, C v4 vocabulary kept, `get` prefix dropped, no `H3` stutter | `GridDisk`, `Parent`, `IsValidCell`→`Cell.IsValid` |
| Public types/constants | Go-style; curated constants get Go names | `MaxResolution`, `NumPentagons`, `MaxCellBoundaryVerts` |
| Error values | `Err` + C code name sans `E_` | `E_NOT_NEIGHBORS`→`ErrNotNeighbors` |
| Acronyms | `IJ`, `IJK`, `CCW`, `CW` uppercase in exported names; `Latitude/Longitude` stay `Lat/Lng` (H3 vocabulary) | `CellToLocalIJ`, `CoordIJ` |
| Deprecated aliases | none (no users to migrate; DR-006) | — |

The unexport renames are **case/format changes only** — the identifier remains recognizably
the C name, every declaration keeps its `// Ported from H3 C:` line, and the rename is
executed by a reviewed script (word-boundary regex over non-test + test files) in one
commit, immediately re-validated by the full parity suite. This satisfies "no blind
renaming": nothing about the structure, file layout, or bodies changes.

### 7.2 Finding the C counterpart of a public operation

Three hops, all mechanical:

1. Public doc comment carries `H3 C API: <cname>` (§3.3).
2. `<cname>` → implementing file via `docs/c-api-inventory.csv` (or grep
   `Ported from H3 C: .*::<cname>`).
3. The implementation file name itself encodes the C module (`algos__gridDisk.go`).

Phase 6 adds CI enforcement that hop 1 exists for every C-public function.

---

## 8. Compatibility and migration strategy

Facts: no git tags, no releases, module path without `/vN`, 332 commits by one author, CI
currently red, README says "Status: WIP". `pkg.go.dev` adoption is implausible given the
package exports no operations. **Conclusion: this is pre-release software; backward
compatibility has no beneficiaries.**

| Option | Cost | Verdict |
|---|---|---|
| Preserve current exported surface | Permanently ship `M_PI`, `H3_BC_MASK_NEGATIVE`, broken `CellBoundary` | Rejected |
| Deprecated aliases for one cycle | ~170 alias declarations nobody will ever use | Rejected |
| **Single breaking sweep, tag `v0.1.0` after the API waves land** | One-time; zero external cost | **Recommended** |
| Major-version bump (`/v2`) | Meaningless before `v1` | N/A |

Versioning plan: stay `v0.x` while the API settles; tag **`v0.1.0`** at Phase 5
completion; consider **`v1.0.0`** only after the first real upstream sync (H3 4.4/4.5)
proves the maintenance workflow, since `v1` freezes the API against exactly those changes.
Note for positioning: uber/h3-go's master now contains `x/h3go`, an unreleased pure-Go H3
port by Uber. This project's differentiators — full 75-function C-API coverage, alias-based
zero-copy dst-buffer/`iter.Seq` APIs, and function-level parity testing — should be locked
in and tagged rather than diluted chasing their surface; where a name is free (GridDisk,
Parent, Children, …) we deliberately match uber's vocabulary to keep migration friction
low in both directions.

---

## 9. Validation and testing strategy

Layered so that `CGO_ENABLED=0 go test ./...` always works for normal users; everything
touching C or external modules is opt-in.

### 9.1 Existing layers (kept, one fix)

- **Ported upstream unit tests** (31 files, pure Go, no tags): kept; the assertions that
  exercise C-public functions are additionally re-pointed at the public wrappers as those
  land, so upstream's own test vectors validate the public layer too.
- **cgo parity tests** (224 files): kept as the ground-truth harness; build tags unified
  to `cgo && c2go` (Phase 0) so plain cgo-enabled builds (CI `-race`, gopls) stop
  breaking. Parity tests keep targeting the *internal* functions — they pin the port, and
  they re-validate every mechanical rename and the `CellBoundary` re-port.

### 9.2 New layers

- **Public wrapper tests** (pure Go): per wrapper — happy path, each documented error
  (`errors.Is`), hole-pruning behavior (pentagon disks/edges/vertexes), `Append*` reuse
  semantics (results appended after `len(dst)`; capacity reused; prior contents intact).
- **Allocation assertions** (pure Go, table-driven `testing.AllocsPerRun`): every
  `Append*` API asserts **0 allocs** on the warm path and ≤ 1 (the result) on the cold
  path; `Boundary`, `String`, `MarshalText`, iterators assert their documented budgets.
  These are regular tests, not benchmarks, so CI enforces them; run with `-race` excluded
  from the count assertions where the race runtime allocates (guard with
  `testing.Short()`-independent build check on race via `race` build tag file).
- **Benchmarks**: public wrapper vs direct internal call for the hot set (`GridDisk`,
  `Children`, `PolygonToCells`, `Boundary`, `CompactCells`), using `testing.B.Loop`
  (Go 1.24). Purpose: keep the wrapper tax at "size-computation + clear" (~15–20% at k=2,
  amortizing to ~0 at larger k — Appendix A, E2), and catch escape-analysis regressions.
- **Property/fuzz tests** (pure Go): `FuzzParseCell` (parse↔String round-trip, mode
  rejection), `FuzzLatLngToCell` (res/domain errors only, never panics; result always
  `IsValid` on success), round-trips `Cell↔(lat,lng)@same-res`, `Cell↔CoordIJ↔Cell`,
  `compact(uncompact(x)) == x` on valid sets, `Parent(Children(c)) == c`.
- **Golden API surface test**: a test renders the exported surface (via
  `golang.org/x/tools/go/packages` in a **separate test-only module** under `apitest/`, or
  more simply `go doc -all` normalization in a script) and diffs against a committed
  `docs/api-surface.txt`. Any accidental (un)export fails CI. Root module stays
  dependency-free.
- **No-unsafe gate** (`make check-unsafe`, introduced in Phase 0): two independent
  layers, both cheap, both pure shell + `go list`:

  ```sh
  # Layer 1 — build-selection check (authoritative for the CI platform).
  # For every normal build mode, ask the toolchain which packages/files it
  # would actually compile — including test files — and fail if anything in
  # THIS module imports unsafe. (Stdlib internals are out of scope by design.)
  for cgo in 0 1; do
    for tags in "" "race"; do
      CGO_ENABLED=$cgo go list -tags "$tags" \
        -f '{{.ImportPath}}: {{join .Imports " "}} {{join .TestImports " "}} {{join .XTestImports " "}}' \
        ./... | grep -w unsafe && exit 1
    done
  done

  # Layer 2 — tag allowlist (platform- and GOOS-independent).
  # Any file in the repository importing unsafe MUST carry a build constraint
  # containing c2go; this also catches files layer 1 can never select on the
  # CI platform (e.g. a hypothetical //go:build windows file).
  for f in $(grep -rl '^[[:space:]]*"unsafe"$' --include='*.go' . | grep -v '^./testref'); do
    grep -q 'go:build.*c2go' "$f" || { echo "$f: unsafe outside c2go"; exit 1; }
  done
  ```

  A third, reviewer-facing layer: a `depguard` rule in `.golangci.yml` denying `unsafe`
  (after the Phase-0 retag, lint's default build tags no longer select the parity files,
  so the rule needs no exceptions). The three layers fail independently, so bypassing the
  policy accidentally requires three separate mistakes.
- **C-API completeness gate**: `tools/apiinventory` in verify mode (Phase 6) fails CI if
  any `H3_EXPORT` in the vendored header lacks both a Go port attribution and a public
  wrapper/omission entry. This is the "detect missing upstream functions" test and runs
  pure-Go (it only reads files; `testref/` is fetched in the CI job that needs it).
- **Differential vs uber/h3-go** (optional, cgo, separate module `interop/uberdiff/` with
  its own `go.mod`): randomized cross-checks of every overlapping operation against the
  official binding (which wraps C 4.5.0 — version-dependent expectations are pinned to the
  overlap). Run by a dedicated make target/CI job, never by default. When upstream `x/h3go`
  releases, it can join the same harness cgo-free.

### 9.3 CI matrix (target state)

| Job | Command | Needs |
|---|---|---|
| build+test | `CGO_ENABLED=0 go test ./...` | nothing |
| race | `go test -race ./...` (interop excluded by `c2go` tag) | cgo toolchain |
| lint | golangci-lint (incl. depguard no-unsafe), `gofmt -s`, smrcptr | nothing |
| no-unsafe | `make check-unsafe` (both layers, all four build modes) | nothing |
| allocs | included in build+test (plain tests) | nothing |
| parity | `make -C testref h3-source && make test-c2go` | network + C toolchain |
| api-surface + completeness | golden diff + apiinventory verify | testref fetch |
| uberdiff | `make test-uberdiff` | cgo + network |

---

## 10. Upstream synchronization workflow

Goal: when H3 `4.4.x`/`4.5.x` releases, produce a reviewable, function-scoped update.

1. **Fetch & diff**: `make -C testref H3VER=4.4.0` (download alongside 4.3.0), then
   `git diff --no-index testref/h3-4.3.0/src/h3lib testref/h3-4.4.0/src/h3lib` — the diff
   is naturally per-C-function because upstream code is function-organized.
2. **Map changed C functions → Go files**: for each changed C function name, grep
   `Ported from H3 C: <file>::<name>` (or consult `docs/c-api-inventory.csv`). The
   one-function-per-file layout makes each port update a self-contained diff.
3. **Detect added/removed API**: rerun
   `go run ./tools/apiinventory -h3ver 4.4.0 > docs/c-api-inventory.csv`; new
   `H3_EXPORT`s appear as `MISSING`, removed ones as stale attributions. The CI
   completeness gate then *requires* either a port + wrapper or an omissions entry.
4. **Port function-by-function**: update the Go body next to the C diff; the parity
   harness retargets with `make test-c2go H3VER=4.4.0` — no code changes needed for the
   harness itself (include paths are injected at test time by design).
5. **Update tests**: port upstream's new/changed `testXxx.c` cases (tracked in
   `docs/ported-c-tests.md` as today); extend parity tests for new functions.
6. **Preserve intentional Go deviations**: they are enumerable — `Angle` fields,
   `int32` mapping, fixed-array `CellBoundary`, dst-slice out-params, iterator structs.
   Each is documented at its declaration and listed in this file; a `docs/DEVIATIONS.md`
   one-pager (Phase 6) is the checklist consulted during ports.
7. **Version bump**: update the `h3Version*` constants and the default `-h3ver`; tag.

Automated: fetching, diffing, inventory, completeness gate, parity retarget.
Deliberately manual: the porting itself and wrapper design for new functions — automating
semantic translation is where silent behavioral drift would enter.

---

## 11. Recommended implementation phases

Each phase is a small, independently green, reviewable unit; `make lint`, `make test`
(CGO_ENABLED=0), `go test -race`, and `make test-c2go` must pass at every phase boundary.

**Phase 0 — Repo hygiene (unblocks everything)**
Scope: `make fix-fmt` the 8 unformatted files at HEAD (the current first CI failure);
unify interop/parity build tags to `cgo && c2go` (19 `*_cgo.go`, 224 parity files,
2 cgo-tagged plain tests, `h3lib_vertexGraph_cgo.c`); verify CI goes green; delete
`h3.go.backup`/`h3_test.go.backup` (superseded by this document; preserved in git
history); add the `make check-unsafe` gate and the `depguard` no-unsafe lint rule
(DR-007, §9.2) to CI; rewrite `README.md`/`AGENTS.md`/`CLAUDE.md` to describe the real
layout; bump `go.mod` to `go 1.24` and CI to an explicit matrix (oldest supported +
latest).
Risks: none meaningful (tag edits are mechanical).
Done when: CI green on all jobs, including the no-unsafe gate; gopls clean with cgo
enabled.

**Phase 1 — Type foundation**
Scope: `h3api_types.go` gains `Cell`/`DirectedEdge`/`Vertex` + `h3Index = Cell` alias
(keeping the exported `H3Index = Cell` spelling temporarily so this commit touches one
file); `errors.go` with 15 sentinels + `toErr`; curated public constants.
Tests: alias build + full suites (already proven, Appendix A E1); `errors.Is` table test.
Done when: both suites green with the alias; no public wrappers yet.

**Phase 2 — Unexport sweep (mechanical)**
Scope: scripted word-boundary renames per §7.1 (~30 types, ~121 consts, 16 tables,
`H3Index`→`h3Index`, `H3Error`→`h3Error`, `E_*`→`e*`, `FLAG_*`); handle the
`DIRECTIONS`→`algosDirections` collision; script committed under `tools/unexport/` for
review, then deleted or kept as documentation.
Tests: the point of this phase is that the *entire existing suite including 224 parity
files* re-validates it. Golden API surface file introduced here (locks the now-minimal
surface).
Risks: rename collisions → script detects identifier collisions before applying; lint
churn (revive on unusual lowerCamel like `mPi`) → reviewed `.golangci.yml` exclusions if
needed.
Done when: `go doc ./...` shows only `Angle` family + new types/constants/errors; suites
green.

**Phase 3 — Public API wave 1: indexing, inspection, hierarchy**
Scope: `cell.go`, `latlng.go`, `boundary.go`, `hierarchy.go`, `errors.go` finalization;
`CellBoundary` value-type re-port (fixed `[10]LatLng`, updates
`faceijk__faceIjk{,Pent}ToCellBoundary.go`, `h3index__cellToBoundary.go`,
`directedEdge_directedEdgeToBoundary.go`, `polygon__cellBoundaryCrossesGeoLoop.go`);
String/Parse/MarshalText via the small generics (§6).
Tests: wrapper tests + allocation assertions + fuzz (Parse, LatLngToCell) + re-pointed
upstream vectors (`testH3Api`, `testCellToParent`, …); parity suite re-validates the
boundary re-port; benchmark `Boundary` (expect 6→0 allocs).
Risks: boundary re-port touches real algorithm files — mitigated by existing parity tests
comparing values against C.
Done when: wave API documented + tested + 0-alloc assertions green.

**Phase 4 — Wave 2: traversal, edges, vertexes**
Scope: `traversal.go`, `edge.go`, `vertex.go`, `localij.go`; `Append*` forms with
window+clear+compact pattern (§5.2); hole-pruning policy per §2.1 inventory.
Tests: pentagon-centric pruning tests (the classic H3 traps), allocation assertions,
`GridDisk`/`GridPath` benchmarks, fuzz round-trips (`CoordIJ`, `GridDistance`
symmetry).
Done when: same bar as Phase 3.

**Phase 5 — Wave 3: regions, compact, metrics, iterators**
Scope: `region.go` (`PolygonToCells*`, `CellsToMultiPolygon`), `compact.go`,
`metrics.go`, `iter.go` (`ChildrenSeq`, `CellsAtRes`, polygon-iterator shape per
§12-Q6).
Tests: upstream polygon vectors re-pointed; multipolygon conversion tests
(donuts/nested/transmeridian); alloc assertions incl. iterator 0-alloc; benchmarks
(`PolygonToCells`, `CompactCells`).
Done when: **every** C-public function has wrapper or omissions entry (gate passes);
tag `v0.1.0`.

**Phase 6 — Docs, examples, gates**
Scope: `doc.go` package docs; `Example*` tests (indexing, polyfill, disk, edge walk);
README rewrite with allocation guidance; `docs/DEVIATIONS.md`; apiinventory verify mode +
`H3 C API:` doc-line enforcement; CI matrix per §9.3; `interop/uberdiff` differential
module (optional job).
Done when: pkg.go.dev-ready docs; all CI gates on.

**Phase 7 — First upstream sync rehearsal (validates the whole design)**
Scope: dry-run the §10 workflow against H3 4.4.0 (or newest): fetch, diff, inventory,
port the delta, extend tests.
Done when: sync lands with function-scoped commits; decide on `v1.0.0`.

---

## 12. Open questions and challenged assumptions

Each with a recommendation; items marked ⚠ need your explicit sign-off because they set
user-visible policy.

**Q1 — "Aliases can give type safety." Partly false, and that's fine.** `h3Index = Cell`
provides zero *internal* safety (edges held as `Cell`-typed values internally, §4.2) — it
provides *boundary* safety plus zero-copy. The alternative (three defined types
everywhere internally) would require rewriting ported signatures and bodies, destroying
traceability. Resolution: accept boundary-level typing; it is strictly stronger than both
C and the current state.

**Q2 — "Slice conversions between differently-named element types work." False by spec** —
`[]Cell` ↔ `[]uint64` is not a permitted conversion (only identical element types
convert). This is exactly why the alias (identity, not conversion) is the design keystone,
and why Option B (internal package) collapses into copies or `unsafe`.

**Q3 — "Move everything to `internal/`." Rejected with evidence** (§3.1): methods can't be
declared on cross-package types, alias-to-internal hides documentation, and the parity
tests are white-box. The single package *is* the traceability-preserving choice, not a
compromise.

**Q4 ⚠ — `Angle` in `LatLng` vs `float64` degrees.** Keeping `Angle` is zero-cost and
bug-resistant but diverges from uber (`float64` degrees) and requires
`h3.LatLng{Lat: h3.Deg(37.77), …}` at construction. Recommendation: **keep `Angle`**, add
`LatLngDegs(lat, lng float64)` for ergonomics. Flip side documented: if you would rather
have drop-in uber familiarity, this is the single decision to override *now* — it becomes
expensive after Phase 3 (it re-types the whole geometry surface and forcibly introduces
convert-copies at every polygon/boundary boundary, §5.4).

**Q5 ⚠ — JSON/text marshaling for geometry.** Deferred (no marshalers on
`LatLng`/`Angle`/`GeoPolygon` in v0) because degrees-vs-radians and GeoJSON's
`[lng, lat]` ordering make any default a silent trap. Index types do get
`MarshalText` (hex). Recommendation: revisit as a `geojson` helper package post-v0.1.

**Q6 — Iterator error shape for polygon iteration.** `iter.Seq2[Cell, error]` pollutes
every loop iteration with an error that can occur at most once; a handle
(`it := IterPolygonCells(...); for c := range it.All() { … }; it.Err()`) is uglier but
honest. Recommendation: validate inputs *before* returning the Seq (errors up front,
sequence itself infallible) — matches how the C iterator actually fails (init-time), keeps
`iter.Seq[Cell]`. Verify against `iterErrorPolygonCompact` semantics during Phase 5.

**Q7 — Hiding `H3Error`/`E_*` codes.** Sentinels cover Go ergonomics, but numeric codes
are cross-language stable. Recommendation: hide now, add `func Code(error) (int, bool)`
later if requested — additive, no redesign.

**Q8 ⚠ — `GridDiskDistances` distance type: `[]int32` vs `[]int`.** The algorithm is
`int32` by the port's C-overflow rule; exposing `[]int` would force a widen-copy per call
even on the warm path. Recommendation: `[]int32` (documented rationale); revisit only if
it proves abrasive. (uber sidesteps with `[][]Cell` buckets — that allocates k+1 slices
and loses the flat shape; can be added later as `GridDiskDistancesGrouped` if wanted.)

**Q9 — "Complete C parity would make a poor Go API." True in places** — resolved by the
documented omissions list (§2.2): out-param style, `destroy*`, `describeH3Error`,
`degsToRads`, raw iterator structs, `GeoMultiPolygon` are all *represented differently*,
not missing. The completeness gate checks representation, not name-for-name parity.

**Q10 — Methods vs package functions coexisting.** Policy instead of uber's ad-hoc mix:
operations with one natural receiver are methods; multi-subject or constructor-like
operations are package functions (`LatLngToCell`, `CompactCells`, `GridDisksUnsafe`,
`PolygonToCells`, `CellToLocalIJ`). No dual forms — one obvious way each.

**Q11 — README's "deterministic ordering" promise.** The C contracts are weaker
(gridDisk: "no particular order"). The port must not silently sort (extra cost + parity
divergence). Resolution: document *actual* per-function order (C's), fix the README;
ordering guarantees beyond C's are out of scope.

**Q12 — Naming-convention debt in ported files** (49/75 public functions in `__` files,
`h3Index_getBaseCellNumber.go` casing, attribution format variants). Resolution: do *not*
mass-rename files (churn without benefit — the inventory tool already normalizes);
optionally fix the single casing outlier and standardize future attributions
(`file.c::name`, no `H3_EXPORT()` wrapper) in `CONTRIBUTING.md`.

**Q13 — Uber's `x/h3go` pure-Go port exists** (in uber/h3-go master, unreleased). It
covers indexing+hierarchy+sets so far, `int64`/degrees shapes. Impact: validates the
pure-Go direction; raises the bar on timing. Differentiators to protect: full API
coverage, zero-alloc `Append*`/Seq forms, function-level upstream traceability.

---

## 13. Decision records

**DR-001 — Package layout: single root package, two file layers.**
Context: idiomatic public API over a 274-file mechanically ported C layer; methods, godoc,
zero-copy, parity tests all constrain placement. Options: single package; `internal`
impl package; per-C-module internal packages; generated layer. Decision: single package.
Consequences: zero-copy via alias possible; parity tests stay white-box; export hygiene
must come from the unexport sweep rather than package boundaries. Rejected: `internal`
split (methods/type identity/doc visibility, §3.1-B/C), generation (§3.1-D).

**DR-002 — Public index types: `Cell`, `DirectedEdge`, `Vertex` as defined `uint64`
types; no public `Index`.**
Context: C uses one `H3Index` for three modes; Go users benefit from mode-typed APIs.
Decision: three defined types + unexported `index` constraint for shared helpers; parsing
is per-type and mode-validated. Consequences: compile-time mode safety at the API
boundary; ≤6-element conversions at edge/vertex boundaries; raw conversions remain
possible and documented as unchecked. Rejected: public umbrella `Index` type (forces
conversions, adds an invalid-mode-shaped hole); uber's `int64` representation (sign traps
in mask arithmetic).

**DR-003 — Alias strategy: internal `h3Index = Cell`; `Angle`-based geometry types shared
between layers.**
Context: §5's zero-copy requirement; `[]T` conversion rules. Decision: alias internal name
to the public defined type; keep `Angle`/`LatLng`/`GeoLoop`/`GeoPolygon` as the single
representation both layers use. Consequences: type identity for all hot paths; ported
sources keep compiling verbatim (E1); degrees users construct via `Deg`/`LatLngDegs`.
Rejected: defined `H3Index` + conversions (illegal for slices, §12-Q2); public degrees
types with internal radians (O(n) convert-copies on every polygon/boundary, §5.4).

**DR-004 — Zero-copy slice handling: dst-window + `clear` + in-place hole compaction;
`Append*` naming; fixed-size stack arrays for ≤6 outputs.**
Context: C fills caller-sized buffers, sometimes as hash sets, sometimes with holes.
Decision: per §5.2/§5.3, with sizing functions exposed. Consequences: warm-path 0 allocs
(E2); a documented `clear` cost; pruning policy fixed per function from the C contract
inventory. Rejected: returning hole-y buffers (uber prunes too; holes are a C memory
idiom, not information); pooling/scratch options (premature).

**DR-005 — Generics: only unexported parse/format/marshal helpers (+ optional prune
helper); no generic algorithms; delete `castSlice`.**
Context: §6 table. Decision as stated. Consequences: 3× duplication avoided where it is
error-prone; zero effect on traceability or allocation. Rejected: generifying ported
signatures; public generic API in v0.

**DR-006 — Compatibility policy: one breaking sweep now; `v0.1.0` after Phase 5; `v1.0.0`
gated on a successful upstream sync rehearsal.**
Context: §8 evidence (no tags/users/callable API; red CI). Consequences: no deprecation
layer to build or maintain; freedom to fix `CellBoundary` and the export surface.
Rejected: compatibility wrappers, gradual deprecation (cost without beneficiaries).

**DR-007 — `unsafe`: prohibited in the production package (hard invariant).**
Context: the `.backup` draft's `castSlice` was the only production-`unsafe` candidate;
the alias eliminates its purpose. A full repository audit (Appendix B) confirms `unsafe`
appears today only in six cgo parity-interop files and the dead `.backup` file.
Decision — the invariant, not a preference:

- The production package `h3` is **safe Go only**: no file selected by any normal build —
  `go build ./...`, `go test ./...`, `CGO_ENABLED=0 go test ./...`,
  `go test -race ./...`, on any GOOS/GOARCH — imports or uses `unsafe`.
- No public API requires callers to perform unsafe conversions to reach zero-copy or
  allocation-efficient functionality, exposes an unsafe abstraction, or depends on a
  representation contract that only `unsafe` could honor. Zero-copy `[]Cell`
  interoperability is achieved through **type identity** (`type h3Index = Cell` alias),
  not memory reinterpretation — verified in Appendix A (E1/E2) and per-category in
  Appendix B.3.
- `unsafe` is permitted **only** in test/parity/benchmark/development files that are
  excluded from every normal build by an explicit build constraint containing `c2go`
  (currently: the cgo parity interop, which needs `C.free` and C-array↔Go-slice
  bridging).
- Enforcement is mechanical, three independent layers: the `make check-unsafe` gate
  (build-selection check across all four build modes + platform-independent tag
  allowlist, §9.2) and a `depguard` lint rule.
- **Introducing `unsafe` into production code in the future requires a new reviewed
  decision record** containing: proof that no safe design can meet the requirement,
  benchmarks quantifying the win, the exact representation invariants relied upon and
  their basis in the Go specification, and the tests that pin the behavior. Absent that
  record, the CI gate is the answer.

Rejected: keeping `castSlice` "just in case" (dead unsafe code is pure risk); the softer
"avoid unsafe" phrasing (unenforceable).

---

## Appendix A — Experimental evidence

All experiments were run on this repository (worktree of commit `52d76be`), Go 1.25.4,
darwin/arm64. The experiment code is not committed; each is reproducible as described.

**E1 — Alias feasibility (keystone).** Change `type H3Index uint64` to
`type Cell uint64; type H3Index = Cell` in `h3api_types.go` — no other edits anywhere.
Results: `CGO_ENABLED=0 go build ./...` and `go vet ./...` clean; full pure-Go suite
passes (`ok … 25.5s`); full cgo parity suite passes (`make test-c2go`: `ok … 25.1s`).
Conclusion: the public defined type can be introduced with a one-file change and no
behavioral drift.

**E2 — Dst-buffer wrapper cost.** `AppendGridDisk`-shaped wrapper (size → grow → `clear`
→ `gridDisk` → return) vs direct internal call, k=2 (19 cells):
`testing.AllocsPerRun`: cold (nil dst) **1.0 alloc** (the result), warm (reused buffer)
**0.0 allocs**. Benchmarks: wrapper 384 ns/op, direct internal (with the mandatory
buffer clear the C contract requires anyway) 325 ns/op — the ~58 ns delta is the
`maxGridDiskSize` call + slice bookkeeping, amortizing toward zero for larger k.

**E3 — Wrapper inlining.** `go test -gcflags='-m'`: thin wrappers (`Resolution`-style
scalar wrapper, `toErr`-style table lookup) report `can inline`; `getResolution` inlines
into its ~10 internal callers today, confirming the pattern.

**E4 — `CellBoundary` allocation hotspot.** `cellToBoundary` on a res-15 hexagon:
fresh `CellBoundary`: **6.0 allocs/call** (the per-vertex grow loop in
`faceijk__faceIjkToCellBoundary.go`); with a caller-preallocated len-10 `Verts` that the
escape analyzer keeps on the stack: **0.0 allocs**. C performs 0 heap allocations
(fixed embedded array). Basis for the §4.4 value-type re-port.

**E5 — `iter.Seq` children iterator.** Range-over-func wrapper over
`iterInitParent`/`iterStepChild`, res 5 → res 8 (343 children):
**0.0 allocs/run** for construction + full iteration (Go 1.23+ required; go.mod bumped to
1.23 for the experiment only).

**E6 — CI failure reproduction.** `CGO_ENABLED=1 go test -race -run xxx ./...` fails at
build with `./mathExtensions_cgo.go:11:10: fatal error: 'mathExtensions.c' file not
found` — the `//go:build cgo` (without `c2go`) tagging pulls parity interop into any
cgo-enabled build; matches the three consecutive red CI runs. Basis for Phase 0.

**Spec references.** Slice conversion legality: Go spec §Conversions — a slice is
convertible to another slice type only when the element types are *identical* (or
identical under struct-tag ignoring); named element types with same underlying type are
not identical, hence `[]Cell` ↛ `[]uint64`. Alias declarations (§Type declarations) create
*the same type*, hence `[]h3Index` ≡ `[]Cell` — the design's foundation.

---

## Appendix B — `unsafe` audit

Audit date 2026-07-11, commit `ee2d510` + this document. Method: `grep -rl '"unsafe"'`
over every `*.go` and `*.go.backup` in the repository (excluding `testref/`, which
contains only upstream C sources), plus a sweep for unsafe-adjacent mechanisms
(`go:linkname`, `reflect.SliceHeader`/`StringHeader`: **none found** outside one benign
`reflect.ValueOf` nil-check in the dead `h3_test.go.backup`), plus `go list`-verified
build-selection checks per mode.

### B.1 Every current occurrence

| File | Kind | Build constraint | Selected by a normal build? | Why `unsafe` exists | Fate |
|---|---|---|---|---|---|
| `h3index_cgo.go` | cgo parity interop | `//go:build cgo` → **`cgo && c2go` in Phase 0** | Today: *selected* whenever `CGO_ENABLED=1` (the build then fails on missing C includes — nothing links against it in practice); after Phase 0: **no** | `C.free(unsafe.Pointer(...))` on C-malloc'd buffers; `(*[1<<30]C.H3Index)(unsafe.Pointer(p))[:n:n]` views of C arrays for Go↔C comparison | keep, retag |
| `linkedGeo_cgo.go` | cgo parity interop | same | same | freeing/marshaling the C linked-list polygon structures | keep, retag |
| `polygon_cgo.go` | cgo parity interop | same | same | C struct marshaling | keep, retag |
| `utility_cgo.go` | cgo parity interop | same | same | C string/buffer interop | keep, retag |
| `vertexGraph_cgo.go` | cgo parity interop | same | same | C vertex-graph node interop | keep, retag |
| `polyfill_cgo.go` | cgo parity interop | `//go:build cgo && c2go` (already correct) | **no** (verified: absent from `go list` CgoFiles without `-tags c2go`, present with it) | C polygon interop | keep as-is |
| `h3.go.backup` | obsolete draft | none — `.backup` is not a Go file; `go list` confirms it is never compiled | **no** | the abandoned `castSlice[From, To ~uint64]` / `asH3Array6` slice-reinterpretation helpers | delete in Phase 0 |

For completeness of the classification requested: the remaining 13 `*_cgo.go` interop
files do **not** import `unsafe`; **zero** of the 274 pure-Go implementation files, 224
parity tests, 66 plain test files, or `tools/apiinventory` import it; there is no
generated code and there are no benchmarks yet (planned ones are pure Go, §9.2).

Transitive user exposure: none, in any state. The interop files define only unexported
`…C` wrappers used by parity tests; no exported symbol references them, and after the
Phase-0 retag they are not even compiled unless the user explicitly opts in with
`-tags c2go` *and* cgo *and* the `CGO_CPPFLAGS` include paths from `make test-c2go`. A
library importer cannot transitively depend on them.

### B.2 Build-mode verification (measured, not asserted)

| Command | Files importing `unsafe` selected today | After Phase-0 retag |
|---|---|---|
| `CGO_ENABLED=0 go test ./...` | 0 (verified: `go list` CgoFiles is empty) | 0 |
| `go build ./...` (cgo default-on) | 5 (build fails anyway on missing includes) | **0** |
| `go test ./...` (cgo default-on) | 5 (same) | **0** |
| `go test -race ./...` | 5 (same — this is the red CI job) | **0** |

The `cgo && c2go` exclusion mechanism is already proven in-repo: `polyfill_cgo.go`
carries it today and is verifiably absent from every default build's file list.

### B.3 Planned public API: per-category safety verification

Every category of the proposed surface, checked against the actual planned
implementation (the Appendix A experiment code compiled with **no** `unsafe` import):

| Category | Mechanism | `unsafe`? |
|---|---|---|
| Scalar conversions (`Cell(x)`, `DirectedEdge(v)`, `int(res)`) | language conversions between types with identical underlying types | no |
| `[]Cell` inputs/outputs incl. all `Append*` APIs | **type identity** — `[]h3Index` *is* `[]Cell` (alias); `slices.Grow` + `clear` + in-place compaction (E2) | no |
| Polygon/geometry inputs (`GeoPolygon`, `GeoLoop`, `LatLng`) | public types *are* the internal types; passed by pointer | no |
| `CellBoundary` | value struct with fixed `[10]LatLng` array; plain indexing (E4) | no |
| Directed-edge / vertex outputs | stack `[6]h3Index` + element-wise converting loop (§5.3) | no |
| Iterators (`ChildrenSeq`, `CellsAtRes`) | range-over-func closures over ported iterator structs (E5) | no |
| Text parsing/formatting/marshaling | `strconv.AppendUint`/manual hex over stack buffers; unexported `~uint64` generics | no |
| Multipolygon conversion (`CellsToMultiPolygon`) | pointer-walk of the linked structure + exact-size slice fill | no |

No production feature in the proposed design requires `unsafe`; consequently no safe
alternative with a performance penalty needs to be traded off anywhere. The one historical
candidate — bridging `[]Cell` to a differently-defined index type — is dissolved by the
alias rather than worked around.

