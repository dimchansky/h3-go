package h3

// maxPolygonToCellsSize returns the number of cells to allocate space for
// when performing a polygonToCells on the given GeoJSON-like data structure.
//
// The size is the maximum of either the number of points in the geoloop or the
// number of cells in the bounding box of the geoloop.
// Ported from H3 C: algos.c::maxPolygonToCellsSize
func maxPolygonToCellsSize(geoPolygon *GeoPolygon, res int32, flags uint32, out *int64) H3Error {
	flagErr := validatePolygonFlags(flags)
	if flagErr != E_SUCCESS {
		return flagErr
	}

	// Get the bounding box for the GeoJSON-like struct
	var bbox BBox
	geoloop := geoPolygon.Geoloop
	bboxFromGeoLoop(geoloop, &bbox)

	var numHexagons int64
	estimateErr := bboxHexEstimate(&bbox, res, &numHexagons)
	if estimateErr != E_SUCCESS {
		return estimateErr
	}

	// This algorithm assumes that the number of vertices is usually less than
	// the number of hexagons, but when it's wrong, this will keep it from
	// failing
	totalVerts := int64(len(geoloop))
	for i := 0; i < len(geoPolygon.Holes); i++ {
		totalVerts += int64(len(geoPolygon.Holes[i]))
	}
	if numHexagons < totalVerts {
		numHexagons = totalVerts
	}

	// When the polygon is very small, near an icosahedron edge and is an odd
	// resolution, the line tracing needs an extra buffer than the estimator
	// function provides (but beefing that up to cover causes most situations to
	// overallocate memory)
	numHexagons += POLYGON_TO_CELLS_BUFFER
	*out = numHexagons
	return E_SUCCESS
}
