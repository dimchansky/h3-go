//go:build cgo

package c2go

/*
// Interop wrapper only; the original C sources are compiled via separate
// build-tagged C shim files in this package (see h3lib_*.c with //go:build c2go).
#include <stdint.h>
#include <stdbool.h>
#include "h3api.h"
#include "algos.h"

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
