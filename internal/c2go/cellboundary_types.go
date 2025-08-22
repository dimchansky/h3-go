package c2go

// CellBoundary is a polyline of vertices (closed or open depending on context).
type CellBoundary struct {
	NumVerts int
	Verts    []LatLng
}
