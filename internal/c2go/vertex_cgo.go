//go:build cgo && c2go

package c2go

/*
#cgo CFLAGS: -DH3_HAVE_VLA
#include <stdint.h>
#include "h3api.h"
#include "vertex.h"

// External wrapper to access the static vertexRotations function
extern H3Error vertexRotations_wrapper(H3Index cell, int *out);
*/
import "C"

// vertexRotationsC wraps the C vertexRotations function.
func vertexRotationsC(cell H3Index, out *int32) H3Error {
	var cOut C.int
	err := H3Error(C.vertexRotations_wrapper(C.H3Index(cell), &cOut))
	*out = int32(cOut)
	return err
}

// vertexNumForDirectionC wraps the C vertexNumForDirection function.
func vertexNumForDirectionC(origin H3Index, direction Direction) int32 {
	return int32(C.vertexNumForDirection(C.H3Index(origin), C.Direction(direction)))
}

// directionForVertexNumC wraps the C directionForVertexNum function.
func directionForVertexNumC(origin H3Index, vertexNum int32) Direction {
	return Direction(C.directionForVertexNum(C.H3Index(origin), C.int(vertexNum)))
}
