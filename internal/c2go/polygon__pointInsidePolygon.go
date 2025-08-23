package c2go

// pointInsidePolygon checks if a point is contained in a GeoPolygon.
// Ported from polygon.c::pointInsidePolygon
// bboxes must contain one bbox for outer geoloop and one per hole.
// Ported from H3 C: polygon.c::pointInsidePolygon
func pointInsidePolygon(poly GeoPolygon, bboxes []BBox, coord *LatLng) bool {
	// primary geoloop
	contains := pointInsideGeoLoop(poly.Geoloop, &bboxes[0], coord)
	if contains && len(poly.Holes) > 0 {
		for i := 0; i < len(poly.Holes); i++ {
			if pointInsideGeoLoop(poly.Holes[i], &bboxes[i+1], coord) {
				return false
			}
		}
	}
	return contains
}
