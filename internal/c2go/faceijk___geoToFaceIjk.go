package c2go

// _geoToFaceIjk encodes a coordinate on the sphere to the FaceIJK address of the containing cell.
// First converts to hex2d coordinates, then converts to IJK coordinates.
// Ported from H3 C: faceijk.c::_geoToFaceIjk
func _geoToFaceIjk(g *LatLng, res int, h *FaceIJK) {
	// first convert to hex2d
	var v Vec2d
	_geoToHex2d(g, res, &h.Face, &v)

	// then convert to ijk+
	_hex2dToCoordIJK(&v, &h.Coord)
}