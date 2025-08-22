package c2go

import "math"

// _geoAzimuthRads determines the azimuth to p2 from p1 in radians.
// Ported from H3 C: latLng.c::_geoAzimuthRads
func _geoAzimuthRads(p1, p2 *LatLng) float64 {
	return math.Atan2(math.Cos(p2.Lat)*math.Sin(p2.Lng-p1.Lng),
		math.Cos(p1.Lat)*math.Sin(p2.Lat)-
			math.Sin(p1.Lat)*math.Cos(p2.Lat)*math.Cos(p2.Lng-p1.Lng))
}
