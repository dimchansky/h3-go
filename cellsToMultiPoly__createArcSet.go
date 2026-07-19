package h3

// createArcSet builds the ArcSet for a cell set: per-cell CCW
// doubly-linked edge loops (each its own union-find component) plus the
// open-addressing hash buckets keyed by hashEdge with linear probing.
// Ported from H3 C: cellsToMultiPoly.c::createArcSet.
func createArcSet(cells []h3Index, numCells int64, arcset *arcSet) h3Error {
	numArcs := getNumEdges(cells, numCells)
	numBuckets := numArcs * hashTableMultiplier

	arcset.numArcs = numArcs
	arcset.numBuckets = numBuckets
	arcset.arcs = make([]arc, numArcs)

	arcset.buckets = make([]*arc, numBuckets)

	j := int64(0)
	for i := int64(0); i < numCells; i++ {
		var numEdges int64
		err := cellToEdgeArcs(cells[i], arcset.arcs[j:], &numEdges)
		if err != eSuccess {
			// NEVER in C.
			destroyArcSet(arcset)
			return err
		}
		j += numEdges
	}

	for i := int64(0); i < arcset.numArcs; i++ {
		// hash edge to initial bucket
		j := int64(hashEdge(arcset.arcs[i].id, uint64(arcset.numBuckets)))

		// linear probe to find next open bucket. wraps around if needed.
		for arcset.buckets[j] != nil {
			j = (j + 1) % arcset.numBuckets
		}
		arcset.buckets[j] = &arcset.arcs[i]
	}

	return eSuccess
}
