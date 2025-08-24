package c2go

// cellAreaRads2 computes the area of an H3 cell in radians^2.
//
// The area is calculated by breaking the cell into spherical triangles and
// summing up their areas. Note that some H3 cells (hexagons and pentagons)
// are irregular, and have more than 6 or 5 sides.
// Ported from H3 C: latLng.c::cellAreaRads2
func cellAreaRads2(cell H3Index) (float64, H3Error) {
	var c LatLng
	var cb CellBoundary

	err := cellToLatLng(cell, &c)
	if err != E_SUCCESS {
		return 0.0, err
	}

	err = cellToBoundary(cell, &cb)
	if err != E_SUCCESS {
		// Uncoverable because cellToLatLng will have returned an error already
		return 0.0, err
	}

	area := 0.0
	for i := int32(0); i < cb.NumVerts; i++ {
		j := (i + 1) % cb.NumVerts
		area += triangleArea(&cb.Verts[i], &cb.Verts[j], &c)
	}

	return area, E_SUCCESS
}
