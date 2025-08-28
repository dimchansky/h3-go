package h3

import "math"

// _v2dAlmostEquals compares two 2D vectors within FLT_EPSILON threshold.
// Ported from H3 C: vec2d.c::_v2dAlmostEquals (uses float epsilon)
func _v2dAlmostEquals(v1, v2 *Vec2d) bool {
	const fltEpsilon = 1.19209290e-07 // FLT_EPSILON
	return math.Abs(v1.X-v2.X) < fltEpsilon && math.Abs(v1.Y-v2.Y) < fltEpsilon
}
