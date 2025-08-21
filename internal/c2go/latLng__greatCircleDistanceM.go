package c2go

// greatCircleDistanceM returns the great circle distance in meters
// between two spherical coordinates (radians).
// Ported from H3 C: latLng.c::H3_EXPORT(greatCircleDistanceM)
func greatCircleDistanceM(a, b LatLng) float64 {
    return greatCircleDistanceKm(a, b) * 1000.0
}

