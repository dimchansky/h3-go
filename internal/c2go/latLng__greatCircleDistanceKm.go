package c2go

// greatCircleDistanceKm returns the great circle distance in kilometers
// between two spherical coordinates (radians).
// Ported from H3 C: latLng.c::H3_EXPORT(greatCircleDistanceKm)
func greatCircleDistanceKm(a, b LatLng) float64 {
	const earthRadiusKm = 6371.007180918475
	return greatCircleDistanceRads(a, b) * earthRadiusKm
}
