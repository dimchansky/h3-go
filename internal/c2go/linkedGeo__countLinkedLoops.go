package c2go

// countLinkedLoops counts the number of linked loops in a polygon.
// Iterates through the linked list of loops, counting each node.
// Returns the total count of loops in the polygon.
// Ported from H3 C: linkedGeo.c::countLinkedLoops
func countLinkedLoops(polygon *LinkedGeoPolygon) int {
	loop := polygon.First
	count := 0
	for loop != nil {
		count++
		loop = loop.Next
	}
	return count
}
