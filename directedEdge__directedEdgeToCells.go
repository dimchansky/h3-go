package h3

// directedEdgeToCells returns the origin and destination cell pair for a directed edge.
//
// This function extracts both the origin and destination cells from a directed edge
// H3 index by calling getDirectedEdgeOrigin and getDirectedEdgeDestination sequentially.
// The results are stored in the provided originDestination slice at indices 0 and 1.
//
// Ported from H3 C: directedEdge.c::directedEdgeToCells.
func directedEdgeToCells(edge h3Index, originDestination []h3Index) h3Error {
	// Get the origin cell from the directed edge
	origin, originResult := getDirectedEdgeOrigin(edge)
	if originResult != eSuccess {
		return originResult
	}
	originDestination[0] = origin

	// Get the destination cell from the directed edge
	destinationResult := getDirectedEdgeDestination(edge, &originDestination[1])
	if destinationResult != eSuccess {
		return destinationResult
	}

	return eSuccess
}
