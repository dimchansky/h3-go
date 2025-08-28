package h3

import (
	"math"
)

// _geoAzimuthRads determines the azimuth to p2 from p1 in radians.
// Ported from H3 C: latLng.c::_geoAzimuthRads
func _geoAzimuthRads(p1, p2 *LatLng) float64 {
	// Use Angle operations and compute sin/cos efficiently with SinCos.
	dl := p2.Lng - p1.Lng
	sinDl, cosDl := dl.SinCos()

	sinLat1, cosLat1 := p1.Lat.SinCos()
	sinLat2, cosLat2 := p2.Lat.SinCos()

	return math.Atan2(cosLat2*sinDl, cosLat1*sinLat2-sinLat1*cosLat2*cosDl)
}
