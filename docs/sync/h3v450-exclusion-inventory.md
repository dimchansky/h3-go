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

### Compile-level: wrappers/shims binding symbols that do not exist at 4.5.0

| File | Excluded test functions | Owner |
|---|---|---|
| `latLng_h3v44_cgo.go` (wrappers: `triangleEdgeLengthsToAreaC`, `triangleAreaC`; the azimuth wrappers retired with I-A) | — | I-B #30 |
| `h3lib_vec3d_c2go.c` (compiles 4.4.0's `vec3d.c`, which faceijk.c's `_geoToClosestFace` links against — needed by the 4.4.0 C library itself even though the Go vec3d wrappers retired with I-A) | — | I-I #36 (deleted at cutover) |
| `latLng__triangleArea_parity_test.go` | `Test_triangleArea_ParityWithC` | I-B #30 |
| `latLng__triangleEdgeLengthsToArea_parity_test.go` | `Test_triangleEdgeLengthsToArea_ParityWithC` | I-B #30 |
| `vertexGraph_cgo.go` | — | I-C #34 (whole domain retired) |
| `h3lib_vertexGraph_cgo.c` (shim for the deleted `vertexGraph.c`) | — | I-C #34 |
| `algos_vertexgraph_cgo.go` (wrappers: `_vertexGraphToLinkedGeoC`, `h3SetToVertexGraphC`, `h3SetToVertexGraphCForParity`) | — | I-C #34 |
| `vertexGraph__hashVertex_parity_test.go` | `Test_hashVertex_parity` | I-C #34 |
| `vertexGraph__initVertexGraph_parity_test.go` | `Test_initVertexGraph_parity` | I-C #34 |
| `vertexGraph__initVertexNode_parity_test.go` | `Test__initVertexNode_parity` | I-C #34 |
| `vertexGraph_addVertexNode_parity_test.go` | `Test_addVertexNode_parity`, `Test_addVertexNode_hash_collision_parity` | I-C #34 |
| `vertexGraph_destroyVertexGraph_parity_test.go` | `Test_destroyVertexGraph_parity` | I-C #34 |
| `vertexGraph_findNodeForEdge_parity_test.go` | `Test_findNodeForEdge_parity`, `Test_findNodeForEdge_linked_list_traversal_parity` | I-C #34 |
| `vertexGraph_findNodeForVertex_parity_test.go` | `Test_findNodeForVertex_parity`, `Test_findNodeForVertex_multiple_edges_same_vertex_parity` | I-C #34 |
| `vertexGraph_firstVertexNode_parity_test.go` | `Test_firstVertexNode_parity` | I-C #34 |
| `vertexGraph_removeVertexNode_parity_test.go` | `Test_removeVertexNode_parity`, `Test_removeVertexNode_not_found_parity`, `Test_removeVertexNode_empty_graph_parity` | I-C #34 |
| `algos_h3SetToVertexGraph_parity_test.go` | `Test_h3SetToVertexGraph_parity` | I-C #34 |
| `algos__vertexGraphToLinkedGeo_parity_test.go` | `Test_vertexGraphToLinkedGeo_parity` | I-C #34 |

### Behavior-level: parity valid at 4.4.0 but changed upstream at 4.5.0

| File | Excluded test functions | 4.5.0 divergence (record ref) | Owner |
|---|---|---|---|
| `algos_cellsToLinkedMultiPolygon_parity_test.go` | `Test_cellsToLinkedMultiPolygon_parity` | invalid-cell input: Go(4.4.0)=E_FAILED vs C(4.5.0)=E_CELL_INVALID (§7.1) | I-C #34 |
| `latLng__cellAreaRads2_parity_test.go` | `Test_cellAreaRads2_parity` | area algorithm change: 1.4e-15 diff at res 0 exceeds the test's tight tolerance (§7.4) | I-B #30 |

## Adaptive (not excluded — compiles differently per tree)

- `h3lib_coordijk_c2go.c`: `__has_include("coordijk.c")` — compiles the
  4.4.0 implementation; empty translation unit at 4.5.0 (the `static
  inline` header definitions serve every including TU). No test skips.

## Present only in the 4.5.0 configuration (`h3v450`)

- `h3lib_area_c2go.c` (shim for the new `area.c`)
- `h3lib_cellsToMultiPoly_c2go.c` (shim for the new `cellsToMultiPoly.c`)

- `vec3d_cgo.go`, `h3index_vec3_cgo.go` + `vec3d_vec3Ops_parity_h3v450_test.go`,
  `h3index_cellToVec3_parity_h3v450_test.go` (landed with I-A #29)

Version-specific wrappers and parity tests for the remaining new 4.5.0
symbols arrive with their owning issues (I-B/I-C/I-D/I-E) behind
`h3v450`.

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
~1e-15-relative area change) and deliberately stay ungated for coverage;
`Test_latLngToCell/cellToLatLng/cellToBoundary` and the exhaustive
parity suites pass unchanged at 4.5.0 on the tested inputs (the Vec3
refactor produced no observable difference there — record §14 risk 1
evidence).
