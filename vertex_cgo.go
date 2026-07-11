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
func vertexRotationsC(cell h3Index, out *int32) h3Error {
	var cOut C.int
	err := h3Error(C.vertexRotations_wrapper(C.H3Index(cell), &cOut))
	*out = int32(cOut)
	return err
}

// vertexNumForDirectionC wraps the C vertexNumForDirection function.
func vertexNumForDirectionC(origin h3Index, direction direction) int32 {
	return int32(C.vertexNumForDirection(C.H3Index(origin), C.Direction(direction)))
}

// directionForVertexNumC wraps the C directionForVertexNum function.
func directionForVertexNumC(origin h3Index, vertexNum int32) direction {
	return direction(C.directionForVertexNum(C.H3Index(origin), C.int(vertexNum)))
}

// cellToVertexC wraps the C cellToVertex function.
func cellToVertexC(cell h3Index, vertexNum int32) (h3Index, h3Error) {
	var out C.H3Index
	err := h3Error(C.cellToVertex(C.H3Index(cell), C.int(vertexNum), &out))
	return h3Index(out), err
}

// cellToVertexesC wraps the C cellToVertexes function.
func cellToVertexesC(cell h3Index, vertexes *[6]h3Index) h3Error {
	var cVertexes [6]C.H3Index
	err := h3Error(C.cellToVertexes(C.H3Index(cell), &cVertexes[0]))
	if err == eSuccess {
		for i := 0; i < 6; i++ {
			(*vertexes)[i] = h3Index(cVertexes[i])
		}
	}
	return err
}

// vertexToLatLngC wraps the C vertexToLatLng function.
func vertexToLatLngC(vertex h3Index, coord *LatLng) h3Error {
	var cCoord C.LatLng
	err := h3Error(C.vertexToLatLng(C.H3Index(vertex), &cCoord))
	if err == eSuccess {
		coord.Lat = Rad(float64(cCoord.lat))
		coord.Lng = Rad(float64(cCoord.lng))
	}
	return err
}

// isValidVertexC wraps the C isValidVertex function.
func isValidVertexC(vertex h3Index) bool {
	return C.isValidVertex(C.H3Index(vertex)) != 0
}
