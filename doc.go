// Package h3 is a pure-Go implementation of Uber's H3 hexagonal hierarchical
// geospatial indexing system, behaviorally equivalent to H3 C v4.5.0.
//
// The production library is safe Go only: no cgo and no unsafe. Behavioral
// equivalence is enforced by an opt-in parity test suite (build tags
// cgo && c2go) that compares every ported function against the original C
// objects compiled from pristine upstream sources.
//
// # Cells, resolutions, and the grid
//
// H3 tessellates the sphere with hexagonal cells at 16 resolutions, from 0
// (coarsest, 122 base cells) to [MaxResolution] (finest). Each finer
// resolution subdivides the grid roughly sevenfold ([NumCells] gives the
// exact count). The grid is built by projecting the sphere onto a regular
// icosahedron: 12 cells per resolution — one at each icosahedron vertex —
// are pentagons rather than hexagons ([Cell.IsPentagon], [Pentagons]).
// Successive resolutions alternate between two orientations: Class II
// (even) and Class III (odd), the latter rotated ~19.1°
// ([Cell.IsResClassIII]).
//
// # Pentagons and icosahedron distortion
//
// Pentagons and icosahedron faces are the two systematic irregularities of
// the grid, and they surface in the API. Fast "Unsafe" traversal variants
// fail with [ErrPentagon] when they meet a pentagon or its distortion area;
// grid paths and local IJ conversions can fail near pentagons with
// [ErrFailed]; and cell or edge boundaries gain extra "distortion vertices"
// where they cross icosahedron faces ([Cell.Boundary],
// [DirectedEdge.Boundary]).
//
// # Hierarchy is logical, not geometric
//
// Parent/child relationships ([Cell.Parent], [Cell.Children]) are defined
// on the index hierarchy, not by geometric containment: a child's boundary
// is not required to lie within its parent's boundary, and a cell's center
// ([Cell.LatLng]) is not guaranteed to be its geographic centroid. Exact
// subdivision of a hexagon into hexagons is impossible, so each level's
// children only approximately cover the parent. Use polygon operations
// (for example [PolygonToCellsExperimental] with [ContainmentFull]) when
// geometric containment matters.
//
// # Index types
//
// The three H3 index kinds are distinct uint64 types: [Cell] (a hexagon or
// pentagon at some resolution), [DirectedEdge] (a directed connection
// between adjacent cells), and [Vertex] (a topological cell corner). Raw
// conversions like Cell(0x8928308280fffff) are legal but unchecked; use
// IsValid, or parse validated values with [ParseCell], [ParseDirectedEdge],
// and [ParseVertex]. All three implement fmt.Stringer (canonical lowercase
// hex) and encoding.TextMarshaler/TextUnmarshaler, so they work with
// encoding/json out of the box.
//
// # Coordinates
//
// Geographic coordinates use [LatLng], whose fields are [Angle] values
// (stored in radians). Construct them explicitly — [LatLngDegs] for degrees,
// [NewLatLng] with [Deg] or [Rad] angles — so degree/radian mix-ups cannot
// compile. Accessors convert on the way out: ll.Lat.Deg(), ll.Lng.Rad().
// Polygon inputs ([GeoLoop], [GeoPolygon]) use implicitly closed loops of
// these radians-backed coordinates and handle antimeridian crossings
// automatically. Local IJ coordinates ([CoordIJ]) are origin-relative and
// not stable across H3 versions; see [CellToLocalIJ].
//
// # Earth model and measurements
//
// All kilometer and meter results — areas, edge lengths, great-circle
// distances, and their per-resolution averages — are spherical
// approximations computed on a sphere of WGS84 authalic radius
// 6371.007180918475 km, not ellipsoidal geodesic values. Exact per-cell and
// per-edge measurements ([Cell.AreaKm2], [DirectedEdge.LengthKm]) account
// for the actual, distortion-aware boundary; the HexagonAreaAvg and
// HexagonEdgeLengthAvg families return per-resolution hexagon averages.
//
// # Ordering and stability
//
// Result ordering is documented per function and never exceeds the H3 C
// public contract: some results are explicitly unordered ([Cell.GridDisk],
// [PolygonToCells], [CompactCells]), some have a documented structure
// (increasing ring distance for the unsafe disk variants, canonical child
// order for [Cell.Children]), and none are silently sorted. Grid paths and
// local IJ coordinate spaces are not guaranteed to be stable across H3
// versions; do not persist them across library upgrades.
//
// # Errors
//
// Operations return sentinel errors matching the H3 C error codes
// ([ErrPentagon], [ErrCellInvalid], [ErrResolutionDomain], ...); match them
// with errors.Is. Each sentinel's documentation names representative
// operations that return it. Pure bit accessors (Resolution,
// BaseCellNumber, IsValid, String) do not fail, and parse/unmarshal syntax
// errors wrap strconv errors rather than sentinels (see [ParseCell]).
//
// # Allocation control
//
// Collection-returning operations come in two forms: a convenience form
// that allocates its result (GridDisk, Children, PolygonToCells, ...) and
// an Append* form (AppendGridDisk, AppendChildren, AppendPolygonToCells,
// ...) that appends into a caller-provided buffer. Append* forms add no
// allocation for the result itself when capacity suffices, but some may
// allocate internal working storage — for example [AppendCompactCells] for
// nontrivial inputs, or the gridDisk family's distance scratch on its
// pentagon fallback path — so the per-function documentation is
// authoritative. The [Cell.ChildrenSeq] and [CellsAtRes] iterators yield
// cells one at a time and allocate nothing (assertion-backed);
// [PolygonToCellsExperimentalSeq] may allocate internal iterator state.
// [CellBoundary] is a fixed-size value type; obtaining a boundary performs
// no heap allocation.
//
// # Relationship to H3 C
//
// The implementation is a function-by-function port of the C library; every
// public operation's doc comment carries an "H3 C API:" line naming its C
// counterpart, and docs/c-api-inventory.csv maps the entire C API surface.
// All 79 public functions of H3 C v4.5.0 are covered. Intentional behavior
// differences (Go-idiomatic hole pruning, validated parsing, ...) are
// documented in docs/DEVIATIONS.md.
//
// The repository also provides a pure-Go, upstream-compatible h3
// command-line utility built on this package:
//
//	go install github.com/dimchansky/h3-go/cmd/h3@latest
package h3
