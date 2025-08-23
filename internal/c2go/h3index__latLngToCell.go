package c2go

import "math"

// latLngToCell encodes a coordinate pair into an H3Index at the given resolution.
// Ported from H3 C: h3Index.c::latLngToCell
func latLngToCell(g *LatLng, res int, out *H3Index) H3Error {
	if res < 0 || res > MAX_H3_RES {
		return E_RES_DOMAIN
	}
	if math.IsInf(g.Lat, 0) || math.IsNaN(g.Lat) || math.IsInf(g.Lng, 0) || math.IsNaN(g.Lng) {
		return E_LATLNG_DOMAIN
	}

	var fijk FaceIJK
	_geoToFaceIjk(g, res, &fijk)
	*out = _faceIjkToH3(&fijk, res)
	// ALWAYS(*out) in C - check if result is truthy
	if *out != 0 {
		return E_SUCCESS
	} else {
		return E_FAILED
	}
}
