//go:build cgo && c2go

package h3

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

// cellToVertexC wraps the C cellToVertex function.
func cellToVertexC(cell H3Index, vertexNum int32) (H3Index, H3Error) {
	var out C.H3Index
	err := H3Error(C.cellToVertex(C.H3Index(cell), C.int(vertexNum), &out))
	return H3Index(out), err
}

// cellToVertexesC wraps the C cellToVertexes function.
func cellToVertexesC(cell H3Index, vertexes *[6]H3Index) H3Error {
	var cVertexes [6]C.H3Index
	err := H3Error(C.cellToVertexes(C.H3Index(cell), &cVertexes[0]))
	if err == E_SUCCESS {
		for i := 0; i < 6; i++ {
			(*vertexes)[i] = H3Index(cVertexes[i])
		}
	}
	return err
}

// vertexToLatLngC wraps the C vertexToLatLng function.
func vertexToLatLngC(vertex H3Index, coord *LatLng) H3Error {
	var cCoord C.LatLng
	err := H3Error(C.vertexToLatLng(C.H3Index(vertex), &cCoord))
	if err == E_SUCCESS {
		coord.Lat = Rad(float64(cCoord.lat))
		coord.Lng = Rad(float64(cCoord.lng))
	}
	return err
}

// isValidVertexC wraps the C isValidVertex function.
func isValidVertexC(vertex H3Index) bool {
	return C.isValidVertex(C.H3Index(vertex)) != 0
}
