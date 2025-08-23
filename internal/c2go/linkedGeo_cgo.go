//go:build cgo

package c2go

/*
#include <stdint.h>
#include "h3api.h"
#include "linkedGeo.h"
// Prototype for the original C helper in linkedGeo.c
int countLinkedCoords(LinkedGeoLoop* loop);
*/
import "C"
import "unsafe"

// countLinkedCoordsC wraps the C countLinkedCoords function for parity testing.
func countLinkedCoordsC(loop *LinkedGeoLoop) int {
	// Create a minimal C LinkedGeoLoop for testing
	var firstNode *C.LinkedLatLng
	var currentGoCoord = loop.First
	var prevCNode *C.LinkedLatLng

	// Build the C linked list from Go linked list
	for currentGoCoord != nil {
		cNode := (*C.LinkedLatLng)(C.malloc(C.size_t(C.sizeof_LinkedLatLng)))
		cNode.vertex.lat = C.double(currentGoCoord.Vertex.Lat)
		cNode.vertex.lng = C.double(currentGoCoord.Vertex.Lng)
		cNode.next = nil

		if firstNode == nil {
			firstNode = cNode
		} else {
			prevCNode.next = cNode
		}

		prevCNode = cNode
		currentGoCoord = currentGoCoord.Next
	}

	// Create C LinkedGeoLoop
	cLoop := (*C.LinkedGeoLoop)(C.malloc(C.size_t(C.sizeof_LinkedGeoLoop)))
	cLoop.first = firstNode
	cLoop.last = prevCNode
	cLoop.next = nil

	// Call C function
	result := int(C.countLinkedCoords(cLoop))

	// Clean up C memory
	currentCNode := firstNode
	for currentCNode != nil {
		nextNode := currentCNode.next
		C.free(unsafe.Pointer(currentCNode))
		currentCNode = nextNode
	}
	C.free(unsafe.Pointer(cLoop))

	return result
}
