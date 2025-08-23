//go:build cgo

package c2go

/*
#include <stdint.h>
#include <stdbool.h>
#include <stdlib.h>
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

// Inline C helper to call _adjustPentVertOverage
static inline int adjust_pent_vert_overage_c(FaceIJK* fijk, int res) {
    return _adjustPentVertOverage(fijk, res);
}
*/
import "C"

// _geoToClosestFaceC calls the original C implementation.
func _geoToClosestFaceC(g *LatLng, face *int, sqd *float64) {
	var cg C.LatLng
	cg.lat = C.double(g.Lat)
	cg.lng = C.double(g.Lng)

	var cface C.int
	var csqd C.double

	C.geo_to_closest_face_c(&cg, &cface, &csqd)

	*face = int(cface)
	*sqd = float64(csqd)
}

// _geoToHex2dC calls the original C implementation.
func _geoToHex2dC(g *LatLng, res int, face *int, v *Vec2d) {
	var cg C.LatLng
	cg.lat = C.double(g.Lat)
	cg.lng = C.double(g.Lng)

	var cface C.int
	var cv C.Vec2d

	C.geo_to_hex2d_c(&cg, C.int(res), &cface, &cv)

	*face = int(cface)
	v.X = float64(cv.x)
	v.Y = float64(cv.y)
}

// _hex2dToGeoC calls the original C implementation.
func _hex2dToGeoC(v *Vec2d, face int, res int, substrate int, g *LatLng) {
	var cv C.Vec2d
	cv.x = C.double(v.X)
	cv.y = C.double(v.Y)

	var cg C.LatLng

	C.hex2d_to_geo_c(&cv, C.int(face), C.int(res), C.int(substrate), &cg)

	g.Lat = float64(cg.lat)
	g.Lng = float64(cg.lng)
}

// _faceIjkToGeoC calls the original C implementation.
func _faceIjkToGeoC(h *FaceIJK, res int, g *LatLng) {
	var ch C.FaceIJK
	ch.face = C.int(h.Face)
	ch.coord.i = C.int(h.Coord.I)
	ch.coord.j = C.int(h.Coord.J)
	ch.coord.k = C.int(h.Coord.K)

	var cg C.LatLng

	C.face_ijk_to_geo_c(&ch, C.int(res), &cg)

	g.Lat = float64(cg.lat)
	g.Lng = float64(cg.lng)
}

// _adjustOverageClassIIC calls the original C implementation.
func _adjustOverageClassIIC(fijk *FaceIJK, res int, pentLeading4 int, substrate int) Overage {
	var cFijk C.FaceIJK
	cFijk.face = C.int(fijk.Face)
	cFijk.coord.i = C.int(fijk.Coord.I)
	cFijk.coord.j = C.int(fijk.Coord.J)
	cFijk.coord.k = C.int(fijk.Coord.K)

	result := int(C._adjustOverageClassII(&cFijk, C.int(res), C.int(pentLeading4), C.int(substrate)))

	// Update the Go struct with results
	fijk.Face = int(cFijk.face)
	fijk.Coord.I = int32(cFijk.coord.i)
	fijk.Coord.J = int32(cFijk.coord.j)
	fijk.Coord.K = int32(cFijk.coord.k)

	return Overage(result)
}

// _faceIjkToVertsC calls the original C implementation.
func _faceIjkToVertsC(fijk *FaceIJK, res *int, fijkVerts []FaceIJK) {
	var cFijk C.FaceIJK
	cFijk.face = C.int(fijk.Face)
	cFijk.coord.i = C.int(fijk.Coord.I)
	cFijk.coord.j = C.int(fijk.Coord.J)
	cFijk.coord.k = C.int(fijk.Coord.K)

	var cRes C.int = C.int(*res)

	// Create C array for output
	var cVerts [6]C.FaceIJK
	C._faceIjkToVerts(&cFijk, &cRes, &cVerts[0])

	// Update resolution
	*res = int(cRes)

	// Copy results back to Go slice
	for i := 0; i < 6; i++ {
		fijkVerts[i].Face = int(cVerts[i].face)
		fijkVerts[i].Coord.I = int32(cVerts[i].coord.i)
		fijkVerts[i].Coord.J = int32(cVerts[i].coord.j)
		fijkVerts[i].Coord.K = int32(cVerts[i].coord.k)
	}
}

// _faceIjkPentToVertsC calls the original C implementation.
func _faceIjkPentToVertsC(fijk *FaceIJK, res *int, fijkVerts []FaceIJK) {
	var cFijk C.FaceIJK
	cFijk.face = C.int(fijk.Face)
	cFijk.coord.i = C.int(fijk.Coord.I)
	cFijk.coord.j = C.int(fijk.Coord.J)
	cFijk.coord.k = C.int(fijk.Coord.K)

	var cRes C.int = C.int(*res)

	// Create C array for output
	var cVerts [5]C.FaceIJK
	C._faceIjkPentToVerts(&cFijk, &cRes, &cVerts[0])

	// Update resolution
	*res = int(cRes)

	// Copy results back to Go slice
	for i := 0; i < 5; i++ {
		fijkVerts[i].Face = int(cVerts[i].face)
		fijkVerts[i].Coord.I = int32(cVerts[i].coord.i)
		fijkVerts[i].Coord.J = int32(cVerts[i].coord.j)
		fijkVerts[i].Coord.K = int32(cVerts[i].coord.k)
	}
}

// _adjustPentVertOverageC calls the original C implementation.
func _adjustPentVertOverageC(fijk *FaceIJK, res int) Overage {
	var cFijk C.FaceIJK
	cFijk.face = C.int(fijk.Face)
	cFijk.coord.i = C.int(fijk.Coord.I)
	cFijk.coord.j = C.int(fijk.Coord.J)
	cFijk.coord.k = C.int(fijk.Coord.K)

	result := int(C.adjust_pent_vert_overage_c(&cFijk, C.int(res)))

	// Update the Go struct with results
	fijk.Face = int(cFijk.face)
	fijk.Coord.I = int32(cFijk.coord.i)
	fijk.Coord.J = int32(cFijk.coord.j)
	fijk.Coord.K = int32(cFijk.coord.k)

	return Overage(result)
}
