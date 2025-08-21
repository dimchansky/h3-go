// Package h3 provides a pure-Go implementation of Uber's H3 hexagonal hierarchical
// geospatial index. This package targets behavioral equivalence with H3 C v4.3.0
// while offering Go-first APIs designed for performance and minimal allocations.
//
// Design note: any function that returns a collection will accept an optional
// destination buffer `dst []T` (which may be nil). Implementations will attempt
// to reuse `dst` capacity to avoid allocations and return the resulting slice.
// Callers must not rely on the backing array surviving beyond the call and
// the package will not retain references to caller-provided buffers.
package h3

import "errors"

// Version metadata for reference implementation parity.
const (
	// MaxResolution is the maximum supported H3 resolution (inclusive).
	// H3 defines resolutions in [0..15].
	MaxResolution = 15
)

// Sentinel errors used throughout the package.
var (
	ErrFailed                = errors.New("the operation failed")
	ErrDomain                = errors.New("argument was outside of acceptable range")
	ErrLatLngDomain          = errors.New("latitude or longitude were outside of acceptable range")
	ErrResolutionDomain      = errors.New("resolution argument was outside of acceptable range")
	ErrCellInvalid           = errors.New("H3Index cell argument was not valid")
	ErrDirectedEdgeInvalid   = errors.New("H3Index directed edge argument was not valid")
	ErrUndirectedEdgeInvalid = errors.New("H3Index undirected edge argument was not valid")
	ErrVertexInvalid         = errors.New("H3Index vertex argument was not valid")
	ErrPentagon              = errors.New("pentagon distortion was encountered")
	ErrDuplicateInput        = errors.New("duplicate input in arguments")
	ErrNotNeighbors          = errors.New("H3Index cell arguments were not neighbors")
	ErrResolutionMismatch    = errors.New("H3Index cell arguments had incompatible resolutions")
	ErrMemoryAlloc           = errors.New("memory allocation failed")
	ErrMemoryBounds          = errors.New("provided memory bounds were not large enough")
	ErrOptionInvalid         = errors.New("mode or flags argument was not valid")
)

// Core opaque index types (H3Index is a 64-bit unsigned integer).
type (
	// Cell identifies a single hexagon (or pentagon) at a given resolution.
	Cell uint64

	// DirectedEdge identifies a directed edge between two cells.
	DirectedEdge uint64

	// Vertex identifies a single topological vertex shared by three cells.
	Vertex uint64

	// CoordIJ are axial coordinates on the hex grid (120° axes).
	CoordIJ struct {
		I, J int
	}

	// LatLng are geographic coordinates in degrees.
	LatLng struct {
		Lat, Lng float64
	}

	// GeoLoop is an ordered list of LatLng making a loop (closed implicitly).
	GeoLoop []LatLng

	// GeoPolygon is an outer loop with optional holes.
	GeoPolygon struct {
		Outer GeoLoop
		Holes []GeoLoop
	}

	// CellBoundary is an ordered list of LatLng forming a cell's boundary.
	// Note: len(CellBoundary) will never exceed the maximum vertex count for the cell.
	CellBoundary []LatLng
)

// ContainmentMode specifies polygon containment behavior for PolygonToCells.
// Values will mirror H3 semantics; placeholders until implemented.
type ContainmentMode int

const (
	// ContainmentModeDefault is a placeholder; exact values TBD per API mapping.
	ContainmentModeDefault ContainmentMode = iota
)

// Internal numeric tolerances (subject to tuning during implementation).
const (
	// EpsRad is a small epsilon in radians used for angular comparisons.
	EpsRad = 1e-12

	// EpsDeg is a small epsilon in degrees used for geographic comparisons.
	EpsDeg = 1e-9
)

// --- API shape examples (to be implemented) ---
// The following are illustrative signatures to set expectations for the dst-buffer
// pattern. The concrete set will be finalized in api.md and implemented incrementally.
//
// func LatLngToCell(latlng LatLng, res int) (Cell, error)
// func CellToLatLng(c Cell) (LatLng, error)
// func CellToBoundary(dst []LatLng, c Cell) ([]LatLng, error)
// func KRing(dst []Cell, origin Cell, k int) ([]Cell, error)
// func PolygonToCells(dst []Cell, poly GeoPolygon, res int) ([]Cell, error)
