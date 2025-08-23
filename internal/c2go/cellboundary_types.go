package c2go

// CellBoundary is a polyline of vertices (closed or open depending on context).
// Mirrors the CellBoundary struct from h3api.h
type CellBoundary struct {
	NumVerts int32    // number of vertices (matches C int)
	Verts    []LatLng // vertices in ccw order
}
