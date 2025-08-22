package c2go

import "math"

// triangleEdgeLengthsToArea computes spherical triangle area on unit sphere.
// Ported from H3 C: latLng.c::triangleEdgeLengthsToArea
func triangleEdgeLengthsToArea(a, b, c float64) float64 {
	s := (a + b + c) * 0.5
	a = (s - a) * 0.5
	b = (s - b) * 0.5
	c = (s - c) * 0.5
	s = s * 0.5
	return 4 * math.Atan(math.Sqrt(math.Tan(s)*math.Tan(a)*math.Tan(b)*math.Tan(c)))
}
