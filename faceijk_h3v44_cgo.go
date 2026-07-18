//go:build cgo && c2go && !h3v450

package h3

// Wrappers for faceijk.c internals that exist only in the H3 4.4.0 tree:
// _geoToClosestFace, _geoToHex2d, _hex2dToGeo, and _faceIjkToGeo were
// replaced by the Vec3 projection pipeline in 4.5.0
// (docs/sync/4.4.0-to-4.5.0.md §5.2). This file and its parity tests
// retire with I-A; docs/sync/h3v450-exclusion-inventory.md tracks the
// exclusion.

/*
#include "faceijk.h"

// Inline C helper to call _geoToClosestFace
static inline void geo_to_closest_face_c(const LatLng* g, int* face, double* sqd) {
    _geoToClosestFace(g, face, sqd);
}

// Inline C helper to call _geoToHex2d
static inline void geo_to_hex2d_c(const LatLng* g, int res, int* face, Vec2d* v) {
    _geoToHex2d(g, res, face, v);
}

// Inline C helper to call _hex2dToGeo
static inline void hex2d_to_geo_c(const Vec2d* v, int face, int res, int substrate, LatLng* g) {
    _hex2dToGeo(v, face, res, substrate, g);
}

// Inline C helper to call _faceIjkToGeo
static inline void face_ijk_to_geo_c(const FaceIJK* h, int res, LatLng* g) {
    _faceIjkToGeo(h, res, g);
}
*/
import "C"

// _geoToClosestFaceC calls the original C implementation.
func _geoToClosestFaceC(g *LatLng, face *int32, sqd *float64) {
	var cg C.LatLng
	cg.lat = C.double(g.Lat.Rad())
	cg.lng = C.double(g.Lng.Rad())

	var cface C.int
	var csqd C.double

	C.geo_to_closest_face_c(&cg, &cface, &csqd)

	*face = int32(cface)
	*sqd = float64(csqd)
}

// _geoToHex2dC calls the original C implementation.
func _geoToHex2dC(g *LatLng, res int32, face *int32, v *vec2d) {
	var cg C.LatLng
	cg.lat = C.double(g.Lat.Rad())
	cg.lng = C.double(g.Lng.Rad())

	var cface C.int
	var cv C.Vec2d

	C.geo_to_hex2d_c(&cg, C.int(res), &cface, &cv)

	*face = int32(cface)
	v.X = float64(cv.x)
	v.Y = float64(cv.y)
}

// _hex2dToGeoC calls the original C implementation.
func _hex2dToGeoC(v *vec2d, face int32, res int32, substrate int32, g *LatLng) {
	var cv C.Vec2d
	cv.x = C.double(v.X)
	cv.y = C.double(v.Y)

	var cg C.LatLng

	C.hex2d_to_geo_c(&cv, C.int(face), C.int(res), C.int(substrate), &cg)

	g.Lat = Rad(float64(cg.lat))
	g.Lng = Rad(float64(cg.lng))
}

// debugGeoToFaceIjkC calls the original C _geoToFaceIjk implementation.
func debugGeoToFaceIjkC(g *LatLng, res int32, fijk *faceIJK) {
	var cg C.LatLng
	cg.lat = C.double(g.Lat.Rad())
	cg.lng = C.double(g.Lng.Rad())

	var cfijk C.FaceIJK
	C._geoToFaceIjk(&cg, C.int(res), &cfijk)

	fijk.Face = int32(cfijk.face)
	fijk.Coord.I = int32(cfijk.coord.i)
	fijk.Coord.J = int32(cfijk.coord.j)
	fijk.Coord.K = int32(cfijk.coord.k)
}

// _faceIjkToGeoC calls the original C implementation.
func _faceIjkToGeoC(h *faceIJK, res int32, g *LatLng) {
	var ch C.FaceIJK
	ch.face = C.int(h.Face)
	ch.coord.i = C.int(h.Coord.I)
	ch.coord.j = C.int(h.Coord.J)
	ch.coord.k = C.int(h.Coord.K)

	var cg C.LatLng

	C.face_ijk_to_geo_c(&ch, C.int(res), &cg)

	g.Lat = Rad(float64(cg.lat))
	g.Lng = Rad(float64(cg.lng))
}
