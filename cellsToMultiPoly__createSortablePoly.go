package h3

// createSortablePoly creates a SortablePolygon from a given
// SortableLoop. The "outer ring" SortableLoop is first in memory,
// followed by its holes. Later, we sort the Polygons by the size of
// their outer loops. sloop is the sub-slice of the loop set starting
// at the polygon's outer loop (C: a SortableLoop* into that array).
// Ported from H3 C: cellsToMultiPoly.c::createSortablePoly.
func createSortablePoly(sloop []sortableLoop, numHoles int64, spoly *sortablePoly) h3Error {
	var holes []GeoLoop
	if numHoles > 0 {
		holes = make([]GeoLoop, numHoles)
		for k := int64(0); k < numHoles; k++ {
			holes[k] = sloop[k+1].loop
		}
	}

	spoly.poly.GeoLoop = sloop[0].loop
	spoly.poly.Holes = holes
	spoly.outerArea = sloop[0].area // area of outer loop

	return eSuccess
}
