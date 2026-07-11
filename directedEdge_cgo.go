//go:build cgo && c2go

package h3

/*
// Interop wrapper only; the original C sources are compiled via separate
// build-tagged C shim files in this package (see h3lib_*.c with //go:build c2go).
#include <stdint.h>
#include <stdbool.h>
#include "h3api.h"
#include "directedEdge.h"
// Normalize C bool to int for cgo comparisons when needed (toolchain-safe)
static int h3_bool_to_int(_Bool b) { return b ? 1 : 0; }

// Wrapper function to call getDirectedEdgeOrigin
static H3Error getDirectedEdgeOrigin_c_wrapper(H3Index edge, H3Index *out) {
    return getDirectedEdgeOrigin(edge, out);
}

// Wrapper function to call isValidDirectedEdge
static int isValidDirectedEdge_c_wrapper(H3Index edge) {
    return isValidDirectedEdge(edge);
}

// Wrapper function to call originToDirectedEdges
static H3Error originToDirectedEdges_c_wrapper(H3Index origin, H3Index *edges) {
    return originToDirectedEdges(origin, edges);
}

// Wrapper function to call getDirectedEdgeDestination
static H3Error getDirectedEdgeDestination_c_wrapper(H3Index edge, H3Index *out) {
    return getDirectedEdgeDestination(edge, out);
}

// Wrapper function to call directedEdgeToBoundary
static H3Error directedEdgeToBoundary_c_wrapper(H3Index edge, CellBoundary *cb) {
    return directedEdgeToBoundary(edge, cb);
}

// Wrapper function to call directedEdgeToCells
static H3Error directedEdgeToCells_c_wrapper(H3Index edge, H3Index *originDestination) {
    return directedEdgeToCells(edge, originDestination);
}

// Wrapper function to call cellsToDirectedEdge
static H3Error cellsToDirectedEdge_c_wrapper(H3Index origin, H3Index destination, H3Index *out) {
    return cellsToDirectedEdge(origin, destination, out);
}

// Wrapper function to call areNeighborCells
static H3Error areNeighborCells_c_wrapper(H3Index origin, H3Index destination, int *out) {
    return areNeighborCells(origin, destination, out);
}
*/
import "C"

// getDirectedEdgeOriginC calls the original C implementation.
func getDirectedEdgeOriginC(edge h3Index) (h3Index, h3Error) {
	var out C.H3Index
	err := h3Error(C.getDirectedEdgeOrigin_c_wrapper(C.H3Index(edge), &out))
	return h3Index(out), err
}

// isValidDirectedEdgeC calls the original C implementation.
func isValidDirectedEdgeC(edge h3Index) bool {
	return C.isValidDirectedEdge_c_wrapper(C.H3Index(edge)) != 0
}

// originToDirectedEdgesC calls the original C implementation.
func originToDirectedEdgesC(origin h3Index, edges []h3Index) h3Error {
	// Convert to C array
	cEdges := make([]C.H3Index, 6)
	err := h3Error(C.originToDirectedEdges_c_wrapper(C.H3Index(origin), &cEdges[0]))
	// Copy back to Go slice
	for i := 0; i < 6; i++ {
		edges[i] = h3Index(cEdges[i])
	}
	return err
}

// getDirectedEdgeDestinationC calls the original C implementation.
func getDirectedEdgeDestinationC(edge h3Index) (h3Index, h3Error) {
	var out C.H3Index
	err := h3Error(C.getDirectedEdgeDestination_c_wrapper(C.H3Index(edge), &out))
	return h3Index(out), err
}

// directedEdgeToBoundaryC calls the original C implementation.
func directedEdgeToBoundaryC(edge h3Index, cb *CellBoundary) h3Error {
	var cCb C.CellBoundary
	cCb.numVerts = 0

	err := h3Error(C.directedEdgeToBoundary_c_wrapper(C.H3Index(edge), &cCb))

	if err == eSuccess {
		// Copy results back to Go struct
		cb.numVerts = int32(cCb.numVerts)

		// Copy vertices from C to Go
		for i := int32(0); i < cb.numVerts; i++ {
			cb.verts[i].Lat = Rad(float64(cCb.verts[i].lat))
			cb.verts[i].Lng = Rad(float64(cCb.verts[i].lng))
		}
	}

	return err
}

// directedEdgeToCellsC calls the original C implementation.
func directedEdgeToCellsC(edge h3Index, originDestination []h3Index) h3Error {
	// Convert to C array
	var cCells [2]C.H3Index
	err := h3Error(C.directedEdgeToCells_c_wrapper(C.H3Index(edge), &cCells[0]))
	// Copy back to Go slice
	if err == eSuccess {
		originDestination[0] = h3Index(cCells[0])
		originDestination[1] = h3Index(cCells[1])
	}
	return err
}

// cellsToDirectedEdgeC calls the original C implementation.
func cellsToDirectedEdgeC(origin, destination h3Index) (h3Index, h3Error) {
	var out C.H3Index
	err := h3Error(C.cellsToDirectedEdge_c_wrapper(C.H3Index(origin), C.H3Index(destination), &out))
	return h3Index(out), err
}

// areNeighborCellsC calls the original C implementation.
func areNeighborCellsC(origin, destination h3Index) (bool, h3Error) {
	var out C.int
	err := h3Error(C.areNeighborCells_c_wrapper(C.H3Index(origin), C.H3Index(destination), &out))
	return out != 0, err
}
