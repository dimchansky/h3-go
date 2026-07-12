package h3

// This file contains types that mirror the H3 C API types defined in h3api.h

// h3Index mirrors C H3Index (uint64) from h3api.h.
//
// It is a type ALIAS of Cell, so the mechanically ported implementation and
// the public API share one type: []h3Index and []Cell are identical types,
// which is what makes every slice-producing algorithm zero-copy at the public
// boundary (docs/public-api-architecture.md, DR-003). Where the C code uses
// h3Index for directed-edge or vertex values, the public wrappers convert to
// DirectedEdge/Vertex at the boundary.
type h3Index = Cell

// LatLng mirrors the C struct from h3api.h
// Latitude and longitude stored as angle.Angle (internally in radians).
type LatLng struct {
	Lat Angle // latitude
	Lng Angle // longitude
}

// CellBoundary is the boundary of a cell or directed edge, in ccw order.
//
// It mirrors the C CellBoundary struct from h3api.h: a fixed-size array of
// up to MaxCellBoundaryVerts vertices plus a count. It is a value type;
// copying is cheap and involves no heap allocation. The zero value is an
// empty boundary. Access the vertices with Len, At, or Verts.
type CellBoundary struct {
	numVerts int32                        // number of vertices (matches C int)
	verts    [MaxCellBoundaryVerts]LatLng // vertices in ccw order
}

// GeoLoop is a simple loop of LatLng coordinates (closed implicitly).
// Mirrors the GeoLoop type from h3api.h.
type GeoLoop []LatLng

// GeoPolygon is an outer loop with optional holes.
// Mirrors the GeoPolygon struct from h3api.h.
type GeoPolygon struct {
	GeoLoop GeoLoop
	Holes   []GeoLoop
}

// CoordIJ represents IJ coordinates from h3api.h (axial coordinates).
// I and J use int32 because H3 C defines the full local-IJ domain with C int;
// this preserves its range and overflow behavior consistently on every Go
// platform and avoids silent narrowing at the implementation boundary.
type CoordIJ struct {
	I int32 // i component
	J int32 // j component
}

// linkedLatLng mirrors the C struct from h3api.h - a coordinate node in a linked geo structure.
type linkedLatLng struct {
	Vertex LatLng
	Next   *linkedLatLng
}

// linkedGeoLoop mirrors the C struct from h3api.h - a loop node in a linked geo structure.
type linkedGeoLoop struct {
	First *linkedLatLng
	Last  *linkedLatLng
	Next  *linkedGeoLoop
}

// linkedGeoPolygon mirrors the C struct from h3api.h - a polygon node in a linked geo structure.
type linkedGeoPolygon struct {
	First *linkedGeoLoop
	Last  *linkedGeoLoop
	Next  *linkedGeoPolygon
}
