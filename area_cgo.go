//go:build cgo && c2go && h3v450

package h3

// Bridges to the H3 4.5.0 area implementation (area.c / area.h):
// geoLoopAreaRads2 is extern in area.h; cellArea* are public and bridged
// version-neutrally in latLng_cgo.go.

/*
#include "h3api.h"
#include "area.h"
*/
import "C"
import "unsafe"

// geoLoopAreaRads2C calls the original C implementation.
func geoLoopAreaRads2C(loop GeoLoop) (float64, h3Error) {
	var cLoop C.GeoLoop
	if len(loop) > 0 {
		cLoop.verts = (*C.LatLng)(unsafe.Pointer(&loop[0]))
		cLoop.numVerts = C.int(len(loop))
	}
	var out C.double
	err := h3Error(C.geoLoopAreaRads2(cLoop, &out))
	return float64(out), err
}
