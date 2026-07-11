package h3

// This file contains types that mirror the H3 C API types defined in h3api.h

// H3Index mirrors C H3Index (uint64) from h3api.h.
//
// It is a type ALIAS of Cell, so the mechanically ported implementation and
// the public API share one type: []H3Index and []Cell are identical types,
// which is what makes every slice-producing algorithm zero-copy at the public
// boundary (docs/public-api-architecture.md, DR-003). Where the C code uses
// H3Index for directed-edge or vertex values, the public wrappers convert to
// DirectedEdge/Vertex at the boundary.
type H3Index = Cell

// LatLng mirrors the C struct from h3api.h
// Latitude and longitude stored as angle.Angle (internally in radians).
type LatLng struct {
	Lat Angle // latitude
	Lng Angle // longitude
}

// CellBoundary mirrors the CellBoundary struct from h3api.h
// A polyline of vertices (closed or open depending on context).
type CellBoundary struct {
	NumVerts int32    // number of vertices (matches C int)
	Verts    []LatLng // vertices in ccw order
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
// Uses int32 to match H3 C implementation exactly (including overflow behavior).
type CoordIJ struct {
	I int32 // i component
	J int32 // j component
}

// LinkedLatLng mirrors the C struct from h3api.h - a coordinate node in a linked geo structure.
type LinkedLatLng struct {
	Vertex LatLng
	Next   *LinkedLatLng
}

// LinkedGeoLoop mirrors the C struct from h3api.h - a loop node in a linked geo structure.
type LinkedGeoLoop struct {
	First *LinkedLatLng
	Last  *LinkedLatLng
	Next  *LinkedGeoLoop
}

// LinkedGeoPolygon mirrors the C struct from h3api.h - a polygon node in a linked geo structure.
type LinkedGeoPolygon struct {
	First *LinkedGeoLoop
	Last  *LinkedGeoLoop
	Next  *LinkedGeoPolygon
}
