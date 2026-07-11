package h3

// addLinkedCoord adds a new coordinate to the current loop's linked list.
// This function allocates a new LinkedLatLng node, copies the vertex data,
// and manages the loop's first/last pointers to maintain the linked list structure.
// If the loop is empty, the coordinate becomes both first and last.
// Otherwise, it's appended after the current last coordinate.
// Ported from H3 C: linkedGeo.c::addLinkedCoord.
func addLinkedCoord(loop *LinkedGeoLoop, vertex *LatLng) *LinkedLatLng {
	// Create a new coordinate node (equivalent to H3_MEMORY(malloc)(sizeof(*coord)))
	coord := &LinkedLatLng{
		Vertex: *vertex, // Copy the vertex value
		Next:   nil,
	}

	// Get the current last coordinate
	last := loop.Last
	if last == nil {
		// Assert that loop.first is also nil (equivalent to assert(loop->first == NULL))
		if loop.First != nil {
			panic("addLinkedCoord: loop.First must be nil when loop.Last is nil")
		}
		loop.First = coord
	} else {
		last.Next = coord
	}
	loop.Last = coord
	return coord
}
