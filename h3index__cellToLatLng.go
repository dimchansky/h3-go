package h3

// cellToLatLng determines the spherical coordinates of the center point of an H3 index.
// Ported from H3 C: h3Index.c::cellToLatLng.
func cellToLatLng(h3 H3Index, g *LatLng) H3Error {
	var fijk FaceIJK
	e := _h3ToFaceIjk(h3, &fijk)
	if e != E_SUCCESS {
		return e
	}
	_faceIjkToGeo(&fijk, getResolution(h3), g)
	return E_SUCCESS
}
