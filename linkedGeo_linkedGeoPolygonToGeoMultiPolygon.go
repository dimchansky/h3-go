package h3

// linkedGeoPolygonToGeoMultiPolygon converts a LinkedGeoPolygon to a
// GeoMultiPolygon.
//
// An empty chain (head node with no loops and no `next`) produces an
// empty output (numPolygons=0) and returns eSuccess. Every non-empty
// polygon node must have an outer loop, and every loop must have >= 3
// vertices; otherwise, eFailed is returned.
//
// On error, any (partial) output is cleaned up via
// destroyGeoMultiPolygon(). On success the caller owns the output (in
// Go the garbage collector frees it).
// Ported from H3 C: linkedGeo.c::linkedGeoPolygonToGeoMultiPolygon.
func linkedGeoPolygonToGeoMultiPolygon(linked *linkedGeoPolygon, out *geoMultiPolygon) h3Error {
	out.NumPolygons = 0
	out.Polygons = nil

	// Empty chain (head has no loops and no next) is valid: 0 polygons
	if linked.First == nil && linked.Next == nil {
		return eSuccess
	}

	numPolygons := countLinkedPolygons(linked)

	polygons := make([]GeoPolygon, numPolygons)

	out.Polygons = polygons
	out.NumPolygons = numPolygons

	lpoly := linked
	for i := int32(0); lpoly != nil && i < numPolygons; i++ {
		// ALWAYS(i < numPolygons) in C.
		if lpoly.First == nil {
			destroyGeoMultiPolygon(out)
			return eFailed
		}
		err := linkedGeoPolygonToGeoPolygon(lpoly, &polygons[i])
		if err != eSuccess {
			destroyGeoMultiPolygon(out)
			return err
		}
		lpoly = lpoly.Next
	}

	return eSuccess
}
