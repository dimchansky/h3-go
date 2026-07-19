package h3

// cellsToMultiPolygon creates a GeoMultiPolygon from a set of H3 cells.
//
// This function converts a set of H3 cells into a GeoMultiPolygon
// representing the region they cover. Note the difference with
// cellsToLinkedMultiPolygon, which returns a linked-list
// LinkedGeoPolygon. A GeoMultiPolygon provides the sizes of its
// elements and supports direct indexing.
//
// Polygons follow the right hand rule, with the outer loop oriented
// counter-clockwise, and the inner loops oriented clockwise.
//
// Polygons within a GeoMultiPolygon are ordered by decreasing area of
// the outer loop.
//
// Note that for polygons with multiple loops (one outer loop + at least
// one hole), *any* loop can serve as the outer loop and still produce
// the *same* valid polygon. We use the convention of choosing as the
// outer loop the one that would give the largest area "outside" of that
// outer loop. This results in what users would probably expect: a
// polygon for the land within a state/province with a large lake would
// have the outer loop be the state's boundary, instead of the lake's
// boundary.
//
// cells must be valid cells at the same resolution with no duplicates;
// violations fail with the validateCellSet errors. The C function is
// binding-only/tentative upstream, so it stays internal in this port
// (docs/sync/4.4.0-to-4.5.0.md §13.2).
// Ported from H3 C: cellsToMultiPoly.c::cellsToMultiPolygon.
func cellsToMultiPolygon(cells []h3Index, numCells int64, out *geoMultiPolygon) h3Error {
	err := checkCellsToMultiPolyOverflow(numCells, hashTableMultiplier)
	if err != eSuccess {
		return err
	}

	err = validateCellSet(cells, numCells)
	if err != eSuccess {
		return err
	}

	if numCells == 0 {
		out.NumPolygons = 0
		out.Polygons = nil
		return eSuccess
	}

	// arcset initializes with separate doubly-linked loops for each cell,
	// each in their own connected component
	var arcset arcSet
	err = createArcSet(cells, numCells, &arcset)
	if err != eSuccess {
		return err
	}

	// Cancel out pairs of edges, updating the doubly-linked loops and merging
	// them into a single connected component
	err = cancelArcPairs(arcset)
	if err != eSuccess {
		// NEVER in C.
		destroyArcSet(&arcset)
		return err
	}

	/*
	   Extract all loops and sort them by:
	     1) their connected component, and then by
	     2) the loop area.
	   This makes loops for each polygon contiguous in memory.
	   Within each polygon, the sorting makes the loop with the smallest
	   enclosed area come first (accounting for loop winding direction),
	   which is what we take to be the outer loop for that polygon.
	*/
	var loopset sortableLoopSet
	err = createSortableLoopSet(arcset, &loopset)
	if err != eSuccess {
		destroyArcSet(&arcset)
		return err
	}

	// Extract polygons, since loops are contiguous in SortableLoopSet memory.
	// Polygons sorted by outer loop area, decreasing.
	err = createMultiPolygon(loopset, out)
	if err != eSuccess {
		destroySortableLoopSet(&loopset)
		destroyArcSet(&arcset)
		return err
	}

	destroyArcSet(&arcset)
	destroySortableLoopSetShallow(&loopset)

	return eSuccess
}
