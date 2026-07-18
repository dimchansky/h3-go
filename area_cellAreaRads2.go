package h3

// cellAreaRads2 computes the area of an H3 cell in radians^2 using
// geoLoopAreaRads2 over the cell boundary.
// Ported from H3 C: area.c::cellAreaRads2.
func cellAreaRads2(cell h3Index) (float64, h3Error) {
	var cb CellBoundary
	err := cellToBoundary(cell, &cb)
	if err != eSuccess {
		return 0, err
	}

	loop := GeoLoop(cb.verts[:cb.numVerts])
	out, err := geoLoopAreaRads2(loop)
	// NEVER(err) in C - geoLoopAreaRads2 cannot fail
	if err != eSuccess {
		return 0, err
	}

	return out, eSuccess
}
