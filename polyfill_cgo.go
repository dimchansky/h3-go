//go:build cgo && c2go

package h3

/*
#cgo CPPFLAGS: -I${SRCDIR}/testref/h3-4.3.0/src/h3lib/include

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
*/
import "C"
import (
	"unsafe"
)

// polygonToCellsExperimentalC calls the C polygonToCellsExperimental function
func polygonToCellsExperimentalC(polygon *GeoPolygon, res int32, flags uint32,
	size int64, out []H3Index) H3Error {
	if len(out) == 0 || int64(len(out)) < size {
		return E_MEMORY_BOUNDS
	}

	// Use a simple approach similar to other polygon functions
	// Convert to C structures and call directly

	// Build C GeoLoop for main loop
	mainLoop := polygon.GeoLoop
	if len(mainLoop) == 0 {
		// Empty polygon case
		return H3Error(C.polygonToCellsExperimental_c_wrapper(
			nil, 0, nil, nil, 0, C.int(res), C.uint32_t(flags),
			C.int64_t(size), (*C.H3Index)(unsafe.Pointer(&out[0]))))
	}

	// For now, simplify and only handle polygons without holes
	// to avoid CGO pointer issues
	if len(polygon.Holes) > 0 {
		// TODO: Implement hole handling once basic case works
		return H3Error(C.polygonToCellsExperimental_c_wrapper(
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

	return H3Error(result)
}

// cellToBBoxC calls the C cellToBBox function
func cellToBBoxC(cell H3Index, coverChildren bool) (BBox, H3Error) {
	var cBBox C.BBox
	err := H3Error(C.cellToBBox(C.H3Index(cell), &cBBox, C.bool(coverChildren)))

	goBBox := BBox{
		North: Angle(cBBox.north),
		South: Angle(cBBox.south),
		East:  Angle(cBBox.east),
		West:  Angle(cBBox.west),
	}

	return goBBox, err
}

// bboxToCellBoundaryC calls the C bboxToCellBoundary function
func bboxToCellBoundaryC(bbox *BBox) CellBoundary {
	cBBox := C.BBox{
		north: C.double(bbox.North),
		south: C.double(bbox.South),
		east:  C.double(bbox.East),
		west:  C.double(bbox.West),
	}

	cBoundary := C.bboxToCellBoundary(&cBBox)

	// Convert C CellBoundary to Go CellBoundary
	boundary := CellBoundary{
		NumVerts: int32(cBoundary.numVerts),
		Verts:    make([]LatLng, int32(cBoundary.numVerts)),
	}

	for i := int32(0); i < boundary.NumVerts; i++ {
		boundary.Verts[i] = LatLng{
			Lat: Angle(cBoundary.verts[i].lat),
			Lng: Angle(cBoundary.verts[i].lng),
		}
	}

	return boundary
}
