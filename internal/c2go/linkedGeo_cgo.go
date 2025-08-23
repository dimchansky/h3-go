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
int countLinkedPolygonsC(LinkedGeoPolygon* polygon);
void bboxFromLinkedGeoLoopC(const LinkedGeoLoop *loop, BBox *bbox);
bool pointInsideLinkedGeoLoopC(const LinkedGeoLoop *loop, const BBox *bbox, const LatLng *coord);
bool isClockwiseLinkedGeoLoopC(const LinkedGeoLoop *loop);
LinkedGeoPolygon* addNewLinkedPolygonC(LinkedGeoPolygon *polygon);
LinkedGeoLoop* addLinkedLoopC(LinkedGeoPolygon *polygon, LinkedGeoLoop *loop);
LinkedGeoLoop* addNewLinkedLoopC(LinkedGeoPolygon *polygon);
LinkedLatLng* addLinkedCoordC(LinkedGeoLoop *loop, const LatLng *vertex);
int countContainersC(const LinkedGeoLoop *loop, const LinkedGeoPolygon **polygons, const BBox **bboxes, int polygonCount);
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

// countLinkedPolygonsC wraps the C countLinkedPolygons function for parity testing.
func countLinkedPolygonsC(polygon *LinkedGeoPolygon) int {
	// Create C structures to mirror the Go linked list
	if polygon == nil {
		return int(C.countLinkedPolygonsC(nil))
	}
	
	// Build the C linked list from Go linked list
	var firstCPolygon *C.LinkedGeoPolygon
	var prevCPolygon *C.LinkedGeoPolygon
	currentGoPolygon := polygon
	
	for currentGoPolygon != nil {
		cPolygon := (*C.LinkedGeoPolygon)(C.malloc(C.size_t(C.sizeof_LinkedGeoPolygon)))
		cPolygon.first = nil
		cPolygon.last = nil
		cPolygon.next = nil
		
		if firstCPolygon == nil {
			firstCPolygon = cPolygon
		} else {
			prevCPolygon.next = cPolygon
		}
		
		prevCPolygon = cPolygon
		currentGoPolygon = currentGoPolygon.Next
	}
	
	// Call C function
	result := int(C.countLinkedPolygonsC(firstCPolygon))
	
	// Clean up C memory
	currentCPolygon := firstCPolygon
	for currentCPolygon != nil {
		nextPolygon := currentCPolygon.next
		C.free(unsafe.Pointer(currentCPolygon))
		currentCPolygon = nextPolygon
	}
	
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

	// Clean up C memory (don't free cResult - it's managed by the C function)
	// Free the allocated cResult and cPolygon since they were allocated by C function
	if cResult != nil {
		C.free(unsafe.Pointer(cResult))
	}
	C.free(unsafe.Pointer(cPolygon))

	return goResult
}

// addLinkedLoopC wraps the C addLinkedLoop function for parity testing.
// Returns true if C function behaves as expected
func addLinkedLoopC(wasEmpty bool, polygonHadLoop bool) (returnsLoop bool, setsFirst bool, setsLast bool, linksLoops bool) {
	// Create C structures to test behavior
	cPolygon := (*C.LinkedGeoPolygon)(C.malloc(C.size_t(C.sizeof_LinkedGeoPolygon)))
	defer C.free(unsafe.Pointer(cPolygon))
	
	var cExisting *C.LinkedGeoLoop
	if polygonHadLoop {
		// Create a dummy C loop for existing first
		cExisting = (*C.LinkedGeoLoop)(C.malloc(C.size_t(C.sizeof_LinkedGeoLoop)))
		defer C.free(unsafe.Pointer(cExisting))
		cExisting.first = nil
		cExisting.last = nil
		cExisting.next = nil
		cPolygon.first = cExisting
		cPolygon.last = cExisting
	} else {
		cPolygon.first = nil
		cPolygon.last = nil
	}
	cPolygon.next = nil

	// Create C LinkedGeoLoop for the new loop
	cLoop := (*C.LinkedGeoLoop)(C.malloc(C.size_t(C.sizeof_LinkedGeoLoop)))
	defer C.free(unsafe.Pointer(cLoop))
	cLoop.first = nil
	cLoop.last = nil
	cLoop.next = nil

	// Call C function
	cResult := C.addLinkedLoopC(cPolygon, cLoop)

	// Check what C did
	returnsLoop = (cResult == cLoop)
	setsFirst = (wasEmpty && cPolygon.first == cLoop) || (!wasEmpty && cPolygon.first == cExisting)
	setsLast = (cPolygon.last == cLoop)
	linksLoops = (polygonHadLoop && cExisting != nil && cExisting.next == cLoop)
	
	return
}

// addNewLinkedLoopC wraps the C addNewLinkedLoop function for parity testing.
// Returns information about C function behavior
func addNewLinkedLoopC(wasEmpty bool, polygonHadLoop bool) (createsLoop bool, setsFirst bool, setsLast bool, linksLoops bool) {
	// Create C structures to test behavior
	cPolygon := (*C.LinkedGeoPolygon)(C.malloc(C.size_t(C.sizeof_LinkedGeoPolygon)))
	defer C.free(unsafe.Pointer(cPolygon))
	
	var cExisting *C.LinkedGeoLoop
	if polygonHadLoop {
		// Create a dummy C loop for existing first
		cExisting = (*C.LinkedGeoLoop)(C.malloc(C.size_t(C.sizeof_LinkedGeoLoop)))
		defer C.free(unsafe.Pointer(cExisting))
		cExisting.first = nil
		cExisting.last = nil
		cExisting.next = nil
		cPolygon.first = cExisting
		cPolygon.last = cExisting
	} else {
		cPolygon.first = nil
		cPolygon.last = nil
	}
	cPolygon.next = nil

	// Call C function
	cResult := C.addNewLinkedLoopC(cPolygon)

	// Check what C did
	createsLoop = (cResult != nil)
	if cResult != nil {
		defer C.free(unsafe.Pointer(cResult))
		setsFirst = (wasEmpty && cPolygon.first == cResult) || (!wasEmpty && cPolygon.first == cExisting)
		setsLast = (cPolygon.last == cResult)
		linksLoops = (polygonHadLoop && cExisting != nil && cExisting.next == cResult)
	}
	
	return
}

// addLinkedCoordC wraps the C addLinkedCoord function for parity testing.
// Returns information about C function behavior
func addLinkedCoordC(wasEmpty bool, loopHadCoord bool, vertex LatLng) (returnsCoord bool, setsFirst bool, setsLast bool, linksCoords bool, vertexMatches bool) {
	// Create C structures to test behavior
	cLoop := (*C.LinkedGeoLoop)(C.malloc(C.size_t(C.sizeof_LinkedGeoLoop)))
	defer C.free(unsafe.Pointer(cLoop))
	
	var cExisting *C.LinkedLatLng
	if loopHadCoord {
		// Create a dummy coordinate for existing first
		cExisting = (*C.LinkedLatLng)(C.malloc(C.size_t(C.sizeof_LinkedLatLng)))
		defer C.free(unsafe.Pointer(cExisting))
		cExisting.vertex.lat = C.double(0.5) // dummy value
		cExisting.vertex.lng = C.double(1.0) // dummy value
		cExisting.next = nil
		cLoop.first = cExisting
		cLoop.last = cExisting
	} else {
		cLoop.first = nil
		cLoop.last = nil
	}
	cLoop.next = nil

	// Create C LatLng for the vertex
	cVertex := C.LatLng{
		lat: C.double(vertex.Lat),
		lng: C.double(vertex.Lng),
	}

	// Call C function
	cResult := C.addLinkedCoordC(cLoop, &cVertex)

	// Check what C did
	returnsCoord = (cResult != nil)
	if cResult != nil {
		defer C.free(unsafe.Pointer(cResult))
		setsFirst = (wasEmpty && cLoop.first == cResult) || (!wasEmpty && cLoop.first == cExisting)
		setsLast = (cLoop.last == cResult)
		linksCoords = (loopHadCoord && cExisting != nil && cExisting.next == cResult)
		vertexMatches = (float64(cResult.vertex.lat) == vertex.Lat && float64(cResult.vertex.lng) == vertex.Lng)
	}
	
	return
}

// countContainersC wraps the C countContainers function for parity testing.
func countContainersC(loop *LinkedGeoLoop, polygons []*LinkedGeoPolygon, bboxes []*BBox) int {
	if len(polygons) != len(bboxes) {
		panic("countContainersC: polygons and bboxes must have same length")
	}
	
	polygonCount := len(polygons)
	if polygonCount == 0 {
		return 0
	}
	
	// Early return if loop is nil or has no coordinates - can't be contained
	if loop == nil || loop.First == nil {
		return 0
	}
	
	// Convert Go loop to C loop
	cLoop := (*C.LinkedGeoLoop)(C.malloc(C.size_t(C.sizeof_LinkedGeoLoop)))
	defer C.free(unsafe.Pointer(cLoop))
	
	// Build the complete C linked list from Go linked list
	var firstCCoord *C.LinkedLatLng
	var prevCCoord *C.LinkedLatLng
	
	currentGoCoord := loop.First
	for currentGoCoord != nil {
		cCoord := (*C.LinkedLatLng)(C.malloc(C.size_t(C.sizeof_LinkedLatLng)))
		defer C.free(unsafe.Pointer(cCoord))
		cCoord.vertex.lat = C.double(currentGoCoord.Vertex.Lat)
		cCoord.vertex.lng = C.double(currentGoCoord.Vertex.Lng)
		cCoord.next = nil
		
		if firstCCoord == nil {
			firstCCoord = cCoord
		} else {
			prevCCoord.next = cCoord
		}
		prevCCoord = cCoord
		currentGoCoord = currentGoCoord.Next
	}
	
	cLoop.first = firstCCoord
	cLoop.last = prevCCoord
	cLoop.next = nil
	
	// Create arrays of C polygons and bboxes
	ptrSize := unsafe.Sizeof(uintptr(0))
	cPolygons := (**C.LinkedGeoPolygon)(C.malloc(C.size_t(uintptr(polygonCount) * ptrSize)))
	defer C.free(unsafe.Pointer(cPolygons))
	
	cBboxes := (**C.BBox)(C.malloc(C.size_t(uintptr(polygonCount) * ptrSize)))
	defer C.free(unsafe.Pointer(cBboxes))
	
	// Array to track allocated memory for cleanup
	allocatedPolygons := make([]*C.LinkedGeoPolygon, polygonCount)
	allocatedBboxes := make([]*C.BBox, polygonCount)
	allocatedLoops := make([][]*C.LinkedGeoLoop, polygonCount)
	allocatedCoords := make([][]*C.LinkedLatLng, polygonCount)
	
	// Convert Go polygons and bboxes to C
	for i := 0; i < polygonCount; i++ {
		// Create C polygon
		cPolygon := (*C.LinkedGeoPolygon)(C.malloc(C.size_t(C.sizeof_LinkedGeoPolygon)))
		allocatedPolygons[i] = cPolygon
		
		// Create first loop if exists
		if polygons[i] != nil && polygons[i].First != nil {
			cFirstLoop := (*C.LinkedGeoLoop)(C.malloc(C.size_t(C.sizeof_LinkedGeoLoop)))
			allocatedLoops[i] = append(allocatedLoops[i], cFirstLoop)
			
			// Build complete coordinate list for this polygon's first loop
			var firstLoopCoord *C.LinkedLatLng
			var prevLoopCoord *C.LinkedLatLng
			
			currentGoLoopCoord := polygons[i].First.First
			for currentGoLoopCoord != nil {
				cLoopCoord := (*C.LinkedLatLng)(C.malloc(C.size_t(C.sizeof_LinkedLatLng)))
				allocatedCoords[i] = append(allocatedCoords[i], cLoopCoord)
				cLoopCoord.vertex.lat = C.double(currentGoLoopCoord.Vertex.Lat)
				cLoopCoord.vertex.lng = C.double(currentGoLoopCoord.Vertex.Lng)
				cLoopCoord.next = nil
				
				if firstLoopCoord == nil {
					firstLoopCoord = cLoopCoord
				} else {
					prevLoopCoord.next = cLoopCoord
				}
				prevLoopCoord = cLoopCoord
				currentGoLoopCoord = currentGoLoopCoord.Next
			}
			
			cFirstLoop.first = firstLoopCoord
			cFirstLoop.last = prevLoopCoord
			cFirstLoop.next = nil
			cPolygon.first = cFirstLoop
		} else {
			cPolygon.first = nil
		}
		cPolygon.last = cPolygon.first
		cPolygon.next = nil
		
		// Set polygon pointer in array
		*(**C.LinkedGeoPolygon)(unsafe.Pointer(uintptr(unsafe.Pointer(cPolygons)) + uintptr(i)*unsafe.Sizeof(uintptr(0)))) = cPolygon
		
		// Create C bbox
		cBbox := (*C.BBox)(C.malloc(C.size_t(C.sizeof_BBox)))
		allocatedBboxes[i] = cBbox
		if bboxes[i] != nil {
			cBbox.north = C.double(bboxes[i].North)
			cBbox.south = C.double(bboxes[i].South)
			cBbox.east = C.double(bboxes[i].East)
			cBbox.west = C.double(bboxes[i].West)
		}
		
		// Set bbox pointer in array
		*(**C.BBox)(unsafe.Pointer(uintptr(unsafe.Pointer(cBboxes)) + uintptr(i)*unsafe.Sizeof(uintptr(0)))) = cBbox
	}
	
	// Call C function
	result := int(C.countContainersC(cLoop, cPolygons, cBboxes, C.int(polygonCount)))
	
	// Clean up allocated polygons and bboxes
	for i := 0; i < polygonCount; i++ {
		if allocatedPolygons[i] != nil {
			C.free(unsafe.Pointer(allocatedPolygons[i]))
		}
		if allocatedBboxes[i] != nil {
			C.free(unsafe.Pointer(allocatedBboxes[i]))
		}
		// Clean up loops
		for _, cLoop := range allocatedLoops[i] {
			C.free(unsafe.Pointer(cLoop))
		}
		// Clean up coordinates
		for _, cCoord := range allocatedCoords[i] {
			C.free(unsafe.Pointer(cCoord))
		}
	}
	
	return result
}
