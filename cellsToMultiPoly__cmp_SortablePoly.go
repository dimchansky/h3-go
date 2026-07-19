package h3

// cmp_SortablePoly orders polygons by area of their outer loop, in
// descending order.
// Ported from H3 C: cellsToMultiPoly.h::cmp_SortablePoly.
func cmp_SortablePoly(a, b *sortablePoly) int32 {
	// Sort by area of outer loop, in descending order
	if a.outerArea > b.outerArea {
		return -1
	}
	if a.outerArea < b.outerArea {
		return 1
	}

	return 0 // equal area
}
