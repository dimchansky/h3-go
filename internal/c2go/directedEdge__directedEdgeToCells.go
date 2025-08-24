package c2go

// directedEdgeToCells returns the origin and destination cell pair for a directed edge.
//
// This function extracts both the origin and destination cells from a directed edge
// H3 index by calling getDirectedEdgeOrigin and getDirectedEdgeDestination sequentially.
// The results are stored in the provided originDestination slice at indices 0 and 1.
//
// Ported from H3 C: directedEdge.c::directedEdgeToCells
func directedEdgeToCells(edge H3Index, originDestination []H3Index) H3Error {
	// Get the origin cell from the directed edge
	origin, originResult := getDirectedEdgeOrigin(edge)
	if originResult != E_SUCCESS {
		return originResult
	}
	originDestination[0] = origin

	// Get the destination cell from the directed edge
	destinationResult := getDirectedEdgeDestination(edge, &originDestination[1])
	if destinationResult != E_SUCCESS {
		return destinationResult
	}

	return E_SUCCESS
}
