package c2go

// destroyLinkedMultiPolygon frees all allocated memory for a linked geo structure.
// The caller is responsible for freeing memory allocated to input polygon struct.
// The function frees all polygons in the linked list except the first (input) polygon,
// and all loops and coordinates within each polygon.
// Ported from H3 C: linkedGeo.c::destroyLinkedMultiPolygon
func destroyLinkedMultiPolygon(polygon *LinkedGeoPolygon) {
	if polygon == nil {
		return
	}

	// flag to skip the input polygon
	skip := true
	var nextPolygon *LinkedGeoPolygon
	var nextLoop *LinkedGeoLoop

	for currentPolygon := polygon; currentPolygon != nil; currentPolygon = nextPolygon {
		for currentLoop := currentPolygon.First; currentLoop != nil; currentLoop = nextLoop {
			destroyLinkedGeoLoop(currentLoop)
			nextLoop = currentLoop.Next
			// In Go, we don't need to explicitly free memory - garbage collection handles it
			// We just need to break the references to help the GC
			currentLoop.Next = nil
		}
		nextPolygon = currentPolygon.Next
		if skip {
			// do not free the input polygon - just clear its references
			skip = false
			currentPolygon.First = nil
			currentPolygon.Last = nil
		} else {
			// Clear references for non-input polygons
			currentPolygon.First = nil
			currentPolygon.Last = nil
			currentPolygon.Next = nil
		}
	}

	// Clear the input polygon's next reference
	polygon.Next = nil
}
