package h3

// countLinkedPolygons counts the number of polygons in a linked list.
// It traverses the linked list starting from the given polygon and counts
// each polygon until it reaches the end (nil).
// Ported from H3 C: linkedGeo.c::countLinkedPolygons
func countLinkedPolygons(polygon *LinkedGeoPolygon) int32 {
	count := int32(0)
	for polygon != nil {
		count++
		polygon = polygon.Next
	}
	return count
}
