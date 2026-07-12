# Intentional deviations from H3 C

This is the checklist consulted during upstream synchronization
(docs/public-api-architecture.md §10): each item below is a *deliberate*
difference between this library and H3 C v4.4.0 (the current parity target).
Anything not listed here is expected to match C behavior exactly and is
enforced by the cgo parity suite. The `h3` command-line utility has its own
compatibility policy, documented in [cli-compatibility.md](cli-compatibility.md).

## Representation

1. **`Angle` type in geometry structs.** `LatLng.Lat/Lng` and the internal
   `bbox` fields are `Angle` (float64 radians inside) instead of raw
   `double` radians. Ported code unwraps with `.Rad()` at arithmetic leaves.
   Rationale: makes degree/radian confusion uncompilable; zero-cost because
   the public and internal representation are the same (DR-003).
2. **`CellBoundary` is encapsulated.** Same memory layout as the C struct
   (`numVerts` + fixed `[10]LatLng` array), but the fields are unexported and
   accessed via `Len`/`At`/`Verts`. (An earlier port draft used a growable
   slice; that was reverted to the C-faithful fixed array, fixing a
   6-allocations-per-call hotspot.)
3. **C `int` → Go `int32`, C `int64_t` → Go `int64`** in all ported code, to
   preserve 32-bit overflow behavior. Public wrappers take/return `int`
   (bounds-checked before narrowing) except where values exceed int32 range
   (`int64` for cell counts and child positions) or where a conversion copy
   would be forced (`[]int32` grid-disk distances).
4. **Linked-geo structures are internal.** `CellsToMultiPolygon` returns
   `[]GeoPolygon` (slices), not the C linked lists; `destroy*` functions are
   meaningless under GC and are not exposed.

## Behavior

5. **Hole pruning.** C fills caller-sized buffers and may leave `H3_NULL`
   holes (hash-set placement, pentagon slots). Public Go APIs prune: GridDisk
   / GridDiskDistances(+Safe) / GridRing (cells and distances in tandem),
   DirectedEdges / Vertexes / IcosahedronFaces (fixed-size pentagon slots),
   Compact (trailing nulls). Unsafe variants are dense ring-walk ordered in C
   already and are returned as-is. Internal ported functions keep exact C
   semantics — pruning happens only at the public boundary.
6. **Validated parsing.** C `stringToH3` accepts any hex string;
   `ParseCell`/`ParseDirectedEdge`/`ParseVertex` validate the index mode and
   reject syntax errors instead of swallowing them.
7. **Errors are sentinels, not codes.** The 15 C error codes surface as
   package-level `Err*` values (message text from C `describeH3Error`)
   matched with `errors.Is`. The numeric codes are internal.
8. **Resolution/k bounds are pre-checked** in public wrappers before int32
   narrowing, so absurd values (e.g. `res = 2^40+9`) cannot wrap into valid
   range. C relies on callers passing `int`.
9. **`gridRing` with negative k is memory-safe.** C `gridRing` reaches
   `memset(out, 0, 6*k*sizeof(H3Index))` after `gridRingUnsafe` rejects
   k < 0 — a negative size wrapped to a huge `size_t`, i.e. undefined
   behavior (it happens to be a no-op with macOS's end-pointer memset and a
   segfault with glibc's). The Go port zeroes by ranging over the caller's
   slice, so invalid k simply returns the domain error. Worth reporting
   upstream. The parity test therefore checks the Go behavior only for this
   input (C-UB inputs are out of parity scope by policy).
10. **Iterators are range-over-func.** The C iterator structs
   (`IterCellsChildren`, ...) are internal; `Cell.ChildrenSeq`, `CellsAtRes`,
   and `PolygonToCellsExperimentalSeq` expose `iter.Seq[Cell]`. Input
   validation happens before the sequence is returned; invalid input yields
   an empty sequence (C's null-iterator contract) or an upfront error.

## Naming

11. Ported implementation identifiers keep C names, unexported (mechanical
    case change only; see tools/unexport). Attribution comments
    (`Ported from H3 C: <file>::<name>`) are never rewritten. Public
    wrappers carry `H3 C API: <name>` doc lines; `make check-api` enforces
    that every C public function is either referenced this way or listed in
    the omissions table in tools/apiinventory/main.go
    (describeH3Error, degsToRads, radsToDegs, destroyLinkedMultiPolygon).
