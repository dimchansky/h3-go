// Package v2d provides a minimal, C-like 2D vector API focused on performance.
package v2d

import "math"

// Float32Epsilon is the IEEE-754 single-precision machine epsilon.
// It numerically equals C's FLT_EPSILON (≈1.19209290e-07).
const Float32Epsilon = 1.19209290e-07

// Vec2d represents a 2D floating-point vector (double precision).
type Vec2d struct {
	X float64 // x component
	Y float64 // y component
}

// Mag returns the magnitude (length) of the vector.
// Matches the C version: sqrt(x*x + y*y) for maximum speed.
func (v Vec2d) Mag() float64 { return math.Sqrt(v.X*v.X + v.Y*v.Y) }

// Len2 returns the squared magnitude (avoids sqrt; useful in hot paths).
func (v Vec2d) Len2() float64 { return v.X*v.X + v.Y*v.Y }

// AlmostEqual reports whether two vectors are almost equal using
// Float32Epsilon as a strict absolute threshold (FLT_EPSILON equivalent).
func (v Vec2d) AlmostEqual(u Vec2d) bool {
	return abs64(v.X-u.X) < Float32Epsilon &&
		abs64(v.Y-u.Y) < Float32Epsilon
}

// Intersect computes the intersection point of two infinite lines (p0–p1) and (p2–p3).
// Preconditions: lines do intersect; intersection is not at endpoints; lines are not parallel.
// No checks are performed (behavior mirrors the original C code).
func Intersect(p0, p1, p2, p3 Vec2d) Vec2d {
	s1x := p1.X - p0.X
	s1y := p1.Y - p0.Y
	s2x := p3.X - p2.X
	s2y := p3.Y - p2.Y

	den := -s2x*s1y + s1x*s2y // no zero-division checks (as in C)
	t := (s2x*(p0.Y-p2.Y) - s2y*(p0.X-p2.X)) / den

	return Vec2d{X: p0.X + t*s1x, Y: p0.Y + t*s1y}
}

// abs64 is an inlined absolute value helper for float64.
func abs64(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
