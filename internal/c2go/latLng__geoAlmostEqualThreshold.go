package c2go

import "math"

// geoAlmostEqualThreshold determines if two spherical coordinates are within
// the given threshold distance (in radians) of each other.
// Ported from H3 C: latLng.c::geoAlmostEqualThreshold
func geoAlmostEqualThreshold(p1, p2 LatLng, threshold float64) bool {
    return math.Abs(p1.lat-p2.lat) < threshold && math.Abs(p1.lng-p2.lng) < threshold
}

