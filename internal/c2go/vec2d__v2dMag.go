package c2go

import "math"

// _v2dMag calculates the magnitude of a 2D cartesian vector.
// Ported from H3 C: vec2d.c::_v2dMag
func _v2dMag(v Vec2d) float64 { return math.Sqrt(v.X*v.X + v.Y*v.Y) }

