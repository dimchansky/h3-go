package c2go

// destroyLinkedGeoLoop frees all allocated memory for coordinates in a linked geo loop.
// The caller is responsible for freeing memory allocated to input loop struct.
// Mirrors the behavior of the C destroyLinkedGeoLoop function.
// Ported from H3 C: linkedGeo.c::destroyLinkedGeoLoop
func destroyLinkedGeoLoop(loop *LinkedGeoLoop) {
	if loop == nil {
		return
	}

	var nextCoord *LinkedLatLng
	for currentCoord := loop.First; currentCoord != nil; currentCoord = nextCoord {
		nextCoord = currentCoord.Next
		// In Go, we don't need to explicitly free memory - garbage collection handles it
		// We just need to break the references to help the GC
		currentCoord.Next = nil
	}

	// Clear the loop's coordinate references
	loop.First = nil
	loop.Last = nil
}
