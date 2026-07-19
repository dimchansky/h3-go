package h3

// countLoops counts the number of distinct loops in an ArcSet.
// Ported from H3 C: cellsToMultiPoly.c::countLoops.
func countLoops(arcset arcSet) int64 {
	arcs := arcset.arcs
	resetVisited(arcset)
	numLoops := int64(0)

	for i := int64(0); i < arcset.numArcs; i++ {
		arc := &arcs[i]
		if !arc.isVisited && !arc.isRemoved {
			numLoops++
			start := arc.id

			for {
				arc.isVisited = true
				arc = arc.next
				if arc.id == start {
					break
				}
			}
		}
	}

	return numLoops
}
