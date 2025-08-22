//go:build cgo && c2go

package c2go

/*
#include <stdint.h>
#include <stdbool.h>
#include "h3api.h"
#include "polygon.h"
*/
import "C"

// validatePolygonFlagsC calls the original C implementation.
func validatePolygonFlagsC(flags uint32) uint32 {
    return uint32(C.validatePolygonFlags(C.uint(flags)))
}

// lineCrossesLineC calls the original C implementation.
func lineCrossesLineC(a1, a2, b1, b2 LatLng) bool {
    var ca1, ca2, cb1, cb2 C.LatLng
    ca1.lat = C.double(a1.Lat); ca1.lng = C.double(a1.Lng)
    ca2.lat = C.double(a2.Lat); ca2.lng = C.double(a2.Lng)
    cb1.lat = C.double(b1.Lat); cb1.lng = C.double(b1.Lng)
    cb2.lat = C.double(b2.Lat); cb2.lng = C.double(b2.Lng)
    if C.lineCrossesLine(&ca1, &ca2, &cb1, &cb2) { return true } else { return false }
}
