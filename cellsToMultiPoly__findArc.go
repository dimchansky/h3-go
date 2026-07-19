package h3

// findArc looks up the Arc for a directed edge in the ArcSet's hash
// buckets (hash + linear probe). Returns nil if the edge is not in the
// set.
// Ported from H3 C: cellsToMultiPoly.c::findArc.
func findArc(arcset arcSet, e h3Index) *arc {
	j := int64(hashEdge(e, uint64(arcset.numBuckets)))

	// hash + linear probe to find edge
	for arcset.buckets[j] != nil && arcset.buckets[j].id != e {
		j = (j + 1) % arcset.numBuckets
	}

	// returns NULL if edge not found
	return arcset.buckets[j]
}
