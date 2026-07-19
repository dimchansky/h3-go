package h3

// getRoot is part of the union-find data structure. It finds the id of
// the connected component this arc/edge is a part of, compressing the
// path as it recurses.
// Ported from H3 C: cellsToMultiPoly.c::getRoot.
func getRoot(arc *arc) *arc {
	parent := arc.parent

	if parent == arc {
		return arc
	}
	parent = getRoot(parent)
	arc.parent = parent
	return parent
}
