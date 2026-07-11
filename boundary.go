package h3

// Len returns the number of vertices in the boundary.
func (b *CellBoundary) Len() int { return int(b.numVerts) }

// At returns the i-th boundary vertex, in counterclockwise order.
// It panics if i is out of range [0, Len()).
func (b *CellBoundary) At(i int) LatLng {
	if i < 0 || i >= int(b.numVerts) {
		panic("h3: CellBoundary index out of range")
	}
	return b.verts[i]
}

// Verts returns the boundary vertices in counterclockwise order as a slice
// aliasing the boundary's internal storage. The slice is valid as long as b
// is; it involves no allocation.
func (b *CellBoundary) Verts() []LatLng { return b.verts[:b.numVerts] }
