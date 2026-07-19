package h3

// unionArcs is part of the union-find data structure. It merges two
// arcs/edges into a single connected component (union by rank).
// Ported from H3 C: cellsToMultiPoly.c::unionArcs.
func unionArcs(a, b *arc) {
	a = getRoot(a)
	b = getRoot(b)

	if a.rank < b.rank {
		// swap so `a` has bigger rank
		a, b = b, a
	}

	if a != b {
		// `a` has bigger rank
		a.rank += b.rank
		b.parent = a
	}
}
