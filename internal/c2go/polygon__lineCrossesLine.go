package c2go

// lineCrossesLine checks for cartesian line segment intersection.
// Ported from H3 C: polygon.c::lineCrossesLine
func lineCrossesLine(a1, a2, b1, b2 *LatLng) bool {
	denom := ((b2.Lng - b1.Lng) * (a2.Lat - a1.Lat)) - ((b2.Lat - b1.Lat) * (a2.Lng - a1.Lng))
	if denom == 0 {
		return false
	}
	test := ((b2.Lat-b1.Lat)*(a1.Lng-b1.Lng) - (b2.Lng-b1.Lng)*(a1.Lat-b1.Lat)) / denom
	if test < 0 || test > 1 {
		return false
	}
	test = ((a2.Lat-a1.Lat)*(a1.Lng-b1.Lng) - (a2.Lng-a1.Lng)*(a1.Lat-b1.Lat)) / denom
	if test < 0 || test > 1 {
		return false
	}
	return true
}
