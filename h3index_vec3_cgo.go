//go:build cgo && c2go

package h3

// Bridges to the internal cellToVec3/vec3ToCell added to h3Index.c in
// H3 4.5.0 (declared in h3Index.h there; absent from the 4.4.0 tree).

/*
#include "h3api.h"
#include "h3Index.h"
#include "vec3d.h"
*/
import "C"

// cellToVec3C calls the original C implementation.
func cellToVec3C(h h3Index) (vec3d, h3Error) {
	var cv C.Vec3d
	err := h3Error(C.cellToVec3(C.H3Index(h), &cv))
	return vec3d{X: float64(cv.x), Y: float64(cv.y), Z: float64(cv.z)}, err
}

// vec3ToCellC calls the original C implementation.
func vec3ToCellC(v vec3d, res int32) (h3Index, h3Error) {
	var cv C.Vec3d
	cv.x = C.double(v.X)
	cv.y = C.double(v.Y)
	cv.z = C.double(v.Z)
	var out C.H3Index
	err := h3Error(C.vec3ToCell(&cv, C.int(res), &out))
	return h3Index(out), err
}
