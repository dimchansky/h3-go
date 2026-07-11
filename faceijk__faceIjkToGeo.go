package h3

// _faceIjkToGeo converts faceIJK coordinates to geographic coordinates.
// Ported from H3 C: faceijk.c::_faceIjkToGeo.
func _faceIjkToGeo(h *faceIJK, res int32, g *LatLng) {
	var v vec2d
	_ijkToHex2d(&h.Coord, &v)
	_hex2dToGeo(&v, h.Face, res, 0, g)
}
