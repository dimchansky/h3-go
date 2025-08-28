package h3

// cellToBoundary determines the cell boundary in spherical coordinates for an H3 index.
// Pentagon cells are handled with 5 vertices, hexagon cells with 6 vertices.
// All boundaries are returned in counter-clockwise order.
// Ported from H3 C: h3Index.c::cellToBoundary
func cellToBoundary(h3 H3Index, cb *CellBoundary) H3Error {
	var fijk FaceIJK
	err := _h3ToFaceIjk(h3, &fijk)
	if err != E_SUCCESS {
		return H3Error(err)
	}

	if isPentagon(h3) {
		_faceIjkPentToCellBoundary(&fijk, getResolution(h3), 0, NUM_PENT_VERTS, cb)
	} else {
		_faceIjkToCellBoundary(&fijk, getResolution(h3), 0, NUM_HEX_VERTS, cb)
	}

	return E_SUCCESS
}
