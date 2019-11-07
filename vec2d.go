package h3

//lint:file-ignore U1000 Ignore all unused code
// 2D floating point vector functions.

import (
	"math"
)

// Vec2d is 2D floating-point vector
type Vec2d struct {
	x float64 // x component
	y float64 // y component
}

// _v2dMag calculates the magnitude of a 2D cartesian vector.
// `v`: The 2D cartesian vector.
// Returns the magnitude of the vector.
func _v2dMag(v *Vec2d) float64 { return math.Sqrt(square(v.x) + square(v.y)) }

// _v2dIntersect finds the intersection between two lines. Assumes that the lines intersect
// and that the intersection is not at an endpoint of either line.
// `p0`: The first endpoint of the first line.
// `p1`: The second endpoint of the first line.
// `p2`: The first endpoint of the second line.
// `p3`: The second endpoint of the second line.
// `inter`: The intersection point.
func _v2dIntersect(p0 *Vec2d, p1 *Vec2d, p2 *Vec2d, p3 *Vec2d, inter *Vec2d) {
	var s1, s2 Vec2d
	s1.x = p1.x - p0.x
	s1.y = p1.y - p0.y
	s2.x = p3.x - p2.x
	s2.y = p3.y - p2.y

	t := (s2.x*(p0.y-p2.y) - s2.y*(p0.x-p2.x)) /
		(-s2.x*s1.y + s1.x*s2.y)

	inter.x = p0.x + (t * s1.x)
	inter.y = p0.y + (t * s1.y)
}

// _v2dEquals checks Whether two 2D vectors are equal. Does not consider possible false
// negatives due to floating-point errors.
// `v1`: The first vector to compare
// `v2`: The second vector to compare
// Returns whether the vectors are equal
func _v2dEquals(v1 *Vec2d, v2 *Vec2d) bool {
	return v1.x == v2.x && v1.y == v2.y
}
