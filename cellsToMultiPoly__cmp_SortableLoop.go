package h3

// cmp_SortableLoop orders loops by (connected component root asc, loop
// area asc); used to make each polygon's loops contiguous with its
// outer loop (smallest enclosed area) first.
// Ported from H3 C: cellsToMultiPoly.h::cmp_SortableLoop.
func cmp_SortableLoop(a, b *sortableLoop) int32 {
	// first, sort on connected component
	if a.root < b.root {
		return -1
	}
	if a.root > b.root {
		return 1
	}

	// second, sort on area of loops
	if a.area < b.area {
		return -1
	}
	if a.area > b.area {
		return 1
	}

	return 0 // same root and equal area
}
