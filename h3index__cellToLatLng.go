package h3

// cellToLatLng determines the spherical coordinates of the center point of an H3 index.
// Ported from H3 C: h3Index.c::cellToLatLng.
func cellToLatLng(h3 h3Index, g *LatLng) h3Error {
	var fijk faceIJK
	e := _h3ToFaceIjk(h3, &fijk)
	if e != eSuccess {
		return e
	}
	_faceIjkToGeo(&fijk, getResolution(h3), g)
	return eSuccess
}
