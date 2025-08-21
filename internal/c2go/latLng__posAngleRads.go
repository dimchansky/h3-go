package c2go

import "math"

// _posAngleRads normalizes radians to a value between 0.0 and 2π.
// Ported from H3 C v4.3.0: latLng.c::_posAngleRads
func _posAngleRads(rads float64) float64 {
	tmp := rads
	if rads < 0.0 {
		tmp = rads + 2*math.Pi
	}
	if rads >= 2*math.Pi {
		tmp -= 2 * math.Pi
	}
	return tmp
}
