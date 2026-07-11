package h3

// vertexToLatLng gets the geocoordinates of an H3 vertex.
// Converts a vertex H3 index to geographic coordinates (latitude and longitude).
// Ported from H3 C: vertex.c::vertexToLatLng.
func vertexToLatLng(vertex H3Index, coord *LatLng) H3Error {
	// Get the vertex number and owner from the vertex
	vertexNum := getReservedBits(vertex)
	owner := vertex
	owner = setMode(owner, H3_CELL_MODE)
	owner = setReservedBits(owner, 0)

	// Get the single vertex from the boundary
	var gb CellBoundary
	var fijk FaceIJK
	fijkError := _h3ToFaceIjk(owner, &fijk)
	if fijkError != E_SUCCESS {
		return fijkError
	}
	res := getResolution(owner)

	if isPentagon(owner) {
		_faceIjkPentToCellBoundary(&fijk, res, vertexNum, 1, &gb)
	} else {
		_faceIjkToCellBoundary(&fijk, res, vertexNum, 1, &gb)
	}

	// Copy from boundary to output coord
	*coord = gb.Verts[0]
	return E_SUCCESS
}
