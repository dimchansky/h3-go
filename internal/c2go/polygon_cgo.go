//go:build cgo

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
	ca1.lat = C.double(a1.Lat)
	ca1.lng = C.double(a1.Lng)
	ca2.lat = C.double(a2.Lat)
	ca2.lng = C.double(a2.Lng)
	cb1.lat = C.double(b1.Lat)
	cb1.lng = C.double(b1.Lng)
	cb2.lat = C.double(b2.Lat)
	cb2.lng = C.double(b2.Lng)
	if C.lineCrossesLine(&ca1, &ca2, &cb1, &cb2) {
		return true
	} else {
		return false
	}
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
	if C.pointInsideGeoLoop(&cg, &cb, &cp) {
		return true
	} else {
		return false
	}
}

// toCGeoPolygon converts Go GeoPolygon to C GeoPolygon; caller must free via returned function.
func toCGeoPolygon(poly GeoPolygon) (C.GeoPolygon, func()) {
	var cp C.GeoPolygon
	// Outer loop
	outer, freeOuter := toCGeoLoop(poly.Geoloop)
	cp.geoloop = outer
	// Holes
	n := len(poly.Holes)
	cp.numHoles = C.int(n)
	var freeHoles []func()
	var holesMem unsafe.Pointer
	if n > 0 {
		holesMem = C.malloc(C.size_t(n) * C.size_t(C.sizeof_GeoLoop))
		holesArr := (*[1 << 30]C.GeoLoop)(holesMem)[:n:n]
		for i, h := range poly.Holes {
			ch, fh := toCGeoLoop(h)
			holesArr[i] = ch
			freeHoles = append(freeHoles, fh)
		}
		cp.holes = (*C.GeoLoop)(holesMem)
	} else {
		cp.holes = nil
	}
	freeFn := func() {
		for _, fh := range freeHoles {
			fh()
		}
		if holesMem != nil {
			C.free(unsafe.Pointer(holesMem))
		}
		freeOuter()
	}
	return cp, freeFn
}

// pointInsidePolygonC calls C pointInsidePolygon for a GeoPolygon.
func pointInsidePolygonC(poly GeoPolygon, bboxes []BBox, p LatLng) bool {
	cp, freePoly := toCGeoPolygon(poly)
	defer freePoly()
	// Prepare bboxes array
	nb := len(bboxes)
	var cbptr *C.BBox
	var bbMem unsafe.Pointer
	if nb > 0 {
		bbMem = C.malloc(C.size_t(nb) * C.size_t(C.sizeof_BBox))
		arr := (*[1 << 30]C.BBox)(bbMem)[:nb:nb]
		for i, b := range bboxes {
			arr[i].north = C.double(b.North)
			arr[i].south = C.double(b.South)
			arr[i].east = C.double(b.East)
			arr[i].west = C.double(b.West)
		}
		cbptr = (*C.BBox)(bbMem)
		defer C.free(unsafe.Pointer(bbMem))
	}
	var cpnt C.LatLng
	cpnt.lat = C.double(p.Lat)
	cpnt.lng = C.double(p.Lng)
	if C.pointInsidePolygon(&cp, cbptr, &cpnt) {
		return true
	} else {
		return false
	}
}

// cellBoundaryCrossesGeoLoopC calls the C implementation.
func cellBoundaryCrossesGeoLoopC(geoloop GeoLoop, loopBBox BBox, boundary CellBoundary, boundaryBBox BBox) bool {
	cg, freeFn := toCGeoLoop(geoloop)
	defer freeFn()
	var cLoopBBox C.BBox
	cLoopBBox.north = C.double(loopBBox.North)
	cLoopBBox.south = C.double(loopBBox.South)
	cLoopBBox.east = C.double(loopBBox.East)
	cLoopBBox.west = C.double(loopBBox.West)
	var cBoundaryBBox C.BBox
	cBoundaryBBox.north = C.double(boundaryBBox.North)
	cBoundaryBBox.south = C.double(boundaryBBox.South)
	cBoundaryBBox.east = C.double(boundaryBBox.East)
	cBoundaryBBox.west = C.double(boundaryBBox.West)
	var cb C.CellBoundary
	n := boundary.NumVerts
	if n > 0 {
		if n > int(C.MAX_CELL_BNDRY_VERTS) {
			n = int(C.MAX_CELL_BNDRY_VERTS)
		}
		cb.numVerts = C.int(n)
		for i := 0; i < n; i++ {
			cb.verts[i].lat = C.double(boundary.Verts[i].Lat)
			cb.verts[i].lng = C.double(boundary.Verts[i].Lng)
		}
	} else {
		cb.numVerts = 0
	}
	if C.cellBoundaryCrossesGeoLoop(&cg, &cLoopBBox, &cb, &cBoundaryBBox) {
		return true
	} else {
		return false
	}
}

// cellBoundaryInsidePolygonC calls C cellBoundaryInsidePolygon for a GeoPolygon.
func cellBoundaryInsidePolygonC(poly GeoPolygon, bboxes []BBox, boundary CellBoundary, boundaryBBox BBox) bool {
	cp, freePoly := toCGeoPolygon(poly)
	defer freePoly()
	// bboxes
	nb := len(bboxes)
	var cbptr *C.BBox
	var bbMem unsafe.Pointer
	if nb > 0 {
		bbMem = C.malloc(C.size_t(nb) * C.size_t(C.sizeof_BBox))
		arr := (*[1 << 30]C.BBox)(bbMem)[:nb:nb]
		for i, b := range bboxes {
			arr[i].north = C.double(b.North)
			arr[i].south = C.double(b.South)
			arr[i].east = C.double(b.East)
			arr[i].west = C.double(b.West)
		}
		cbptr = (*C.BBox)(bbMem)
		defer C.free(unsafe.Pointer(bbMem))
	}
	// boundary bbox
	var cBoundaryBBox C.BBox
	cBoundaryBBox.north = C.double(boundaryBBox.North)
	cBoundaryBBox.south = C.double(boundaryBBox.South)
	cBoundaryBBox.east = C.double(boundaryBBox.East)
	cBoundaryBBox.west = C.double(boundaryBBox.West)
	// boundary
	var cb C.CellBoundary
	n := boundary.NumVerts
	if n > 0 {
		if n > int(C.MAX_CELL_BNDRY_VERTS) {
			n = int(C.MAX_CELL_BNDRY_VERTS)
		}
		cb.numVerts = C.int(n)
		for i := 0; i < n; i++ {
			cb.verts[i].lat = C.double(boundary.Verts[i].Lat)
			cb.verts[i].lng = C.double(boundary.Verts[i].Lng)
		}
	}
	if C.cellBoundaryInsidePolygon(&cp, cbptr, &cb, &cBoundaryBBox) {
		return true
	} else {
		return false
	}
}

// cellBoundaryCrossesPolygonC calls C cellBoundaryCrossesPolygon for a GeoPolygon.
func cellBoundaryCrossesPolygonC(poly GeoPolygon, bboxes []BBox, boundary CellBoundary, boundaryBBox BBox) bool {
	cp, freePoly := toCGeoPolygon(poly)
	defer freePoly()
	// bboxes
	nb := len(bboxes)
	var cbptr *C.BBox
	var bbMem unsafe.Pointer
	if nb > 0 {
		bbMem = C.malloc(C.size_t(nb) * C.size_t(C.sizeof_BBox))
		arr := (*[1 << 30]C.BBox)(bbMem)[:nb:nb]
		for i, b := range bboxes {
			arr[i].north = C.double(b.North)
			arr[i].south = C.double(b.South)
			arr[i].east = C.double(b.East)
			arr[i].west = C.double(b.West)
		}
		cbptr = (*C.BBox)(bbMem)
		defer C.free(unsafe.Pointer(bbMem))
	}
	// boundary bbox
	var cBoundaryBBox C.BBox
	cBoundaryBBox.north = C.double(boundaryBBox.North)
	cBoundaryBBox.south = C.double(boundaryBBox.South)
	cBoundaryBBox.east = C.double(boundaryBBox.East)
	cBoundaryBBox.west = C.double(boundaryBBox.West)
	// boundary
	var cb C.CellBoundary
	n := boundary.NumVerts
	if n > 0 {
		if n > int(C.MAX_CELL_BNDRY_VERTS) {
			n = int(C.MAX_CELL_BNDRY_VERTS)
		}
		cb.numVerts = C.int(n)
		for i := 0; i < n; i++ {
			cb.verts[i].lat = C.double(boundary.Verts[i].Lat)
			cb.verts[i].lng = C.double(boundary.Verts[i].Lng)
		}
	}
	if C.cellBoundaryCrossesPolygon(&cp, cbptr, &cb, &cBoundaryBBox) {
		return true
	} else {
		return false
	}
}
