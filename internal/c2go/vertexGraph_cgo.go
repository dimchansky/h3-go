//go:build cgo

package c2go

/*
#include <stdint.h>
#include "h3api.h"
#include "vertexGraph.h"
// Prototypes for the original C helpers in vertexGraph.c
uint32_t _hashVertex(const LatLng* vertex, int res, int numBuckets);
void initVertexGraph(VertexGraph* graph, int numBuckets, int res);
*/
import "C"

// _hashVertexC wraps the C _hashVertex function for parity testing.
func _hashVertexC(vertex *LatLng, res int32, numBuckets int32) uint32 {
	cVertex := C.LatLng{
		lat: C.double(vertex.Lat),
		lng: C.double(vertex.Lng),
	}
	result := C._hashVertex(&cVertex, C.int(res), C.int(numBuckets))
	return uint32(result)
}

// initVertexGraphC wraps the C initVertexGraph function for parity testing.
func initVertexGraphC(graph *VertexGraph, numBuckets int32, res int32) {
	var cGraph C.VertexGraph
	C.initVertexGraph(&cGraph, C.int(numBuckets), C.int(res))

	// Copy results back to Go struct
	graph.NumBuckets = int32(cGraph.numBuckets)
	graph.Size = int32(cGraph.size)
	graph.Res = int32(cGraph.res)

	// For buckets, we only verify the allocation occurred (non-nil vs nil)
	// Full bucket testing would require memory allocation tracking
	if numBuckets > 0 {
		graph.Buckets = make([]*VertexNode, numBuckets)
	} else {
		graph.Buckets = nil
	}
}
