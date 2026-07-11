package h3

// addLinkedLoop adds an existing linked loop to the current polygon.
// This function manages the polygon's first/last pointers to maintain the
// linked list structure. If the polygon is empty, the loop becomes both first
// and last. Otherwise, it's appended after the current last loop.
// Ported from H3 C: linkedGeo.c::addLinkedLoop.
func addLinkedLoop(polygon *linkedGeoPolygon, loop *linkedGeoLoop) *linkedGeoLoop {
	last := polygon.Last
	if last == nil {
		// Assert that polygon.first is also nil (equivalent to assert(polygon->first == NULL))
		if polygon.First != nil {
			panic("addLinkedLoop: polygon.First must be nil when polygon.Last is nil")
		}
		polygon.First = loop
	} else {
		last.Next = loop
	}
	polygon.Last = loop
	return loop
}

// addNewLinkedLoop creates a new linked loop and adds it to the current polygon.
// The function allocates a new linkedGeoLoop with zero-initialized fields,
// then uses addLinkedLoop to properly link it into the polygon's loop chain.
// This is a convenience function that combines allocation and linking.
// Ported from H3 C: linkedGeo.c::addNewLinkedLoop.
func addNewLinkedLoop(polygon *linkedGeoPolygon) *linkedGeoLoop {
	// Create a new loop (equivalent to H3_MEMORY(calloc)(1, sizeof(*loop)))
	loop := &linkedGeoLoop{
		First: nil,
		Last:  nil,
		Next:  nil,
	}

	// Use addLinkedLoop to properly link it into the polygon
	return addLinkedLoop(polygon, loop)
}
