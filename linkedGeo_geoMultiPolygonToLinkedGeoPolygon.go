package h3

// geoMultiPolygonToLinkedGeoPolygon converts a GeoMultiPolygon to a
// LinkedGeoPolygon.
//
// The first polygon is placed in the caller-owned `out` node. Every
// loop must have >= 3 vertices; otherwise eFailed is returned.
//
// On error, the output is cleaned up via destroyLinkedMultiPolygon().
// On success, the caller owns the output (in Go the garbage collector
// frees it).
// Ported from H3 C: linkedGeo.c::geoMultiPolygonToLinkedGeoPolygon.
func geoMultiPolygonToLinkedGeoPolygon(mpoly *geoMultiPolygon, out *linkedGeoPolygon) h3Error {
	*out = linkedGeoPolygon{}

	currentPoly := out
	for i := int32(0); i < mpoly.NumPolygons; i++ {
		if i > 0 {
			newPoly := &linkedGeoPolygon{}
			currentPoly.Next = newPoly
			currentPoly = newPoly
		}

		err := geoPolygonToLinkedGeoLoops(&mpoly.Polygons[i], currentPoly)
		if err != eSuccess {
			destroyLinkedMultiPolygon(out)
			return err
		}
	}

	return eSuccess
}
