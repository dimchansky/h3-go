package h3

// cellsToDirectedEdge returns a directed edge H3 index based on the provided origin and destination.
// Creates a directed edge from the origin cell to the destination cell by determining
// the direction and encoding it in the reserved bits with h3DirectededgeMode.
// Ported from H3 C: directedEdge.c::cellsToDirectedEdge.
func cellsToDirectedEdge(origin h3Index, destination h3Index) (h3Index, h3Error) {
	// Determine the IJK direction from the origin to the destination
	direction := _directionForNeighbor(origin, destination)

	// The direction will be invalid if the cells are not neighbors
	if direction == invalidDigit {
		return 0, eNotNeighbors
	}

	// Create the edge index for the neighbor direction
	output := origin
	output = setMode(output, h3DirectededgeMode)
	output = setReservedBits(output, int32(direction))

	return output, eSuccess
}
