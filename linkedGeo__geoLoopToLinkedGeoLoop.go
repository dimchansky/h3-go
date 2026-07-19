package h3

// geoLoopToLinkedGeoLoop populates a LinkedGeoLoop with vertices from a
// GeoLoop. loop must be zeroed (C: calloc-zeroed). Returns eSuccess or
// eFailed (< 3 verts); eMemoryAlloc is unreachable in Go.
// Ported from H3 C: linkedGeo.c::geoLoopToLinkedGeoLoop.
func geoLoopToLinkedGeoLoop(src GeoLoop, loop *linkedGeoLoop) h3Error {
	if len(src) < 3 {
		return eFailed
	}
	for i := 0; i < len(src); i++ {
		coord := &linkedLatLng{Vertex: src[i], Next: nil}
		if loop.Last == nil {
			loop.First = coord
		} else {
			loop.Last.Next = coord
		}
		loop.Last = coord
	}
	return eSuccess
}
