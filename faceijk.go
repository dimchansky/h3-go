package h3

//lint:file-ignore U1000 Ignore all unused code

// FaceIJK holds face number and ijk coordinates on that face-centered coordinate system
type FaceIJK struct {
	face  int      // face number
	coord CoordIJK // ijk coordinates on that face
}

// _geoToFaceIjk encodes a coordinate on the sphere to the FaceIJK address of the containing
// cell at the specified resolution.
//
// `g`: The spherical coordinates to encode.
// `res`: The desired H3 resolution for the encoding.
// `h`: The FaceIJK address of the containing cell at resolution res.
func _geoToFaceIjk(g *GeoCoord, res int, h *FaceIJK) {
	/*
	   // first convert to hex2d
	   Vec2d v;
	   _geoToHex2d(g, res, &h->face, &v);

	   // then convert to ijk+
	   _hex2dToCoordIJK(&v, &h->coord);
	*/
	panic("not implemented")
}
