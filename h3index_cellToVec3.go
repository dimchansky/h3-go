package h3

// cellToVec3 determines the 3D cartesian coordinates of the center of an
// H3 cell.
// Ported from H3 C: h3Index.c::cellToVec3.
func cellToVec3(h3 h3Index, v *vec3d) h3Error {
	var fijk faceIJK
	e := _h3ToFaceIjk(h3, &fijk)
	if e != eSuccess {
		return e
	}
	_faceIjkToVec3(&fijk, getResolution(h3), v)
	return eSuccess
}
