package h3

// isValidDirectedEdge determines if the provided h3Index is a valid directed edge index.
//
// This function validates a directed edge by checking:
// 1. direction is valid (not center digit, not invalid digit)
// 2. Origin cell is valid when extracted from the edge
// 3. Pentagon cells don't have K-axis directed edges (deleted subsequence)
//
// Ported from H3 C: directedEdge.c::isValidDirectedEdge.
func isValidDirectedEdge(edge h3Index) bool {
	neighborDirection := direction(getReservedBits(edge))
	if neighborDirection <= centerDigit || neighborDirection >= numDigits {
		return false
	}

	origin, originErr := getDirectedEdgeOrigin(edge)
	if originErr != eSuccess {
		return false
	}

	if isPentagon(origin) && neighborDirection == kAxesDigit {
		return false
	}

	return isValidCell(origin)
}
