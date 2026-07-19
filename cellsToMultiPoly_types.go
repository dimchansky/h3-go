package h3

// geoMultiPolygon mirrors the GeoMultiPolygon struct from h3api.h: a
// simplified multiple-polygon type. It is internal in this port —
// the public CellsToMultiPolygon API returns []GeoPolygon
// (docs/DEVIATIONS.md §4); the ported 4.5.0 multipolygon pipeline and
// the area helpers consume this C-shaped form.
// Ported from H3 C: h3api.h.in::GeoMultiPolygon.
type geoMultiPolygon struct {
	NumPolygons int32
	Polygons    []GeoPolygon
}

// hashTableMultiplier is the arc hash-table sizing factor.
// After rough search, 10 seems to minimize compute time for large sets.
// Ported from H3 C: cellsToMultiPoly.h::HASH_TABLE_MULTIPLIER.
const hashTableMultiplier = 10

// arc is one directed cell edge in the arc-cancellation multipolygon
// algorithm: a node in a doubly-linked CCW loop of edges and in a
// union-find forest of connected components. The C struct's Arc*
// pointers are Go pointers into the arcSet's arcs slice (allocated
// once, never reallocated, so element pointers stay valid).
// Ported from H3 C: cellsToMultiPoly.h::Arc.
type arc struct {
	id h3Index

	isVisited bool
	isRemoved bool

	// For doubly-arced list of edges in loop.
	next *arc
	prev *arc

	// For union-find datastructure
	// https://en.wikipedia.org/wiki/Disjoint-set_data_structure
	parent *arc
	rank   int64
}

// arcSet is the full set of arcs for a cell set plus the
// open-addressing hash buckets for fast edge/arc lookup.
// Ported from H3 C: cellsToMultiPoly.h::ArcSet.
type arcSet struct {
	numArcs int64
	arcs    []arc

	// hash buckets for fast edge/arc lookup
	numBuckets int64
	buckets    []*arc
}

// sortableLoop is a boundary loop tagged with its connected-component
// root and enclosed area for sorting. C's GeoLoop{numVerts, verts}
// maps to the slice-backed GeoLoop of this port.
// Ported from H3 C: cellsToMultiPoly.h::SortableLoop.
type sortableLoop struct {
	root h3Index
	area float64

	loop GeoLoop
}

// sortableLoopSet is the sorted set of all boundary loops.
// Ported from H3 C: cellsToMultiPoly.h::SortableLoopSet.
type sortableLoopSet struct {
	numLoops int64
	sloops   []sortableLoop
}

// sortablePoly is a polygon tagged with its outer-loop area for the
// final area-descending polygon sort.
// Ported from H3 C: cellsToMultiPoly.h::SortablePoly.
type sortablePoly struct {
	outerArea float64
	poly      GeoPolygon
}
