package c2go

// maxFaceCount returns the max number of possible icosahedron faces an H3 index may intersect.
// A pentagon always intersects 5 faces, a hexagon never intersects more than 2 (but may only intersect 1).
// Ported from H3 C: h3Index.c::maxFaceCount
func maxFaceCount(h3 H3Index, out *int32) H3Error {
	// a pentagon always intersects 5 faces, a hexagon never intersects more
	// than 2 (but may only intersect 1)
	if isPentagon(h3) {
		*out = 5
	} else {
		*out = 2
	}
	return E_SUCCESS
}
