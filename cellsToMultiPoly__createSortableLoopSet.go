package h3

import "sort"

// createSortableLoopSet creates the set of all SortableLoops and sorts
// them. The comparison function makes all polygons (outer loop and
// holes) contiguous in memory, so that the outer loop "contains" the
// holes.
// Ported from H3 C: cellsToMultiPoly.c::createSortableLoopSet.
func createSortableLoopSet(arcset arcSet, loopset *sortableLoopSet) h3Error {
	numLoops := countLoops(arcset)
	resetVisited(arcset)
	arcs := arcset.arcs

	sloops := make([]sortableLoop, numLoops)

	j := int64(0)
	for i := int64(0); i < arcset.numArcs; i++ {
		if !arcs[i].isVisited && !arcs[i].isRemoved {
			err := createSortableLoop(&arcs[i], &sloops[j])
			if err != eSuccess {
				// Free any verts already allocated in previous loops
				partialLoopSet := sortableLoopSet{numLoops: j, sloops: sloops}
				destroySortableLoopSet(&partialLoopSet)
				return err
			}
			j++
		}
	}

	// The comparison function makes all polygons (outer loop and holes)
	// contiguous in memory, so that the outer loop "contains" the holes.
	sort.Slice(sloops, func(a, b int) bool {
		return cmp_SortableLoop(&sloops[a], &sloops[b]) < 0
	})

	loopset.numLoops = numLoops
	loopset.sloops = sloops

	return eSuccess
}
