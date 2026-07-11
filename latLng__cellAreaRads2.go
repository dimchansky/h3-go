package h3

// cellAreaRads2 computes the area of an H3 cell in radians^2.
//
// The area is calculated by breaking the cell into spherical triangles and
// summing up their areas. Note that some H3 cells (hexagons and pentagons)
// are irregular, and have more than 6 or 5 sides.
// Ported from H3 C: latLng.c::cellAreaRads2.
func cellAreaRads2(cell h3Index) (float64, h3Error) {
	var c LatLng
	var cb CellBoundary

	err := cellToLatLng(cell, &c)
	if err != eSuccess {
		return 0.0, err
	}

	err = cellToBoundary(cell, &cb)
	if err != eSuccess {
		// Uncoverable because cellToLatLng will have returned an error already
		return 0.0, err
	}

	area := 0.0
	for i := int32(0); i < cb.numVerts; i++ {
		j := (i + 1) % cb.numVerts
		area += triangleArea(&cb.verts[i], &cb.verts[j], &c)
	}

	return area, eSuccess
}
