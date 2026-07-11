package h3

// isValidIndex returns whether or not an H3 index is valid for any mode
// (cell, directed edge, or vertex).
// Ported from H3 C: h3Index.c::isValidIndex.
func isValidIndex(h h3Index) bool {
	return isValidCell(h) || isValidDirectedEdge(h) || isValidVertex(h)
}
