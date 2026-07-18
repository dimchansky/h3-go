//go:build cgo && c2go && h3v450

package h3

// Bridges to the H3 4.5.0 area implementation (area.c / area.h / adder.h):
// geoLoopAreaRads2, geoPolygonAreaRads2, and geoMultiPolygonAreaRads2 are
// extern in area.h; kadd is static inline in adder.h (each including TU
// gets its own definition); cagnoli is reached through the same-TU
// wrapper in h3lib_area_c2go.c; cellArea* are public and bridged
// version-neutrally in latLng_cgo.go.

/*
#include <stdlib.h>
#include "h3api.h"
#include "area.h"
#include "adder.h"

double h3goTest_cagnoli(LatLng x, LatLng y);

// kadd mutates in place; return the updated state by value for the bridge.
static Adder kadd_c_wrapper(Adder a, double x) {
    kadd(&a, x);
    return a;
}
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

// kaddC calls the original C implementation on an adder state.
func kaddC(a adder, x float64) adder {
	var ca C.Adder
	ca.sum = C.double(a.sum)
	ca._c = C.double(a.c)
	ca = C.kadd_c_wrapper(ca, C.double(x))
	return adder{sum: float64(ca.sum), c: float64(ca._c)}
}

// _cagnoliC calls the original C implementation (same-TU wrapper).
func _cagnoliC(x, y LatLng) float64 {
	var cx, cy C.LatLng
	cx.lat = C.double(x.Lat.Rad())
	cx.lng = C.double(x.Lng.Rad())
	cy.lat = C.double(y.Lat.Rad())
	cy.lng = C.double(y.Lng.Rad())
	return float64(C.h3goTest_cagnoli(cx, cy))
}

// areaToCGeoLoop converts a Go GeoLoop into a C GeoLoop backed by
// C-allocated memory; the caller must invoke the returned free function.
func areaToCGeoLoop(loop GeoLoop) (C.GeoLoop, func()) {
	var cl C.GeoLoop
	n := len(loop)
	cl.numVerts = C.int(n)
	if n == 0 {
		return cl, func() {}
	}
	mem := C.malloc(C.size_t(n) * C.size_t(C.sizeof_LatLng))
	verts := (*[1 << 30]C.LatLng)(mem)[:n:n]
	for i, ll := range loop {
		verts[i].lat = C.double(ll.Lat.Rad())
		verts[i].lng = C.double(ll.Lng.Rad())
	}
	cl.verts = (*C.LatLng)(mem)
	return cl, func() { C.free(mem) }
}

// areaToCGeoPolygon converts a Go GeoPolygon into a C GeoPolygon; the
// caller must invoke the returned free function.
func areaToCGeoPolygon(poly GeoPolygon) (C.GeoPolygon, func()) {
	var cp C.GeoPolygon
	outer, freeOuter := areaToCGeoLoop(poly.GeoLoop)
	cp.geoloop = outer
	n := len(poly.Holes)
	cp.numHoles = C.int(n)
	frees := []func(){freeOuter}
	if n > 0 {
		mem := C.malloc(C.size_t(n) * C.size_t(C.sizeof_GeoLoop))
		holes := (*[1 << 30]C.GeoLoop)(mem)[:n:n]
		for i, h := range poly.Holes {
			ch, fh := areaToCGeoLoop(h)
			holes[i] = ch
			frees = append(frees, fh)
		}
		cp.holes = (*C.GeoLoop)(mem)
		frees = append(frees, func() { C.free(mem) })
	}
	return cp, func() {
		for _, f := range frees {
			f()
		}
	}
}

// geoPolygonAreaRads2C calls the original C implementation.
func geoPolygonAreaRads2C(poly GeoPolygon) (float64, h3Error) {
	cp, freeFn := areaToCGeoPolygon(poly)
	defer freeFn()
	var out C.double
	err := h3Error(C.geoPolygonAreaRads2(cp, &out))
	return float64(out), err
}

// geoMultiPolygonAreaRads2C calls the original C implementation.
func geoMultiPolygonAreaRads2C(mpoly geoMultiPolygon) (float64, h3Error) {
	var cm C.GeoMultiPolygon
	n := int(mpoly.NumPolygons)
	cm.numPolygons = C.int(n)
	var frees []func()
	if n > 0 {
		mem := C.malloc(C.size_t(n) * C.size_t(C.sizeof_GeoPolygon))
		polys := (*[1 << 30]C.GeoPolygon)(mem)[:n:n]
		for i := 0; i < n; i++ {
			cp, fp := areaToCGeoPolygon(mpoly.Polygons[i])
			polys[i] = cp
			frees = append(frees, fp)
		}
		cm.polygons = (*C.GeoPolygon)(mem)
		frees = append(frees, func() { C.free(mem) })
	}
	defer func() {
		for _, f := range frees {
			f()
		}
	}()
	var out C.double
	err := h3Error(C.geoMultiPolygonAreaRads2(cm, &out))
	return float64(out), err
}
