//go:build cgo && c2go

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
