package h3

// createSortableLoop creates, starting from a given Arc, a SortableLoop
// that contains that Arc. SortableLoops are sorted by the root (i.e.,
// connected component) and then by the area contained by the loop. We
// use this to merge all loops in a connected component into a single
// polygon. We use the area values to determine which loop will be the
// "outer" loop of the polygon.
// Ported from H3 C: cellsToMultiPoly.c::createSortableLoop.
func createSortableLoop(arc *arc, sloop *sortableLoop) h3Error {
	var gb CellBoundary
	start := arc.id

	var numVerts int64
	var verts []LatLng

	numVerts = 0
	for {
		// This is an overestimate for numVerts.
		// Most cell edges will just need one vert (we don't use the last
		// vertex in the edge).
		// For even resolutions, all cell edges need just one vert.
		// Over-allocate for now; shrink to actual number later (C: realloc).
		numVerts += 2
		arc = arc.next
		if arc.id == start {
			break
		}
	}

	verts = make([]LatLng, numVerts)

	numVerts = 0
	j := int64(0)
	for {
		err := directedEdgeToBoundary(arc.id, &gb)
		if err != eSuccess {
			// NEVER in C.
			return err
		}

		for i := int64(0); i < int64(gb.numVerts)-1; i++ {
			verts[j] = gb.verts[i]
			j++
		}
		numVerts += int64(gb.numVerts) - 1
		arc.isVisited = true
		arc = arc.next
		if arc.id == start {
			break
		}
	}

	// This memory ends up in GeoMultiPolygon, to be freed by caller of
	// cellsToMultiPolygon() (C: realloc down to the actual count).
	verts = verts[:numVerts]

	sloop.root = getRoot(arc).id
	sloop.loop = GeoLoop(verts)
	sloop.area, _ = geoLoopAreaRads2(sloop.loop)

	return eSuccess
}
