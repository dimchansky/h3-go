# h3 (pure Go) — TODO

**Goal:** Re-implement Uber’s H3 (v4.3.0) in **pure Go** (no cgo, no external deps), with a high-performance, allocation-aware API.

**Ground truth:** Behavior must match H3 C @ tag **v4.3.0**.
- Reference: https://github.com/uber/h3 (v4.3.0)
- API docs: https://h3geo.org/docs/api/indexing
- Go bindings (naming inspo only): https://github.com/uber/h3-go

> Work style rule: Keep this TODO up to date. Add/adjust sub-tasks as you discover dependencies. Check items off only when code + tests are merged.

---

## 0) Repo bootstrap
- [x] Initialize repo: `module github.com/dimchansky/h3` (Go ≥ 1.22).
- [x] Add `LICENSE` (Apache-2.0); add `NOTICE` with attribution to H3.
- [x] Add CI (GitHub Actions): `go build`, `go test`, race, `go vet`, `golangci-lint`.
- [x] Add `Makefile`: `make test`, `make bench`, `make lint`, `make gen` (for tables), `make ref` (build C oracle CLI).

## 1) Public API inventory
- [x] Create `api.md` listing **all public functions** grouped by domain:
  - Indexing / Resolution / Base cells
  - Cell ↔ LatLng transforms (cell center, boundary)
  - Neighbors / Grid distance / Rings (k-ring, hexRange, etc.)
  - Directed/Undirected edges & vertices
  - IJ (axial) ↔ Cell conversions
  - Polygon to cells, containment
  - Compaction / Uncompaction
- For **each** function record:
  - [ ] Go-friendly signature using the **dst-buffer** pattern (e.g. `func KRing(dst []Cell, origin Cell, k int) ([]Cell, error)`).
  - [ ] Error conditions (incl. pentagon caveats, resolution bounds).
  - [ ] Output size bounds/estimates (to help caller preallocate).
  - [ ] Deterministic ordering of results.
  - [ ] Dependencies on internal helpers/tables.

## 2) Dependency mapping (C → Go) 
- [x] Analyzed H3 C internals and designed Go package structure
- [x] Implemented bottom-up plan: indexbits → tables → validate → public API
- [x] Documented dst-buffer patterns and error handling approach

## 3) Implementation plan (bottom-up)
### 3.1 Low-level primitives
- [x] H3Index pack/unpack (mode, res, base cell, 15×3-bit digits).
- [x] Static tables: base cells, neighbors, face mappings, rotations, pentagon flags (initial population in `internal/tables`).
- [x] Angle helpers: deg↔rad, normalization, clamping; tolerances.
- [x] Coordinate systems: icosahedron face/IJ transforms (Geo→FaceIJK and FaceIJK→H3 scaffolding with tests; parity-backed).
- [ ] Pentagon handling: canonical rotations, missing neighbors (FaceIJK rotations and cross-face traversal).

### 3.2 Mid-level helpers
- [ ] LatLng ↔ local face coords; great-circle math as needed by H3.
- [ ] Cell center & boundary computation.
- [ ] Neighbor step / grid distance; ring traversal with pentagon rules.
- [ ] IJ (axial) conversions (cell ↔ CoordIJ).

- ### 3.3 Public functions (initial wave)
- [x] `LatLngToCell` (implemented via FaceIJK path; oracle parity tests added)
- [ ] `Cell.ToLatLng` (center) 
- [ ] `Cell.ToBoundary`
- [ ] `Cell.IsNeighborOf`
- [ ] `Cell.DistanceTo`
- [ ] `Cell.KRing` / `HexRange` / `HexRangeDistances` / `HexRing`
- [x] `Cell.IsValid`, `Cell.Resolution`, `Cell.BaseCell`, `Cell.IsPentagon`
- [ ] `ToChildren` / `ToParent` / `Compact` / `Uncompact`
- [ ] `PolygonToCells`
- [ ] Edge APIs: `CellsToDirectedEdge`, `DirectedEdgeToCells`, `DirectedEdgeToBoundary`, etc.
- [ ] Vertex APIs: `CellToVertexes`, `CellToVertex`, `VertexToLatLng`, etc.
- [ ] IJ APIs: `CellToIJ`, `IJToCell`

> For every function above: follow the **dst-buffer** pattern and document ordering.

-## 4) Testing & correctness
- [ ] Implement C→Go error code mapping in `testref/` CLI and Go test harness (use the table documented in `api.md`); unknown codes map to `ErrFailed` (or wrapped) with logging.

- [x] Add oracle-backed parity suites for Geo→FaceIJK and FaceIJK→H3; randomized cases with caps.
- [ ] Mirror H3 C tests where available; otherwise design table tests for:
  - Invalid inputs (out of range, invalid indices)
  - Resolution extremes (0..MaxResolution)
  - Pentagons, poles, antimeridian
  - Deterministic ordering
- [ ] Build **external C oracle CLI** (no cgo) under `testref/`:
  - [x] `make ref` builds `testref/h3ref` from H3 v4.3.0 with a simple CLI protocol; commands mirror H3 C names.
  - [x] Go tests invoke `exec.Command("testref/h3ref", ...)` to get oracle outputs.
  - [ ] Cache small golden datasets where stable.
  - [x] `testref/README.md` documents oracle architecture and protocol.
- [ ] Add fuzz tests for reversible transforms (Cell ↔ LatLng ↔ Cell; edges/vertices roundtrips).
- [ ] Define numeric tolerances; start with `1e-12` (radians), `1e-9` (degrees).

## 5) Performance & allocations
- [ ] Benchmarks: `BenchmarkKRing`, `BenchmarkPolygonToCells`, `BenchmarkCellToBoundary`, etc.
- [x] Ensure functions reuse `dst` when capacity permits across implemented APIs; validate via tests.
- [ ] Avoid per-element `append` when size bounds known; prefer indexed fills.
- [x] No global mutable state; read-only tables only.
- [ ] Consider small, carefully measured inlining for hot helpers.

## 6) Docs & examples
- [ ] Package docs: index layout, pentagon caveats, dst-buffer pattern.
- [ ] Usage examples (compile as `Example...` tests).
- [x] `README.md` references oracle-backed tests and roadmap.
- [x] `api.md` serves as the canonical reference for signatures.

## 7) Lint/CI/quality gates
- [x] Enable `-race` in tests (test-only).
- [x] `go vet`, `golangci-lint` (incl. `staticcheck`, `ineffassign`, `gocritic`).
- [x] Added `smrcptr` for consistent receiver type checking.
- [ ] Coverage target: non-trivial code ≥ 85% where realistic.

## 8) License & attribution
- [x] Include Apache-2.0 license; add NOTICE crediting Uber H3.
- [x] Re-implement algorithms; do not copy C code verbatim. Cite sources.

---

## Milestone 1 Completed ✅

Status: Repository bootstrapped with core indexbits, angles, initial tables, and basic public metadata APIs. Oracle scaffolding in place. Tests compile and basic suites pass.

### Completed Items:
- [x] Repo bootstrap (go.mod, Makefile, LICENSE)
- [x] Internal packages: indexbits, angles, tables (initial), faceijk (Geo↔FaceIJK, FaceIJK→H3 scaffolding)  
- [x] Public metadata APIs: IsValid, Resolution, BaseCell, IsPentagon
- [x] LatLngToCell end-to-end path (Geo→FaceIJK→H3) with validation
- [x] Oracle-backed parity tests for FaceIJK transforms and LatLngToCell
- [x] README, api.md synced with approach; TODO maintained
- [x] Benchmark stubs for index packing/unpacking

### Next Milestone: Cell geometry & neighbors
- [ ] Implement `Cell.ToLatLng` (center)
- [ ] Implement `Cell.ToBoundary` (hex/pentagon; CCW canonical order)
- [ ] Implement neighbor check and grid distance
- [ ] Implement KRing/HexRange/HexRing and distance variants
- [ ] Expand tables and pentagon rotation handling
- [ ] Add table integrity and boundary parity tests

---

## Open questions / decisions
- [x] Exact result ordering per function (documented in api.md and tested).
- [ ] Codegen strategy for large tables (`go:generate` vs. committed).
- [x] Public constants for bounds (e.g., `MaxResolution = 15`) and size estimates per API.
- [ ] Build tags (e.g., future `unsafe` or SIMD paths) — default remains pure Go.

---

## Architecture Notes

The codebase follows a bottom-up dependency structure:
- **Low-level**: `internal/indexbits` (bit layout), `internal/tables` (constants), `internal/v2d` (2D vectors) ✅
- **Mid-level**: `internal/angles` (conversions), `internal/coordijk` (IJK coordinate system) ✅  
- **Coordinate systems**: `internal/faceijk` (face transforms, geo↔H3) ✅
- **High-level**: rings, polygons, edges/vertices (TODO)
- **Public API**: dst-buffer pattern, deterministic ordering ✅

All code compiles and tests pass. The foundation is solid for implementing geometric operations.

**Recent improvements:**
- Refactored CoordIJK to use consistent pointer receivers with method chaining
- Added comprehensive smrcptr linting for receiver type consistency
- Replaced Vec2d with performance-focused internal/v2d package
- Enhanced CI with fmt checking and smrcptr validation
- Added fix-fmt command for automatic code formatting
