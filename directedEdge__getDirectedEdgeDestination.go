package h3

// getDirectedEdgeDestination returns the destination hexagon from the directed edge H3Index.
//
// The function extracts the direction bits from the edge index, obtains the origin cell,
// and then uses h3NeighborRotations to find the neighboring cell in that direction.
//
// Ported from H3 C: directedEdge.c::getDirectedEdgeDestination.
func getDirectedEdgeDestination(edge H3Index, out *H3Index) H3Error {
	direction := Direction(getReservedBits(edge))
	rotations := int32(0)
	// Note: This call is also checking for H3_DIRECTEDEDGE_MODE
	origin, originResult := getDirectedEdgeOrigin(edge)
	if originResult != E_SUCCESS {
		return originResult
	}
	return h3NeighborRotations(origin, direction, &rotations, out)
}
