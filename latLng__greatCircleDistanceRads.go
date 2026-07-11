package h3

import (
	"math"
)

// greatCircleDistanceRads returns the great circle distance in radians
// between two spherical coordinates (radians).
// Ported from H3 C: latLng.c::H3_EXPORT(greatCircleDistanceRads).
func greatCircleDistanceRads(a, b *LatLng) float64 {
	dLat := (b.Lat - a.Lat).Mul(0.5)
	dLng := (b.Lng - a.Lng).Mul(0.5)
	sinLat := dLat.Sin()
	sinLng := dLng.Sin()
	A := sinLat*sinLat + a.Lat.Cos()*b.Lat.Cos()*sinLng*sinLng
	return 2 * math.Atan2(math.Sqrt(A), math.Sqrt(1-A))
}
