package c2go

import "math"

// greatCircleDistanceRads returns the great circle distance in radians
// between two spherical coordinates (radians).
// Ported from H3 C: latLng.c::H3_EXPORT(greatCircleDistanceRads)
func greatCircleDistanceRads(a, b LatLng) float64 {
	sinLat := math.Sin((b.Lat - a.Lat) * 0.5)
	sinLng := math.Sin((b.Lng - a.Lng) * 0.5)
	A := sinLat*sinLat + math.Cos(a.Lat)*math.Cos(b.Lat)*sinLng*sinLng
	return 2 * math.Atan2(math.Sqrt(A), math.Sqrt(1-A))
}
