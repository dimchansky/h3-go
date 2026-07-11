package h3

// getDirectedEdgeOrigin extracts the origin cell from a directed edge H3 index.
//
// This function validates that the input is a directed edge index and then
// extracts the origin cell by changing the mode to h3CellMode and clearing
// the reserved bits that store the edge direction.
// Ported from H3 C: directedEdge.c::getDirectedEdgeOrigin.
func getDirectedEdgeOrigin(edge h3Index) (h3Index, h3Error) {
	if getMode(edge) != h3DirectededgeMode {
		return 0, eDirEdgeInvalid
	}

	origin := edge
	origin = setMode(origin, h3CellMode)
	origin = setReservedBits(origin, 0)

	return origin, eSuccess
}
