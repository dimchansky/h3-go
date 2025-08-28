package h3

// _cellsToDirectedEdge returns a directed edge H3 index based on the provided origin and destination.
// Creates a directed edge from the origin cell to the destination cell by determining
// the direction and encoding it in the reserved bits with H3_DIRECTEDEDGE_MODE.
// Ported from H3 C: directedEdge.c::cellsToDirectedEdge
func _cellsToDirectedEdge(origin H3Index, destination H3Index) (H3Index, H3Error) {
	// Determine the IJK direction from the origin to the destination
	direction := _directionForNeighbor(origin, destination)

	// The direction will be invalid if the cells are not neighbors
	if direction == INVALID_DIGIT {
		return 0, E_NOT_NEIGHBORS
	}

	// Create the edge index for the neighbor direction
	output := origin
	output = setMode(output, H3_DIRECTEDEDGE_MODE)
	output = setReservedBits(output, int32(direction))

	return output, E_SUCCESS
}
