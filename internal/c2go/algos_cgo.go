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
