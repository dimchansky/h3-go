package h3

//lint:file-ignore U1000 Ignore all unused code
// 2D floating point vector functions.

import (
	"math"
)

// Vec2d is 2D floating-point vector
type Vec2d struct {
	X float64 // X component
	Y float64 // Y component
}

// v2dMag calculates the magnitude of a 2D cartesian vector.
// `v` is the 2D cartesian vector.
// Returns the magnitude of the vector.
func v2dMag(v *Vec2d) float64 { return math.Sqrt(square(v.X) + square(v.Y)) }

// v2dIntersect finds the intersection between two lines. Assumes that the lines intersect
// and that the intersection is not at an endpoint of either line.
// `p0` is the first endpoint of the first line.
// `p1` is the second endpoint of the first line.
// `p2` is the first endpoint of the second line.
// `p3` is the second endpoint of the second line.
// `inter` is the intersection point.
func v2dIntersect(p0 *Vec2d, p1 *Vec2d, p2 *Vec2d, p3 *Vec2d, inter *Vec2d) {
	var s1, s2 Vec2d
	s1.X = p1.X - p0.X
	s1.Y = p1.Y - p0.Y
	s2.X = p3.X - p2.X
	s2.Y = p3.Y - p2.Y

	t := (s2.X*(p0.Y-p2.Y) - s2.Y*(p0.X-p2.X)) /
		(-s2.X*s1.Y + s1.X*s2.Y)

	inter.X = p0.X + (t * s1.X)
	inter.Y = p0.Y + (t * s1.Y)
}

// v2dEquals checks Whether two 2D vectors are equal. Does not consider possible false
// negatives due to floating-point errors.
// `v1` is the first vector to compare
// `v2` is the second vector to compare
// Returns whether the vectors are equal
func v2dEquals(v1 *Vec2d, v2 *Vec2d) bool {
	return v1.X == v2.X && v1.Y == v2.Y
}
