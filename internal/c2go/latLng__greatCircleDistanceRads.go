package c2go

import "math"

// greatCircleDistanceRads returns the great circle distance in radians
// between two spherical coordinates (radians).
// Ported from H3 C: latLng.c::H3_EXPORT(greatCircleDistanceRads)
func greatCircleDistanceRads(a, b LatLng) float64 {
    sinLat := math.Sin((b.lat - a.lat) * 0.5)
    sinLng := math.Sin((b.lng - a.lng) * 0.5)
    A := sinLat*sinLat + math.Cos(a.lat)*math.Cos(b.lat)*sinLng*sinLng
    return 2 * math.Atan2(math.Sqrt(A), math.Sqrt(1-A))
}

