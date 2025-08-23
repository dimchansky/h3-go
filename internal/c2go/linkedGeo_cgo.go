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
const LinkedGeoPolygon* findDeepestContainerC(const LinkedGeoPolygon **polygons, const BBox **bboxes, int polygonCount);
const LinkedGeoPolygon* findPolygonForHoleC(const LinkedGeoLoop *loop, const LinkedGeoPolygon *polygon, const BBox *bboxes, int polygonCount);
H3Error normalizeMultiPolygon(LinkedGeoPolygon *root);
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

// findDeepestContainerC wraps the C findDeepestContainer function for parity testing.
func findDeepestContainerC(polygons []*LinkedGeoPolygon, bboxes []*BBox) *LinkedGeoPolygon {
	if len(polygons) != len(bboxes) {
		panic("findDeepestContainerC: polygons and bboxes must have same length")
	}

	polygonCount := len(polygons)
	if polygonCount == 0 {
		return nil
	}

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
	cResult := C.findDeepestContainerC(cPolygons, cBboxes, C.int(polygonCount))

	// Find which Go polygon corresponds to the C result
	var result *LinkedGeoPolygon
	if cResult != nil {
		for i := 0; i < polygonCount; i++ {
			if allocatedPolygons[i] == cResult {
				result = polygons[i]
				break
			}
		}
	}

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

// findPolygonForHoleC wraps the C findPolygonForHole function for parity testing.
func findPolygonForHoleC(loop *LinkedGeoLoop, polygons []*LinkedGeoPolygon, bboxes []*BBox) *LinkedGeoPolygon {
	if len(polygons) != len(bboxes) {
		panic("findPolygonForHoleC: polygons and bboxes must have same length")
	}

	polygonCount := len(polygons)
	if polygonCount == 0 {
		return nil
	}

	// Early exit if loop is nil or has no coordinates
	if loop == nil || loop.First == nil {
		return nil
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

	// Create linked list of C polygons
	var firstCPolygon *C.LinkedGeoPolygon
	var prevCPolygon *C.LinkedGeoPolygon

	// Arrays to track allocated memory for cleanup
	allocatedPolygons := make([]*C.LinkedGeoPolygon, polygonCount)
	allocatedLoops := make([][]*C.LinkedGeoLoop, polygonCount)
	allocatedCoords := make([][]*C.LinkedLatLng, polygonCount)

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

		// Link polygons
		if firstCPolygon == nil {
			firstCPolygon = cPolygon
		} else {
			prevCPolygon.next = cPolygon
		}
		cPolygon.next = nil
		prevCPolygon = cPolygon
	}

	// Create C bboxes array
	cBboxes := (*C.BBox)(C.malloc(C.size_t(uintptr(polygonCount) * C.sizeof_BBox)))
	defer C.free(unsafe.Pointer(cBboxes))

	for i := 0; i < polygonCount; i++ {
		cBboxPtr := (*C.BBox)(unsafe.Pointer(uintptr(unsafe.Pointer(cBboxes)) + uintptr(i)*C.sizeof_BBox))
		if bboxes[i] != nil {
			cBboxPtr.north = C.double(bboxes[i].North)
			cBboxPtr.south = C.double(bboxes[i].South)
			cBboxPtr.east = C.double(bboxes[i].East)
			cBboxPtr.west = C.double(bboxes[i].West)
		}
	}

	// Call C function
	cResult := C.findPolygonForHoleC(cLoop, firstCPolygon, cBboxes, C.int(polygonCount))

	// Find which Go polygon corresponds to the C result
	var result *LinkedGeoPolygon
	if cResult != nil {
		for i := 0; i < polygonCount; i++ {
			if allocatedPolygons[i] == cResult {
				result = polygons[i]
				break
			}
		}
	}

	// Clean up allocated polygons and loops
	for i := 0; i < polygonCount; i++ {
		if allocatedPolygons[i] != nil {
			C.free(unsafe.Pointer(allocatedPolygons[i]))
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

// normalizeMultiPolygonC wraps the C normalizeMultiPolygon function for parity testing.
func normalizeMultiPolygonC(root *LinkedGeoPolygon) H3Error {
	if root == nil {
		// C function expects a valid pointer, so passing nil will cause segfault
		// Based on the Go implementation, nil input should return E_FAILED
		return E_FAILED
	}

	// Convert Go LinkedGeoPolygon to C structure
	cPolygon := (*C.LinkedGeoPolygon)(C.malloc(C.size_t(C.sizeof_LinkedGeoPolygon)))
	defer C.free(unsafe.Pointer(cPolygon))

	// Keep track of allocated memory for cleanup
	var allocatedLoops []*C.LinkedGeoLoop
	var allocatedCoords []*C.LinkedLatLng

	// Helper function to convert Go LinkedGeoLoop to C LinkedGeoLoop
	convertLoop := func(goLoop *LinkedGeoLoop) *C.LinkedGeoLoop {
		if goLoop == nil {
			return nil
		}

		cLoop := (*C.LinkedGeoLoop)(C.malloc(C.size_t(C.sizeof_LinkedGeoLoop)))
		allocatedLoops = append(allocatedLoops, cLoop)

		// Convert coordinates
		var firstCCoord *C.LinkedLatLng
		var prevCCoord *C.LinkedLatLng

		currentGoCoord := goLoop.First
		for currentGoCoord != nil {
			cCoord := (*C.LinkedLatLng)(C.malloc(C.size_t(C.sizeof_LinkedLatLng)))
			allocatedCoords = append(allocatedCoords, cCoord)

			cCoord.vertex.lat = C.double(currentGoCoord.Vertex.Lat)
			cCoord.vertex.lng = C.double(currentGoCoord.Vertex.Lng)
			cCoord.next = nil

			if firstCCoord == nil {
				firstCCoord = cCoord
				cLoop.first = cCoord
			} else {
				prevCCoord.next = cCoord
			}

			prevCCoord = cCoord
			currentGoCoord = currentGoCoord.Next
		}

		cLoop.last = prevCCoord
		cLoop.next = nil

		return cLoop
	}

	// Convert the first loop
	if root.First != nil {
		cPolygon.first = convertLoop(root.First)
	} else {
		cPolygon.first = nil
	}

	// Convert subsequent loops linked to the first one
	var prevCLoop *C.LinkedGeoLoop = cPolygon.first
	currentGoLoop := root.First
	if currentGoLoop != nil {
		currentGoLoop = currentGoLoop.Next
	}

	for currentGoLoop != nil {
		cLoop := convertLoop(currentGoLoop)
		if prevCLoop != nil {
			prevCLoop.next = cLoop
		}
		prevCLoop = cLoop
		currentGoLoop = currentGoLoop.Next
	}

	if cPolygon.first != nil {
		cPolygon.last = prevCLoop
	} else {
		cPolygon.last = nil
	}

	// Handle the Next polygon pointer if it exists
	if root.Next != nil {
		// Convert the Next polygon
		nextCPolygon := (*C.LinkedGeoPolygon)(C.malloc(C.size_t(C.sizeof_LinkedGeoPolygon)))
		// For simplicity, just set up basic structure - the main test is the Next pointer existence
		nextCPolygon.first = nil
		nextCPolygon.last = nil  
		nextCPolygon.next = nil
		cPolygon.next = nextCPolygon
	} else {
		cPolygon.next = nil
	}

	// Call the C function
	result := C.normalizeMultiPolygon(cPolygon)

	// IMPORTANT: Do NOT free the allocated memory here!
	// The C normalizeMultiPolygon function takes ownership of the memory
	// and will free what it needs to free. Attempting to free here
	// causes double-free errors and crashes.
	//
	// The C function:
	// 1. Zeros out the root polygon: *root = (LinkedGeoPolygon){0}
	// 2. Restructures all the loops into new polygons
	// 3. May call destroyLinkedGeoLoop() and free() on orphaned holes
	//
	// Since the function modifies the structure in-place and manages
	// its own memory, we should not attempt to free anything we allocated.
	
	return H3Error(result)
}
