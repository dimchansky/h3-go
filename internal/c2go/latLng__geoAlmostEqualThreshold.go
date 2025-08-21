package c2go

import "math"

// geoAlmostEqualThreshold determines if two spherical coordinates are within
// the given threshold distance (in radians) of each other.
// Ported from H3 C: latLng.c::geoAlmostEqualThreshold
func geoAlmostEqualThreshold(p1, p2 LatLng, threshold float64) bool {
	return math.Abs(p1.Lat-p2.Lat) < threshold && math.Abs(p1.Lng-p2.Lng) < threshold
}
