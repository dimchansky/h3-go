package h3

// _faceIjkToGeo converts FaceIJK coordinates to geographic coordinates.
// Ported from H3 C: faceijk.c::_faceIjkToGeo
func _faceIjkToGeo(h *FaceIJK, res int32, g *LatLng) {
	var v Vec2d
	_ijkToHex2d(&h.Coord, &v)
	_hex2dToGeo(&v, h.Face, res, 0, g)
}
