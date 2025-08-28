package h3

import "math"

// geoAlmostEqual determines if two spherical coordinates are within the
// standard epsilon distance of each other.
// Ported from H3 C: latLng.c::geoAlmostEqual
func geoAlmostEqual(p1, p2 *LatLng) bool {
	const epsilonDeg = 1e-9
	const epsilonRad = epsilonDeg * (math.Pi / 180.0)
	return geoAlmostEqualThreshold(p1, p2, epsilonRad)
}
