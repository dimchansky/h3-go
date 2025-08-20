# API Surface (Proposed) — h3 (pure Go)

This document lists the **public API** proposed for the pure-Go implementation of H3 (behavior matching H3 C v4.3.0).  
The API follows a **dst-buffer pattern** for collection-returning functions to minimize allocations and enable predictable performance.

> Status: proposal — evolves with implementation. Keep it in sync with code.

---

## Conventions

- **Indices**: `Cell`, `DirectedEdge`, `Vertex` are `uint64`-backed opaque types.
- **Units**: Public functions use degrees for `LatLng` inputs/outputs. Internals use radians.
- **Errors**: See package-level sentinel errors in `types.go`. Validate inputs (ranges, resolution bounds) and return the most specific error.
- **Deterministic ordering**:
  - When returning a **set** of indices, return **ascending by numeric index** unless the H3 spec mandates a canonical order (e.g., ring order, boundary winding). Document the order per function below.
- **dst-buffer pattern**: For `dst []T` parameters:
  - If `cap(dst)` is sufficient, reuse it (typically `dst = dst[:needed]`).
  - If not, allocate once to `needed` and return that slice.
  - Never retain references to caller-provided buffers after return.
- **Tolerances**: Angle comparisons use `EpsRad` (1e-12) internally; degrees comparisons use `EpsDeg` (1e-9).

---

## Types introduced by this API

```go
// From types.go (may be extended as needed).
type (
    Cell         uint64
    DirectedEdge uint64
    Vertex       uint64

    CoordIJ struct{ I, J int }
    LatLng  struct{ Lat, Lng float64 }

    GeoLoop    []LatLng
    GeoPolygon struct{ Outer GeoLoop; Holes []GeoLoop }

    CellBoundary []LatLng
)

// Additional helper type for range-with-distance results.
type CellDistance struct {
    Cell     Cell
    Distance int // grid distance (0..k)
}
```

---

## Indexing & Coordinate transforms

### Point ↔ Cell
```go
func LatLngToCell(p LatLng, res int) (Cell, error)                 // Order: n/a
func (c Cell) ToLatLng() (LatLng, error)                           // Order: n/a
func (c Cell) ToBoundary(dst []LatLng) ([]LatLng, error)           // Order: boundary-winding, counterclockwise, starting at a canonical vertex
```
- Errors: `ErrLatLngDomain`, `ErrResolutionDomain`, `ErrCellInvalid`, `ErrPentagon` (boundary around pentagons when relevant).
- Output bounds: `CellToBoundary` returns up to 6 or 5 vertices (hex or pentagon). Caller may preallocate `dst` of length 7 to be safe.

### Index metadata
```go
func (c Cell) IsValid() bool
func (c Cell) IsPentagon() (bool, error)
func (c Cell) Resolution() (int, error)
func (c Cell) BaseCell() (int, error) // 0..121
```
- Errors: `ErrCellInvalid` when index mode or fields are malformed.

---

## Hierarchy (parent/children, compaction)

```go
func ToParent(c Cell, parentRes int) (Cell, error) // Errors: ErrResolutionDomain, ErrCellInvalid
func ToChildren(dst []Cell, c Cell, childRes int) ([]Cell, error) // Order: ascending by Cell
func Compact(dst []Cell, cells []Cell) ([]Cell, error)            // Order: ascending by Cell
func Uncompact(dst []Cell, cells []Cell, res int) ([]Cell, error) // Order: ascending by Cell
```
- Output bounds: `ToChildren` ⇒ at most `7^(childRes-parentRes)` (pentagon paths differ). Expose a helper if needed.

---

## Neighborhoods & distances

```go
func (a Cell) IsNeighborOf(b Cell) (bool, error)

func (a Cell) DistanceTo(b Cell) (int, error) // hex grid distance (>= 0)

func (c Cell) KRing(dst []Cell, k int) ([]Cell, error) // Order: ascending by Cell

func (c Cell) HexRange(dst []Cell, k int) ([]Cell, error) // Synonym of KRing semantics (bounded traversal)

func (c Cell) HexRangeDistances(dst []CellDistance, k int) ([]CellDistance, error) // Order: stable by (distance, Cell)

func (c Cell) HexRing(dst []Cell, k int) ([]Cell, error) // Cells exactly at distance k; Order: ring-walk order (documented)
```
- Errors: `ErrCellInvalid`, `ErrPentagon` (when traversal crosses pentagon distortions), `ErrResolutionDomain` if needed.
- Helpers:
```go
func MaxKRingSize(k int) int // 3*k*(k+1) + 1
```

---

## Directed/Undirected Edges

```go
func CellsToDirectedEdge(a, b Cell) (DirectedEdge, error) // Error: ErrNotNeighbors, ErrCellInvalid
func DirectedEdgeToCells(e DirectedEdge) (origin Cell, destination Cell, err error)

func DirectedEdgesFromCell(dst []DirectedEdge, origin Cell) ([]DirectedEdge, error) // Order: ascending by edge index

func DirectedEdgeToBoundary(dst []LatLng, e DirectedEdge) ([]LatLng, error) // 2 points (or segment boundary polyline per H3 spec)
```
- Note: H3 models **directed** edges. Undirected semantics can be derived by sorting endpoints.

---

## Vertices

```go
func CellToVertexes(dst []Vertex, c Cell) ([]Vertex, error) // Up to 6 or 5
func CellToVertex(c Cell, vertexNum int) (Vertex, error)    // vertexNum in [0..numVerts-1]
func VertexToLatLng(v Vertex) (LatLng, error)
```
- Order: vertex ordering is canonical (counterclockwise around the cell starting from canonical vertex 0).

---

## IJ (local axial coords) — experimental in H3 C

Anchored at a chosen origin cell.

```go
func CellToLocalIJ(origin, c Cell) (CoordIJ, error)
func LocalIJToCell(origin Cell, ij CoordIJ) (Cell, error)
```
- Errors: `ErrResolutionMismatch` (if origin and c are at different resolutions where unsupported), `ErrPentagon`, `ErrCellInvalid`.

---

## Polygons & containment

```go
// Fill a polygon with cells at resolution `res`.
func PolygonToCells(dst []Cell, poly GeoPolygon, res int, mode ContainmentMode) ([]Cell, error) // Order: ascending by Cell

// Optional reverse operation: assemble rings (multi-polygons) from a set of cells.
func CellsToMultiPolygon(cells []Cell) ([]GeoPolygon, error) // Order: polygons sorted by min index; rings CCW for outers, CW for holes
```
- Errors: `ErrResolutionDomain`, `ErrOptionInvalid` (mode), `ErrLatLngDomain` (invalid polygon coords), `ErrDuplicateInput` (if applicable).
- Helpers:
```go
func MaxPolygonToCellsSize(poly GeoPolygon, res int, mode ContainmentMode) (int, error)
```

---

## Utilities & helpers

```go
// Counts total cells at a resolution.
func NumCellsAtResolution(res int) (int64, error) // Formerly getNumCells in H3

// String conversions (optional nicety; hex form)
func (c Cell) String() string
func ParseCell(s string) (Cell, error)
```
- `String`/`ParseCell` are not required by H3 spec but are helpful for tooling and tests.

---

## Ordering summary

- **Sets**: ascending numeric order by default (cells, edges, vertices).
- **Rings**: `HexRing` returns a ring-walk starting from a canonical neighbor and proceeding CCW (document exact rule in docstring).
- **Boundaries**: `CellToBoundary` returns vertices in CCW order starting from the canonical vertex for the cell’s orientation/face.

---

## Output size hints (for preallocation)

- `CellToBoundary`: at most 7 points (safe cap); typically 6 (hex) or 5 (pentagon).
- `KRing/HexRange`: at most `MaxKRingSize(k)` elements.
- `DirectedEdgesFromCell`: up to 6 (or 5 for pentagon).
- `CellToVertexes`: up to 6 (or 5).
- `ToChildren`: up to `7^(Δres)`, with pentagon caveats.
- `PolygonToCells`: depends on area and resolution; use `MaxPolygonToCellsSize` to pre-size where feasible.

---

## Error mapping notes

- Validate **resolution**: `0 <= res <= MaxResolution` → else `ErrResolutionDomain`.
- Validate **lat/lng**:
  - `Lat ∈ [-90, 90]`, `Lng ∈ (-180, 180]` (normalize longitudes where appropriate). Out-of-range → `ErrLatLngDomain`.
- Validate **index modes/fields**: malformed bit patterns → `ErrCellInvalid` (or corresponding type error for edges/vertices).
- Pentagon special cases may yield `ErrPentagon` where H3 C surfaces it.

---

## Future extensions (non-breaking)

- Build tags for optional acceleration backends (still pure Go by default).
- Additional helpers mirroring H3 utility functions (e.g., `GreatCircleDistance`, if needed by users).


## Error code mapping (C → Go)

To maintain behavioral parity with H3 C v4.3.0 and simplify tests that use the external C oracle,
we adopt the following **numeric error code → Go error** mapping (borrowed from h3-go, with a corrected
name for `ErrResolutionMismatch`). The pure-Go library itself returns the sentinel errors directly;
the numeric codes are only used by the test oracle/adapter.

```go
// C error code -> Go error sentinel mapping (parity with H3 C v4.3.0).
var cErrToGo = map[uint32]error{
    0:  nil, // Success
    1:  ErrFailed,
    2:  ErrDomain,
    3:  ErrLatLngDomain,
    4:  ErrResolutionDomain,
    5:  ErrCellInvalid,
    6:  ErrDirectedEdgeInvalid,
    7:  ErrUndirectedEdgeInvalid,
    8:  ErrVertexInvalid,
    9:  ErrPentagon,
    10: ErrDuplicateInput,
    11: ErrNotNeighbors,
    12: ErrResolutionMismatch, // corrected name
    13: ErrMemoryAlloc,
    14: ErrMemoryBounds,
    15: ErrOptionInvalid,
}
```
