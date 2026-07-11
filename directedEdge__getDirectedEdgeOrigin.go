package h3

// getDirectedEdgeOrigin extracts the origin cell from a directed edge H3 index.
//
// This function validates that the input is a directed edge index and then
// extracts the origin cell by changing the mode to H3_CELL_MODE and clearing
// the reserved bits that store the edge direction.
// Ported from H3 C: directedEdge.c::getDirectedEdgeOrigin.
func getDirectedEdgeOrigin(edge H3Index) (H3Index, H3Error) {
	if getMode(edge) != H3_DIRECTEDEDGE_MODE {
		return 0, E_DIR_EDGE_INVALID
	}

	origin := edge
	origin = setMode(origin, H3_CELL_MODE)
	origin = setReservedBits(origin, 0)

	return origin, E_SUCCESS
}
