package h3

// linkedGeoLoopToGeoLoop converts a LinkedGeoLoop to a GeoLoop by
// counting coords and copying. Returns eSuccess, eFailed (< 3 verts),
// or eMemoryAlloc (unreachable in Go: allocation cannot fail).
// Ported from H3 C: linkedGeo.c::linkedGeoLoopToGeoLoop.
func linkedGeoLoopToGeoLoop(linked *linkedGeoLoop, out *GeoLoop) h3Error {
	numVerts := countLinkedCoords(linked)
	if numVerts < 3 {
		return eFailed
	}
	verts := make([]LatLng, numVerts)

	coord := linked.First
	for i := int32(0); coord != nil && i < numVerts; i++ {
		// ALWAYS(i < numVerts) in C.
		verts[i] = coord.Vertex
		coord = coord.Next
	}
	*out = GeoLoop(verts)
	return eSuccess
}
