package h3

// geoAlmostEqualThreshold determines if two spherical coordinates are within
// the given threshold distance (in radians) of each other.
// Ported from H3 C: latLng.c::geoAlmostEqualThreshold
func geoAlmostEqualThreshold(p1, p2 *LatLng, threshold float64) bool {
	return (p1.Lat-p2.Lat).Abs().Rad() < threshold && (p1.Lng-p2.Lng).Abs().Rad() < threshold
}
