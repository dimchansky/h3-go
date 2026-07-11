package h3

// countLinkedCoords counts the number of coordinates in a linked geo loop.
// Iterates through the linked list of coordinates, counting each node.
// Returns the total count of coordinates in the loop.
// Ported from H3 C: linkedGeo.c::countLinkedCoords.
func countLinkedCoords(loop *linkedGeoLoop) int32 {
	coord := loop.First
	count := int32(0)
	for coord != nil {
		count++
		coord = coord.Next
	}
	return count
}
