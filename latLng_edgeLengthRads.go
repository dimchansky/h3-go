package h3

// edgeLengthRads returns the length of a directed edge in radians.
// The length is calculated by summing the great circle distances between
// consecutive vertices of the edge boundary.
// Ported from H3 C: latLng.c::H3_EXPORT(edgeLengthRads).
func edgeLengthRads(edge h3Index, length *float64) h3Error {
	var cb CellBoundary

	err := directedEdgeToBoundary(edge, &cb)
	if err != eSuccess {
		return err
	}

	*length = 0.0
	for i := int32(0); i < cb.numVerts-1; i++ {
		*length += greatCircleDistanceRads(&cb.verts[i], &cb.verts[i+1])
	}

	return eSuccess
}
