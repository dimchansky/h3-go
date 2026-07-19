package h3

// destroySortableLoopSet is the helper function to free memory
// allocated for a SortableLoopSet. Frees all vertex arrays in the
// loops, then the loops array itself (in Go: nils the references; the
// garbage collector frees the memory).
// Ported from H3 C: cellsToMultiPoly.h::destroySortableLoopSet.
func destroySortableLoopSet(loopset *sortableLoopSet) {
	if loopset.sloops != nil {
		for i := int64(0); i < loopset.numLoops; i++ {
			if loopset.sloops[i].loop != nil {
				loopset.sloops[i].loop = nil
			}
		}
	}
	destroySortableLoopSetShallow(loopset)
}
