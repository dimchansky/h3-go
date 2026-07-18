package h3

// cellToLatLng determines the spherical coordinates of the center point of an H3 index.
// Ported from H3 C: h3Index.c::cellToLatLng.
func cellToLatLng(h3 h3Index, g *LatLng) h3Error {
	var v vec3d
	e := cellToVec3(h3, &v)
	if e != eSuccess {
		return e
	}
	*g = vec3ToLatLng(v)
	return eSuccess
}
