package h3

// _faceIjkToVec3 determines the center point in 3D coordinates of a cell
// given by a FaceIJK address at a specified resolution.
// Ported from H3 C: faceijk.c::_faceIjkToVec3.
func _faceIjkToVec3(h *faceIJK, res int32, g *vec3d) {
	var v vec2d
	_ijkToHex2d(&h.Coord, &v)
	_hex2dToVec3(&v, h.Face, res, 0, g)
}
