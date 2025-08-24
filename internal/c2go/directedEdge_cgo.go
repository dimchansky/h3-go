//go:build cgo

package c2go

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
*/
import "C"

// getDirectedEdgeOriginC calls the original C implementation.
func getDirectedEdgeOriginC(edge H3Index) (H3Index, H3Error) {
	var out C.H3Index
	err := H3Error(C.getDirectedEdgeOrigin_c_wrapper(C.H3Index(edge), &out))
	return H3Index(out), err
}

// isValidDirectedEdgeC calls the original C implementation.
func isValidDirectedEdgeC(edge H3Index) bool {
	return C.isValidDirectedEdge_c_wrapper(C.H3Index(edge)) != 0
}

// originToDirectedEdgesC calls the original C implementation.
func originToDirectedEdgesC(origin H3Index, edges []H3Index) H3Error {
	// Convert to C array
	cEdges := make([]C.H3Index, 6)
	err := H3Error(C.originToDirectedEdges_c_wrapper(C.H3Index(origin), &cEdges[0]))
	// Copy back to Go slice
	for i := 0; i < 6; i++ {
		edges[i] = H3Index(cEdges[i])
	}
	return err
}

// getDirectedEdgeDestinationC calls the original C implementation.
func getDirectedEdgeDestinationC(edge H3Index) (H3Index, H3Error) {
	var out C.H3Index
	err := H3Error(C.getDirectedEdgeDestination_c_wrapper(C.H3Index(edge), &out))
	return H3Index(out), err
}
