//go:build cgo

package c2go

/*
#include <stdint.h>
#include "h3api.h"
// Prototype for the original C helper in vertexGraph.c
uint32_t _hashVertex(const LatLng* vertex, int res, int numBuckets);
*/
import "C"

// _hashVertexC wraps the C _hashVertex function for parity testing.
func _hashVertexC(vertex *LatLng, res int, numBuckets int) uint32 {
	cVertex := C.LatLng{
		lat: C.double(vertex.Lat),
		lng: C.double(vertex.Lng),
	}
	result := C._hashVertex(&cVertex, C.int(res), C.int(numBuckets))
	return uint32(result)
}
