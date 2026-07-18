package h3

import "math"

// vec3ToCell encodes a coordinate on the sphere to the H3 index of the
// containing cell at the specified resolution.
//
// Vec3d v is expected to be on the unit sphere.
// Ported from H3 C: h3Index.c::vec3ToCell.
func vec3ToCell(v *vec3d, res int32, out *h3Index) h3Error {
	if res < 0 || res > maxH3Res {
		return eResDomain
	}
	if math.IsInf(v.X, 0) || math.IsNaN(v.X) ||
		math.IsInf(v.Y, 0) || math.IsNaN(v.Y) ||
		math.IsInf(v.Z, 0) || math.IsNaN(v.Z) {
		return eDomain
	}

	var fijk faceIJK
	_vec3ToFaceIjk(*v, res, &fijk)
	*out = _faceIjkToH3(&fijk, res)
	// ALWAYS(*out) in C - check if result is truthy
	if *out != 0 {
		return eSuccess
	} else {
		return eFailed
	}
}
