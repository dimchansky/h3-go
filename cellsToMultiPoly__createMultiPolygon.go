package h3

import "sort"

// createMultiPolygon groups the sorted loop set into polygons (each
// distinct connected-component root's loops are contiguous, outer loop
// first) and sorts the polygons by their outer loop area, descending.
// For example, in a multipolygon representing the USA, the continental
// US will come before any of the Hawaiian islands. If the loop set is
// empty (the input tiles the globe), emits the 8-octant globe
// multipolygon.
// Ported from H3 C: cellsToMultiPoly.c::createMultiPolygon.
func createMultiPolygon(loopset sortableLoopSet, mpoly *geoMultiPolygon) h3Error {
	if loopset.numLoops == 0 {
		return createGlobeMultiPolygon(mpoly)
	}

	numPolys := countPolys(loopset)
	spolys := make([]sortablePoly, numPolys)

	sloop := loopset.sloops
	i := int64(0) // index of first loop in polygon (outer loop)
	p := int64(0) // index of polygon we're working on
	// j: index + 1 of last loop in polygon (last hole + 1)
	for j := int64(0); j <= loopset.numLoops; j++ {
		if j == loopset.numLoops || sloop[i].root != sloop[j].root {
			// We've reached the end of the loops in the polygon, so
			// now construct a polygon from the start of those loops.
			err := createSortablePoly(sloop[i:], (j-i)-1, &spolys[p])
			if err != eSuccess {
				destroySortablePolys(spolys, p)
				return err
			}
			p++
			i = j
		}
	}

	// Sort polygons by their outer loop area. For example, in a multipolygon
	// representing the USA, the continental US will come before any of the
	// Hawaiian islands
	sort.Slice(spolys, func(a, b int) bool {
		return cmp_SortablePoly(&spolys[a], &spolys[b]) < 0
	})

	mpoly.Polygons = make([]GeoPolygon, numPolys)

	mpoly.NumPolygons = int32(numPolys)
	for i := int64(0); i < numPolys; i++ {
		mpoly.Polygons[i] = spolys[i].poly
	}

	return eSuccess
}
