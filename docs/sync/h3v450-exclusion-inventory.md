# `!h3v450` exclusion inventory (temporary, lives until the I-I cutover)

The migration harness (issue #27, record
[4.4.0-to-4.5.0.md](4.4.0-to-4.5.0.md) §15.1) supports both upstream
trees via an `H3VER`-derived build tag: `make test-c2go H3VER=4.5.0`
passes `-tags c2go,h3v450`, and files valid for only one tree shape
carry `h3v450`/`!h3v450` constraints. **This inventory lists every file
and every test function excluded from the 4.5.0 configuration**, each
mapped to the implementation issue that replaces or retires it — so a
green `make test-c2go H3VER=4.5.0` run cannot silently hide unaccounted
skips. Each issue that removes entries updates this file; it must be
**empty (and deleted) at the I-I cutover** (#36).

## Excluded from the 4.5.0 configuration (`!h3v450`)

| File | Excluded test functions | Owner |
|---|---|---|
| `h3lib_vec3d_c2go.c` (compiles 4.4.0's `vec3d.c`, which faceijk.c's `_geoToClosestFace` links against — needed by the 4.4.0 C library itself even though the Go vec3d wrappers retired with I-A) | — | I-I #36 (deleted at cutover) |

No Go files and no test functions are excluded from the 4.5.0
configuration any more: I-C #34 retired the entire vertexGraph parity
domain together with its subject (upstream deleted `vertexGraph.c/h`
and the Go port deleted the corresponding `vertexGraph*.go` /
`algos_h3SetToVertexGraph*` / `algos__vertexGraphToLinkedGeo*` files,
their parity tests, and the `h3lib_vertexGraph_cgo.c` shim), and
flipped `algos_cellsToLinkedMultiPolygon_parity_test.go` to the
4.5.0-only configuration (the Go library now implements the 4.5.0
arc-based algorithm).

## Adaptive (not excluded — compiles differently per tree)

- `h3lib_coordijk_c2go.c`: `__has_include("coordijk.c")` — compiles the
  4.4.0 implementation; empty translation unit at 4.5.0 (the `static
  inline` header definitions serve every including TU). No test skips.
- `h3lib_faceijk_c2go.c` / `h3lib_localij_c2go.c`: version-neutral
  includes plus `#if __has_include("area.c")`-guarded test-only
  wrappers (`h3goTest_*`) that expose the 4.5.0 file-static Vec3
  pipeline helpers and `gridPathCellsInterpolate` to the parity
  harness in their own translation units. No 4.4.0 impact.
- `parity_float_helpers_test.go`: ungated shared helpers (both
  configurations).

## Present only in the 4.5.0 configuration (`h3v450`)

- `h3lib_area_c2go.c` (shim for the new `area.c`)
- `h3lib_cellsToMultiPoly_c2go.c` (shim for the new
  `cellsToMultiPoly.c` + same-TU `h3goTest_*` wrappers for its
  file-statics)

- `vec3d_cgo.go`, `h3index_vec3_cgo.go`, `faceijk_vec3_cgo.go` +
  `vec3d_vec3Ops_parity_test.go`, `h3index_cellToVec3_parity_test.go`,
  `faceijk_vec3Pipeline_parity_test.go` (I-A #29)
- `area_cgo.go` + `area_geoLoopAreaRads2_parity_test.go` (geoLoop/
  geoPolygon/geoMultiPolygon/cellAreaRads2/kadd/cagnoli parity; I-B #30)
- `localij_interpolate_cgo.go` +
  `localij_gridPathCellsInterpolate_parity_test.go` (interpolate parity
  and the exact 4.5.0 gridPathCells pairs; I-D #31)
- `directedEdge_reverse_cgo.go` +
  `directedEdge_reverseDirectedEdge_parity_test.go` (I-E #32)
- `cellsToMultiPoly_cgo.go` + `cellsToMultiPoly_parity_test.go` and
  the flipped `algos_cellsToLinkedMultiPolygon_parity_test.go`
  (arc-based multipolygon machinery, linkedGeo conversions, and the
  rewritten cellsToLinkedMultiPolygon; I-C #34)

Every version-specific parity domain planned by the record is now in
place; the remaining `!h3v450` entry above is harness-plumbing only and
is deleted at the I-I cutover.

**Resolved by I-C #34** (deleted with the replaced implementation): the
whole vertexGraph parity domain — `h3lib_vertexGraph_cgo.c`,
`vertexGraph_cgo.go`, `algos_vertexgraph_cgo.go`, the eleven
`vertexGraph*_parity_test.go` files,
`algos_h3SetToVertexGraph_parity_test.go`, and
`algos__vertexGraphToLinkedGeo_parity_test.go` — and the behavior-level
exclusion of `algos_cellsToLinkedMultiPolygon_parity_test.go` (now
h3v450-gated instead).

**Resolved by I-B #30** (deleted with the replaced implementation): the
triangle-area wrappers and parity tests (`latLng_h3v44_cgo.go`,
`latLng__triangleArea_parity_test.go`,
`latLng__triangleEdgeLengthsToArea_parity_test.go`) and the gated
4.4.0-era `latLng__cellAreaRads2_parity_test.go` (replaced by the
h3v450 area parity above).

**Resolved by I-A #29** (files deleted with their replaced pipeline, so
they are no longer exclusions): the old vec3d Go harness (the 4.4.0
`vec3d_cgo.go`, three vec3d parity tests, two vec3d pure tests),
`faceijk_h3v44_cgo.go` with the four old-projection parity tests, and
the two latLng azimuth parity tests. The `h3lib_vec3d_c2go.c` C shim
stays (row above): the 4.4.0 C library needs the translation unit
regardless of the Go side.

## Still exercised at 4.5.0 despite upstream changes

`latLng__cellAreaKm2_parity_test.go` / `latLng__cellAreaM2_parity_test.go`
pass in both configurations (their absolute tolerances absorb the
~1e-15-relative area change) and deliberately stay ungated for coverage.
`Test_latLngToCell/cellToLatLng/cellToBoundary` and the exhaustive
parity suites pass unchanged at 4.5.0 within their documented
comparison disciplines — cell indexes and error codes exactly,
cellToLatLng coordinates within the suite's 1e-10 rad tolerance,
boundary vertices within 1e-12 rad, and the upstream fixtures within
their 1e-6-degree threshold (latLngToCell fixture indexes exactly). At
those tolerances the Vec3 refactor produced no observable difference
(record §14 risk 1 evidence); bit-identity of the continuous outputs is
not claimed.

Note on the 4.4.0 configuration after I-C: the Go library implements
the 4.5.0 multipolygon algorithm, so the multipolygon parity domain is
h3v450-only by necessity (the 4.4.0 C oracle implements the retired
vertexGraph algorithm). `make test-c2go H3VER=4.4.0` still runs every
other suite, including the coordijk/faceijk/latLng domains shared by
both trees.
