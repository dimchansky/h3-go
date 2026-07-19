package h3

// _vec3ToFaceIjk encodes a Vec3d coordinate to the FaceIJK address of the
// containing cell at the specified resolution.
//
// Vec3d p is expected to be on the unit sphere.
// Ported from H3 C: faceijk.c::_vec3ToFaceIjk.
func _vec3ToFaceIjk(p vec3d, res int32, h *faceIJK) {
	// first convert to hex2d
	var v vec2d
	_vec3ToHex2d(&p, res, &h.Face, &v)

	// then convert to ijk+
	_hex2dToCoordIJK(&v, &h.Coord)
}
