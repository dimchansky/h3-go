package h3

// directedEdgeToBoundary provides the coordinates defining the directed edge.
// Gets the boundary coordinates that define a directed edge between two cells.
// The boundary may contain additional distortion vertices if the edge crosses
// an icosahedral face edge.
// Ported from H3 C: directedEdge.c::directedEdgeToBoundary.
func directedEdgeToBoundary(edge h3Index, cb *CellBoundary) h3Error {
	// Get the origin and neighbor direction from the edge
	dir := getReservedBits(edge)
	origin, originResult := getDirectedEdgeOrigin(edge)
	if originResult != eSuccess {
		return originResult
	}

	// Get the start vertex for the edge
	startVertex := vertexNumForDirection(origin, direction(dir))
	if startVertex == invalidVertexNum {
		// This is not actually an edge (i.e. no valid direction),
		// so return no vertices.
		cb.NumVerts = 0
		return eDirEdgeInvalid
	}

	// Get the geo boundary for the appropriate vertexes of the origin. Note
	// that while there are always 2 topological vertexes per edge, the
	// resulting edge boundary may have an additional distortion vertex if it
	// crosses an edge of the icosahedron.
	var fijk faceIJK
	fijkResult := _h3ToFaceIjk(origin, &fijk)
	if fijkResult != eSuccess {
		return fijkResult
	}
	res := getResolution(origin)
	isPent := isPentagon(origin)

	if isPent {
		_faceIjkPentToCellBoundary(&fijk, res, startVertex, 2, cb)
	} else {
		_faceIjkToCellBoundary(&fijk, res, startVertex, 2, cb)
	}
	return eSuccess
}
