package h3

// destroySortableLoopSetShallow is the helper function to free the
// SortableLoopSet array without freeing vertex data. Used when vertex
// ownership has been transferred to the output GeoMultiPolygon.
// Ported from H3 C: cellsToMultiPoly.h::destroySortableLoopSetShallow.
func destroySortableLoopSetShallow(loopset *sortableLoopSet) {
	if loopset.sloops != nil {
		loopset.sloops = nil
	}
}
