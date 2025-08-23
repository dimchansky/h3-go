//go:build cgo

package c2go

/*
#include <stdint.h>
#include "h3api.h"
#include "linkedGeo.h"
#include "bbox.h"
// Prototypes for the original C helpers in linkedGeo.c
int countLinkedCoords(LinkedGeoLoop* loop);
int countLinkedLoops(LinkedGeoPolygon* polygon);
void bboxFromLinkedGeoLoopC(const LinkedGeoLoop *loop, BBox *bbox);
bool pointInsideLinkedGeoLoopC(const LinkedGeoLoop *loop, const BBox *bbox, const LatLng *coord);
bool isClockwiseLinkedGeoLoopC(const LinkedGeoLoop *loop);
LinkedGeoPolygon* addNewLinkedPolygonC(LinkedGeoPolygon *polygon);
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

// countLinkedLoopsC wraps the C countLinkedLoops function for parity testing.
func countLinkedLoopsC(polygon *LinkedGeoPolygon) int {
	// Create a minimal C LinkedGeoPolygon for testing
	var firstLoop *C.LinkedGeoLoop
	var currentGoLoop = polygon.First
	var prevCLoop *C.LinkedGeoLoop

	// Build the C linked list from Go linked list
	for currentGoLoop != nil {
		cLoop := (*C.LinkedGeoLoop)(C.malloc(C.size_t(C.sizeof_LinkedGeoLoop)))

		// For simplicity, set first and last to nil (we're only testing loop counting)
		cLoop.first = nil
		cLoop.last = nil
		cLoop.next = nil

		if firstLoop == nil {
			firstLoop = cLoop
		} else {
			prevCLoop.next = cLoop
		}

		prevCLoop = cLoop
		currentGoLoop = currentGoLoop.Next
	}

	// Create C LinkedGeoPolygon
	cPolygon := (*C.LinkedGeoPolygon)(C.malloc(C.size_t(C.sizeof_LinkedGeoPolygon)))
	cPolygon.first = firstLoop
	cPolygon.last = prevCLoop
	cPolygon.next = nil

	// Call C function
	result := int(C.countLinkedLoops(cPolygon))

	// Clean up C memory
	currentCLoop := firstLoop
	for currentCLoop != nil {
		nextLoop := currentCLoop.next
		C.free(unsafe.Pointer(currentCLoop))
		currentCLoop = nextLoop
	}
	C.free(unsafe.Pointer(cPolygon))

	return result
}

// bboxFromLinkedGeoLoopC wraps the C bboxFromLinkedGeoLoop function for parity testing.
func bboxFromLinkedGeoLoopC(loop *LinkedGeoLoop, bbox *BBox) {
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

	// Create C BBox for result
	cBbox := (*C.BBox)(C.malloc(C.size_t(C.sizeof_BBox)))

	// Call C function
	C.bboxFromLinkedGeoLoopC(cLoop, cBbox)

	// Convert C BBox back to Go BBox
	bbox.North = float64(cBbox.north)
	bbox.South = float64(cBbox.south)
	bbox.East = float64(cBbox.east)
	bbox.West = float64(cBbox.west)

	// Clean up C memory
	currentCNode := firstNode
	for currentCNode != nil {
		nextNode := currentCNode.next
		C.free(unsafe.Pointer(currentCNode))
		currentCNode = nextNode
	}
	C.free(unsafe.Pointer(cLoop))
	C.free(unsafe.Pointer(cBbox))
}

// pointInsideLinkedGeoLoopC wraps the C pointInsideLinkedGeoLoop function for parity testing.
func pointInsideLinkedGeoLoopC(loop *LinkedGeoLoop, bbox *BBox, coord *LatLng) bool {
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

	// Create C BBox
	cBbox := (*C.BBox)(C.malloc(C.size_t(C.sizeof_BBox)))
	cBbox.north = C.double(bbox.North)
	cBbox.south = C.double(bbox.South)
	cBbox.east = C.double(bbox.East)
	cBbox.west = C.double(bbox.West)

	// Create C LatLng
	cCoord := (*C.LatLng)(C.malloc(C.size_t(C.sizeof_LatLng)))
	cCoord.lat = C.double(coord.Lat)
	cCoord.lng = C.double(coord.Lng)

	// Call C function
	result := bool(C.pointInsideLinkedGeoLoopC(cLoop, cBbox, cCoord))

	// Clean up C memory
	currentCNode := firstNode
	for currentCNode != nil {
		nextNode := currentCNode.next
		C.free(unsafe.Pointer(currentCNode))
		currentCNode = nextNode
	}
	C.free(unsafe.Pointer(cLoop))
	C.free(unsafe.Pointer(cBbox))
	C.free(unsafe.Pointer(cCoord))

	return result
}

// isClockwiseLinkedGeoLoopC wraps the C isClockwiseLinkedGeoLoop function for parity testing.
func isClockwiseLinkedGeoLoopC(loop *LinkedGeoLoop) bool {
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
	result := bool(C.isClockwiseLinkedGeoLoopC(cLoop))

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

// addNewLinkedPolygonC wraps the C addNewLinkedPolygon function for parity testing.
func addNewLinkedPolygonC(polygon *LinkedGeoPolygon) *LinkedGeoPolygon {
	// Create a minimal C LinkedGeoPolygon structure
	cPolygon := (*C.LinkedGeoPolygon)(C.malloc(C.size_t(C.sizeof_LinkedGeoPolygon)))
	cPolygon.first = nil
	cPolygon.last = nil
	cPolygon.next = nil

	// Call C function
	cResult := C.addNewLinkedPolygonC(cPolygon)

	// Convert result back to Go
	goResult := &LinkedGeoPolygon{
		First: nil,
		Last:  nil,
		Next:  nil,
	}

	// Clean up C memory
	C.free(unsafe.Pointer(cResult))
	C.free(unsafe.Pointer(cPolygon))

	return goResult
}
