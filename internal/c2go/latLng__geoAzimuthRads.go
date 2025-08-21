package c2go

import "math"

// _geoAzimuthRads determines the azimuth to p2 from p1 in radians.
// Ported from H3 C: latLng.c::_geoAzimuthRads
func _geoAzimuthRads(p1, p2 LatLng) float64 {
    return math.Atan2(math.Cos(p2.lat)*math.Sin(p2.lng-p1.lng),
        math.Cos(p1.lat)*math.Sin(p2.lat)-
            math.Sin(p1.lat)*math.Cos(p2.lat)*math.Cos(p2.lng-p1.lng))
}

