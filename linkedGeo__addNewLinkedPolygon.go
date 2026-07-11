package h3

// addNewLinkedPolygon adds a new linked polygon to the current polygon.
// The function creates a new linkedGeoPolygon, links it to the provided polygon
// as the next element, and returns a pointer to the new polygon. This is used
// to extend linked lists of polygons in multi-polygon structures.
// Ported from H3 C: linkedGeo.c::addNewLinkedPolygon.
func addNewLinkedPolygon(polygon *linkedGeoPolygon) *linkedGeoPolygon {
	// Assert that polygon.next is nil (equivalent to assert(polygon->next == NULL))
	if polygon.Next != nil {
		panic("addNewLinkedPolygon: polygon.Next must be nil")
	}

	// Create a new polygon (equivalent to H3_MEMORY(calloc)(1, sizeof(*next)))
	next := &linkedGeoPolygon{
		First: nil,
		Last:  nil,
		Next:  nil,
	}

	// Link the new polygon (equivalent to polygon->next = next)
	polygon.Next = next

	// Return pointer to the new polygon
	return next
}
