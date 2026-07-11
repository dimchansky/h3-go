//go:build cgo && c2go

package h3

/*
#include <stdint.h>
#include <stdbool.h>
#include "vec2d.h"

static int h3_bool_to_int_vec2d(_Bool b) { return b ? 1 : 0; }
*/
import "C"

// v2dMagC calls the original C implementation for magnitude.
func v2dMagC(v vec2d) float64 {
	var cv C.Vec2d
	cv.x = C.double(v.X)
	cv.y = C.double(v.Y)
	return float64(C._v2dMag(&cv))
}

// v2dIntersectC calls the original C implementation for line intersection.
func v2dIntersectC(p0, p1, p2, p3 vec2d) vec2d {
	var c0, c1, c2, c3, ci C.Vec2d
	c0.x = C.double(p0.X)
	c0.y = C.double(p0.Y)
	c1.x = C.double(p1.X)
	c1.y = C.double(p1.Y)
	c2.x = C.double(p2.X)
	c2.y = C.double(p2.Y)
	c3.x = C.double(p3.X)
	c3.y = C.double(p3.Y)
	C._v2dIntersect(&c0, &c1, &c2, &c3, &ci)
	return vec2d{X: float64(ci.x), Y: float64(ci.y)}
}

// v2dAlmostEqualsC calls the original C implementation for equality check.
func v2dAlmostEqualsC(v1, v2 vec2d) bool {
	var c1, c2 C.Vec2d
	c1.x = C.double(v1.X)
	c1.y = C.double(v1.Y)
	c2.x = C.double(v2.X)
	c2.y = C.double(v2.Y)
	return C.h3_bool_to_int_vec2d(C._v2dAlmostEquals(&c1, &c2)) != 0
}
