package h3

// countPolys counts the number of polygons in a sorted SortableLoopSet:
// one per distinct connected-component root.
// Ported from H3 C: cellsToMultiPoly.c::countPolys.
func countPolys(loopset sortableLoopSet) int64 {
	numPolys := int64(0)

	cur := h3Null
	for i := int64(0); i < loopset.numLoops; i++ {
		if loopset.sloops[i].root != cur {
			numPolys++
			cur = loopset.sloops[i].root
		}
	}

	return numPolys
}
