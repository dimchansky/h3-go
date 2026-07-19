package h3

// cancelArcPairs cancels out pairs of edges in the ArcSet, marking them
// as isRemoved. Updates the doubly-linked loop list to maintain valid
// loops. Merges the connected components of edge pairs; each connected
// component denotes a separate polygon (outer loop and holes).
// Ported from H3 C: cellsToMultiPoly.c::cancelArcPairs.
func cancelArcPairs(arcset arcSet) h3Error {
	for i := int64(0); i < arcset.numArcs; i++ {
		a := &arcset.arcs[i]

		if a.isRemoved {
			// Already removed, so we can skip.
			continue
		}

		var reversedEdge h3Index
		err := reverseDirectedEdge(a.id, &reversedEdge)
		if err != eSuccess {
			// NEVER in C.
			return err
		}

		b := findArc(arcset, reversedEdge)
		if b == nil {
			// The reversed edge was *not* in the set, so there's nothing to do.
			continue
		}

		// If we're at this point, then the two loops overlap at edges
		// `a` and `b`, which are opposites of each other.
		// Remove the two edges, and merge the loops to maintain
		// valid doubly-linked loops. Note that the two loops might be the
		// *same* loop, and the logic is the same either way.

		// mark both as removed
		a.isRemoved = true
		b.isRemoved = true

		// stitch together loops at removal site
		a.next.prev = b.prev
		a.prev.next = b.next
		b.next.prev = a.prev
		b.prev.next = a.next

		// update parent to merge into a single connected component
		unionArcs(a, b)
	}

	return eSuccess
}
