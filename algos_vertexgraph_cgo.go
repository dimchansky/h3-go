//go:build cgo && c2go && !h3v450

package h3

// Wrappers for the vertex-graph half of algos.c that exists only in the
// H3 4.4.0 tree: h3SetToVertexGraph and _vertexGraphToLinkedGeo (and
// vertexGraph.h itself) were deleted in 4.5.0 with the cellsToMultiPolygon
// rewrite (docs/sync/4.4.0-to-4.5.0.md §5.2). This file and its parity
// tests retire with I-C; docs/sync/h3v450-exclusion-inventory.md tracks
// the exclusion.

/*
// Interop wrapper only; the original C sources are compiled via separate
// build-tagged C shim files in this package (see h3lib_*.c with //go:build c2go).
#include <stdint.h>
#include "h3api.h"
#include "algos.h"
#include "vertexGraph.h"

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
*/
import "C"

// _vertexGraphToLinkedGeoC calls the original C implementation.
// This function is complex to test due to the need to convert entire graph structures.
// The parity test will focus on specific behavior verification.
func _vertexGraphToLinkedGeoC(graph *C.VertexGraph, out *C.LinkedGeoPolygon) {
	C._vertexGraphToLinkedGeo_c_wrapper(graph, out)
}

// h3SetToVertexGraphC calls the original C implementation.
func h3SetToVertexGraphC(h3Set []h3Index, graph *C.VertexGraph) h3Error {
	if len(h3Set) == 0 {
		return eSuccess
	}
	return h3Error(C.h3SetToVertexGraph_c_wrapper(
		(*C.H3Index)(&h3Set[0]),
		C.int(len(h3Set)),
		graph,
	))
}

// vertexGraphCResult holds the basic properties of a C vertexGraph for parity testing.
type vertexGraphCResult struct {
	Err        h3Error
	Size       int32
	NumBuckets int32
	Res        int32
}

// h3SetToVertexGraphCForParity calls the C implementation and returns basic graph properties.
// Used for parity testing without exposing C types to test files.
func h3SetToVertexGraphCForParity(h3Set []h3Index) vertexGraphCResult {
	if len(h3Set) == 0 {
		return vertexGraphCResult{Err: eSuccess, Size: 0, NumBuckets: 0, Res: 0}
	}

	var cGraph C.VertexGraph
	err := h3Error(C.h3SetToVertexGraph_c_wrapper(
		(*C.H3Index)(&h3Set[0]),
		C.int(len(h3Set)),
		&cGraph,
	))

	result := vertexGraphCResult{
		Err: err,
	}

	// Only read properties and clean up if the operation succeeded
	if err == eSuccess {
		result.Size = int32(cGraph.size)
		result.NumBuckets = int32(cGraph.numBuckets)
		result.Res = int32(cGraph.res)
		C.destroyVertexGraph(&cGraph)
	}

	return result
}
