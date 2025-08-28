package h3

// edgeLengthRads returns the length of a directed edge in radians.
// The length is calculated by summing the great circle distances between
// consecutive vertices of the edge boundary.
// Ported from H3 C: latLng.c::H3_EXPORT(edgeLengthRads)
func edgeLengthRads(edge H3Index, length *float64) H3Error {
	var cb CellBoundary

	err := directedEdgeToBoundary(edge, &cb)
	if err != E_SUCCESS {
		return err
	}

	*length = 0.0
	for i := int32(0); i < cb.NumVerts-1; i++ {
		*length += greatCircleDistanceRads(&cb.Verts[i], &cb.Verts[i+1])
	}

	return E_SUCCESS
}
