package h3

// resetVisited clears the isVisited flag on every arc in the set.
// Ported from H3 C: cellsToMultiPoly.c::resetVisited.
func resetVisited(arcset arcSet) {
	for i := int64(0); i < arcset.numArcs; i++ {
		arcset.arcs[i].isVisited = false
	}
}
