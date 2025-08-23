package c2go

// countLinkedCoords counts the number of coordinates in a linked geo loop.
// Iterates through the linked list of coordinates, counting each node.
// Returns the total count of coordinates in the loop.
// Ported from H3 C: linkedGeo.c::countLinkedCoords
func countLinkedCoords(loop *LinkedGeoLoop) int {
	coord := loop.First
	count := 0
	for coord != nil {
		count++
		coord = coord.Next
	}
	return count
}
