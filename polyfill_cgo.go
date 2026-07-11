//go:build cgo && c2go

package h3

/*

#include <stdlib.h>
#include <stdbool.h>
#include "h3api.h"
#include "bbox.h"

// Forward declaration for polygonToCellsExperimental
H3Error polygonToCellsExperimental(const GeoPolygon *polygon, int res, uint32_t flags,
                                   int64_t size, H3Index *out);

// Forward declaration for cellToBBox
H3Error cellToBBox(H3Index cell, BBox *out, bool coverChildren);

// C wrapper for polygonToCellsExperimental that handles C GeoPolygon conversion
static H3Error polygonToCellsExperimental_c_wrapper(
    const LatLng *geoloop, int numVerts,
    const LatLng **holes, const int *holeSizes, int numHoles,
    int res, uint32_t flags, int64_t size, H3Index *out) {

    // Build main geoloop
    GeoLoop mainLoop = {.numVerts = numVerts, .verts = (LatLng*)geoloop};

    // Build hole loops
    GeoLoop *holeLoops = NULL;
    if (numHoles > 0) {
        holeLoops = (GeoLoop*)malloc(numHoles * sizeof(GeoLoop));
        for (int i = 0; i < numHoles; i++) {
            holeLoops[i].numVerts = holeSizes[i];
            holeLoops[i].verts = (LatLng*)holes[i];
        }
    }

    // Build polygon
    GeoPolygon polygon = {.geoloop = mainLoop, .numHoles = numHoles, .holes = holeLoops};

    H3Error result = polygonToCellsExperimental(&polygon, res, flags, size, out);

    if (holeLoops) {
        free(holeLoops);
    }

    return result;
}

// C wrapper for maxPolygonToCellsSizeExperimental that handles C GeoPolygon conversion
static H3Error maxPolygonToCellsSizeExperimental_c_wrapper(
    const LatLng *geoloop, int numVerts,
    const LatLng **holes, const int *holeSizes, int numHoles,
    int res, uint32_t flags, int64_t *out) {

    // Build main geoloop
    GeoLoop mainLoop = {.numVerts = numVerts, .verts = (LatLng*)geoloop};

    // Build hole loops
    GeoLoop *holeLoops = NULL;
    if (numHoles > 0) {
        holeLoops = (GeoLoop*)malloc(numHoles * sizeof(GeoLoop));
        for (int i = 0; i < numHoles; i++) {
            holeLoops[i].numVerts = holeSizes[i];
            holeLoops[i].verts = (LatLng*)holes[i];
        }
    }

    // Build polygon
    GeoPolygon polygon = {.geoloop = mainLoop, .numHoles = numHoles, .holes = holeLoops};

    H3Error result = maxPolygonToCellsSizeExperimental(&polygon, res, flags, out);

    if (holeLoops) {
        free(holeLoops);
    }

    return result;
}
*/
import "C"
import (
	"unsafe"
)

// polygonToCellsExperimentalC calls the C polygonToCellsExperimental function
func polygonToCellsExperimentalC(polygon *GeoPolygon, res int32, flags uint32,
	size int64, out []h3Index) h3Error {
	if len(out) == 0 || int64(len(out)) < size {
		return eMemoryBounds
	}

	// Use a simple approach similar to other polygon functions
	// Convert to C structures and call directly

	// Build C GeoLoop for main loop
	mainLoop := polygon.GeoLoop
	if len(mainLoop) == 0 {
		// Empty polygon case
		return h3Error(C.polygonToCellsExperimental_c_wrapper(
			nil, 0, nil, nil, 0, C.int(res), C.uint32_t(flags),
			C.int64_t(size), (*C.H3Index)(unsafe.Pointer(&out[0]))))
	}

	// For now, simplify and only handle polygons without holes
	// to avoid CGO pointer issues
	if len(polygon.Holes) > 0 {
		// TODO: Implement hole handling once basic case works
		return h3Error(C.polygonToCellsExperimental_c_wrapper(
			(*C.LatLng)(unsafe.Pointer(&mainLoop[0])),
			C.int(len(mainLoop)),
			nil, nil, 0,
			C.int(res), C.uint32_t(flags), C.int64_t(size),
			(*C.H3Index)(unsafe.Pointer(&out[0]))))
	}

	// Simple case: polygon without holes
	result := C.polygonToCellsExperimental_c_wrapper(
		(*C.LatLng)(unsafe.Pointer(&mainLoop[0])),
		C.int(len(mainLoop)),
		nil, nil, 0,
		C.int(res),
		C.uint32_t(flags),
		C.int64_t(size),
		(*C.H3Index)(unsafe.Pointer(&out[0])))

	return h3Error(result)
}

// cellToBBoxC calls the C cellToBBox function
func cellToBBoxC(cell h3Index, coverChildren bool) (bbox, h3Error) {
	var cBBox C.BBox
	err := h3Error(C.cellToBBox(C.H3Index(cell), &cBBox, C.bool(coverChildren)))

	goBBox := bbox{
		North: Angle(cBBox.north),
		South: Angle(cBBox.south),
		East:  Angle(cBBox.east),
		West:  Angle(cBBox.west),
	}

	return goBBox, err
}

// bboxToCellBoundaryC calls the C bboxToCellBoundary function
func bboxToCellBoundaryC(bbox *bbox) CellBoundary {
	cBBox := C.BBox{
		north: C.double(bbox.North),
		south: C.double(bbox.South),
		east:  C.double(bbox.East),
		west:  C.double(bbox.West),
	}

	cBoundary := C.bboxToCellBoundary(&cBBox)

	// Convert C CellBoundary to Go CellBoundary
	boundary := CellBoundary{numVerts: int32(cBoundary.numVerts)}

	for i := int32(0); i < boundary.numVerts; i++ {
		boundary.verts[i] = LatLng{
			Lat: Angle(cBoundary.verts[i].lat),
			Lng: Angle(cBoundary.verts[i].lng),
		}
	}

	return boundary
}

// maxPolygonToCellsSizeExperimentalC calls the C maxPolygonToCellsSizeExperimental function
func maxPolygonToCellsSizeExperimentalC(polygon *GeoPolygon, res int32, flags uint32) (int64, h3Error) {
	if len(polygon.GeoLoop) == 0 {
		return 0, eSuccess
	}

	// Convert Go GeoPolygon to C GeoPolygon
	mainLoop := polygon.GeoLoop

	// For now, simplify and only handle polygons without holes
	// to avoid CGO pointer issues
	var out C.int64_t
	var result C.H3Error

	if len(polygon.Holes) > 0 {
		// TODO: Implement hole handling once basic case works
		result = C.maxPolygonToCellsSizeExperimental_c_wrapper(
			(*C.LatLng)(unsafe.Pointer(&mainLoop[0])),
			C.int(len(mainLoop)),
			nil, nil, 0,
			C.int(res), C.uint32_t(flags), &out)
	} else {
		// Simple case: polygon without holes
		result = C.maxPolygonToCellsSizeExperimental_c_wrapper(
			(*C.LatLng)(unsafe.Pointer(&mainLoop[0])),
			C.int(len(mainLoop)),
			nil, nil, 0,
			C.int(res), C.uint32_t(flags), &out)
	}

	return int64(out), h3Error(result)
}
