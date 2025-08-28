//go:build cgo

package h3

/*
// Interop wrapper only; the original C sources are compiled via separate
// build-tagged C shim files in this package (see h3lib_*.c with //go:build c2go).
#include <stdint.h>
#include <stdbool.h>
#include "h3api.h"
#include "algos.h"
#include "vertexGraph.h"

// Wrapper function to call h3NeighborRotations
static H3Error h3NeighborRotations_c_wrapper(H3Index origin, Direction dir, int *rotations, H3Index *out) {
    return h3NeighborRotations(origin, dir, rotations, out);
}

// Wrapper function to call directionForNeighbor
static Direction directionForNeighbor_c_wrapper(H3Index origin, H3Index destination) {
    return directionForNeighbor(origin, destination);
}

// Wrapper function to call _gridDiskDistancesInternal
static H3Error _gridDiskDistancesInternal_c_wrapper(H3Index origin, int k, H3Index *out,
                                                    int *distances, int64_t maxIdx, int curK) {
    return _gridDiskDistancesInternal(origin, k, out, distances, maxIdx, curK);
}

// Wrapper function to call maxGridDiskSize
static H3Error maxGridDiskSize_c_wrapper(int k, int64_t *out) {
    return maxGridDiskSize(k, out);
}

// Wrapper function to call gridDiskDistancesUnsafe
static H3Error gridDiskDistancesUnsafe_c_wrapper(H3Index origin, int k, H3Index *out, int *distances) {
    return gridDiskDistancesUnsafe(origin, k, out, distances);
}

// Wrapper function to call gridDiskDistances
static H3Error gridDiskDistances_c_wrapper(H3Index origin, int k, H3Index *out, int *distances) {
    return gridDiskDistances(origin, k, out, distances);
}

// Wrapper function to call gridDisk
static H3Error gridDisk_c_wrapper(H3Index origin, int k, H3Index *out) {
    return gridDisk(origin, k, out);
}

// Wrapper function to call gridDiskUnsafe
static H3Error gridDiskUnsafe_c_wrapper(H3Index origin, int k, H3Index *out) {
    return gridDiskUnsafe(origin, k, out);
}

// Wrapper function to call gridDiskDistancesSafe
static H3Error gridDiskDistancesSafe_c_wrapper(H3Index origin, int k, H3Index *out, int *distances) {
    return gridDiskDistancesSafe(origin, k, out, distances);
}

// Wrapper function to call maxGridRingSize
static H3Error maxGridRingSize_c_wrapper(int k, int64_t *out) {
    return maxGridRingSize(k, out);
}

// Wrapper function to call gridRingUnsafe
static H3Error gridRingUnsafe_c_wrapper(H3Index origin, int k, H3Index *out) {
    return gridRingUnsafe(origin, k, out);
}

// Wrapper function to call _gridRingInternal
static H3Error _gridRingInternal_c_wrapper(H3Index origin, int k, H3Index *out) {
    return _gridRingInternal(origin, k, out);
}

// Wrapper function to call gridRing
static H3Error gridRing_c_wrapper(H3Index origin, int k, H3Index *out) {
    return gridRing(origin, k, out);
}

// Wrapper function to call gridDisksUnsafe
static H3Error gridDisksUnsafe_c_wrapper(H3Index *h3Set, int length, int k, H3Index *out) {
    return gridDisksUnsafe(h3Set, length, k, out);
}

// Wrapper function to call _getEdgeHexagons
static H3Error _getEdgeHexagons_c_wrapper(const GeoLoop *geoloop, int64_t numHexagons, int res,
                                           int64_t *numSearchHexes, H3Index *search, H3Index *found) {
    return _getEdgeHexagons(geoloop, numHexagons, res, numSearchHexes, search, found);
}

// Wrapper function to call maxPolygonToCellsSize
static H3Error maxPolygonToCellsSize_c_wrapper(const GeoPolygon *geoPolygon, int res, uint32_t flags, int64_t *out) {
    return maxPolygonToCellsSize(geoPolygon, res, flags, out);
}

// Wrapper function to call polygonToCells
static H3Error polygonToCells_c_wrapper(const GeoPolygon *geoPolygon, int res, uint32_t flags, H3Index *out) {
    return polygonToCells(geoPolygon, res, flags, out);
}

// Forward declaration for _vertexGraphToLinkedGeo
void _vertexGraphToLinkedGeo(VertexGraph *graph, LinkedGeoPolygon *out);

// Wrapper function to call _vertexGraphToLinkedGeo
static void _vertexGraphToLinkedGeo_c_wrapper(VertexGraph *graph, LinkedGeoPolygon *out) {
    _vertexGraphToLinkedGeo(graph, out);
}

// Wrapper function to call h3SetToVertexGraph
static H3Error h3SetToVertexGraph_c_wrapper(const H3Index *h3Set, const int numHexes, VertexGraph *graph) {
    return h3SetToVertexGraph(h3Set, numHexes, graph);
}

// Wrapper function to call cellsToLinkedMultiPolygon
static H3Error cellsToLinkedMultiPolygon_c_wrapper(const H3Index *h3Set, const int numHexes, LinkedGeoPolygon *out) {
    return cellsToLinkedMultiPolygon(h3Set, numHexes, out);
}
*/
import "C"

// h3NeighborRotationsC calls the original C implementation.
func h3NeighborRotationsC(origin H3Index, dir Direction, rotations *int32) (H3Index, H3Error) {
	var out C.H3Index
	cRotations := C.int(*rotations)
	err := H3Error(C.h3NeighborRotations_c_wrapper(C.H3Index(origin), C.Direction(dir), &cRotations, &out))
	*rotations = int32(cRotations) // Update the rotations value
	return H3Index(out), err
}

// directionForNeighborC calls the original C implementation.
func directionForNeighborC(origin H3Index, destination H3Index) Direction {
	return Direction(C.directionForNeighbor_c_wrapper(C.H3Index(origin), C.H3Index(destination)))
}

// _gridDiskDistancesInternalC calls the original C implementation.
func _gridDiskDistancesInternalC(origin H3Index, k int32, out []H3Index, distances []int32, maxIdx int64, curK int32) H3Error {
	if len(out) == 0 || len(distances) == 0 {
		return E_FAILED
	}
	return H3Error(C._gridDiskDistancesInternal_c_wrapper(
		C.H3Index(origin),
		C.int(k),
		(*C.H3Index)(&out[0]),
		(*C.int)(&distances[0]),
		C.int64_t(maxIdx),
		C.int(curK),
	))
}

// maxGridDiskSizeC calls the original C implementation.
func maxGridDiskSizeC(k int32, out *int64) H3Error {
	return H3Error(C.maxGridDiskSize_c_wrapper(C.int(k), (*C.int64_t)(out)))
}

// gridDiskDistancesUnsafeC calls the original C implementation.
func gridDiskDistancesUnsafeC(origin H3Index, k int32, out []H3Index, distances []int32) H3Error {
	if len(out) == 0 {
		return E_FAILED
	}
	var distPtr *C.int
	if len(distances) > 0 {
		distPtr = (*C.int)(&distances[0])
	}
	return H3Error(C.gridDiskDistancesUnsafe_c_wrapper(
		C.H3Index(origin),
		C.int(k),
		(*C.H3Index)(&out[0]),
		distPtr,
	))
}

// gridDiskDistancesC calls the original C implementation.
func gridDiskDistancesC(origin H3Index, k int32, out []H3Index, distances []int32) H3Error {
	if len(out) == 0 {
		return E_FAILED
	}
	var distPtr *C.int
	if len(distances) > 0 {
		distPtr = (*C.int)(&distances[0])
	}
	return H3Error(C.gridDiskDistances_c_wrapper(
		C.H3Index(origin),
		C.int(k),
		(*C.H3Index)(&out[0]),
		distPtr,
	))
}

// gridDiskC calls the original C implementation.
func gridDiskC(origin H3Index, k int32, out []H3Index) H3Error {
	if len(out) == 0 {
		return E_FAILED
	}
	return H3Error(C.gridDisk_c_wrapper(
		C.H3Index(origin),
		C.int(k),
		(*C.H3Index)(&out[0]),
	))
}

// gridDiskUnsafeC calls the original C implementation.
func gridDiskUnsafeC(origin H3Index, k int32, out []H3Index) H3Error {
	if len(out) == 0 {
		return E_FAILED
	}
	return H3Error(C.gridDiskUnsafe_c_wrapper(
		C.H3Index(origin),
		C.int(k),
		(*C.H3Index)(&out[0]),
	))
}

// gridDiskDistancesSafeC calls the original C implementation.
func gridDiskDistancesSafeC(origin H3Index, k int32, out []H3Index, distances []int32) H3Error {
	if len(out) == 0 {
		return E_FAILED
	}
	var distPtr *C.int
	if len(distances) > 0 {
		distPtr = (*C.int)(&distances[0])
	}
	return H3Error(C.gridDiskDistancesSafe_c_wrapper(
		C.H3Index(origin),
		C.int(k),
		(*C.H3Index)(&out[0]),
		distPtr,
	))
}

// maxGridRingSizeC calls the original C implementation.
func maxGridRingSizeC(k int32, out *int64) H3Error {
	return H3Error(C.maxGridRingSize_c_wrapper(C.int(k), (*C.int64_t)(out)))
}

// gridRingUnsafeC calls the original C implementation.
func gridRingUnsafeC(origin H3Index, k int32, out []H3Index) H3Error {
	if len(out) == 0 {
		return E_FAILED
	}
	return H3Error(C.gridRingUnsafe_c_wrapper(
		C.H3Index(origin),
		C.int(k),
		(*C.H3Index)(&out[0]),
	))
}

// _gridRingInternalC calls the original C implementation.
func _gridRingInternalC(origin H3Index, k int32, out []H3Index) H3Error {
	if len(out) == 0 {
		return E_FAILED
	}
	return H3Error(C._gridRingInternal_c_wrapper(
		C.H3Index(origin),
		C.int(k),
		(*C.H3Index)(&out[0]),
	))
}

// gridRingC calls the original C implementation.
func gridRingC(origin H3Index, k int32, out []H3Index) H3Error {
	if len(out) == 0 {
		return E_FAILED
	}
	return H3Error(C.gridRing_c_wrapper(
		C.H3Index(origin),
		C.int(k),
		(*C.H3Index)(&out[0]),
	))
}

// gridDisksUnsafeC calls the original C implementation.
func gridDisksUnsafeC(h3Set []H3Index, k int32, out []H3Index) H3Error {
	if len(h3Set) == 0 || len(out) == 0 {
		return E_FAILED
	}
	return H3Error(C.gridDisksUnsafe_c_wrapper(
		(*C.H3Index)(&h3Set[0]),
		C.int(len(h3Set)),
		C.int(k),
		(*C.H3Index)(&out[0]),
	))
}

// _getEdgeHexagonsC calls the original C implementation.
func _getEdgeHexagonsC(geoloop []LatLng, numHexagons int64, res int32, numSearchHexes *int64, search []H3Index, found []H3Index) H3Error {
	if len(search) == 0 || len(found) == 0 {
		return E_FAILED
	}

	cGeoloop, freeFn := toCGeoLoop(geoloop)
	defer freeFn()

	cNumSearchHexes := C.int64_t(*numSearchHexes)
	err := H3Error(C._getEdgeHexagons_c_wrapper(
		&cGeoloop,
		C.int64_t(numHexagons),
		C.int(res),
		&cNumSearchHexes,
		(*C.H3Index)(&search[0]),
		(*C.H3Index)(&found[0]),
	))
	*numSearchHexes = int64(cNumSearchHexes)
	return err
}

// maxPolygonToCellsSizeC calls the original C implementation.
func maxPolygonToCellsSizeC(geoPolygon *GeoPolygon, res int32, flags uint32, out *int64) H3Error {
	cGeoPolygon, freeFn := toCGeoPolygon(*geoPolygon)
	defer freeFn()

	return H3Error(C.maxPolygonToCellsSize_c_wrapper(
		&cGeoPolygon,
		C.int(res),
		C.uint32_t(flags),
		(*C.int64_t)(out),
	))
}

// polygonToCellsC calls the original C implementation.
func polygonToCellsC(geoPolygon *GeoPolygon, res int32, flags uint32, out []H3Index) H3Error {
	if len(out) == 0 {
		return E_FAILED
	}

	cGeoPolygon, freeFn := toCGeoPolygon(*geoPolygon)
	defer freeFn()

	return H3Error(C.polygonToCells_c_wrapper(
		&cGeoPolygon,
		C.int(res),
		C.uint32_t(flags),
		(*C.H3Index)(&out[0]),
	))
}

// _vertexGraphToLinkedGeoC calls the original C implementation.
// This function is complex to test due to the need to convert entire graph structures.
// The parity test will focus on specific behavior verification.
func _vertexGraphToLinkedGeoC(graph *C.VertexGraph, out *C.LinkedGeoPolygon) {
	C._vertexGraphToLinkedGeo_c_wrapper(graph, out)
}

// h3SetToVertexGraphC calls the original C implementation.
func h3SetToVertexGraphC(h3Set []H3Index, graph *C.VertexGraph) H3Error {
	if len(h3Set) == 0 {
		return E_SUCCESS
	}
	return H3Error(C.h3SetToVertexGraph_c_wrapper(
		(*C.H3Index)(&h3Set[0]),
		C.int(len(h3Set)),
		graph,
	))
}

// cellsToLinkedMultiPolygonC calls the original C implementation.
func cellsToLinkedMultiPolygonC(h3Set []H3Index, out *C.LinkedGeoPolygon) H3Error {
	if len(h3Set) == 0 {
		return E_SUCCESS
	}
	return H3Error(C.cellsToLinkedMultiPolygon_c_wrapper(
		(*C.H3Index)(&h3Set[0]),
		C.int(len(h3Set)),
		out,
	))
}
