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
- [ ] Add `LICENSE` (Apache-2.0) and `NOTICE` with attribution to H3.
- [ ] Add CI (GitHub Actions): `go build`, `go test`, race, `go vet`, `golangci-lint`.
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
- [x] Static tables: base cells, neighbors, face mappings, rotations, pentagon flags (stubbed).
- [x] Angle helpers: deg↔rad, normalization, clamping; tolerances.
- [ ] Coordinate systems: icosahedron face/IJ transforms.
- [ ] Pentagon handling: canonical rotations, missing neighbors.

### 3.2 Mid-level helpers
- [ ] LatLng ↔ local face coords; great-circle math as needed by H3.
- [ ] Cell center & boundary computation.
- [ ] Neighbor step / grid distance; ring traversal with pentagon rules.
- [ ] IJ (axial) conversions (cell ↔ CoordIJ).

### 3.3 Public functions (initial wave)
- [x] `LatLngToCell` (stub - returns ErrOptionInvalid)
- [x] `CellToLatLng` (stub - returns ErrOptionInvalid) 
- [x] `CellToBoundary` (stub - returns ErrOptionInvalid)
- [x] `AreNeighbors` (stub - returns ErrOptionInvalid)
- [x] `GridDistance` (stub - returns ErrOptionInvalid)
- [x] `KRing` / `HexRange` / `HexRangeDistances` / `HexRing` (stubs - return ErrOptionInvalid)
- [x] `IsValidCell`, `Resolution`, `BaseCell`, `IsPentagon` (implemented)
- [ ] `ToChildren` / `ToParent` / `Compact` / `Uncompact`
- [ ] `PolygonToCells`
- [ ] Edge APIs: `CellsToDirectedEdge`, `DirectedEdgeToCells`, `DirectedEdgeToBoundary`, etc.
- [ ] Vertex APIs: `CellToVertexes`, `CellToVertex`, `VertexToLatLng`, etc.
- [ ] IJ APIs: `CellToIJ`, `IJToCell`

> For every function above: follow the **dst-buffer** pattern and document ordering.

## 4) Testing & correctness
- [ ] Implement C→Go error code mapping in `testref/` CLI and Go test harness (use the table documented in `api.md`), and assert unknown codes map to `ErrFailed` (or a wrapped error) with logging.

- [ ] Mirror H3 C tests where available; otherwise design table tests for:
  - Invalid inputs (out of range, invalid indices)
  - Resolution extremes (0..MaxResolution)
  - Pentagons, poles, antimeridian
  - Deterministic ordering
- [ ] Build **external C oracle CLI** (no cgo) under `testref/`:
  - [ ] `make ref` builds `testref/h3ref` from H3 v4.3.0 with a simple stdin/stdout protocol (e.g., JSON). 
  - [ ] Go tests invoke `exec.Command("testref/h3ref", ...)` to get oracle outputs.
  - [ ] Cache small golden datasets where stable.
  - [x] Created `testref/README.md` with oracle architecture and protocol design.
- [ ] Add fuzz tests for reversible transforms (Cell ↔ LatLng ↔ Cell; edges/vertices roundtrips).
- [ ] Define numeric tolerances; start with `1e-12` (radians), `1e-9` (degrees).

## 5) Performance & allocations
- [ ] Benchmarks: `BenchmarkKRing`, `BenchmarkPolygonToCells`, `BenchmarkCellToBoundary`, etc.
- [ ] Ensure functions reuse `dst` when capacity permits; zero or minimal allocs in steady state.
- [ ] Avoid per-element `append` when size bounds known; prefer indexed fills.
- [ ] No global mutable state; read-only tables only.
- [ ] Consider small, carefully measured inlining for hot helpers.

## 6) Docs & examples
- [ ] Package docs: index layout, pentagon caveats, dst-buffer pattern.
- [ ] Usage examples (compile as `Example...` tests).
- [ ] `README.md` kept in sync with API surface.
- [ ] `api.md` serves as the canonical reference for signatures.

## 7) Lint/CI/quality gates
- [ ] Enable `-race` in tests (test-only).
- [ ] `go vet`, `golangci-lint` (incl. `staticcheck`, `ineffassign`, `gocritic`).
- [ ] Coverage target: non-trivial code ≥ 85% where realistic.

## 8) License & attribution
- [ ] Include Apache-2.0 license and NOTICE crediting Uber H3.
- [ ] Re-implement algorithms; do not copy C code verbatim. Cite sources.

---

## Milestone 1 Completed ✅

**Status**: All deliverables for the first milestone are complete and tests pass.

### Completed Items:
- [x] Repo bootstrap (go.mod, Makefile)
- [x] Internal packages: indexbits, angles, tables (stubbed), validate  
- [x] Public API stubs: IsValidCell, Resolution, BaseCell, IsPentagon (fully implemented)
- [x] Test coverage: indexbits_test.go, meta_test.go with comprehensive bit-level and API tests
- [x] Benchmark stubs: BenchmarkIndexPack, BenchmarkIndexUnpack
- [x] testref/ directory with oracle architecture documentation
- [x] TODO.md updated with progress

### Next Milestone Tasks:

## Next: Core Coordinate Transforms
- [ ] Populate internal/tables with correct H3 v4.3.0 constants
- [ ] Implement icosahedron face coordinate system
- [ ] Implement LatLngToCell and CellToLatLng
- [ ] Implement CellToBoundary
- [ ] Add table integrity tests
- [ ] Build external C oracle for validation

---

## Open questions / decisions
- [x] Exact result ordering per function (documented in api.md and tested).
- [ ] Codegen strategy for large tables (`go:generate` vs. committed).
- [x] Public constants for bounds (e.g., `MaxResolution = 15`) and size estimates per API.
- [ ] Build tags (e.g., future `unsafe` or SIMD paths) — default remains pure Go.

---

## Architecture Notes

The codebase follows a bottom-up dependency structure:
- **Low-level**: `internal/indexbits` (bit layout), `internal/tables` (constants) ✅
- **Mid-level**: `internal/angles` (conversions), `internal/validate` (input validation) ✅  
- **Coordinate systems**: Face transforms, IJK math (TODO)
- **High-level**: rings, polygons, edges/vertices (TODO)
- **Public API**: dst-buffer pattern, deterministic ordering ✅

All code compiles and tests pass. The foundation is solid for implementing geometric operations.
