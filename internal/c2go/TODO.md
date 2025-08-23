Scope
- Track C→Go function conversions performed in `internal/c2go`.
- Record dependencies and testing status for each function.

Design rule (performance, Go‑friendly APIs)
- Use the dst‑buffer slice pattern for any function producing collections.
  - Prefer signatures like `func fn(dst []T, args ...) ([]T, H3Error)`.
  - Reuse capacity in `dst` when sufficient; allocate only when `cap(dst)` is too small.
  - Keep element ordering identical to the C reference.

Signature mirroring (C pointers → Go pointers)
- When the C function takes struct pointers (e.g., `const LatLng*`, `const BBox*`), mirror that in Go with pointer parameters (`*LatLng`, `*BBox`, `*Vec2d`, `*Vec3d`, etc.).
  - Rationale: preserves in/out semantics, avoids accidental copies, and keeps the port traceable to the C API.
  - Const semantics: if the C pointer is `const`, keep the Go pointer parameter but do not mutate the pointed value.

Iteration Workflow (standard operating procedure)
- Select target: pick a small, self-contained C function (or a tight cluster) from the planned list.
- Prepare interop: if the C function is not exported or uses structs, add minimal cgo helpers to call it directly (use C structs like C.LatLng/C.BBox, etc.); avoid splitting scalars.
- Transpile: implement a faithful Go port in `internal/c2go/<cfile>__<function>.go`, mirroring names/signatures (unexported) and behavior (NO build tag for pure Go).
- Document: follow established comment style with function description, technical details, and `// Ported from H3 C: <file>::<function>` attribution line.
- Parity test: add `internal/c2go/<cfile>__<function>_parity_test.go` with `//go:build cgo` that compares Go vs C output; use tight, justified tolerances for floats and safe bool handling.
- Sanity run: `make test-c2go` and ensure all parity tests pass.
- Format code: run `make fix-fmt` prior to committing.
- Update tracker: update this TODO to mark the function DONE and list the next planned items (do this BEFORE each commit).
- Commit: commit the minimal, focused changes with a message stating the ported function(s), parity, and TODO update.

Source Reference
- Local H3 C source code available at: `testref/h3-4.3.0/src/h3lib/lib/*.c`
- When porting C code, reference the local sources directly for accurate implementation
- Port as closely as possible to the C implementation, preserving algorithm details and edge cases

Build Tag Rules
- **Pure Go implementations** (`<cfile>__<function>.go`): NO build tags
- **CGO interop files** (`*_cgo.go`): Must have `//go:build cgo` as first line
- **Parity test files** (`*_parity_test.go`): Must have `//go:build cgo` as first line
- This ensures pure Go code is always available while CGO dependencies are isolated

Comment Style Requirements
- **Function description**: Brief summary of what the function does
- **Technical details**: Include important implementation notes, algorithm explanations, or edge case handling
- **Source attribution**: Must include `// Ported from H3 C: <file>::<function>` line
- **Example format**:
  ```go
  // _functionName calculates something important.
  // Additional technical details about the implementation.
  // Ported from H3 C: fileName.c::_functionName
  func _functionName(params) returnType {
  ```
- **Constants**: Should be declared in `<cfile>_constants.go` files, not inline

Interop helper note
- When performing Go→C struct conversions in `*_cgo.go`, prefer creating small helper functions to convert common types (e.g., `toCGeoLoop`, `toCGeoPolygon`, `toCBBox`). This avoids duplication when multiple functions need the same transformation and makes memory management (malloc/free) clearer and safer.

Checklist
- [x] Scaffold package and docs (README.md)
- [x] Add first dependency-free function: `mathExtensions.c::_ipow`
- [x] Port `latLng.c::_posAngleRads` (angle normalization) — parity test via `_posAngleRadsC`
 - [x] Port `latLng.c::constrainLng` and `constrainLat`
 - [x] Port `latLng.c::degsToRads` and `radsToDegs`
 - [x] Port `latLng.c::geoAlmostEqualThreshold` and `geoAlmostEqual`
 - [x] Port `latLng.c::H3_EXPORT(greatCircleDistanceRads/Km/M)`
 - [x] Port `latLng.c::_geoAzimuthRads` and `_geoAzDistanceRads`
 - [x] Port `latLng.c::normalizeLng`, `triangleEdgeLengthsToArea`, `triangleArea`, `_setGeoRads`, `setGeoDegs`
 - [x] Port `vec2d.c::_v2dMag`, `_v2dIntersect`, `_v2dAlmostEquals`
- [ ] Port additional small helpers from `mathExtensions.c` and `latLng.c`
- [ ] Identify next targets with minimal dependencies; add TODO chains in code
- [ ] Mirror select C tests (apps/testapps, fuzzers) where practical

Non-negotiable constraints (repo policy)
- Do not copy or vendor H3 `.c` or `.h` files into this repo. Only include original files via build-tagged shims (e.g., `h3lib_*.c`) and let `CGO_CPPFLAGS` provide include paths.
- Keep all c2go work scoped to `internal/c2go` only; do not introduce dependencies on other Go packages.

Functions
- mathExtensions.c
  - [x] `_ipow(int64_t base, int64_t exp)` — implemented; parity-tested via `_ipowC` wrapper
- latLng.c
  - [x] `_posAngleRads(double rads)` — DONE; parity via `_posAngleRadsC`
  - [x] `constrainLng(double lng)` — DONE
  - [x] `constrainLat(double lat)` — DONE
  - [x] `H3_EXPORT(degsToRads)(double degrees)` — DONE
  - [x] `H3_EXPORT(radsToDegs)(double radians)` — DONE
  - [x] `geoAlmostEqualThreshold(const LatLng*, const LatLng*, double)` — DONE
  - [x] `geoAlmostEqual(const LatLng*, const LatLng*)` — DONE
  - [x] `H3_EXPORT(greatCircleDistanceRads)(const LatLng*, const LatLng*)` — DONE
  - [x] `H3_EXPORT(greatCircleDistanceKm)(const LatLng*, const LatLng*)` — DONE
  - [x] `H3_EXPORT(greatCircleDistanceM)(const LatLng*, const LatLng*)` — DONE
  - [x] `_geoAzimuthRads(const LatLng*, const LatLng*)` — DONE
  - [x] `_geoAzDistanceRads(const LatLng*, double, double, LatLng*)` — DONE
  - [x] `normalizeLng(const double, const LongitudeNormalization)` — DONE
  - [x] `triangleEdgeLengthsToArea(double, double, double)` — DONE
  - [x] `triangleArea(const LatLng*, const LatLng*, const LatLng*)` — DONE
  - [x] `_setGeoRads(LatLng*, double, double)` and `setGeoDegs(LatLng*, double, double)` — DONE

- vec2d.c
  - [x] `_v2dMag(const Vec2d*)` — DONE
  - [x] `_v2dIntersect(const Vec2d*, const Vec2d*, const Vec2d*, const Vec2d*, Vec2d*)` — DONE
  - [x] `_v2dAlmostEquals(const Vec2d*, const Vec2d*)` — DONE

- vec3d.c
  - [x] `_pointSquareDist(const Vec3d*, const Vec3d*)` — DONE (Go port + C parity via vec3d shim; no C logic duplication)
  - [x] `_geoToVec3d(const LatLng*, Vec3d*)` — DONE (Go port + C parity via vec3d shim)

- coordijk.c
  - [x] `_ijkAdd(const CoordIJK*, const CoordIJK*, CoordIJK*)` — DONE (Go port + C parity)
  - [x] `_ijkSub(const CoordIJK*, const CoordIJK*, CoordIJK*)` — DONE (Go port + C parity)
  - [x] `_setIJK(CoordIJK*, int, int, int)` — DONE (Go port + C parity)
  - [x] `_ijkMatches(const CoordIJK*, const CoordIJK*)` — DONE (Go port + C parity)
  - [x] `_ijkScale(CoordIJK*, int)` — DONE (Go port + C parity)
  - [x] `_ijkNormalize(CoordIJK*)` — DONE (Go port + C parity)
  - [x] `ijkDistance(const CoordIJK*, const CoordIJK*)` — DONE (Go port + C parity)
  - [x] `_ijkRotate60ccw(CoordIJK*)` — DONE (Go port + C parity)
  - [x] `_ijkRotate60cw(CoordIJK*)` — DONE (Go port + C parity)
  - [x] `_unitIjkToDigit(const CoordIJK*)` — DONE (Go port + C parity)
  - [x] `_neighbor(CoordIJK*, Direction)` — DONE (Go port + C parity)
  - [x] `_rotate60ccw(Direction)` — DONE (Go port + C parity)
  - [x] `_rotate60cw(Direction)` — DONE (Go port + C parity)
  - [x] `ijkToCube(CoordIJK*)` — DONE (Go port + C parity)
  - [x] `cubeToIjk(CoordIJK*)` — DONE (Go port + C parity)
  - [x] `_ijkToHex2d(const CoordIJK*, Vec2d*)` — DONE (Go port + C parity)
  - [x] `_hex2dToCoordIJK(const Vec2d*, CoordIJK*)` — DONE (Go port + C parity)

- h3Index.c
  - [x] `H3_GET_RESERVED_BITS/H3_SET_RESERVED_BITS` — DONE (Go ports + cgo wrappers)
  - [x] `H3_GET_INDEX_DIGIT/H3_SET_INDEX_DIGIT` — DONE (Go ports + cgo wrappers)
  - [x] `_h3LeadingNonZeroDigit` — DONE (Go port + direct C call parity)

Next up (planned)
- polygon.c (next small helpers):
  - [x] `bboxFromGeoLoop(const GeoLoop*, BBox*)` — DONE
  - [x] `pointInsideGeoLoop(const GeoLoop*, const BBox*, const LatLng*)` — DONE
  - [x] `pointInsidePolygon(const GeoPolygon*, const BBox*, const LatLng*)` — DONE
  - [x] `cellBoundaryCrossesGeoLoop(const GeoLoop*, const BBox*, const CellBoundary*, const BBox*)` — DONE
  - [x] `cellBoundaryInsidePolygon(...)` — DONE
  - [x] `cellBoundaryCrossesPolygon(...)` — DONE
- h3Index.c (more utilities):
  - [x] `H3_GET_RESERVED_BITS/H3_SET_RESERVED_BITS` ports
  - [x] `H3_GET_INDEX_DIGIT/H3_SET_INDEX_DIGIT` ports
  - [x] `_h3LeadingNonZeroDigit` (interop + Go port)
  - [x] `H3_GET_MODE/H3_SET_MODE` — DONE (Go ports + cgo wrappers + parity)
  - [x] `H3_GET_HIGH_BIT/H3_SET_HIGH_BIT` — DONE (Go ports + cgo wrappers + parity)
  - [x] `_zeroIndexDigits(H3Index, int, int)` — DONE (Go port + C parity); note: C allows start=0 (overlaps base cell bits) and treats end>15 as no-op
  - [x] `H3_EXPORT(isResClassIII)(H3Index)` — DONE (Go port returns int for parity)
  - [x] `H3_EXPORT(isPentagon)(H3Index)` — DONE (Go port + C parity via tables.IsPentagonBaseCell and _h3LeadingNonZeroDigit)
  - [x] `setH3Index(H3Index*, int, int, Direction)` — DONE (Go port + parity)
  - [x] `makeDirectChild(H3Index, int)` — DONE (Go port + parity via inline C wrapper)
  - [x] `H3_EXPORT(cellToCenterChild)(H3Index, int, H3Index*)` — DONE (Go port + C parity)
  - [x] `H3_EXPORT(cellToParent)(H3Index, int, H3Index*)` — DONE (Go port + C parity)
  - [x] `_hasChildAtRes(H3Index, int)` — DONE (Go port + C parity via shim)
  - [x] `H3_EXPORT(cellToChildrenSize)(H3Index, int, int64_t*)` — DONE (Go port + C parity)
  - [x] `H3_EXPORT(pentagonCount)(void)` — DONE (Go port + C parity)
  - [x] `_isBaseCellPentagon(int)` — DONE (Go parity test via baseCells shim; local mirror in constants.go)
  - [x] `_isBaseCellPolarPentagon(int)` — DONE (Go port + C parity)
  - [x] `H3_EXPORT(res0CellCount)(void)` — DONE (Go port + C parity)
  - [x] `H3_EXPORT(cellToChildPos)(H3Index, int, int64_t*)` — DONE (Go port + C parity)
  - [x] `H3_EXPORT(childPosToCell)(int64_t, H3Index, int, H3Index*)` — DONE (Go port + C parity)
  - [x] `validateChildPos` (static) — DONE (Go port; covered indirectly by parity)
  - [x] `isResolutionClassIII(int res)` — DONE (Go port + C parity)

Planned small targets
- vec3d.c:
  - [x] `_geoToVec3d(const LatLng*, Vec3d*)` — DONE this iteration; converts LL to 3D vector
- coordijk.c:
  - [x] `_ijkAdd(const CoordIJK*, const CoordIJK*, CoordIJK*)` — DONE
  - [x] `_ijkSub(const CoordIJK*, const CoordIJK*, CoordIJK*)` — DONE
  - [x] `_setIJK(CoordIJK*, int, int, int)` — DONE; sets IJK coordinate components
  - [x] `_ijkMatches(const CoordIJK*, const CoordIJK*)` — DONE; compares IJK coordinates for equality
  - [x] `_ijkScale(CoordIJK*, int)` — DONE; scales IJK coordinates by factor
  - [x] `_ijkNormalize(CoordIJK*)` — DONE; normalizes IJK by removing negatives
  - [x] `ijkDistance(const CoordIJK*, const CoordIJK*)` — DONE; computes distance between IJK coordinates
  - [x] `_ijkRotate60ccw(CoordIJK*)` — DONE; rotates IJK coordinates 60° counter-clockwise
  - [x] `_ijkRotate60cw(CoordIJK*)` — DONE; rotates IJK coordinates 60° clockwise
  - [x] `_unitIjkToDigit(const CoordIJK*)` — DONE; converts unit IJK coordinate to direction digit
  - [x] `_neighbor(CoordIJK*, Direction)` — DONE; applies direction vector to IJK coordinates
  - [x] `_rotate60ccw(Direction)` — DONE; rotates direction digit 60° counter-clockwise
  - [x] `_rotate60cw(Direction)` — DONE; rotates direction digit 60° clockwise

Notes this iteration
- Verified available vec3d symbols via header; opted to plan `_geoToVec3d` as a safe, dependency-light target.

Execution plan per function
- Extend `<cfile>_cgo.go` with direct calls using C structs (C.BBox/C.LatLng/GeoLoop); avoid scalar params.
- Implement faithful Go ports in `<cfile>__<function>.go` with same unexported names.
- Add `<cfile>__<function>_parity_test.go` with realistic tolerances; compare bools directly.
Conventions
- One function per Go file: `<cfile>__<function>.go`
- cgo interop per C module: `<cfile>_cgo.go` (build tag `cgo && c2go`), includes `"<cfile>.c"` by name only.
- C wrappers named distinctly (e.g., `_ipow_c_wrapper`), Go helpers mirror with a `C` suffix (e.g., `_ipowC`).
- Parity tests named `<cfile>__<function>_parity_test.go` with `//go:build c2go`. Reserve `<cfile>__<function>_test.go` for plain-Go tests later.
- Do not hardcode H3 version in code; include directories provided via `CGO_CPPFLAGS` in `make test-c2go` with `H3VER`.

Issue encountered: linking large C modules (resolved)
- Symptom: Including `latLng.c` directly from a single cgo file pulled in references to many other H3 C symbols (e.g., `cellToBoundary`, `directedEdgeToBoundary`), causing undefined symbol errors or duplicate static symbol definitions when multiple C files were included together.
- Constraints: We must not copy or reimplement C logic in the cgo layer.
- Solution applied:
  - Compile original H3 C modules as separate translation units via build-tagged C shim files in this package (e.g., `h3lib_latLng_c2go.c` contains `#include "latLng.c"`, `h3lib_h3Index_c2go.c` contains `#include "h3Index.c"`, etc.). This avoids static symbol redefinition (e.g., duplicate `DIRECTIONS`) by keeping each module isolated.
  - Keep `mathExtensions.c` referenced only once (through `mathExtensions_cgo.go`) to avoid duplicate `_ipow`.
  - Pass include directories via `CGO_CPPFLAGS` in `make test-c2go` using `H3VER` (no version in code).
  - Add `CGO_CFLAGS=-ffunction-sections -fdata-sections` and `CGO_LDFLAGS=-Wl,-dead_strip` so the linker drops unused functions referenced by the compiled C objects.
  - Tests call tiny Go wrappers around the original C symbols; cgo usage remains only in non-test files with the `c2go` tag. No C code is copied into Go files.

Guideline going forward when a C function links in extra deps
- First, add or reuse a per-module C shim file `h3lib_<module>_c2go.c` that `#include`s the needed original C file(s) by name.
- Ensure each original C file is compiled in its own shim to avoid static symbol redefinitions across modules.
- If link still fails due to further deps, add shims for those modules too (one TU per module).
- Keep the Go interop file (`<cfile>_cgo.go`) limited to declaring wrappers that forward to C; never copy C logic.


Note on C bool interop (cgo)
- Preferred: include `<stdbool.h>` and compare C.bool return values directly to `0` in Go (`C.fn(...) != 0`). This is valid because `_Bool` is an integer type in C99.
- Toolchain caveat: some cgo toolchains reject direct comparison `C.bool != 0`. In those cases, use a tiny inline C helper to normalize to `int` (e.g., `static int h3_bool_to_int(_Bool b) { return b ? 1 : 0; }`) and compare its result to `0` from Go.

- faceijk.c
  - [x] `_geoToClosestFace(const LatLng*, int*, double*)` — DONE (Go port + C parity); finds closest icosahedral face center
  - [x] `_geoToHex2d(const LatLng*, int, int*, Vec2d*)` — DONE (Go port + C parity); converts geo to 2D hex coordinates on face
  - [x] `_hex2dToGeo(const Vec2d*, int, int, int, LatLng*)` — DONE (Go port + C parity); converts 2D hex coordinates back to geo

Next up (planned)
- faceijk.c (additional functions for future consideration):
  - [x] `_faceIjkToGeo(const FaceIJK*, int, LatLng*)` — DONE (Go port + C parity); converts FaceIJK coordinates to geographic coordinates
- latLng.c
  - [x] `normalizeLng(const double, const LongitudeNormalization)` — DONE (Go port + C parity); normalizes longitude values based on strategy
  - [x] `triangleEdgeLengthsToArea(double, double, double)` — DONE (Go port + C parity); calculates spherical triangle area from edge lengths
  - [x] `triangleArea(const LatLng*, const LatLng*, const LatLng*)` — DONE (Go port + C parity); calculates spherical triangle area from vertices
  - [x] `_setGeoRads(LatLng*, double, double)` — DONE (Go port + C parity); sets LatLng components in radians
  - [x] `setGeoDegs(LatLng*, double, double)` — DONE (Go port + C parity); sets LatLng components in degrees
- vec2d.c
  - [x] `_v2dMag(const Vec2d*)` — DONE (Go port + C parity); calculates magnitude of 2D vector
  - [x] `_v2dIntersect(const Vec2d*, const Vec2d*, const Vec2d*, const Vec2d*, Vec2d*)` — DONE (Go port + C parity); finds intersection between two lines
  - [x] `_v2dAlmostEquals(const Vec2d*, const Vec2d*)` — DONE (Go port + C parity); compares 2D vectors for equality
 - vec3d.c
  - [x] `_geoToVec3d(const LatLng*, Vec3d*)` — DONE (Go port + C parity); converts geographic coordinates to 3D unit vector
  - [x] `_square(double)` — DONE (Go port + C parity); returns square of a number
- h3Index.c
  - [ ] Next: `cellToChildren` iterator-based parity (depends on iterators); may be deferred
  - [ ] Other tiny getters: isValidCell (larger), error codes mapping (already wrapped)

Execution plan per function
- Extend `<cfile>_cgo.go` with direct calls to the original C functions using C structs where applicable (no scalar explosion; use C.LatLng / C.Vec2d).
- Implement the Go port in `<cfile>__<function>.go` (faithful translation; same unexported name where possible).
- Add parity test `<cfile>__<function>_parity_test.go` under `//go:build c2go` with tight but realistic tolerances.
- If C returns `bool`, prefer `<stdbool.h>` and compare `!= 0`. If toolchain warns, route through a tiny C helper returning `int`.

## Missing Functions from Gap Analysis

The following 96 functions were identified from comprehensive H3 source analysis but are not yet tracked in TODO:

  - [ ] `_adjustPentVertOverage` — TODO: needs C→Go port + parity test
  - [x] `_baseCellIsCwOffset` — DONE (Go port + C parity)
  - [x] `_baseCellToCCWrot60` — DONE (Go port + C parity)
  - [x] `_baseCellToFaceIjk` — DONE (Go port + C parity)
  - [x] `_downAp3` — DONE (Go port + C parity)
  - [x] `_downAp3r` — DONE (Go port + C parity)
  - [x] `_downAp7` — DONE (Go port + C parity)
  - [x] `_downAp7r` — DONE (Go port + C parity)
  - [ ] `_faceIjkPentToVerts` — TODO: needs C→Go port + parity test
  - [x] `_faceIjkToBaseCell` — DONE (Go port + C parity)
  - [x] `_faceIjkToBaseCellCCWrot60` — DONE (Go port + C parity)
  - [ ] `_faceIjkToH3` — TODO: needs C→Go port + parity test
  - [ ] `_faceIjkToVerts` — TODO: needs C→Go port + parity test
  - [x] `_firstOneIndex` — DONE (Go port + C parity)
  - [x] `_geoToFaceIjk` — DONE (Go port + C parity)
  - [x] `_getBaseCellDirection` — DONE (Go port + C parity)
  - [x] `_getBaseCellNeighbor` — DONE (Go port + C parity)
  - [x] `_getResDigit` — DONE (Go port + C parity)
  - [ ] `_gridRingInternal` — TODO: needs C→Go port + parity test
  - [x] `_h3Rotate60ccw` — DONE (Go port + C parity)
  - [x] `_h3Rotate60cw` — DONE (Go port + C parity)
  - [x] `_h3RotatePent60ccw` — DONE [parity test]
  - [x] `_h3RotatePent60cw` — DONE [parity test]
  - [ ] `_h3ToFaceIjk` — TODO: needs C→Go port + parity test
  - [ ] `_h3ToFaceIjkWithInitializedFijk` — TODO: needs C→Go port + parity test
  - [x] `_hasAll7AfterRes` — DONE (Go port + C parity)
  - [x] `_hasAny7UptoRes` — DONE (Go port + C parity)
  - [x] `_hasDeletedSubsequence` — DONE (Go port + C parity)
  - [x] `_hasGoodTopBits` — DONE (Go port + C parity)
  - [x] `_hashVertex` — DONE (Go port + C parity)
  - [ ] `_hexRadiusKm` — TODO: needs C→Go port + parity test
  - [x] `_ijkNormalizeCouldOverflow` — DONE (Go port + C parity)
  - [x] `_incrementResDigit` — DONE: iterators__incrementResDigit.go with parity test
  - [x] `_iterInitParent` — DONE: iterators__iterInitParent.go with parity test
  - [x] `_null_iter` — DONE: iterators__null_iter.go with parity test
  - [x] `_upAp7` — DONE (Go port + C parity)
  - [x] `_upAp7Checked` — DONE (Go port + C parity)  
  - [x] `_upAp7r` — DONE (Go port + C parity)
  - [x] `_upAp7rChecked` — DONE (Go port + C parity)
  - [ ] `_vertexGraphToLinkedGeo` — TODO: needs C→Go port + parity test
  - [x] `baseCellNumToCell` — DONE (Go port + C parity)
  - [x] `bboxCenter` — DONE (Go port + C parity)
  - [x] `bboxContains` — DONE (Go port + C parity)
  - [x] `bboxContainsBBox` — DONE (Go port + C parity)
  - [x] `bboxEquals` — DONE (Go port + C parity)
  - [x] `bboxesFromGeoPolygon` — DONE: polygon__bboxesFromGeoPolygon.go with parity test
  - [ ] `bboxHexEstimate` — TODO: needs C→Go port + parity test
  - [x] `bboxIsTransmeridian` — DONE (Go port + C parity)
  - [x] `bboxOverlapsBBox` — DONE (Go port + C parity)
  - [ ] `bboxToCellBoundary` — TODO: needs C→Go port + parity test
  - [ ] `cellToBBox` — TODO: needs C→Go port + parity test
  - [ ] `cellToLocalIjk` — TODO: needs C→Go port + parity test
  - [x] `countLinkedCoords` — DONE (Go port + C parity)
  - [x] `countLinkedLoops` — DONE (Go port + C parity)
  - [ ] `countLinkedPolygons` — TODO: needs C→Go port + parity test
  - [x] `cubeRound` — DONE (Go port + C parity)
  - [ ] `destroyLinkedGeoLoop` — TODO: needs C→Go port + parity test
  - [ ] `destroyVertexGraph` — TODO: needs C→Go port + parity test
  - [ ] `directionForNeighbor` — TODO: needs C→Go port + parity test
  - [ ] `directionForVertexNum` — TODO: needs C→Go port + parity test
  - [ ] `getAverageCellArea` — TODO: needs C→Go port + parity test
  - [x] `ijkToIj` — DONE [parity test]
  - [x] `ijToIjk` — DONE [parity test]
  - [x] `initVertexGraph` — DONE (Go port + C parity)
  - [ ] `iterDestroyPolygon` — TODO: needs C→Go port + parity test
  - [ ] `iterDestroyPolygonCompact` — TODO: needs C→Go port + parity test
  - [ ] `iterInitBaseCellNum` — TODO: needs C→Go port + parity test
  - [ ] `iterInitRes` — TODO: needs C→Go port + parity test
  - [ ] `iterStepChild` — TODO: needs C→Go port + parity test
  - [ ] `iterStepPolygon` — TODO: needs C→Go port + parity test
  - [ ] `iterStepPolygonCompact` — TODO: needs C→Go port + parity test
  - [ ] `iterStepRes` — TODO: needs C→Go port + parity test
  - [ ] `localIjkToCell` — TODO: needs C→Go port + parity test
  - [ ] `nextCell` — TODO: needs C→Go port + parity test
  - [ ] `normalizeMultiPolygon` — TODO: needs C→Go port + parity test
  - [ ] `removeVertexNode` — TODO: needs C→Go port + parity test
  - [x] `scaleBBox` — DONE [parity test]
  - [x] `validatePolygonFlags` — DONE [parity test]
  - [ ] `vertexNumForDirection` — TODO: needs C→Go port + parity test
  - [ ] `vertexRotations` — TODO: needs C→Go port + parity test
