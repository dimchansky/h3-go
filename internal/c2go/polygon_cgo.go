//go:build cgo && c2go

package c2go

/*
#include <stdint.h>
#include <stdbool.h>
#include "h3api.h"
#include "polygon.h"
*/
import "C"
import "unsafe"

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

// toCGeoLoop builds a C GeoLoop from Go LatLng slice. Caller must call freeFn.
func toCGeoLoop(verts []LatLng) (C.GeoLoop, func()) {
    var loop C.GeoLoop
    n := len(verts)
    loop.numVerts = C.int(n)
    if n == 0 {
        loop.verts = nil
        return loop, func() {}
    }
    bytes := C.size_t(n) * C.size_t(C.sizeof_LatLng)
    mem := C.malloc(bytes)
    // Fill array
    arr := (*[1 << 30]C.LatLng)(mem)[:n:n]
    for i, v := range verts {
        arr[i].lat = C.double(v.Lat)
        arr[i].lng = C.double(v.Lng)
    }
    loop.verts = (*C.LatLng)(mem)
    freeFn := func() { C.free(unsafe.Pointer(mem)) }
    return loop, freeFn
}

// bboxFromGeoLoopC calls C bboxFromGeoLoop on a Go slice of LatLng.
func bboxFromGeoLoopC(loop []LatLng) BBox {
    cg, freeFn := toCGeoLoop(loop)
    defer freeFn()
    var cb C.BBox
    C.bboxFromGeoLoop(&cg, &cb)
    return BBox{North: float64(cb.north), South: float64(cb.south), East: float64(cb.east), West: float64(cb.west)}
}

// pointInsideGeoLoopC calls C pointInsideGeoLoop for a point within a GeoLoop + BBox.
func pointInsideGeoLoopC(loop []LatLng, bbox BBox, p LatLng) bool {
    cg, freeFn := toCGeoLoop(loop)
    defer freeFn()
    var cb C.BBox
    cb.north = C.double(bbox.North)
    cb.south = C.double(bbox.South)
    cb.east = C.double(bbox.East)
    cb.west = C.double(bbox.West)
    var cp C.LatLng
    cp.lat = C.double(p.Lat)
    cp.lng = C.double(p.Lng)
    if C.pointInsideGeoLoop(&cg, &cb, &cp) { return true } else { return false }
}
