//go:build cgo && c2go

package h3

/*
#cgo CFLAGS: -I${SRCDIR}/testref/h3-4.3.0/src/h3lib/include -I${SRCDIR}/testref/h3-4.3.0/src/apps/applib/include
#cgo CFLAGS: -std=c99

#include <stdint.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "h3api.h"
#include "utility.h"

// Callback type for H3Index iteration
typedef void (*h3_callback)(H3Index);

// CGO wrappers for print functions
static void h3Print_wrapper(H3Index h) {
    h3Print(h);
}

static void h3Println_wrapper(H3Index h) {
    h3Println(h);
}

static void coordIjkPrint_wrapper(const CoordIJK *c) {
    coordIjkPrint(c);
}

static void geoToStringRads_wrapper(const LatLng *p, char *str) {
    geoToStringRads(p, str);
}

static void geoToStringDegs_wrapper(const LatLng *p, char *str) {
    geoToStringDegs(p, str);
}

static void geoToStringDegsNoFmt_wrapper(const LatLng *p, char *str) {
    geoToStringDegsNoFmt(p, str);
}

static void geoPrint_wrapper(const LatLng *p) {
    geoPrint(p);
}

static void geoPrintln_wrapper(const LatLng *p) {
    geoPrintln(p);
}

static void geoPrintNoFmt_wrapper(const LatLng *p) {
    geoPrintNoFmt(p);
}

static void geoPrintlnNoFmt_wrapper(const LatLng *p) {
    geoPrintlnNoFmt(p);
}

static void cellBoundaryPrint_wrapper(const CellBoundary *b) {
    cellBoundaryPrint(b);
}

static void cellBoundaryPrintln_wrapper(const CellBoundary *b) {
    cellBoundaryPrintln(b);
}

static void bboxPrint_wrapper(const BBox *bbox) {
    bboxPrint(bbox);
}

static void bboxPrintln_wrapper(const BBox *bbox) {
    bboxPrintln(bbox);
}

// Global variables to store Go callback results
static H3Index* callback_buffer = NULL;
static int callback_index = 0;
static int callback_capacity = 0;

// C callback function that stores H3Index values
static void store_h3index(H3Index h) {
    if (callback_buffer != NULL && callback_index < callback_capacity) {
        callback_buffer[callback_index++] = h;
    }
}

// Wrapper for iterateAllIndexesAtResPartial
static int iterateAllIndexesAtResPartial_wrapper(int res, int baseCells, H3Index* out, int outSize) {
    callback_buffer = out;
    callback_index = 0;
    callback_capacity = outSize;
    iterateAllIndexesAtResPartial(res, store_h3index, baseCells);
    int result = callback_index;
    callback_buffer = NULL;
    callback_index = 0;
    callback_capacity = 0;
    return result;
}

// Wrapper for iterateAllIndexesAtRes
static int iterateAllIndexesAtRes_wrapper(int res, H3Index* out, int outSize) {
    callback_buffer = out;
    callback_index = 0;
    callback_capacity = outSize;
    iterateAllIndexesAtRes(res, store_h3index);
    int result = callback_index;
    callback_buffer = NULL;
    callback_index = 0;
    callback_capacity = 0;
    return result;
}

// Wrapper for iterateBaseCellIndexesAtRes
static int iterateBaseCellIndexesAtRes_wrapper(int res, int baseCell, H3Index* out, int outSize) {
    callback_buffer = out;
    callback_index = 0;
    callback_capacity = outSize;
    iterateBaseCellIndexesAtRes(res, store_h3index, baseCell);
    int result = callback_index;
    callback_buffer = NULL;
    callback_index = 0;
    callback_capacity = 0;
    return result;
}
*/
import "C"
import (
	"unsafe"
)

// iterateAllIndexesAtResPartialC calls the C implementation of iterateAllIndexesAtResPartial
func iterateAllIndexesAtResPartialC(res int32, baseCells int32, out []h3Index) int {
	if len(out) == 0 {
		return 0
	}
	count := C.iterateAllIndexesAtResPartial_wrapper(
		C.int(res),
		C.int(baseCells),
		(*C.H3Index)(unsafe.Pointer(&out[0])),
		C.int(len(out)),
	)
	return int(count)
}

// iterateAllIndexesAtResC calls the C implementation of iterateAllIndexesAtRes
func iterateAllIndexesAtResC(res int32, out []h3Index) int {
	if len(out) == 0 {
		return 0
	}
	count := C.iterateAllIndexesAtRes_wrapper(
		C.int(res),
		(*C.H3Index)(unsafe.Pointer(&out[0])),
		C.int(len(out)),
	)
	return int(count)
}

// iterateBaseCellIndexesAtResC calls the C implementation of iterateBaseCellIndexesAtRes
func iterateBaseCellIndexesAtResC(res int32, baseCell int32, out []h3Index) int {
	if len(out) == 0 {
		return 0
	}
	count := C.iterateBaseCellIndexesAtRes_wrapper(
		C.int(res),
		C.int(baseCell),
		(*C.H3Index)(unsafe.Pointer(&out[0])),
		C.int(len(out)),
	)
	return int(count)
}

// h3PrintC calls the C implementation of h3Print
func h3PrintC(h h3Index) {
	C.h3Print_wrapper(C.H3Index(h))
}

// h3PrintlnC calls the C implementation of h3Println
func h3PrintlnC(h h3Index) {
	C.h3Println_wrapper(C.H3Index(h))
}

// coordIjkPrintC calls the C implementation of coordIjkPrint
func coordIjkPrintC(c *coordIJK) {
	cCoord := C.CoordIJK{
		i: C.int(c.I),
		j: C.int(c.J),
		k: C.int(c.K),
	}
	C.coordIjkPrint_wrapper(&cCoord)
}

// geoToStringRadsC calls the C implementation of geoToStringRads
func geoToStringRadsC(p *LatLng) string {
	cLatLng := C.LatLng{
		lat: C.double(p.Lat.Rad()),
		lng: C.double(p.Lng.Rad()),
	}
	str := C.malloc(buffSize)
	defer C.free(str)
	C.geoToStringRads_wrapper(&cLatLng, (*C.char)(str))
	return C.GoString((*C.char)(str))
}

// geoToStringDegsC calls the C implementation of geoToStringDegs
func geoToStringDegsC(p *LatLng) string {
	cLatLng := C.LatLng{
		lat: C.double(p.Lat.Rad()),
		lng: C.double(p.Lng.Rad()),
	}
	str := C.malloc(buffSize)
	defer C.free(str)
	C.geoToStringDegs_wrapper(&cLatLng, (*C.char)(str))
	return C.GoString((*C.char)(str))
}

// geoToStringDegsNoFmtC calls the C implementation of geoToStringDegsNoFmt
func geoToStringDegsNoFmtC(p *LatLng) string {
	cLatLng := C.LatLng{
		lat: C.double(p.Lat.Rad()),
		lng: C.double(p.Lng.Rad()),
	}
	str := C.malloc(buffSize)
	defer C.free(str)
	C.geoToStringDegsNoFmt_wrapper(&cLatLng, (*C.char)(str))
	return C.GoString((*C.char)(str))
}

// geoPrintC calls the C implementation of geoPrint
func geoPrintC(p *LatLng) {
	cLatLng := C.LatLng{
		lat: C.double(p.Lat.Rad()),
		lng: C.double(p.Lng.Rad()),
	}
	C.geoPrint_wrapper(&cLatLng)
}

// geoPrintlnC calls the C implementation of geoPrintln
func geoPrintlnC(p *LatLng) {
	cLatLng := C.LatLng{
		lat: C.double(p.Lat.Rad()),
		lng: C.double(p.Lng.Rad()),
	}
	C.geoPrintln_wrapper(&cLatLng)
}

// geoPrintNoFmtC calls the C implementation of geoPrintNoFmt
func geoPrintNoFmtC(p *LatLng) {
	cLatLng := C.LatLng{
		lat: C.double(p.Lat.Rad()),
		lng: C.double(p.Lng.Rad()),
	}
	C.geoPrintNoFmt_wrapper(&cLatLng)
}

// geoPrintlnNoFmtC calls the C implementation of geoPrintlnNoFmt
func geoPrintlnNoFmtC(p *LatLng) {
	cLatLng := C.LatLng{
		lat: C.double(p.Lat.Rad()),
		lng: C.double(p.Lng.Rad()),
	}
	C.geoPrintlnNoFmt_wrapper(&cLatLng)
}

// cellBoundaryPrintC calls the C implementation of cellBoundaryPrint
func cellBoundaryPrintC(b *CellBoundary) {
	// C CellBoundary has a fixed array of 10 LatLng vertices
	var cBoundary C.CellBoundary
	cBoundary.numVerts = C.int(b.numVerts)

	// Copy vertices to the fixed-size C array
	for i := 0; i < len(b.verts) && i < 10; i++ {
		// We need to access the array elements directly
		// This is a workaround since we can't easily assign to the fixed array
		cVertPtr := (*C.LatLng)(unsafe.Pointer(uintptr(unsafe.Pointer(&cBoundary.verts[0])) + uintptr(i)*unsafe.Sizeof(C.LatLng{})))
		cVertPtr.lat = C.double(b.verts[i].Lat.Rad())
		cVertPtr.lng = C.double(b.verts[i].Lng.Rad())
	}
	C.cellBoundaryPrint_wrapper(&cBoundary)
}

// cellBoundaryPrintlnC calls the C implementation of cellBoundaryPrintln
func cellBoundaryPrintlnC(b *CellBoundary) {
	// C CellBoundary has a fixed array of 10 LatLng vertices
	var cBoundary C.CellBoundary
	cBoundary.numVerts = C.int(b.numVerts)

	// Copy vertices to the fixed-size C array
	for i := 0; i < len(b.verts) && i < 10; i++ {
		// We need to access the array elements directly
		// This is a workaround since we can't easily assign to the fixed array
		cVertPtr := (*C.LatLng)(unsafe.Pointer(uintptr(unsafe.Pointer(&cBoundary.verts[0])) + uintptr(i)*unsafe.Sizeof(C.LatLng{})))
		cVertPtr.lat = C.double(b.verts[i].Lat.Rad())
		cVertPtr.lng = C.double(b.verts[i].Lng.Rad())
	}
	C.cellBoundaryPrintln_wrapper(&cBoundary)
}

// bboxPrintC calls the C implementation of bboxPrint
func bboxPrintC(bbox *bbox) {
	cBBox := C.BBox{
		north: C.double(bbox.North.Rad()),
		south: C.double(bbox.South.Rad()),
		east:  C.double(bbox.East.Rad()),
		west:  C.double(bbox.West.Rad()),
	}
	C.bboxPrint_wrapper(&cBBox)
}

// bboxPrintlnC calls the C implementation of bboxPrintln
func bboxPrintlnC(bbox *bbox) {
	cBBox := C.BBox{
		north: C.double(bbox.North.Rad()),
		south: C.double(bbox.South.Rad()),
		east:  C.double(bbox.East.Rad()),
		west:  C.double(bbox.West.Rad()),
	}
	C.bboxPrintln_wrapper(&cBBox)
}
