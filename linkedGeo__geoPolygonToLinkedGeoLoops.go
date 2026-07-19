package h3

// geoPolygonToLinkedGeoLoops converts a single GeoPolygon to
// LinkedGeoLoops within a LinkedGeoPolygon.
// Ported from H3 C: linkedGeo.c::geoPolygonToLinkedGeoLoops.
func geoPolygonToLinkedGeoLoops(poly *GeoPolygon, currentPoly *linkedGeoPolygon) h3Error {
	err := addLinkedGeoLoop(poly.GeoLoop, currentPoly)
	if err != eSuccess {
		return err
	}

	for i := 0; i < len(poly.Holes); i++ {
		err = addLinkedGeoLoop(poly.Holes[i], currentPoly)
		if err != eSuccess {
			return err
		}
	}

	return eSuccess
}
