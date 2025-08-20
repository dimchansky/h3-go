# h3 (pure Go) — TODO

**Goal:** Re-implement Uber’s H3 (v4.3.0) in **pure Go** (no cgo, no external deps), with a high-performance, allocation-aware API.

**Ground truth:** Behavior must match H3 C @ tag **v4.3.0**.
- Reference: https://github.com/uber/h3 (v4.3.0)
- API docs: https://h3geo.org/docs/api/indexing
- Go bindings (naming inspo only): https://github.com/uber/h3-go

> Work style rule: Keep this TODO up to date. Add/adjust sub-tasks as you discover dependencies. Check items off only when code + tests are merged.

---

## 0) Repo bootstrap
- [ ] Initialize repo: `module github.com/<you>/h3` (Go ≥ 1.22).
- [ ] Add `LICENSE` (Apache-2.0) and `NOTICE` with attribution to H3.
- [ ] Add CI (GitHub Actions): `go build`, `go test`, race, `go vet`, `golangci-lint`.
- [ ] Add `Makefile`: `make test`, `make bench`, `make lint`, `make gen` (for tables), `make ref` (build C oracle CLI).

## 1) Public API inventory
- [ ] Create `api.md` listing **all public functions** grouped by domain:
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
- [ ] For every public function, map dependencies to H3 C internals (v4.3.0): math, bit layout helpers, rotation tables, pentagon logic, etc.
- [ ] Produce a **bottom-up plan** that starts with primitives (bit packing/extraction, base-cell metadata) and climbs to public APIs.
- [ ] Capture the map in `internal/design/deps.md` (optional) and mirror it as a check-list below.

## 3) Implementation plan (bottom-up)
### 3.1 Low-level primitives
- [ ] H3Index pack/unpack (mode, res, base cell, 15×3-bit digits).
- [ ] Static tables: base cells, neighbors, face mappings, rotations, pentagon flags.
- [ ] Angle helpers: deg↔rad, normalization, clamping; tolerances.
- [ ] Coordinate systems: icosahedron face/IJ transforms.
- [ ] Pentagon handling: canonical rotations, missing neighbors.

### 3.2 Mid-level helpers
- [ ] LatLng ↔ local face coords; great-circle math as needed by H3.
- [ ] Cell center & boundary computation.
- [ ] Neighbor step / grid distance; ring traversal with pentagon rules.
- [ ] IJ (axial) conversions (cell ↔ CoordIJ).

### 3.3 Public functions (initial wave)
- [ ] `LatLngToCell`
- [ ] `CellToLatLng`
- [ ] `CellToBoundary`
- [ ] `AreNeighbors`
- [ ] `GridDistance`
- [ ] `KRing` / `HexRange` / `HexRangeDistances`
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

## Open questions / decisions
- [ ] Exact result ordering per function (documented and tested).
- [ ] Codegen strategy for large tables (`go:generate` vs. committed).
- [ ] Public constants for bounds (e.g., `MaxResolution = 15`) and size estimates per API.
- [ ] Build tags (e.g., future `unsafe` or SIMD paths) — default remains pure Go.

---

## Dependency breakdown (to be refined as we explore)
- Low-level: bit layout, tables → Mid-level: coords, boundary, neighbor step → High-level: rings, polygons, edges/vertices → Top-level: compact/uncompact, IJ mappings.

Keep this section evolving with links to code and checkmarks as pieces land.
