package h3

// getDirectedEdgeDestination returns the destination hexagon from the directed edge h3Index.
//
// The function extracts the direction bits from the edge index, obtains the origin cell,
// and then uses h3NeighborRotations to find the neighboring cell in that direction.
//
// Ported from H3 C: directedEdge.c::getDirectedEdgeDestination.
func getDirectedEdgeDestination(edge h3Index, out *h3Index) h3Error {
	direction := direction(getReservedBits(edge))
	rotations := int32(0)
	// Note: This call is also checking for h3DirectededgeMode
	origin, originResult := getDirectedEdgeOrigin(edge)
	if originResult != eSuccess {
		return originResult
	}
	return h3NeighborRotations(origin, direction, &rotations, out)
}
