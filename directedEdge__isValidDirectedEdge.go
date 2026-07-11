package h3

// isValidDirectedEdge determines if the provided H3Index is a valid directed edge index.
//
// This function validates a directed edge by checking:
// 1. Direction is valid (not center digit, not invalid digit)
// 2. Origin cell is valid when extracted from the edge
// 3. Pentagon cells don't have K-axis directed edges (deleted subsequence)
//
// Ported from H3 C: directedEdge.c::isValidDirectedEdge.
func isValidDirectedEdge(edge H3Index) bool {
	neighborDirection := Direction(getReservedBits(edge))
	if neighborDirection <= CENTER_DIGIT || neighborDirection >= NUM_DIGITS {
		return false
	}

	origin, originErr := getDirectedEdgeOrigin(edge)
	if originErr != E_SUCCESS {
		return false
	}

	if isPentagon(origin) && neighborDirection == K_AXES_DIGIT {
		return false
	}

	return isValidCell(origin)
}
