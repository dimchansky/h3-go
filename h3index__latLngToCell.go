package h3

import "math"

// latLngToCell encodes a coordinate pair into an h3Index at the given resolution.
// Ported from H3 C: h3Index.c::latLngToCell.
func latLngToCell(g *LatLng, res int32, out *h3Index) h3Error {
	if res < 0 || res > maxH3Res {
		return eResDomain
	}
	latRad := g.Lat.Rad()
	lngRad := g.Lng.Rad()
	if math.IsInf(latRad, 0) || math.IsNaN(latRad) || math.IsInf(lngRad, 0) || math.IsNaN(lngRad) {
		return eLatlngDomain
	}

	var fijk faceIJK
	_geoToFaceIjk(g, res, &fijk)
	*out = _faceIjkToH3(&fijk, res)
	// ALWAYS(*out) in C - check if result is truthy
	if *out != 0 {
		return eSuccess
	} else {
		return eFailed
	}
}
